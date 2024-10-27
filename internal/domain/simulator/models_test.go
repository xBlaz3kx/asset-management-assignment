package simulator

import (
	"testing"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/stretchr/testify/suite"
)

type simulatorTestSuite struct {
	suite.Suite
}

func (s *simulatorTestSuite) SetupTest() {
}

func (s *simulatorTestSuite) TearDownSuite() {
}

func (s *simulatorTestSuite) TestConfiguration_GenerateRandomMeasurement() {
	tests := []struct {
		name     string
		cfg      Configuration
		expected measurements.Measurement
	}{
		{
			name: "Random measurement",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actual := tt.cfg.GenerateRandomMeasurement()
			s.Assert().Equal(tt.expected, actual)
		})
	}
}

func (s *simulatorTestSuite) TestConfiguration_Validate() {
	tests := []struct {
		name     string
		cfg      Configuration
		expected error
	}{
		{
			name: "Valid configuration",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actual := tt.cfg.Validate()
			s.Assert().Equal(tt.expected, actual)
		})
	}
}

func TestSimulator(t *testing.T) {
	suite.Run(t, new(simulatorTestSuite))
}
