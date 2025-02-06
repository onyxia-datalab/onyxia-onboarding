package domain

import (
	"context"
)

type NamespaceService interface {
	CreateNamespace(ctx context.Context, name string) error
}

type OnboardingRequest struct {
	Group    *string // Use pointer to indicate optional value
	UserName string
}

type OnboardingUsecase interface {
	Onboard(ctx context.Context, req OnboardingRequest) error
}
