package http

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	measurements "asset-measurements-assignment/internal/domain/measurements/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type measurementsHandlerTestSuite struct {
	suite.Suite
	router                 *gin.Engine
	mockMeasurementService *measurements.MockService
}

func (s *measurementsHandlerTestSuite) SetupSuite() {
	s.router = gin.Default()
	handler := NewMeasurementsGinHandler(s.mockMeasurementService)
	handler.RegisterRoutes(s.router)
}

func (s *measurementsHandlerTestSuite) SetupTest() {
	s.mockMeasurementService = measurements.NewMockService(s.T())
}

func (s *measurementsHandlerTestSuite) TestGetLatest() {
	tests := []struct {
		name         string
		responseBody string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: "",
			assetId:      "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "No measurement found",
			responseBody: "",
			assetId:      "2",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Asset not found",
			responseBody: "",
			assetId:      "3",
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Database error",
			responseBody: "",
			assetId:      "4",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/measurements/latest", tt.assetId)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *measurementsHandlerTestSuite) TestGetAvgWithinTimeInterval() {
	tests := []struct {
		name          string
		requestParams string
		responseBody  string
		assetId       string
		expectedCode  int
	}{
		{
			name:          "Success",
			requestParams: "",
			responseBody:  "",
			assetId:       "1",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "No entries in the time range",
			requestParams: "",
			responseBody:  "",
			assetId:       "2",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "From time is greater than to time",
			requestParams: "",
			responseBody:  "",
			assetId:       "3",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "From time is not provided",
			requestParams: "",
			responseBody:  "",
			assetId:       "4",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "To time is not provided",
			requestParams: "",
			responseBody:  "",
			assetId:       "5",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "Unsupported sort order",
			requestParams: "",
			responseBody:  "",
			assetId:       "6",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "Database error",
			requestParams: "",
			responseBody:  "",
			assetId:       "7",
			expectedCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/measurements/avg?%s", tt.assetId, tt.requestParams)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *measurementsHandlerTestSuite) TestGetWithinTimeInterval() {
	tests := []struct {
		name          string
		requestParams string
		responseBody  string
		assetId       string
		expectedCode  int
	}{
		{
			name:          "Success",
			requestParams: "",
			responseBody:  "",
			assetId:       "1",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "No entries in the time range",
			requestParams: "",
			responseBody:  "",
			assetId:       "2",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "From time is greater than to time",
			requestParams: "",
			responseBody:  "",
			assetId:       "3",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "From time is not provided",
			requestParams: "",
			responseBody:  "",
			assetId:       "4",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "To time is not provided",
			requestParams: "",
			responseBody:  "",
			assetId:       "5",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "Database error",
			requestParams: "",
			responseBody:  "",
			assetId:       "6",
			expectedCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/measurements?%s", tt.assetId, tt.requestParams)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.NotEmpty(w.Body)
		})
	}
}

func TestMeasurementsHandler(t *testing.T) {
	suite.Run(t, new(measurementsHandlerTestSuite))
}
