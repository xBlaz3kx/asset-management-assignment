package asset_simulation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type assetConfigurationTestSuite struct {
	suite.Suite
}

func (s *assetConfigurationTestSuite) TestConfiguration_Validate() {
	tests := []struct {
		name string
		cfg  Configuration
		err  bool
	}{
		{
			name: "Valid configuration with battery type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "battery",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            -300,
				MaxPowerStep:        1,
			},
			err: false,
		},
		{
			name: "Valid configuration with solar type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "solar",
				MeasurementInterval: time.Second,
				MaxPower:            -1000,
				MinPower:            -300,
				MaxPowerStep:        1,
			},
			err: false,
		},
		{
			name: "Valid configuration with wind type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "wind",
				MeasurementInterval: time.Second,
				MaxPower:            -1000,
				MinPower:            -300,
				MaxPowerStep:        1,
			},
			err: false,
		},
		{
			name: "Invalid configuration with battery type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "battery",
				MeasurementInterval: time.Second,
				MaxPower:            -1000,
				MinPower:            10000,
				MaxPowerStep:        1,
			},
			err: true,
		},
		{
			name: "Invalid configuration with solar type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "solar",
				MeasurementInterval: time.Second,
				MaxPower:            -0,
				MinPower:            -500,
				MaxPowerStep:        1,
			},
			err: true,
		},
		{
			name: "Invalid configuration with wind type",
			cfg: Configuration{
				AssetId:             "1234",
				MeasurementInterval: time.Second,
				Type:                "wind",
				MaxPower:            -500,
				MinPower:            -1000,
				MaxPowerStep:        1,
			},
			err: true,
		},
		{
			name: "Invalid asset type",
			cfg: Configuration{
				AssetId:             "1234",
				MeasurementInterval: time.Second,
				Type:                "abcdef",
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        1,
			},
			err: true,
		},
		{
			name: "Missing measurement interval",
			cfg: Configuration{
				AssetId:      "1234",
				Type:         "wind",
				MaxPower:     -1000,
				MinPower:     -500,
				MaxPowerStep: 1,
			},
			err: true,
		},
		{
			name: "Missing asset type",
			cfg: Configuration{
				AssetId:             "1234",
				Type:                "",
				MeasurementInterval: time.Second,
				MaxPower:            1000,
				MinPower:            500,
				MaxPowerStep:        1,
			},
			err: true,
		},
		{
			name: "Missing asset ID",
			cfg: Configuration{
				AssetId:             "",
				Type:                "wind",
				MeasurementInterval: time.Second,
				MaxPower:            -1000,
				MinPower:            -500,
				MaxPowerStep:        1,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.err {
				s.Errorf(err, "%s", tt.name)
			} else {
				s.NoErrorf(err, "%s", tt.name)
			}
		})
	}
}

func TestAssetConfiguration(t *testing.T) {
	suite.Run(t, new(assetConfigurationTestSuite))
}
