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
