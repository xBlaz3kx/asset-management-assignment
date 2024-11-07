package service

import (
	"context"
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain/assets"
	"asset-measurements-assignment/internal/domain/measurements"
	measurementMocks "asset-measurements-assignment/internal/domain/measurements/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type measurementServiceTestSuite struct {
	suite.Suite
	service                    *measurementsService
	assetRepositoryMock        *assets.MockRepository
	measurementsRepositoryMock *measurementMocks.MockRepository
}

func (s *measurementServiceTestSuite) SetupSuite() {
	s.assetRepositoryMock = assets.NewMockRepository(s.T())
	s.measurementsRepositoryMock = measurementMocks.NewMockRepository(s.T())

	s.service = &measurementsService{
		obs:             observability.NewNoopObservability(),
		repository:      s.measurementsRepositoryMock,
		assetRepository: s.assetRepositoryMock,
	}
}
func (s *measurementServiceTestSuite) SetupTest() {
	s.assetRepositoryMock = assets.NewMockRepository(s.T())
	s.measurementsRepositoryMock = measurementMocks.NewMockRepository(s.T())
}

func (s *measurementServiceTestSuite) TestGetLatestAssetMeasurement() {
	tests := []struct {
		name    string
		assetId string
		err     bool
	}{
		{
			name:    "Get latest asset measurement",
			assetId: "1",
			err:     false,
		},
		{
			name:    "Asset not found",
			assetId: "2",
			err:     true,
		},
		{
			name:    "Failed to get latest asset measurement",
			assetId: "3",
			err:     true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			expectedMeasurement := &measurements.Measurement{
				Time: time.Now(),
			}
			switch tt.name {
			case "Get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil).Once()
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(expectedMeasurement, nil).Once()
			case "Asset not found":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(nil, assets.ErrAssetNotFound).Once()
			case "Failed to get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil).Once()
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(nil, errors.New("some error")).Once()
			}

			measurement, err := s.service.GetLatestAssetMeasurement(context.TODO(), tt.assetId)
			if tt.err {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(expectedMeasurement)
				s.Equal(expectedMeasurement, measurement)
			}
		})
	}
}

func (s *measurementServiceTestSuite) TestGetAssetMeasurements() {
	tests := []struct {
		name    string
		assetId string
		err     bool
	}{
		{
			name:    "Get latest asset measurement",
			assetId: "1",
			err:     false,
		},
		{
			name:    "Asset not found",
			assetId: "2",
			err:     true,
		},
		{
			name:    "Failed to get latest asset measurement",
			assetId: "3",
			err:     true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			expectedMeasurement := &measurements.Measurement{}
			switch tt.name {
			case "Get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil).Once()
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(expectedMeasurement, nil)
			case "Asset not found":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(nil, assets.ErrAssetNotFound)
			case "Failed to get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil)
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(nil, errors.New("some error"))
			}

			measurement, err := s.service.GetLatestAssetMeasurement(context.TODO(), tt.assetId)
			if tt.err {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(expectedMeasurement)
				s.Equal(expectedMeasurement, measurement)
			}
		})
	}
}

func (s *measurementServiceTestSuite) TestGetAssetMeasurementsAveraged() {
	tests := []struct {
		name    string
		assetId string
		err     bool
	}{
		{
			name:    "Get latest asset measurement",
			assetId: "1",
			err:     false,
		},
		{
			name:    "Asset not found",
			assetId: "2",
			err:     true,
		},
		{
			name:    "Failed to get latest asset measurement",
			assetId: "3",
			err:     true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			expectedMeasurement := &measurements.Measurement{}
			switch tt.name {
			case "Get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil).Once()
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(expectedMeasurement, nil).Once()
			case "Asset not found":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(nil, assets.ErrAssetNotFound).Once()
			case "Failed to get latest asset measurement":
				s.assetRepositoryMock.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{}, nil).Once()
				s.measurementsRepositoryMock.EXPECT().GetLatestAssetMeasurement(mock.Anything, tt.assetId).Return(nil, errors.New("some error")).Once()
			}

			measurement, err := s.service.GetLatestAssetMeasurement(context.TODO(), tt.assetId)
			if tt.err {
				s.Error(err)
			} else {
				s.NoError(err)
				s.NotNil(expectedMeasurement)
				s.Equal(expectedMeasurement, measurement)
			}
		})
	}
}

func TestMeasurementService(t *testing.T) {
	t.Skip()
	suite.Run(t, new(measurementServiceTestSuite))
}
