package controller

import (
	"context"
	"errors"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/internal/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/context"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ✅ Mock `OnboardingUsecase`
type MockOnboardingUsecase struct {
	mock.Mock
}

var _ domain.OnboardingUsecase = (*MockOnboardingUsecase)(nil)

func (m *MockOnboardingUsecase) Onboard(ctx context.Context, req domain.OnboardingRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// ✅ Test Setup Function
func setupController(
	mockUsecase *MockOnboardingUsecase,
	mockUserCtx interfaces.UserContextReader,
) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: mockUsecase,
		UserContextReader: mockUserCtx,
	}
}

// ✅ Test Cases

func TestOnboardingController_Onboard_Success_NoGroup(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx, _ := usercontext.NewMockUserContext(&domain.User{
		Username: "test-user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	})

	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(nil)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Set: false}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)
	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}

func TestOnboardingController_Onboard_GetUserFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx, _ := usercontext.NewMockUserContext(nil) // ❌ GetUser fails

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)
	mockUsecase.AssertNotCalled(t, "Onboard")
}

func TestOnboardingController_Onboard_GroupValidationFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx, _ := usercontext.NewMockUserContext(&domain.User{
		Username: "test-user",
		Groups:   []string{"other-group"}, // ❌ Does not match "test-group"
		Roles:    []string{"role1"},
	})

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardUnauthorized{}, res)
	mockUsecase.AssertNotCalled(t, "Onboard")
}

func TestOnboardingController_Onboard_OnboardingFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx, _ := usercontext.NewMockUserContext(&domain.User{
		Username: "test-user",
		Groups:   []string{"test-group"},
		Roles:    []string{"role1"},
	})

	mockUsecase.On("Onboard", mock.Anything, mock.Anything).
		Return(errors.New("onboarding service error"))

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)

	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}
