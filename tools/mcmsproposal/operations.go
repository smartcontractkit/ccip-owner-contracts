package mcmsproposal

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Operation struct {
	To           common.Address
	Data         []byte
	Value        *big.Int
	ContractType string   `json:"contractType"`
	Tags         []string `json:"tags"`
}

type ChainOperation struct {
	ChainIdentifier string `json:"chainIdentifier"`
	Operation
}
