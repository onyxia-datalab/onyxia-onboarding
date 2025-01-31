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

	message, err := c.OnboardingUsecase.Onboard(domain.OnboardingRequest{Group: req.Group})
	if err != nil {
		return nil, err
	}

	return struct {
		api.OnboardRes
		Message string `json:"message"`
	}{Message: message}, nil
}
