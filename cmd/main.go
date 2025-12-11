package main

import (
	"bankroll_simulator_betstamp/server"
	"log"
	"time"
)

func main() {
	log.Printf("Server started at %s", time.Now())
	srv := server.NewServer()
	srv.Run(":8080")
}
