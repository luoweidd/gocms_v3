package service

import (
	"fmt"
	"log"
	"time"

	"gocms_v3/app/db"
	middleware2 "gocms_v3/app/middleware"
	"gocms_v3/app/model"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Password string   `json:"password" binding:"required,min=6"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email" binding:"omitempty,email"`
	Phone    string   `json:"phone"`
	Roles    []string `json:"roles"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string   `json:"token"`
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Roles    []string `json:"roles"`
}

// Login 用户登录
func (s *UserService) Login(req LoginRequest) (*LoginResponse, error) {
	log.Printf("[SERVICE] 用户登录请求 - 用户名: %s", req.Username)

	var user model.User
	if err := db.GetDB().Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
		log.Printf("[SERVICE-WARN] 用户登录失败 - 用户名: %s, 原因: 用户不存在或已禁用", req.Username)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	if !checkPassword(user.Password, req.Password) {
		log.Printf("[SERVICE-WARN] 用户登录失败 - 用户名: %s, 原因: 密码错误", req.Username)
		return nil, fmt.Errorf("用户名或密码错误")
	}

	now := time.Now()
	db.GetDB().Model(&user).Update("last_login_at", now)

	token, err := middleware2.GenerateToken(user.ID, user.Username, []string(user.Roles))
	if err != nil {
		log.Printf("[SERVICE-ERROR] 生成令牌失败 - 用户ID: %d, 错误: %v", user.ID, err)
		return nil, fmt.Errorf("生成令牌失败")
	}

	log.Printf("[SERVICE] 用户登录成功 - 用户ID: %d, 用户名: %s", user.ID, user.Username)
	return &LoginResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Roles:    []string(user.Roles),
	}, nil
}

// Register 用户注册
func (s *UserService) Register(req RegisterRequest) (*model.User, error) {
	log.Printf("[SERVICE] 用户注册请求 - 用户名: %s", req.Username)

	var count int64
	db.GetDB().Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		log.Printf("[SERVICE-WARN] 用户注册失败 - 用户名: %s, 原因: 用户名已存在", req.Username)
		return nil, fmt.Errorf("用户名已存在")
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("[SERVICE-ERROR] 密码加密失败 - 用户名: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("密码加密失败")
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
		Roles:    model.JSONString(req.Roles),
	}

	if err := db.GetDB().Create(&user).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 用户注册失败 - 用户名: %s, 错误: %v", req.Username, err)
		return nil, fmt.Errorf("注册失败")
	}

	log.Printf("[SERVICE] 用户注册成功 - 用户ID: %d, 用户名: %s", user.ID, user.Username)
	return &user, nil
}

// ConvertRoles 将 JSONString 转为 []string
func ConvertRoles(j model.JSONString) []string {
	return []string(j)
}

// hashPassword 使用bcrypt标准库哈希密码
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// checkPassword 验证密码
func checkPassword(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
