package assets

import "github.com/go-playground/validator/v10"

type AssetType string

const (
	AssetTypeBattery = AssetType("battery")
	AssetTypeSolar   = AssetType("solar")
	AssetTypeWind    = AssetType("wind")
)

type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" validate:"required,min=4,max=100"`
	Description string    `json:"description" validate:"omitempty,min=4,max=100"`
	Type        AssetType `json:"type" validate:"required,oneof=battery solar wind"`
	Enabled     bool      `json:"enabled"`
}

func (a *Asset) Validate() error {
	return validator.New().Struct(a)
}
