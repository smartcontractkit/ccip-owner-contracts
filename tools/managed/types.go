package managed

import "github.com/smartcontractkit/ccip-owner-contracts/tools/executable"

type MCMSProposal interface {
	// TODO: maybe enforce a validate function here
	ToExecutableMCMSProposal() (executable.ExecutableMCMSProposal, error)
	AddSignature(sig executable.Signature) error
	Validate() error
}

type MCMSProposalType string

const (
	MCMSOnly          MCMSProposalType = "mcms-only"
	MCMSWithTimelock  MCMSProposalType = "mcms-with-timelock"
	MCMSWithMultisend MCMSProposalType = "mcms-with-multisend"
)

var MCMSProposalTypeMap = map[string]MCMSProposalType{
	"mcms-only":           MCMSOnly,
	"mcms-with-timelock":  MCMSWithTimelock,
	"mcms-with-multisend": MCMSWithMultisend,
}
