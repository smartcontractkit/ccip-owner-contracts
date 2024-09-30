package proposal

import (
	"errors"

	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
)

type ProposalType string

const (
	// MCMSProposalType is a proposal type for the MCMS contract.
	MCMS ProposalType = "MCMS"
	// MCMSWithTimelock is a proposal type for the MCMS contract with timelock.
	MCMSWithTimelock ProposalType = "MCMSWithTimelock"
)

var StringToProposalType = map[string]ProposalType{
	"MCMS":             MCMS,
	"MCMSWithTimelock": MCMSWithTimelock,
}

type Proposal interface {
	ToExecutor(sim bool) (*mcms.Executor, error)
	AddSignature(signature mcms.Signature)
	Validate() error
}

func LoadProposal(proposalType ProposalType, filePath string) (Proposal, error) {
	switch proposalType {
	case MCMS:
		return mcms.NewProposalFromFile(filePath)
	case MCMSWithTimelock:
		return timelock.NewMCMSWithTimelockProposalFromFile(filePath)
	default:
		return nil, errors.New("unknown proposal type")
	}
}
