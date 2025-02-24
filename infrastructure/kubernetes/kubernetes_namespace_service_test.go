package kubernetes

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
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

	result, err := service.CreateNamespace(context.Background(), "test-namespace", nil)

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

	result, err := service.CreateNamespace(context.Background(), "test-namespace", nil)

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

	result, err := service.CreateNamespace(context.Background(), "test-namespace", newAnnotations)

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

	result, err := service.CreateNamespace(context.Background(), "error-namespace", nil)

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
		map[string]string{"new-key": "new-value"},
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
