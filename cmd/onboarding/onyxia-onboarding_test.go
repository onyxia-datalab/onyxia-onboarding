package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// setupRouter initializes the router for the application.
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})
	return router
}

func TestHelloWorld(t *testing.T) {
	router := setupRouter() // No need to import since it’s in the same package

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	expected := "Hello World!"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rr.Body.String())
	}
}
