// Package gethwrappers keeps track of the golang wrappers of the solidity contracts
package gethwrappers

import (
	"github.com/ethereum/go-ethereum/common"
)

//go:generate ./generate.sh

// AbigenLog is an interface for abigen generated log topics
type AbigenLog interface {
	Topic() common.Hash
}
