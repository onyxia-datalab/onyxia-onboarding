package usecase

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

type onboardingUsecase struct {
	namespaceService  interfaces.NamespaceService
	namespace         domain.Namespace
	quotas            domain.Quotas
	userContextReader interfaces.UserContextReader
}

func NewOnboardingUsecase(
	namespaceService interfaces.NamespaceService,
	namespace domain.Namespace,
	quotas domain.Quotas,
	userContextReader interfaces.UserContextReader,

) *onboardingUsecase {
	return &onboardingUsecase{
		namespaceService:  namespaceService,
		namespace:         namespace,
		quotas:            quotas,
		userContextReader: userContextReader,
	}
}

func (s *onboardingUsecase) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	namespace := s.getNamespace(req)

	if err := s.createNamespace(ctx, namespace); err != nil {
		return err
	}

	if err := s.applyQuotas(ctx, namespace, req); err != nil {
		return err
	}

	return nil
}

func (s *onboardingUsecase) getNamespace(req domain.OnboardingRequest) string {
	if req.Group != nil {
		return s.namespace.GroupNamespacePrefix + *req.Group
	}
	return s.namespace.NamespacePrefix + req.UserName
}
