package service

import (
	"context"

	"asset-measurements-assignment/internal/domain/simulator"
	"asset-measurements-assignment/internal/simulator/asset_simulation"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type configService struct {
	obs        observability.Observability
	repository simulator.Repository
	manager    *asset_simulation.AssetSimulatorManager
	publisher  asset_simulation.Publisher
}

func (c *configService) StartWorkersFromDatabaseConfigurations(ctx context.Context) error {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "config.service.StartWorkersFromDatabaseConfigurations")
	defer cancel()
	logger.Info("Starting workers from database configurations")

	// Fetch assets from repository
	configs, err := c.GetConfigurations(ctx)
	if err != nil {
		return err
	}

	// Create workers from configurations
	for _, config := range configs {
		configuration := toAssetConfig(config)
		worker, err := asset_simulation.NewSimpleAssetSimulator(c.obs, configuration, c.publisher)
		if err != nil {
			c.obs.Log().Error("Failed to create worker", zap.Error(err))
			continue
		}

		c.manager.AddWorker(worker)
	}

	c.manager.StartWorkers(ctx)
	return nil
}

func (c *configService) GetConfigurations(ctx context.Context) ([]simulator.Configuration, error) {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "config.service.GetConfigurations")
	defer cancel()
	logger.Info("Getting configurations")

	return c.repository.GetConfigurations(ctx)
}

func (c *configService) GetAssetConfiguration(ctx context.Context, assetId string) (*simulator.Configuration, error) {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "config.service.GetAssetConfiguration")
	defer cancel()
	logger.Info("Getting asset configuration", zap.String("assetId", assetId))

	return c.repository.GetAssetConfiguration(ctx, assetId)
}

func (c *configService) CreateConfiguration(ctx context.Context, configuration simulator.Configuration) error {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "config.service.CreateConfiguration")
	defer cancel()
	logger.Info("Creating configuration", zap.Any("configuration", configuration))

	err := configuration.Validate()
	if err != nil {
		return err
	}

	err = c.repository.CreateConfiguration(ctx, configuration)
	if err != nil {
		return err
	}

	err2 := c.recreateWorker(configuration)
	if err2 != nil {
		logger.With(zap.Error(err2)).Error("Failed to recreate worker")
	}

	return nil
}

// recreateWorker removes the worker from the manager and creates a new worker with the new configuration
func (c *configService) recreateWorker(configuration simulator.Configuration) error {
	c.manager.RemoveWorker(configuration.AssetId)

	cfg := toAssetConfig(configuration)
	worker, err := asset_simulation.NewSimpleAssetSimulator(c.obs, cfg, c.publisher)
	if err != nil {
		return err
	}

	// Add and start worker from the new configuration
	c.manager.AddAndStartWorker(context.Background(), worker)
	return nil
}

func (c *configService) DeleteConfiguration(ctx context.Context, assetId string, configurationId string) error {
	ctx, cancel, logger := c.obs.LogSpan(ctx, "config.service.DeleteConfiguration")
	defer cancel()
	logger.Info("Deleting configuration", zap.String("configurationId", configurationId))

	err := c.repository.DeleteConfiguration(ctx, configurationId)
	if err != nil {
		return err
	}

	// Get previous configuration version for asset
	configuration, err := c.repository.GetAssetConfiguration(ctx, assetId)
	if err != nil {
		return err
	}

	// If there is configuration left for asset, only remove the worker
	if configuration == nil {
		c.manager.RemoveWorker(assetId)
		return nil
	}

	workerErr := c.recreateWorker(*configuration)
	if workerErr != nil {
		logger.With(zap.Error(workerErr)).Error("Failed to recreate worker")
	}

	return nil
}

func NewConfigService(obs observability.Observability, repository simulator.Repository, manager *asset_simulation.AssetSimulatorManager, publisher asset_simulation.Publisher) simulator.ConfigService {
	return &configService{
		repository: repository,
		obs:        obs,
		manager:    manager,
		publisher:  publisher,
	}
}

func toAssetConfig(configuration simulator.Configuration) asset_simulation.Configuration {
	cfg := asset_simulation.Configuration{
		Type:                configuration.Type,
		AssetId:             configuration.AssetId,
		MinPower:            configuration.MinPower,
		MaxPower:            configuration.MaxPower,
		MaxPowerStep:        configuration.MaxPowerStep,
		MeasurementInterval: configuration.MeasurementInterval,
	}
	return cfg
}
