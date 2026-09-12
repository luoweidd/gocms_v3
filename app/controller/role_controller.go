package controller

import (
	"fmt"
	"strconv"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
)

// PermissionNode 权限节点
type PermissionNode struct {
	ID       uint             `json:"id"`
	Name     string           `json:"name"`
	Code     string           `json:"code"`
	Type     string           `json:"type"` // menu/button
	ParentID uint             `json:"parent_id"`
	Children []PermissionNode `json:"children,omitempty"`
}

// ==================== 角色控制器 ====================

type RoleController struct{}

func NewRoleController() *RoleController {
	return &RoleController{}
}

// ListRoles 获取角色列表（P0 修复：补充缺失的角色管理接口）
func (c *RoleController) ListRoles(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var roles []model.Role
	offset := (page - 1) * pageSize

	query := db.GetDB().Model(&model.Role{})

	// 搜索条件
	if keyword := ctx.Query("keyword"); keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status := ctx.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&roles).Error; err != nil {
		response.Error(ctx, 500, "查询失败")
		return
	}

	var total int64
	if err := query.Model(&model.Role{}).Count(&total).Error; err != nil {
		response.Error(ctx, 500, "查询总数失败")
		return
	}

	response.Success(ctx, gin.H{
		"list":      roles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetRole 获取角色详情（P0 修复：补充缺失的角色管理接口）
func (c *RoleController) GetRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var role model.Role
	if err := db.GetDB().Where("id = ?", id).First(&role).Error; err != nil {
		response.Error(ctx, 404, "角色不存在")
		return
	}

	response.Success(ctx, role)
}

// CreateRole 创建角色（P0 修复：补充缺失的角色管理接口）
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	// 检查 code 是否已存在
	var count int64
	db.GetDB().Model(&model.Role{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		response.Error(ctx, 409, "角色编码已存在")
		return
	}

	role := model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := db.GetDB().Create(&role).Error; err != nil {
		response.Error(ctx, 500, "创建失败")
		return
	}

	response.SuccessWithMsg(ctx, "创建成功", gin.H{"id": role.ID})
}

// UpdateRole 更新角色（P0 修复：补充缺失的角色管理接口）
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var role model.Role
	if err := db.GetDB().Where("id = ?", id).First(&role).Error; err != nil {
		response.Error(ctx, 404, "角色不存在")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      int    `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	// 如果修改了 code，检查是否已存在
	if req.Code != "" && req.Code != role.Code {
		var count int64
		db.GetDB().Model(&model.Role{}).Where("code = ? AND id != ?", req.Code, id).Count(&count)
		if count > 0 {
			response.Error(ctx, 409, "角色编码已存在")
			return
		}
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Status >= 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		db.GetDB().Model(&role).Updates(updates)
	}

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// DeleteRole 删除角色（P0 修复：补充缺失的角色管理接口）
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	// 检查是否有关联的用户
	var userCount int64
	db.GetDB().Model(&model.User{}).Where("roles LIKE ?", "%"+strconv.FormatUint(id, 10)+"%").Count(&userCount)
	if userCount > 0 {
		response.Error(ctx, 403, "该角色下有关联用户，无法删除")
		return
	}

	if err := db.GetDB().Delete(&model.Role{}, id).Error; err != nil {
		response.Error(ctx, 500, "删除失败")
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}

// GetPermissionTree 获取权限树（P2 修复：补充前端 getPermissionTree 调用的接口）
func (c *RoleController) GetPermissionTree(ctx *gin.Context) {
	// 从菜单表获取所有启用的菜单作为权限节点
	var menus []model.Menu
	if err := db.GetDB().Where("status = 1").Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		fmt.Printf("[权限树] 查询失败: %v\n", err)
		response.Error(ctx, 500, "查询权限树失败")
		return
	}

	// 构建权限树
	nodeMap := make(map[uint]*PermissionNode)
	var roots []PermissionNode

	for _, m := range menus {
		// 安全获取 ParentID 值
		getParentID := func(pid *uint) uint {
			if pid != nil {
				return *pid
			}
			return 0
		}
		node := &PermissionNode{
			ID:       m.ID,
			Name:     m.Name,
			Code:     m.Permission,
			Type:     mapMenuType(m.MenuType),
			ParentID: getParentID(m.ParentID),
		}
		nodeMap[m.ID] = node

		if m.ParentID == nil || *m.ParentID == 0 {
			roots = append(roots, *node)
		}
	}

	// 构建父子关系
	for _, node := range nodeMap {
		if node.ParentID != 0 {
			if parent, exists := nodeMap[node.ParentID]; exists {
				parent.Children = append(parent.Children, *node)
			}
		}
	}

	response.Success(ctx, roots)
}

// mapMenuType 映射菜单类型为权限类型
func mapMenuType(menuType int) string {
	switch menuType {
	case 1:
		return "menu"
	case 2:
		return "directory"
	case 3:
		return "button"
	default:
		return "menu"
	}
}
