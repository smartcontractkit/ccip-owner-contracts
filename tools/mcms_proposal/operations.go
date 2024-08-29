package mcms_proposal

import "github.com/ethereum/go-ethereum/common"

type Operation struct {
	To           common.Address
	Data         string
	Value        uint64
	ContractType string   `json:"contractType"`
	Tags         []string `json:"tags"`
}

type ChainOperation struct {
	ChainIdentifier string `json:"chainIdentifier"`
	Operation
}
