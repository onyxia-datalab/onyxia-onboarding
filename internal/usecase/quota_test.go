package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplyQuotas_Success(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	mockService.On("ApplyResourceQuotas", mock.Anything, userNamespace, &quotas.Default).
		Return(interfaces.QuotaCreated, nil)

	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.NoError(t, err)
	mockService.AssertCalled(
		t,
		"ApplyResourceQuotas",
		mock.Anything,
		userNamespace,
		&quotas.Default,
	)
}

func TestApplyQuotas_AlreadyUpToDate(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	mockService.On("ApplyResourceQuotas", mock.Anything, userNamespace, &quotas.Default).
		Return(interfaces.QuotaUnchanged, nil)

	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.NoError(t, err)
	mockService.AssertCalled(
		t,
		"ApplyResourceQuotas",
		mock.Anything,
		userNamespace,
		&quotas.Default,
	)
}

func TestApplyQuotas_QuotasDisabled(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{Enabled: false}
	usecase := setupPrivateUsecase(mockService, quotas)

	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.NoError(t, err)
	mockService.AssertNotCalled(t, "ApplyResourceQuotas")
}

func TestApplyQuotas_QuotaUpdated(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	mockService.On("ApplyResourceQuotas", mock.Anything, userNamespace, &quotas.Default).
		Return(interfaces.QuotaUpdated, nil)

	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.NoError(t, err)
	mockService.AssertCalled(
		t,
		"ApplyResourceQuotas",
		mock.Anything,
		userNamespace,
		&quotas.Default,
	)
}

func TestApplyQuotas_QuotaIgnored(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	mockService.On("ApplyResourceQuotas", mock.Anything, userNamespace, &quotas.Default).
		Return(interfaces.QuotaIgnored, nil)

	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.NoError(t, err)
	mockService.AssertCalled(
		t,
		"ApplyResourceQuotas",
		mock.Anything,
		userNamespace,
		&quotas.Default,
	)
}
func TestApplyQuotas_Failure(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	mockService.On("ApplyResourceQuotas", mock.Anything, userNamespace, &quotas.Default).
		Return(interfaces.QuotaApplicationResult(""), errors.New("failed to apply quotas"))
	err := usecase.applyQuotas(
		context.Background(),
		userNamespace,
		domain.OnboardingRequest{UserName: testUserName},
	)

	assert.Error(t, err)
	mockService.AssertCalled(
		t,
		"ApplyResourceQuotas",
		mock.Anything,
		userNamespace,
		&quotas.Default,
	)
}

func TestGetQuota_GroupQuota(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled:      true,
		GroupEnabled: true,
		Group:        domain.Quota{MemoryRequest: "12Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	groupName := testGroupName
	req := domain.OnboardingRequest{Group: &groupName, UserName: testUserName}

	quota := usecase.getQuota(context.Background(), req, groupNamespace)

	assert.Equal(t, &quotas.Group, quota)
}

func TestGetGroupQuota_FallbackToDefault(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled:      true,
		GroupEnabled: false,                               // ❌ Group quotas disabled
		Default:      domain.Quota{MemoryRequest: "10Gi"}, // ✅ Default quota exists
		Group:        domain.Quota{MemoryRequest: "20Gi"}, // 🚨 Should not be used
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	groupName := testGroupName
	req := domain.OnboardingRequest{UserName: testUserName, Group: &groupName}

	quota := usecase.getGroupQuota(context.Background(), req, userNamespace)

	// ✅ Expected: Fallback to `quotas.Default`
	assert.Equal(
		t,
		&quotas.Default,
		quota,
		"Expected fallback to default quota when group quotas are disabled",
	)
}

func TestGetQuota_UserQuota(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled:     true,
		UserEnabled: true,
		User:        domain.Quota{MemoryRequest: "11Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{Group: nil, UserName: testUserName}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	assert.Equal(t, &quotas.User, quota)
}

func TestGetQuota_DefaultQuota(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{Group: nil, UserName: testUserName}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	assert.Equal(t, &quotas.Default, quota)
}

func TestGetQuota_RoleQuota(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Roles: map[string]domain.Quota{
			"admin": {MemoryRequest: "16Gi"},
		},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{
		UserName:  testUserName,
		UserRoles: []string{"admin"}, // ✅ Only one role, should be used
	}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	expectedQuota := quotas.Roles["admin"]
	assert.Equal(t, &expectedQuota, quota, "Expected 'admin' role quota")
}

func TestGetQuota_RoleQuota_AppliesFirstMatch(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Roles: map[string]domain.Quota{
			"admin":     {MemoryRequest: "16Gi"},
			"developer": {MemoryRequest: "14Gi"},
		},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{
		UserName:  testUserName,
		UserRoles: []string{"developer", "admin"}, // ✅ "developer" should be used
	}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	expectedQuota := quotas.Roles["developer"] // ✅ Copy value before taking address
	assert.Equal(t, &expectedQuota, quota, "Expected the first matching role's quota")
}

func TestGetQuota_UserQuota_WhenNoRoleMatches(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled:     true,
		UserEnabled: true,
		User:        domain.Quota{MemoryRequest: "12Gi"},
		Roles: map[string]domain.Quota{
			"admin":     {MemoryRequest: "16Gi"},
			"developer": {MemoryRequest: "14Gi"},
		},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{
		UserName:  testUserName,
		UserRoles: []string{"nonexistent-role"}, // ❌ Role is not in the quota map
	}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	expectedQuota := quotas.User
	assert.Equal(t, &expectedQuota, quota, "Expected fallback to user quota when no role matches")
}

func TestGetQuota_DefaultQuota_WhenNoRoleAndUserQuotaDisabled(t *testing.T) {
	mockService := new(MockNamespaceService)
	quotas := domain.Quotas{
		Enabled: true,
		Default: domain.Quota{MemoryRequest: "10Gi"},
		User:    domain.Quota{MemoryRequest: "12Gi"},
	}
	usecase := setupPrivateUsecase(mockService, quotas)

	req := domain.OnboardingRequest{
		UserName:  testUserName,
		UserRoles: []string{}, // ✅ No roles provided
	}

	quota := usecase.getQuota(context.Background(), req, userNamespace)

	expectedQuota := quotas.Default
	assert.Equal(t, &expectedQuota, quota, "Expected default quota when no role/user quota applies")
}
