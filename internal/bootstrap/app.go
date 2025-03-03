package bootstrap

import (
	"fmt"
	"log/slog"

	usercontext "github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/context"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/infrastructure/kubernetes"
	"github.com/onyxia-datalab/onyxia-onboarding/internal/interfaces"
)

type Application struct {
	Env               *Env
	K8sClient         *kubernetes.KubernetesClient
	UserContextReader interfaces.UserContextReader
	UserContextWriter interfaces.UserContextWriter
}

func NewApplication() (*Application, error) {
	userReader, userWriter := usercontext.NewUserContext()

	InitLogger(userReader)

	env, err := NewEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load environment: %w", err)

	}

	k8sClient, err := kubernetes.NewKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes client: %w", err)
	}

	app := &Application{
		Env:               env,
		K8sClient:         k8sClient,
		UserContextReader: userReader,
		UserContextWriter: userWriter,
	}

	slog.Info("Application initialized successfully")

	return app, nil
}
