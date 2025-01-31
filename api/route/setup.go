package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func Setup(env *bootstrap.Env, r *chi.Mux) {

	r.Group(func(r chi.Router) {
		SetupOnboardingRoutes(env, r)
	})
}
