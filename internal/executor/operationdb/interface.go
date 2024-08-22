package operationdb

import (
	o "github.com/smartcontractkit/ccip-owner-contracts/internal/executor/operation"
)

type Status int

const (
	NotStarted  Status = iota
	Executing   Status = iota
	Success     Status = iota
	Failed      Status = iota
	DidNotStart Status = iota
	Unknown     Status = iota
)

// can add boilerplate getters to make safer
type Operation struct {
	ID           string
	Status       Status
	Dependencies []string
	Operation    o.Operation
	Timestamp    int64 // unix timestamp
	Receipt      interface{}
}

type Store interface {
	AddOps([]Operation) error
	GetOperationByID(string) (Operation, error)
	WriteToStore(Operation) (bool, error)
	GetOperationsByStatus(Status) ([]Operation, error)
	GetOperationStatusByID(string) (Status, error)
	UpdateTimestamp() error
}
