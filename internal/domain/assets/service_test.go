package assets

import (
	"testing"

	assetMocks "asset-measurements-assignment/internal/domain/assets/mocks"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type assetServiceTestSuite struct {
	suite.Suite
	service        Service
	repositoryMock *assetMocks.MockRepository
}

func (s *assetServiceTestSuite) SetupTest() {
	s.repositoryMock = assetMocks.NewMockRepository(s.T())
	s.service = NewService(observability.NewNoopObservability(), s.repositoryMock)
}

func (s *assetServiceTestSuite) TearDownSuite() {
}

func (s *assetServiceTestSuite) TestCreateAsset() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestUpdateAsset() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestDeleteAsset() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestGetAsset() {
	tests := []struct {
		name string
	}{}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

		})
	}
}

func (s *assetServiceTestSuite) TestGetAssets() {
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
