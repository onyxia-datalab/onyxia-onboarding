package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
)

func main() {

	app := bootstrap.App()
	env := app.Env

	// Temp code
	envJSON, _ := json.MarshalIndent(env, "", "  ")
	log.Println("Loaded environment:", string(envJSON))

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   env.Security.CORSAllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Requested-With"},
		ExposedHeaders:   []string{"Link", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
