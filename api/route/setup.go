package route

import (
	"context"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func Setup(env *bootstrap.Env, r *chi.Mux) {

	auth, err := middleware.OidcMiddleware(context.Background(),
		env.OIDC.IssuerURI,
		env.OIDC.ClientID,
		env.OIDC.UsernameClaim,
	)

	if err != nil {
		panic(err)
	}

	r.Group(func(r chi.Router) {
		SetupOnboardingRoutes(env, r, auth)
	})
}
