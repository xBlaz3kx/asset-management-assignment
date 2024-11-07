package service

import (
	"context"

	"asset-measurements-assignment/internal/domain/assets"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type measurementsService struct {
	obs             observability.Observability
	repository      measurements.Repository
	assetRepository assets.Repository
}

func NewMeasurementsService(obs observability.Observability, assetRepository assets.Repository, repository measurements.Repository) measurements.Service {
	return &measurementsService{
		obs:             obs,
		repository:      repository,
		assetRepository: assetRepository,
	}
}

// GetLatestAssetMeasurement returns the latest measurement for the given asset.
func (m *measurementsService) GetLatestAssetMeasurement(ctx context.Context, assetID string) (*measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx, "measurements.service.GetLatestAssetMeasurement", zap.String("assetId", assetID))
	defer cancel()
	logger.Info("Getting latest asset measurement")

	// Verify if asset exists
	_, err := m.assetRepository.GetAsset(ctx, assetID)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get asset")
		return nil, err
	}

	measurement, err := m.repository.GetLatestAssetMeasurement(ctx, assetID)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get latest asset measurement")
		return nil, err
	}

	return measurement, nil
}

// GetAssetMeasurements returns the measurements in an interval for the given asset.
func (m *measurementsService) GetAssetMeasurements(ctx context.Context, assetID string, timeRange measurements.TimeRange) ([]measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx, "measurements.service.GetAssetMeasurements", zap.String("assetId", assetID))
	defer cancel()
	logger.Info("Getting latest asset measurement")

	err := timeRange.Validate()
	if err != nil {
		logger.With(zap.Error(err)).Error("Invalid time range")
		return nil, assets.ErrTimeRangeViolation
	}

	// Verify if asset exists
	_, err = m.assetRepository.GetAsset(ctx, assetID)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get asset")
		return nil, err
	}

	measurement, err := m.repository.GetAssetMeasurements(ctx, assetID, timeRange)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get asset measurements")
		return nil, err
	}

	return measurement, nil
}

// GetAssetMeasurementsAveraged returns the average power from measurements for the given asset.
func (m *measurementsService) GetAssetMeasurementsAveraged(ctx context.Context, assetID string, params measurements.AssetMeasurementAveragedParams) ([]measurements.Measurement, error) {
	ctx, cancel, logger := m.obs.LogSpan(ctx,
		"measurements.service.GetAssetMeasurementsAveraged",
		zap.String("assetId", assetID),
		zap.Any("query", params),
	)
	defer cancel()
	logger.Info("Getting asset measurement averages")

	err := params.Validate()
	if err != nil {
		return nil, assets.ErrTimeRangeViolation
	}

	// Verify if asset exists
	_, err = m.assetRepository.GetAsset(ctx, assetID)
	if err != nil {
		return nil, err
	}

	measurements, err := m.repository.GetAssetMeasurementsAveraged(ctx, assetID, params)
	if err != nil {
		logger.With(zap.Error(err)).Error("Failed to get asset measurements")
		return nil, err
	}

	return measurements, nil
}
