package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	"github.com/stretchr/testify/assert"
)

// ✅ Fake Namespace Service (Mocks `CreateNamespace`)
type FakeNamespaceService struct {
	ShouldFail bool
}

func (f *FakeNamespaceService) CreateNamespace(ctx context.Context, name string) error {
	if f.ShouldFail {
		return errors.New("failed to create namespace")
	}
	return nil
}

// ✅ Test: Group name is required
func TestOnboard_EmptyGroupName(t *testing.T) {
	service := NewOnboardingUsecase(&FakeNamespaceService{})

	err := service.Onboard(context.Background(), domain.OnboardingRequest{Group: ""})
	assert.Error(t, err)
	assert.Equal(t, "❌ Group name is required", err.Error())
}

// ✅ Test: Successful onboarding
func TestOnboard_Success(t *testing.T) {
	service := NewOnboardingUsecase(&FakeNamespaceService{ShouldFail: false})

	err := service.Onboard(context.Background(), domain.OnboardingRequest{Group: "test-group"})
	assert.NoError(t, err)
}

// ✅ Test: Namespace creation failure
func TestOnboard_FailedNamespaceCreation(t *testing.T) {
	service := NewOnboardingUsecase(&FakeNamespaceService{ShouldFail: true})

	err := service.Onboard(context.Background(), domain.OnboardingRequest{Group: "test-group"})
	assert.Error(t, err)
	assert.Equal(t, "failed to create namespace", err.Error())

}
