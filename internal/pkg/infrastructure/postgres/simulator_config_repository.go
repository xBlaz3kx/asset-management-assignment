package postgres

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SimulatorConfiguration represents the entity of the configuration in the database.
type SimulatorConfiguration struct {
	gorm.Model

	Version string

	AssetId string

	Type string

	MeasurementInterval time.Duration

	MaxPower float64

	MinPower float64

	MaxPowerStep float64
}

type SimulatorConfigurationRepository struct {
	obs observability.Observability
	db  *gorm.DB
}

func NewSimulatorConfigurationRepository(obs observability.Observability, db *gorm.DB) *SimulatorConfigurationRepository {
	return &SimulatorConfigurationRepository{
		obs: obs,
		db:  db,
	}
}

func (s *SimulatorConfigurationRepository) GetAssetConfiguration(ctx context.Context, assetId string) (*simulator.Configuration, error) {
	ctx, cancel := s.obs.Span(ctx, "configuration.repository.GetAssetConfiguration", zap.String("assetId", assetId))
	defer cancel()

	var dbConfig SimulatorConfiguration
	result := s.db.WithContext(ctx).Where("asset_id = ?", assetId).First(&dbConfig)
	if result.Error != nil {
		return nil, result.Error
	}

	cfg := toConfiguration(dbConfig)
	return &cfg, nil
}

func (s *SimulatorConfigurationRepository) GetConfigurations(ctx context.Context) ([]simulator.Configuration, error) {
	ctx, cancel := s.obs.Span(ctx, "configuration.repository.GetConfigurations")
	defer cancel()

	var dbConfigs []SimulatorConfiguration
	result := s.db.WithContext(ctx).Find(&dbConfigs)
	if result.Error != nil {
		return nil, result.Error
	}

	var configs []simulator.Configuration
	for _, dbConfig := range dbConfigs {
		configs = append(configs, toConfiguration(dbConfig))
	}

	return configs, nil
}

func toDBConfiguration(config simulator.Configuration) SimulatorConfiguration {
	return SimulatorConfiguration{
		Version:             config.Version,
		AssetId:             config.AssetId,
		Type:                string(config.Type),
		MeasurementInterval: config.MeasurementInterval,
		MaxPower:            config.MaxPower,
		MinPower:            config.MinPower,
		MaxPowerStep:        config.MaxPowerStep,
	}
}

func toConfiguration(dbConfig SimulatorConfiguration) simulator.Configuration {
	return simulator.Configuration{
		Version:             dbConfig.Version,
		AssetId:             dbConfig.AssetId,
		Type:                simulator.AssetType(dbConfig.Type),
		MeasurementInterval: dbConfig.MeasurementInterval,
		MaxPower:            dbConfig.MaxPower,
		MinPower:            dbConfig.MinPower,
		MaxPowerStep:        dbConfig.MaxPowerStep,
	}
}
