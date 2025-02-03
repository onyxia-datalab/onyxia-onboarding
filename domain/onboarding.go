package domain

import (
	"context"
)

type NamespaceService interface {
	CreateNamespace(ctx context.Context, name string) error
}

type OnboardingRequest struct {
	Group string
}

type OnboardingService interface {
	Onboard(ctx context.Context, req OnboardingRequest) error
}
