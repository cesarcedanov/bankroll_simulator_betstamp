package service

import (
	"bankroll_simulator_betstamp/model"
	"bankroll_simulator_betstamp/storage"
	"math/rand"
	"sort"
)

// AmericanOddsToImpliedProb
func AmericanOddsToImpliedProb(odds int) float64 {
	if odds < 0 {
		// Negative odds (favorite)
		return float64(-odds) / (float64(-odds) + 100)
	}
	// Positive odds (underdog)
	return 100 / (float64(odds) + 100)
}

// CalculatePayout
func CalculatePayout(stake float64, odds int) float64 {
	if odds < 0 {
		// Negative odds
		return stake * (100.0 / float64(-odds))
	}
	// Positive odds
	return stake * (float64(odds) / 100.0)
}

// RunSingleIteration
func RunSingleIteration(req model.SimulationRequest, rng *rand.Rand) JobResult {
	bankroll := req.StartingBankroll
	// Implied Probability
	pImplied := AmericanOddsToImpliedProb(req.Odds)
	//Applying an Edge
	pTrue := pImplied * (1 + req.Edge)

	busted := false

	for bet := 0; bet < req.NumBets; bet++ {
		// Determine stake
		var stake float64
		if req.BetSizing.Mode == "flat" {
			stake = req.BetSizing.Amount
		} else {
			stake = bankroll * req.BetSizing.Fraction
		}

		// Check if can place bet
		if stake > bankroll || bankroll <= 0 {
			busted = true
			bankroll = 0
			break
		}

		// Simulate win/loss
		if rng.Float64() < pTrue {
			// Win
			payout := CalculatePayout(stake, req.Odds)
			bankroll += payout
		} else {
			// Loss
			bankroll -= stake
		}

		if bankroll <= 0 {
			busted = true
			bankroll = 0
			break
		}
	}

	return JobResult{
		FinalBankroll: bankroll,
		Busted:        busted,
	}
}

// RunSimulation
func RunSimulation(sim *model.Simulation, storage *storage.Storage) {
	numWorkers := 8
	if sim.Request.Iterations < 100 {
		numWorkers = 2
	}

	pool := NewWorkerPool(numWorkers)
	pool.Start()

	// Submit jobs
	go func() {
		for i := 0; i < sim.Request.Iterations; i++ {
			pool.Submit(Job{
				IterationID: i,
				Request:     sim.Request,
			})
		}
	}()

	// Collect results
	finalBankrolls := make([]float64, 0, sim.Request.Iterations)
	bustCount := 0

	for i := 0; i < sim.Request.Iterations; i++ {
		result := <-pool.results
		finalBankrolls = append(finalBankrolls, result.FinalBankroll)
		if result.Busted {
			bustCount++
		}
	}

	pool.Stop()

	// Calculate statistics
	sort.Float64s(finalBankrolls)

	riskOfRuin := float64(bustCount) / float64(sim.Request.Iterations)

	p10 := finalBankrolls[int(float64(len(finalBankrolls))*0.1)]
	p50 := finalBankrolls[int(float64(len(finalBankrolls))*0.5)]
	p90 := finalBankrolls[int(float64(len(finalBankrolls))*0.9)]

	sum := 0.0
	for _, v := range finalBankrolls {
		sum += v
	}
	expectedFinal := sum / float64(len(finalBankrolls))

	roi := (expectedFinal - sim.Request.StartingBankroll) / sim.Request.StartingBankroll

	// Update simulation with results
	sim.Mutex.Lock()
	sim.Result = &model.SimulationResult{
		SimulationId: sim.Id,
		RiskOfRuin:   riskOfRuin,
		FinalBankrollDistribution: model.FinalBankrollDistribution{
			P10: p10,
			P50: p50,
			P90: p90,
		},
		ExpectedFinalBankroll: expectedFinal,
		ROI:                   roi,
		Status:                model.Status_Done,
	}

	sim.Status = model.Status_Done
	sim.Mutex.Unlock()
}
