package usecase

import (
	"context"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type onboardingUsecase struct {
	namespaceService     domain.NamespaceService
	namespacePrefix      string
	groupNamespacePrefix string
}

func NewOnboardingUsecase(namespaceService domain.NamespaceService, namespacePrefix string,
	groupNamespacePrefix string) domain.OnboardingUsecase {
	return &onboardingUsecase{namespaceService: namespaceService,
		namespacePrefix: namespacePrefix, groupNamespacePrefix: groupNamespacePrefix}
}

func (s *onboardingUsecase) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	var namespace string

	if req.Group != nil {
		namespace = s.groupNamespacePrefix + *req.Group // Group namespace
	} else {
		namespace = s.namespacePrefix + req.UserName
	}

	log.Printf("🚀 Creating namespace: %s", namespace)

	if err := s.namespaceService.CreateNamespace(ctx, namespace); err != nil {
		log.Printf("❌ Failed to create namespace (%s): %v", namespace, err)
		return err
	}

	log.Printf("✅ Successfully created namespace: %s", namespace)
	return nil

}
