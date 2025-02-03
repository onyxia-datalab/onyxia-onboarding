package controller

import (
	"context"
	"errors"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
)

// ✅ Fake implementation of OnboardingService to simulate real behavior
type FakeOnboardingService struct {
	OnboardFunc func(ctx context.Context, req domain.OnboardingRequest) error
}

func (f *FakeOnboardingService) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	return f.OnboardFunc(ctx, req)
}

// ✅ Test: Successful onboarding
func TestOnboardingController_Onboard_Success(t *testing.T) {
	service := &FakeOnboardingService{
		OnboardFunc: func(ctx context.Context, req domain.OnboardingRequest) error {
			if req.Group == "" {
				return errors.New("group cannot be empty")
			}
			return nil
		},
	}

	controller := NewOnboardingController(service)
	req := api.OnboardingRequest{Group: "test-group"}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)
}

// ✅ Test: Onboarding failure (forbidden)
func TestOnboardingController_Onboard_Forbidden(t *testing.T) {
	service := &FakeOnboardingService{
		OnboardFunc: func(ctx context.Context, req domain.OnboardingRequest) error {
			return errors.New("forbidden")
		},
	}

	controller := NewOnboardingController(service)
	req := api.OnboardingRequest{Group: "test-group"}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)
}

// ✅ Test: Empty group should return an error
func TestOnboardingController_Onboard_EmptyGroup(t *testing.T) {
	service := &FakeOnboardingService{
		OnboardFunc: func(ctx context.Context, req domain.OnboardingRequest) error {
			if req.Group == "" {
				return errors.New("group name is required")
			}
			return nil
		},
	}

	controller := NewOnboardingController(service)
	req := api.OnboardingRequest{Group: ""}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)
}
