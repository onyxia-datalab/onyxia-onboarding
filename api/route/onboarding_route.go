package route

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/api/controller"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure"
	"github.com/onyxia-datalab/onyxia-onboarding/usecase"
)

func SetupOnboardingRoutes(app *bootstrap.Application, router chi.Router, auth api.SecurityHandler, getUser func(ctx context.Context) (string, error)) {

	namespaceCreator := infrastructure.NewKubernetesNamespaceService(app.K8sClient.Clientset)

	onboardingUsecase := usecase.NewOnboardingUsecase(namespaceCreator, app.Env.Service.NamespacePrefix, app.Env.Service.GroupNamespacePrefix)
	controller := controller.NewOnboardingController(onboardingUsecase, getUser)

	srv, err := api.NewServer(controller, auth)
	if err != nil {
		panic(err)
	}

	router.Mount("/", http.HandlerFunc(srv.ServeHTTP))
}
