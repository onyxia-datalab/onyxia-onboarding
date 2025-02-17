package route

import (
	"context"
	"fmt"
	"net/http"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/api/middleware"
	oas "github.com/onyxia-datalab/onyxia-onboarding/api/oas"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func Setup(app *bootstrap.Application, r *chi.Mux) error {

	auth, err := middleware.OidcMiddleware(context.Background(),
		app.Env.AuthenticationMode,
		middleware.OIDCConfig(app.Env.OIDC),
		app.UserContextWriter,
	)

	if err != nil {
		return fmt.Errorf("failed to initialize OIDC middleware: %w", err)
	}

	onboardingController := SetupOnboardingController(app)

	handler := &MyHandler{onboardImpl: onboardingController.Onboard}

	srv, err := oas.NewServer(
		handler,
		auth,
	)

	if err != nil {
		return fmt.Errorf("failed to create api server: %w", err)
	}

	r.Mount("/", http.HandlerFunc(srv.ServeHTTP))
	return nil
}
