package domain

import (
	"context"
)

type OnboardingRequest struct {
	Group    *string // Use pointer to indicate optional value
	UserName string
}

type OnboardingUsecase interface {
	Onboard(ctx context.Context, req OnboardingRequest) error
}
