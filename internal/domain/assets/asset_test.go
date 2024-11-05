package assets

import (
	"testing"

	"asset-measurements-assignment/internal/domain"
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
			name: "Battery",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeBattery,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Heater",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeHeater,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Solar panel",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeSolar,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Motor",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeMotor,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Hydro turbine",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeHydroTurbine,
				Enabled:     false,
			},
			err: false,
		},
		{
			name: "Valid asset",
			asset: Asset{
				ID:          uuid.New().String(),
				Name:        "Test1234",
				Description: "Some description",
				Type:        domain.AssetTypeHeatTurbine,
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
				Type:        domain.AssetTypeBattery,
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
				Type:        domain.AssetTypeBattery,
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
				Type:        domain.AssetTypeBattery,
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
