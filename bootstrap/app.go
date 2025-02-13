package bootstrap

import (
	"log/slog"

	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"
)

type Application struct {
	Env       *Env
	K8sClient *kubernetes.KubernetesClient
}

func App() Application {
	InitLogger()
	app := &Application{}

	app.Env = NewEnv()

	k8sClient := kubernetes.NewKubernetesClient()

	app.K8sClient = k8sClient

	slog.Info("Application initialized") // ✅ Logs will now work globally

	return *app
}
