package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/api/controller"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/infrastructure"
	"github.com/onyxia-datalab/onyxia-onboarding/usecase"
)

// SetupOnboardingRoutes initializes onboarding-related routes.
func SetupOnboardingRoutes(app *bootstrap.Application, router chi.Router, auth api.SecurityHandler) {

	namespaceCreator := infrastructure.NewKubernetesNamespaceService(app.K8sClient.Clientset)

	onboardingUsecase := usecase.NewOnboardingUsecase(namespaceCreator)
	controller := controller.NewOnboardingController(onboardingUsecase)

	// Create the ogen server with the handler
	srv, err := api.NewServer(controller, auth)
	if err != nil {
		panic(err)
	}

	router.Mount("/", http.HandlerFunc(srv.ServeHTTP))
}
