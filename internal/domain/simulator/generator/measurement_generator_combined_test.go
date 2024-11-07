package generator

import (
	"testing"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/stretchr/testify/suite"
)

type combinedGeneratorTestSuite struct {
	suite.Suite
	generator *CombinedMeasurementGenerator
}

func (s *combinedGeneratorTestSuite) SetupSuite() {
	s.generator = NewCombined(batteryCfg)
}

func (s *combinedGeneratorTestSuite) TestGenerateMeasurement() {
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
			name: "Battery full",
			previousMeasurement: &measurements.Measurement{
				Power: measurements.Power{
					Value: 30,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 100,
			},
		},
		{
			name: "Battery empty",
			previousMeasurement: &measurements.Measurement{
				Power: measurements.Power{
					Value: 0,
					Unit:  measurements.UnitWatt,
				},
				StateOfEnergy: 0,
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
				case "Battery full":
					s.Require().NotNil(tt.previousMeasurement)
					// If generated measurement is positive, it should be zero
					if measurement.Power.Value > 0 {
						s.Zero(measurement.Power.Value)
					} else {
						// todo
						// s.InDelta(tt.previousMeasurement.Power.Value, measurement.Power.Value, s.generator.cfg.MaxPowerStep)
					}
				case "Battery empty":
					s.Require().NotNil(tt.previousMeasurement)
					if measurement.Power.Value < 0 {
						s.Zero(measurement.Power.Value)
					} else {
						s.InDelta(tt.previousMeasurement.Power.Value, measurement.Power.Value, s.generator.cfg.MaxPowerStep)
					}
				}
			}
		})
	}
}

func (s *combinedGeneratorTestSuite) TestGetEnergyType() {
	s.EqualValues(domain.EnergyTypeCombined, s.generator.GetEnergyType())
}

func TestCombinedGenerator(t *testing.T) {
	suite.Run(t, new(combinedGeneratorTestSuite))
}
