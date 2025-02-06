package infrastructure

import (
	"bytes"
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

// ✅ Test: Create Namespace Successfully
func TestCreateNamespace_Success(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	err := service.CreateNamespace(context.Background(), "test-namespace")
	assert.NoError(t, err)

	_, err = clientset.CoreV1().Namespaces().Get(context.Background(), "test-namespace", metav1.GetOptions{})
	assert.NoError(t, err)
}

// ✅ Test: Namespace Already Exists
func TestCreateNamespace_AlreadyExists(t *testing.T) {
	// ✅ Capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	defer log.SetOutput(nil)

	clientset := fake.NewSimpleClientset(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: "existing-namespace",
		},
	})
	service := NewKubernetesNamespaceService(clientset)

	err := service.CreateNamespace(context.Background(), "existing-namespace")
	assert.NoError(t, err)
	logOutput := logBuffer.String()
	assert.Contains(t, logOutput, "⚠️ Namespace 'existing-namespace' already exists")
}

// ✅ Test: Simulated API Failure
func TestCreateNamespace_Failure(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	service := NewKubernetesNamespaceService(clientset)

	clientset.PrependReactor("create", "namespaces",
		func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
			return true, nil, errors.New("simulated API failure")
		})

	err := service.CreateNamespace(context.Background(), "error-namespace")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated API failure")
}
