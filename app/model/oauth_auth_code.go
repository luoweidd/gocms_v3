package model

import "time"

// OAuthAuthCode OAuth授权码模型（对应 oauth_auth_codes 表）
type OAuthAuthCode struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint      `gorm:"index;not null" json:"user_id"`    // 用户ID
	Provider    string    `gorm:"size:32;not null" json:"provider"` // 第三方平台: wechat/qq/github
	OpenID      string    `gorm:"size:64;not null" json:"open_id"`  // 第三方OpenID
	Code        string    `gorm:"size:64;uniqueIndex" json:"code"`  // 授权码
	RedirectURI string    `gorm:"size:255" json:"redirect_uri"`     // 回调地址
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`       // 过期时间
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 指定表名
func (OAuthAuthCode) TableName() string {
	return "oauth_auth_codes"
}

// OAuthAuthCodeCreateRequest 创建授权码请求
type OAuthAuthCodeCreateRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Provider    string `json:"provider" binding:"required,oneof=wechat qq github"`
	OpenID      string `json:"open_id" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
	ExpireHours int    `json:"expire_hours"`
}

// OAuthAuthCodeVO 授权码视图对象
type OAuthAuthCodeVO struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Provider    string    `json:"provider"`
	OpenID      string    `json:"open_id"`
	RedirectURI string    `json:"redirect_uri"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}
