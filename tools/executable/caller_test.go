package executable

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	"github.com/stretchr/testify/assert"
)

func setupSimulatedBackend() (*configwrappers.WrappedManyChainMultisig, *bind.TransactOpts, *simulated.Backend, error) {
	key, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	if err != nil {
		return nil, nil, nil, err
	}

	auth.GasLimit = uint64(8000000)

	genesisAlloc := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(1e18)},
	}
	blockGasLimit := uint64(8000000)
	sim := simulated.New(genesisAlloc, blockGasLimit)

	// Deploy a ManyChainMultiSig contract
	_, _, mcms, err := configwrappers.DeployWrappedManyChainMultisig(auth, sim.Client())
	if err != nil {
		return nil, nil, nil, err
	}
	sim.Commit()

	return mcms, auth, sim, nil
}

func TestCaller_GetConfigs_Success_SingleChain(t *testing.T) {
	mcms, auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x1"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx, err := mcms.SetConfig(auth, &config)
	assert.NotNil(t, tx)
	assert.NoError(t, err)
	sim.Commit()

	// Create a Caller
	caller, err := NewCaller(map[string]common.Address{
		"chain1": mcms.Address(),
	}, map[string]ContractDeployBackend{
		"chain1": sim.Client(),
	})
	assert.NoError(t, err)

	// Get the config
	configs, err := caller.GetConfigs()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(configs))
	assert.Equal(t, config.ToRawConfig(), configs["chain1"])
}

func TestCaller_GetOpCounts_Success_SingleChain(t *testing.T) {
	mcms, auth, sim, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x1"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx, err := mcms.SetConfig(auth, &config)
	assert.NotNil(t, tx)
	assert.NoError(t, err)
	sim.Commit()

	// Create a Caller
	caller, err := NewCaller(map[string]common.Address{
		"chain1": mcms.Address(),
	}, map[string]ContractDeployBackend{
		"chain1": sim.Client(),
	})
	assert.NoError(t, err)

	// Get the config
	opCounts, err := caller.GetCurrentOpCounts()
	opCount := opCounts["chain1"]
	assert.NoError(t, err)
	assert.Equal(t, 1, len(opCounts))
	assert.Equal(t, int64(0), opCount.Int64())
}

func TestCaller_GetConfigs_Success_MultipleChains(t *testing.T) {
	mcms1, auth1, sim1, err := setupSimulatedBackend()
	assert.NoError(t, err)

	mcms2, auth2, sim2, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config1 := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x1"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx1, err := mcms1.SetConfig(auth1, &config1)
	assert.NotNil(t, tx1)
	assert.NoError(t, err)
	sim1.Commit()

	// Define a valid config
	config2 := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x2"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx2, err := mcms2.SetConfig(auth2, &config2)
	assert.NotNil(t, tx2)
	assert.NoError(t, err)
	sim2.Commit()

	// Create a Caller
	caller, err := NewCaller(map[string]common.Address{
		"chain1": mcms1.Address(),
		"chain2": mcms2.Address(),
	}, map[string]ContractDeployBackend{
		"chain1": sim1.Client(),
		"chain2": sim2.Client(),
	})
	assert.NoError(t, err)

	// Get the config
	configs, err := caller.GetConfigs()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(configs))
	assert.Equal(t, config1.ToRawConfig(), configs["chain1"])
	assert.Equal(t, config2.ToRawConfig(), configs["chain2"])
}

func TestCaller_GetOpCounts_Success_MultipleChains(t *testing.T) {
	mcms1, auth1, sim1, err := setupSimulatedBackend()
	assert.NoError(t, err)

	mcms2, auth2, sim2, err := setupSimulatedBackend()
	assert.NoError(t, err)

	// Define a valid config
	config1 := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x1"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx1, err := mcms1.SetConfig(auth1, &config1)
	assert.NotNil(t, tx1)
	assert.NoError(t, err)
	sim1.Commit()

	// Define a valid config
	config2 := configwrappers.Config{
		Quorum: 1,
		Signers: []common.Address{
			common.HexToAddress("0x2"),
		},
		GroupSigners: []configwrappers.Config{},
	}

	// Set the config
	tx2, err := mcms2.SetConfig(auth2, &config2)
	assert.NotNil(t, tx2)
	assert.NoError(t, err)
	sim2.Commit()

	// Create a Caller
	caller, err := NewCaller(map[string]common.Address{
		"chain1": mcms1.Address(),
		"chain2": mcms2.Address(),
	}, map[string]ContractDeployBackend{
		"chain1": sim1.Client(),
		"chain2": sim2.Client(),
	})
	assert.NoError(t, err)

	// Get the config
	opCounts, err := caller.GetCurrentOpCounts()
	opCount1 := opCounts["chain1"]
	opCount2 := opCounts["chain2"]
	assert.NoError(t, err)
	assert.Equal(t, 2, len(opCounts))
	assert.Equal(t, int64(0), opCount1.Int64())
	assert.Equal(t, int64(0), opCount2.Int64())
}
