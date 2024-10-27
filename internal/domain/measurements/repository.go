package measurements

import (
	"context"
	"time"
)

type Repository interface {
	AddMeasurement(ctx context.Context, assetId string, measurement Measurement) error
	GetLatestAssetMeasurement(ctx context.Context, assetID string) (*Measurement, error)
	GetAssetMeasurements(ctx context.Context, assetID string, timeRange TimeRange) ([]Measurement, error)
	GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params AssetMeasurementAveragedParams) ([]Measurement, error)
}

type TimeRange struct {
	From *time.Time
	To   *time.Time
}

type AssetMeasurementAveragedParams struct {
	TimeRange
	GroupBy string
	Sort    string
}
