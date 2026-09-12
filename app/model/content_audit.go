package model

import "time"

// ContentAudit 内容审核
type ContentAudit struct {
	BaseModel
	AuditNo              string     `gorm:"size:50;uniqueIndex" json:"audit_no"`
	ContentType          string     `gorm:"size:20;not null;index:idx_content_type" json:"content_type"`
	ContentID            uint       `gorm:"not null;index:idx_content_id" json:"content_id"`
	ContentTitle         string     `gorm:"size:255" json:"content_title"`
	SubmitterID          uint       `gorm:"not null" json:"submitter_id"`
	SubmitterName        string     `gorm:"size:100" json:"submitter_name"`
	AuditorID            uint       `json:"auditor_id"`
	AuditorName          string     `gorm:"size:100" json:"auditor_name"`
	Status               int8       `gorm:"not null;default:0" json:"status"` // 0待审 1通过 2驳回 3需修改
	AuditComment         string     `gorm:"type:text" json:"audit_comment"`
	PreviousStatus       int8       `json:"previous_status"`
	NextStatus           int8       `json:"next_status"`
	RequiredModification string     `gorm:"type:text" json:"required_modification"`
	AuditedAt            *time.Time `json:"audited_at"`

	// 关联对象（不存储到数据库，用于预加载）
	Article *Article `gorm:"-" json:"article,omitempty"`
	Video   *Video   `gorm:"-" json:"video,omitempty"`
	Comment *Comment `gorm:"-" json:"comment,omitempty"`
}

func (ContentAudit) TableName() string {
	return "content_audits"
}

// ContentAuditListRequest 审核列表查询请求
type ContentAuditListRequest struct {
	Page        int    `form:"page" binding:"min=1"`
	PageSize    int    `form:"page_size" binding:"max=200"`
	ContentType string `form:"content_type"`
	Status      *int8  `form:"status"`
	Keyword     string `form:"keyword"`
	AuditorID   uint   `form:"auditor_id"`
}

// AuditStats 审核统计
type AuditStats struct {
	Total       int64 `json:"total"`
	Pending     int64 `json:"pending"`
	Passed      int64 `json:"passed"`
	Rejected    int64 `json:"rejected"`
	NeedsModify int64 `json:"needs_modify"`
}

// BatchAuditRequest 批量审核请求
type BatchAuditRequest struct {
	IDs                  []uint `json:"ids" binding:"required"`
	Action               string `json:"action" binding:"required"` // approve/reject
	AuditComment         string `json:"audit_comment"`
	RequiredModification string `json:"required_modification"`
}
