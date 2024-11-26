package assets

import "context"

type Repository interface {
	CreateAsset(ctx context.Context, asset Asset) (*Asset, error)
	UpdateAsset(ctx context.Context, assetId string, asset Asset) (*Asset, error)
	DeleteAsset(ctx context.Context, assetId string) error
	GetAsset(ctx context.Context, assetId string) (*Asset, error)
	GetAssets(ctx context.Context, query AssetQuery) ([]Asset, error)
}

type AssetQuery struct {
	// Filter by asset name
	Enabled *bool `form:"enabled"`

	// Filter by asset type
	Type *string `form:"type" binding:"omitempty,oneof=battery solar wind"`
}
