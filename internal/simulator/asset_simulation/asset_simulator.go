package asset_simulation

import (
	"context"
	"math/rand/v2"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/pkg/errors"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, measurement measurements.Measurement, assetId string) error
}

var ErrMinPowerGreaterThanMaxPower = errors.New("minPower is greater than maxPower")

type simpleAssetSimulator struct {
	obs          observability.Observability
	simulatorCfg Configuration
	stopChan     chan bool
	publisher    Publisher
	isRunning    bool
	// Last generated measurement
	previousMeasurement *measurements.Measurement
}

func NewSimpleAssetSimulator(obs observability.Observability, configuration Configuration, publisher Publisher) (AssetSimulator, error) {
	err := configuration.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate configuration")
	}

	return &simpleAssetSimulator{
		obs:          obs,
		stopChan:     make(chan bool),
		isRunning:    false,
		simulatorCfg: configuration,
		publisher:    publisher,
	}, nil
}

// Start creates a ticker and generates a measurement at each tick.
// The measurement is then published via a Publisher.
func (s *simpleAssetSimulator) Start(ctx context.Context) error {
	s.obs.Log().Debug("Starting simple asset simulator", zap.Any("config", s.simulatorCfg))
	ticker := time.NewTicker(s.simulatorCfg.MeasurementInterval)
	defer ticker.Stop()

	s.isRunning = true
	defer func() {
		s.isRunning = false
	}()

	for {
		select {
		case <-s.stopChan:
			return nil
		case <-ticker.C:
			s.publishMessage(ctx)
		case <-ctx.Done():
			if !errors.Is(ctx.Err(), context.Canceled) {
				s.obs.Log().With(zap.Error(ctx.Err())).Error("Context error")
				return ctx.Err()
			}

			return nil
		}
	}
}

// publishMessage generates a random measurement and publishes it via a Publisher.
func (s *simpleAssetSimulator) publishMessage(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	timeoutCtx, cancel2 := s.obs.Span(timeoutCtx, "asset-simulator.simple.publishMessage")
	defer cancel()
	defer cancel2()

	// Generate measurement
	measurement, err := s.generateRandomMeasurement()
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to generate random measurement")
		return
	}

	// Store the measurement as previous measurement
	s.previousMeasurement = measurement

	// Publish message
	err = s.publisher.Publish(timeoutCtx, *measurement, s.GetId())
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to publish measurement")
	}
}

// generateRandomMeasurement generates a random measurement for the asset based on the provided configuration
// and the previous measurement.
func (s *simpleAssetSimulator) generateRandomMeasurement() (*measurements.Measurement, error) {
	if s.previousMeasurement == nil {
		// This is the first measurement
		return &measurements.Measurement{
			Power: measurements.Power{
				Value: 0,
				Unit:  "W",
			},
			StateOfEnergy: 0,
			Time:          time.Now(),
		}, nil
	}

	// Generate a random step value
	if s.simulatorCfg.MaxPowerStep <= 0 {
		// Generate a random power step that wont exceed the MinPower or MaxPower

	}

	// Determine the sign of the power step (randomly)
	sign := rand.IntN(2)
	if sign == 0 {
		// Negative step
	} else {
		// Positive step
	}

	stepValue := rand.Float64() * s.simulatorCfg.MaxPowerStep

	if s.previousMeasurement.Power.Value+stepValue > s.simulatorCfg.MaxPower {

	}

	switch s.simulatorCfg.Type {
	case "battery":
	case "solar":
	case "wind":
	default:
	}

	// Depending on the type of asset, the power can increase or decrease
	return nil, nil
}

func (s *simpleAssetSimulator) IsRunning() bool {
	return s.isRunning
}

func (s *simpleAssetSimulator) Stop() error {
	close(s.stopChan)
	return nil
}

func (s *simpleAssetSimulator) GetId() string {
	return s.simulatorCfg.AssetId
}
