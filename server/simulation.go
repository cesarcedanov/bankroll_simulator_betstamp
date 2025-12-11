package server

import (
	"bankroll_simulator_betstamp/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (s *Server) HandleCreateSimulation(w http.ResponseWriter, r *http.Request) {
	var req model.SimulationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.StartingBankroll <= 0 {
		http.Error(w, "Starting bankroll must be greater than 0", http.StatusBadRequest)
		return
	}
	if req.NumBets <= 0 || req.Iterations <= 0 {
		http.Error(w, "Number of bets and iterations must be greater than 0", http.StatusBadRequest)
		return
	}
	if strings.ToLower(req.BetSizing.Mode) != "flat" && strings.ToLower(req.BetSizing.Mode) != "fractional" {
		http.Error(w, "bet_sizing.mode must be 'flat' or 'fractional'", http.StatusBadRequest)
		return
	}

	newSimulation := &model.Simulation{
		Id:        fmt.Sprintf("sim_%s", uuid.New().String()),
		UserId:    req.UserId,
		Request:   req,
		Status:    model.Status_Running,
		CreatedAt: time.Now(),
	}

	s.Storage.SaveSimulation(newSimulation)

	json.NewEncoder(w).Encode(model.SimulationResponse{
		SimulationId: newSimulation.Id,
		Status:       newSimulation.Status,
	})
}

func (s *Server) HandleSimulationResult(w http.ResponseWriter, r *http.Request) {
	simulationId := chi.URLParam(r, "id")
	simulation, exists := s.Storage.GetSimulation(simulationId)
	if !exists {
		http.Error(w, "Simulation not found", http.StatusNotFound)
		return
	}
	simulation.Mutex.RLock()
	defer simulation.Mutex.RUnlock()

	// TODO: Check simulation status and return results
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(simulation)
}
