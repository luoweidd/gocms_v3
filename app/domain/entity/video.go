package entity

import (
	"time"
)

// Video 视频领域实体
type Video struct {
	ID              uint       `json:"id"`
	TenantID        uint       `json:"tenant_id" gorm:"default:0"`
	Title           string     `json:"title" gorm:"size:255;not null"`
	Slug            string     `json:"slug" gorm:"uniqueIndex;size:255"`
	Desc            string     `json:"description" gorm:"type:text"`
	CoverImage      string     `json:"cover_image" gorm:"size:255"`
	VideoURL        string     `json:"video_url" gorm:"size:512;not null"`
	VideoSize       int64      `json:"video_size" gorm:"default:0"` // 文件大小(字节)
	Duration        int        `json:"duration" gorm:"default:0"`   // 时长(秒)
	Status          int        `json:"status" gorm:"default:0"`     // 0待处理 1已发布 2审核中 3拒绝
	PublishAt       *time.Time `json:"publish_at"`
	ViewCount       int        `json:"view_count" gorm:"default:0"`
	Likes           int64      `json:"likes" gorm:"default:0"`
	Dislikes        int64      `json:"dislikes" gorm:"default:0"`
	UserID          uint       `json:"user_id"`
	CategoryID      uint       `json:"category_id"`
	Tags            []string   `json:"tags" gorm:"type:text"`
	IsTop           bool       `json:"is_top" gorm:"default:false"`
	Version         int        `json:"version" gorm:"default:1"`
	PreviousVersion uint       `json:"previous_version_id" gorm:"default:0"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Video) TableName() string {
	return "videos"
}

// IsPublished 检查视频是否已发布
func (v *Video) IsPublished() bool {
	return v.Status == 1 && v.PublishAt != nil && v.PublishAt.Before(time.Now())
}

// CreateVideoParams 创建视频参数
type CreateVideoParams struct {
	Title      string
	Slug       string
	Desc       string
	CoverImage string
	VideoURL   string
	VideoSize  int64
	Duration   int
	Status     int
	PublishAt  *time.Time
	CategoryID uint
	Tags       []string
	UserID     uint
	TenantID   uint
	IsTop      bool
}

// UpdateVideoParams 更新视频参数
type UpdateVideoParams struct {
	Title      *string
	Slug       *string
	Desc       *string
	CoverImage *string
	VideoURL   *string
	VideoSize  *int64
	Duration   *int
	Status     *int
	PublishAt  *time.Time
	CategoryID *uint
	Tags       *[]string
	IsTop      *bool
}

// VideoCategory 视频分类领域实体
type VideoCategory struct {
	ID        uint      `json:"id"`
	TenantID  uint      `json:"tenant_id" gorm:"default:0"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	ParentID  uint      `json:"parent_id" gorm:"default:0"`
	SortOrder int       `json:"sort_order" gorm:"default:0"`
	Path      string    `json:"path" gorm:"size:255"`
	Status    int       `json:"status" gorm:"default:1"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (VideoCategory) TableName() string {
	return "video_categories"
}
