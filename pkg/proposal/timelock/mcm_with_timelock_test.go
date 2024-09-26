package timelock

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"
	mcm_errors "github.com/smartcontractkit/ccip-owner-contracts/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/stretchr/testify/assert"
)

var TestAddress = common.HexToAddress("0x1234567890abcdef")
var TestChain1 = mcms.ChainIdentifier(3379446385462418246)
var TestChain2 = mcms.ChainIdentifier(16015286601757825753)
var TestChain3 = mcms.ChainIdentifier(10344971235874465080)

func TestValidate_ValidProposal(t *testing.T) {
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: TestAddress,
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:           TestAddress,
						Value:        big.NewInt(0),
						Data:         common.Hex2Bytes("0x"),
						ContractType: "Sample contract",
						Tags:         []string{"tag1", "tag2"},
					},
				},
			},
		},
		Schedule,
		"1h",
	)

	assert.NoError(t, err)
	assert.NotNil(t, proposal)
}

func TestValidate_InvalidOperation(t *testing.T) {
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: TestAddress,
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:           TestAddress,
						Value:        big.NewInt(0),
						Data:         common.Hex2Bytes("0x"),
						ContractType: "Sample contract",
						Tags:         []string{"tag1", "tag2"},
					},
				},
			},
		},
		"invalid",
		"1h",
	)

	assert.Error(t, err)
	assert.Nil(t, proposal)
	assert.IsType(t, &mcm_errors.ErrInvalidTimelockOperation{}, err)
}

func TestValidate_InvalidMinDelaySchedule(t *testing.T) {
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: TestAddress,
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:           TestAddress,
						Value:        big.NewInt(0),
						Data:         common.Hex2Bytes("0x"),
						ContractType: "Sample contract",
						Tags:         []string{"tag1", "tag2"},
					},
				},
			},
		},
		Schedule,
		"invalid",
	)

	assert.Error(t, err)
	assert.Nil(t, proposal)
	assert.Equal(t, err.Error(), "time: invalid duration \"invalid\"")
}

func TestValidate_InvalidMinDelayBypassShouldBeValid(t *testing.T) {
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 1,
				MCMAddress:      TestAddress,
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: TestAddress,
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:           TestAddress,
						Value:        big.NewInt(0),
						Data:         common.Hex2Bytes("0x"),
						ContractType: "Sample contract",
						Tags:         []string{"tag1", "tag2"},
					},
				},
			},
		},
		Bypass,
		"invalid",
	)

	assert.NoError(t, err)
	assert.NotNil(t, proposal)
}

