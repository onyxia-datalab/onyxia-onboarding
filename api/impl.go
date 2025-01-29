package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Ensure `Server` implements `ServerInterface`
var _ ServerInterface = (*Server)(nil)

// Server struct
type Server struct{}

// NewServer initializes a new API server
func NewServer() *Server {
	return &Server{}
}

func (s *Server) Onboard(w http.ResponseWriter, r *http.Request) {
	var req OnboardingRequest

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Log received request
	log.Printf("Received onboarding request for group: %s", *req.Group)

	// Respond with 200 OK
	w.WriteHeader(http.StatusOK)
}
