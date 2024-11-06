package asset_simulation

import (
	"context"
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type runnerTestSuite struct {
	suite.Suite
	obs           observability.Observability
	publisherMock *MockPublisher
}

func (s *runnerTestSuite) SetupSuite() {
	s.obs = observability.NewNoopObservability()
	s.publisherMock = NewMockPublisher(s.T())
}

func (s *runnerTestSuite) SetupTest() {}

func (s *runnerTestSuite) TestNewSimpleAssetSimulator() {
	tests := []struct {
		name          string
		configuration simulator.Configuration
		err           bool
	}{
		{
			name: "Valid configuration",
			configuration: simulator.Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			err: false,
		},
		{
			name: "No asset id",
			configuration: simulator.Configuration{
				AssetId:             "",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			err: true,
		},
		{
			name: "Invalid measurement interval",
			configuration: simulator.Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Millisecond,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewRunner(s.obs, tt.configuration.AssetId, tt.configuration.MeasurementInterval, nil, s.publisherMock)
			if tt.err {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().NotNil(simulator)
				s.Assert().Implements((*Runner)(nil), simulator)
			}
		})
	}
}

func (s *runnerTestSuite) Test_publishMessage() {
	tests := []struct {
		name          string
		configuration simulator.Configuration
	}{
		{
			name: "Successfully published message",
			configuration: simulator.Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
		},
		{
			name: "MinPower greater than MaxPower",
			configuration: simulator.Configuration{
				AssetId:             "2",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1150,
				MaxPowerStep:        0,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator := &runner{
				obs:       s.obs,
				stopChan:  make(chan bool),
				publisher: s.publisherMock,
			}

			switch tt.name {
			case "Successfully published message":
				s.publisherMock.EXPECT().
					Publish(mock.Anything, mock.Anything, tt.configuration.AssetId).
					RunAndReturn(func(ctx context.Context, measurement measurements.Measurement, s string) error {
						// todo validate measurement
						return nil
					}).Return(nil)
			case "MinPower greater than MaxPower":
				defer s.publisherMock.AssertNotCalled(t, "Publish")
			}

			simulator.publishMessage(context.Background())
		})
	}
}

func (s *runnerTestSuite) TestStart() {
	tests := []struct {
		name          string
		configuration simulator.Configuration
	}{
		{
			name: "Successful start",
			configuration: simulator.Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
		},
		{
			name: "MinPower greater than MaxPower",
			configuration: simulator.Configuration{
				AssetId:             "2",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1150,
				MaxPowerStep:        0,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator := &runner{
				obs:       s.obs,
				stopChan:  make(chan bool),
				publisher: s.publisherMock,
			}

			switch tt.name {
			case "Successful start":
				s.publisherMock.EXPECT().
					Publish(mock.Anything, mock.Anything, tt.configuration.AssetId).
					RunAndReturn(func(ctx context.Context, measurement measurements.Measurement, assetId string) error {
						// todo validate interval
						s.Assert().Equal(tt.configuration.AssetId, assetId)
						return nil
					}).
					Return(nil)

			case "MinPower greater than MaxPower":
				defer s.publisherMock.AssertNotCalled(t, "Publish")
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.configuration.MeasurementInterval*3)
			_ = simulator.Start(ctx)
			cancel()
		})
	}
}

func TestSimpleAssetSimulator(t *testing.T) {
	suite.Run(t, new(runnerTestSuite))
}
