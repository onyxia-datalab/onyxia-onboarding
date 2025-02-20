package controller

import (
	"context"
	"errors"
	"testing"

	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/infrastructure/context"
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

// ✅ Updated MockUserContextReader to use GetUser instead of GetUsername, GetGroups, GetRoles
type MockUserContextReader struct {
	mock.Mock
}

func (m *MockUserContextReader) GetUser(ctx context.Context) (*domain.User, bool) {
	args := m.Called(ctx)
	user, _ := args.Get(0).(*domain.User)
	return user, args.Bool(1)
}

func (m *MockUserContextReader) GetUsername(ctx context.Context) (string, bool) {
	user, ok := m.GetUser(ctx)
	if !ok {
		return "", false
	}
	return user.Username, true
}

func (m *MockUserContextReader) GetGroups(ctx context.Context) ([]string, bool) {
	user, ok := m.GetUser(ctx)
	if !ok {
		return nil, false
	}
	return user.Groups, true
}

func (m *MockUserContextReader) GetRoles(ctx context.Context) ([]string, bool) {
	user, ok := m.GetUser(ctx)
	if !ok {
		return nil, false
	}
	return user.Roles, true
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

func TestOnboardingController_Onboard_Success_NoGroup(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)

	mockUsecase.On("Onboard", mock.Anything, mock.Anything).Return(nil)
	mockUserCtx.On("GetUser", mock.Anything).Return(&domain.User{
		Username: "test-user",
		Groups:   []string{"group1", "group2"},
		Roles:    []string{"role1"},
	}, true)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Set: false}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.NoError(t, err)
	assert.IsType(t, &api.OnboardOK{}, res)
	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}

func TestOnboardingController_Onboard_GetUserFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUserCtx.On("GetUser", mock.Anything).Return(nil, false) // ❌ GetUser fails

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)
	mockUsecase.AssertNotCalled(t, "Onboard")
}

func TestOnboardingController_Onboard_GroupValidationFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUserCtx.On("GetUser", mock.Anything).Return(&domain.User{
		Username: "test-user",
		Groups:   []string{"other-group"},
		Roles:    []string{"role1"},
	}, true)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardUnauthorized{}, res)
	mockUsecase.AssertNotCalled(t, "Onboard")
}

func TestOnboardingController_Onboard_OnboardingFails(t *testing.T) {
	mockUsecase := new(MockOnboardingUsecase)
	mockUserCtx := new(MockUserContextReader)
	mockUsecase.On("Onboard", mock.Anything, mock.Anything).
		Return(errors.New("onboarding service error"))
	mockUserCtx.On("GetUser", mock.Anything).Return(&domain.User{
		Username: "test-user",
		Groups:   []string{"test-group"},
		Roles:    []string{"role1"},
	}, true)

	controller := setupController(mockUsecase, mockUserCtx)
	req := api.OnboardingRequest{Group: api.OptString{Value: "test-group", Set: true}}

	res, err := controller.Onboard(context.Background(), &req)

	assert.Error(t, err)
	assert.IsType(t, &api.OnboardForbidden{}, res)

	mockUsecase.AssertCalled(t, "Onboard", mock.Anything, mock.Anything)
}
