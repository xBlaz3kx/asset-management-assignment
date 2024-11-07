package generator

import (
	"testing"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/stretchr/testify/suite"
)

type consumerGeneratorTestSuite struct {
	suite.Suite
	generator *ConsumerMeasurementGenerator
}

func (s *consumerGeneratorTestSuite) SetupSuite() {
	s.generator = NewConsumer(motorCfg)
}

func (s *consumerGeneratorTestSuite) TestGenerateMeasurement() {
	tests := []struct {
		name                string
		previousMeasurement *measurements.Measurement
		err                 bool
	}{
		{
			name:                "First measurement",
			previousMeasurement: nil,
		},
		{
			name: "Generate measurement with positive MaxPowerStep",
			previousMeasurement: &measurements.Measurement{
				Power: measurements.Power{
					Value: 0,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 0,
			},
		},
		{
			name: "Generate measurement with negative MaxPowerStep",
			previousMeasurement: &measurements.Measurement{
				Power: measurements.Power{
					Value: 30,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 30,
			},
		},
		{
			name: "State of charge is 100",
			previousMeasurement: &measurements.Measurement{
				Power: measurements.Power{
					Value: 30,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 100,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				s.generator.previousMeasurement = nil
				s.generator.cfg.MaxPowerStep = batteryCfg.MaxPowerStep
			})

			if tt.name == "Generate measurement with negative MaxPowerStep" {
				s.generator.cfg.MaxPowerStep = 0
			}

			s.generator.previousMeasurement = tt.previousMeasurement

			measurement, err := s.generator.GenerateMeasurement()
			if tt.err {
				s.Error(err)
				s.Nil(measurement)
			} else {
				s.NoError(err)

				// Todo determine state of energy change?
				// Determine output based on the test case
				switch tt.name {
				case "First measurement":
					s.EqualValues(0, measurement.Power.Value)
				case "Generate measurement with positive MaxPowerStep":
					s.Require().NotNil(tt.previousMeasurement)
					s.InDelta(tt.previousMeasurement.Power.Value, measurement.Power.Value, s.generator.cfg.MaxPowerStep)
				case "Generate measurement with negative MaxPowerStep":
					s.Require().NotNil(tt.previousMeasurement)
					s.InDelta(tt.previousMeasurement.Power.Value, measurement.Power.Value, s.generator.cfg.MaxPowerStep)
				case "State of charge is 100":
					s.Require().NotNil(tt.previousMeasurement)
					s.InDelta(tt.previousMeasurement.Power.Value, measurement.Power.Value, s.generator.cfg.MaxPowerStep)
					s.Equal(tt.previousMeasurement.StateOfEnergy, measurement.StateOfEnergy)
				}
			}
		})
	}
}

func (s *consumerGeneratorTestSuite) TestGetEnergyType() {
	s.EqualValues(domain.EnergyTypeConsumer, s.generator.GetEnergyType())
}

func TestConsumerGenerator(t *testing.T) {
	suite.Run(t, new(consumerGeneratorTestSuite))
}
