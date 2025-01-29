package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/onyxia-datalab/onyxia-onboarding/api"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

// setupRouter initializes the router for the application.
func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	server := api.NewServer()
	api.HandlerFromMux(server, r)
	return r
}

func TestOnboardingEndpoint(t *testing.T) {
	router := setupRouter() // No need to import since it’s in the same package

	requestBody, _ := json.Marshal(api.OnboardingRequest{
		Group: strPtr("test-group"),
	})

	req, err := http.NewRequest("POST", "/onboarding", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code) // Ensure 200 OK
}
