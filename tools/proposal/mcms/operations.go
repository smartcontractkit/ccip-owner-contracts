package mcms

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainIdentifier uint64

type Operation struct {
	To           common.Address `json:"to"`
	Data         []byte         `json:"data"`
	Value        *big.Int       `json:"value"`
	ContractType string         `json:"contractType"`
	Tags         []string       `json:"tags"`
}

type ChainOperation struct {
	ChainIdentifier `json:"chainIdentifier"`
	Operation
}
