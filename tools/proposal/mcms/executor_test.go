package mcms

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/configwrappers"
	owner_errors "github.com/smartcontractkit/ccip-owner-contracts/tools/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/gethwrappers"
	"github.com/stretchr/testify/assert"
)

func setupSimulatedBackendWithMCMS(numSigners uint64) ([]*ecdsa.PrivateKey, []*bind.TransactOpts, *backends.SimulatedBackend, *gethwrappers.ManyChainMultiSig, error) {
	// Generate a private key
	keys := make([]*ecdsa.PrivateKey, numSigners)
	auths := make([]*bind.TransactOpts, numSigners)
	for i := uint64(0); i < numSigners; i++ {
		key, _ := crypto.GenerateKey()
		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, nil, nil, nil, err
		}
		auth.GasLimit = uint64(8000000)
		keys[i] = key
		auths[i] = auth
	}

	// Setup a simulated backend
	genesisAlloc := map[common.Address]core.GenesisAccount{}
	for _, auth := range auths {
		genesisAlloc[auth.From] = core.GenesisAccount{Balance: big.NewInt(1e18)}
	}
	blockGasLimit := uint64(8000000)
	sim := backends.NewSimulatedBackend(genesisAlloc, blockGasLimit)

	// Deploy a ManyChainMultiSig contract with any of the signers
	_, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auths[0], sim)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the contract to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Check the receipt status
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, nil, nil, nil, errors.New("contract deployment failed")
	}

	// Set a valid config
	signers := make([]common.Address, numSigners)
	for i, auth := range auths {
		signers[i] = auth.From
	}

	// Set the config
	config := &configwrappers.Config{
		Quorum:       uint8(numSigners),
		Signers:      signers,
		GroupSigners: []configwrappers.Config{},
	}
	quorums, parents, signersAddresses, signerGroups := config.ExtractSetConfigInputs()

	tx, err = mcms.SetConfig(auths[0], signersAddresses, signerGroups, quorums, parents, false)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return keys, auths, sim, mcms, nil
}

func TestExecutor_ExecuteE2E_SingleChainSingleSignerSingleTX_Success(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transaction
	role, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)
	grantRoleData, err := timelockAbi.Pack("grantRole", role, mcms.Address())
	assert.NoError(t, err)

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:    timelock.Address(),
					Value: big.NewInt(0),
					Data:  grantRoleData,
				},
			},
		},
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcms.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// // Check the state of the MCMS contract
	newOpCount, err := mcms.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.NotNil(t, newOpCount)
	assert.Equal(t, uint64(1), newOpCount.Uint64())

	// Check the state of the timelock contract
	proposerCount, err := timelock.GetRoleMemberCount(&bind.CallOpts{}, role)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), proposerCount)
	proposer, err := timelock.GetRoleMember(&bind.CallOpts{}, role, big.NewInt(0))
	assert.NoError(t, err)
	assert.Equal(t, mcms.Address(), proposer)
}

func TestExecutor_ExecuteE2E_SingleChainMultipleSignerSingleTX_Success(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(3)
	assert.NoError(t, err)
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)
	for i := 0; i < 3; i++ {
		assert.NotNil(t, keys[i])
		assert.NotNil(t, auths[i])
	}

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transaction
	role, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)
	grantRoleData, err := timelockAbi.Pack("grantRole", role, mcms.Address())
	assert.NoError(t, err)

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: []ChainOperation{
			{
				ChainIdentifier: TestChain1,
				Operation: Operation{
					To:    timelock.Address(),
					Value: big.NewInt(0),
					Data:  grantRoleData,
				},
			},
		},
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	for i := 0; i < 3; i++ {
		sig, err := crypto.Sign(hash.Bytes(), keys[i])
		assert.NoError(t, err)

		// Construct a signature
		sigObj, err := NewSignatureFromBytes(sig)
		assert.NoError(t, err)
		proposal.Signatures = append(proposal.Signatures, sigObj)
	}

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcms.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check the state of the MCMS contract
	newOpCount, err := mcms.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.NotNil(t, newOpCount)
	assert.Equal(t, uint64(1), newOpCount.Uint64())

	// Check the state of the timelock contract
	proposerCount, err := timelock.GetRoleMemberCount(&bind.CallOpts{}, role)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), proposerCount)
	proposer, err := timelock.GetRoleMember(&bind.CallOpts{}, role, big.NewInt(0))
	assert.NoError(t, err)
	assert.Equal(t, mcms.Address(), proposer)
}

func TestExecutor_ExecuteE2E_SingleChainSingleSignerMultipleTX_Success(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	executorRole, err := timelock.EXECUTORROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]ChainOperation, 4)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		data, err := timelockAbi.Pack("grantRole", role, mcms.Address())
		assert.NoError(t, err)
		operations[i] = ChainOperation{
			ChainIdentifier: TestChain1,
			Operation: Operation{
				To:    timelock.Address(),
				Value: big.NewInt(0),
				Data:  data,
			},
		}
	}

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: operations,
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcms.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	for i := 0; i < 4; i++ {
		// Execute the proposal
		tx, err = executor.ExecuteOnChain(auths[0], i)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		sim.Commit()

		// Wait for the transaction to be mined
		receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
		assert.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// Check the state of the MCMS contract
	newOpCount, err := mcms.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.NotNil(t, newOpCount)
	assert.Equal(t, uint64(4), newOpCount.Uint64())

	// Check the state of the timelock contract
	for _, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		roleCount, err := timelock.GetRoleMemberCount(&bind.CallOpts{}, role)
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1), roleCount)
		roleMember, err := timelock.GetRoleMember(&bind.CallOpts{}, role, big.NewInt(0))
		assert.NoError(t, err)
		assert.Equal(t, mcms.Address(), roleMember)
	}
}

