package controller

import (
	"strconv"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MenuController struct {
	menuService *service.MenuService
}

func NewMenuController() *MenuController {
	return &MenuController{
		menuService: service.NewMenuService(),
	}
}

// Create 创建菜单
func (ctrl *MenuController) Create(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required,min=1,max=100"`
		Path       string `json:"path" binding:"required"`
		Component  string `json:"component"`
		Icon       string `json:"icon"`
		SortOrder  int    `json:"sort_order"`
		ParentID   uint   `json:"parent_id"`
		MenuType   int    `json:"menu_type" binding:"required,oneof=1 2 3"`
		Permission string `json:"permission"`
		Visible    int    `json:"visible"`
		Status     int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	menu := &model.Menu{
		Name:       req.Name,
		Path:       req.Path,
		Component:  req.Component,
		Icon:       req.Icon,
		SortOrder:  req.SortOrder,
		ParentID:   &req.ParentID,
		MenuType:   req.MenuType,
		Permission: req.Permission,
		Visible:    req.Visible,
		Status:     req.Status,
	}

	if err := ctrl.menuService.CreateMenu(menu); err != nil {
		zap.L().Error("创建菜单失败", zap.Error(err))
		response.Error(c, 500, "创建失败")
		return
	}

	zap.L().Info("创建菜单成功", zap.String("name", req.Name))
	response.Success(c, menu)
}

// Update 更新菜单
func (ctrl *MenuController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var req struct {
		Name       *string `json:"name"`
		Path       *string `json:"path"`
		Component  *string `json:"component"`
		Icon       *string `json:"icon"`
		SortOrder  *int    `json:"sort_order"`
		ParentID   *uint   `json:"parent_id"`
		MenuType   *int    `json:"menu_type"`
		Permission *string `json:"permission"`
		Visible    *int    `json:"visible"`
		Status     *int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误: "+err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Path != nil {
		updates["path"] = *req.Path
	}
	if req.Component != nil {
		updates["component"] = *req.Component
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.MenuType != nil {
		updates["menu_type"] = *req.MenuType
	}
	if req.Permission != nil {
		updates["permission"] = *req.Permission
	}
	if req.Visible != nil {
		updates["visible"] = *req.Visible
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := ctrl.menuService.UpdateMenu(uint(id), updates); err != nil {
		zap.L().Error("更新菜单失败", zap.Error(err))
		response.Error(c, 500, "更新失败")
		return
	}

	zap.L().Info("更新菜单成功", zap.String("id", idStr))
	response.Success(c, nil)
}

// Delete 删除菜单
func (ctrl *MenuController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.menuService.DeleteMenu(uint(id)); err != nil {
		zap.L().Error("删除菜单失败", zap.Error(err))
		response.Error(c, 500, "删除失败")
		return
	}

	zap.L().Info("删除菜单成功", zap.String("id", idStr))
	response.Success(c, nil)
}

// GetByID 获取菜单详情
func (ctrl *MenuController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	menu, err := ctrl.menuService.GetMenuByID(uint(id))
	if err != nil {
		zap.L().Error("获取菜单详情失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, menu)
}

// GetList 获取菜单列表
func (ctrl *MenuController) GetList(c *gin.Context) {
	var req struct {
		Status   *int `form:"status"`
		MenuType *int `form:"menu_type"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误")
		return
	}

	menus, err := ctrl.menuService.ListMenus(req.Status, req.MenuType)
	if err != nil {
		zap.L().Error("获取菜单列表失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, menus)
}

// GetTree 获取菜单树
func (ctrl *MenuController) GetTree(c *gin.Context) {
	tree, err := ctrl.menuService.GetMenuTree()
	if err != nil {
		zap.L().Error("获取菜单树失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, tree)
}

// GetPermissions 获取权限列表
func (ctrl *MenuController) GetPermissions(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		role = "admin"
	}

	permissions, err := ctrl.menuService.GetPermissionsByRole(role)
	if err != nil {
		zap.L().Error("获取权限列表失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, permissions)
}
