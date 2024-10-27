package simulator

import (
	"time"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/pkg/errors"
)

var ErrMinPowerGreaterThanMaxPower = errors.New("minPower is greater than maxPower")

type AssetType string

const (
	AssetTypeBattery = AssetType("battery")
	AssetTypeSolar   = AssetType("solar")
	AssetTypeWind    = AssetType("wind")
)

type Configuration struct {
	Id                  string        `json:"id"`
	AssetId             string        `json:"assetId"`
	Version             string        `json:"version"`
	Type                AssetType     `json:"type"`
	MeasurementInterval time.Duration `json:"measurementInterval"`
	MaxPower            float64       `json:"maxPower"`
	MinPower            float64       `json:"minPower"`
	MaxPowerStep        float64       `json:"maxPowerStep"`
}

func (c *Configuration) Validate() error {
	// Check if minPower is less than maxPower
	if c.MinPower > c.MaxPower {
		return ErrMinPowerGreaterThanMaxPower
	}

	// Validate the type
	switch c.Type {
	case AssetTypeBattery, AssetTypeSolar, AssetTypeWind:
		// Allowed
	default:
		return errors.Errorf("invalid asset type: %s", c.Type)
	}

	return nil
}

func (c *Configuration) GenerateRandomMeasurement() measurements.Measurement {
	maxPower := measurements.Power{
		Value: c.MaxPower,
		Unit:  measurements.UnitWatt,
	}
	minPower := measurements.Power{
		Value: c.MaxPower,
		Unit:  measurements.UnitWatt,
	}

	return measurements.NewRandomMeasurement(maxPower, minPower, c.MaxPowerStep)
}
