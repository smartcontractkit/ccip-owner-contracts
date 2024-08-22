package operationdb

import (
	"reflect"
	"testing"

	op "github.com/smartcontractkit/ccip-owner-contracts/internal/executor/operation"
)

func TestBasicMap_AddOps(t *testing.T) {
	om := NewMapdb()

	// Test adding operations
	operations := []Operation{
		{ID: "0x0001_0", Status: NotStarted, Dependencies: []string{}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_1", Status: NotStarted, Dependencies: []string{"0x0001_0"}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_2", Status: NotStarted, Dependencies: []string{"0x0001_1"}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
	}

	err := om.AddOps(operations)
	if err != nil {
		t.Errorf("Failed to add operations: %v", err)
	}

	// Test adding duplicate operation
	duplicateOp := Operation{ID: "0x0001_0", Status: NotStarted, Dependencies: []string{}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil}
	err = om.AddOps([]Operation{duplicateOp})
	if err == nil {
		t.Errorf("Expected error when adding duplicate operation, but got nil")
	}

	// Test adding empty operations
	err = om.AddOps([]Operation{})
	if err != nil {
		t.Errorf("Failed to add empty operations: %v", err)
	}
}

func TestBasicMap_GetOperationByID(t *testing.T) {
	om := NewMapdb()
	// Add operations
	// Test adding operations
	operations := []Operation{
		{ID: "0x0001_0", Status: NotStarted, Dependencies: []string{}, Operation: nil, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_1", Status: NotStarted, Dependencies: []string{"0x0001_0"}, Operation: nil, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_2", Status: NotStarted, Dependencies: []string{"0x0001_1"}, Operation: nil, Timestamp: 0, Receipt: nil},
	}
	err := om.AddOps(operations)
	if err != nil {
		t.Errorf("Failed to add operations: %v", err)
	}

	// Test getting existing operation
	expectedOp := operations[0]
	actualOp, err := om.GetOperationByID(expectedOp.ID)
	if err != nil {
		t.Errorf("Failed to get operation by ID: %v", err)
	}
	if !reflect.DeepEqual(expectedOp, actualOp) {
		t.Errorf("Expected operation %v, but got %v", expectedOp, actualOp)
	}

	// Test getting non-existent operation
	_, err = om.GetOperationByID("invalid")
	if err == nil {
		t.Errorf("Expected error when getting non-existent operation, but got nil")
	}
}

// Add more unit tests for other methods...
func TestBasicMap_GetOperationsByStatus(t *testing.T) {
	om := NewMapdb()
	// Add operations

	operations := []Operation{
		{ID: "0x0001_0", Status: NotStarted, Dependencies: []string{}, Operation: nil, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_1", Status: NotStarted, Dependencies: []string{"0x0001_0"}, Operation: nil, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_2", Status: Failed, Dependencies: []string{"0x0001_1"}, Operation: nil, Timestamp: 0, Receipt: nil},
	}
	err := om.AddOps(operations)
	if err != nil {
		t.Errorf("Failed to add operations: %v", err)
	}

	// Test getting operations by status
	expectedOps := []Operation{{ID: "0x0001_2", Status: Failed, Dependencies: []string{"0x0001_1"}, Operation: nil, Timestamp: 0, Receipt: nil}}
	actualOps, err := om.GetOperationsByStatus(Failed)
	if err != nil {
		t.Errorf("Failed to get operations by status: %v", err)
	}
	if !reflect.DeepEqual(expectedOps, actualOps) {
		t.Errorf("Expected operations %v, but got %v", expectedOps, actualOps)
	}

	actualOps, err = om.GetOperationsByStatus(NotStarted)
	if err != nil {
		t.Errorf("Failed to get operations by status: %v", err)
	}
	if len(actualOps) != 2 {
		t.Errorf("Expected 2 operations, but got %v", len(actualOps))
	}
}

func TestBasicMap_GetOperationStatusByID(t *testing.T) {
	om := NewMapdb()
	// Add operations
	// Test adding operations
	operations := []Operation{
		{ID: "0x0001_0", Status: NotStarted, Dependencies: []string{}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_1", Status: NotStarted, Dependencies: []string{"0x0001_0"}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
		{ID: "0x0001_2", Status: Failed, Dependencies: []string{"0x0001_1"}, Operation: op.EvmOperation{}, Timestamp: 0, Receipt: nil},
	}
	err := om.AddOps(operations)
	if err != nil {
		t.Errorf("Failed to add operations: %v", err)
	}

	// Test getting status of existing operation
	expectedStatus := NotStarted
	actualStatus, err := om.GetOperationStatusByID(operations[0].ID)
	if err != nil {
		t.Errorf("Failed to get operation status by ID: %v", err)
	}
	if expectedStatus != actualStatus {
		t.Errorf("Expected status %v, but got %v", expectedStatus, actualStatus)
	}

	// Test getting status of non-existent operation
	_, err = om.GetOperationStatusByID("invalid")
	if err == nil {
		t.Errorf("Expected error when getting status of non-existent operation, but got nil")
	}
}
