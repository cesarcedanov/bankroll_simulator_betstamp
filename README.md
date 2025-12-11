# Bankroll Risk Simulation Backend

A high-performance **Monte Carlo simulation service** for estimating bankroll risk in betting strategies.  
Built with **Go**, leveraging concurrent worker pools and a clean, RESTful API.

---

## Features

- **Monte Carlo Simulations** — Run thousands of betting scenarios to analyze bankroll volatility and long-term risk.
- **High-Throughput Worker Pool** — Efficient concurrent execution using goroutines and bounded channels.
- **Multiple Bet Sizing Modes** — Supports fixed (flat) betting and fractional (Kelly-style) strategies.
- **Risk Metrics** — Calculates probability of ruin, expected ROI, and percentile bankroll outcomes.
- **RESTful API** — Simple HTTP endpoints for creating and retrieving simulations.
- **Memory Safe** — No goroutine leaks, predictable resource use, scalable to 100k+ iterations.

---

## Installation

### Prerequisites
- Go **1.19+**

### Setup

```bash
# Clone the repository
git clone <repository-url>
cd bankroll_simulation_betstamp

# Install dependencies
go mod download

# Run the service
go run cmd/*.go
```
The server starts at:
http://localhost:8080/betstamp/



## API Endpoints

### Health Check
**GET** localhost:8080/betstamp/health


### Create Simulation
**POST** localhost:8080/betstamp/simulations
Creates and runs a new bankroll simulation.

```
Request Body
{
  "user_id": "u1",
  "starting_bankroll": 10000,
  "odds": -110,
  "edge": 0.02,
  "num_bets": 1000,
  "bet_sizing": {
    "mode": "fractional",
    "fraction": 0.02
  },
  "iterations": 20000
}

Response
{
  "simulation_id": "sim_a1b2c3d4",
  "status": "running"
}
```


### Get Simulation Results
**GET** localhost:8080/betstamp/simulations/:id/result\
Retrieves results for a completed simulation.
```
Response
{
  "simulation_id": "sim_a1b2c3d4",
  "risk_of_ruin": 0.18,
  "final_bankroll_distribution": {
    "p10": 7800,
    "p50": 10550,
    "p90": 14200
  },
  "expected_final_bankroll": 11230.50,
  "roi": 0.123,
  "status": "completed"
}
```


### Get User Simulations
**GET** 
localhost:8080/betstamp/users/:id/simulations
Lists all simulations created by a user.
```
Response
[
  {
    "simulation_id": "sim_a1b2c3d4",
    "status": "completed",
    "created_at": "2024-12-11T10:30:00Z"
  },
  {
    "simulation_id": "sim_e5f6g7h8",
    "status": "running",
    "created_at": "2024-12-11T11:45:00Z"
  }
]
```

## Example Usage
1. Conservative Fractional Betting
```
curl -X POST http://localhost:8080/betstamp/simulations \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "starting_bankroll": 10000,
    "odds": -110,
    "edge": 0.02,
    "num_bets": 1000,
    "bet_sizing": {
      "mode": "fractional",
      "fraction": 0.02
    },
    "iterations": 20000
  }'
```

2. Flat Betting Strategy
```
curl -X POST http://localhost:8080/betstamp/simulations \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user456",
    "starting_bankroll": 5000,
    "odds": 150,
    "edge": 0.05,
    "num_bets": 500,
    "bet_sizing": {
      "mode": "flat",
      "amount": 100
    },
    "iterations": 10000
  }'
```

3. Retrieve Results
```
curl http://localhost:8080/betstamp/simulations/sim_a1b2c3d4/result
```

4. Get All User Simulations
```
curl http://localhost:8080/betstamp/users/user123/simulations
```


## Simulation Mathematics
- **American Odds → Implied Probability**

Favorites (negative odds):
```
p = |odds| / (|odds| + 100)

Example:
-110 → 110 / (110 + 100) = 0.524
```

Underdogs (positive odds):
```
p = 100 / (odds + 100)

Example:
+150 → 100 / 250 = 0.40
```

- **Edge Application**
```
true_probability = implied_probability × (1 + edge)

Example:
0.524 × 1.02 = 0.534
```

- **Payout Calculation**

```
Negative odds:
payout = stake × (100 / |odds|)

Example:
$100 @ -110 → $90.91 profit


Positive odds:
payout = stake × (odds / 100)

Example:
$100 @ +150 → $150 profit
```