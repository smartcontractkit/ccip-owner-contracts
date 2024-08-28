package configwrappers

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/stretchr/testify/assert"
)

func setupSimulatedBackend() (*bind.TransactOpts, *simulated.Backend, error) {
	key, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		return nil, nil, err
	}

	auth.GasLimit = uint64(8000000)

	genesisAlloc := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(1e18)},
	}
	blockGasLimit := uint64(8000000)
	sim := simulated.New(genesisAlloc, blockGasLimit)

	return auth, sim, nil
}

func TestDeploy_Success(t *testing.T) {
	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	address, tx, mcms, err := DeployWrappedManyChainMultisig(auth, sim.Client())
	sim.Commit()

	assert.NoError(t, err)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
	assert.NotNil(t, mcms)

	// Get the config from the contract and validate its default values
	config, err := mcms.GetConfig(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, &Config{Quorum: 0, Signers: []common.Address{}, GroupSigners: []Config{}}, config)
}

func TestSetConfig_Success(t *testing.T) {
	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)
	config := NewConfig(1, []common.Address{common.HexToAddress("0x1")}, []Config{})

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim.Client())
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

func TestSetConfig_Failure_InvalidConfig(t *testing.T) {
	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim.Client())
	assert.NoError(t, err)
	assert.NotNil(t, address)
	assert.NotNil(t, tx)
	sim.Commit()

	// Invalid config to trigger failure
	invalidConfig := &Config{}
	wrappedMcmsObj := &WrappedManyChainMultisig{*mcms}
	setConfigTx, err := wrappedMcmsObj.SetConfig(auth, invalidConfig)

	assert.Error(t, err)
	assert.Equal(t, "invalid MCMS config: Quorum must be greater than 0", err.Error())
	assert.Nil(t, setConfigTx)
}

func TestGetConfig_Success(t *testing.T) {
	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	address, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auth, sim.Client())
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
	_, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Simulate failure by using an invalid address
	mcmsObj, err := gethwrappers.NewManyChainMultiSig(common.HexToAddress("0x01"), sim.Client())
	assert.NoError(t, err)

	// Simulate failure by using an invalid call option
	wrappedMcmsObj := &WrappedManyChainMultisig{*mcmsObj}
	config, err := wrappedMcmsObj.GetConfig(nil)

	assert.Error(t, err)
	assert.Nil(t, config)
}
