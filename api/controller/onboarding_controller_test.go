package controller

import (
	"context"
	"errors"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
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

type MockUserContextReader struct {
	mock.Mock
}

func (m *MockUserContextReader) GetUser(ctx context.Context) (string, bool) {
	args := m.Called(ctx)
	return args.String(0), args.Bool(1)
}

func (m *MockUserContextReader) GetGroups(ctx context.Context) ([]string, bool) {
	return nil, true
}

func (m *MockUserContextReader) GetRoles(ctx context.Context) ([]string, bool) {
	return nil, true
}

func setupController(
	mockUsecase *MockOnboardingUsecase,
	userCtx usercontext.UserContextReader,
) *OnboardingController {
	return &OnboardingController{
		OnboardingUsecase: mockUsecase,
		UserContextReader: userCtx,
	}
}
func TestOnboardingController_Onboard_Success(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(nil)
	mockUserCtx.On("GetUser", mock.Anything).Return("test-user", true)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)
	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}

func TestOnboardingController_Onboard_GetUserFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUserCtx.On("GetUser", mock.Anything).Return("", false)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)
	mockUsecase.AssertNotCalled(t, "Onboard")
}

func TestOnboardingController_Onboard_OnboardingFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).
		Return(errors.New("onboarding service error"))
	mockUserCtx.On("GetUser", mock.Anything).Return("test-user", true)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)

	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}
