package service

import (
	"context"
	"testing"
	"time"

	"asset-measurements-assignment/internal/domain/assets"
	assetsMocks "asset-measurements-assignment/internal/domain/assets"
	"asset-measurements-assignment/internal/domain/measurements"
	measurementMocks "asset-measurements-assignment/internal/domain/measurements/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/xBlaz3kx/DevX/observability"
)

type consumerServiceTestSuite struct {
	suite.Suite
	service         ConsumerService
	repository      *measurementMocks.MockRepository
	assetRepository *assetsMocks.MockRepository
}

func (s *consumerServiceTestSuite) SetupTest() {
	// Will regenerate the mocks for each subtest.
	s.repository = measurementMocks.NewMockRepository(s.T())
	s.assetRepository = assetsMocks.NewMockRepository(s.T())
	s.service = NewConsumerService(observability.NewNoopObservability(), s.repository, s.assetRepository)
}

func (s *consumerServiceTestSuite) TestConsumerService_AddMeasurement() {
	currentTime := time.Now()
	tests := []struct {
		name        string
		assetId     string
		measurement measurements.Measurement
		err         bool
	}{
		{
			name:    "Added measurement",
			assetId: "1",
			measurement: measurements.Measurement{
				Power:         measurements.Power{},
				StateOfEnergy: 0.0,
				Time:          currentTime,
			},
			err: false,
		},
		{
			name:    "Asset doesnt exist",
			assetId: "2",
			measurement: measurements.Measurement{
				Power:         measurements.Power{},
				StateOfEnergy: 0.0,
				Time:          currentTime,
			},
			err: true,
		},
		{
			name:    "Asset disabled",
			assetId: "3",
			measurement: measurements.Measurement{
				Power:         measurements.Power{},
				StateOfEnergy: 0.0,
				Time:          currentTime,
			},
			err: false,
		},
		{
			name:    "Repository error",
			assetId: "4",
			measurement: measurements.Measurement{
				Power:         measurements.Power{},
				StateOfEnergy: 0.0,
				Time:          currentTime,
			},
			err: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "Added measurement":
				s.assetRepository.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{Enabled: true}, nil)
				s.repository.EXPECT().AddMeasurement(mock.Anything, tt.assetId, tt.measurement).Return(nil)
			case "Asset doesnt exist":
				s.assetRepository.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(nil, errors.New("asset not found"))
			case "Asset disabled":
				s.assetRepository.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{Enabled: false}, nil)
			case "Repository error":
				s.assetRepository.EXPECT().GetAsset(mock.Anything, tt.assetId).Return(&assets.Asset{Enabled: true}, nil)
				s.repository.EXPECT().AddMeasurement(mock.Anything, tt.assetId, tt.measurement).Return(errors.New("repository error"))
			}

			err := s.service.AddMeasurement(context.Background(), tt.assetId, tt.measurement)
			if tt.err {
				s.Assert().Error(err)
			} else {
				s.Assert().NoError(err)
			}
		})
	}
}

func TestConsumerService(t *testing.T) {
	suite.Run(t, new(consumerServiceTestSuite))
}
