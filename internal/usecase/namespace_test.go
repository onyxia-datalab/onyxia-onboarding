package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/domain"
	usercontext "github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/context"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
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

func TestGetNamespaceAnnotations_Disabled(t *testing.T) {
	usecase := setupPrivateUsecase(new(MockNamespaceService), domain.Quotas{})
	usecase.namespace.Annotation.Enabled = false

	annotations := usecase.getNamespaceAnnotations(context.Background())

	assert.Nil(t, annotations, "Expected nil when annotations are disabled")
}

func TestGetNamespaceAnnotations_StaticOnly(t *testing.T) {
	usecase := setupPrivateUsecase(new(MockNamespaceService), domain.Quotas{})
	usecase.namespace.Annotation.Enabled = true
	usecase.namespace.Annotation.Static = map[string]string{
		"static-key": "static-value",
	}

	annotations := usecase.getNamespaceAnnotations(context.Background())

	assert.NotNil(t, annotations)
	assert.Equal(t, "static-value", annotations["static-key"])
}

func TestGetNamespaceAnnotations_LastLoginTimestamp(t *testing.T) {
	usecase := setupPrivateUsecase(new(MockNamespaceService), domain.Quotas{})
	usecase.namespace.Annotation.Enabled = true
	usecase.namespace.Annotation.Dynamic.LastLoginTimestamp = true

	annotations := usecase.getNamespaceAnnotations(context.Background())

	assert.NotNil(t, annotations)
	assert.Contains(t, annotations, "onyxia_last_login_timestamp")

	timestamp, err := time.ParseDuration(annotations["onyxia_last_login_timestamp"] + "ms")
	assert.NoError(t, err)
	assert.Greater(t, timestamp.Milliseconds(), int64(0))
}

func TestGetNamespaceAnnotations_UserAttributes(t *testing.T) {
	mockUserCtx, _ := usercontext.NewFakeUserContext(&domain.User{
		Attributes: map[string]any{
			"user-attr1": "value1",
			"user-attr2": "value2",
		},
	})

	usecase := setupPrivateUsecase(new(MockNamespaceService), domain.Quotas{})
	usecase.namespace.Annotation.Enabled = true
	usecase.namespace.Annotation.Dynamic.UserAttributes = []string{"user-attr1", "user-attr2"}
	usecase.userContextReader = mockUserCtx

	annotations := usecase.getNamespaceAnnotations(context.Background())

	assert.NotNil(t, annotations)
	assert.Equal(t, "value1", annotations["user-attr1"])
	assert.Equal(t, "value2", annotations["user-attr2"])
}

func TestGetNamespaceAnnotations_AllAnnotations(t *testing.T) {
	mockUserCtx, _ := usercontext.NewFakeUserContext(&domain.User{
		Attributes: map[string]any{
			"user-attr1": "value1",
		},
	})

	usecase := setupPrivateUsecase(new(MockNamespaceService), domain.Quotas{})
	usecase.namespace.Annotation.Enabled = true
	usecase.namespace.Annotation.Static = map[string]string{
		"static-key": "static-value",
	}
	usecase.namespace.Annotation.Dynamic.LastLoginTimestamp = true
	usecase.namespace.Annotation.Dynamic.UserAttributes = []string{"user-attr1"}
	usecase.userContextReader = mockUserCtx

	annotations := usecase.getNamespaceAnnotations(context.Background())

	assert.NotNil(t, annotations)
	assert.Equal(t, "static-value", annotations["static-key"])
	assert.Contains(t, annotations, "onyxia_last_login_timestamp")
	assert.Equal(t, "value1", annotations["user-attr1"])
}
