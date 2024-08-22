package operation

import "github.com/ethereum/go-ethereum/common"

type EvmOperation struct {
	ChainSelector uint64
	To            common.Address
	Data          string
	Value         uint64
}

func (e EvmOperation) GetOpType() string {
	return "EvmOperation"
}

func (e EvmOperation) GetOpChainID() uint64 {
	return e.ChainSelector
}

func (e EvmOperation) GetOpData() interface{} {
	return e
}
