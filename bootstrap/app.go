package bootstrap

import "github.com/onyxia-datalab/onyxia-onboarding/infrastructure/kubernetes"

type Application struct {
	Env       *Env
	K8sClient *kubernetes.KubernetesClient
}

func App() Application {
	app := &Application{}

	app.Env = NewEnv()

	k8sClient := kubernetes.NewKubernetesClient()

	app.K8sClient = k8sClient
	return *app
}
