package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

func (s *onboardingUsecase) createNamespace(ctx context.Context, name string) error {
	result, err := s.namespaceService.CreateNamespace(
		ctx,
		name,
		s.getNamespaceAnnotations(ctx),
		s.namespace.NamespaceLabels,
	)

	slog.Info("result create Namespace", slog.String("result", string(result)))

	if err != nil {
		slog.ErrorContext(ctx, "❌ Failed to create namespace",
			slog.String("namespace", name),
			slog.Any("error", err),
		)
		return err
	}

	switch result {
	case interfaces.NamespaceCreated:
		slog.InfoContext(ctx, "✅ Successfully created namespace",
			slog.String("namespace", name),
		)
	case interfaces.NamespaceAlreadyExists:
		slog.WarnContext(ctx, "⚠️ Namespace already exists",
			slog.String("namespace", name),
		)
	}

	return nil
}

func (s *onboardingUsecase) getNamespaceAnnotations(
	ctx context.Context,
) map[string]string {
	if !s.namespace.Annotation.Enabled {
		return nil
	}

	annotations := s.namespace.Annotation.Static
	if annotations == nil {
		annotations = make(map[string]string)
	}

	if s.namespace.Annotation.Dynamic.LastLoginTimestamp {
		annotations["onyxia_last_login_timestamp"] = fmt.Sprint(time.Now().UnixMilli())
	}

	if attributes, ok := s.userContextReader.GetAttributes(ctx); ok {
		for _, attr := range s.namespace.Annotation.Dynamic.UserAttributes {
			annotations[attr] = fmt.Sprint(attributes[attr])
		}
	}
	return annotations
}
