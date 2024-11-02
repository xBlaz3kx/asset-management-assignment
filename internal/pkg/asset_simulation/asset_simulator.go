package asset_simulation

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, measurement measurements.Measurement, assetId string) error
}

var ErrMinPowerGreaterThanMaxPower = errors.New("minPower is greater than maxPower")

type Configuration struct {
	AssetId             string        `json:"assetId" validate:"required,gte=1"`
	MeasurementInterval time.Duration `json:"measurementInterval" validate:"required"`
	MaxPower            float64       `json:"maxPower" validate:"required,gte=0"`
	MinPower            float64       `json:"minPower" validate:"required,gte=0"`
	MaxPowerStep        float64       `json:"maxPowerStep"`
}

func (c *Configuration) Validate() error {
	err := validator.New().Struct(c)
	if err != nil {
		return err
	}

	// Check if minPower is less than maxPower
	if c.MinPower > c.MaxPower {
		return ErrMinPowerGreaterThanMaxPower
	}

	// Sanity check the measurement interval
	if c.MeasurementInterval <= time.Millisecond*100 {
		return errors.New("interval must be greater than 0")
	}

	return nil
}

// GenerateRandomMeasurement generates a random measurement based on the configuration.
func (c *Configuration) GenerateRandomMeasurement() (*measurements.Measurement, error) {
	maxPower := measurements.Power{
		Value: c.MaxPower,
		Unit:  measurements.UnitWatt,
	}
	minPower := measurements.Power{
		Value: c.MinPower,
		Unit:  measurements.UnitWatt,
	}

	return measurements.NewRandomMeasurement(minPower, maxPower, c.MaxPowerStep)
}

type simpleAssetSimulator struct {
	obs          observability.Observability
	simulatorCfg Configuration
	stopChan     chan bool
	publisher    Publisher
}

func NewSimpleAssetSimulator(obs observability.Observability, configuration Configuration, publisher Publisher) (AssetSimulator, error) {
	err := configuration.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate configuration")
	}

	return &simpleAssetSimulator{
		obs:          obs,
		stopChan:     make(chan bool),
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

func (s *simpleAssetSimulator) publishMessage(ctx context.Context) {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	timeoutCtx, cancel2 := s.obs.Span(timeoutCtx, "asset-simulator.simple.publishMessage")
	defer cancel()
	defer cancel2()

	// Generate measurement
	measurement, err := s.simulatorCfg.GenerateRandomMeasurement()
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to generate random measurement")
		return
	}

	// Publish message
	err = s.publisher.Publish(timeoutCtx, *measurement, s.GetId())
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
