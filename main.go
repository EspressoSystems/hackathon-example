package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"hackathon-example/types_utils"

	"github.com/EspressoSystems/espresso-sequencer-go/types/common"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type HotShotResponse struct {
	Proof        *json.RawMessage     `json:"proof"`
	Transactions []common.Transaction `json:"transactions"`
}

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
			fmt.Println("FetchBlock response received")
			continue
		}

		var hotShotResponse HotShotResponse
		if err := json.Unmarshal(availBody, &hotShotResponse); err != nil {
			fmt.Println("Error parsing JSON:", err)
			continue
		}

		for _, tx := range hotShotResponse.Transactions {
			_, _, _, messages, err := types_utils.ParseHotShotPayload(tx.Payload)
			if err != nil {
				fmt.Println("Error parsing hotshot payload:", err)
				continue
			}
			for _, message := range messages {
				var messageWithMetadata types_utils.MessageWithMetadata
				err = rlp.DecodeBytes(message, &messageWithMetadata)
				if err != nil {
					fmt.Println("Error parsing message:", err)
					continue
				}
				zero := big.NewInt(int64(cfg.ChainID))
				t, err := types_utils.ParseL2Transactions(messageWithMetadata.Message, zero)
				if err != nil {
					fmt.Println("Error parsing L2 transactions:", err)
					continue
				}

				var tx *types.Transaction = t[0]
				var data = tx.Data()

				var dataHash eth_common.Hash = crypto.Keccak256Hash(data)

				fmt.Printf("__________________________________________________________\n")
				fmt.Printf("Transaction Type: %s\n", types_utils.GetTransactionTypeDescription(tx.Type()))
				fmt.Printf("Transaction Destination: %s\n", tx.To())
				fmt.Printf("Transaction Value: %s\n", tx.Value())
				fmt.Printf("Transaction Data: %s\n", data)
				fmt.Printf("Transaction Data Hash: %s\n", dataHash)
				fmt.Printf("__________________________________________________________\n")

				var EIP155Signer = types.NewLondonSigner(big.NewInt(int64(cfg.ChainID)))
				var ArbitrumSigner = types.NewArbitrumSigner(EIP155Signer)

				from, err := types.Sender(ArbitrumSigner, tx)
				if err != nil {
					fmt.Println("Error parsing sender:", err)
				}
				fmt.Printf("Transaction From: %s\n", from)
			}
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
