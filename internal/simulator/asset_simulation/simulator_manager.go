package asset_simulation

import (
	"context"
	"sync"

	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Runner interface {
	Start(ctx context.Context) error
	Stop() error
	GetId() string
	IsRunning() bool
}

// AssetSimulatorManager manages runners for asset simulation.
// It supports adding, removing, starting and stopping runners in a thread-safe way at runtime.
type AssetSimulatorManager struct {
	wg      sync.WaitGroup
	mu      sync.Mutex
	obs     observability.Observability
	workers map[string]Runner
}

func NewAssetSimulatorManager(obs observability.Observability) *AssetSimulatorManager {
	return &AssetSimulatorManager{
		obs:     obs,
		workers: make(map[string]Runner),
	}
}

// AddAndStartWorker adds a worker to the manager and starts it
func (wm *AssetSimulatorManager) AddAndStartWorker(ctx context.Context, worker Runner) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.obs.Log().Debug("Adding worker", zap.String("workerId", worker.GetId()))

	// Stop the worker and remove it
	previousWorker, ok := wm.workers[worker.GetId()]
	if ok && previousWorker.IsRunning() {
		err := previousWorker.Stop()
		if err != nil {
			wm.obs.Log().With(zap.Error(err)).Error("Unable to stop worker")
		}
	}

	// Start worker in a goroutine
	wm.workers[worker.GetId()] = worker
	wm.wg.Add(1)
	go func() {
		defer wm.wg.Done()

		wm.obs.Log().Debug("Starting worker", zap.String("workerId", worker.GetId()))
		err := worker.Start(ctx)
		if err != nil {
			wm.obs.Log().With(zap.Error(err)).Error("Unable to start worker")
		}
	}()
}

// AddWorker adds a worker to the manager
func (wm *AssetSimulatorManager) AddWorker(worker Runner) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.obs.Log().Debug("Adding worker", zap.String("workerId", worker.GetId()))

	wm.workers[worker.GetId()] = worker
}

// RemoveWorker removes a worker by its ID
func (wm *AssetSimulatorManager) RemoveWorker(workerId string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.obs.Log().Debug("Removing worker", zap.String("workerId", workerId))

	// Stop the worker before removing it
	worker, ok := wm.workers[workerId]
	if ok && worker.IsRunning() {
		err := worker.Stop()
		if err != nil {
			wm.obs.Log().With(zap.Error(err)).Error("Unable to stop worker")
		}
	}

	delete(wm.workers, workerId)
}

// GetWorker returns a worker by its ID
func (wm *AssetSimulatorManager) GetWorker(workerId string) (Runner, bool) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	worker, ok := wm.workers[workerId]
	return worker, ok
}

// GetWorkers returns all workers
func (wm *AssetSimulatorManager) GetWorkers() map[string]Runner {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.workers
}

// StartWorkers starts all workers, each worker in a separate goroutine
func (wm *AssetSimulatorManager) StartWorkers(ctx context.Context) {
	workers := wm.GetWorkers()
	wm.wg.Add(len(workers))

	wm.obs.Log().Info("Starting workers", zap.Int("count", len(workers)))

	for _, worker := range workers {
		// Start worker in a goroutine
		go func() {
			defer wm.wg.Done()

			wm.obs.Log().Debug("Starting worker", zap.String("workerId", worker.GetId()))
			err := worker.Start(ctx)
			if err != nil {
				wm.obs.Log().With(zap.Error(err)).Error("Unable to start worker")
			}
		}()
	}
}

// StopAll stops all workers
func (wm *AssetSimulatorManager) StopAll() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.obs.Log().Info("Stopping all workers", zap.Int("count", len(wm.workers)))

	for _, worker := range wm.workers {
		err := worker.Stop()
		if err != nil {
			wm.obs.Log().With(zap.Error(err)).Error("Unable to stop worker")
			continue
		}
	}

	wm.wg.Wait()
}

// RemoveAll stops and removes all workers
func (wm *AssetSimulatorManager) RemoveAll() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.obs.Log().Info("Stopping and removing all workers", zap.Int("count", len(wm.workers)))

	for _, worker := range wm.workers {
		wm.RemoveWorker(worker.GetId())
	}

	wm.wg.Wait()
}
