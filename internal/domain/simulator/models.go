package simulator

import (
	"time"

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
	if c.AssetId == "" {
		return errors.New("assetId is required")
	}

	// Check if minPower is less than maxPower
	if c.MinPower > c.MaxPower {
		return ErrMinPowerGreaterThanMaxPower
	}

	// Sanity check the measurement interval
	if c.MeasurementInterval <= time.Millisecond*100 {
		return errors.New("measurementInterval must be greater than 0")
	}

	// Validate the type
	switch c.Type {
	case AssetTypeBattery, AssetTypeSolar, AssetTypeWind:
		// Allowed
		return nil
	default:
		return errors.Errorf("invalid asset type: %s", c.Type)
	}
}