func TestExecutor_ExecuteE2E_SingleChainMultipleSignerMultipleTX_Success(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(3)
	assert.NoError(t, err)
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)
	for i := 0; i < 3; i++ {
		assert.NotNil(t, keys[i])
		assert.NotNil(t, auths[i])
	}

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	executorRole, err := timelock.EXECUTORROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]ChainOperation, 4)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		data, err := timelockAbi.Pack("grantRole", role, mcms.Address())
		assert.NoError(t, err)
		operations[i] = ChainOperation{
			ChainIdentifier: TestChain1,
			Operation: Operation{
				To:    timelock.Address(),
				Value: big.NewInt(0),
				Data:  data,
			},
		}
	}

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: operations,
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	for i := 0; i < 3; i++ {
		sig, err := crypto.Sign(hash.Bytes(), keys[i])
		assert.NoError(t, err)

		// Construct a signature
		sigObj, err := NewSignatureFromBytes(sig)
		assert.NoError(t, err)
		proposal.Signatures = append(proposal.Signatures, sigObj)
	}

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcms.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	for i := 0; i < 4; i++ {
		// Execute the proposal
		tx, err = executor.ExecuteOnChain(auths[0], i)
		assert.NoError(t, err)
		assert.NotNil(t, tx)
		sim.Commit()

		// Wait for the transaction to be mined
		receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
		assert.NoError(t, err)
		assert.NotNil(t, receipt)
		assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
	}

	// Check the state of the MCMS contract
	newOpCount, err := mcms.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.NotNil(t, newOpCount)
	assert.Equal(t, uint64(4), newOpCount.Uint64())

	// Check the state of the timelock contract
	for _, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		roleCount, err := timelock.GetRoleMemberCount(&bind.CallOpts{}, role)
		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1), roleCount)
		roleMember, err := timelock.GetRoleMember(&bind.CallOpts{}, role, big.NewInt(0))
		assert.NoError(t, err)
		assert.Equal(t, mcms.Address(), roleMember)
	}
}

func TestExecutor_ExecuteE2E_SingleChainMultipleSignerMultipleTX_FailureMissingQuorum(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(3)
	assert.NoError(t, err)
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)
	for i := 0; i < 3; i++ {
		assert.NotNil(t, keys[i])
		assert.NotNil(t, auths[i])
	}

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	executorRole, err := timelock.EXECUTORROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]ChainOperation, 4)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		data, err := timelockAbi.Pack("grantRole", role, mcms.Address())
		assert.NoError(t, err)
		operations[i] = ChainOperation{
			ChainIdentifier: TestChain1,
			Operation: Operation{
				To:    timelock.Address(),
				Value: big.NewInt(0),
				Data:  data,
			},
		}
	}

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: operations,
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	for i := 0; i < 2; i++ {
		sig, err := crypto.Sign(hash.Bytes(), keys[i])
		assert.NoError(t, err)

		// Construct a signature
		sigObj, err := NewSignatureFromBytes(sig)
		assert.NoError(t, err)
		proposal.Signatures = append(proposal.Signatures, sigObj)
	}

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.False(t, quorumMet)
	assert.Error(t, err)
	// assert error is of type ErrQuorumNotMet
	assert.IsType(t, &owner_errors.ErrQuorumNotMet{}, err)
}

func TestExecutor_ExecuteE2E_SingleChainMultipleSignerMultipleTX_FailureInvalidSigner(t *testing.T) {
	keys, auths, sim, mcms, err := setupSimulatedBackendWithMCMS(3)
	assert.NoError(t, err)
	assert.NotNil(t, sim)
	assert.NotNil(t, mcms)
	for i := 0; i < 3; i++ {
		assert.NotNil(t, keys[i])
		assert.NotNil(t, auths[i])
	}

	// Generate a new key
	newKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	keys[2] = newKey

	// Deploy a timelock contract for testing
	addr, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		mcms.Address(),
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
		[]common.Address{},
	)
	assert.NoError(t, err)
	assert.NotNil(t, addr)
	assert.NotNil(t, tx)
	assert.NotNil(t, timelock)
	sim.Commit()

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	executorRole, err := timelock.EXECUTORROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]ChainOperation, 4)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole, executorRole} {
		data, err := timelockAbi.Pack("grantRole", role, mcms.Address())
		assert.NoError(t, err)
		operations[i] = ChainOperation{
			ChainIdentifier: TestChain1,
			Operation: Operation{
				To:    timelock.Address(),
				Value: big.NewInt(0),
				Data:  data,
			},
		}
	}

	// Construct a proposal
	proposal := MCMSProposal{
		Version:              "1.0",
		ValidUntil:           2004259681,
		Signatures:           []Signature{},
		OverridePreviousRoot: false,
		ChainMetadata: map[ChainIdentifier]ChainMetadata{
			TestChain1: {
				NonceOffset: 0,
				MCMAddress:  mcms.Address(),
			},
		},
		Transactions: operations,
	}

	// Construct executor
	executor, err := proposal.ToExecutor(map[ChainIdentifier]ContractDeployBackend{TestChain1: sim})
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	for i := 0; i < 3; i++ {
		sig, err := crypto.Sign(hash.Bytes(), keys[i])
		assert.NoError(t, err)

		// Construct a signature
		sigObj, err := NewSignatureFromBytes(sig)
		assert.NoError(t, err)
		proposal.Signatures = append(proposal.Signatures, sigObj)
	}

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures()
	assert.False(t, quorumMet)
	assert.Error(t, err)
	// assert error is of type ErrQuorumNotMet
	assert.IsType(t, &owner_errors.ErrInvalidSignature{}, err)
}
