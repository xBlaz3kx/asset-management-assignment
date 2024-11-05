package postgres

import (
	"context"
	"time"

	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/assets"
	"github.com/google/uuid"
	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Asset represents an asset entity in the database.
type Asset struct {
	ID          string `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Name        string         `gorm:"unique"`
	Description string
	Type        string
	Enabled     bool
}

func (u *Asset) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	return
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

// CreateAsset creates an asset in the database.
func (a *AssetRepository) CreateAsset(ctx context.Context, asset assets.Asset) error {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.CreateAsset", zap.Any("asset", asset))
	defer cancel()

	// Convert domain asset to database asset
	dbAsset := toDBAsset(asset)

	// Create asset in the database
	result := a.db.WithContext(ctx).Create(&dbAsset)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// UpdateAsset updates an asset in the database.
func (a *AssetRepository) UpdateAsset(ctx context.Context, assetId string, asset assets.Asset) error {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.UpdateAsset", zap.String("assetId", assetId))
	defer cancel()

	// Convert domain asset to database asset
	dbAsset := toDBAsset(asset)
	dbAsset.ID = assetId

	// Update asset in the database
	result := a.db.WithContext(ctx).Updates(&dbAsset)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// DeleteAsset deletes an asset from the database.
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

// GetAsset retrieves an asset from the database.
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

// GetAssets retrieves assets based on filters from the database.
func (a *AssetRepository) GetAssets(ctx context.Context, query assets.AssetQuery) ([]assets.Asset, error) {
	ctx, cancel := a.obs.Span(ctx, "asset.repository.GetAssets", zap.Any("query", query))
	defer cancel()

	db := a.db.WithContext(ctx)
	if query.Enabled != nil {
		db = db.Where("enabled = ?", *query.Enabled)
	}

	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}

	// Get assets from the database
	var dbAssets []Asset
	result := db.Find(&dbAssets)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert database assets to domain assets
	assets := toDomainAssets(dbAssets)
	return assets, nil
}

func toDomainAsset(dbAsset Asset) assets.Asset {
	return assets.Asset{
		ID:          dbAsset.ID,
		Name:        dbAsset.Name,
		Description: dbAsset.Description,
		Type:        domain.AssetType(dbAsset.Type),
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
