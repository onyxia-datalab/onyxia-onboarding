package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

// ✅ Test: Create Namespace Successfully
func TestCreateNamespace_Success(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	result, err := service.CreateNamespace(context.Background(), "test-namespace", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceCreated, result)

	_, err = clientset.CoreV1().
		Namespaces().
		Get(context.Background(), "test-namespace", metav1.GetOptions{})
	assert.NoError(t, err)
}

// ✅ Test: Namespace Already Exists (No Annotation Change)
func TestCreateNamespace_AlreadyExists(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
	})
	service := NewKubernetesNamespaceService(clientset)

	result, err := service.CreateNamespace(context.Background(), "test-namespace", nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceAlreadyExists, result)
}

// ✅ Test: Namespace Already Exists (No Annotations Given)
func TestCreateNamespace_AlreadyExists_NoAnnotations(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
	})
	service := NewKubernetesNamespaceService(clientset)

	result, err := service.CreateNamespace(
		context.Background(),
		"test-namespace",
		map[string]string{},
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceAlreadyExists, result)
}

// ✅ Test: Update Annotations When Namespace Exists
func TestCreateNamespace_UpdateAnnotations(t *testing.T) {
	existingAnnotations := map[string]string{"old-key": "old-value"}
	newAnnotations := map[string]string{"new-key": "new-value"}

	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-namespace",
			Annotations: existingAnnotations,
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor(
		"patch",
		"namespaces",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			patchAction, ok := action.(k8stesting.PatchAction)
			assert.True(t, ok, "Expected PatchAction")

			var patch map[string]map[string]map[string]string
			err := json.Unmarshal(patchAction.GetPatch(), &patch)
			assert.NoError(t, err)
			assert.Equal(t, newAnnotations, patch["metadata"]["annotations"])

			return true, &v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test-namespace",
					Annotations: newAnnotations,
				},
			}, nil
		},
	)

	result, err := service.CreateNamespace(
		context.Background(),
		"test-namespace",
		newAnnotations,
		nil,
	)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceAnnotationsUpdated, result)
}

// ❌ Test: Simulated API Failure (Create)
func TestCreateNamespace_Failure(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("create", "namespaces",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("simulated API failure")
		})

	result, err := service.CreateNamespace(context.Background(), "error-namespace", nil, nil)

	assert.Error(t, err)
	assert.Equal(t, interfaces.NamespaceCreationResult(""), result)
	assert.Contains(t, err.Error(), "simulated API failure")
}

// ❌ Test: Simulated API Failure (Patch)
func TestCreateNamespace_FailurePatch(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
	})
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("patch", "namespaces",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("failed to patch annotations")
		})

	result, err := service.CreateNamespace(
		context.Background(),
		"test-namespace",
		map[string]string{"new-key": "new-value"}, nil,
	)

	assert.Error(t, err)
	assert.Equal(t, interfaces.NamespaceCreationResult(""), result)
	assert.Contains(t, err.Error(), "failed to patch annotations")
}

// ✅ Test: Apply Resource Quotas Successfully
func TestApplyResourceQuotas_Success(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi", CPURequest: "10"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaCreated, result)
}

// ✅ Test: Quota Already Exists with Unchanged Values
func TestApplyResourceQuotas_UnchangedQuota(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("10Gi"),
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi"} // Same values as existing

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaUnchanged, result)
}

// ✅ Test: Quota is Ignored Due to Annotation
func TestApplyResourceQuotas_IgnoredQuota(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      QuotaName,
			Namespace: "test-namespace",
			Annotations: map[string]string{
				IgnoreQuotaAnnotation: "true",
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaIgnored, result)
}

