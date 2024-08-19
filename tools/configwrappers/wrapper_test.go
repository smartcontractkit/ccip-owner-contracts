package configwrappers

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/stretchr/testify/assert"
)

func setupSimulatedBackend() (*bind.TransactOpts, *backends.SimulatedBackend) {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	auth.GasLimit = uint64(8000000)

	genesisAlloc := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(10000000000)},
	}
	blockGasLimit := uint64(8000000)
	// sim := simulated.NewBackend(genesisAlloc, blockGasLimit)
	sim := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	return auth, sim
}

func TestDeployAndConfigureManyChainMultisig_Success(t *testing.T) {
	auth, sim := setupSimulatedBackend()
	config := NewConfig(1, []common.Address{common.HexToAddress("0x1")}, []Config{})

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, config)
	assert.NoError(t, err)
	sim.Commit()

	assert.NotNil(t, address)
	assert.NotNil(t, tx)
	assert.NotNil(t, setConfigTx)
}

func TestDeployAndConfigureManyChainMultisig_Failure(t *testing.T) {
	auth, sim := setupSimulatedBackend()

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	// Invalid config to trigger failure
	invalidConfig := &Config{}
	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, invalidConfig)

	assert.Error(t, err)
	assert.Nil(t, setConfigTx)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
}

func TestSetConfig_Success(t *testing.T) {
	auth, sim := setupSimulatedBackend()
	config := &Config{}

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, config)
	assert.NoError(t, err)
	sim.Commit()

	assert.NotNil(t, address)
	assert.NotNil(t, tx)
	assert.NotNil(t, setConfigTx)
}

func TestSetConfig_Failure(t *testing.T) {
	auth, sim := setupSimulatedBackend()

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	// Invalid config to trigger failure
	invalidConfig := &Config{}
	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, invalidConfig)

	assert.Error(t, err)
	assert.Nil(t, setConfigTx)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
}

func TestGetConfig_Success(t *testing.T) {
	auth, sim := setupSimulatedBackend()

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	config, err := wrappedMcmsObj.GetConfig(&bind.CallOpts{})

	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
}

func TestGetConfig_Failure(t *testing.T) {
	auth, sim := setupSimulatedBackend()

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim)
	assert.NoError(t, err)
	sim.Commit()

	// Simulate failure by using an invalid call option
	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	config, err := wrappedMcmsObj.GetConfig(nil)

	assert.Error(t, err)
	assert.Nil(t, config)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
}
