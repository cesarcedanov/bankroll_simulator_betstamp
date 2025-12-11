package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) HandleUserSimulations(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "id")
	w.Write([]byte("UserId:" + userId))
}
