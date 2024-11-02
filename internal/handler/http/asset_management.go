package http

import (
	"net/http"

	"asset-measurements-assignment/internal/domain/assets"
	"github.com/gin-gonic/gin"
)

type AssetGinHandler struct {
	service assets.Service
}

func NewAssetGinHandler(service assets.Service) *AssetGinHandler {
	return &AssetGinHandler{service: service}
}

func (d *AssetGinHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/assets", d.CreateAsset)
	router.GET("/assets", d.GetAssets)
	router.GET("/assets/:assetId", d.GetAssetById)
	router.PUT("/assets/:assetId", d.UpdateAsset)
	router.DELETE("/assets/:assetId", d.DeleteAsset)
}

func (d *AssetGinHandler) GetAssetById(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	getAsset, err := d.service.GetAsset(reqCtx, assetId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, getAsset)
}

func (d *AssetGinHandler) GetAssets(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var query assets.AssetQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	getAssets, err := d.service.GetAssets(reqCtx, query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, getAssets)
}

func (d *AssetGinHandler) CreateAsset(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var asset assets.Asset
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	err := d.service.CreateAsset(reqCtx, asset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}

func (d *AssetGinHandler) UpdateAsset(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	var asset assets.Asset
	if err := ctx.ShouldBindJSON(&asset); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	err := d.service.UpdateAsset(reqCtx, assetId, asset)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (d *AssetGinHandler) DeleteAsset(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	err := d.service.DeleteAsset(reqCtx, assetId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
