package ngac

import "time"

// Role model
type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	Level       int       `gorm:"default:1" json:"level"`
	Path        string    `gorm:"size:255" json:"path"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int       `gorm:"default:1" json:"status"`
	Metadata    string    `gorm:"type:text" json:"metadata"`
}

func (Role) TableName() string { return "ngac_roles" }

// Permission model
type Permission struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name                  string    `gorm:"uniqueIndex:uk_name_code;size:191;not null" json:"name"`
	Code                  string    `gorm:"uniqueIndex:uk_name_code;uniqueIndex:uk_code;size:191;not null" json:"code"`
	Description           string    `gorm:"size:255" json:"description"`
	ResourceType          string    `gorm:"size:50;not null" json:"resource_type"`
	Action                string    `gorm:"size:50;not null" json:"action"`
	SubjectType           string    `gorm:"size:50;not null" json:"subject_type"`
	ResourceObjectClassID uint      `json:"resource_object_class_id"`
	Allowed               bool      `gorm:"default:true" json:"allowed"`
	Priority              int       `gorm:"default:0" json:"priority"`
	Efficiency            string    `gorm:"size:50" json:"efficiency"`
	RuleExpr              string    `gorm:"type:text" json:"rule_expr"`
	Metadata              string    `gorm:"type:text" json:"metadata"`
}

func (Permission) TableName() string { return "ngac_permissions" }

// RolePermission model - 修复字段名为 PermissionID
type RolePermission struct {
	ID           uint `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID       uint `gorm:"column:role_id;uniqueIndex:idx_role_perm;not null" json:"role_id"`
	PermissionID uint `gorm:"column:permission_id;uniqueIndex:idx_role_perm;not null" json:"permission_id"`
}

func (RolePermission) TableName() string { return "ngac_role_permissions" }

// UserRole model
type UserRole struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID uint `gorm:"uniqueIndex:idx_user_role;not null" json:"user_id"`
	RoleID uint `gorm:"uniqueIndex:idx_user_role;not null" json:"role_id"`
}

func (UserRole) TableName() string { return "ngac_user_roles" }

// ObjectClass model
type ObjectClass struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	Level       int       `gorm:"default:1" json:"level"`
	Path        string    `gorm:"size:255" json:"path"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Status      int       `gorm:"default:1" json:"status"`
	Attributes  string    `gorm:"type:text" json:"attributes"`
}

func (ObjectClass) TableName() string { return "ngac_object_classes" }

// ObjectAttribute model
type ObjectAttribute struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"size:255;uniqueIndex:idx_obj_attr;not null" json:"name"`
	Value         string `gorm:"size:500;not null" json:"value"`
	Type          string `gorm:"size:20" json:"type"`
	ResourceType  string `gorm:"size:50" json:"resource_type"`
	ResourceID    uint   `gorm:"index" json:"resource_id"`
	ObjectClassID uint   `json:"object_class_id"`
	ParentID      *uint  `gorm:"index" json:"parent_id"`
}

func (ObjectAttribute) TableName() string { return "ngac_object_attributes" }

// UserAttribute model
type UserAttribute struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint   `gorm:"index;not null" json:"user_id"`
	Name     string `gorm:"size:100;not null" json:"name"`
	Value    string `gorm:"size:500;not null" json:"value"`
	Type     string `gorm:"size:20" json:"type"`
	ParentID *uint  `json:"parent_id"`
}

func (UserAttribute) TableName() string { return "ngac_user_attributes" }

// ContextAttribute model
type ContextAttribute struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"uniqueIndex;size:100;not null" json:"name"`
	Value       string `gorm:"size:500;not null" json:"value"`
	Type        string `gorm:"size:20" json:"type"`
	IsGlobal    bool   `gorm:"default:false" json:"is_global"`
	Description string `gorm:"size:255" json:"description"`
}

func (ContextAttribute) TableName() string { return "ngac_context_attributes" }

// AccessLog model
type AccessLog struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint   `json:"user_id"`
	Username       string `gorm:"size:50" json:"username"`
	ResourceType   string `gorm:"size:50" json:"resource_type"`
	ResourceID     uint   `json:"resource_id"`
	Action         string `gorm:"size:50" json:"action"`
	SubjectAttrs   string `gorm:"type:text" json:"subject_attrs"`
	ObjectAttrs    string `gorm:"type:text" json:"object_attrs"`
	ContextAttrs   string `gorm:"type:text" json:"context_attrs"`
	Decision       string `gorm:"size:10" json:"decision"`
	EvaluationTime int64  `json:"evaluation_time"`
	Reason         string `gorm:"size:255" json:"reason"`
	RequestIP      string `gorm:"size:45" json:"request_ip"`
	UserAgent      string `gorm:"size:255" json:"user_agent"`
}

func (AccessLog) TableName() string { return "ngac_access_logs" }

// PDPEngineConfig model
type PDPEngineConfig struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string `gorm:"uniqueIndex;size:100;not null" json:"config_key"`
	ConfigValue string `gorm:"type:text;not null" json:"config_value"`
	Description string `gorm:"size:255" json:"description"`
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`
}

func (PDPEngineConfig) TableName() string { return "ngac_pdp_config" }

// ========== VO / Request DTOs ==========

type RoleTreeVO struct {
	ID       uint         `json:"id"`
	Name     string       `json:"name"`
	Code     string       `json:"code"`
	Level    int          `json:"level"`
	Sort     int          `json:"sort"`
	Status   int          `json:"status"`
	Children []RoleTreeVO `json:"children"`
}

