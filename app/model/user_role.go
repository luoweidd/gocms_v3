package model

import "time"

// UserRole 用户-角色关联模型（对应 user_roles 表，基础 RBAC 角色系统）
type UserRole struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:uk_user_role;not null" json:"user_id"`
	RoleID    uint      `gorm:"uniqueIndex:uk_user_role;not null" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
