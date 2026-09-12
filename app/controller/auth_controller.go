package controller

import (
	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/model/ngac"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	UserService *service.UserService
}

func NewAuthController() *AuthController {
	return &AuthController{
		UserService: service.NewUserService(),
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req service.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	result, err := c.UserService.Login(req)
	if err != nil {
		// 登录失败返回401认证错误码，而非500服务器内部错误
		response.Error(ctx, 401, "用户名或密码错误")
		return
	}

	response.Success(ctx, result)
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	var req service.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	result, err := c.UserService.Register(req)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, result)
}

// Logout 用户登出（前端删除 Token 即可）
func (c *AuthController) Logout(ctx *gin.Context) {
	response.SuccessWithMsg(ctx, "已退出登录", nil)
}

// GetUserInfo 获取当前用户信息
func (c *AuthController) GetUserInfo(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var user model.User

	if err := db.GetDB().First(&user, userID).Error; err != nil {
		response.Error(ctx, 404, "用户不存在")
		return
	}

	// 移除密码字段
	user.Password = ""
	response.Success(ctx, user)
}

// GetUserPermissions 获取当前用户权限列表
func (c *AuthController) GetUserPermissions(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	// 查询用户角色
	var userRoles []ngac.UserRole
	if err := db.GetDB().Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		response.Error(ctx, 500, "查询用户角色失败")
		return
	}

	if len(userRoles) == 0 {
		response.Success(ctx, gin.H{
			"permissions": []string{},
			"roles":       []string{},
		})
		return
	}

	// 提取角色 ID
	var roleIDs []uint
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}

	// 查询角色名称和权限码
	var roleNames []string
	var permissions []string

	for _, roleID := range roleIDs {
		var role ngac.Role
		if err := db.GetDB().Where("id = ?", roleID).First(&role).Error; err == nil {
			roleNames = append(roleNames, role.Name)
			permissions = append(permissions, role.Code)
		}
	}

	response.Success(ctx, gin.H{
		"permissions": permissions,
		"roles":       roleNames,
	})
}