// Constructs a simulated backend with a ManyChainMultiSig contract and a RBACTimelock contract
// The Admin of the RBACTimelock is itself and the RBACTimelock owns the ManyChainMultiSig
func setupSimulatedBackendWithMCMSAndTimelock(numSigners uint64) ([]*ecdsa.PrivateKey, []*bind.TransactOpts, *backends.SimulatedBackend, *gethwrappers.ManyChainMultiSig, *gethwrappers.RBACTimelock, error) {
	// Generate a private key
	keys := make([]*ecdsa.PrivateKey, numSigners)
	auths := make([]*bind.TransactOpts, numSigners)
	for i := uint64(0); i < numSigners; i++ {
		key, _ := crypto.GenerateKey()
		auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
		if err != nil {
			return nil, nil, nil, nil, nil, err
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
	mcmAddr, tx, mcms, err := gethwrappers.DeployManyChainMultiSig(auths[0], sim)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the contract to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Check the receipt status
	if receipt.Status != types.ReceiptStatusSuccessful {
		return nil, nil, nil, nil, nil, errors.New("contract deployment failed")
	}

	// Set a valid config
	signers := make([]common.Address, numSigners)
	for i, auth := range auths {
		signers[i] = auth.From
	}

	// Set the config
	config := &config.Config{
		Quorum:       uint8(numSigners),
		Signers:      signers,
		GroupSigners: []config.Config{},
	}
	quorums, parents, signersAddresses, signerGroups := config.ExtractSetConfigInputs()

	tx, err = mcms.SetConfig(auths[0], signersAddresses, signerGroups, quorums, parents, false)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Deploy a timelock contract for testing
	_, tx, timelock, err := gethwrappers.DeployRBACTimelock(
		auths[0],
		sim,
		big.NewInt(0),
		auths[0].From, // Temporarily set the admin to the first signer
		[]common.Address{mcmAddr},
		[]common.Address{mcmAddr, auths[0].From},
		[]common.Address{mcmAddr},
		[]common.Address{mcmAddr},
	)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the contract to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Transfer the ownership of the ManyChainMultiSig to the timelock
	tx, err = mcms.TransferOwnership(auths[0], timelock.Address())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Construct payload for Accepting the ownership of the ManyChainMultiSig
	mcmsAbi, err := gethwrappers.ManyChainMultiSigMetaData.GetAbi()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	acceptOwnershipData, err := mcmsAbi.Pack("acceptOwnership")
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Accept the ownership of the ManyChainMultiSig
	tx, err = timelock.BypasserExecuteBatch(auths[0], []gethwrappers.RBACTimelockCall{
		{
			Target: mcms.Address(),
			Data:   acceptOwnershipData,
			Value:  big.NewInt(0),
		},
	})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Give the timelock admin rights
	role, err := timelock.ADMINROLE(&bind.CallOpts{})
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	tx, err = timelock.GrantRole(auths[0], role, timelock.Address())
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Revoking the admin rights of the first signer
	tx, err = timelock.RevokeRole(auths[0], role, auths[0].From)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Mine a block
	sim.Commit()

	// Wait for the transaction to be mined
	_, err = bind.WaitMined(auths[0].Context, sim, tx)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return keys, auths, sim, mcms, timelock, nil
}

func TestE2E_ValidScheduleAndExecuteProposalOneTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transaction to grant EOA the PROPOSER role
	role, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)
	grantRoleData, err := timelockAbi.Pack("grantRole", role, auths[0].From)
	assert.NoError(t, err)

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, role, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:    timelock.Address(),
						Value: big.NewInt(0),
						Data:  grantRoleData,
					},
				},
			},
		},
		Schedule,
		"5s",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check all the logs
	var operationId common.Hash
	for _, log := range receipt.Logs {
		event, err := timelock.ParseCallScheduled(*log)
		if err == nil {
			operationId = event.Id
		}
	}

	// Validate Contract State and verify operation was scheduled
	grantRoleCall := []gethwrappers.RBACTimelockCall{
		{
			Target: timelock.Address(),
			Value:  big.NewInt(0),
			Data:   grantRoleData,
		},
	}

	isOperation, err := timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperation)
	isOperationPending, err := timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationPending)
	isOperationReady, err := timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)

	// sleep for 5 seconds and then mine a block
	time.Sleep(5 * time.Second)
	sim.Commit()

	// Check that the operation is now ready
	isOperationReady, err = timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationReady)

	// Execute the operation
	tx, err = timelock.ExecuteBatch(auths[0], grantRoleCall, ZERO_HASH, ZERO_HASH)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err = bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check that the operation is done
	isOperationDone, err := timelock.IsOperationDone(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationDone)

	// Check that the operation is no longer pending
	isOperationPending, err = timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationPending)

	// Validate Contract State and verify role was granted
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, role, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
}

func TestE2E_ValidScheduleAndCancelProposalOneTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transaction to grant EOA the PROPOSER role
	role, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)
	grantRoleData, err := timelockAbi.Pack("grantRole", role, auths[0].From)
	assert.NoError(t, err)

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, role, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:    timelock.Address(),
						Value: big.NewInt(0),
						Data:  grantRoleData,
					},
				},
			},
		},
		Schedule,
		"5s",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check all the logs
	var operationId common.Hash
	for _, log := range receipt.Logs {
		event, err := timelock.ParseCallScheduled(*log)
		if err == nil {
			operationId = event.Id
		}
	}

	// Check operation state and see that it was scheduled
	isOperation, err := timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperation)
	isOperationPending, err := timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationPending)
	isOperationReady, err := timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)

	// Get and validate the current operation count
	currOpCount, err := mcmsObj.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, currOpCount.Int64(), int64(len(proposal.Transactions)))

	// Generate a new proposal to cancel the operation
	// Update the proposal Operation to Cancel
	// Update the proposal ChainMetadata StartingOpCount to the current operation count
	proposal.Operation = Cancel
	proposal.ChainMetadata[TestChain1] = mcms.ChainMetadata{
		StartingOpCount: currOpCount.Uint64(),
		MCMAddress:      mcmsObj.Address(),
	}

	// Construct executor
	executor, err = proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err = executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err = crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err = mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err = executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err = mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err = bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Verify operation state and confirm it was cancelled
	isOperation, err = timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperation)
	isOperationPending, err = timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationPending)
	isOperationReady, err = timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)
}

