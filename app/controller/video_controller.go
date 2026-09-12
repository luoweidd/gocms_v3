package controller

import (
	"fmt"
	"gocms_v3/app/response"
	"gocms_v3/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	VideoService *service.VideoService
}

func NewVideoController() *VideoController {
	return &VideoController{
		VideoService: service.NewVideoService(),
	}
}

// Create 创建视频
func (c *VideoController) Create(ctx *gin.Context) {
	var req service.CreateVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	video, err := c.VideoService.CreateVideo(req)
	if err != nil {
		response.Error(ctx, 500, "创建失败")
		return
	}

	response.Success(ctx, video)
}

// Update 更新视频
func (c *VideoController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	var req service.UpdateVideoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.VideoService.UpdateVideo(id, req); err != nil {
		response.Error(ctx, 500, "更新失败")
		return
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// Delete 删除视频
func (c *VideoController) Delete(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	if err := c.VideoService.DeleteVideo(id); err != nil {
		response.Error(ctx, 500, "删除失败")
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetList 获取视频列表
func (c *VideoController) GetList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	categoryIDStr := ctx.Query("category_id")
	categoryID := uint(0)
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	genre := ctx.Query("genre")
	region := ctx.Query("region")
	language := ctx.Query("language")
	yearStr := ctx.Query("year")
	year := 0
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	// 解析 status 参数（nil 表示不过滤）
	statusStr := ctx.Query("status")
	var status *int
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			status = &s
		}
	}

	// 解析 keyword 参数
	keyword := ctx.Query("keyword")

	// 解析 author_id 参数
	authorIDStr := ctx.Query("author_id")
	authorID := uint(0)
	if authorIDStr != "" {
		if id, err := strconv.ParseUint(authorIDStr, 10, 32); err == nil {
			authorID = uint(id)
		}
	}

	// 解析排序参数
	sortBy := ctx.DefaultQuery("sort_by", "created_at")
	sortOrder := ctx.DefaultQuery("sort_order", "desc")

	params := service.VideoQueryParams{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: categoryID,
		Status:     status,
		Keyword:    keyword,
		AuthorID:   authorID,
		Genre:      genre,
		Region:     region,
		Language:   language,
		Year:       year,
		Director:   ctx.Query("director"),
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	pageData, err := c.VideoService.GetVideoList(params)
	if err != nil {
		fmt.Printf("[VIDEO] GetVideoList error: %v\n", err)
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, pageData)
}

// GetDetail 获取视频详情
func (c *VideoController) GetDetail(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	video, err := c.VideoService.GetVideoByID(id)
	if err != nil {
		response.Error(ctx, 404, "视频不存在")
		return
	}

	response.Success(ctx, video)
}

// GetPublished 获取已发布视频列表
func (c *VideoController) GetPublished(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	categoryIDStr := ctx.Query("category_id")
	categoryID := uint(0)
	if categoryIDStr != "" {
		if id, err := strconv.ParseUint(categoryIDStr, 10, 32); err == nil {
			categoryID = uint(id)
		}
	}

	genre := ctx.Query("genre")
	region := ctx.Query("region")
	language := ctx.Query("language")
	yearStr := ctx.Query("year")
	year := 0
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	pageData, err := c.VideoService.GetPublishedVideos(page, pageSize, categoryID, genre, region, language, year)
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, pageData)
}

// Publish 发布视频
func (c *VideoController) Publish(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	isTop := false
	req := service.UpdateVideoRequest{
		Status: intPtr(1),
		IsTop:  &isTop,
	}
	if err := c.VideoService.UpdateVideo(id, req); err != nil {
		response.Error(ctx, 500, "发布失败")
		return
	}

	response.SuccessWithMsg(ctx, "发布成功", nil)
}

// Unpublish 下架视频
func (c *VideoController) Unpublish(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	isTop := false
	req := service.UpdateVideoRequest{
		Status: intPtr(2),
		IsTop:  &isTop,
	}
	if err := c.VideoService.UpdateVideo(id, req); err != nil {
		response.Error(ctx, 500, "下架失败")
		return
	}

	response.SuccessWithMsg(ctx, "下架成功", nil)
}

// Top 置顶视频
func (c *VideoController) Top(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	isTop := true
	req := service.UpdateVideoRequest{IsTop: &isTop}
	if err := c.VideoService.UpdateVideo(id, req); err != nil {
		response.Error(ctx, 500, "置顶失败")
		return
	}

	response.SuccessWithMsg(ctx, "置顶成功", nil)
}

// UnTop 取消置顶
func (c *VideoController) UnTop(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	isTop := false
	req := service.UpdateVideoRequest{IsTop: &isTop}
	if err := c.VideoService.UpdateVideo(id, req); err != nil {
		response.Error(ctx, 500, "取消置顶失败")
		return
	}

	response.SuccessWithMsg(ctx, "已取消置顶", nil)
}

// Rate 评分
func (c *VideoController) Rate(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	ratingStr := ctx.DefaultQuery("rating", "5")
	rating, err := strconv.ParseFloat(ratingStr, 64)
	if err != nil || rating < 0 || rating > 10 {
		response.Error(ctx, 400, "评分无效（0-10）")
		return
	}

	if err := c.VideoService.SubmitRating(id, rating); err != nil {
		response.Error(ctx, 500, "评分失败")
		return
	}

	response.SuccessWithMsg(ctx, "评分成功", nil)
}

// Upload 上传视频文件
func (c *VideoController) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, 400, "请选择视频文件")
		return
	}

	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetUint("tenant_id")

	url, err := service.UploadFile(file, userID, tenantID)
	if err != nil {
		response.Error(ctx, 500, "上传失败")
		return
	}

	response.Success(ctx, gin.H{
		"url":       url,
		"file_size": file.Size,
		"filename":  file.Filename,
	})
}

// parseIDParam 解析路径中的 ID 参数
func parseIDParam(ctx *gin.Context) (uint, bool) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return 0, false
	}
	return uint(id), true
}

// intPtr 返回 int 指针
func intPtr(i int) *int {
	return &i
}

// float64Ptr 返回 float64 指针
func float64Ptr(f float64) *float64 {
	return &f
}
