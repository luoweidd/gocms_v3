package model

// VideoCategory 视频分类模型（树形结构）
type VideoCategory struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Slug        string `gorm:"size:100;index" json:"slug"`                      // URL别名
	Pid         *uint  `gorm:"column:pid;index" json:"pid"`                     // 父分类ID（nil表示顶级）
	Sort        int    `gorm:"column:sort;default:0" json:"sort"`               // 排序（越小越靠前）
	Icon        string `gorm:"size:100" json:"icon"`                            // 分类图标
	Description string `gorm:"size:500" json:"description"`                     // 分类描述
	Status      int    `gorm:"default:1;index" json:"status"`                   // 0禁用 1启用
	VideoCount  int    `gorm:"column:video_count;default:0" json:"video_count"` // 视频数量（冗余统计）

	// 关联
	Children []VideoCategory `gorm:"foreignKey:Pid" json:"children,omitempty"`
	Parent   *VideoCategory  `gorm:"foreignKey:Pid" json:"parent,omitempty"`
}

func (VideoCategory) TableName() string {
	return "video_categories"
}
