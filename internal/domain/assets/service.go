package assets

import (
	"context"

	"github.com/xBlaz3kx/DevX/observability"
	"go.uber.org/zap"
)

type Service interface {
	CreateAsset(ctx context.Context, asset Asset) error
	UpdateAsset(ctx context.Context, assetId string, asset Asset) error
	DeleteAsset(ctx context.Context, assetId string) error
	GetAsset(ctx context.Context, assetId string) (*Asset, error)
	GetAssets(ctx context.Context, query AssetQuery) ([]Asset, error)
}

type service struct {
	obs        observability.Observability
	repository Repository
}

func (s *service) CreateAsset(ctx context.Context, asset Asset) error {
	ctx, cancel, logger := s.obs.LogSpan(ctx, "assets.service.CreateAsset")
	defer cancel()
	logger.Info("Creating asset", zap.Any("asset", asset))

	// Validate the asset
	err := asset.Validate()
	if err != nil {
		return err
	}

	return s.repository.CreateAsset(ctx, asset)
}

func (s *service) UpdateAsset(ctx context.Context, assetId string, asset Asset) error {
	ctx, cancel, logger := s.obs.LogSpan(ctx, "assets.service.UpdateAsset")
	defer cancel()
	logger.Info("Updating an asset", zap.Any("assetId", assetId))

	// Validate the asset
	err := asset.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateAsset(ctx, assetId, asset)
}

func (s *service) DeleteAsset(ctx context.Context, assetId string) error {
	ctx, cancel, logger := s.obs.LogSpan(ctx, "assets.service.DeleteAsset")
	defer cancel()
	logger.Info("Deleting an asset", zap.String("assetId", assetId))

	return s.repository.DeleteAsset(ctx, assetId)
}

func (s *service) GetAsset(ctx context.Context, assetId string) (*Asset, error) {
	ctx, cancel, logger := s.obs.LogSpan(ctx, "assets.service.GetAsset")
	defer cancel()
	logger.Info("Getting an asset", zap.String("assetId", assetId))

	return s.repository.GetAsset(ctx, assetId)
}

func (s *service) GetAssets(ctx context.Context, query AssetQuery) ([]Asset, error) {
	ctx, cancel, logger := s.obs.LogSpan(ctx, "assets.service.GetAssets")
	defer cancel()
	logger.Info("Getting assets", zap.Any("query", query))

	return s.repository.GetAssets(ctx, query)
}

func NewService(obs observability.Observability, repository Repository) Service {
	return &service{
		repository: repository,
		obs:        obs,
	}
}
