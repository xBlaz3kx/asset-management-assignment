package simulator

import (
	"context"
)

type ConfigService interface {
	StartWorkersFromDatabaseConfigurations(ctx context.Context) error
	GetConfigurations(ctx context.Context) ([]Configuration, error)
	GetAssetConfiguration(ctx context.Context, assetId string) (*Configuration, error)
	CreateConfiguration(ctx context.Context, configuration Configuration) error
	DeleteConfiguration(ctx context.Context, assetId string, configurationId string) error
}
