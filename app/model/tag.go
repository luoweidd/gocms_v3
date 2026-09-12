package model

import "time"

// Tag 标签模型
type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex:uk_tag_name;size:50;not null" json:"name"`
	Description string `gorm:"size:255;default:'';column:description" json:"description"`
	Status      int    `gorm:"default:1" json:"status"` // 1启用 0禁用
}

func (Tag) TableName() string {
	return "tags"
}

// VideoTag 视频-标签关联表模型 (many2many join table)
type VideoTag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	VideoID   uint      `gorm:"uniqueIndex:uk_video_tags_video_tag;not null;column:video_id" json:"video_id"`
	TagID     uint      `gorm:"uniqueIndex:uk_video_tags_video_tag;index:idx_video_tags_tag_id;not null;column:tag_id" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (VideoTag) TableName() string {
	return "video_tags"
}

// ArticleTag 文章-标签关联表模型 (many2many join table)
type ArticleTag struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint      `gorm:"uniqueIndex:uk_article_tags_article_tag;not null;column:article_id" json:"article_id"`
	TagID     uint      `gorm:"uniqueIndex:uk_article_tags_article_tag;index:idx_article_tags_tag_id;not null;column:tag_id" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (ArticleTag) TableName() string {
	return "article_tags"
}
