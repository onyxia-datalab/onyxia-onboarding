package usecase

import (
	"context"
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

func (s *onboardingUsecase) createNamespace(ctx context.Context, namespace string) error {
	result, err := s.namespaceService.CreateNamespace(ctx, namespace)

	if err != nil {
		slog.Error("❌ Failed to create namespace",
			slog.String("namespace", namespace),
			slog.Any("error", err),
		)
		return err
	}

	switch result {
	case interfaces.NamespaceCreated:
		slog.Info("✅ Successfully created namespace",
			slog.String("namespace", namespace),
		)
	case interfaces.NamespaceAlreadyExists:
		slog.Warn("⚠️ Namespace already exists",
			slog.String("namespace", namespace),
		)
	}

	return nil
}
