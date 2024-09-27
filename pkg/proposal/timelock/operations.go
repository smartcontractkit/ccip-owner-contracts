package timelock

import "github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"

type BatchChainOperation struct {
	ChainIdentifier mcms.ChainIdentifier `json:"chainIdentifier"`
	Batch           []mcms.Operation     `json:"batch"`
}
