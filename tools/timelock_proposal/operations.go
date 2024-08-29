package timelock_proposal

import "github.com/smartcontractkit/ccip-owner-contracts/tools/mcms_proposal"

type BatchChainOperation struct {
	ChainIdentifier string                    `json:"chainIdentifier"`
	Batch           []mcms_proposal.Operation `json:"batch"`
}
