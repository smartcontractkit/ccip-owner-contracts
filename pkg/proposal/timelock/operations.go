package timelock

import "github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"

type BatchChainOperation struct {
	ChainIdentifier mcms.ChainIdentifier `json:"chainIdentifier"`
	Salt            [32]byte             `json:"salt"`
	Batch           []mcms.Operation     `json:"batch"`
}
