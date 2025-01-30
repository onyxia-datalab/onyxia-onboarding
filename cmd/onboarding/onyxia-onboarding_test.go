package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/onyxia-datalab/onyxia-onboarding/bootstrap"
	"github.com/stretchr/testify/assert"
)

// TestSetupRouter ensures the router initializes properly.
func TestSetupRouter(t *testing.T) {
	r := setupRouter()
	if r == nil {
		t.Fatal("Router setup failed, got nil")
	}
}

// TestHelloWorld verifies the "/" endpoint.
func TestHelloWorld(t *testing.T) {
	r := setupRouter()

	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Hello World!", strings.TrimSpace(rr.Body.String()))
}

// TestCORS ensures CORS headers are correctly applied.
func TestCORS(t *testing.T) {
	r := setupRouter()

	req, err := http.NewRequest("OPTIONS", "/", nil)
	assert.NoError(t, err)

	req.Header.Set("Origin", "http://example.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "Content-Type")
}

// setupRouter initializes a test router with middleware and routes.
func setupRouter() *chi.Mux {
	app := bootstrap.App()
	env := app.Env

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

	return r
}
