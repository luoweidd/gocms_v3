package model

import (
	"time"
)

// Comment 评论模型（文章/视频通用）
type Comment struct {
	BaseModel
	ArticleID   uint       `gorm:"index" json:"article_id"`           // 关联文章ID
	VideoID     uint       `gorm:"index" json:"video_id"`             // 关联视频ID
	UserID      uint       `gorm:"index" json:"user_id"`              // 评论用户ID
	Nickname    string     `gorm:"size:100" json:"nickname"`          // 评论者昵称(支持匿名)
	Email       string     `gorm:"size:100" json:"email"`             // 评论者邮箱
	Content     string     `gorm:"type:text;not null" json:"content"` // 评论内容
	Status      int        `gorm:"default:0;index" json:"status"`     // 0待审核 1已通过 2已拒绝 3已删除
	ParentID    uint       `gorm:"default:0;index" json:"parent_id"`  // 父评论ID(回复)
	ReplyCount  int        `gorm:"default:0" json:"reply_count"`      // 回复数量
	Likes       int64      `gorm:"default:0" json:"likes"`            // 点赞数
	AuditBy     uint       `gorm:"default:0" json:"audit_by"`         // 审核人ID
	AuditRemark string     `gorm:"size:500" json:"audit_remark"`      // 审核备注
	IPAddress   string     `gorm:"size:45" json:"ip_address"`         // IP地址
	UserAgent   string     `gorm:"size:500" json:"user_agent"`        // 用户代理
	AuditedAt   *time.Time `json:"audited_at"`                        // 审核时间
}

func (Comment) TableName() string {
	return "comments"
}

// CommentAuditRequest 评论审核请求
type CommentAuditRequest struct {
	Status      int    `json:"status" binding:"required,oneof=1 2"`
	AuditRemark string `json:"audit_remark"`
}

// CommentCreateRequest 评论创建请求
type CommentCreateRequest struct {
	ArticleID *uint  `json:"article_id"`
	VideoID   *uint  `json:"video_id"`
	Nickname  string `json:"nickname" binding:"max=100"`
	Email     string `json:"email" binding:"email,max=100"`
	Content   string `json:"content" binding:"required,min=1,max=2000"`
	ParentID  uint   `json:"parent_id"`
}

// CommentListQuery 评论列表查询参数
type CommentListQuery struct {
	ArticleID uint   `form:"article_id"`
	VideoID   uint   `form:"video_id"`
	Status    *int   `form:"status"`
	UserID    uint   `form:"user_id"`
	Page      int    `form:"page" binding:"min=1"`
	PageSize  int    `form:"page_size" binding:"min=1,max=100"`
	SortBy    string `form:"sort_by"`
}

// CommentListItem 评论列表项(带用户信息)
type CommentListItem struct {
	ID         uint              `json:"id"`
	ArticleID  uint              `json:"article_id"`
	VideoID    uint              `json:"video_id"`
	UserID     uint              `json:"user_id"`
	Nickname   string            `json:"nickname"`
	Email      string            `json:"email"`
	Content    string            `json:"content"`
	Status     int               `json:"status"`
	ParentID   uint              `json:"parent_id"`
	ReplyCount int               `json:"reply_count"`
	Likes      int64             `json:"likes"`
	AuditBy    uint              `json:"audit_by"`
	AuditedAt  *time.Time        `json:"audited_at"`
	IPAddress  string            `json:"ip_address"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	UserInfo   *UserInfo         `json:"user_info,omitempty" gorm:"-"`
	Replies    []CommentListItem `json:"replies,omitempty" gorm:"-"`
}

// UserInfo 用户信息(用于评论列表)
type UserInfo struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
