package onboarding

import (
	"context"
	"fmt"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/security"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type kubernetesOnboarder struct {
	namespaceCreator NamespaceCreator
	namespacePrefix  string
}

type NamespaceCreator interface {
	Create(ctx context.Context, ns *v1.Namespace, opts metav1.CreateOptions) (*v1.Namespace, error)
}

func NewKubernetesOnboarder(namespaceCreator NamespaceCreator, namespacePrefix string) *kubernetesOnboarder {
	return &kubernetesOnboarder{
		namespaceCreator: namespaceCreator,
		namespacePrefix:  namespacePrefix,
	}
}

var _ oas.Handler = (*kubernetesOnboarder)(nil)

func (s *kubernetesOnboarder) Onboard(ctx context.Context, req *oas.OnboardingRequest) (oas.OnboardRes, error) {
	user := ctx.Value(security.UserContextKey)
	if user == nil {
		return &oas.OnboardUnauthorized{}, nil
	}

	_, err := s.namespaceCreator.Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s%s", s.namespacePrefix, user.(string)),
		},
	}, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		err = nil
	}

	return &oas.OnboardOK{}, err
}
