package usecase

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

type onboardingUsecase struct {
	namespaceService     interfaces.NamespaceService
	namespacePrefix      string
	groupNamespacePrefix string
	quotas               domain.Quotas
}

func NewOnboardingUsecase(
	namespaceService interfaces.NamespaceService,
	namespacePrefix, groupNamespacePrefix string,
	quotas domain.Quotas,
) *onboardingUsecase {
	return &onboardingUsecase{
		namespaceService:     namespaceService,
		namespacePrefix:      namespacePrefix,
		groupNamespacePrefix: groupNamespacePrefix,
		quotas:               quotas,
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
		return s.groupNamespacePrefix + *req.Group
	}
	return s.namespacePrefix + req.UserName
}
