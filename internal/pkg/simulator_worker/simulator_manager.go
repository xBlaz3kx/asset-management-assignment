package simulator_worker

import (
	"context"
	"sync"
)

type AssetSimulator interface {
	Start(ctx context.Context) error
	Stop() error
	GetId() string
}

type AssetSimulatorManager struct {
	wg      sync.WaitGroup
	mu      sync.Mutex
	workers map[string]AssetSimulator
}

func NewAssetSimulatorManager() *AssetSimulatorManager {
	return &AssetSimulatorManager{
		workers: make(map[string]AssetSimulator),
	}
}

// AddWorker adds a worker to the manager
func (wm *AssetSimulatorManager) AddWorker(worker AssetSimulator) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.workers[worker.GetId()] = worker
}

// RemoveWorker removes a worker by its ID
func (wm *AssetSimulatorManager) RemoveWorker(workerId string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	delete(wm.workers, workerId)
}

// GetWorker returns a worker by its ID
func (wm *AssetSimulatorManager) GetWorker(workerId string) (AssetSimulator, bool) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	worker, ok := wm.workers[workerId]
	return worker, ok
}

// GetWorkers returns all workers
func (wm *AssetSimulatorManager) GetWorkers() map[string]AssetSimulator {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.workers
}

// StartWorkers starts all workers, each worker in a separate goroutine
func (wm *AssetSimulatorManager) StartWorkers(ctx context.Context) {
	workers := wm.GetWorkers()
	wm.wg.Add(len(workers))

	for _, worker := range workers {
		// Start worker in a goroutine
		go func() {
			defer wm.wg.Done()
			err := worker.Start(ctx)
			if err != nil {
				return
			}
		}()
	}
}

// StopAll stops all workers
func (wm *AssetSimulatorManager) StopAll() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for _, worker := range wm.workers {
		err := worker.Stop()
		if err != nil {
			continue
		}
	}

	wm.wg.Wait()
}
