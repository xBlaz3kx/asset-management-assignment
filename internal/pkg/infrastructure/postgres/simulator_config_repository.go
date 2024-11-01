package postgres

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/google/uuid"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SimulatorConfiguration represents the Simulator Configuration entity.
type SimulatorConfiguration struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Version of the configuration
	Version string

	// AssetId that measurements are generated for
	AssetId string

	// Type of the asset
	Type string

	// MeasurementInterval is the interval between measurements
	MeasurementInterval time.Duration

	// MaxPower is the maximum power value that can be generated
	MaxPower float64

	// MinPower is the minimum power value that can be generated
	MinPower float64

	// MaxPowerStep is the maximum step between power values
	MaxPowerStep float64
}

func (u *SimulatorConfiguration) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
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

// GetAssetConfiguration returns the configuration for the asset with the given ID.
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

// GetConfigurations returns configurations for all assets.
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
