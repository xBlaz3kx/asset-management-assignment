package measurements

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

type Repository interface {
	AddMeasurement(ctx context.Context, assetId string, measurement Measurement) error
	GetLatestAssetMeasurement(ctx context.Context, assetID string) (*Measurement, error)
	GetAssetMeasurements(ctx context.Context, assetID string, timeRange TimeRange) ([]Measurement, error)
	GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params AssetMeasurementAveragedParams) ([]Measurement, error)
}

type TimeRange struct {
	From *time.Time `form:"from" binding:"required"`
	To   *time.Time `form:"to" binding:"required"`
}

var ErrInvalidTimeRange = errors.New("invalid time range")

func (t TimeRange) Validate() error {
	if t.From != nil {
		// From cannot be in the future
		if t.From.After(time.Now()) {
			return ErrInvalidTimeRange
		}

		// If To is provided, From cannot be after To
		if t.To != nil && t.From.After(*t.To) {
			return ErrInvalidTimeRange
		}
	}

	if t.To != nil {
		// To cannot be in the future
		if t.To.After(time.Now()) {
			return ErrInvalidTimeRange
		}
	}

	return nil
}

type AssetMeasurementAveragedParams struct {
	TimeRange
	GroupBy string `form:"groupBy" binding:"required"`
	Sort    string `form:"sort" binding:"required"`
}
