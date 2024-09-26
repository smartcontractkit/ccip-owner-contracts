package main

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/timelock"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
)

func TestGenerateTestProposal(t *testing.T) {
	testChain := mcms.ChainIdentifier(chain_selectors.TEST_1000.Selector)
	tp, err := timelock.NewMCMSWithTimelockProposal(
		"1.0",
		2004259681,
		[]mcms.Signature{},
		false,
		map[mcms.ChainIdentifier]mcms.ChainMetadata{
			testChain: {
				StartingOpCount: 0,
				MCMAddress:      common.HexToAddress("0xabc"),
			},
		},
		map[mcms.ChainIdentifier]common.Address{
			testChain: common.HexToAddress("0xabd"),
		},
		"Sample description",
		[]timelock.BatchChainOperation{
			{
				ChainIdentifier: testChain,
				Batch: []mcms.Operation{
					{
						To:    common.HexToAddress("0xabd"),
						Value: big.NewInt(0),
						Data:  []byte("hello"),
					},
				},
			},
		},
		timelock.Schedule,
		"5s",
	)
	require.NoError(t, err)
	require.NoError(t, WriteProposalToFile(tp, "test.json"))
	proposal, err := proposal.LoadProposal(proposal.MCMSWithTimelock, "test.json")
	require.NoError(t, err)

	// Parse the derivation path
	path, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	require.NoError(t, err)

	require.NoError(t, SignLedger(path, proposal))
}
