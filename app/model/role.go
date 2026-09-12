package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型（唯一索引由迁移文件控制，GORM不做处理）
type Role struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Code        string         `gorm:"size:100" json:"code"`
	NameEn      string         `gorm:"size:100" json:"name_en"`
	Description string         `gorm:"size:500" json:"description"`
	DataScope   string         `gorm:"size:50" json:"data_scope"` // all/own
	Status      int            `gorm:"default:1" json:"status"`   // 1启用 0禁用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
