package usecase

import (
	"context"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/infrastructure/context"
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
	annotations map[string]string,
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

var mockUserContextReader, _ = usercontext.NewMockUserContext(&domain.User{
	Username: testUserName,
	Groups:   []string{testGroupName},
	Roles:    []string{"role1"},
	Attributes: map[string]any{
		"attr1": "value1",
	},
})

func setupUsecase(
	mockService *MockNamespaceService,
	quotas domain.Quotas,
) domain.OnboardingUsecase {
	return NewOnboardingUsecase(
		mockService,
		domain.Namespace{
			NamespacePrefix:      namespacePrefix,
			GroupNamespacePrefix: groupNamespacePref,
			Annotation: domain.Annotation{
				Enabled: false,
				Static:  nil,
			},
		},
		quotas,
		mockUserContextReader,
	)
}

func setupPrivateUsecase(
	mockService *MockNamespaceService,
	quotas domain.Quotas,
) *onboardingUsecase {
	return &onboardingUsecase{
		namespaceService: mockService,
		namespace: domain.Namespace{
			NamespacePrefix:      namespacePrefix,
			GroupNamespacePrefix: groupNamespacePref,
			Annotation: domain.Annotation{
				Enabled: false,
				Static:  nil,
			},
		},
		quotas:            quotas,
		userContextReader: mockUserContextReader,
	}
}
