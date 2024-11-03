package asset_simulation

import (
	"context"
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type simpleAssetSimulatorTestSuite struct {
	suite.Suite
	obs           observability.Observability
	publisherMock *MockPublisher
}

func (s *simpleAssetSimulatorTestSuite) SetupSuite() {
	s.obs = observability.NewNoopObservability()
	s.publisherMock = NewMockPublisher(s.T())
}

func (s *simpleAssetSimulatorTestSuite) SetupTest() {}

func (s *simpleAssetSimulatorTestSuite) TestNewSimpleAssetSimulator() {
	tests := []struct {
		name          string
		configuration Configuration
		err           bool
	}{
		{
			name: "Valid configuration",
			configuration: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			err: false,
		},
		{
			name: "Invalid configuration",
			configuration: Configuration{
				AssetId:             "",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewSimpleAssetSimulator(s.obs, tt.configuration, s.publisherMock)
			if tt.err {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().NotNil(simulator)
				s.Assert().Implements((*AssetSimulator)(nil), simulator)
			}
		})
	}
}

func (s *simpleAssetSimulatorTestSuite) Test_publishMessage() {
	tests := []struct {
		name          string
		configuration Configuration
	}{
		{
			name: "Successfully published message",
			configuration: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
		},
		{
			name: "MinPower greater than MaxPower",
			configuration: Configuration{
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
			simulator := &simpleAssetSimulator{
				obs:          s.obs,
				simulatorCfg: tt.configuration,
				stopChan:     make(chan bool),
				publisher:    s.publisherMock,
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

func (s *simpleAssetSimulatorTestSuite) TestStart() {
	tests := []struct {
		name          string
		configuration Configuration
	}{
		{
			name: "Successful start",
			configuration: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
		},
		{
			name: "MinPower greater than MaxPower",
			configuration: Configuration{
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
			simulator := &simpleAssetSimulator{
				obs:          s.obs,
				simulatorCfg: tt.configuration,
				stopChan:     make(chan bool),
				publisher:    s.publisherMock,
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

func (s *simpleAssetSimulatorTestSuite) TestConfiguration_Validate() {
	tests := []struct {
		name    string
		cfg     Configuration
		isValid bool
	}{
		{
			name: "Valid configuration",
			cfg: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			isValid: true,
		},
		{
			name: "MinPower greater than MaxPower",
			cfg: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1100,
				MaxPowerStep:        0,
			},
			isValid: false,
		},
		{
			name: "Invalid measurement interval",
			cfg: Configuration{
				AssetId:             "1",
				MeasurementInterval: time.Millisecond,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			isValid: false,
		},
		{
			name: "Missing assetId",
			cfg: Configuration{
				AssetId:             "",
				MeasurementInterval: 150 * time.Millisecond,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			isValid: false,
		},
		{
			name: "Missing measurement interval",
			cfg: Configuration{
				AssetId:             "123",
				MeasurementInterval: 0,
				MaxPower:            1000,
				MinPower:            1,
				MaxPowerStep:        0,
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.isValid {
				s.Assert().NoError(err)
			} else {
				s.Assert().Error(err)
			}
		})
	}
}

func (s *simpleAssetSimulatorTestSuite) TestConfiguration_GenerateRandomMeasurement() {
	s.T().Skip("Not implemented")

	tests := []struct {
		name     string
		cfg      Configuration
		expected measurements.Measurement
		err      bool
	}{
		{
			name: "Random measurement",
			cfg: Configuration{
				MaxPower:     1000,
				MinPower:     1,
				MaxPowerStep: 0,
			},
			err: false,
		},
		{
			name: "Random measurement with step",
			cfg: Configuration{
				MaxPower:     1000,
				MinPower:     1,
				MaxPowerStep: 10,
			},
			err: false,
		},
		{
			name: "MinPower equals MaxPower",
			cfg: Configuration{
				MaxPower:     1000,
				MinPower:     1000,
				MaxPowerStep: 10,
			},
			err: true,
		},
		{
			name: "MinPower greater than MaxPower",
			cfg: Configuration{
				MaxPower:     1000,
				MinPower:     1100,
				MaxPowerStep: 10,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actual, err := tt.cfg.GenerateRandomMeasurement()
			if tt.err {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
				s.Assert().NotNil(actual)
				s.Assert().Equal(tt.expected.Power.Unit, actual.Power.Unit)
				s.Assert().True(tt.expected.Power.Value >= tt.cfg.MinPower && tt.expected.Power.Value <= tt.cfg.MaxPower)
			}
		})
	}
}

func TestSimpleAssetSimulator(t *testing.T) {
	suite.Run(t, new(simpleAssetSimulatorTestSuite))
}
