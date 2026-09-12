package model

// ArticleCategory 文章分类模型（对应 article_categories 表）
type ArticleCategory struct {
	BaseModel
	Name        string             `gorm:"uniqueIndex:uk_article_category_name;size:50;not null" json:"name"` // 分类名称
	Icon        string             `gorm:"size:100" json:"icon"`                                              // 图标类名
	SortOrder   int                `gorm:"column:sort_order;default:0" json:"sort"`                           // 排序
	Status      int                `gorm:"default:1;index:idx_article_category_status" json:"status"`         // 1启用 0禁用
	Description string             `gorm:"size:500" json:"description"`                                       // 描述
	Children    []*ArticleCategory `gorm:"-" json:"children,omitempty"`                                       // 子分类
}

// TableName 指定表名
func (ArticleCategory) TableName() string {
	return "article_categories"
}
