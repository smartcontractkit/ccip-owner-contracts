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

var CheckQuorumCmd = &cobra.Command{
	Use:   "check-quorum",
	Short: "Determines whether the provided signatures meet the quorum to set the root",
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

		quorumMet, err := e.CheckQuorum(clientBackend, auth, mcms.ChainIdentifier(uintChainSelector))
		if err != nil {
			return err
		}

		if quorumMet {
			log.Printf("Signature Quorum met!")
		} else {
			log.Printf("Signature Quorum not met!")
		}

		return nil
	},
}
