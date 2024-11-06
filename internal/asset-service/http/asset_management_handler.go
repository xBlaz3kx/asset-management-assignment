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

// swagger:route GET /assets/{assetId} asset getAssetById
// Get asset by id
// ---
//
//	responses:
//	  200: Asset
//	  404: errorResponse
//	  500: errorResponse
func (d *AssetGinHandler) GetAssetById(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	getAsset, err := d.service.GetAsset(reqCtx, assetId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, d.toAsset(*getAsset))
}

// swagger:route GET /assets asset getAssets
// Get assets
// ---
//
//	responses:
//	  200: []Asset
//	  400: errorResponse
//	  404: errorResponse
//	  500: errorResponse
func (d *AssetGinHandler) GetAssets(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var query GetAssetQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	getAssets, err := d.service.GetAssets(reqCtx, query.ToAssetQuery())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	// Map to API model
	ctx.JSON(http.StatusOK, d.toApiModels(getAssets))
}

func (d *AssetGinHandler) toApiModels(getAssets []assets.Asset) []Asset {
	var response []Asset
	for _, asset := range getAssets {
		response = append(response, d.toAsset(asset))
	}

	return response
}

func (d *AssetGinHandler) toAsset(asset assets.Asset) Asset {
	return Asset{
		Id:          asset.ID,
		Name:        asset.Name,
		Description: asset.Description,
		Type:        string(asset.Type),
		Enabled:     asset.Enabled,
	}
}

// swagger:route POST /assets asset createAsset
// Create asset
// ---
//
//	Parameters:
//	 + name: createAsset
//	   in: body
//	   required: true
//	   type: CreateAssetRequest
//
//	responses:
//	201: Asset
//	400: errorResponse
//	500: errorResponse
func (d *AssetGinHandler) CreateAsset(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()

	var req CreateAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	asset, err := d.service.CreateAsset(reqCtx, req.toDomainAsset())
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, d.toAsset(*asset))
}

// swagger:route PUT /assets/{assetId} asset updateAsset
// Update asset
// ---
//
//	Parameters:
//	 + name: updateAsset
//	   in: body
//	   required: true
//	   type: UpdateAssetRequest
//	responses:
//	 200: Asset
//	 400: errorResponse
//	 404: errorResponse
//	 500: errorResponse
func (d *AssetGinHandler) UpdateAsset(ctx *gin.Context) {
	reqCtx := ctx.Request.Context()
	assetId := ctx.Param("assetId")

	var request UpdateAssetRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(badRequest(err))
		return
	}

	asset, err := d.service.UpdateAsset(reqCtx, assetId, request.toAsset(assetId))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, d.toAsset(*asset))
}

// swagger:route DELETE /assets/{assetId} asset deleteAsset
// Delete an asset
// ---
//
//	responses:
//	 204: Asset
//	 404: errorResponse
//	 500: errorResponse
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
