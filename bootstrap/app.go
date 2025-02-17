package bootstrap

import (
	"log/slog"
	"os"

	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"
)

type Application struct {
	Env       *Env
	K8sClient *kubernetes.KubernetesClient
}

func App() Application {
	InitLogger()
	app := &Application{}

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
