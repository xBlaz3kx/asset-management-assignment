package service

import (
	"context"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type measurementsService struct {
	obs        observability.Observability
	repository measurements.Repository
}

func NewService(obs observability.Observability, repository measurements.Repository) measurements.Service {
	return &measurementsService{
		obs:        obs,
		repository: repository,
	}
}

func (m *measurementsService) GetLatestAssetMeasurement(ctx context.Context, assetID string) (*measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx, "measurements.service.GetLatestAssetMeasurement", zap.String("assetId", assetID))
	defer cancel()
	logger.Info("Getting latest asset measurement")

	measurement, err := m.repository.GetLatestAssetMeasurement(ctx, assetID)
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

// todo time range
func (m *measurementsService) GetAssetMeasurements(ctx context.Context, assetID string) ([]measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx, "measurements.service.GetAssetMeasurements", zap.String("assetId", assetID))
	defer cancel()
	logger.Info("Getting latest asset measurement")

	measurement, err := m.repository.GetAssetMeasurements(ctx, assetID, measurements.TimeRange{})
	if err != nil {
		return nil, err
	}

	return measurement, nil
}

func (m *measurementsService) GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params measurements.AssetMeasurementAveragedParams) ([]measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx, "measurements.service.GetAssetMeasurementsAveraged", zap.String("assetId", assetID), zap.Any("query", params))
	defer cancel()
	logger.Info("Getting asset measurement averages")

	measurements, err := m.repository.GetAssetMeasurementsAveraged(ctx, assetID, params)
	if err != nil {
		return nil, err
	}

	return measurements, nil
}
