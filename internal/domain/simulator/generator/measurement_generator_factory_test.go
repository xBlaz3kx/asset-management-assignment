package generator

import (
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	batteryCfg = simulator.Configuration{
		Id:                  uuid.New().String(),
		AssetId:             uuid.New().String(),
		Version:             "1.0.0",
		Type:                domain.AssetTypeBattery,
		MeasurementInterval: time.Second,
		MaxPower:            100,
		MinPower:            -50,
		MaxPowerStep:        10,
	}

	motorCfg = simulator.Configuration{
		Id:                  uuid.New().String(),
		AssetId:             uuid.New().String(),
		Version:             "1.0.0",
		Type:                domain.AssetTypeMotor,
		MeasurementInterval: time.Second,
		MaxPower:            100,
		MinPower:            -50,
		MaxPowerStep:        10,
	}

	heaterCfg = simulator.Configuration{
		Id:                  uuid.New().String(),
		AssetId:             uuid.New().String(),
		Version:             "1.0.0",
		Type:                domain.AssetTypeHeater,
		MeasurementInterval: time.Second,
		MaxPower:            100,
		MinPower:            -50,
		MaxPowerStep:        10,
	}

	solarCfg = simulator.Configuration{
		Id:                  uuid.New().String(),
		AssetId:             uuid.New().String(),
		Version:             "1.0.0",
		Type:                domain.AssetTypeSolar,
		MeasurementInterval: time.Second,
		MaxPower:            100,
		MinPower:            -50,
		MaxPowerStep:        10,
	}

	windCfg = simulator.Configuration{
		Id:                  uuid.New().String(),
		AssetId:             uuid.New().String(),
		Version:             "1.0.0",
		Type:                domain.AssetTypeWind,
		MeasurementInterval: time.Second,
		MaxPower:            100,
		MinPower:            -50,
		MaxPowerStep:        10,
	}
)

func TestGetGeneratorFromConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		cfg          simulator.Configuration
		expectedType domain.EnergyType
		err          bool
	}{
		{
			name:         "Combined",
			cfg:          batteryCfg,
			expectedType: domain.EnergyTypeCombined,
		},
		{
			name:         "Motor Consumer",
			cfg:          motorCfg,
			expectedType: domain.EnergyTypeConsumer,
		},
		{
			name:         "Heater Consumer",
			cfg:          heaterCfg,
			expectedType: domain.EnergyTypeConsumer,
		},
		{
			name:         "Solar Producer",
			cfg:          solarCfg,
			expectedType: domain.EnergyTypeProducer,
		},
		{
			name:         "Wind Producer",
			cfg:          windCfg,
			expectedType: domain.EnergyTypeProducer,
		},
		{
			name: "Unknown",
			cfg: simulator.Configuration{
				Type: "UnknownType",
			},
			expectedType: "",
			err:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generatorFromConfiguration, err := GetGeneratorFromConfiguration(tt.cfg)
			if tt.err {
				assert.Error(t, err)
				assert.Nil(t, generatorFromConfiguration)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, generatorFromConfiguration)
				assert.Equal(t, tt.expectedType, generatorFromConfiguration.GetEnergyType())
				assert.Implements(t, (*MeasurementGenerator)(nil), generatorFromConfiguration)
			}
		})
	}
}
