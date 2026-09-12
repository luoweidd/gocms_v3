package entity

import (
	"time"
)

// Article 文章领域实体
type Article struct {
	ID              uint       `json:"id"`
	TenantID        uint       `json:"tenant_id" gorm:"default:0"`
	Title           string     `json:"title" gorm:"size:255;not null"`
	Slug            string     `json:"slug" gorm:"uniqueIndex;size:255"`
	Content         string     `json:"content" gorm:"type:longtext"`
	Summary         string     `json:"summary" gorm:"type:text"`
	CoverImage      string     `json:"cover_image" gorm:"size:255"`
	Status          int        `json:"status" gorm:"default:0"` // 0草稿 1已发布 2已撤回
	PublishAt       *time.Time `json:"publish_at"`
	ViewCount       int        `json:"view_count" gorm:"default:0"`
	Likes           int64      `json:"likes" gorm:"default:0"`
	UserID          uint       `json:"user_id"`
	CategoryID      uint       `json:"category_id"`
	Tags            []string   `json:"tags" gorm:"type:text"`
	IsTop           bool       `json:"is_top" gorm:"default:false"`
	AllowComment    bool       `json:"allow_comment" gorm:"default:true"`
	Version         int        `json:"version" gorm:"default:1"`             // 内容版本号
	PreviousVersion uint       `json:"previous_version_id" gorm:"default:0"` // 上一个版本ID

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Article) TableName() string {
	return "articles"
}

// IsPublished 检查文章是否已发布
func (a *Article) IsPublished() bool {
	return a.Status == 1 && a.PublishAt != nil && a.PublishAt.Before(time.Now())
}

// CreateArticleParams 创建文章参数
type CreateArticleParams struct {
	Title        string
	Slug         string
	Content      string
	Summary      string
	CoverImage   string
	Status       int
	PublishAt    *time.Time
	CategoryID   uint
	Tags         []string
	UserID       uint
	TenantID     uint
	IsTop        bool
	AllowComment bool
}

// UpdateArticleParams 更新文章参数
type UpdateArticleParams struct {
	Title        *string
	Slug         *string
	Content      *string
	Summary      *string
	CoverImage   *string
	Status       *int
	PublishAt    *time.Time
	CategoryID   *uint
	Tags         *[]string
	IsTop        *bool
	AllowComment *bool
}

// ArticleCategory 文章分类领域实体
type ArticleCategory struct {
	ID        uint      `json:"id"`
	TenantID  uint      `json:"tenant_id" gorm:"default:0"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	ParentID  uint      `json:"parent_id" gorm:"default:0"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	Path      string    `json:"path" gorm:"size:255"` // 路径，用于快速查询父子关系
	Status    int       `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}

// GetParentID 获取父分类ID
func (c *ArticleCategory) GetParentID() uint {
	return c.ParentID
}

// GetPath 获取路径
func (c *ArticleCategory) GetPath() string {
	return c.Path
}

// ArticleVersion 文章版本历史实体（内容版本管理）
type ArticleVersion struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id" gorm:"uniqueIndex:idx_article_version"`
	Version   int       `json:"version"`
	Content   string    `json:"content" gorm:"type:longtext"`
	Summary   string    `json:"summary" gorm:"type:text"`
	Title     string    `json:"title" gorm:"size:255"`
	CreatedAt time.Time `json:"created_at"`
}

func (ArticleVersion) TableName() string {
	return "article_versions"
}
