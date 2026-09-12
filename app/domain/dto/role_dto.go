package dto

import (
	"time"
)

// ==================== Role DTOs ====================

// RoleDTO 角色数据传输对象
type RoleDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	NameEn      string    `json:"name_en"`
	Description string    `json:"description"`
	DataScope   string    `json:"data_scope"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoleListQuery 角色列表查询DTO
type RoleListQuery struct {
	Keyword string `form:"keyword"`
	Status  *int   `form:"status"`
	Page    int    `form:"page" binding:"required,min=1"`
	Size    int    `form:"page_size" binding:"required,min=1,max=100"`
}

// CreateRoleDTO 创建角色DTO
type CreateRoleDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Code        string `json:"code" binding:"required,min=2,max=50"`
	NameEn      string `json:"name_en"`
	Description string `json:"description"`
	DataScope   string `json:"data_scope"`
	Status      int    `json:"status"`
}

// UpdateRoleDTO 更新角色DTO
type UpdateRoleDTO struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	NameEn      *string `json:"name_en"`
	Description *string `json:"description"`
	DataScope   *string `json:"data_scope"`
	Status      *int    `json:"status"`
}

// ==================== Tag DTOs ====================

// TagDTO 标签数据传输对象
type TagDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TagListQuery 标签列表查询DTO
type TagListQuery struct {
	Keyword string `form:"keyword"`
	Status  *int   `form:"status"`
	Page    int    `form:"page" binding:"required,min=1"`
	Size    int    `form:"page_size" binding:"required,min=1,max=100"`
}

// CreateTagDTO 创建标签DTO
type CreateTagDTO struct {
	Name   string `json:"name" binding:"required,min=1,max=50"`
	Status int    `json:"status"`
}

// UpdateTagDTO 更新标签DTO
type UpdateTagDTO struct {
	Name   *string `json:"name"`
	Status *int    `json:"status"`
}
