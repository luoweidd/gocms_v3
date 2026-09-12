package controller

import (
	"strconv"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
)

// ==================== 标签控制器 ====================

type TagController struct{}

func NewTagController() *TagController {
	return &TagController{}
}

// ListTags 获取标签列表（P0 修复：补充缺失的标签管理接口）
func (c *TagController) ListTags(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	var tags []model.Tag
	offset := (page - 1) * pageSize

	query := db.GetDB().Model(&model.Tag{})

	// 搜索条件
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	var total int64
	if err := query.Model(&model.Tag{}).Count(&total).Error; err != nil {
		response.Error(ctx, 500, "查询总数失败")
		return
	}

	response.Success(ctx, gin.H{
		"list":      tags,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetTag 获取标签详情（P0 修复：补充缺失的标签管理接口）
func (c *TagController) GetTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var tag model.Tag
	if err := db.GetDB().Where("id = ?", id).First(&tag).Error; err != nil {
		response.Error(ctx, 404, "标签不存在")
		return
	}

	response.Success(ctx, tag)
}

// CreateTag 创建标签（P0 修复：补充缺失的标签管理接口）
func (c *TagController) CreateTag(ctx *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		Status int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	// 检查名称是否已存在
	var count int64
	db.GetDB().Model(&model.Tag{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		response.Error(ctx, 409, "标签名称已存在")
		return
	}

	tag := model.Tag{
		Name:   req.Name,
		Status: req.Status,
	}

	if err := db.GetDB().Create(&tag).Error; err != nil {
		response.Error(ctx, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(ctx, "创建成功", gin.H{"id": tag.ID})
}

// UpdateTag 更新标签（P0 修复：补充缺失的标签管理接口）
func (c *TagController) UpdateTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var tag model.Tag
	if err := db.GetDB().Where("id = ?", id).First(&tag).Error; err != nil {
		response.Error(ctx, 404, "标签不存在")
		return
	}

	var req struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" && req.Name != tag.Name {
		// 检查新名称是否已存在
		var count int64
		db.GetDB().Model(&model.Tag{}).Where("name = ? AND id != ?", req.Name, id).Count(&count)
		if count > 0 {
			response.Error(ctx, 409, "标签名称已存在")
			return
		}
		updates["name"] = req.Name
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		db.GetDB().Model(&tag).Updates(updates)
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// DeleteTag 删除标签（P0 修复：补充缺失的标签管理接口）
func (c *TagController) DeleteTag(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	if err := db.GetDB().Delete(&model.Tag{}, id).Error; err != nil {
		response.Error(ctx, 500, "删除失败")
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}
