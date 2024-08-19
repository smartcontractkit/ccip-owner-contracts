package configwrappers

import (
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	genesisAlloc := map[common.Address]types.Account{
		auth.From: {Balance: big.NewInt(1e18)},
	}
	blockGasLimit := uint64(8000000)
	sim := simulated.NewBackend(genesisAlloc, simulated.WithBlockGasLimit(blockGasLimit))

	return auth, sim, nil
}

func TestDeployAndConfigureManyChainMultisig_Success(t *testing.T) {
	var wg sync.WaitGroup
	resultChan := make(chan struct {
		address common.Address
		txs     []*types.Transaction
		mcms    *WrappedManyChainMultisig
		err     error
	}, 1)

	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config := NewConfig(1, []common.Address{common.HexToAddress("0x1")}, []Config{})

	// Call the DeployAndConfigureManyChainMultisig function asynchronously
	wg.Add(1)
	go func() {
		defer wg.Done()

		address, txs, mcms, err := DeployAndConfigureManyChainMultisig(auth, sim.Client(), config)

		// Send the results to the channel
		resultChan <- struct {
			address common.Address
			txs     []*types.Transaction
			mcms    *WrappedManyChainMultisig
			err     error
		}{address, txs, mcms, err}
	}()

	// While async function is running, continue to mine blocks
	for i := 0; i < 5; i++ {
		time.Sleep(250 * time.Millisecond)
		sim.Commit()
	}

	// Wait for the async function to complete
	wg.Wait()
	close(resultChan)

	// Read the result from the channel
	result := <-resultChan

	// Validate the result
	assert.NoError(t, result.err)
	assert.Len(t, result.txs, 2)
	assert.NotNil(t, result.address)

	// Validate the contract deployment and configuration
	configFromContract, err := result.mcms.GetConfig(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, config, configFromContract)
}

func TestDeployAndConfigureManyChainMultisig_Failure_InvalidConfig(t *testing.T) {
	var wg sync.WaitGroup
	resultChan := make(chan struct {
		address common.Address
		txs     []*types.Transaction
		mcms    *WrappedManyChainMultisig
		err     error
	}, 1)

	auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config := &Config{GroupSigners: []Config{}, Quorum: 0, Signers: []common.Address{}}

	// Call the DeployAndConfigureManyChainMultisig function asynchronously
	wg.Add(1)
	go func() {
		defer wg.Done()

		address, txs, mcms, err := DeployAndConfigureManyChainMultisig(auth, sim.Client(), config)

		// Send the results to the channel
		resultChan <- struct {
			address common.Address
			txs     []*types.Transaction
			mcms    *WrappedManyChainMultisig
			err     error
		}{address, txs, mcms, err}
	}()

	// While async function is running, continue to mine blocks
	for i := 0; i < 5; i++ {
		time.Sleep(250 * time.Millisecond)
		sim.Commit()
	}

	// Wait for the async function to complete
	wg.Wait()
	close(resultChan)

	// Read the result from the channel
	result := <-resultChan

	// Validate the result
	assert.Error(t, result.err)
	assert.Len(t, result.txs, 1)
	assert.NotNil(t, result.address)
	assert.Equal(t, "invalid MCMS config: Quorum must be greater than 0", result.err.Error())

	// Validate the contract deployment and configuration
	configFromContract, err := result.mcms.GetConfig(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, config, configFromContract)
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
