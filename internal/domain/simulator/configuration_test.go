package simulator

import (
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type simulatorTestSuite struct {
	suite.Suite
}

func (s *simulatorTestSuite) TestConfiguration_Validate() {
	assetId := uuid.New().String()
	tests := []struct {
		name string
		cfg  Configuration
		err  bool
	}{
		{
			name: "Valid configuration",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             assetId,
				Version:             "1.0.0",
				Type:                domain.AssetTypeBattery,
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
			err: false,
		},
		{
			name: "Invalid measurement interval",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             assetId,
				Version:             "1.0.0",
				Type:                domain.AssetTypeBattery,
				MeasurementInterval: 0,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
			err: true,
		},
		{
			name: "MinPower greater than MaxPower",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             assetId,
				Version:             "1.0.0",
				Type:                domain.AssetTypeBattery,
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            1500,
				MaxPowerStep:        0,
			},
			err: true,
		},
		{
			name: "Invalid asset type",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             assetId,
				Version:             "1.0.0",
				Type:                "car",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
			err: true,
		},
		{
			name: "Empty asset type",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             assetId,
				Version:             "1.0.0",
				Type:                "",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
			err: true,
		},
		{
			name: "Missing assetId",
			cfg: Configuration{
				Id:                  uuid.New().String(),
				AssetId:             "",
				Version:             "1.0.0",
				Type:                domain.AssetTypeBattery,
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        0,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			actual := tt.cfg.Validate()
			if tt.err {
				s.Error(actual)
			} else {
				s.NoError(actual)
			}
		})
	}
}

func TestSimulator(t *testing.T) {
	suite.Run(t, new(simulatorTestSuite))
}
