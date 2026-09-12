package controller

import (
	"strconv"

	"gocms_v3/app/db"
	"gocms_v3/app/dto"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ==================== 用户控制器 ====================

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

// GetUserInfo 获取当前用户信息
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")

	var user model.User
	if err := db.GetDB().First(&user, userID).Error; err != nil {
		response.Error(ctx, 404, "用户不存在")
		return
	}

	user.Password = ""
	response.Success(ctx, user)
}

// UpdateUser 更新用户信息
func (c *UserController) UpdateUser(ctx *gin.Context) {
	userID := ctx.GetUint("user_id")
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	db.GetDB().Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"avatar":   req.Avatar,
	})

	response.SuccessWithMsg(ctx, "更新成功", nil)
}

// ListUsers 获取用户列表
func (c *UserController) ListUsers(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	var users []model.User
	offset := (page - 1) * pageSize

	db.GetDB().Offset(offset).Limit(pageSize).Find(&users)

	var total int64
	db.GetDB().Model(&model.User{}).Count(&total)

	response.Success(ctx, gin.H{
		"list":     users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// CreateUser 创建用户
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.UserCreateRequest

	// 绑定 JSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// 获取原始请求体用于调试
		rawBody, _ := ctx.GetRawData()
		response.Error(ctx, 400, "参数错误: "+err.Error()+", 请求数据: "+string(rawBody))
		return
	}

	// 检查用户名是否已存在
	var count int64
	db.GetDB().Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		response.Error(ctx, 409, "用户名已存在")
		return
	}

	// 处理角色：如果提供了 role_id，将其转换为 roles 数组
	if req.RoleID != nil && *req.RoleID > 0 {
		// 先检查是否已有 roles
		if req.Roles == nil {
			req.Roles = []string{}
		}
		// 查询角色名称
		var role model.Role
		if err := db.GetDB().First(&role, *req.RoleID).Error; err == nil {
			req.Roles = append(req.Roles, role.Name)
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(ctx, 500, "密码加密失败")
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Roles:    req.Roles,
		Status:   req.Status,
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		response.Error(ctx, 500, "创建失败: "+err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "创建成功", nil)
}

// GetByID 获取用户详情（P1 修复：补充缺失的用户详情接口）
func (c *UserController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var user model.User
	if err := db.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		response.Error(ctx, 404, "用户不存在")
		return
	}

	// 清除密码字段
	user.Password = ""
	response.Success(ctx, user)
}

// DeleteUser 删除用户（P1 修复：补充缺失的用户删除接口）
func (c *UserController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, 400, "参数错误")
		return
	}

	// 防止删除自己
	userID := ctx.GetUint("user_id")
	if uint(id) == userID {
		response.Error(ctx, 403, "不能删除当前登录用户")
		return
	}

	if err := db.GetDB().Delete(&model.User{}, id).Error; err != nil {
		response.Error(ctx, 500, "删除失败")
		return
	}

	response.SuccessWithMsg(ctx, "删除成功", nil)
}
