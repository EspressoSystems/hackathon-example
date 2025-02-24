package common

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/ccoveille/go-safecast"

	"github.com/EspressoSystems/espresso-sequencer-go/types/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const LEN_SIZE int = 8
const INDEX_SIZE int = 8

// Parameters imported from the espresso-sequencer-go package
type HotShotResponse struct {
	Proof        *json.RawMessage     `json:"proof"`
	Transactions []common.Transaction `json:"transactions"`
}

type MessageIndex uint64

func ParseHotShotPayload(payload []byte) (signature []byte, userDataHash []byte, indices []uint64, messages [][]byte, err error) {
	if len(payload) < LEN_SIZE {
		return nil, nil, nil, nil, errors.New("payload too short to parse signature size")
	}

	// Extract the signature size
	signatureSize, err := safecast.ToInt(binary.BigEndian.Uint64(payload[:LEN_SIZE]))
	if err != nil {
		return nil, nil, nil, nil, errors.New("could not convert signature size to int")
	}

	currentPos := LEN_SIZE

	if len(payload[currentPos:]) < signatureSize {
		return nil, nil, nil, nil, errors.New("payload too short for signature")
	}

	// Extract the signature
	signature = payload[currentPos : currentPos+signatureSize]
	currentPos += signatureSize

	indices = []uint64{}
	messages = [][]byte{}

	// Take keccak256 hash of the rest of payload
	userDataHash = crypto.Keccak256(payload[currentPos:])
	// Parse messages
	for {
		if currentPos == len(payload) {
			break
		}
		if len(payload[currentPos:]) < LEN_SIZE+INDEX_SIZE {
			return nil, nil, nil, nil, errors.New("remaining bytes")
		}

		// Extract the index
		index := binary.BigEndian.Uint64(payload[currentPos : currentPos+INDEX_SIZE])
		currentPos += INDEX_SIZE

		// Extract the message size
		messageSize, err := safecast.ToInt(binary.BigEndian.Uint64(payload[currentPos : currentPos+LEN_SIZE]))
		if err != nil {
			return nil, nil, nil, nil, errors.New("could not convert message size to int")
		}
		currentPos += LEN_SIZE

		if len(payload[currentPos:]) < messageSize {
			return nil, nil, nil, nil, errors.New("message size mismatch")
		}

		// Extract the message
		message := payload[currentPos : currentPos+messageSize]
		currentPos += messageSize

		indices = append(indices, index)
		messages = append(messages, message)
	}

	return signature, userDataHash, indices, messages, nil
}
