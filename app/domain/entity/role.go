package entity

import (
	"time"
)

// Role 角色领域实体
type Role struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:100"`
	NameEn      string    `json:"name_en" gorm:"size:100"`
	Description string    `json:"description" gorm:"size:500"`
	DataScope   string    `json:"data_scope" gorm:"size:50"` // all/own
	Status      int       `json:"status" gorm:"default:1"`   // 1启用 0禁用
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

// IsEnabled 检查角色是否启用
func (r *Role) IsEnabled() bool {
	return r.Status == 1
}

// CreateRoleParams 创建角色参数
type CreateRoleParams struct {
	Name        string
	Code        string
	NameEn      string
	Description string
	DataScope   string
	Status      int
}

// UpdateRoleParams 更新角色参数
type UpdateRoleParams struct {
	Name        *string
	Code        *string
	NameEn      *string
	Description *string
	DataScope   *string
	Status      *int
}

// RoleListQuery 角色列表查询参数
type RoleListQuery struct {
	Keyword string
	Status  *int
	Page    int
	PageSQL int
}
