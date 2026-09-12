package entity

import (
	"fmt"
	"time"
)

// User 用户领域实体
type User struct {
	ID             uint       `json:"id"`
	TenantID       uint       `json:"tenant_id" gorm:"default:0"` // 多租户支持
	Username       string     `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash   string     `json:"-" gorm:"size:128;not null"` // 密码哈希值，不返回
	Nickname       string     `json:"nickname" gorm:"size:50"`
	Email          string     `json:"email" gorm:"uniqueIndex;size:100"`
	EmailEncrypted string     `json:"-" gorm:"type:blob"` // 加密邮箱
	Phone          string     `json:"phone" gorm:"uniqueIndex;size:20"`
	PhoneEncrypted string     `json:"-" gorm:"type:blob"` // 加密手机号
	Avatar         string     `json:"avatar" gorm:"size:255"`
	Status         int        `json:"status" gorm:"default:1"` // 1启用 0禁用
	Roles          []string   `json:"roles" gorm:"type:text"`  // JSON数组存储角色
	LastLoginAt    *time.Time `json:"last_login_at"`

	// OAuth2第三方登录绑定
	WechatOpenID string `json:"wechat_openid" gorm:"size:64"`
	QQOpenID     string `json:"qq_openid" gorm:"size:64"`
	GitHubOpenID string `json:"github_openid" gorm:"size:64"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// IsEnabled 检查用户是否启用
func (u *User) IsEnabled() bool {
	return u.Status == 1
}

// HasRole 检查用户是否有指定角色
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// CreateUserParams 创建用户参数
type CreateUserParams struct {
	Username     string
	PasswordHash string
	Nickname     string
	Email        string
	Phone        string
	Avatar       string
	Roles        []string
	TenantID     uint
}

// UpdateUserParams 更新用户参数
type UpdateUserParams struct {
	Nickname *string
	Email    *string
	Phone    *string
	Avatar   *string
	Status   *int
	Roles    *[]string
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireDigit   bool
	RequireSpecial bool
}

// DefaultPasswordPolicy 默认密码策略
var DefaultPasswordPolicy = PasswordPolicy{
	MinLength:    8,
	RequireUpper: true,
	RequireLower: true,
	RequireDigit: true,
}

// ValidatePassword 验证密码强度
func (p *PasswordPolicy) ValidatePassword(password string) error {
	if len(password) < p.MinLength {
		return fmt.Errorf("密码长度不能少于%d位", p.MinLength)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}

	if p.RequireUpper && !hasUpper {
		return fmt.Errorf("密码必须包含大写字母")
	}
	if p.RequireLower && !hasLower {
		return fmt.Errorf("密码必须包含小写字母")
	}
	if p.RequireDigit && !hasDigit {
		return fmt.Errorf("密码必须包含数字")
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