type ObjectClassTreeVO struct {
	ID       uint                `json:"id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	Level    int                 `json:"level"`
	Sort     int                 `json:"sort"`
	Status   int                 `json:"status"`
	Children []ObjectClassTreeVO `json:"children"`
}

type RuleMatchInfo struct {
	PermissionID   uint   `json:"permission_id"`
	PermissionName string `json:"permission_name"`
	Code           string `json:"code"`
	Action         string `json:"action"`
	Allowed        bool   `json:"allowed"`
	Priority       int    `json:"priority"`
	Reason         string `json:"reason"`
}

type AuthzResponse struct {
	Allowed        bool            `json:"allowed"`
	Decision       string          `json:"decision"`
	Reason         string          `json:"reason"`
	RulesMatched   []RuleMatchInfo `json:"rules_matched"`
	EvaluationTime int64           `json:"evaluation_time"`
}

type CreateRoleReq struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description *string `json:"description"`
	ParentID    *uint   `json:"parent_id"`
	Sort        *int    `json:"sort"`
	Status      *int    `json:"status" binding:"omitempty,oneof=0 1"`
	Metadata    *string `json:"metadata"`
}

type UpdateRoleReq struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	ParentID    *uint   `json:"parent_id"`
	Sort        *int    `json:"sort"`
	Status      *int    `json:"status" binding:"omitempty,oneof=0 1"`
	Metadata    *string `json:"metadata"`
}

type RoleListReq struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Name     string `form:"name"`
	Code     string `form:"code"`
	Status   *int   `form:"status"`
}

type CreatePermissionReq struct {
	Name                  string  `json:"name" binding:"required"`
	Code                  string  `json:"code" binding:"required"`
	Description           *string `json:"description"`
	ResourceType          string  `json:"resource_type" binding:"required"`
	Action                string  `json:"action" binding:"required"`
	SubjectType           *string `json:"subject_type"`
	ResourceObjectClassID *uint   `json:"resource_object_class_id"`
	Allowed               *bool   `json:"allowed"`
	Priority              *int    `json:"priority"`
	Efficiency            *string `json:"efficiency"`
	RuleExpr              *string `json:"rule_expr"`
	Metadata              *string `json:"metadata"`
}

type UpdatePermissionReq struct {
	Name                  *string `json:"name"`
	Code                  *string `json:"code"`
	Description           *string `json:"description"`
	ResourceType          *string `json:"resource_type"`
	Action                *string `json:"action"`
	SubjectType           *string `json:"subject_type"`
	ResourceObjectClassID *uint   `json:"resource_object_class_id"`
	Allowed               *bool   `json:"allowed"`
	Priority              *int    `json:"priority"`
	Efficiency            *string `json:"efficiency"`
	RuleExpr              *string `json:"rule_expr"`
	Metadata              *string `json:"metadata"`
}

type PermissionListReq struct {
	Page         int    `form:"page" binding:"min=1"`
	PageSize     int    `form:"page_size" binding:"min=1,max=100"`
	ResourceType string `form:"resource_type"`
	Action       string `form:"action"`
}

type AssignRolesReq struct {
	UserID  uint   `json:"user_id" binding:"required"`
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

type AuthzRequest struct {
	UserID       uint              `json:"user_id"`
	Username     string            `json:"username"`
	ResourceType string            `json:"resource_type" binding:"required"`
	ResourceID   uint              `json:"resource_id"`
	Action       string            `json:"action" binding:"required"`
	SubjectAttrs map[string]string `json:"subject_attrs"`
	ObjectAttrs  map[string]string `json:"object_attrs"`
	ContextAttrs map[string]string `json:"context_attrs"`
}

type CreateUserAttributeReq struct {
	UserID   uint    `json:"user_id" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Value    string  `json:"value" binding:"required"`
	Type     *string `json:"type"`
	ParentID *uint   `json:"parent_id"`
}

type CreateObjectClassReq struct {
	Name        string  `json:"name" binding:"required"`
	Code        string  `json:"code" binding:"required"`
	Description *string `json:"description"`
	ParentID    *uint   `json:"parent_id"`
	Sort        *int    `json:"sort"`
	Status      *int    `json:"status" binding:"omitempty,oneof=0 1"`
	Attributes  *string `json:"attributes"`
}

type CreateContextAttributeReq struct {
	Name        string  `json:"name" binding:"required"`
	Value       string  `json:"value" binding:"required"`
	Type        *string `json:"type"`
	IsGlobal    *bool   `json:"is_global"`
	Description *string `json:"description"`
}

type CreateObjectAttributeReq struct {
	Name          string  `json:"name" binding:"required"`
	Value         string  `json:"value" binding:"required"`
	Type          *string `json:"type"`
	ResourceType  string  `json:"resource_type" binding:"required"`
	ResourceID    uint    `json:"resource_id"`
	ObjectClassID uint    `json:"object_class_id"`
	ParentID      *uint   `json:"parent_id"`
}

type BatchAuthzReq struct {
	UserID       uint                   `json:"user_id" binding:"required"`
	ResourceType string                 `json:"resource_type" binding:"required"`
	ResourceID   uint                   `json:"resource_id"`
	Action       string                 `json:"action" binding:"required"`
	SubjectAttrs map[string]interface{} `json:"subject_attrs"`
	ObjectAttrs  map[string]interface{} `json:"object_attrs"`
	ContextAttrs map[string]interface{} `json:"context_attrs"`
}
