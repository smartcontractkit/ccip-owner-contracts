package operationdb

import (
	"errors"
	"strings"
	"sync"
)

type BasicMap struct {
	mu         sync.Mutex
	operations map[string]Operation
}

func NewMapdb() *BasicMap {
	return &BasicMap{
		operations: make(map[string]Operation),
	}
}

func (om *BasicMap) AddOps(operations []Operation) error {
	failedOps := make([]string, 0)
	for _, o := range operations {
		// if the key exists in the operations BasicMap, then the operation is already present
		// [M2]TODO: Update Op in the BasicMap if status is different
		_, exists := om.operations[o.ID]
		if exists {
			failedOps = append(failedOps, o.ID)
		}
		om.operations[o.ID] = o
	}

	if len(failedOps) == 0 {
		return nil
	}
	return errors.New("Failed to add operations: " + strings.Join(failedOps, ", "))
}

func (om *BasicMap) GetOperationByID(id string) (Operation, error) {
	om.mu.Lock()
	op, ok := om.operations[id]
	if !ok {
		return Operation{}, errors.New("Operation not found")
	}
	om.mu.Unlock()
	return op, nil
}

func (om *BasicMap) WriteToStore(o Operation) (bool, error) {
	om.mu.Lock()
	om.operations[o.ID] = o
	om.mu.Unlock()
	return true, nil
}

func (om *BasicMap) GetOperationsByStatus(s Status) ([]Operation, error) {
	om.mu.Lock()
	ops := make([]Operation, 0)
	for _, o := range om.operations {
		if o.Status == s {
			ops = append(ops, o)
		}
	}
	om.mu.Unlock()
	return ops, nil
}

func (om *BasicMap) GetOperationStatusByID(opID string) (Status, error) {
	om.mu.Lock()
	if _, ok := om.operations[opID]; !ok {
		return Unknown, errors.New("Operation not found")
	}
	o := om.operations[opID]
	om.mu.Unlock()
	return o.Status, nil
}
