package assets

import (
	"asset-measurements-assignment/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type Asset struct {
	ID          string           `json:"id"`
	Name        string           `json:"name" validate:"required,min=4,max=100"`
	Description string           `json:"description" validate:"omitempty,min=4,max=100"`
	Type        domain.AssetType `json:"type" validate:"required"`
	Enabled     bool             `json:"enabled"`
}

func (a *Asset) Validate() error {
	if !domain.IsValidAssetType(a.Type) {
		return errors.New("invalid asset type")
	}

	return validator.New().Struct(a)
}
