package route

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/api/controller"
	api "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/onyxia-datalab/onyxia-onboarding/usecase"
)

// SetupOnboardingRoutes initializes onboarding-related routes.
func SetupOnboardingRoutes(env *bootstrap.Env, router chi.Router, auth api.SecurityHandler) {

	usecase := usecase.NewOnboardingUsecase()
	controller := controller.NewOnboardingController(usecase)

	// Create the ogen server with the handler
	srv, err := api.NewServer(controller, auth)
	if err != nil {
		panic(err)
	}

	router.Mount("/", http.HandlerFunc(srv.ServeHTTP))
}
