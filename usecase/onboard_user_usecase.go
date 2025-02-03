package usecase

import (
	"context"
	"errors"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type onboardingUsecase struct {
	namespaceService domain.NamespaceService
}

func NewOnboardingUsecase(namespaceService domain.NamespaceService) domain.OnboardingService {
	return &onboardingUsecase{namespaceService: namespaceService}
}

func (s *onboardingUsecase) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	if req.Group == "" {
		return errors.New("❌ Group name is required")
	}

	log.Printf("🚀 Onboarding user to group: %s", req.Group)

	if err := s.namespaceService.CreateNamespace(ctx, req.Group); err != nil {
		log.Printf("❌ Failed to create namespace: %v", err)
		return err
	}

	log.Println("✅ Onboarding completed successfully")
	return nil
}
