package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	slogchi "github.com/samber/slog-chi"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/api/route"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func main() {

	app := bootstrap.App()
	env := app.Env

	// Get logger (slog.Default() is already set by bootsrap.App())
	logger := slog.Default()

	r := chi.NewRouter()

	//Logger middleware needs to be at top
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)

	r.Use(middleware.Heartbeat("/"))
	r.Use(cors.New(cors.Options{
		AllowedOrigins: env.Security.CORSAllowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"Origin",
			"X-Requested-With",
		},
		ExposedHeaders:   []string{"Link", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	if err := route.Setup(&app, r); err != nil {
		slog.Error("failed to set up routes: %v", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Server starting...", slog.String("address", ":8080"))

	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("failed to listen and serve",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
