package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) HandleUserSimulations(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	userSimulations := s.Storage.GetUserSimulations(userId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userSimulations)
}
