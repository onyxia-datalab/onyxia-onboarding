package controller

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
)

type OnboardingController struct {
	OnboardingUsecase domain.OnboardingUsecase
	UserContextReader interfaces.UserContextReader
}

func NewOnboardingController(
	onboardingUsecase domain.OnboardingUsecase,
	userContextReader interfaces.UserContextReader,
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

	user, ok := c.UserContextReader.GetUser(ctx)
	if !ok || user == nil {
		err := fmt.Errorf("user not found in context")
		slog.Error("❌ Failed to retrieve user from context", slog.Any("error", err))
		return &api.OnboardForbidden{}, err
	}

	slog.InfoContext(ctx, "🔵 User identified")

	// Extract optional value from OptString
	var groupPtr *string
	if req.Group.Set { // Check if value is set
		groupPtr = &req.Group.Value

		// ✅ Check if the requested group is in user's groups
		if !slices.Contains(user.Groups, *groupPtr) {
			err := fmt.Errorf("user does not have access to group: %s", *groupPtr)
			slog.ErrorContext(ctx, "❌ Unauthorized group access",
				slog.String("group", *groupPtr),
				slog.Any("userGroups", user.Groups),
				slog.Any("error", err),
			)
			return &api.OnboardUnauthorized{}, err
		}
	}

	err := c.OnboardingUsecase.Onboard(ctx, domain.OnboardingRequest{
		Group:     groupPtr,
		UserName:  user.Username,
		UserRoles: user.Roles,
	})
	if err != nil {
		slog.ErrorContext(ctx, "❌ Onboarding failed",
			slog.Any("error", err),
		)
		return &api.OnboardForbidden{}, err
	}

	slog.InfoContext(ctx, "✅ Onboarding successful")
	return &api.OnboardOK{}, nil
}
