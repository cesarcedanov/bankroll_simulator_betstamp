package storage

import (
	"bankroll_simulator_betstamp/model"
	"sync"
)

// Storage will be used to store simulations in memory
type Storage struct {
	simulations map[string]*model.Simulation
	userSims    map[string][]string
	mu          sync.RWMutex
}

func NewStorage() *Storage {
	return &Storage{
		simulations: make(map[string]*model.Simulation),
		userSims:    make(map[string][]string),
	}
}

// SaveSimulation will expect a simulation and save it to the in-memory storage
func (s *Storage) SaveSimulation(sim *model.Simulation) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Add simulation to in-memory storage
	s.simulations[sim.Id] = sim
	// Add simulation id to user's simulations'
	s.userSims[sim.UserId] = append(s.userSims[sim.UserId], sim.Id)
}

// GetSimulation will return a simulation from the in-memory storage
func (s *Storage) GetSimulation(id string) (*model.Simulation, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Get simulation from in-memory storage
	sim, exists := s.simulations[id]
	return sim, exists
}

// GetUserSimulations will return all simulations for a given user ID
func (s *Storage) GetUserSimulations(userId string) []*model.Simulation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	// Find all the Simulation IDs for the given user ID
	simIDs := s.userSims[userId]
	result := make([]*model.Simulation, 0, len(simIDs))
	for _, id := range simIDs {
		// If the Simulation ID exists in the in-memory storage, add it to the result
		if sim, exists := s.simulations[id]; exists {
			result = append(result, sim)
		}
	}
	return result
}
