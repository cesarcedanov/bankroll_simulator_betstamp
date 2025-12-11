package server

import (
	"bankroll_simulator_betstamp/storage"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Storage *storage.Storage
}

func NewServer() *Server {
	return &Server{
		Storage: storage.NewStorage(),
	}
}

func (s *Server) Run(addr string) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("BetStamp API"))
	})

	r.Post(urLWithPrefix("simulations"), s.HandleCreateSimulation)
	r.Get(urLWithPrefix("simulations/{id}/result"), s.HandleSimulationResult)
	r.Get(urLWithPrefix("users/{id}/simulations"), s.HandleUserSimulations)

	log.Println("Server running on", addr)

	http.ListenAndServe(addr, r)
}

func urLWithPrefix(url string) string {
	return fmt.Sprintf("/betstamp/%s", url)
}
