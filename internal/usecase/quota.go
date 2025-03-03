package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

func (s *onboardingUsecase) applyQuotas(
	ctx context.Context,
	namespace string,
	req domain.OnboardingRequest,
) error {
	if !s.quotas.Enabled {
		slog.WarnContext(ctx, "⚠️ Quotas are disabled, skipping quota application",
			slog.String("namespace", namespace),
		)
		return nil
	}

	quotaToApply := s.getQuota(ctx, req, namespace)

	result, err := s.namespaceService.ApplyResourceQuotas(ctx, namespace, quotaToApply)
	if err != nil {
		slog.ErrorContext(ctx, "❌ Failed to apply quotas",
			slog.String("namespace", namespace),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to apply quotas to namespace (%s): %w", namespace, err)
	}

	switch result {
	case interfaces.QuotaCreated:
		slog.InfoContext(ctx, "✅ Created new resource quota",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaUpdated:
		slog.InfoContext(ctx, "✅ Updated resource quota",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaUnchanged:
		slog.WarnContext(ctx, "⚠️ Resource quota is already up-to-date",
			slog.String("namespace", namespace),
		)
	case interfaces.QuotaIgnored:
		slog.WarnContext(ctx, "⚠️ Quota ignored due to annotation",
			slog.String("namespace", namespace),
		)
	}

	return nil
}

func (s *onboardingUsecase) getQuota(
	ctx context.Context,
	req domain.OnboardingRequest,
	namespace string,
) *domain.Quota {
	// ✅ If a group is set, check if group quotas are enabled
	if req.Group != nil {
		return s.getGroupQuota(ctx, req, namespace)
	}

	// ✅ Otherwise, apply user quota (roles first)
	return s.getUserQuota(ctx, req, namespace)
}

func (s *onboardingUsecase) getGroupQuota(
	ctx context.Context,
	req domain.OnboardingRequest,
	namespace string,
) *domain.Quota {
	if s.quotas.GroupEnabled {
		slog.InfoContext(ctx, "🔹 Applying group quota",
			slog.String("namespace", namespace),
			slog.String("group", *req.Group),
		)
		return &s.quotas.Group
	}
	return &s.quotas.Default
}

func (s *onboardingUsecase) getUserQuota(
	ctx context.Context,
	req domain.OnboardingRequest,
	namespace string,
) *domain.Quota {
	for _, role := range req.UserRoles {
		if quota, exists := s.quotas.Roles[role]; exists {
			slog.InfoContext(ctx, "🔹 Applying role-based user quota",
				slog.String("namespace", namespace),
				slog.String("role", role),
			)
			return &quota
		}
	}

	// ✅ Fallback to user quota (if enabled)
	if s.quotas.UserEnabled {
		slog.InfoContext(ctx, "🔹 Applying user quota",
			slog.String("namespace", namespace),
		)
		return &s.quotas.User
	}

	// ✅ Fallback to default quota
	slog.InfoContext(ctx, "🔹 Applying default quota",
		slog.String("namespace", namespace),
	)
	return &s.quotas.Default
}
