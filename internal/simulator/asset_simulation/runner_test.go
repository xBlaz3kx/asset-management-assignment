package asset_simulation

import (
	"context"
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type runnerTestSuite struct {
	suite.Suite
	obs           observability.Observability
	publisherMock *MockPublisher
	mockGenerator *MockMeasurementGenerator
}

func (s *runnerTestSuite) SetupSuite() {
	s.obs = observability.NewNoopObservability()
}

func (s *runnerTestSuite) SetupTest() {
	s.publisherMock = NewMockPublisher(s.T())
	s.mockGenerator = NewMockMeasurementGenerator(s.T())
}

func (s *runnerTestSuite) TestNewRunner() {
	tests := []struct {
		name      string
		id        string
		interval  time.Duration
		generator MeasurementGenerator
		publisher Publisher
		err       bool
	}{
		{
			name:      "Valid runner",
			id:        "1",
			interval:  time.Second,
			generator: NewMockMeasurementGenerator(s.T()),
			publisher: NewMockPublisher(s.T()),
			err:       false,
		},
		{
			name:      "No asset id",
			id:        "",
			interval:  time.Second,
			generator: NewMockMeasurementGenerator(s.T()),
			publisher: NewMockPublisher(s.T()),
			err:       true,
		},
		{
			name:      "Invalid measurement interval",
			id:        "2",
			interval:  time.Millisecond,
			generator: NewMockMeasurementGenerator(s.T()),
			publisher: NewMockPublisher(s.T()),
			err:       true,
		},
		{
			name:      "Nil publisher",
			id:        "3",
			interval:  time.Millisecond,
			generator: NewMockMeasurementGenerator(s.T()),
			publisher: nil,
			err:       true,
		},
		{
			name:      "Nil generator",
			id:        "3",
			interval:  time.Millisecond,
			generator: nil,
			publisher: NewMockPublisher(s.T()),
			err:       true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewRunner(s.obs, tt.id, tt.interval, tt.generator, tt.publisher)
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
	currentTime := time.Now()

	tests := []struct {
		name        string
		id          string
		measurement measurements.Measurement
	}{
		{
			name: "Successfully published message",
			id:   "1",
			measurement: measurements.Measurement{
				Power: measurements.Power{
					Value: 110,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 100,
				Time:          currentTime,
			},
		},
		{
			name: "Generator unable to generate measurement",
			id:   "2",
			// Irrelevant measurement, as the generator will return an error
			measurement: measurements.Measurement{},
		},
		{
			name: "Publisher unable to publish message",
			id:   "3",
			measurement: measurements.Measurement{
				Power: measurements.Power{
					Value: 100,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 100,
				Time:          currentTime,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewRunner(s.obs, tt.id, time.Second, s.mockGenerator, s.publisherMock)
			s.Require().NoError(err)

			switch tt.name {
			case "Successfully published message":
				s.mockGenerator.EXPECT().GenerateMeasurement().Return(&tt.measurement, nil).Once()
				s.publisherMock.EXPECT().Publish(mock.Anything, tt.measurement, tt.id).Return(nil)
			case "Generator unable to generate measurement":
				s.mockGenerator.EXPECT().GenerateMeasurement().Return(nil, errors.New("generator unable to generate measurement")).Once()
			case "Publisher unable to publish message":
				s.mockGenerator.EXPECT().GenerateMeasurement().Return(&tt.measurement, nil).Once()
				s.publisherMock.EXPECT().Publish(mock.Anything, tt.measurement, tt.id).Return(errors.New("publisher unable to publish message"))
			}

			simulator.(*runner).publishMessage(context.Background())

			if tt.name == "Generator unable to generate measurement" {
				s.publisherMock.AssertNotCalled(t, "Publish")
			}
		})
	}
}

func (s *runnerTestSuite) TestStart() {
	tests := []struct {
		name     string
		id       string
		interval time.Duration
	}{
		{
			name:     "Successful start",
			id:       "1",
			interval: time.Second,
		},
		{
			name:     "Runner not started",
			id:       "2",
			interval: time.Second,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewRunner(s.obs, tt.id, tt.interval, s.mockGenerator, s.publisherMock)
			s.Require().NoError(err)

			switch tt.name {
			case "Successful start":
				measurement := measurements.Measurement{}
				s.mockGenerator.EXPECT().GetEnergyType().Return(domain.EnergyTypeConsumer)
				s.mockGenerator.EXPECT().GenerateMeasurement().Return(&measurement, nil)
				s.publisherMock.EXPECT().Publish(mock.Anything, measurement, tt.id).Return(nil)
			}

			if tt.name == "Successful start" {
				// Iterate over the interval 3 times
				// This should produce 3 calls to the generator and publisher
				ctx, cancel := context.WithTimeout(context.Background(), tt.interval*3)
				defer cancel()

				// Run the simulator in a goroutine
				go func() {
					err = simulator.Start(ctx)
					s.ErrorIs(err, context.DeadlineExceeded)
				}()

				go func() {
					// Check if it is running after the first interval
					time.Sleep(tt.interval)
					s.True(simulator.IsRunning())
				}()
			}

			time.Sleep(tt.interval * 4)
			s.False(simulator.IsRunning())

			s.mockGenerator.AssertNumberOfCalls(t, "GetEnergyType", 1)
			s.mockGenerator.AssertNumberOfCalls(t, "GenerateMeasurement", 3)
			s.publisherMock.AssertNumberOfCalls(t, "Publish", 3)
		})
	}
}

func (s *runnerTestSuite) TestStop() {
	s.T().Skip("Not fully working yet")
	tests := []struct {
		name     string
		id       string
		interval time.Duration
	}{
		{
			name:     "Started and stopped",
			id:       "1",
			interval: time.Second,
		},
		{
			name:     "Runner not started",
			id:       "2",
			interval: time.Second,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			simulator, err := NewRunner(s.obs, tt.id, tt.interval, s.mockGenerator, s.publisherMock)
			s.Require().NoError(err)

			switch tt.name {
			case "Started and stopped":
				measurement := measurements.Measurement{}
				s.mockGenerator.EXPECT().GetEnergyType().Return(domain.EnergyTypeConsumer)
				s.mockGenerator.EXPECT().GenerateMeasurement().Return(&measurement, nil)
				s.publisherMock.EXPECT().Publish(mock.Anything, measurement, tt.id).Return(nil)

				ctx, cancel := context.WithTimeout(context.Background(), tt.interval*3)
				defer cancel()

				// Run the simulator in a goroutine
				go func() {
					err = simulator.Start(ctx)
					s.NoError(err)
				}()

				go func() {
					// Check if it is running after the first interval
					time.Sleep(tt.interval)
					s.True(simulator.IsRunning())
					// Stop the runner using the .Stop method
					err = simulator.Stop()
					s.NoError(err)
				}()

				time.Sleep(tt.interval * 2)
				s.False(simulator.IsRunning())

				s.mockGenerator.AssertNumberOfCalls(t, "GetEnergyType", 1)
				s.mockGenerator.AssertNumberOfCalls(t, "GenerateMeasurement", 1)
				s.publisherMock.AssertNumberOfCalls(t, "Publish", 1)
			case "Runner not started":
				err = simulator.Stop()
				s.ErrorIs(err, errors.New("runner not started"))
			}
		})
	}
}

func TestRunner(t *testing.T) {
	suite.Run(t, new(runnerTestSuite))
}
