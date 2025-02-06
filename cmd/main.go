package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/api/route"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func main() {

	app := bootstrap.App()
	env := app.Env

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(middleware.Heartbeat("/"))
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   env.Security.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Link", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	route.Setup(&app, r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
