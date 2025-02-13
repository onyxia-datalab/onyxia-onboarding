package route

import (
	"context"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/api/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/domain/usercontext"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func Setup(app *bootstrap.Application, r *chi.Mux) {

	userContextReader, userContextWriter := usercontext.NewUserContext()

	auth, err := middleware.OidcMiddleware(context.Background(),
		app.Env.AuthenticationMode,
		middleware.OIDCConfig(app.Env.OIDC),
		userContextWriter,
	)

	if err != nil {
		panic(err)
	}

	r.Group(func(r chi.Router) {
		SetupOnboardingRoutes(app, r, auth, userContextReader)
	})
}
