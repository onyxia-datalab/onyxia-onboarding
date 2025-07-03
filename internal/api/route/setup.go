package route

import (
	"context"
	"fmt"
	"net/http"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/internal/api/middleware"
	oas "github.com/onyxia-datalab/onyxia-onboarding/internal/api/oas"

	"github.com/onyxia-datalab/onyxia-onboarding/internal/bootstrap"
)

func Setup(app *bootstrap.Application) (http.Handler, error) {

	auth, err := middleware.OidcMiddleware(context.Background(),
		app.Env.AuthenticationMode,
		middleware.OIDCConfig(app.Env.OIDC),
		app.UserContextWriter,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC middleware: %w", err)
	}

	onboardingController := SetupOnboardingController(app)

	handler := &MyHandler{onboardImpl: onboardingController.Onboard}

	srv, err := oas.NewServer(
		handler,
		auth,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create api server: %w", err)
	}

	return srv, nil
}
