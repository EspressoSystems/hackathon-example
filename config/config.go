package config

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"
)

type Config struct {
	HotShotURL      string        `json:"hotshot_url"`
	ChainID         uint64        `json:"chain_id"`
	PollingInterval time.Duration `json:"polling_interval"`
	From            string        `json:"from"`
	Value           *big.Int      `json:"value"`
}

func LoadConfig() (Config, error) {
	data, err := os.ReadFile("config/config.json")
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("error parsing config file: %w", err)
	}
	return cfg, nil
}
