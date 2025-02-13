package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HotShotURL      string        `json:"hotshot_url"`
	ChainID         uint64        `json:"chain_id"`
	PollingInterval time.Duration `json:"polling_interval"`
}

func loadConfig() Config {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		os.Exit(1)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		fmt.Println("Error parsing config file:", err)
		os.Exit(1)
	}
	return cfg
}

func main() {
	cfg := loadConfig()
	fmt.Printf("Loaded Config: %+v\n", cfg)

	ticker := time.NewTicker(cfg.PollingInterval * time.Second / 2)
	defer ticker.Stop()

	var temp_height uint64

	for range ticker.C {
		blockHeight, err := fetchBlockHeight(cfg)
		if err != nil {
			fmt.Println("Error fetching block height:", err)
			continue
		}

		if blockHeight != temp_height {
			temp_height = blockHeight
			fmt.Println("Latest Block Height:", blockHeight)
		}

		availBody, err := fetchTransactions(cfg, blockHeight)

		if err != nil {
			fmt.Println("Error fetching availability:", err)
			continue
		}
		if strings.Contains(string(availBody), "FetchBlock") {
			continue
		}
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, availBody, "", "  "); err == nil {
			fmt.Printf("Hotshot Namespace Transactions:\n%s\n", prettyJSON.String())
		} else {
			fmt.Printf("Hotshot Namespace Transactions: %s\n", availBody)
		}
	}
}

func fetchBlockHeight(cfg Config) (uint64, error) {
	resp, err := http.Get(cfg.HotShotURL + "/status/block-height")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	blockHeight, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
	if err != nil {
		return 0, err
	}
	return blockHeight, nil
}

func fetchTransactions(cfg Config, blockHeight uint64) ([]byte, error) {
	availURL := fmt.Sprintf("%s/availability/block/%d/namespace/%d", cfg.HotShotURL, blockHeight, cfg.ChainID)
	resp, err := http.Get(availURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
