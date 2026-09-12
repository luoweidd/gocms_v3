package model

import "time"

// SystemMessage 系统消息表
type SystemMessage struct {
	BaseModel
	Title         string     `gorm:"column:title;type:varchar(255);not null;index" json:"title"`
	Content       string     `gorm:"column:content;type:text" json:"content"`
	MessageType   string     `gorm:"column:message_type;type:varchar(20);index" json:"message_type"` // system/publish/comment_reply/audit_passed/audited_failed
	SenderID      uint       `gorm:"column:sender_id;type:bigint unsigned;not null;default:0" json:"sender_id"`
	SenderName    string     `gorm:"column:sender_name;type:varchar(100)" json:"sender_name"`
	Priority      int8       `gorm:"column:priority;type:tinyint;not null;default:0" json:"priority"` // 0普通 1重要 2紧急
	TargetType    string     `gorm:"column:target_type;type:varchar(20)" json:"target_type"`          // all/user/role
	TargetUserIDs string     `gorm:"column:target_user_ids;type:text" json:"target_user_ids"`         // JSON数组
	IsReadDefault int8       `gorm:"column:is_read_default;type:tinyint;not null;default:0" json:"is_read_default"`
	ExpireAt      *time.Time `gorm:"column:expire_at;type:datetime" json:"expire_at"`
	Status        int8       `gorm:"column:status;type:tinyint;not null;default:1;index" json:"status"` // 1已发布 0草稿 2已撤回
}

// TableName 指定表名
func (SystemMessage) TableName() string {
	return "system_messages"
}

// UserNotification 用户通知表
type UserNotification struct {
	BaseModel
	UserID    uint       `gorm:"column:user_id;type:bigint unsigned;not null;index;uniqueIndex:uk_user_message" json:"user_id"`
	MessageID uint       `gorm:"column:message_id;type:bigint unsigned;not null;index;uniqueIndex:uk_user_message" json:"message_id"`
	IsRead    int8       `gorm:"column:is_read;type:tinyint;not null;default:0" json:"is_read"`
	ReadAt    *time.Time `gorm:"column:read_at;type:datetime" json:"read_at"`
}

// TableName 指定表名
func (UserNotification) TableName() string {
	return "user_notifications"
}

// SystemMessageQuery 系统消息查询参数
type SystemMessageQuery struct {
	Keyword     string `form:"keyword"`
	MessageType string `form:"message_type"`
	Status      *int8  `form:"status"`
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
}

// SystemMessageResponse 系统消息响应
type SystemMessageResponse struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Content       string     `json:"content"`
	MessageType   string     `json:"message_type"`
	SenderID      uint       `json:"sender_id"`
	SenderName    string     `json:"sender_name"`
	Priority      int8       `json:"priority"`
	TargetType    string     `json:"target_type"`
	TargetUserIDs string     `json:"target_user_ids"`
	IsReadDefault int8       `json:"is_read_default"`
	ExpireAt      *time.Time `json:"expire_at"`
	Status        int8       `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
}

// UserNotificationQuery 用户通知查询参数
type UserNotificationQuery struct {
	UserID      uint   `form:"user_id"`
	IsRead      *int8  `form:"is_read"`
	MessageType string `form:"message_type"`
	Page        int    `form:"page" binding:"required"`
	PageSize    int    `form:"page_size" binding:"required"`
}

// UserNotificationResponse 用户通知响应
type UserNotificationResponse struct {
	ID        uint        `json:"id"`
	UserID    uint        `json:"user_id"`
	MessageID uint        `json:"message_id"`
	Message   MessageInfo `json:"message"`
	IsRead    int8        `json:"is_read"`
	ReadAt    *time.Time  `json:"read_at"`
	CreatedAt time.Time   `json:"created_at"`
}

// MessageInfo 消息概要信息
type MessageInfo struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	Priority    int8      `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
}

// NotificationStats 通知统计
type NotificationStats struct {
	TotalCount   int64 `json:"total_count"`
	UnreadCount  int64 `json:"unread_count"`
	ReadCount    int64 `json:"read_count"`
	SystemCount  int64 `json:"system_count"`
	PublishCount int64 `json:"publish_count"`
	CommentCount int64 `json:"comment_count"`
	AuditCount   int64 `json:"audit_count"`
}

// CreateUserNotificationRequest 创建用户通知请求
type CreateUserNotificationRequest struct {
	UserID    uint `form:"user_id" binding:"required"`
	MessageID uint `form:"message_id" binding:"required"`
}

// MarkAsReadRequest 标记已读请求
type MarkAsReadRequest struct {
	IsRead int8 `json:"is_read" form:"is_read" binding:"required"`
}

// PublishMessageRequest 发布消息请求
type PublishMessageRequest struct {
	Title         string `form:"title" binding:"required"`
	Content       string `form:"content" binding:"required"`
	MessageType   string `form:"message_type"`
	Priority      int8   `form:"priority"`
	TargetType    string `form:"target_type"`
	TargetUserIDs []uint `form:"target_user_ids"`
	IsReadDefault int8   `form:"is_read_default"`
	ExpireAt      string `form:"expire_at"`
}

// BatchMarkAsReadRequest 批量标记已读请求
type BatchMarkAsReadRequest struct {
	IDs []uint `form:"ids" binding:"required"`
}
