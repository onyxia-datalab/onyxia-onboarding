package bootstrap

import (
	"log/slog"
	"os"

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

func NewApplication() Application {
	userReader, userWriter := usercontext.NewUserContext()

	InitLogger(userReader)

	env, err := NewEnv()
	if err != nil {
		slog.Error("Failed to load environment", slog.Any("error", err))
		os.Exit(1)
	}

	k8sClient, err := kubernetes.NewKubernetesClient()
	if err != nil {
		slog.Error("Failed to initialize Kubernetes client", slog.Any("error", err))
		os.Exit(1)
	}

	app := &Application{
		Env:               env,
		K8sClient:         k8sClient,
		UserContextReader: userReader,
		UserContextWriter: userWriter,
	}

	slog.Info("Application initialized successfully")

	return *app
}
