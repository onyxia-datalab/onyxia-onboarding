package kubernetes

import (
	"context"
	"fmt"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

const QuotaName string = "onyxia-quota"
const IgnoreQuotaAnnotation string = "onyxia.sh/ignore"

type KubernetesNamespaceService struct {
	clientset k8s.Interface
}

func NewKubernetesNamespaceService(clientset k8s.Interface) interfaces.NamespaceService {
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
		return "", fmt.Errorf("failed to create namespace: %w", err)
	}

	return interfaces.NamespaceCreated, nil
}

func (s *KubernetesNamespaceService) ApplyResourceQuotas(
	ctx context.Context,
	namespace string,
	quota *domain.Quota,
) (interfaces.QuotaApplicationResult, error) {
	quotasClient := s.clientset.CoreV1().ResourceQuotas(namespace)

	hardLimits, err := convertQuotaToResourceMap(*quota)

	if err != nil {
		return "", fmt.Errorf("error converting quota to ResourceQuota: %w", err)
	}

	// ✅ If no valid quotas exist, return early
	if len(hardLimits) == 0 {
		return interfaces.QuotaUnchanged, nil
	}

	resourceQuota := &v1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      QuotaName,
			Namespace: namespace,
			Labels: map[string]string{
				"created-by": "onyxia",
			},
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: hardLimits,
		}}

	existingQuota, err := quotasClient.Get(ctx, QuotaName, metav1.GetOptions{})

	if err == nil {
		// Ignore quota if marked as ignored
		if ignore, ok := existingQuota.Annotations[IgnoreQuotaAnnotation]; ok && ignore == "true" {
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
			return "", fmt.Errorf(
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
			return "", fmt.Errorf("failed to create resource quota: %w", err)
		}
		return interfaces.QuotaCreated, nil
	}

	return "", fmt.Errorf(
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

func convertQuotaToResourceMap(quota domain.Quota) (map[v1.ResourceName]resource.Quantity, error) {
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
			quantity, err := resource.ParseQuantity(value)
			if err != nil {
				return nil, err
			}
			hardLimits[key] = quantity
		}
	}

	return hardLimits, nil
}
