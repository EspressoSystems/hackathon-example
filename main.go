package main

import (
	"log"
	"time"

	"hackathon-example/config"
)

func main() {
	InitLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Listening to rollup at chain ID: %d", cfg.ChainID)
	log.Printf("Checking for transactions with value larger than %s, initiated by wallet %s", cfg.Value.String(), cfg.From)

	// Frequent polling is required to ensure transactions are captured as soon as
	// they are available, before moving to the next block.
	ticker := time.NewTicker(cfg.PollingInterval * time.Second / 2)
	defer ticker.Stop()

	var temp_height uint64 = 0

	for range ticker.C {
		ProcessTransactions(cfg, &temp_height)
	}
}
