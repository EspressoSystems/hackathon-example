package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"hackathon-example/common"
	"hackathon-example/config"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func ProcessTransactions(cfg config.Config, temp_height *uint64) {
	blockHeight, err := fetchBlockHeight(cfg)
	if err != nil {
		log.Println("Error fetching block height:", err)
		return
	}

	// Only log block height if it has changed
	if blockHeight != *temp_height {
		log.Println("Latest Block Height:", blockHeight)
		*temp_height = blockHeight
	}

	transactionsBody, err := fetchTransactions(cfg, blockHeight)

	if err != nil {
		log.Println("Error fetching transactions:", err)
		return
	}

	// If the response includes the word "FetchBlock", skip it.
	// This is a response from the HotShot API when awaiting for the transactions to be available.
	if strings.Contains(string(transactionsBody), "FetchBlock") {
		return
	}

	var hotShotResponse common.HotShotResponse
	if err := json.Unmarshal(transactionsBody, &hotShotResponse); err != nil {
		log.Println("Error parsing JSON:", err)
		return
	}

	for _, tx := range hotShotResponse.Transactions {
		_, _, _, messages, err := common.ParseHotShotPayload(tx.Payload)
		if err != nil {
			log.Println("Error parsing hotshot payload:", err)
			continue
		}
		for _, message := range messages {
			processTransactionMessage(cfg, message)
		}
	}
}

func processTransactionMessage(cfg config.Config, message []byte) {
	var messageWithMetadata common.MessageWithMetadata
	err := rlp.DecodeBytes(message, &messageWithMetadata)
	if err != nil {
		log.Println("Error parsing message:", err)
		return
	}

	transactions, err := common.ParseL2Transactions(messageWithMetadata.Message, big.NewInt(int64(cfg.ChainID)))
	if err != nil {
		log.Println("Error parsing L2 transactions:", err)
		return
	}

	verifyTransactionDetails(cfg, transactions)
}

func verifyTransactionDetails(cfg config.Config, transactions types.Transactions) {
	var EIP155Signer = types.NewLondonSigner(big.NewInt(int64(cfg.ChainID)))
	var ArbitrumSigner = types.NewArbitrumSigner(EIP155Signer)

	for _, tx := range transactions {
		from, err := types.Sender(ArbitrumSigner, tx)

		if err != nil {
			log.Println("Error parsing sender:", err)
			continue
		}

		var value = tx.Value()

		log.Println("_________________________________________________________________________________________________")
		if from.String() == cfg.From && value.Cmp(cfg.Value) >= 0 {
			log.Printf("Alarming transaction detected: %s\n", tx.Hash().String())
			log.Printf("Transaction Type: %s\n", GetTransactionTypeDescription(tx.Type()))
			log.Printf("Transaction Value: %s\n", value)
			log.Printf("Transaction Destination: %s\n", tx.To())
		} else {
			log.Printf("Block transaction: %s\n", tx.Hash().String())
		}
		log.Println("_________________________________________________________________________________________________")
	}
}

func fetchBlockHeight(cfg config.Config) (uint64, error) {
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

func fetchTransactions(cfg config.Config, blockHeight uint64) ([]byte, error) {
	availURL := fmt.Sprintf("%s/availability/block/%d/namespace/%d", cfg.HotShotURL, blockHeight, cfg.ChainID)
	resp, err := http.Get(availURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Map transaction type numbers to descriptive strings
func GetTransactionTypeDescription(txType uint8) string {
	switch txType {
	case types.LegacyTxType:
		return "Legacy Transaction"
	case types.AccessListTxType:
		return "Access List Transaction"
	case types.DynamicFeeTxType:
		return "Dynamic Fee Transaction"
	case types.ArbitrumDepositTxType:
		return "Arbitrum Deposit Transaction"
	case types.ArbitrumUnsignedTxType:
		return "Arbitrum Unsigned Transaction"
	case types.ArbitrumContractTxType:
		return "Arbitrum Contract Transaction"
	case types.ArbitrumRetryTxType:
		return "Arbitrum Retry Transaction"
	case types.ArbitrumSubmitRetryableTxType:
		return "Arbitrum Submit Retryable Transaction"
	case types.ArbitrumInternalTxType:
		return "Arbitrum Internal Transaction"
	case types.ArbitrumLegacyTxType:
		return "Arbitrum Legacy Transaction"
	default:
		return "Unknown Transaction Type"
	}
}
