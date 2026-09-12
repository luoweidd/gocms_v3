package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型（所有表都继承此结构）
type BaseModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (BaseModel) TableName() string {
	return ""
}