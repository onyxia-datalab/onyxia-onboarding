package kubernetes

import (
	"context"
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

	result, err := service.CreateNamespace(context.Background(), "test-namespace")

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceCreated, result) // ✅ Check the result type

	_, err = clientset.CoreV1().
		Namespaces().
		Get(context.Background(), "test-namespace", metav1.GetOptions{})
	assert.NoError(t, err) // ✅ Ensure namespace was actually created
}

// ✅ Test: Namespace Already Exists
func TestCreateNamespace_AlreadyExists(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
	})
	service := NewKubernetesNamespaceService(clientset)

	result, err := service.CreateNamespace(context.Background(), "test-namespace")

	assert.NoError(t, err)
	assert.Equal(t, interfaces.NamespaceAlreadyExists, result) // ✅ Check correct return value
}

// ❌ Test: Simulated API Failure
func TestCreateNamespace_Failure(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("create", "namespaces",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("simulated API failure")
		})

	result, err := service.CreateNamespace(context.Background(), "error-namespace")

	assert.Error(t, err)
	assert.Equal(
		t,
		interfaces.NamespaceCreationResult(""),
		result,
	) // ✅ Ensure correct error type is returned
	assert.Contains(t, err.Error(), "simulated API failure")
}

// ✅ Test: Apply Resource Quotas Successfully
func TestApplyResourceQuotas_Success(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: "test-namespace"},
	})
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi", CPURequest: "10"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaCreated, result) // ✅ Ensure quota was created

	createdQuota, err := clientset.CoreV1().
		ResourceQuotas("test-namespace").
		Get(context.Background(), QuotaName, metav1.GetOptions{})
	assert.NoError(t, err)
	assert.Equal(t, QuotaName, createdQuota.Name)
	assert.Equal(t, "onyxia", createdQuota.Labels["created-by"])

	memoryQuantity, exists := createdQuota.Spec.Hard[v1.ResourceRequestsMemory]
	assert.True(t, exists)
	assert.Equal(t, "10Gi", memoryQuantity.String())
}

// ✅ Test: Update Existing Quota (Different Values)
func TestApplyResourceQuotas_UpdateExistingQuota(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("8Gi"),
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi", CPURequest: "10"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaUpdated, result) // ✅ Ensure update is detected

	updatedQuota, err := clientset.CoreV1().
		ResourceQuotas("test-namespace").
		Get(context.Background(), QuotaName, metav1.GetOptions{})
	assert.NoError(t, err)

	memoryQuantity, exists := updatedQuota.Spec.Hard[v1.ResourceRequestsMemory]
	assert.True(t, exists)
	assert.Equal(t, "10Gi", memoryQuantity.String()) // ✅ Check new value
}

// ✅ Test: Quota Already Up-to-Date
func TestApplyResourceQuotas_NoUpdateNeeded(t *testing.T) {
	clientset := fake.NewSimpleClientset(&v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{Name: QuotaName, Namespace: "test-namespace"},
		Spec: v1.ResourceQuotaSpec{
			Hard: map[v1.ResourceName]resource.Quantity{
				v1.ResourceRequestsMemory: resource.MustParse("10Gi"),
				v1.ResourceRequestsCPU:    resource.MustParse("10"),
			},
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	quota := &domain.Quota{MemoryRequest: "10Gi", CPURequest: "10"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.NoError(t, err)
	assert.Equal(t, interfaces.QuotaUnchanged, result) // ✅ Ensure no update was needed
}

// ❌ Test: Quota Application Fails (namespace do not exists)
func TestApplyResourceQuotas_Failure(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("create", "resourcequotas",
		func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, errors.New("failed to create resource quota")
		})

	quota := &domain.Quota{MemoryRequest: "10Gi"}

	result, err := service.ApplyResourceQuotas(context.Background(), "test-namespace", quota)

	assert.Error(t, err)
	assert.Equal(
		t,
		interfaces.QuotaApplicationResult(""),
		result,
	) // ✅ Ensure correct error type is returned
	assert.Contains(t, err.Error(), "failed to create resource quota")
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
	assert.Equal(t, interfaces.QuotaIgnored, result) // ✅ Ensure quota was ignored
}
