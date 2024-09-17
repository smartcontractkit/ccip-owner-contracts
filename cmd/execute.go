package cmd

import (
	"log"
	"math/big"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

var ExecuteOperationCmd = &cobra.Command{
	Use:   "execute",
	Short: "Performs an operation execution on the MCMS. Root must be set first.",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		proposal, err := BuildExecutableProposal(proposalPath)
		if err != nil {
			return err
		}

		clientBackend, err := ethclient.Dial(rpc)
		if err != nil {
			return err
		}

		e, err := mcms.NewProposalExecutor(proposal, false)
		if err != nil {
			return err
		}

		uintChainSelector, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			return err
		}
		cs := mcms.ChainIdentifier(uintChainSelector)

		ecdsa, err := crypto.HexToECDSA(pk)
		if err != nil {
			return err
		}

		bigIntChainID := big.NewInt(int64(uintChainSelector))

		for i, op := range e.Proposal.Transactions {
			if op.ChainIdentifier == cs {
				auth, err := bind.NewKeyedTransactorWithChainID(ecdsa, bigIntChainID)
				if err != nil {
					return err
				}

				transaction, err := e.ExecuteOnChain(clientBackend, auth, i)
				if err != nil {
					return err
				}

				log.Printf("Update pending: 0x%x\n", transaction.Hash())
			}
		}

		return nil
	},
}

func init() {
	ExecuteOperationCmd.PersistentFlags().StringVar(&pk, "pk", "", "private key that will send the transactions in the proposal")
}
