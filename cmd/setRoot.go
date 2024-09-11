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

var SetMerkleCmd = &cobra.Command{
	Use:   "set-root",
	Short: "Sets the Merkle Root on the MCM Contract",
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

		e, err := mcms.NewProposalExecutor(proposal)
		if err != nil {
			return err
		}

		uintChainSelector, err := strconv.ParseUint(chainSelector, 10, 64)
		if err != nil {
			return err
		}

		ecdsa, err := crypto.HexToECDSA(pk)
		if err != nil {
			return err
		}

		bigIntChainID := big.NewInt(int64(uintChainSelector))

		auth, err := bind.NewKeyedTransactorWithChainID(ecdsa, bigIntChainID)
		if err != nil {
			return err
		}

		transaction, err := e.SetRootOnChain(clientBackend, auth, mcms.ChainIdentifier(uintChainSelector))
		if err != nil {
			return err
		}

		log.Printf("Transaction sent: %s", transaction.Hash().Hex())
		return nil
	},
}

func init() {
	SetMerkleCmd.PersistentFlags().StringVar(&pk, "pk", "", "private key that will send the transactions in the proposal")
}
