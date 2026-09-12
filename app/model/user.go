package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// User 用户模型（后台管理员）
type User struct {
	BaseModel
	TenantID       *uint      `gorm:"index;column:tenant_id" json:"tenant_id"` // 所属租户ID
	Username       string     `gorm:"uniqueIndex:uk_user_username;size:50;not null" json:"username"`
	Password       string     `gorm:"size:128;not null" json:"-"`              // 不返回密码
	Nickname       string     `gorm:"size:50;column:nickname" json:"nickname"` // 昵称
	Email          string     `gorm:"uniqueIndex:uk_user_email;size:100;index:idx_user_email" json:"email"`
	Phone          string     `gorm:"size:20;column:phone" json:"phone"`                                       // 手机号
	Avatar         string     `gorm:"size:255;column:avatar" json:"avatar"`                                    // 头像URL
	WechatOpenID   *string    `gorm:"size:64;index:idx_user_wechat;column:wechat_openid" json:"wechat_openid"` // 微信OpenID
	QQOpenID       *string    `gorm:"size:64;index:idx_user_qq;column:qq_openid" json:"qq_openid"`             // QQ OpenID
	GithubOpenID   *string    `gorm:"size:64;index:idx_user_github;column:github_openid" json:"github_openid"` // GitHub OpenID
	PhoneEncrypted []byte     `gorm:"column:phone_encrypted" json:"phone_encrypted"`                           // 加密手机号
	EmailEncrypted []byte     `gorm:"column:email_encrypted" json:"email_encrypted"`                           // 加密邮箱
	Status         int        `gorm:"default:1;index:idx_user_status" json:"status"`                           // 1启用 0禁用
	Roles          JSONString `gorm:"type:text;column:roles" json:"roles"`                                     // JSON数组存储角色
	LastLoginAt    *time.Time `gorm:"column:last_login_at" json:"last_login_at"`                               // 最后登录时间
}

func (User) TableName() string {
	return "users"
}

// JSONString 是用于 GORM 的 JSON 字符串字段的自定义类型
type JSONString []string

// Scan 实现 sql.Scanner interface
func (j *JSONString) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan JSONString: unsupported type %T", value)
	}
	trimmed := bytes
	// 去掉外层方括号和引号，兼容数据库存的普通字符串
	str := string(trimmed)
	if str == "" {
		*j = []string{}
		return nil
	}
	// 如果是 JSON 数组格式 ["xxx","yyy"]
	if len(str) >= 2 && str[0] == '[' && str[len(str)-1] == ']' {
		return json.Unmarshal(bytes, j)
	}
	// 否则作为单个字符串
	*j = []string{str}
	return nil
}

// Value 实现 driver.Valuer interface
func (j JSONString) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
