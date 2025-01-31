package usecase

import (
	"errors"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type onboardingUsecase struct {
}

func NewOnboardingUsecase() domain.OnboardingService {
	return &onboardingUsecase{}
}

func (s *onboardingUsecase) Onboard(req domain.OnboardingRequest) (string, error) {
	if req.Group == "" {
		return "", errors.New("group name is required")
	}
	return "User onboarded successfully", nil
}
