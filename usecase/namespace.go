package usecase

import (
	"context"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

func (s *onboardingUsecase) createNamespace(ctx context.Context, namespace string) error {
	result, err := s.namespaceService.CreateNamespace(ctx, namespace)

	if err != nil {
		log.Printf("❌ Failed to create namespace (%s): %v", namespace, err)
		return err
	}

	switch result {
	case interfaces.NamespaceCreated:
		log.Printf("✅ Successfully created namespace: %s", namespace)
	case interfaces.NamespaceAlreadyExists:
		log.Printf("⚠️ Namespace already exists: %s", namespace)
	}

	return nil
}
