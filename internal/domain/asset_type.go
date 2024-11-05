package domain

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
