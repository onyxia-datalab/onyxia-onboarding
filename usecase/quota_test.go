package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/onyxia-datalab/onyxia-onboarding/interfaces"
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
