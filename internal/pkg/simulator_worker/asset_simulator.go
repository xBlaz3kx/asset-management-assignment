package simulator_worker

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, measurement measurements.Measurement, assetId string) error
}

type simpleAssetSimulator struct {
	obs          observability.Observability
	simulatorCfg simulator.Configuration
	stopChan     chan bool
	publisher    Publisher
}

func NewSimpleAssetSimulator(obs observability.Observability, configuration simulator.Configuration, publisher Publisher) AssetSimulator {
	return &simpleAssetSimulator{
		obs:          obs,
		stopChan:     make(chan bool),
		simulatorCfg: configuration,
		publisher:    publisher,
	}
}

// Start creates a ticker and generates a measurement at each tick.
// The measurement is then published via a Publisher.
func (s *simpleAssetSimulator) Start(ctx context.Context) error {
	s.obs.Log().Debug("Starting simulator", zap.Any("config", s.simulatorCfg))
	ticker := time.NewTicker(s.simulatorCfg.MeasurementInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return nil
		case <-ticker.C:
			s.publishMessage(ctx)

		case <-ctx.Done():
			return nil
		}
	}
}

func (s *simpleAssetSimulator) publishMessage(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	timeoutCtx, cancel2 := s.obs.Span(timeoutCtx, "simpleAssetSimulator.publishMessage")
	defer cancel()
	defer cancel2()

	// Generate measurement
	measurement := s.simulatorCfg.GenerateRandomMeasurement()

	// Publish message
	err := s.publisher.Publish(timeoutCtx, measurement, s.GetId())
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to publish measurement")
	}
}

func (s *simpleAssetSimulator) Stop() error {
	close(s.stopChan)
	return nil
}

func (s *simpleAssetSimulator) GetId() string {
	return s.simulatorCfg.AssetId
}
