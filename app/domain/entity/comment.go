package entity

import (
	"time"
)

// Comment 评论领域实体
type Comment struct {
	ID         uint      `json:"id"`
	TenantID   uint      `json:"tenant_id" gorm:"default:0"`
	UserID     uint      `json:"user_id"`
	TargetType string    `json:"target_type" gorm:"size:20"` // article, video
	TargetID   uint      `json:"target_id"`
	ParentID   uint      `json:"parent_id" gorm:"default:0"` // 父评论ID
	Content    string    `json:"content" gorm:"type:text;not null"`
	Status     int       `json:"status" gorm:"default:1"` // 0待审 1通过 2拒绝
	LikeCount  int64     `json:"like_count" gorm:"default:0"`
	IPAddress  string    `json:"ip_address" gorm:"size:45"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}

// IsApproved 检查评论是否已审核
func (c *Comment) IsApproved() bool {
	return c.Status == 1
}

// CreateCommentParams 创建评论参数
type CreateCommentParams struct {
	UserID     uint
	TargetType string
	TargetID   uint
	ParentID   uint
	Content    string
	IPAddress  string
	TenantID   uint
}

// OAuth2AuthCode OAuth2授权码实体
type OAuth2AuthCode struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Provider    string    `json:"provider" gorm:"size:20"` // wechat, qq, github
	OpenID      string    `json:"openid" gorm:"size:64"`
	Code        string    `json:"code" gorm:"uniqueIndex;size:64"`
	RedirectURI string    `json:"redirect_uri" gorm:"size:255"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (OAuth2AuthCode) TableName() string {
	return "oauth_auth_codes"
}

// Tenant 租户实体
type Tenant struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Domain    string    `json:"domain" gorm:"uniqueIndex;size:100"`
	APIKey    string    `json:"api_key" gorm:"uniqueIndex;size:64"`
	Status    int       `json:"status" gorm:"default:1"`  // 1启用 0禁用
	PlanType  string    `json:"plan_type" gorm:"size:20"` // free, standard, enterprise
	MaxUsers  int       `json:"max_users" gorm:"default:10"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// ContentVersion 内容版本实体
type ContentVersion struct {
	ID         uint      `json:"id"`
	TargetType string    `json:"target_type" gorm:"size:20"` // article, video
	TargetID   uint      `json:"target_id"`
	Version    int       `json:"version"`
	Data       string    `json:"data" gorm:"type:longtext"` // JSON格式的内容快照
	Comment    string    `json:"comment" gorm:"size:255"`   // 版本变更说明
	UserID     uint      `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (ContentVersion) TableName() string {
	return "content_versions"
}

// CronJob 定时任务实体
type CronJob struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name" gorm:"size:100;not null"`
	Expression string     `json:"expression" gorm:"size:50"` // cron表达式
	Handler    string     `json:"handler" gorm:"size:100"`   // 处理器标识
	Status     int        `json:"status" gorm:"default:1"`   // 1启用 0禁用
	LastRunAt  *time.Time `json:"last_run_at"`
	NextRunAt  *time.Time `json:"next_run_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (CronJob) TableName() string {
	return "cron_jobs"
}

// Webhook 回调实体
type Webhook struct {
	ID            uint       `json:"id"`
	URL           string     `json:"url" gorm:"size:512;not null"`
	Secret        string     `json:"secret" gorm:"size:128"`
	Events        []string   `json:"events" gorm:"type:text"` // article.created, video.published等
	Status        int        `json:"status" gorm:"default:1"` // 1启用 0禁用
	LastTriggered *time.Time `json:"last_triggered_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Webhook) TableName() string {
	return "webhooks"
}

// Language 国际化语言实体
type Language struct {
	Code      string    `json:"code" gorm:"primarykey;size:10"` // zh-CN, en-US
	Name      string    `json:"name" gorm:"size:50"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
}

func (Language) TableName() string {
	return "languages"
}

// Translation 翻译实体
type Translation struct {
	ID        uint      `json:"id"`
	Key       string    `json:"key" gorm:"uniqueIndex:idx_language_key;size:255"`
	Language  string    `json:"language" gorm:"uniqueIndex:idx_language_key;size:10"`
	Value     string    `json:"value" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Translation) TableName() string {
	return "translations"
}

// FileUploadRecord 文件上传记录实体
type FileUploadRecord struct {
	ID         uint      `json:"id"`
	TenantID   uint      `json:"tenant_id" gorm:"default:0"`
	FileName   string    `json:"file_name" gorm:"size:255"`
	FilePath   string    `json:"file_path" gorm:"size:512"`
	FileSize   int64     `json:"file_size"`
	ChunkCount int       `json:"chunk_count"` // 分片数
	MimeType   string    `json:"mime_type" gorm:"size:100"`
	UserID     uint      `json:"user_id"`
	Status     int       `json:"status" gorm:"default:0"` // 0上传中 1完成 2失败
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (FileUploadRecord) TableName() string {
	return "file_upload_records"
}
