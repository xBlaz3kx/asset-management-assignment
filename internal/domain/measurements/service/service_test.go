package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type measurementServiceTestSuite struct {
	suite.Suite
}

func (s *measurementServiceTestSuite) SetupTest() {
}

func (s *measurementServiceTestSuite) TearDownSuite() {
}

func (s *measurementServiceTestSuite) TestGetLatestAssetMeasurement() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *measurementServiceTestSuite) TestGetAssetMeasurements() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *measurementServiceTestSuite) TestGetAssetMeasurementsAveraged() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestMeasurementService(t *testing.T) {
	suite.Run(t, new(measurementServiceTestSuite))
}
