package controller

import (
	"strconv"

	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

type VideoCategoryController struct {
	CategoryService *service.VideoCategoryService
}

func NewVideoCategoryController() *VideoCategoryController {
	return &VideoCategoryController{
		CategoryService: service.NewVideoCategoryService(),
	}
}

// Create 创建分类
func (c *VideoCategoryController) Create(ctx *gin.Context) {
	var req service.CreateVideoCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	category, err := c.CategoryService.Create(req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, category)
}

// Update 更新分类
func (c *VideoCategoryController) Update(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	var req service.UpdateVideoCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	category, err := c.CategoryService.Update(id, req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, category)
}

// Delete 删除分类
func (c *VideoCategoryController) Delete(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	if err := c.CategoryService.Delete(id); err != nil {
		response.Error(ctx, 400, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetTree 获取分类树
func (c *VideoCategoryController) GetTree(ctx *gin.Context) {
	includeDisabled := ctx.Query("include_disabled") == "true"

	tree, err := c.CategoryService.GetTree(includeDisabled)
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, tree)
}

// GetByID 获取分类详情
func (c *VideoCategoryController) GetByID(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	category, err := c.CategoryService.GetByID(id)
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, category)
}

// GetBySlug 根据Slug获取分类
func (c *VideoCategoryController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		response.Error(ctx, 400, "参数错误")
		return
	}

	category, err := c.CategoryService.GetBySlug(slug)
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, category)
}

// GetCategoryVideos 获取分类下的视频数量
func (c *VideoCategoryController) GetCategoryVideos(ctx *gin.Context) {
	id, ok := parseIDParam(ctx)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 使用VideoService查询该分类下的视频
	videoSvc := service.NewVideoService()
	status := 1 // 只显示已发布的
	params := service.VideoQueryParams{
		Page:       page,
		PageSize:   pageSize,
		CategoryID: id,
		Status:     &status,
	}

	result, err := videoSvc.GetVideoList(params)
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, result)
}
