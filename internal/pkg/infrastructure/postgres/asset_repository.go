package postgres

import (
	"context"

	"asset-measurements-assignment/internal/domain/assets"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Asset represents an asset entity in the database.
type Asset struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Type        string
	Enabled     bool
}

type AssetRepository struct {
	obs observability.Observability
	db  *gorm.DB
}

func NewAssetRepository(obs observability.Observability, db *gorm.DB) *AssetRepository {
	return &AssetRepository{
		obs: obs,
		db:  db,
	}
}

func (a *AssetRepository) CreateAsset(ctx context.Context, asset assets.Asset) error {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.CreateAsset", zap.Any("asset", asset))
	defer cancel()

	// Convert domain asset to database asset
	dbAsset := Asset{
		Name:        asset.Name,
		Description: asset.Description,
		Type:        string(asset.Type),
		Enabled:     asset.Enabled,
	}

	// Create asset in the database
	result := a.db.WithContext(ctx).Create(&dbAsset)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a *AssetRepository) UpdateAsset(ctx context.Context, assetId string, asset assets.Asset) error {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.UpdateAsset", zap.String("assetId", assetId))
	defer cancel()

	// Convert domain asset to database asset
	dbAsset := toDBAsset(asset)
	//dbAsset.ID = assetId

	// Update asset in the database
	result := a.db.WithContext(ctx).Save(&dbAsset)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a *AssetRepository) DeleteAsset(ctx context.Context, assetId string) error {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.DeleteAsset", zap.String("assetId", assetId))
	defer cancel()

	// Soft delete asset from the database
	result := a.db.WithContext(ctx).Delete(&Asset{}, assetId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (a *AssetRepository) GetAsset(ctx context.Context, assetId string) (*assets.Asset, error) {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.GetAsset", zap.String("assetId", assetId))
	defer cancel()

	// Get asset from the database
	var dbAsset Asset
	result := a.db.WithContext(ctx).Where("id = ?", assetId).First(&dbAsset)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert database asset to domain asset
	asset := toDomainAsset(dbAsset)
	return &asset, nil
}

func (a *AssetRepository) GetAssets(ctx context.Context, query assets.AssetQuery) ([]assets.Asset, error) {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.GetAssets", zap.Any("query", query))
	defer cancel()

	// Get assets from the database
	var dbAssets []Asset
	result := a.db.WithContext(ctx).Find(&dbAssets)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert database assets to domain assets
	assets := toDomainAssets(dbAssets)
	return assets, nil
}

func toDomainAsset(dbAsset Asset) assets.Asset {
	return assets.Asset{
		Name:        dbAsset.Name,
		Description: dbAsset.Description,
		Type:        assets.AssetType(dbAsset.Type),
		Enabled:     dbAsset.Enabled,
	}
}

func toDomainAssets(dbAssets []Asset) []assets.Asset {
	var assets []assets.Asset
	for _, dbAsset := range dbAssets {
		assets = append(assets, toDomainAsset(dbAsset))
	}
	return assets
}

func toDBAsset(asset assets.Asset) Asset {
	return Asset{
		Name:        asset.Name,
		Description: asset.Description,
		Type:        string(asset.Type),
		Enabled:     asset.Enabled,
	}
}
