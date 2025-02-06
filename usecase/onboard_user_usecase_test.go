package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ✅ Ensure `MockNamespaceService` implements `NamespaceService`
var _ domain.NamespaceService = (*MockNamespaceService)(nil)

// ✅ Mock `NamespaceService` using Testify
type MockNamespaceService struct {
	mock.Mock
}

// ✅ Implement `CreateNamespace`
func (m *MockNamespaceService) CreateNamespace(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

// ✅ Setup function for `OnboardingUsecase`
func setupUsecase(mockService *MockNamespaceService) domain.OnboardingUsecase {
	return NewOnboardingUsecase(mockService, "user-", "group-")
}

// ✅ Test: Group is provided → Should create a group namespace
func TestOnboardingUsecase_Onboard_WithGroup(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupUsecase(mockService)

	groupName := "test-group"
	expectedNamespace := "group-" + groupName

	// ✅ Mock `CreateNamespace`
	mockService.On("CreateNamespace", mock.Anything, expectedNamespace).Return(nil)

	req := domain.OnboardingRequest{Group: &groupName, UserName: "test-user"}
	err := usecase.Onboard(context.Background(), req)

	assert.NoError(t, err)

	// ✅ Ensure `CreateNamespace` was called with the correct namespace
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, expectedNamespace)
}

// ✅ Test: Group is `nil` → Should create a user namespace
func TestOnboardingUsecase_Onboard_WithoutGroup(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupUsecase(mockService)

	expectedNamespace := "user-test-user"

	// ✅ Mock `CreateNamespace`
	mockService.On("CreateNamespace", mock.Anything, expectedNamespace).Return(nil)

	req := domain.OnboardingRequest{Group: nil, UserName: "test-user"}
	err := usecase.Onboard(context.Background(), req)

	assert.NoError(t, err)

	// ✅ Ensure `CreateNamespace` was called with the correct namespace
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, expectedNamespace)
}

// ❌ Test: `CreateNamespace` fails → Should return an error
func TestOnboardingUsecase_Onboard_CreateNamespaceFails(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupUsecase(mockService)

	groupName := "test-group"
	expectedNamespace := "group-" + groupName
	expectedError := errors.New("namespace creation failed")

	// ❌ Mock `CreateNamespace` failure
	mockService.On("CreateNamespace", mock.Anything, expectedNamespace).Return(expectedError)

	req := domain.OnboardingRequest{Group: &groupName, UserName: "test-user"}
	err := usecase.Onboard(context.Background(), req)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	// ✅ Ensure `CreateNamespace` was called
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, expectedNamespace)
}
