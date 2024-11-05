package simulator

import (
	"math"
	"time"

	"asset-measurements-assignment/internal/domain"
	"github.com/pkg/errors"
)

var ErrMinPowerGreaterThanMaxPower = errors.New("minPower is greater than maxPower")

type Configuration struct {
	Id                  string           `json:"id"`
	AssetId             string           `json:"assetId"`
	Version             string           `json:"version"`
	Type                domain.AssetType `json:"type"`
	MeasurementInterval time.Duration    `json:"measurementInterval"`
	MaxPower            float64          `json:"maxPower"`
	MinPower            float64          `json:"minPower"`
	MaxPowerStep        float64          `json:"maxPowerStep"`
}

func (c *Configuration) Validate() error {
	if c.AssetId == "" {
		return errors.New("assetId is required")
	}

	// Check if minPower is less than maxPower
	if math.Abs(c.MinPower) > math.Abs(c.MaxPower) {
		return ErrMinPowerGreaterThanMaxPower
	}

	// Sanity check the measurement interval
	if c.MeasurementInterval < time.Millisecond*100 {
		return errors.New("measurementInterval must be greater than 0")
	}

	// Validate the type
	if !domain.IsValidAssetType(c.Type) {
		return errors.New("invalid asset type")
	}

	return nil
}
