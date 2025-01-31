package usecase

import (
	"errors"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type onboardingUsecase struct {
}

func NewOnboardingUsecase() domain.OnboardingService {
	return &onboardingUsecase{}
}

func (s *onboardingUsecase) Onboard(req domain.OnboardingRequest) error {
	if req.Group == "" {
		return errors.New("group name is required")
	}

	log.Printf("Onboarding user to group: %s", req.Group)
	return nil
}
