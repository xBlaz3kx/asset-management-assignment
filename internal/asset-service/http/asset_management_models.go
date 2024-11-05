package http

import (
	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/assets"
)

type GetAssetQuery struct {
	// Filter by asset name
	Enabled *bool `form:"enabled"`

	// Filter by asset type
	Type *string `form:"type" binding:"omitempty,oneof=battery solar wind"`
}

func (q GetAssetQuery) ToAssetQuery() assets.AssetQuery {
	return assets.AssetQuery{
		Enabled: q.Enabled,
		Type:    q.Type,
	}
}

type Asset struct {
	Id          string `json:"id" `
	Name        string `json:"name" validate:"required,min=4,max=100"`
	Description string `json:"description" validate:"omitempty,min=4,max=100"`
	Type        string `json:"type" validate:"required,oneof=battery solar wind"`
	Enabled     bool   `json:"enabled"`
}

type CreateAssetRequest struct {
	Name        string `json:"name" validate:"required,min=4,max=100"`
	Description string `json:"description" validate:"omitempty,min=4,max=100"`
	Type        string `json:"type" validate:"required,oneof=battery solar wind"`
	Enabled     bool   `json:"enabled"`
}

func (r *CreateAssetRequest) toDomainAsset() assets.Asset {
	return assets.Asset{
		Name:        r.Name,
		Description: r.Description,
		Type:        domain.AssetType(r.Type),
		Enabled:     r.Enabled,
	}
}

type UpdateAssetRequest struct {
	Name        string `json:"name" validate:"required,min=4,max=100"`
	Description string `json:"description" validate:"omitempty,min=4,max=100"`
	Type        string `json:"type" validate:"required,oneof=battery solar wind"`
	Enabled     bool   `json:"enabled"`
}

func (r *UpdateAssetRequest) toAsset(assetId string) assets.Asset {
	return assets.Asset{
		ID:          assetId,
		Name:        r.Name,
		Description: r.Description,
		Type:        domain.AssetType(r.Type),
		Enabled:     r.Enabled,
	}
}
