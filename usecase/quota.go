package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

func (s *onboardingUsecase) applyQuotas(
	ctx context.Context,
	namespace string,
	req domain.OnboardingRequest,
) error {
	if !s.quotas.Enabled {
		slog.Warn("⚠️ Quotas are disabled, skipping quota application",
			slog.String("namespace", namespace),
		)
		return nil
	}

	quotaToApply := s.getQuota(req, namespace)

	if quotaToApply == nil {
		slog.Warn("⚠️ No applicable quota found",
			slog.String("namespace", namespace),
		)
		return nil
	}

	result, err := s.namespaceService.ApplyResourceQuotas(ctx, namespace, quotaToApply)
	if err != nil {
		slog.Error("❌ Failed to apply quotas",
			slog.String("namespace", namespace),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to apply quotas to namespace (%s): %w", namespace, err)
	}

	switch result {
	case interfaces.QuotaCreated:
		slog.Info("✅ Created new resource quota",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaUpdated:
		slog.Info("✅ Updated resource quota",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaUnchanged:
		slog.Warn("⚠️ Resource quota is already up-to-date",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaIgnored:
		slog.Warn("⚠️ Quota ignored due to annotation",
			slog.String("namespace", namespace),
		)
	}

	return nil
}

func (s *onboardingUsecase) getQuota(req domain.OnboardingRequest, namespace string) *domain.Quota {
	switch {
	case req.Group != nil && s.quotas.GroupEnabled:
		slog.Info("🔹 Applying group quota",
			slog.String("namespace", namespace),
		)
		return &s.quotas.Group
	case s.quotas.UserEnabled:
		slog.Info("🔹 Applying user quota",
			slog.String("namespace", namespace),
		)
		return &s.quotas.User
	default:
		slog.Info("🔹 Applying default quota",
			slog.String("namespace", namespace),
		)
		return &s.quotas.Default
	}
}
