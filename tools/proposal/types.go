package proposal

import "github.com/smartcontractkit/ccip-owner-contracts/tools/proposal/mcms"

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
	ToExecutor() (*mcms.Executor, error)
	AddSignature(signature mcms.Signature)
	Validate() error
}
