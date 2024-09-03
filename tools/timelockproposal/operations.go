package timelockproposal

import "github.com/smartcontractkit/ccip-owner-contracts/tools/mcmsproposal"

type BatchChainOperation struct {
	ChainIdentifier mcmsproposal.ChainIdentifier `json:"chainIdentifier"`
	Batch           []mcmsproposal.Operation     `json:"batch"`
}
