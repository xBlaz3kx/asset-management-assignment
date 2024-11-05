package asset_simulation

import (
	"fmt"
	"math"
	"time"

	"asset-measurements-assignment/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Configuration struct {
	AssetId             string           `json:"assetId" validate:"required,gte=1"`
	Type                domain.AssetType `json:"type" validate:"required"`
	MeasurementInterval time.Duration    `json:"measurementInterval" validate:"required"`
	MaxPower            float64          `json:"maxPower" validate:"required"`
	MinPower            float64          `json:"minPower" validate:"required"`
	MaxPowerStep        float64          `json:"maxPowerStep"`
}

func (c *Configuration) Validate() error {
	err := validator.New().Struct(c)
	if err != nil {
		return err
	}

	// Validate min and max power based on asset type
	switch c.Type.GetEnergyType() {
	case domain.EnergyTypeCombined:
	case domain.EnergyTypeConsumer:
		// Must be positive
		if c.MinPower < 0.0 || c.MaxPower < 0.0 {
			return errors.New("minPower and maxPower must be positive for consumer assets")
		}
	case domain.EnergyTypeProducer:
		if c.MinPower > 0.0 || c.MaxPower > 0.0 {
			return errors.New("minPower and maxPower must be negative for producer assets")
		}
	default:
		return fmt.Errorf("unsupported asset type: %s", c.Type)
	}

	// Check if minPower is less than maxPower
	if math.Abs(c.MinPower) > math.Abs(c.MaxPower) {
		return ErrMinPowerGreaterThanMaxPower
	}

	// Sanity check the measurement interval
	if c.MeasurementInterval <= time.Millisecond*100 {
		return errors.New("interval must be greater than 0")
	}

	return nil
}
