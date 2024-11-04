package simulator

import "context"

type Repository interface {
	GetAssetConfiguration(ctx context.Context, assetId string) (*Configuration, error)
	GetConfigurations(ctx context.Context) ([]Configuration, error)
	CreateConfiguration(ctx context.Context, configuration Configuration) error
	DeleteConfiguration(ctx context.Context, configurationId string) error
}
