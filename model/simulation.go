package model

type Simulation struct {
	Id int
}

type BetSizing struct {
	Mode     string  `json:"mode"` // "flat" or "fractional"
	Amount   float64 `json:"amount,omitempty"`
	Fraction float64 `json:"fraction,omitempty"`
}

type CreateSimulationRequest struct {
	UserId           string     `json:"userId"`
	StartingBankroll int        `json:"startingBankroll"`
	Odds             float32    `json:"odds"`
	Edge             float32    `json:"edge"`
	NumBets          int        `json:"numBets"`
	BetSizing        *BetSizing `json:"betSizing"`
	Iterations       int        `json:"iterations"`
}

type CreateSimulationResponse struct {
	SimulationId string `json:"simulationId"`
	Status       string `json:"status"`
}
