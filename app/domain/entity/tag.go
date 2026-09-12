package entity

import (
	"time"
)

// Tag 标签领域实体
type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"uniqueIndex;size:50;not null"`
	Status    int       `json:"status" gorm:"default:1"` // 1启用 0禁用
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Tag) TableName() string {
	return "tags"
}

// IsEnabled 检查标签是否启用
func (t *Tag) IsEnabled() bool {
	return t.Status == 1
}

// CreateTagParams 创建标签参数
type CreateTagParams struct {
	Name   string
	Status int
}

// UpdateTagParams 更新标签参数
type UpdateTagParams struct {
	Name   *string
	Status *int
}

// TagListQuery 标签列表查询参数
type TagListQuery struct {
	Keyword string
	Status  *int
	Page    int
	PageSQL int
}
