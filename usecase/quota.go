package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

func (s *onboardingUsecase) applyQuotas(
	ctx context.Context,
	namespace string,
	req domain.OnboardingRequest,
) error {
	if !s.quotas.Enabled {
		log.Printf(
			"⚠️ Quotas are disabled, skipping quota application for namespace: %s",
			namespace,
		)
		return nil
	}

	quotaToApply := s.getQuota(req, namespace)

	if quotaToApply == nil {
		log.Printf("⚠️ No applicable quota found for namespace: %s", namespace)
		return nil
	}

	result, err := s.namespaceService.ApplyResourceQuotas(ctx, namespace, quotaToApply)
	if err != nil {
		return fmt.Errorf("failed to apply quotas to namespace (%s): %w", namespace, err)
	}

	switch result {
	case interfaces.QuotaCreated:
		log.Printf("✅ Created new resource quota for namespace: %s", namespace)
	case interfaces.QuotaUpdated:
		log.Printf("✅ Updated resource quota for namespace: %s", namespace)
	case interfaces.QuotaUnchanged:
		log.Printf("⚠️ Resource quota is already up-to-date for namespace: %s", namespace)
	case interfaces.QuotaIgnored:
		log.Printf("⚠️ Quota ignored due to annotation in namespace: %s", namespace)
	}

	return nil
}

func (s *onboardingUsecase) getQuota(req domain.OnboardingRequest, namespace string) *domain.Quota {
	switch {
	case req.Group != nil && s.quotas.GroupEnabled:
		log.Printf("🔹 Applying group quota for namespace: %s", namespace)
		return &s.quotas.Group
	case s.quotas.UserEnabled:
		log.Printf("🔹 Applying user quota for namespace: %s", namespace)
		return &s.quotas.User
	default:
		log.Printf("🔹 Applying default quota for namespace: %s", namespace)
		return &s.quotas.Default
	}
}
