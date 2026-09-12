package model

import "time"

// APISignKey API签名密钥模型（对应 api_sign_keys 表）
type APISignKey struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	AppID     string    `gorm:"size:64;not null;uniqueIndex:uk_api_sign_app_id" json:"app_id"` // 应用ID
	SecretKey string    `gorm:"size:128;not null" json:"-"`                                    // 密钥（不返回）
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`                                    // 过期时间
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (APISignKey) TableName() string {
	return "api_sign_keys"
}

// APISignKeyCreateRequest 创建密钥请求
type APISignKeyCreateRequest struct {
	AppID      string `json:"app_id" binding:"required,max:64"`
	SecretKey  string `json:"secret_key" binding:"required,max:128"`
	ExpireDays int    `json:"expire_days"`
}

// APISignKeyVO 密钥视图对象（隐藏 secret_key）
type APISignKeyVO struct {
	ID        uint      `json:"id"`
	AppID     string    `json:"app_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired 检查是否过期
func (k *APISignKey) IsExpired() bool {
	return k.ExpiresAt.Before(time.Now())
}
