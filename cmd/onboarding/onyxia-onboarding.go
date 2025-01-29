package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			log.Printf("Failed to write response: %v", err)
		}
	})

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
