package service

import (
	"context"

	"asset-measurements-assignment/internal/domain/assets"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type ConsumerService interface {
	AddMeasurement(ctx context.Context, assetId string, measurement measurements.Measurement) error
}

type consumerService struct {
	obs             observability.Observability
	repository      measurements.Repository
	assetRepository assets.Repository
}

func NewConsumerService(obs observability.Observability, repository measurements.Repository, assetRepository assets.Repository) ConsumerService {
	return &consumerService{
		obs:             obs,
		repository:      repository,
		assetRepository: assetRepository,
	}
}

// AddMeasurement adds a new measurement to the database if the asset is enabled.
func (c *consumerService) AddMeasurement(ctx context.Context, assetId string, measurement measurements.Measurement) error {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "consumer.service.AddMeasurement", zap.String("assetId", assetId))
	defer cancel()
	logger.Info("Checking if asset exists")

	// Check if asset exists
	asset, err := c.assetRepository.GetAsset(ctx, assetId)
	if err != nil {
		return err
	}

	// Only add measurement if asset is enabled
	if asset.Enabled {
		logger.Info("Adding measurement to asset")
		err = c.repository.AddMeasurement(ctx, assetId, measurement)
		if err != nil {
			return err
		}
	} else {
		logger.Info("Asset is disabled, skipping measurement")
	}

	return nil
}
