package bootstrap

import (
	"log/slog"
	"os"

	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"
)

type Application struct {
	Env               *Env
	K8sClient         *kubernetes.KubernetesClient
	UserContextReader usercontext.UserContextReader
	UserContextWriter usercontext.UserContextWriter
}

func App() Application {
	app := &Application{}

	// Initialize User Context
	userReader, userWriter := usercontext.NewUserContext()
	app.UserContextReader = userReader
	app.UserContextWriter = userWriter

	InitLogger(userReader)

	env, err := NewEnv()
	if err != nil {
		slog.Error("Failed to load environment", slog.Any("error", err))
		os.Exit(1)
	}

	app.Env = env

	k8sClient, err := kubernetes.NewKubernetesClient()

	if err != nil {
		slog.Error("Failed to initialize Kubernetes client", slog.Any("error", err))
		os.Exit(1)
	}
	app.K8sClient = k8sClient

	slog.Info("Application initialized")

	return *app
}
