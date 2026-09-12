package controller

import (
	"strconv"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// MessageCategoryController 消息类别控制器
type MessageCategoryController struct {
	service *service.MessageCategoryService
}

// NewMessageCategoryController 创建消息类别控制器
func NewMessageCategoryController(service *service.MessageCategoryService) *MessageCategoryController {
	return &MessageCategoryController{service: service}
}

// GetCategories 获取消息类别列表
func (c *MessageCategoryController) GetCategories(ginCtx *gin.Context) {
	query := model.MessageCategoryListQuery{
		Page:     1,
		PageSize: 20,
	}

	pageStr := ginCtx.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			query.Page = p
		}
	}

	pageSizeStr := ginCtx.Query("page_size")
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			query.PageSize = ps
		}
	}

	query.Name = ginCtx.Query("name")
	query.Code = ginCtx.Query("code")

	statusStr := ginCtx.Query("status")
	if statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			si := int8(s)
			query.Status = &si
		}
	}

	categories, total, err := c.service.GetCategories(query)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, gin.H{
		"list":     categories,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetCategoryByID 根据ID获取消息类别
func (c *MessageCategoryController) GetCategoryByID(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的类别ID")
		return
	}

	cat, err := c.service.GetCategoryByID(uint(id))
	if err != nil {
		response.Error(ginCtx, 404, err.Error())
		return
	}

	response.Success(ginCtx, cat)
}

// CreateCategory 创建消息类别
func (c *MessageCategoryController) CreateCategory(ginCtx *gin.Context) {
	var req model.CreateMessageCategoryRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	cat, err := c.service.CreateCategory(req)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "创建成功", cat)
}

// UpdateCategory 更新消息类别
func (c *MessageCategoryController) UpdateCategory(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的类别ID")
		return
	}

	var req model.UpdateMessageCategoryRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	cat, err := c.service.UpdateCategory(uint(id), req)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "更新成功", cat)
}

// DeleteCategory 删除消息类别
func (c *MessageCategoryController) DeleteCategory(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的类别ID")
		return
	}

	if err := c.service.DeleteCategory(uint(id)); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "删除成功", nil)
}
