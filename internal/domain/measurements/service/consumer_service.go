package service

import (
	"context"

	"asset-measurements-assignment/internal/domain/assets"
	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/xBlaz3kx/DevX/observability"
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

func (c *consumerService) AddMeasurement(ctx context.Context, assetId string, measurement measurements.Measurement) error {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "consumer.service.AddMeasurement")
	defer cancel()
	logger.Info("Adding measurement")

	// Check if asset exists & is enabled
	asset, err := c.assetRepository.GetAsset(ctx, assetId)
	if err != nil {
		return err
	}

	// Only add measurement if asset is enabled
	if asset.Enabled {
		err = c.repository.AddMeasurement(ctx, assetId, measurement)
		if err != nil {
			return err
		}
	}

	return nil
}
