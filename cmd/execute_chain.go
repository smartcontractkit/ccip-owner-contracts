package cmd

import (
	"errors"
	"log"
	"math/big"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ExecuteChainCmd = &cobra.Command{
	Use:   "execute-chain",
	Short: "Executes all operations for a given chain in an MCMS Proposal. Root must be set first.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get private key
		pk, err := LoadPrivateKey()
		if err != nil {
			return err
		}

		// Load proposal
		proposal, err := LoadProposal(ProposalType, ProposalPath)
		if err != nil {
			return err
		}

		// Dial the RPC
		clientBackend, err := ethclient.Dial(Rpc)
		if err != nil {
			return err
		}

		// Convert proposal to executor
		e, err := proposal.ToExecutor(false)
		if err != nil {
			return err
		}

		// Get EVM chain ID
		chain, exists := chain_selectors.ChainBySelector(ChainSelector)
		if !exists {
			return errors.New("chain not found")
		}

		// Create a new transactor
		auth, err := bind.NewKeyedTransactorWithChainID(pk, big.NewInt(int64(chain.EvmChainID)))
		if err != nil {
			return err
		}

		for i, op := range e.Proposal.Transactions {
			if op.ChainIdentifier == mcms.ChainIdentifier(ChainSelector) {
				transaction, err := e.ExecuteOnChain(clientBackend, auth, i)
				if err != nil {
					return err
				}

				log.Printf("Transaction sent: %s", transaction.Hash().Hex())

				// Wait for transaction to be mined
				receipt, err := bind.WaitMined(auth.Context, clientBackend, transaction)
				if err != nil {
					return err
				}

				// Check if the transaction was successful
				if receipt.Status != types.ReceiptStatusSuccessful {
					return errors.New("transaction failed")
				}

				log.Printf("Transaction mined: %s", receipt.TxHash.Hex())
			}
		}

		return nil
	},
}
