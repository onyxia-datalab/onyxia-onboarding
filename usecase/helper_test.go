package usecase

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
	"github.com/stretchr/testify/mock"
)

// ✅ Shared Test Constants
const (
	testUserName       = "test-user"
	testGroupName      = "test-group"
	defaultNamespace   = "user-test-user"
	userNamespace      = "user-test-user"
	groupNamespace     = "projet-test-group"
	namespacePrefix    = "user-"
	groupNamespacePref = "projet-"
)

// ✅ Mock `NamespaceService`
type MockNamespaceService struct {
	mock.Mock
}

var _ interfaces.NamespaceService = (*MockNamespaceService)(nil)

func (m *MockNamespaceService) CreateNamespace(
	ctx context.Context,
	name string,
) (interfaces.NamespaceCreationResult, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(interfaces.NamespaceCreationResult), args.Error(1)
}

func (m *MockNamespaceService) ApplyResourceQuotas(
	ctx context.Context,
	namespace string,
	quota *domain.Quota,
) (interfaces.QuotaApplicationResult, error) {
	args := m.Called(ctx, namespace, quota)
	return args.Get(0).(interfaces.QuotaApplicationResult), args.Error(1)
}

func setupUsecase(
	mockService *MockNamespaceService,
	quotas domain.Quotas,
) domain.OnboardingUsecase {
	return NewOnboardingUsecase(mockService, namespacePrefix, groupNamespacePref, quotas)
}

func setupPrivateUsecase(
	mockService *MockNamespaceService,
	quotas domain.Quotas,
) *onboardingUsecase {
	return &onboardingUsecase{
		namespaceService:     mockService,
		namespacePrefix:      namespacePrefix,
		groupNamespacePrefix: groupNamespacePref,
		quotas:               quotas,
	}
}
