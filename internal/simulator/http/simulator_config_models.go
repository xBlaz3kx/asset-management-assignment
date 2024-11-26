package http

import (
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/simulator"
)

// swagger:model
type Configuration struct {
	Id                  string        `json:"id"`
	AssetId             string        `json:"assetId"`
	Version             string        `json:"version"`
	Type                string        `json:"type"`
	MeasurementInterval time.Duration `json:"measurementInterval"`
	MaxPower            float64       `json:"maxPower"`
	MinPower            float64       `json:"minPower"`
	MaxPowerStep        float64       `json:"maxPowerStep"`
}

// swagger:model
type CreateConfiguration struct {
	Type                string        `json:"type" binding:"required,oneof=battery solar wind motor hydro_turbine heat_turbine"`
	MeasurementInterval time.Duration `json:"measurementInterval" binding:"required,gte=100ms"`
	MaxPower            float64       `json:"maxPower" binding:"required"`
	MinPower            float64       `json:"minPower" binding:"required"`
	MaxPowerStep        float64       `json:"maxPowerStep" binding:"required"`
}

func (c CreateConfiguration) toDomainConfiguration(assetId string) simulator.Configuration {
	cfg := simulator.Configuration{
		AssetId:             assetId,
		Type:                domain.AssetType(c.Type),
		MeasurementInterval: c.MeasurementInterval,
		MaxPower:            c.MaxPower,
		MinPower:            c.MinPower,
		MaxPowerStep:        c.MaxPowerStep,
	}
	return cfg
}

type ConfigurationQuery struct {
}
