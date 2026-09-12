package controller

import (
	"strconv"
	"strings"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// MediaController 媒体控制器
type MediaController struct {
	MediaService *service.MediaService
}

// NewMediaController 创建媒体控制器
func NewMediaController() *MediaController {
	return &MediaController{
		MediaService: service.NewMediaService(),
	}
}

// ========== 媒体资源相关 ==========

// CreateMediaAsset 创建媒体资源
func (c *MediaController) CreateMediaAsset(ctx *gin.Context) {
	var req service.CreateMediaAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 从JWT中获取用户ID
	userID := ctx.GetUint("user_id")
	req.UserID = userID

	// 如果没有提供file_type，根据MIME类型自动判断
	if req.FileType == "" {
		req.FileType = c.MediaService.GetFileTypeFromMIME(req.MimeType)
	}

	asset, err := c.MediaService.CreateMediaAsset(req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, asset)
}

// UpdateMediaAsset 更新媒体资源
func (c *MediaController) UpdateMediaAsset(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req service.UpdateMediaAssetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	asset, err := c.MediaService.UpdateMediaAsset(uint(id), req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, asset)
}

// DeleteMediaAsset 删除媒体资源
func (c *MediaController) DeleteMediaAsset(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.MediaService.DeleteMediaAsset(uint(id)); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// BatchDeleteMediaAssets 批量删除媒体资源
func (c *MediaController) BatchDeleteMediaAssets(ctx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.Error(ctx, 400, "请选择要删除的资源")
		return
	}

	if err := c.MediaService.BatchDeleteMediaAssets(req.IDs); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量删除成功", nil)
}

// GetMediaAssetList 获取媒体资源列表
func (c *MediaController) GetMediaAssetList(ctx *gin.Context) {
	var req model.MediaAssetListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 200 {
		req.PageSize = 20
	}

	// 如果前端没有传 user_id，尝试从 JWT 中获取
	if req.UserID == 0 {
		userID := ctx.GetUint("user_id")
		if userID > 0 {
			req.UserID = uint64(userID)
		}
	}

	result, err := c.MediaService.GetMediaAssetList(req)
	if err != nil {
		response.Error(ctx, 500, "查询失败: "+err.Error())
		return
	}

	response.Success(ctx, result)
}

// GetMediaAssetByID 获取媒体资源详情
func (c *MediaController) GetMediaAssetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	asset, err := c.MediaService.GetMediaAssetByID(uint(id))
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, asset)
}

// ========== 相册相关 ==========

// CreateAlbum 创建相册
func (c *MediaController) CreateAlbum(ctx *gin.Context) {
	var req service.CreateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 从JWT中获取用户ID
	userID := ctx.GetUint("user_id")
	req.UserID = userID

	album, err := c.MediaService.CreateAlbum(req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, album)
}

// UpdateAlbum 更新相册
func (c *MediaController) UpdateAlbum(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req service.UpdateAlbumRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	album, err := c.MediaService.UpdateAlbum(uint(id), req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, album)
}

// DeleteAlbum 删除相册
func (c *MediaController) DeleteAlbum(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.MediaService.DeleteAlbum(uint(id)); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetAlbumList 获取相册列表
func (c *MediaController) GetAlbumList(ctx *gin.Context) {
	var req model.AlbumListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 200 {
		req.PageSize = 20
	}

	// 如果前端没有传 user_id，尝试从 JWT 中获取
	if req.UserID == 0 {
		userID := ctx.GetUint("user_id")
		if userID > 0 {
			req.UserID = uint64(userID)
		}
	}

	result, err := c.MediaService.GetAlbumList(req)
	if err != nil {
		response.Error(ctx, 500, "查询失败: "+err.Error())
		return
	}

	response.Success(ctx, result)
}

// GetAlbumByID 获取相册详情（含资源列表）
func (c *MediaController) GetAlbumByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	album, assets, err := c.MediaService.GetAlbumByID(uint(id))
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, gin.H{
		"album":  album,
		"assets": assets,
	})
}

// AddAssetToAlbum 添加资源到相册
func (c *MediaController) AddAssetToAlbum(ctx *gin.Context) {
	albumIDStr := ctx.Param("id")
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil || albumID == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req struct {
		AssetIDs []uint `json:"asset_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.MediaService.AddAssetToAlbum(uint(albumID), req.AssetIDs); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "添加成功", nil)
}

// RemoveAssetFromAlbum 从相册移除资源
func (c *MediaController) RemoveAssetFromAlbum(ctx *gin.Context) {
	albumIDStr := ctx.Param("id")
	assetIDStr := ctx.Param("asset_id")
	albumID, err1 := strconv.ParseUint(albumIDStr, 10, 32)
	assetID, err2 := strconv.ParseUint(assetIDStr, 10, 32)
	if err1 != nil || err2 != nil || albumID == 0 || assetID == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.MediaService.RemoveAssetFromAlbum(uint(albumID), uint(assetID)); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "移除成功", nil)
}

// GetMediaStats 获取媒体统计信息
func (c *MediaController) GetMediaStats(ctx *gin.Context) {
	// 先从查询参数中获取 user_id
	var reqUserID *uint64
	if userIDStr := ctx.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 64); err == nil && uid > 0 {
			reqUserID = &uid
		}
	}

	// 如果没有查询参数，从 JWT 中获取用户ID
	if reqUserID == nil {
		userIDUint := ctx.GetUint("user_id")
		if userIDUint > 0 {
			uid64 := uint64(userIDUint)
			reqUserID = &uid64
		}
	}

	stats, err := c.MediaService.GetMediaStats(reqUserID)
	if err != nil {
		response.Error(ctx, 500, "查询失败: "+err.Error())
		return
	}

	response.Success(ctx, stats)
}

// Helper functions

// parseIDFromString 从URL参数中解析ID
func parseIDFromString(idStr string) (uint, error) {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		return 0, strconv.ErrSyntax
	}
	return uint(id), nil
}

// extractFileExtension 提取文件扩展名
func extractFileExtension(filename string) string {
	ext := strings.Split(filename, ".")
	if len(ext) > 1 {
		return strings.ToLower(ext[len(ext)-1])
	}
	return ""
}
