package controller

import (
	"context"
	"log"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type OnboardingController struct {
	OnboardingUsecase domain.OnboardingUsecase
	getUser           func(ctx context.Context) (string, error)
}

func NewOnboardingController(onboardingUsecase domain.OnboardingUsecase, getUser func(ctx context.Context) (string, error)) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: onboardingUsecase,
		getUser:           getUser,
	}
}

func (c *OnboardingController) Onboard(ctx context.Context, req *api.OnboardingRequest) (api.OnboardRes, error) {
	log.Printf("🟢 Received Onboarding Request")

	userName, err := c.getUser(ctx)
	if err != nil {
		log.Printf("❌ Failed to retrieve user from context: %v", err)
		return &api.OnboardForbidden{}, err
	}
	log.Printf("🔵 User identified: %s", userName)

	// Extract optional value from OptString
	var groupPtr *string
	if req.Group.Set { // Check if value is set
		groupPtr = &req.Group.Value
		log.Printf("📌 Group provided: %s", req.Group.Value)
	}

	err = c.OnboardingUsecase.Onboard(ctx, domain.OnboardingRequest{Group: groupPtr, UserName: userName})

	if err != nil {
		log.Printf("❌ Onboarding failed | User: %s | Group: %v | Error: %v", userName, groupPtr, err)
		return &api.OnboardForbidden{}, err
	}

	log.Printf("✅ Onboarding successful | User: %s | Group: %v", userName, groupPtr)
	return &api.OnboardOK{}, nil
}
