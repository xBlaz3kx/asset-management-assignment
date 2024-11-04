package http

import (
	"net/http"

	"asset-measurements-assignment/internal/domain/measurements"
	"github.com/gin-gonic/gin"
)

type MeasurementsGinHandler struct {
	service measurements.Service
}

func NewMeasurementsGinHandler(service measurements.Service) *MeasurementsGinHandler {
	return &MeasurementsGinHandler{service: service}
}

func (d *MeasurementsGinHandler) RegisterRoutes(router *gin.Engine) {
	rg := router.Group("/assets/:assetId/measurements")
	rg.GET("/latest", d.GetLatest)
	rg.GET("/avg", d.GetAvgWithinTimeInterval)
	rg.GET("", d.GetWithinTimeInterval)
}

func (d *MeasurementsGinHandler) GetLatest(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	assetMeasurement, err := d.service.GetLatestAssetMeasurement(reqCtx, assetId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, assetMeasurement)
}

func (d *MeasurementsGinHandler) GetAvgWithinTimeInterval(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	var query AssetMeasurementAveragedParams
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	assetMeasurementsAveraged, err := d.service.GetAssetMeasurementsAveraged(reqCtx, assetId, query.toDomainModel())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, assetMeasurementsAveraged)
}

func (d *MeasurementsGinHandler) GetWithinTimeInterval(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	var query TimeRange
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	assetMeasurements, err := d.service.GetAssetMeasurements(reqCtx, assetId, query.toDomainModel())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, assetMeasurements)
}
