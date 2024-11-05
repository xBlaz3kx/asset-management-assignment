package domain

type EnergyType string

const (
	EnergyTypeProducer = "producer"
	EnergyTypeConsumer = "consumer"
	EnergyTypeCombined = "combined"
)

type AssetType string

const (
	AssetTypeBattery      = AssetType("battery")
	AssetTypeSolar        = AssetType("solar")
	AssetTypeWind         = AssetType("wind")
	AssetTypeMotor        = AssetType("motor")
	AssetTypeHeater       = AssetType("heater")
	AssetTypeHeatTurbine  = AssetType("heatTurbine")
	AssetTypeHydroTurbine = AssetType("hydroTurbine")
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
