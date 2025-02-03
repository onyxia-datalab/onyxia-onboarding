package route

import (
	"context"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func Setup(app *bootstrap.Application, r *chi.Mux) {

	auth, err := middleware.OidcMiddleware(context.Background(),
		app.Env.AuthenticationMode,
		app.Env.OIDC.IssuerURI,
		app.Env.OIDC.ClientID,
		app.Env.OIDC.UsernameClaim,
	)

	if err != nil {
		panic(err)
	}

	r.Group(func(r chi.Router) {
		SetupOnboardingRoutes(app, r, auth)
	})
}
