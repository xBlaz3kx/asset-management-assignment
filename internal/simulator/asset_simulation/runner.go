package asset_simulation

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/pkg/errors"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Publisher interface {
	Publish(ctx context.Context, measurement measurements.Measurement, assetId string) error
}

type MeasurementGenerator interface {
	GenerateMeasurement() (*measurements.Measurement, error)
	GetEnergyType() domain.EnergyType
}

type runner struct {
	obs       observability.Observability
	generator MeasurementGenerator
	stopChan  chan bool
	interval  time.Duration
	publisher Publisher
	id        string
	isRunning bool
}

// NewRunner Creates a new Runner instance.
func NewRunner(
	obs observability.Observability,
	id string,
	interval time.Duration,
	generator MeasurementGenerator,
	publisher Publisher,
) (Runner, error) {
	if generator == nil {
		return nil, errors.New("generator is required")
	}

	if publisher == nil {
		return nil, errors.New("publisher is required")
	}

	if interval <= time.Millisecond*100 {
		return nil, errors.New("interval must be greater than 100ms")
	}

	if id == "" {
		return nil, errors.New("id is required")
	}

	return &runner{
		obs:       obs,
		stopChan:  make(chan bool),
		isRunning: false,
		generator: generator,
		publisher: publisher,
		id:        id,
		interval:  interval,
	}, nil
}

// Start creates a ticker and generates a measurement at each tick.
// The measurement is then published via a Publisher.
func (s *runner) Start(ctx context.Context) error {
	s.obs.Log().Debug(
		"Starting simulator runner",
		zap.String("id", s.id),
		zap.String("generatorType", string(s.generator.GetEnergyType())),
	)

	if s.interval <= time.Millisecond*100 {
		return errors.New("interval must be greater than 100ms")
	}

	if s.id == "" {
		return errors.New("id is required")
	}

	ticker := time.NewTicker(s.interval)
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
func (s *runner) publishMessage(ctx context.Context) {
	spanCtx, cancel2 := s.obs.Span(ctx, "simulator.runner.publishMessage")
	defer cancel2()

	// Generate measurement
	measurement, err := s.generator.GenerateMeasurement()
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to generate random measurement")
		return
	}

	// Publish message
	err = s.publisher.Publish(spanCtx, *measurement, s.GetId())
	if err != nil {
		s.obs.Log().With(zap.Error(err)).Error("Failed to publish measurement")
	}
}

func (s *runner) IsRunning() bool {
	return s.isRunning
}

func (s *runner) Stop() error {
	if !s.isRunning {
		return errors.New("worker is already stopped")
	}

	// Stop the worker
	s.stopChan <- true
	close(s.stopChan)
	return nil
}

func (s *runner) GetId() string {
	return s.id
}