func TestE2E_ValidBypassProposalOneTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transaction to grant EOA the PROPOSER role
	role, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)
	grantRoleData, err := timelockAbi.Pack("grantRole", role, auths[0].From)
	assert.NoError(t, err)

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, role, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch: []mcms.Operation{
					{
						To:    timelock.Address(),
						Value: big.NewInt(0),
						Data:  grantRoleData,
					},
				},
			},
		},
		Bypass,
		"",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Validate Contract State and verify role was granted
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, role, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
}

func TestE2E_ValidScheduleAndExecuteProposalOneBatchTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]mcms.Operation, 3)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole} {
		data, err := timelockAbi.Pack("grantRole", role, auths[0].From)
		assert.NoError(t, err)
		operations[i] = mcms.Operation{
			To:    timelock.Address(),
			Value: big.NewInt(0),
			Data:  data,
		}
	}

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, proposerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, bypasserRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, cancellerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch:           operations,
			},
		},
		Schedule,
		"5s",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check all the logs
	var operationId common.Hash
	for _, log := range receipt.Logs {
		event, err := timelock.ParseCallScheduled(*log)
		if err == nil {
			operationId = event.Id
		}
	}

	// Validate Contract State and verify operation was scheduled
	grantRoleCalls := []gethwrappers.RBACTimelockCall{
		{
			Target: timelock.Address(),
			Value:  big.NewInt(0),
			Data:   operations[0].Data,
		},
		{
			Target: timelock.Address(),
			Value:  big.NewInt(0),
			Data:   operations[1].Data,
		},
		{
			Target: timelock.Address(),
			Value:  big.NewInt(0),
			Data:   operations[2].Data,
		},
	}

	isOperation, err := timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperation)
	isOperationPending, err := timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationPending)
	isOperationReady, err := timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)

	// sleep for 5 seconds and then mine a block
	time.Sleep(5 * time.Second)
	sim.Commit()

	// Check that the operation is now ready
	isOperationReady, err = timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationReady)

	// Execute the operation
	tx, err = timelock.ExecuteBatch(auths[0], grantRoleCalls, ZERO_HASH, ZERO_HASH)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err = bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check that the operation is done
	isOperationDone, err := timelock.IsOperationDone(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationDone)

	// Check that the operation is no longer pending
	isOperationPending, err = timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationPending)

	// Validate Contract State and verify role was granted
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, proposerRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, bypasserRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, cancellerRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
}

func TestE2E_ValidScheduleAndCancelProposalOneBatchTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]mcms.Operation, 3)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole} {
		data, err := timelockAbi.Pack("grantRole", role, auths[0].From)
		assert.NoError(t, err)
		operations[i] = mcms.Operation{
			To:    timelock.Address(),
			Value: big.NewInt(0),
			Data:  data,
		}
	}

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, proposerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, bypasserRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, cancellerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch:           operations,
			},
		},
		Schedule,
		"5s",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Check all the logs
	var operationId common.Hash
	for _, log := range receipt.Logs {
		event, err := timelock.ParseCallScheduled(*log)
		if err == nil {
			operationId = event.Id
		}
	}

	// Check operation state and see that it was scheduled
	isOperation, err := timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperation)
	isOperationPending, err := timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.True(t, isOperationPending)
	isOperationReady, err := timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)

	// Get and validate the current operation count
	currOpCount, err := mcmsObj.GetOpCount(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, currOpCount.Int64(), int64(len(proposal.Transactions)))

	// Generate a new proposal to cancel the operation
	// Update the proposal Operation to Cancel
	// Update the proposal ChainMetadata StartingOpCount to the current operation count
	proposal.Operation = Cancel
	proposal.ChainMetadata[TestChain1] = mcms.ChainMetadata{
		StartingOpCount: currOpCount.Uint64(),
		MCMAddress:      mcmsObj.Address(),
	}

	// Construct executor
	executor, err = proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err = executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err = crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err = mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err = executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err = executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err = mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err = bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Verify operation state and confirm it was cancelled
	isOperation, err = timelock.IsOperation(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperation)
	isOperationPending, err = timelock.IsOperationPending(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationPending)
	isOperationReady, err = timelock.IsOperationReady(&bind.CallOpts{}, operationId)
	assert.NoError(t, err)
	assert.False(t, isOperationReady)
}

