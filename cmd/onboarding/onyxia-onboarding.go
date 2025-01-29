package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimiddle "github.com/go-chi/chi/v5/middleware"
	oapimiddle "github.com/oapi-codegen/nethttp-middleware"

	"github.com/onyxia-datalab/onyxia-onboarding/api" // Import your generated API package
)

func main() {
	// Create a new Chi router
	r := chi.NewRouter()

	// Load OpenAPI schema for validation
	swagger, err := api.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	// Apply request validation middleware
	r.Use(oapimiddle.OapiRequestValidator(swagger))
	r.Use(chimiddle.Logger)
	// Initialize API server
	server := api.NewServer()
	api.HandlerFromMux(server, r)

	// Start HTTP server
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
