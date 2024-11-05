package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"asset-measurements-assignment/internal/domain/assets"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type assetManagementHandlerTestSuite struct {
	suite.Suite
	router           *gin.Engine
	mockAssetService *assets.MockService
}

func (s *assetManagementHandlerTestSuite) SetupSuite() {
	s.router = gin.Default()
	handler := NewAssetGinHandler(s.mockAssetService)
	handler.RegisterRoutes(s.router)
}

func (s *assetManagementHandlerTestSuite) SetupTest() {
	s.mockAssetService = assets.NewMockService(s.T())
}

func (s *assetManagementHandlerTestSuite) TestGetAssetById() {
	tests := []struct {
		name         string
		responseBody string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Database error",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s/measurements/avg", tt.assetId)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *assetManagementHandlerTestSuite) TestGetAssets() {
	tests := []struct {
		name          string
		requestParams string
		responseBody  string
		expectedCode  int
	}{
		{
			name:          "Success",
			requestParams: "",
			responseBody:  "",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "With query params",
			requestParams: "",
			responseBody:  "",
			expectedCode:  http.StatusOK,
		},
		{
			name:          "Database error",
			requestParams: "",
			responseBody:  "",
			expectedCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets?%s", tt.requestParams)
			req, _ := http.NewRequest("GET", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *assetManagementHandlerTestSuite) TestCreateAsset() {
	tests := []struct {
		name         string
		responseBody string
		requestBody  string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: "",
			requestBody:  "",
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Database error",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets")
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(tt.requestBody)))
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *assetManagementHandlerTestSuite) TestUpdateAsset() {
	tests := []struct {
		name         string
		responseBody string
		requestBody  string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			responseBody: "",
			requestBody:  "",
			assetId:      "",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Database error",
			responseBody: "",
			requestBody:  "",
			assetId:      "",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s", tt.assetId)
			req, _ := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(tt.requestBody)))
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func (s *assetManagementHandlerTestSuite) TestDeleteAsset() {
	tests := []struct {
		name         string
		assetId      string
		expectedCode int
	}{
		{
			name:         "Success",
			assetId:      "1",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Database error",
			assetId:      "2",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/assets/%s", tt.assetId)
			req, _ := http.NewRequest("DELETE", url, nil)
			s.router.ServeHTTP(w, req)

			s.Equal(tt.expectedCode, w.Code)
			s.Empty(w.Body)
		})
	}
}

func TestAssetManagementHandler(t *testing.T) {
	suite.Run(t, new(assetManagementHandlerTestSuite))
}
