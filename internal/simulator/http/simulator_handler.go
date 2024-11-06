package http

import (
	"net/http"

	"asset-measurements-assignment/internal/domain/simulator"
	"github.com/gin-gonic/gin"
)

type SimulatorConfigHandler struct {
	service simulator.ConfigService
}

func NewSimulatorConfigHandler(service simulator.ConfigService) *SimulatorConfigHandler {
	return &SimulatorConfigHandler{service: service}
}

func (d *SimulatorConfigHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/assets/:assetId/config", d.CreateAssetConfig)
	router.GET("/assets/:assetId/config", d.GetCurrentAssetConfig)
	router.DELETE("/assets/:assetId/config/:configId", d.DeleteConfiguration)
}

// swagger:route GET /assets/{assetId}/config simulator getAssetConfig
// Get asset configuration by asset id
// ---
//
//	responses:
//	  204: Configuration
//	  400: errorResponse
//	  404: errorResponse
//	  500: errorResponse
//	  default: errorResponse
func (d *SimulatorConfigHandler) GetCurrentAssetConfig(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	getAsset, err := d.service.GetAssetConfiguration(reqCtx, assetId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, getAsset)
}

// swagger:route POST /assets/{assetId}/config simulator createAssetConfig
// Create asset configuration
// ---
//
//	Parameters:
//	 + name: configurationRequest
//	   in: body
//	   required: true
//	   type: CreateConfiguration
//
//	 responses:
//	   204: Configuration
//	   400: errorResponse
//	   404: errorResponse
//	   500: errorResponse
func (d *SimulatorConfigHandler) CreateAssetConfig(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	var configuration CreateConfiguration
	if err := ctx.ShouldBindJSON(&configuration); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	cfg, err := d.service.CreateConfiguration(reqCtx, configuration.toDomainConfiguration(assetId))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, cfg)
}

// swagger:route DELETE /assets/{assetId}/config/{configId} simulator deleteConfiguration
// Delete asset configuration
// ---
//
//	responses:
//	  204: emptyResponse
//	  404: errorResponse
//	  500: errorResponse
func (d *SimulatorConfigHandler) DeleteConfiguration(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")
	configId := ctx.Param("configId")

	err := d.service.DeleteConfiguration(reqCtx, assetId, configId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
