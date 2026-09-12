package model

import (
	"time"
)

// Tenant 租户模型
type Tenant struct {
	BaseModel
	Name         string     `gorm:"size:100;not null" json:"name"`                        // 租户名称
	Code         string     `gorm:"size:50;uniqueIndex:uk_tenant_code" json:"code"`       // 租户编码
	ContactName  string     `gorm:"size:50" json:"contact_name"`                          // 联系人
	ContactPhone string     `gorm:"size:20" json:"contact_phone"`                         // 联系电话
	ContactEmail string     `gorm:"size:100" json:"contact_email"`                        // 联系邮箱
	Domain       string     `gorm:"size:255;uniqueIndex:uk_tenant_domain" json:"domain"`  // 绑定域名
	APIKey       string     `gorm:"size:64;uniqueIndex:uk_tenant_api_key" json:"api_key"` // API密钥
	PlanType     string     `gorm:"size:20;default:free" json:"plan_type"`                // 套餐类型: free/standard/premium/enterprise
	MaxUsers     int        `gorm:"default:10" json:"max_users"`                          // 最大用户数
	MaxStorage   int64      `gorm:"default:10737418240" json:"max_storage"`               // 最大存储(字节, 10GB)
	UsedStorage  int64      `gorm:"default:0" json:"used_storage"`                        // 已用存储(字节)
	ExpireAt     *time.Time `json:"expire_at"`                                            // 到期时间
	Status       int        `gorm:"default:1" json:"status"`                              // 1启用 0停用 2已过期
	Settings     string     `gorm:"type:text" json:"settings"`                            // 租户配置(JSON)
}

func (Tenant) TableName() string {
	return "tenants"
}

// TenantCreateRequest 创建租户请求
type TenantCreateRequest struct {
	Name         string `json:"name" binding:"required,max:100"`
	Code         string `json:"code" binding:"required,max:50"`
	ContactName  string `json:"contact_name" binding:"max:50"`
	ContactPhone string `json:"contact_phone" binding:"max:20"`
	ContactEmail string `json:"contact_email" binding:"email,max:100"`
	Domain       string `json:"domain" binding:"max:255"`
	PlanType     string `json:"plan_type" binding:"required,oneof=free standard premium enterprise"`
	MaxUsers     int    `json:"max_users"`
	MaxStorage   int64  `json:"max_storage"`
	ExpireAt     string `json:"expire_at"`
}

// TenantUpdateRequest 更新租户请求
type TenantUpdateRequest struct {
	Name         string `json:"name" binding:"max:100"`
	ContactName  string `json:"contact_name" binding:"max:50"`
	ContactPhone string `json:"contact_phone" binding:"max:20"`
	ContactEmail string `json:"contact_email" binding:"email,max:100"`
	Domain       string `json:"domain" binding:"max:255"`
	APIKey       string `json:"api_key" binding:"max:64"`
	PlanType     string `json:"plan_type" binding:"oneof=free standard premium enterprise"`
	MaxUsers     *int   `json:"max_users"`
	MaxStorage   *int64 `json:"max_storage"`
	ExpireAt     string `json:"expire_at"`
	Status       *int   `json:"status"`
	Settings     string `json:"settings"`
}

// TenantListQuery 租户列表查询参数
type TenantListQuery struct {
	Name     string `form:"name"`
	Code     string `form:"code"`
	PlanType string `form:"plan_type"`
	Status   *int   `form:"status"`
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// TenantListItem 租户列表项
type TenantListItem struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Code         string     `json:"code"`
	ContactName  string     `json:"contact_name"`
	ContactPhone string     `json:"contact_phone"`
	ContactEmail string     `json:"contact_email"`
	Domain       string     `json:"domain"`
	APIKey       string     `json:"api_key"`
	PlanType     string     `json:"plan_type"`
	MaxUsers     int        `json:"max_users"`
	MaxStorage   int64      `json:"max_storage"`
	UsedStorage  int64      `json:"used_storage"`
	ExpireAt     *time.Time `json:"expire_at"`
	Status       int        `json:"status"`
	Settings     string     `json:"settings,omitempty"`
	UserCount    int        `json:"user_count"` // 当前用户数
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// TenantStats 租户统计
type TenantStats struct {
	Total    int64            `json:"total"`
	Active   int64            `json:"active"`
	Inactive int64            `json:"inactive"`
	Expired  int64            `json:"expired"`
	ByPlan   map[string]int64 `json:"by_plan"` // 按套餐类型统计
}

// TenantUserItem 租户用户项
type TenantUserItem struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TenantUserRequest 租户用户操作请求
type TenantUserRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}
