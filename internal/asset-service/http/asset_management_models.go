package http

import (
	"asset-measurements-assignment/internal/domain"
	"asset-measurements-assignment/internal/domain/assets"
)

// swagger:parameters getAssets
type GetAssetQuery struct {
	// Filter by asset enabled status
	// required: false
	Enabled *bool `form:"enabled"`

	// Filter by asset type
	// required: false
	Type *string `form:"type" binding:"omitempty,oneof=battery solar wind"`
}

func (q GetAssetQuery) ToAssetQuery() assets.AssetQuery {
	return assets.AssetQuery{
		Enabled: q.Enabled,
		Type:    q.Type,
	}
}

// swagger:model
type Asset struct {
	Id          string `json:"id" `
	Name        string `json:"name" validate:"required,min=4,max=100"`
	Description string `json:"description" validate:"omitempty,min=4,max=100"`
	Type        string `json:"type" validate:"required,oneof=battery solar wind"`
	Enabled     bool   `json:"enabled"`
}

// swagger:model
type CreateAssetRequest struct {
	// Name of the asset
	// required: true
	// min length: 4
	// max length: 100
	Name string `json:"name" validate:"required,min=4,max=100"`

	// Description of the asset
	// required: false
	// min length: 4
	// max length: 100
	Description string `json:"description" validate:"omitempty,min=4,max=100"`

	// Type of the asset
	// required: true
	Type string `json:"type" validate:"required,oneof=battery solar wind"`

	// Enabled status of the asset
	// required: false
	Enabled bool `json:"enabled"`
}

func (r *CreateAssetRequest) toDomainAsset() assets.Asset {
	return assets.Asset{
		Name:        r.Name,
		Description: r.Description,
		Type:        domain.AssetType(r.Type),
		Enabled:     r.Enabled,
	}
}

// swagger:model
type UpdateAssetRequest struct {
	// Name of the asset
	// required: true
	// min length: 4
	// max length: 100
	Name string `json:"name" validate:"required,min=4,max=100"`

	// Description of the asset
	// required: false
	// min length: 4
	// max length: 100
	Description string `json:"description" validate:"omitempty,min=4,max=100"`

	// Type of the asset
	// required: true
	Type string `json:"type" validate:"required,oneof=battery solar wind"`

	// Enabled status of the asset
	// required: false
	Enabled bool `json:"enabled"`
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