// ❌ Test: Failure When Checking for an Existing Quota
func TestApplyResourceQuotas_FailureCheck(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor(
		"get",
		"resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("failed to get quota")
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(t, interfaces.QuotaApplicationResult(""), result)
	assert.Contains(t, err.Error(), "failed to get quota")
}

// ❌ Test: Failure When Creating a Quota
func TestApplyResourceQuotas_FailureCreate(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor(
		"create",
		"resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("failed to create quota")
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(t, interfaces.QuotaApplicationResult(""), result)
	assert.Contains(t, err.Error(), "failed to create quota")
}

func TestApplyResourceQuotas_QuotaUpdated(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse(
					"5Gi",
				),
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor(
		"update",
		"resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, &v1.ResourceQuota{
				ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
				Spec: v1.ResourceQuotaSpec{
					Hard: map[v1.ResourceName]resource.Quantity{
						v1.ResourceRequestsMemory: resource.MustParse("10Gi"),
					},
				},
			}, nil
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"} // 👈 Updated quota

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaUpdated, result)
}

func TestApplyResourceQuotas_UnexpectedGetError(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor(
		"get",
		"resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("unexpected API error") // ✅ Force an unexpected error
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(t, interfaces.QuotaApplicationResult(""), result)
	assert.Contains(t, err.Error(), "unexpected error checking for existing quota")
}
func TestApplyResourceQuotas_FailureUpdate(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse(
					"5Gi",
				), // 👈 Existing quota is different
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	// Simulate an API failure when calling `Update()`
	clientset.PrependReactor(
		"update",
		"resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("failed to update resource quota")
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"} // 👈 Updated quota

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(t, interfaces.QuotaApplicationResult(""), result)
	assert.Contains(t, err.Error(), "failed to update resource quota")
}

func TestApplyResourceQuotas_EmptyQuota(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{} // 👈 Empty quota should return "QuotaUnchanged"

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaUnchanged, result)
}

func TestApplyResourceQuotas_FailureConvertQuota(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	// Simulate a quota that causes conversion failure
	quota := &domain.Quota{
		MemoryRequest: "invalid", // ❌ This will fail conversion
	}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(t, interfaces.QuotaApplicationResult(""), result)
	assert.Contains(t, err.Error(), "error converting quota to ResourceQuota")
}

func TestApplyResourceQuotas_LabelOnCreate(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	clientset.PrependReactor("create", "resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			createAction := action.(k8stesting.CreateAction)
			obj := createAction.GetObject().(*v1.ResourceQuota)

			labels := obj.GetLabels()
			assert.Equal(t, "onyxia", labels["created-by"], "Expected label 'created-by: onyxia'")

			return false, nil, nil // let the clientset handle creation
		},
	)

	_, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)
	assert.NoError(t, err)
}

func TestApplyResourceQuotas_LabelOnUpdate(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      QuotaName,
			Namespace: "test-namespace",
			Labels:    map[string]string{"created-by": "onyxia"},
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("5Gi"),
			},
		},
	})

	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("update", "resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			updateAction := action.(k8stesting.UpdateAction)
			obj := updateAction.GetObject().(*v1.ResourceQuota)

			labels := obj.GetLabels()
			assert.Equal(t, "onyxia", labels["created-by"], "Expected label 'created-by: onyxia'")

			return false, nil, nil // let clientset do the actual update
		},
	)

	quota := &domain.Quota{MemoryRequest: "10Gi"}
	_, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)
	assert.NoError(t, err)
}

func TestQuotasAreDifferent(t *testing.T) {
	existing := &v1.ResourceQuota{
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("5Gi"),
			},
		},
	}

	newQuota := &v1.ResourceQuota{
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("10Gi"), // 👈 Different value
			},
		},
	}

	result := quotasAreDifferent(existing, newQuota)
	assert.True(t, result, "Expected quotas to be different")
}

func TestQuotasAreDifferent_ExtraKeyInExistingQuota(t *testing.T) {
	existing := &v1.ResourceQuota{
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("10Gi"),
				v1.ResourceRequestsCPU: resource.MustParse(
					"5",
				), // 👈 Extra key not present in newQuota
			},
		},
	}

	newQuota := &v1.ResourceQuota{
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("10Gi"),
			},
		},
	}

	result := quotasAreDifferent(existing, newQuota)
	assert.True(t, result, "Expected quotas to be different due to missing key in new quota")
}

func TestConvertQuotaToResourceMap_InvalidQuantity(t *testing.T) {
	quota := domain.Quota{
		MemoryRequest: "invalid-quantity",
	}

	result, err := convertQuotaToResourceMap(quota)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "quantities must match")
}
