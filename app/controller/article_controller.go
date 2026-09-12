package controller

import (
	"fmt"
	"strconv"

	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	ArticleService *service.ArticleService
}

func NewArticleController() *ArticleController {
	return &ArticleController{
		ArticleService: service.NewArticleService(),
	}
}

// Create 创建文章
func (c *ArticleController) Create(ctx *gin.Context) {
	var req service.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	authorID := ctx.GetUint("user_id")
	req.AuthorID = authorID

	if err := c.ArticleService.CreateArticle(req); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建成功", nil)
}

// Update 更新文章
func (c *ArticleController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req service.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.ArticleService.UpdateArticle(uint(id), req); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// Delete 删除文章
func (c *ArticleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := c.ArticleService.DeleteArticle(uint(id)); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetList 获取文章列表（分页）
func (c *ArticleController) GetList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	result, err := c.ArticleService.GetArticleList(page, pageSize)
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, result)
}

// GetDetail 获取文章详情
func (c *ArticleController) GetDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	article, err := c.ArticleService.GetArticleByID(uint(id))
	if err != nil {
		response.Error(ctx, 404, "文章不存在")
		return
	}

	response.Success(ctx, article)
}

// GetCategories 获取分类列表
func (c *ArticleController) GetCategories(ctx *gin.Context) {
	categories, err := service.GetCategoryList()
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, categories)
}

// UploadCover 上传封面图
func (c *ArticleController) UploadCover(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.Error(ctx, 400, "请选择文件")
		return
	}

	userID := ctx.GetUint("user_id")
	tenantID := ctx.GetUint("tenant_id")

	url, err := service.UploadFile(file, userID, tenantID)
	if err != nil {
		response.Error(ctx, 500, "上传失败")
		return
	}

	response.Success(ctx, gin.H{"url": url})
}

// GetTags 获取标签列表
func (c *ArticleController) GetTags(ctx *gin.Context) {
	tags, err := service.GetTagList()
	if err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	response.Success(ctx, tags)
}

// Helper functions

func parseUintParam(ctx *gin.Context, name string) (uint, error) {
	idStr := ctx.Param(name)
	if idStr == "" {
		return 0, fmt.Errorf("参数 %s 不能为空", name)
	}
	var id int
	fmt.Sscanf(idStr, "%d", &id)
	return uint(id), nil
}
