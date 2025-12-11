package server

import (
	"bankroll_simulator_betstamp/model"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) HandleCreateSimulation(w http.ResponseWriter, r *http.Request) {
	var req model.CreateSimulationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	json.NewEncoder(w).Encode(model.CreateSimulationResponse{
		SimulationId: "test_123",
		Status:       model.Status_Running,
	})
}

func (s *Server) HandleSimulationResult(w http.ResponseWriter, r *http.Request) {
	simulationId := chi.URLParam(r, "id")
	w.Write([]byte("simulationId:" + simulationId))
}
