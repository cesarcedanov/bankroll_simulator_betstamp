package model

import (
	"sync"
	"time"
)

type BetSizing struct {
	Mode     string  `json:"mode"` // "flat" or "fractional"
	Amount   float64 `json:"amount,omitempty"`
	Fraction float64 `json:"fraction,omitempty"`
}

type SimulationRequest struct {
	UserId           string    `json:"user_id"`
	StartingBankroll float64   `json:"starting_bankroll"`
	Odds             float64   `json:"odds"`
	Edge             float64   `json:"edge"`
	NumBets          int       `json:"num_bets"`
	BetSizing        BetSizing `json:"bet_sizing"`
	Iterations       int       `json:"iterations"`
}

type SimulationResponse struct {
	SimulationId string `json:"simulation_id"`
	Status       string `json:"status"`
}

type FinalBankrollDistribution struct {
	P10 float64 `json:"p10"`
	P50 float64 `json:"p50"`
	P90 float64 `json:"p90"`
}

type SimulationResult struct {
	SimulationId              string                    `json:"simulation_id"`
	RiskOfRuin                float64                   `json:"risk_of_ruin"`
	FinalBankrollDistribution FinalBankrollDistribution `json:"final_bankroll_distribution"`
	ExpectedFinalBankroll     float64                   `json:"expected_final_bankroll"`
	ROI                       float64                   `json:"roi"`
	Status                    string                    `json:"status"`
}

type Simulation struct {
	Id        string
	UserId    string
	Request   SimulationRequest
	Result    *SimulationResult
	Status    string
	CreatedAt time.Time
	Mutex     sync.RWMutex
}
