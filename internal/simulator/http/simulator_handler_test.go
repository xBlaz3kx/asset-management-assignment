package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"asset-measurements-assignment/internal/domain/simulator"
	simulatorMock "asset-measurements-assignment/internal/domain/simulator/mocks"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type simulatorHandlerTestSuite struct {
	suite.Suite
	router            *gin.Engine
	mockConfigService *simulatorMock.MockConfigService
}

func (s *simulatorHandlerTestSuite) SetupSuite() {
	s.router = gin.Default()
	handler := NewSimulatorConfigHandler(s.mockConfigService)
	handler.RegisterRoutes(s.router)
}

func (s *simulatorHandlerTestSuite) SetupTest() {
	s.mockConfigService = simulatorMock.NewMockConfigService(s.T())
}

func (s *simulatorHandlerTestSuite) TestCreateAssetConfig() {
	tests := []struct {
		name         string
		responseBody string
		requestBody  string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: `{}`,
			requestBody:  `{"type":"wind","measurementInterval":1000000000,"maxPower":-1000.0,"minPower":-1.0,"maxPowerStep":10}`,
			assetId:      "1",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Validation error",
			responseBody: ``,
			requestBody:  `{"type":"battery","measurementInterval":1000000000,"maxPower":-1000.0,"minPower":-1.0,"maxPowerStep":10}`,
			assetId:      "2",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Parse error",
			responseBody: `{"error":"An unknown error occurred","code":0,"description":"Please try again later or contact support"}`,
			requestBody:  `{"type":"battery","measurementInterval":1000000000,"maxPower":-1000.0,"minPower":-1.0,"maxPowerStep":10`,
			assetId:      "2",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Database error",
			responseBody: `{"error":"An unknown error occurred","code":0,"description":"Please try again later or contact support"}`,
			requestBody:  `{"type":"wind","measurementInterval":1000000000,"maxPower":-1000,"minPower":-1,"maxPowerStep":10}`,
			assetId:      "3",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "Success":
				s.mockConfigService.EXPECT().CreateConfiguration(mock.Anything, simulator.Configuration{
					Id:                  tt.assetId,
					Type:                "wind",
					MeasurementInterval: 1000000000,
					MaxPower:            -1000.0,
					MinPower:            -1.0,
					MaxPowerStep:        10,
				}).Return(&simulator.Configuration{}, nil)
			case "Database error":
				s.mockConfigService.EXPECT().
					CreateConfiguration(mock.Anything, simulator.Configuration{
						Id:                  tt.assetId,
						Type:                "wind",
						MeasurementInterval: 1000000000,
						MaxPower:            -1000.0,
						MinPower:            -1.0,
						MaxPowerStep:        10,
					}).
					Return(nil, errors.New("database error"))
			}

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/config", tt.assetId)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(tt.requestBody)))
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *simulatorHandlerTestSuite) TestGetCurrentAssetConfig() {
	tests := []struct {
		name         string
		responseBody string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: ``,
			assetId:      "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Not found error",
			responseBody: `{"error":"Config not found","code":2002,"description":""}`,
			assetId:      "2",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Database error",
			responseBody: `{"error":"An unknown error occurred","code":0,"description":"Please try again later or contact support"}`,
			assetId:      "3",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/config", tt.assetId)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *simulatorHandlerTestSuite) TestDeleteConfiguration() {
	tests := []struct {
		name         string
		responseBody string
		assetId      string
		configId     string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: "",
			assetId:      "1",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Not found error",
			responseBody: `{"error":"Config not found","code":2002,"description":""}`,
			assetId:      "2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Database error",
			responseBody: `{"error":"An unknown error occurred","code":0,"description":"Please try again later or contact support"}`,
			assetId:      "3",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/config/%s", tt.assetId, tt.configId)
			req, _ := http.NewRequest("DELETE", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func TestSimulatorHandler(t *testing.T) {
	t.Skip("Skip test")
	suite.Run(t, new(simulatorHandlerTestSuite))
}
