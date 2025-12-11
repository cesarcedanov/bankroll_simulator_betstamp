package service

import (
	"bankroll_simulator_betstamp/model"
	"math/rand"
	"testing"
)

func TestAmericanOddsConversion(t *testing.T) {
	tests := []struct {
		odds     int
		expected float64
		name     string
	}{
		{-110, 0.524, "Negative odds favorite"},
		{-200, 0.667, "Heavy favorite"},
		{150, 0.400, "Positive odds underdog"},
		{100, 0.500, "Even odds"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AmericanOddsToImpliedProb(tt.odds)
			if abs(result-tt.expected) > 0.01 {
				t.Errorf("AmericanOddsToImpliedProb(%d) = %f, want %f", tt.odds, result, tt.expected)
			}
		})
	}
}

func TestPayoutCalculation(t *testing.T) {
	tests := []struct {
		stake    float64
		odds     int
		expected float64
		name     string
	}{
		{100, -110, 90.91, "Favorite payout"},
		{100, 150, 150.0, "Underdog payout"},
		{50, -200, 25.0, "Heavy favorite payout"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePayout(tt.stake, tt.odds)
			if abs(result-tt.expected) > 0.1 {
				t.Errorf("CalculatePayout(%f, %d) = %f, want %f", tt.stake, tt.odds, result, tt.expected)
			}
		})
	}
}

func TestSingleIteration(t *testing.T) {
	req := model.SimulationRequest{
		StartingBankroll: 1000,
		Odds:             -110,
		Edge:             0.02,
		NumBets:          100,
		BetSizing: model.BetSizing{
			Mode:   "flat",
			Amount: 10,
		},
	}

	rng := rand.New(rand.NewSource(42))
	result := RunSingleIteration(req, rng)

	if result.FinalBankroll < 0 {
		t.Error("Final bankroll should not be negative")
	}

	if result.Busted && result.FinalBankroll != 0 {
		t.Error("Busted simulation should have 0 bankroll")
	}
}

func TestFractionalBetting(t *testing.T) {
	req := model.SimulationRequest{
		StartingBankroll: 1000,
		Odds:             -110,
		Edge:             0.1, // High edge for testing
		NumBets:          50,
		BetSizing: model.BetSizing{
			Mode:     "fractional",
			Fraction: 0.05, // 5% of bankroll
		},
	}

	rng := rand.New(rand.NewSource(42))
	result := RunSingleIteration(req, rng)

	if result.Busted {
		t.Log("Busted with high edge - this can happen with bad luck")
	}
}

func TestBustCondition(t *testing.T) {
	req := model.SimulationRequest{
		StartingBankroll: 100,
		Odds:             -110,
		Edge:             -0.1, // Negative edge to encourage busting
		NumBets:          1000,
		BetSizing: model.BetSizing{
			Mode:   "flat",
			Amount: 20, // Large bets relative to bankroll
		},
	}

	rng := rand.New(rand.NewSource(42))
	result := RunSingleIteration(req, rng)

	t.Logf("Busted: %v, Final: %f", result.Busted, result.FinalBankroll)
}
func BenchmarkSingleIteration(b *testing.B) {
	req := model.SimulationRequest{
		StartingBankroll: 10000,
		Odds:             -110,
		Edge:             0.02,
		NumBets:          1000,
		BetSizing: model.BetSizing{
			Mode:     "fractional",
			Fraction: 0.02,
		},
	}

	rng := rand.New(rand.NewSource(42))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RunSingleIteration(req, rng)
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
