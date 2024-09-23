package cmd

import (
	"os"

	"github.com/smartcontractkit/ccip-owner-contracts/tools/proposal"
	"github.com/spf13/cobra"
)

var (
	Rpc           string
	ProposalPath  string
	ProposalType  proposal.ProposalType
	ChainSelector uint64
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mcms",
	Short: "Tools for on-chain interactions with the MCMS ",
	Long:  ``,
}

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(CheckQuorumCmd)
	RootCmd.AddCommand(ExecuteOperationCmd)
	RootCmd.AddCommand(SetMerkleCmd)
	RootCmd.AddCommand(ExecuteChainCmd)

	RootCmd.PersistentFlags().StringVar(&Rpc, "rpc", "", "rpc to be used in the proposal")
	RootCmd.PersistentFlags().StringVar(&ProposalPath, "proposal", "p", "Absolute file path containing the proposal to be submitted")
	RootCmd.PersistentFlags().Uint64Var(&ChainSelector, "selector", 0, "Chain selector for the command to connect to")

	var proposalTypeStr string
	RootCmd.PersistentFlags().StringVar(&proposalTypeStr, "proposalType", string(proposal.MCMS), "The type of proposal being ingested")
	ProposalType = proposal.StringToProposalType[proposalTypeStr]
}
