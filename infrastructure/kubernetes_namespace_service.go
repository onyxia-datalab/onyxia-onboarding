package infrastructure

import (
	"context"
	"fmt"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const quotaName string = "onyxia-quota"

type KubernetesNamespaceService struct {
	clientset kubernetes.Interface
}

func NewKubernetesNamespaceService(clientset kubernetes.Interface) interfaces.NamespaceService {
	return &KubernetesNamespaceService{
		clientset: clientset,
	}
}

func (s *KubernetesNamespaceService) CreateNamespace(
	ctx context.Context,
	name string,
) (interfaces.NamespaceCreationResult, error) {
	namespacesClient := s.clientset.CoreV1().Namespaces()

	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}

	_, err := namespacesClient.Create(ctx, namespace, metav1.CreateOptions{})

	if errors.IsAlreadyExists(err) {
		return interfaces.NamespaceAlreadyExists, nil
	}

	if err != nil {
		return interfaces.NamespaceError, fmt.Errorf("failed to create namespace: %w", err)
	}

	return interfaces.NamespaceCreated, nil
}

func (s *KubernetesNamespaceService) ApplyResourceQuotas(
	ctx context.Context,
	namespace string,
	quota *domain.Quota,
) (interfaces.QuotaApplicationResult, error) {
	quotasClient := s.clientset.CoreV1().ResourceQuotas(namespace)

	hardLimits := convertQuotaToResourceMap(quota)

	// ✅ If no valid quotas exist, return early
	if len(hardLimits) == 0 {
		return interfaces.QuotaUnchanged, nil
	}

	resourceQuota := &v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      quotaName,
			Namespace: namespace,
			Labels: map[string]string{
				"createdby": "onyxia",
			},
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: hardLimits,
		}}

	existingQuota, err := quotasClient.Get(ctx, quotaName, metav1.GetOptions{})

	if err == nil {
		// Ignore quota if marked as ignored
		if _, ignore := existingQuota.Annotations["onyxia_ignore"]; ignore {
			return interfaces.QuotaIgnored, nil
		}

		// If quota is unchanged, return early
		if !quotasAreDifferent(existingQuota, resourceQuota) {
			return interfaces.QuotaUnchanged, nil
		}

		// Update existing quota
		existingQuota.Spec = resourceQuota.Spec
		_, updateErr := quotasClient.Update(ctx, existingQuota, metav1.UpdateOptions{})
		if updateErr != nil {
			return interfaces.QuotaError, fmt.Errorf(
				"failed to update resource quota: %w",
				updateErr,
			)
		}

		return interfaces.QuotaUpdated, nil
	}

	// If quota doesn't exist, create it
	if errors.IsNotFound(err) {
		_, err = quotasClient.Create(ctx, resourceQuota, metav1.CreateOptions{})
		if err != nil {
			return interfaces.QuotaError, fmt.Errorf("failed to create resource quota: %w", err)
		}
		return interfaces.QuotaCreated, nil
	}

	return interfaces.QuotaError, fmt.Errorf(
		"unexpected error checking for existing quota: %w",
		err,
	)
}

func quotasAreDifferent(existing, newQuota *v1.ResourceQuota) bool {
	if len(existing.Spec.Hard) != len(newQuota.Spec.Hard) {
		return true
	}

	for key, newValue := range newQuota.Spec.Hard {
		existingValue, exists := existing.Spec.Hard[key]
		if !exists || !existingValue.Equal(newValue) {
			return true
		}
	}

	return false
}

func convertQuotaToResourceMap(quota *domain.Quota) map[v1.ResourceName]resource.Quantity {
	quotaEntries := map[v1.ResourceName]string{
		v1.ResourceRequestsMemory:                  quota.MemoryRequest,
		v1.ResourceRequestsCPU:                     quota.CPURequest,
		v1.ResourceLimitsMemory:                    quota.MemoryLimit,
		v1.ResourceLimitsCPU:                       quota.CPULimit,
		v1.ResourceRequestsStorage:                 quota.StorageRequest,
		v1.ResourceName("count/pods"):              quota.MaxPods,
		v1.ResourceRequestsEphemeralStorage:        quota.EphemeralStorageRequest,
		v1.ResourceLimitsEphemeralStorage:          quota.EphemeralStorageLimit,
		v1.ResourceName("requests.nvidia.com/gpu"): quota.GPURequest,
		v1.ResourceName("limits.nvidia.com/gpu"):   quota.GPULimit,
	}

	// ✅ Filter out empty values and create a new immutable map
	hardLimits := make(map[v1.ResourceName]resource.Quantity, len(quotaEntries))
	for key, value := range quotaEntries {
		if value != "" {
			hardLimits[key] = resource.MustParse(value)
		}
	}

	return hardLimits
}
