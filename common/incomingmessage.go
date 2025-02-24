// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	L1MessageType_L2Message             = 3
	L1MessageType_EndOfBlock            = 6
	L1MessageType_L2FundedByL1          = 7
	L1MessageType_RollupEvent           = 8
	L1MessageType_SubmitRetryable       = 9
	L1MessageType_BatchForGasEstimation = 10 // probably won't use this in practice
	L1MessageType_Initialize            = 11
	L1MessageType_EthDeposit            = 12
	L1MessageType_BatchPostingReport    = 13
	L1MessageType_Invalid               = 0xFF
)

const MaxL2MessageSize = 256 * 1024

var PackInternalTxDataBatchPostingReport func(...interface{}) ([]byte, error)

type L1IncomingMessageHeader struct {
	Kind        uint8          `json:"kind"`
	Poster      common.Address `json:"sender"`
	BlockNumber uint64         `json:"blockNumber"`
	Timestamp   uint64         `json:"timestamp"`
	RequestId   *common.Hash   `json:"requestId" rlp:"nilList"`
	L1BaseFee   *big.Int       `json:"baseFeeL1"`
}

type L1IncomingMessage struct {
	Header *L1IncomingMessageHeader `json:"header"`
	L2msg  []byte                   `json:"l2Msg"`

	// Only used for `L1MessageType_BatchPostingReport`
	BatchGasCost *uint64 `json:"batchGasCost,omitempty" rlp:"optional"`
}

func AddressFromReader(rd io.Reader) (common.Address, error) {
	buf := make([]byte, 20)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(buf), nil
}

func HashFromReader(rd io.Reader) (common.Hash, error) {
	buf := make([]byte, 32)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(buf), nil
}

func Uint256FromReader(rd io.Reader) (*big.Int, error) {
	asHash, err := HashFromReader(rd)
	if err != nil {
		return nil, err
	}
	return asHash.Big(), nil
}

func Uint64FromReader(rd io.Reader) (uint64, error) {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(rd, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint64(buf), nil
}

func AddressFrom256FromReader(rd io.Reader) (common.Address, error) {
	h, err := HashFromReader(rd)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(h.Bytes()[12:]), nil
}

func ParseBatchPostingReportMessageFields(rd io.Reader) (*big.Int, common.Address, common.Hash, uint64, *big.Int, uint64, error) {
	batchTimestamp, err := HashFromReader(rd)
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	batchPosterAddr, err := AddressFromReader(rd)
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	dataHash, err := HashFromReader(rd)
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	batchNum, err := HashFromReader(rd)
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	l1BaseFee, err := HashFromReader(rd)
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	extraGas, err := Uint64FromReader(rd)
	if errors.Is(err, io.EOF) {
		// This field isn't always present
		extraGas = 0
		err = nil
	}
	if err != nil {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, err
	}
	batchNumBig := batchNum.Big()
	if !batchNumBig.IsUint64() {
		return nil, common.Address{}, common.Hash{}, 0, nil, 0, fmt.Errorf("batch number %v is not a uint64", batchNumBig)
	}
	return batchTimestamp.Big(), batchPosterAddr, dataHash, batchNumBig.Uint64(), l1BaseFee.Big(), extraGas, nil
}
