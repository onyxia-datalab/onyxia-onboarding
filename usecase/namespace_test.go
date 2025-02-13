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

func TestCreateNamespace_Success(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupPrivateUsecase(mockService, domain.Quotas{})

	mockService.On("CreateNamespace", mock.Anything, userNamespace).
		Return(interfaces.NamespaceCreated, nil)

	err := usecase.createNamespace(context.Background(), userNamespace)

	assert.NoError(t, err)
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, userNamespace)
}

func TestCreateNamespace_AlreadyExists(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupPrivateUsecase(mockService, domain.Quotas{})

	mockService.On("CreateNamespace", mock.Anything, userNamespace).
		Return(interfaces.NamespaceAlreadyExists, nil)

	err := usecase.createNamespace(context.Background(), userNamespace)

	assert.NoError(t, err)
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, userNamespace)
}

func TestCreateNamespace_Failure(t *testing.T) {
	mockService := new(MockNamespaceService)
	usecase := setupPrivateUsecase(mockService, domain.Quotas{})

	mockService.On("CreateNamespace", mock.Anything, userNamespace).
		Return(interfaces.NamespaceCreationResult(""), errors.New("failed to create namespace"))
	err := usecase.createNamespace(context.Background(), userNamespace)

	assert.Error(t, err)
	mockService.AssertCalled(t, "CreateNamespace", mock.Anything, userNamespace)
}
