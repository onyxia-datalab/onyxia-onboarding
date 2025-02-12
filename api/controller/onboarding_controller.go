package controller

import (
	"context"
	"log/slog"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
)

type OnboardingController struct {
	OnboardingUsecase domain.OnboardingUsecase
	getUser           func(ctx context.Context) (string, error)
}

func NewOnboardingController(
	onboardingUsecase domain.OnboardingUsecase,
	getUser func(ctx context.Context) (string, error),
) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: onboardingUsecase,
		getUser:           getUser,
	}
}

func (c *OnboardingController) Onboard(
	ctx context.Context,
	req *api.OnboardingRequest,
) (api.OnboardRes, error) {
	slog.Info("🟢 Received Onboarding Request")

	userName, err := c.getUser(ctx)
	if err != nil {
		slog.Error("❌ Failed to retrieve user from context",
			slog.Any("error", err),
		)
		return &api.OnboardForbidden{}, err
	}
	slog.Info("🔵 User identified", slog.String("user", userName))

	// Extract optional value from OptString
	var groupPtr *string
	if req.Group.Set { // Check if value is set
		groupPtr = &req.Group.Value
		slog.Info("📌 Group provided", slog.String("group", req.Group.Value))
	}

	err = c.OnboardingUsecase.Onboard(
		ctx,
		domain.OnboardingRequest{Group: groupPtr, UserName: userName},
	)

	if err != nil {
		slog.Error("❌ Onboarding failed",
			slog.String("user", userName),
			slog.Any("group", groupPtr),
			slog.Any("error", err),
		)
		return &api.OnboardForbidden{}, err
	}

	slog.Info("✅ Onboarding successful",
		slog.String("user", userName),
		slog.Any("group", groupPtr),
	)
	return &api.OnboardOK{}, nil
}
