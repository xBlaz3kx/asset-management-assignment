package assets

import (
	"context"
	"testing"

	"asset-measurements-assignment/internal/domain"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

var (
	solarPanel = Asset{
		ID:      uuid.New().String(),
		Name:    "Solar panel",
		Type:    domain.AssetTypeSolar,
		Enabled: true,
	}

	windTurbine = Asset{
		ID:      uuid.New().String(),
		Name:    "Wind turbine",
		Type:    domain.AssetTypeWind,
		Enabled: true,
	}

	battery = Asset{
		ID:      uuid.New().String(),
		Name:    "Battery",
		Type:    domain.AssetTypeBattery,
		Enabled: true,
	}
)

type assetServiceTestSuite struct {
	suite.Suite
	service        Service
	repositoryMock *MockRepository
}

func (s *assetServiceTestSuite) SetupSuite() {
	s.repositoryMock = NewMockRepository(s.T())
	s.service = NewService(observability.NewNoopObservability(), s.repositoryMock)
}

func (s *assetServiceTestSuite) SetupTest() {
	s.repositoryMock = NewMockRepository(s.T())
}

func (s *assetServiceTestSuite) TestCreateAsset() {
	s.T().Skip("Not fully working")
	tests := []struct {
		name  string
		asset Asset
		err   bool
	}{
		{
			name:  "Solar panel",
			asset: solarPanel,
		},
		{
			name:  "Wind",
			asset: windTurbine,
		},
		{
			name:  "Battery",
			asset: battery,
		},
		{
			name: "Validation failed",
			asset: Asset{
				Name:    "Battery",
				Type:    "unknown",
				Enabled: true,
			},
			err: true,
		},
		{
			name:  "Database error",
			asset: solarPanel,
			err:   true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Solar panel", "Wind", "Battery":
				s.repositoryMock.EXPECT().CreateAsset(mock.Anything, tt.asset).Return(&tt.asset, nil).Once()
			case "Validation failed":
			case "Database error":
				s.repositoryMock.EXPECT().CreateAsset(mock.Anything, tt.asset).Return(nil, errors.New("some error")).Once()
			}

			asset, err := s.service.CreateAsset(context.Background(), tt.asset)
			if tt.err {
				s.Error(err)
				s.Nil(asset)
			} else {
				s.NoError(err)
				s.NotNil(asset)
				s.EqualValues(tt.asset, asset)
			}
		})
	}
}

func (s *assetServiceTestSuite) TestUpdateAsset() {
	s.T().Skip("Not fully working")
	tests := []struct {
		name  string
		asset Asset
		err   bool
	}{
		{
			name:  "Solar panel",
			asset: solarPanel,
		},
		{
			name:  "Wind",
			asset: windTurbine,
		},
		{
			name:  "Battery",
			asset: battery,
		},
		{
			name: "Validation failed",
			asset: Asset{
				Name:    "Battery",
				Type:    "unknown",
				Enabled: true,
			},
			err: true,
		},
		{
			name:  "Database error",
			asset: solarPanel,
			err:   true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Solar panel", "Wind", "Battery":
				s.repositoryMock.EXPECT().UpdateAsset(mock.Anything, tt.asset.ID, tt.asset).Return(&tt.asset, nil).Once()
			case "Validation failed":
			case "Database error":
				s.repositoryMock.EXPECT().UpdateAsset(mock.Anything, tt.asset.ID, tt.asset).Return(nil, errors.New("some error")).Once()
			}

			asset, err := s.service.UpdateAsset(context.Background(), tt.asset.ID, tt.asset)
			if tt.err {
				s.Error(err)
				s.Nil(asset)
			} else {
				s.NoError(err)
				s.NotNil(asset)
				s.EqualValues(tt.asset, asset)
			}
		})
	}
}

func (s *assetServiceTestSuite) TestDeleteAsset() {
	s.T().Skip("Not implemented")
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestGetAsset() {
	s.T().Skip("Not implemented")
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestGetAssets() {
	s.T().Skip("Not implemented")
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func TestAssetService(t *testing.T) {
	suite.Run(t, new(assetServiceTestSuite))
}
