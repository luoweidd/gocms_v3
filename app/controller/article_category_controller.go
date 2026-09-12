package controller

import (
	"strconv"

	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ArticleCategoryController struct {
	categoryService *service.ArticleCategoryService
}

func NewArticleCategoryController() *ArticleCategoryController {
	return &ArticleCategoryController{
		categoryService: service.NewArticleCategoryService(),
	}
}

// Create 创建文章分类
func (c *ArticleCategoryController) Create(ctx *gin.Context) {
	var req service.CreateArticleCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	category, err := c.categoryService.CreateArticleCategory(req)
	if err != nil {
		zap.L().Error("创建文章分类失败", zap.Error(err))
		response.Error(ctx, 500, "创建失败")
		return
	}

	zap.L().Info("创建文章分类成功", zap.String("name", req.Name))
	response.Success(ctx, category)
}

// Update 更新文章分类
func (c *ArticleCategoryController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req service.UpdateArticleCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	category, err := c.categoryService.UpdateArticleCategory(uint(id), req)
	if err != nil {
		zap.L().Error("更新文章分类失败", zap.Error(err))
		response.Error(ctx, 500, "更新失败")
		return
	}

	zap.L().Info("更新文章分类成功", zap.String("id", idStr))
	response.Success(ctx, category)
}

// Delete 删除文章分类
func (c *ArticleCategoryController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.categoryService.DeleteArticleCategory(uint(id)); err != nil {
		zap.L().Error("删除文章分类失败", zap.Error(err))
		response.Error(ctx, 500, "删除失败")
		return
	}

	zap.L().Info("删除文章分类成功", zap.String("id", idStr))
	response.Success(ctx, nil)
}

// GetList 获取文章分类列表
func (c *ArticleCategoryController) GetList(ctx *gin.Context) {
	var req service.ListArticleCategoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	categories, err := c.categoryService.ListArticleCategory(req)
	if err != nil {
		zap.L().Error("获取文章分类列表失败", zap.Error(err))
		response.Error(ctx, 500, "查询失败")
		return
	}

	zap.L().Info("获取文章分类列表成功")
	response.Success(ctx, categories)
}

// GetTree 获取文章分类树
func (c *ArticleCategoryController) GetTree(ctx *gin.Context) {
	var req service.TreeArticleCategoryRequest
	// 使用ShouldBindQuery而不严格要求binding标签
	if err := ctx.ShouldBindQuery(&req); err != nil {
		zap.L().Warn("绑定查询参数失败（可能没有提供参数）", zap.Error(err))
		// 允许空查询参数，继续执行
	}

	tree, err := c.categoryService.TreeArticleCategory(req)
	if err != nil {
		zap.L().Error("获取文章分类树失败", zap.Error(err))
		response.Error(ctx, 500, "查询失败")
		return
	}

	zap.L().Info("获取文章分类树成功")
	response.Success(ctx, tree)
}

// GetByID 获取文章分类详情
func (c *ArticleCategoryController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	category, err := c.categoryService.GetArticleCategoryByID(uint(id))
	if err != nil {
		zap.L().Error("获取文章分类详情失败", zap.Error(err))
		response.Error(ctx, 500, "查询失败")
		return
	}

	zap.L().Info("获取文章分类详情成功", zap.String("id", idStr))
	response.Success(ctx, category)
}
