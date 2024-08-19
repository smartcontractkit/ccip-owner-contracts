package configwrappers

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
)

type WrappedManyChainMultisig struct {
	gethwrappers.ManyChainMultiSig
}

func DeployAndConfigureManyChainMultisig(auth *bind.TransactOpts, backend bind.ContractBackend, config *Config) (common.Address, []*types.Transaction, *WrappedManyChainMultisig, error) {
	mcmsAddress, tx, mcmsObj, err := gethwrappers.DeployManyChainMultiSig(auth, backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	wrappedMcmsObj := &WrappedManyChainMultisig{*mcmsObj}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, config) // TODO: can the same TransactOpts be used for both transactions?
	if err != nil {
		return common.Address{}, []*types.Transaction{tx}, wrappedMcmsObj, err
	}

	return mcmsAddress, []*types.Transaction{tx, setConfigTx}, wrappedMcmsObj, err
}

func (w *WrappedManyChainMultisig) SetConfig(opts *bind.TransactOpts, config *Config) (*types.Transaction, error) {
	// Validate the config
	if err := config.Validate(); err != nil {
		return nil, err
	}

	groupQuorums, groupParents, signers, signerGroups := config.ExtractSetConfigInputs()
	return w.ManyChainMultiSig.SetConfig(opts, signers, signerGroups, groupQuorums, groupParents, false)
}

func (w *WrappedManyChainMultisig) GetConfig(opts *bind.CallOpts) (*Config, error) {
	config, err := w.ManyChainMultiSig.GetConfig(opts)
	if err != nil {
		return nil, err
	}

	return NewConfigFromRaw(config), nil
}
