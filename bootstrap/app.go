package bootstrap

type Application struct {
	Env       *Env
	K8sClient *KubernetesClient
}

func App() Application {
	app := &Application{}

	app.Env = NewEnv()

	k8sClient := NewKubernetesClient()

	app.K8sClient = k8sClient
	return *app
}
