package model

import (
	"time"
)

// Article 文章模型
type Article struct {
	BaseModel
	Title        string           `gorm:"size:255;not null;column:title" json:"title"`                                                                      // 文章标题
	Summary      string           `gorm:"size:500;column:summary" json:"summary"`                                                                           // 文章摘要
	Content      string           `gorm:"type:text;column:content" json:"content"`                                                                          // HTML内容
	CoverImage   string           `gorm:"size:500;column:cover_image" json:"cover_image"`                                                                   // 封面图URL
	CategoryID   uint             `gorm:"index:idx_article_category;column:category_id" json:"category_id"`                                                 // 分类ID
	AuthorID     uint             `gorm:"index:idx_article_author;column:author_id" json:"author_id"`                                                       // 创建用户ID
	Status       int              `gorm:"default:0;index:idx_article_status;column:status" json:"status"`                                                   // 0草稿 1发布 2下架
	ViewCount    int64            `gorm:"default:0;column:view_count" json:"view_count"`                                                                    // 浏览量
	LikeCount    int64            `gorm:"default:0;column:like_count" json:"like_count"`                                                                    // 点赞数
	CommentCount int64            `gorm:"default:0;column:comment_count" json:"comment_count"`                                                              // 评论数
	PublishedAt  *time.Time       `gorm:"index:idx_article_publish;column:publish_at" json:"published_at"`                                                  // 发布时间
	Tags         []Tag            `gorm:"many2many:article_tags;" json:"tags"`                                                                              // 标签
	Category     *ArticleCategory `gorm:"foreignKey:CategoryID" json:"category"`                                                                            // 所属分类
	Audits       []ContentAudit   `gorm:"foreignKey:ContentID;references:id;constraint:where content_type = 'article';preloadable" json:"audits,omitempty"` // 审核记录
}

func (Article) TableName() string {
	return "articles"
}
