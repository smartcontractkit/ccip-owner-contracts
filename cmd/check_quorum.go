package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"
)

var CheckQuorumCmd = &cobra.Command{
	Use:   "check-quorum",
	Short: "Determines whether the provided signatures meet the quorum to set the root",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load proposal
		proposal, err := LoadProposal(proposalType, proposalPath)
		if err != nil {
			return err
		}

		// Dial the RPC
		clientBackend, err := ethclient.Dial(rpc)
		if err != nil {
			return err
		}

		// Convert proposal to executor
		e, err := proposal.ToExecutor()
		if err != nil {
			return err
		}

		// Check quorum
		quorumMet, err := e.CheckQuorum(clientBackend, mcms.ChainIdentifier(chainSelector))
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
