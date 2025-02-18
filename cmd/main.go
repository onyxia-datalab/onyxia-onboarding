package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/api/route"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func main() {

	app := bootstrap.App()
	env := app.Env

	r := chi.NewRouter()

	r.Use(
		httplog.RequestLogger(
			&httplog.Logger{
				Logger:  slog.Default(),
				Options: httplog.Options{Concise: true, RequestHeaders: true},
			},
		),
	)

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

	address := fmt.Sprintf(":%d", env.Server.Port)

	slog.Info("Server starting...", slog.String("address", address))

	if err := http.ListenAndServe(address, r); err != nil {
		slog.Error("failed to listen and serve",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
