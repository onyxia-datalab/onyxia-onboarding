package controller

import (
	"context"
	"fmt"
	"log/slog"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
)

type OnboardingController struct {
	OnboardingUsecase domain.OnboardingUsecase
	UserContextReader usercontext.UserContextReader
}

func NewOnboardingController(
	onboardingUsecase domain.OnboardingUsecase,
	userContextReader usercontext.UserContextReader,
) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: onboardingUsecase,
		UserContextReader: userContextReader,
	}
}

func (c *OnboardingController) Onboard(
	ctx context.Context,
	req *api.OnboardingRequest,
) (api.OnboardRes, error) {
	slog.Info("🟢 Received Onboarding Request")

	userName, ok := c.UserContextReader.GetUser(ctx)
	if !ok {
		err := fmt.Errorf("user not found in context")
		slog.Error("❌ Failed to retrieve user from context", slog.Any("error", err))
		return &api.OnboardForbidden{}, err
	}

	slog.InfoContext(ctx, "🔵 User identified")

	// Extract optional value from OptString
	var groupPtr *string
	if req.Group.Set { // Check if value is set
		groupPtr = &req.Group.Value
	}

	err := c.OnboardingUsecase.Onboard(
		ctx,
		domain.OnboardingRequest{Group: groupPtr, UserName: userName},
	)
	if err != nil {
		slog.ErrorContext(ctx, "❌ Onboarding failed",
			slog.Any("error", err),
		)
		return &api.OnboardForbidden{}, err
	}

	slog.InfoContext(ctx, "✅ Onboarding successful")
	return &api.OnboardOK{}, nil
}
