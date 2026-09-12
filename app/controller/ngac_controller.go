package controller

import (
	"net/http"
	"strconv"

	"gocms_v3/app/model/ngac"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

type NgacController struct {
	svc *service.NgacService
}

func NewNgacController() *NgacController {
	return &NgacController{svc: service.NewNgacService()}
}

// ==================== 角色管理 ====================

func (c *NgacController) CreateRole(ctx *gin.Context) {
	var req ngac.CreateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	role, err := c.svc.CreateRole(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": role})
}

func (c *NgacController) ListRoles(ctx *gin.Context) {
	var req ngac.RoleListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}
	roles, total, err := c.svc.ListRoles(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": roles, "total": total, "page": req.Page, "pageSize": req.PageSize})
}

func (c *NgacController) GetRole(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	role, err := c.svc.GetRole(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": role})
}

func (c *NgacController) UpdateRole(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req ngac.UpdateRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.UpdateRole(uint(id), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (c *NgacController) DeleteRole(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.DeleteRole(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (c *NgacController) BuildRoleTree(ctx *gin.Context) {
	tree, err := c.svc.BuildRoleTree()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": tree})
}

// ==================== 权限管理 ====================

func (c *NgacController) CreatePermission(ctx *gin.Context) {
	var req ngac.CreatePermissionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	perm, err := c.svc.CreatePermission(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": perm})
}

func (c *NgacController) GetPermission(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	perm, err := c.svc.GetPermission(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": perm})
}

func (c *NgacController) UpdatePermission(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	var req ngac.UpdatePermissionReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.UpdatePermission(uint(id), req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (c *NgacController) DeletePermission(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err := c.svc.DeletePermission(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (c *NgacController) ListPermissions(ctx *gin.Context) {
	var req ngac.PermissionListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}
	perms, total, err := c.svc.ListPermissions(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": perms, "total": total, "page": req.Page, "pageSize": req.PageSize})
}

func (c *NgacController) AssignPermissionsToRole(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Param("role_id"), 10, 64)
	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.AssignPermissionsToRole(uint(roleID), req.PermissionIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "assigned"})
}

func (c *NgacController) RemovePermissionFromRole(ctx *gin.Context) {
	roleID, _ := strconv.ParseUint(ctx.Param("role_id"), 10, 64)
	permID, _ := strconv.ParseUint(ctx.Param("perm_id"), 10, 64)
	if err := c.svc.RemovePermissionFromRole(uint(roleID), uint(permID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "removed"})
}

// ==================== 用户 - 角色管理 ====================

func (c *NgacController) AssignRoles(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	var req ngac.AssignRolesReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.AssignRolesToUser(uint(userID), req.RoleIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "roles assigned"})
}

func (c *NgacController) GetUserRoles(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	roles, err := c.svc.GetUserRoles(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": roles})
}

func (c *NgacController) GetUserPermissions(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	perms, err := c.svc.GetUserPermissions(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": perms})
}

// ==================== 对象类管理 ====================

func (c *NgacController) GetObjectClassTree(ctx *gin.Context) {
	tree, err := c.svc.GetObjectClassTree()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": tree})
}

// ==================== 属性管理 ====================

func (c *NgacController) CreateUserAttribute(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	var req ngac.CreateUserAttributeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.svc.CreateUserAttribute(uint(userID), req.Name, req.Value, req.Type); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "created"})
}

func (c *NgacController) GetSubjectAttributes(ctx *gin.Context) {
	userID, _ := strconv.ParseUint(ctx.Param("user_id"), 10, 64)
	attrs := c.svc.GetSubjectAttributes(uint(userID))
	ctx.JSON(http.StatusOK, gin.H{"data": attrs})
}

// ==================== 上下文属性管理 ====================

func (c *NgacController) CreateContextAttribute(ctx *gin.Context) {
	var req ngac.CreateContextAttributeReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctxAttr, err := c.svc.CreateContextAttribute(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": ctxAttr})
}

func (c *NgacController) GetContextAttributes(ctx *gin.Context) {
	attrs := c.svc.GetContextAttributes()
	ctx.JSON(http.StatusOK, gin.H{"data": attrs})
}

// ==================== 授权接口（PDP） ====================

func (c *NgacController) Authz(ctx *gin.Context) {
	var req ngac.AuthzRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := service.EvaluateAuthz(req)
	ctx.JSON(http.StatusOK, resp)
}
