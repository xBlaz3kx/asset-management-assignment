package assets

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAsset_Validate(t *testing.T) {
	tests := []struct {
		name  string
		asset Asset
		err   bool
	}{
		{
			name: "Valid asset",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        AssetTypeBattery,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Name is missing",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "",
				Description: "Some description",
				Type:        AssetTypeBattery,
				Enabled:     false,
			},
			err: true,
		},
		{
			name: "Name is too short",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "tes",
				Description: "Some description",
				Type:        AssetTypeBattery,
				Enabled:     false,
			},
			err: true,
		},
		{
			name: "Type is missing",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        "",
				Enabled:     false,
			},
			err: true,
		},
		{
			name: "Type invalid",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        "invalid",
				Enabled:     false,
			},
			err: true,
		},
		{
			name: "Description is too short",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "des",
				Type:        AssetTypeBattery,
				Enabled:     false,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()

			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
