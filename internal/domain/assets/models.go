package assets

type AssetType string

const (
	AssetTypeBattery = AssetType("battery")
	AssetTypeSolar   = AssetType("solar")
	AssetTypeWind    = AssetType("wind")
)

type Asset struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        AssetType `json:"type"`
	Enabled     bool      `json:"enabled"`
}

func (a *Asset) Validate() error {
	return nil
}
