package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
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

	log.Println("Running on", addr)
	http.ListenAndServe(addr, r)
}

func urLWithPrefix(url string) string {
	return fmt.Sprintf("/betstamp/%s", url)
}
