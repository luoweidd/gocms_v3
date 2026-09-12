package model

import "time"

// MessageCategory 消息类别表
type MessageCategory struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(100);not null;uniqueIndex" json:"name"`
	Code        string `gorm:"column:code;type:varchar(50);not null;uniqueIndex" json:"code"`
	Description string `gorm:"column:description;type:varchar(255)" json:"description"`
	IconType    string `gorm:"column:icon_type;type:varchar(50)" json:"icon_type"`
	Sort        int    `gorm:"column:sort;type:int;default:0" json:"sort"`
	Status      int8   `gorm:"column:status;type:tinyint;default:1" json:"status"`
}

// TableName 指定表名
func (MessageCategory) TableName() string {
	return "message_categories"
}

// MessageCategoryListQuery 消息类别列表查询参数
type MessageCategoryListQuery struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"page_size" binding:"required"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   *int8  `form:"status"`
}

// MessageCategoryResponse 消息类别响应
type MessageCategoryResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Code         string    `json:"code"`
	Description  string    `json:"description"`
	IconType     string    `json:"icon_type"`
	Sort         int       `json:"sort"`
	Status       int8      `json:"status"`
	MessageCount int64     `json:"message_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateMessageCategoryRequest 创建消息类别请求
type CreateMessageCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	IconType    string `json:"icon_type"`
	Sort        int    `json:"sort"`
	Status      int8   `json:"status"`
}

// UpdateMessageCategoryRequest 更新消息类别请求
type UpdateMessageCategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IconType    string `json:"icon_type"`
	Sort        int    `json:"sort"`
	Status      int8   `json:"status"`
}
