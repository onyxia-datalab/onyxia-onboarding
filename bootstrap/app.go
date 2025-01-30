package bootstrap

type Application struct {
	Env *Env
}

func App() Application {
	app := &Application{}

	app.Env = NewEnv()

	return *app
}
