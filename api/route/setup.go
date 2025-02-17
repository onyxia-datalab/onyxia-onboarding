package route

import (
	"context"
	"net/http"

	middleware "github.com/onyxia-datalab/onyxia-onboarding/api/middleware"
	oas "github.com/onyxia-datalab/onyxia-onboarding/api/oas"
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

	onboardingRoute := NewOnboardingRoute(app, userContextReader)

	handler := &MyHandler{onboardImpl: onboardingRoute}

	srv, err := oas.NewServer(
		handler,
		auth,
	)

	if err != nil {
		panic(err)
	}

	r.Mount("/", http.HandlerFunc(srv.ServeHTTP))
}
