package scheduler

import (
	"context"
	"time"

	"github.com/smartcontractkit/ccip-owner-contracts/internal/executor/operationdb"
	w "github.com/smartcontractkit/ccip-owner-contracts/internal/executor/worker"
)

type Scheduler struct {
	WorkerCount   int
	OpStore       operationdb.Store
	WorkerTimeout time.Duration
}

func NewScheduler(workerCount int, opStore operationdb.Store) *Scheduler {
	return &Scheduler{
		WorkerCount: workerCount,
		OpStore:     opStore,
	}
}

// Step 1: Get all Ops that are not failed or DID_NOT_START
// Step 2: Iterate over all Ops forever
// Step 3: For each Op, check status of dependencies
// Step 4: If all dependencies are SUCCESS, then push job to channel
// Step 5: Start worker to read from channels and execute Op
func (s *Scheduler) Run(ctx context.Context) (bool, error) {
	// Step 1
	ops, err := s.OpStore.GetOperationsByStatus(operationdb.NotStarted)
	if err != nil {
		return false, err
	}

	// Step 2
	opsCounter := 0
	for {
		// reset counter if we've hit limit
		if opsCounter == len(ops) {
			opsCounter = 0
		}

		// the current operation we're considering
		operation := ops[opsCounter]

		// Step 3
		dependencies := operation.Dependencies
		allDepsSuccess := true
		for _, dep := range dependencies {
			status, err := s.OpStore.GetOperationStatusByID(dep)
			if err != nil {
				return false, err
			}
			if status != operationdb.Success {
				allDepsSuccess = false
				break
			}
		}

		// Step 4
		// we don't use a channel here, because we could eventually just trigger lambda
		// should we use a channel?
		if allDepsSuccess {
			go w.SendTxn(ctx, operation, s.OpStore)
		}
		time.Sleep(time.Second)
	}
}

func (s *Scheduler) UpdateWorkerTimeout(timeout time.Duration) {
	s.WorkerTimeout = timeout
}
