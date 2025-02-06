package infrastructure

import (
	"context"
	"fmt"
	"log"

	"github.com/onyxia-datalab/onyxia-onboarding/domain"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubernetesNamespaceService struct {
	clientset kubernetes.Interface
}

func NewKubernetesNamespaceService(clientset kubernetes.Interface) domain.NamespaceService {
	return &KubernetesNamespaceService{
		clientset: clientset,
	}
}

func (s *KubernetesNamespaceService) CreateNamespace(ctx context.Context, name string) error {
	namespacesClient := s.clientset.CoreV1().Namespaces()

	// Attempt to create the namespace
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := namespacesClient.Create(ctx, namespace, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		log.Printf("⚠️ Namespace '%s' already exists, continuing...", name)
		return nil
	}

	if err != nil {
		return fmt.Errorf("❌ Failed to create namespace: %w", err)
	}

	log.Printf("✅ Namespace '%s' created successfully", name)
	return nil
}
