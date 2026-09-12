package model

import "time"

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;size:100;not null" json:"config_key"`
	ConfigValue string    `gorm:"type:text" json:"config_value"`
	ConfigType  string    `gorm:"size:20" json:"config_type"` // string/number/boolean/json
	Description string    `gorm:"size:500" json:"description"`
	IsSystem    bool      `gorm:"default:false" json:"is_system"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
