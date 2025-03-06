package main

import (
	"log"
	"time"

	"hackathon-example/config"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	InitLogger()
	log.Println("Starting the application...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ticker := time.NewTicker(cfg.PollingInterval * time.Second / 2)
	defer ticker.Stop()

	client, err := ethclient.Dial(cfg.CaffNodeURL)

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	var lastBlockNumber uint64

	for range ticker.C {
		searchLatestTransactions(cfg, client, &lastBlockNumber)
	}
}
