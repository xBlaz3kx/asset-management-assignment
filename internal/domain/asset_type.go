package domain

import "github.com/pkg/errors"

type EnergyType string

const (
	EnergyTypeProducer = "producer"
	EnergyTypeConsumer = "consumer"
	EnergyTypeCombined = "combined"
)

type AssetType string

const (
	// Combined
	AssetTypeBattery = AssetType("battery")
	// Consumers
	AssetTypeMotor  = AssetType("motor")
	AssetTypeHeater = AssetType("heater")
	// Producers
	AssetTypeSolar        = AssetType("solar")
	AssetTypeWind         = AssetType("wind")
	AssetTypeHeatTurbine  = AssetType("heat_turbine")
	AssetTypeHydroTurbine = AssetType("hydro_turbine")
)

func IsValidAssetType(t AssetType) bool {
	switch t {
	case AssetTypeBattery, AssetTypeSolar, AssetTypeWind,
		AssetTypeMotor, AssetTypeHeater, AssetTypeHeatTurbine,
		AssetTypeHydroTurbine:
		return true
	default:
		return false
	}
}

func (a AssetType) String() string {
	return string(a)
}

// GetEnergyType returns the energy type based on asset type
func (a AssetType) GetEnergyType() EnergyType {
	switch a {
	case AssetTypeBattery:
		return EnergyTypeCombined
	case AssetTypeSolar, AssetTypeWind,
		AssetTypeHydroTurbine, AssetTypeHeatTurbine:
		return EnergyTypeProducer
	case AssetTypeMotor, AssetTypeHeater:
		return EnergyTypeConsumer
	default:
		return ""
	}
}

func (e EnergyType) ValidateBounds(min, max float64) error {
	switch e {
	case EnergyTypeProducer:
		// Range should be negative; max must be less than min
		if min < max {
			return errors.New("min power should be greater than max power")
		}

	case EnergyTypeConsumer:
		// Range should be positive; min must be less than max
		if max < min {
			return errors.New("max power should be greater than min power")
		}
	case EnergyTypeCombined:
		// Range can be either positive or negative; min must be less than max
		if min > max {
			return errors.New("min power should be less than max power")
		}
	default:
		return errors.New("invalid energy type")
	}

	return nil
}
