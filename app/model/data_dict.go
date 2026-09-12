package model

import "time"

// DataDict 数据字典模型
type DataDict struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	DictType  string    `gorm:"index;size:100;not null" json:"dict_type"`
	DictLabel string    `gorm:"size:100;not null" json:"dict_label"`
	DictValue string    `gorm:"size:255" json:"dict_value"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Status    int       `gorm:"default:1" json:"status"` // 1启用 0禁用
	Remark    string    `gorm:"size:500" json:"remark"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (DataDict) TableName() string {
	return "data_dicts"
}
