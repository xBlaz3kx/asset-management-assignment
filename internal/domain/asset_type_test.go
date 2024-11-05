package domain

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type assetTestSuite struct {
	suite.Suite
}

func (s *assetTestSuite) TestAsset_GetEnergyType() {
	tests := []struct {
		name string
		a    AssetType
		want EnergyType
	}{
		{
			name: "Battery",
			a:    AssetTypeBattery,
			want: EnergyTypeCombined,
		},
		{
			name: "Solar",
			a:    AssetTypeSolar,
			want: EnergyTypeProducer,
		},
		{
			name: "Wind",
			a:    AssetTypeWind,
			want: EnergyTypeProducer,
		},
		{
			name: "Motor",
			a:    AssetTypeMotor,
			want: EnergyTypeConsumer,
		},
		{
			name: "Heater",
			a:    AssetTypeHeater,
			want: EnergyTypeConsumer,
		},
		{
			name: "HeatTurbine",
			a:    AssetTypeHeatTurbine,
			want: EnergyTypeProducer,
		},
		{
			name: "HydroTurbine",
			a:    AssetTypeHydroTurbine,
			want: EnergyTypeProducer,
		},
		{
			name: "Unknown",
			a:    AssetType("unknown"),
			want: "",
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got := tt.a.GetEnergyType()
			s.Equal(tt.want, got)
		})
	}
}

func (s *assetTestSuite) TestIsValidAssetType() {
	tests := []struct {
		name    string
		a       AssetType
		isValid bool
	}{
		{
			name:    "Battery",
			a:       AssetTypeBattery,
			isValid: true,
		},
		{
			name:    "Solar",
			a:       AssetTypeSolar,
			isValid: true,
		},
		{
			name:    "Wind",
			a:       AssetTypeWind,
			isValid: true,
		},
		{
			name:    "Motor",
			a:       AssetTypeMotor,
			isValid: true,
		},
		{
			name:    "Heater",
			a:       AssetTypeHeater,
			isValid: true,
		},
		{
			name:    "HeatTurbine",
			a:       AssetTypeHeatTurbine,
			isValid: true,
		},
		{
			name:    "HydroTurbine",
			a:       AssetTypeHydroTurbine,
			isValid: true,
		},
		{
			name:    "Unknown",
			a:       AssetType("unknown"),
			isValid: false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got := IsValidAssetType(tt.a)
			s.Equal(tt.isValid, got)
		})
	}
}

func TestAsset(t *testing.T) {
	suite.Run(t, new(assetTestSuite))
}
