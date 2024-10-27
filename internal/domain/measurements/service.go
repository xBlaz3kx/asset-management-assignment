package measurements

import (
	"context"
)

type Service interface {
	GetLatestAssetMeasurement(ctx context.Context, assetID string) (*Measurement, error)
	GetAssetMeasurements(ctx context.Context, assetID string) ([]Measurement, error)
	GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params AssetMeasurementAveragedParams) ([]Measurement, error)
}
