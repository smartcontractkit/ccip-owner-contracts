package cmd

import (
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/spf13/cobra"
)

var rpc string
var proposalPath string
var proposalType proposal.ProposalType
var chainSelector uint64

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mcms",
	Short: "Tools for on-chain interactions with the MCMS ",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(CheckQuorumCmd)
	rootCmd.AddCommand(ExecuteOperationCmd)
	rootCmd.AddCommand(SetMerkleCmd)
	rootCmd.AddCommand(ExecuteChainCmd)

	rootCmd.PersistentFlags().StringVar(&rpc, "rpc", "", "rpc to be used in the proposal")
	rootCmd.PersistentFlags().StringVar(&proposalPath, "proposal", "p", "Absolute file path containing the proposal to be submitted")
	rootCmd.PersistentFlags().Uint64Var(&chainSelector, "selector", 0, "Chain selector for the command to connect to")

	var proposalTypeStr string
	rootCmd.PersistentFlags().StringVar(&proposalTypeStr, "proposalType", string(proposal.MCMS), "The type of proposal being ingested")
	proposalType = proposal.StringToProposalType[proposalTypeStr]
}
