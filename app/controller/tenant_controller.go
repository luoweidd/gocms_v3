package controller

import (
	"strconv"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

type TenantController struct {
	tenantService *service.TenantService
}

func NewTenantController() *TenantController {
	return &TenantController{
		tenantService: service.NewTenantService(),
	}
}

// CreateTenant 创建租户 API: POST /tenants
func (ctrl *TenantController) CreateTenant(c *gin.Context) {
	var req model.TenantCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	tenant, err := ctrl.tenantService.CreateTenant(req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, tenant)
}

// GetTenantList 获取租户列表 API: GET /tenants
func (ctrl *TenantController) GetTenantList(c *gin.Context) {
	var query model.TenantListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	pageData, err := ctrl.tenantService.GetTenantList(query)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, pageData)
}

// GetTenant 获取租户详情 API: GET /tenants/:id
func (ctrl *TenantController) GetTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	tenant, err := ctrl.tenantService.GetTenant(uint(id))
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, tenant)
}

// UpdateTenant 更新租户 API: PUT /tenants/:id
func (ctrl *TenantController) UpdateTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var req model.TenantUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if err := ctrl.tenantService.UpdateTenant(uint(id), req); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "更新成功", nil)
}

// DeleteTenant 删除租户 API: DELETE /tenants/:id
func (ctrl *TenantController) DeleteTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.tenantService.DeleteTenant(uint(id)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "删除成功", nil)
}

// UpdateTenantStatus 更新租户状态 API: PUT /tenants/:id/status
func (ctrl *TenantController) UpdateTenantStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if err := ctrl.tenantService.UpdateTenantStatus(uint(id), req.Status); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "状态更新成功", nil)
}

// AddUserToTenant 添加用户到租户 API: POST /tenants/:id/users
func (ctrl *TenantController) AddUserToTenant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if err := ctrl.tenantService.AddUserToTenant(uint(id), req.UserID); err != nil {
		response.Error(c, 400, err.Error())
		return
	}
	response.SuccessWithMsg(c, "添加成功", nil)
}

// GetTenantUsers 获取租户用户列表 API: GET /tenants/:id/users
func (ctrl *TenantController) GetTenantUsers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var query struct {
		Page     int `form:"page" binding:"min=1"`
		PageSize int `form:"page_size" binding:"min=1,max=100"`
	}
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	pageData, err := ctrl.tenantService.GetTenantUsers(uint(id), query.Page, query.PageSize)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, pageData)
}

// RemoveUserFromTenant 从租户移除用户 API: DELETE /tenants/:id/users/:user_id
func (ctrl *TenantController) RemoveUserFromTenant(c *gin.Context) {
	idStr := c.Param("id")
	userIDStr := c.Param("user_id")
	tenantID, err := strconv.ParseUint(idStr, 10, 32)
	userID, err2 := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || err2 != nil || tenantID == 0 || userID == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	if err := ctrl.tenantService.RemoveUserFromTenant(uint(tenantID), uint(userID)); err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.SuccessWithMsg(c, "移除成功", nil)
}

// GetTenantStats 获取租户统计 API: GET /tenants/stats
func (ctrl *TenantController) GetTenantStats(c *gin.Context) {
	stats, err := ctrl.tenantService.GetTenantStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}
