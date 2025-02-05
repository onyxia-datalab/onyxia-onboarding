package route

import (
	"context"
	"errors"

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

	getUserFromContext := func(ctx context.Context) (string, error) {
		user, ok := ctx.Value(middleware.UserContextKey).(string)
		if !ok || user == "" {
			return "", errors.New("user not found in context")
		}
		return user, nil
	}

	r.Group(func(r chi.Router) {
		SetupOnboardingRoutes(app, r, auth, getUserFromContext)
	})
}
