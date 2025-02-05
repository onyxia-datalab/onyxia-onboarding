package controller

import (
	"context"
	"errors"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOnboardingUsecase struct {
	mock.Mock
}

var _ domain.OnboardingUsecase = (*MockOnboardingUsecase)(nil)

func (m *MockOnboardingUsecase) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func fakeGetUser(ctx context.Context) (string, error) {
	return "test-user", nil
}

func fakeGetUserFail(ctx context.Context) (string, error) {
	return "", errors.New("user retrieval failed")
}

func setupController(mockUsecase *MockOnboardingUsecase, getUser func(ctx context.Context) (string, error)) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: mockUsecase,
		getUser:           getUser,
	}
}

// ✅ Test: Onboarding succeeds when `Group` is set
func TestOnboardingController_Onboard_Success(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(nil) // ✅ Define behavior

	controller := setupController(mockUsecase, fakeGetUser)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)

	// ✅ Verify `Onboard` was called
	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}

// ✅ Test: Empty group should still succeed
func TestOnboardingController_Onboard_EmptyGroup(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(nil)

	controller := setupController(mockUsecase, fakeGetUser)
	req := api.OnboardingRequest{Group: api.OptString{Value: "", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)

	// ✅ Verify `Onboard` was called
	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}

// ❌ Test: `getUser` fails → Should return `OnboardForbidden`
func TestOnboardingController_Onboard_GetUserFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)

	controller := setupController(mockUsecase, fakeGetUserFail)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)

	mockUsecase.AssertNotCalled(t, "Onboard")
}

// ❌ Test: `OnboardingUsecase.Onboard` fails → Should return `OnboardForbidden`
func TestOnboardingController_Onboard_OnboardingFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(errors.New("onboarding service error"))

	controller := setupController(mockUsecase, fakeGetUser)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)

	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}
