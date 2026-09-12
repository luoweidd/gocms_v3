package dto

import "gocms_v3/app/domain/entity"

// LoginRequest 登录请求DTO
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginWithOAuth2Request OAuth2登录请求DTO
type LoginWithOAuth2Request struct {
	Provider string `json:"provider" binding:"required,oneof=wechat qq github"` // wechat, qq, github
	Code     string `json:"code" binding:"required"`
	State    string `json:"state"`
}

// RegisterRequest 注册请求DTO
type RegisterRequest struct {
	Username string   `json:"username" binding:"required,min=3,max=50"`
	Password string   `json:"password" binding:"required,min=8"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email" binding:"omitempty,email"`
	Phone    string   `json:"phone" binding:"omitempty"`
	Roles    []string `json:"roles"`
	TenantID uint     `json:"tenant_id"`
}

// LoginResponse 登录响应DTO
type LoginResponse struct {
	Token    string       `json:"token"`
	User     *entity.User `json:"user"`
	OAuth2ID string       `json:"oauth2_id,omitempty"` // OAuth2用户唯一标识
	Provider string       `json:"provider,omitempty"`
}

// RefreshTokenRequest 刷新Token请求DTO
type RefreshTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	NewToken string `json:"new_token"`
}

// ChangePasswordRequest 修改密码请求DTO
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
