package controller

import (
	"context"
	"log"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type OnboardingController struct {
	OnboardingUsecase domain.OnboardingService
}

func NewOnboardingController(onboardingUsecase domain.OnboardingService) *OnboardingController {
	return &OnboardingController{OnboardingUsecase: onboardingUsecase}
}

func (c *OnboardingController) Onboard(ctx context.Context, req *api.OnboardingRequest) (api.OnboardRes, error) {
	log.Println("Received Onboarding Request:", req.Group)

	err := c.OnboardingUsecase.Onboard(domain.OnboardingRequest{Group: req.Group})
	if err != nil {
		return &api.OnboardForbidden{}, err
	}

	return &api.OnboardOK{}, nil
}