func TestE2E_ValidBypassProposalOneBatchTx(t *testing.T) {
	keys, auths, sim, mcmsObj, timelock, err := setupSimulatedBackendWithMCMSAndTimelock(1)
	assert.NoError(t, err)
	assert.NotNil(t, keys[0])
	assert.NotNil(t, auths[0])
	assert.NotNil(t, sim)
	assert.NotNil(t, mcmsObj)
	assert.NotNil(t, timelock)

	// Construct example transactions
	proposerRole, err := timelock.PROPOSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	bypasserRole, err := timelock.BYPASSERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	cancellerRole, err := timelock.CANCELLERROLE(&bind.CallOpts{})
	assert.NoError(t, err)
	timelockAbi, err := gethwrappers.RBACTimelockMetaData.GetAbi()
	assert.NoError(t, err)

	operations := make([]mcms.Operation, 3)
	for i, role := range []common.Hash{proposerRole, bypasserRole, cancellerRole} {
		data, err := timelockAbi.Pack("grantRole", role, auths[0].From)
		assert.NoError(t, err)
		operations[i] = mcms.Operation{
			To:    timelock.Address(),
			Value: big.NewInt(0),
			Data:  data,
		}
	}

	// Validate Contract State and verify role does not exist
	hasRole, err := timelock.HasRole(&bind.CallOpts{}, proposerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, bypasserRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, cancellerRole, auths[0].From)
	assert.NoError(t, err)
	assert.False(t, hasRole)

	// Construct example transaction
	proposal, err := NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			TestChain1: {
				StartingOpCount: 0,
				MCMAddress:      mcmsObj.Address(),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			TestChain1: timelock.Address(),
		},
		"Sample description",
		[]BatchChainOperation{
			{
				ChainIdentifier: TestChain1,
				Batch:           operations,
			},
		},
		Bypass,
		"",
	)
	assert.NoError(t, err)
	assert.NotNil(t, proposal)

	// Gen caller map for easy access
	callers := map[mcms.ChainIdentifier]mcms.ContractDeployBackend{TestChain1: sim}

	// Construct executor
	executor, err := proposal.ToExecutor(true)
	assert.NoError(t, err)
	assert.NotNil(t, executor)

	// Get the hash to sign
	hash, err := executor.SigningHash()
	assert.NoError(t, err)

	// Sign the hash
	sig, err := crypto.Sign(hash.Bytes(), keys[0])
	assert.NoError(t, err)

	// Construct a signature
	sigObj, err := mcms.NewSignatureFromBytes(sig)
	assert.NoError(t, err)
	executor.Proposal.Signatures = append(proposal.Signatures, sigObj)

	// Validate the signatures
	quorumMet, err := executor.ValidateSignatures(callers)
	assert.True(t, quorumMet)
	assert.NoError(t, err)

	// SetRoot on the contract
	tx, err := executor.SetRootOnChain(sim, auths[0], TestChain1)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Validate Contract State and verify root was set
	root, err := mcmsObj.GetRoot(&bind.CallOpts{})
	assert.NoError(t, err)
	assert.Equal(t, root.Root, [32]byte(executor.Tree.Root.Bytes()))
	assert.Equal(t, root.ValidUntil, proposal.ValidUntil)

	// Execute the proposal
	tx, err = executor.ExecuteOnChain(sim, auths[0], 0)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	sim.Commit()

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(auths[0].Context, sim, tx)
	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	// Validate Contract State and verify role was granted
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, proposerRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, bypasserRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
	hasRole, err = timelock.HasRole(&bind.CallOpts{}, cancellerRole, auths[0].From)
	assert.NoError(t, err)
	assert.True(t, hasRole)
}
