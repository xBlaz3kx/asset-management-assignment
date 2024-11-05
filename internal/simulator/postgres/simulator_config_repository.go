package postgres

import (
	"context"
	"strconv"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/simulator"
	"asset-measurements-assignment/internal/pkg/errors"
	"github.com/google/uuid"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SimulatorConfiguration represents the Simulator Configuration entity.
type SimulatorConfiguration struct {
	gorm.Model
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Version of the configuration
	Version int

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
	// Check if there is already a configuration for the asset

	prevCfg := &SimulatorConfiguration{}
	latestConfig := tx.Unscoped().
		Order("created_at desc").
		First(prevCfg, "asset_id = ?", u.AssetId)

	if latestConfig.Error != nil {
		return latestConfig.Error
	}

	if latestConfig.RowsAffected == 0 {
		// If there is no configuration, set the version to 1
		u.Version = 1
		return
	}

	// If there is a configuration, set the version to the previous version + 1
	u.Version = prevCfg.Version + 1
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
	result := s.db.WithContext(ctx).Where("asset_id = ?", assetId).Order("version desc").Find(&dbConfig).Limit(1)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, errors.ErrConfigNotFound
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

func (s *SimulatorConfigurationRepository) CreateConfiguration(ctx context.Context, configuration simulator.Configuration) (*simulator.Configuration, error) {
	ctx, cancel := s.obs.Span(ctx, "configuration.repository.CreateConfiguration", zap.Any("configuration", configuration))
	defer cancel()

	dbConfig := toDBConfiguration(configuration)
	result := s.db.WithContext(ctx).Create(&dbConfig)
	if result.Error != nil {
		return nil, result.Error
	}

	retCfg := toConfiguration(dbConfig)
	return &retCfg, nil
}

func (s *SimulatorConfigurationRepository) DeleteConfiguration(ctx context.Context, configurationId string) error {
	ctx, cancel := s.obs.Span(ctx, "configuration.repository.DeleteConfiguration", zap.String("configurationId", configurationId))
	defer cancel()

	result := s.db.WithContext(ctx).Delete(&SimulatorConfiguration{}, "id = ?", configurationId)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.ErrConfigNotFound
	}

	return nil
}

func toDBConfiguration(config simulator.Configuration) SimulatorConfiguration {
	return SimulatorConfiguration{
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
		Id:                  dbConfig.ID,
		Version:             strconv.Itoa(dbConfig.Version),
		AssetId:             dbConfig.AssetId,
		Type:                domain.AssetType(dbConfig.Type),
		MeasurementInterval: dbConfig.MeasurementInterval,
		MaxPower:            dbConfig.MaxPower,
		MinPower:            dbConfig.MinPower,
		MaxPowerStep:        dbConfig.MaxPowerStep,
	}
}
