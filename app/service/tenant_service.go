package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"gorm.io/gorm"
)

type TenantService struct{}

func NewTenantService() *TenantService {
	return &TenantService{}
}

// generateAPIKey 生成 API Key
func generateAPIKey() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return "tk_" + hex.EncodeToString(b)
}

// CreateTenant 创建租户
func (s *TenantService) CreateTenant(req model.TenantCreateRequest) (*model.Tenant, error) {
	// 验证数据库连接
	if db.GetDB() == nil {
		return nil, errors.New("数据库未初始化")
	}

	// 生成 API Key
	apiKey := generateAPIKey()

	tenant := &model.Tenant{
		Name:         req.Name,
		Code:         req.Code,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		Domain:       req.Domain,
		APIKey:       apiKey,
		PlanType:     req.PlanType,
		MaxUsers: func() int {
			if req.MaxUsers > 0 {
				return req.MaxUsers
			} else {
				m := map[string]int{"free": 10, "standard": 50, "premium": 200, "enterprise": 99999}
				return m[req.PlanType]
			}
		}(),
		MaxStorage: func() int64 {
			if req.MaxStorage > 0 {
				return req.MaxStorage
			} else {
				m := map[string]int64{"free": 10737418240, "standard": 53687091200, "premium": 107374182400, "enterprise": 1073741824000}
				return m[req.PlanType]
			}
		}(),
		Status: 1,
	}

	if req.ExpireAt != "" {
		expire, err := time.Parse("2006-01-02", req.ExpireAt)
		if err == nil {
			tenant.ExpireAt = &expire
		}
	}

	if err := db.GetDB().Create(tenant).Error; err != nil {
		return nil, fmt.Errorf("创建租户失败: %w", err)
	}
	return tenant, nil
}

// GetTenantList 获取租户列表
func (s *TenantService) GetTenantList(query model.TenantListQuery) (*response.PageData, error) {
	var tenants []model.Tenant
	var total int64

	q := db.GetDB().Model(&model.Tenant{}).Scopes(s.setListScope(query))

	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := q.Scopes(s.paginateScope(query.Page, query.PageSize)).Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, err
	}

	// 转换为列表项并填充用户数
	list := make([]model.TenantListItem, len(tenants))
	for i, t := range tenants {
		var userCount int64
		db.GetDB().Model(&model.User{}).Where("tenant_id = ?", t.ID).Count(&userCount)
		list[i] = model.TenantListItem{
			ID:           t.ID,
			Name:         t.Name,
			Code:         t.Code,
			ContactName:  t.ContactName,
			ContactPhone: t.ContactPhone,
			ContactEmail: t.ContactEmail,
			Domain:       t.Domain,
			APIKey:       t.APIKey,
			PlanType:     t.PlanType,
			MaxUsers:     t.MaxUsers,
			MaxStorage:   t.MaxStorage,
			UsedStorage:  t.UsedStorage,
			ExpireAt:     t.ExpireAt,
			Status:       t.Status,
			Settings:     t.Settings,
			UserCount:    int(userCount),
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		}
	}

	return &response.PageData{
		List:     list,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// GetTenant 获取租户详情
func (s *TenantService) GetTenant(id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := db.GetDB().First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

// UpdateTenant 更新租户
func (s *TenantService) UpdateTenant(id uint, req model.TenantUpdateRequest) error {
	tenant, err := s.GetTenant(id)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.ContactEmail != "" {
		updates["contact_email"] = req.ContactEmail
	}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	if req.PlanType != "" {
		updates["plan_type"] = req.PlanType
	}
	if req.MaxUsers != nil {
		updates["max_users"] = *req.MaxUsers
	}
	if req.MaxStorage != nil {
		updates["max_storage"] = *req.MaxStorage
	}
	if req.ExpireAt != "" {
		expire, err := time.Parse("2006-01-02", req.ExpireAt)
		if err == nil {
			updates["expire_at"] = expire
		}
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Settings != "" {
		updates["settings"] = req.Settings
	}

	return db.GetDB().Model(tenant).Updates(updates).Error
}

// DeleteTenant 删除租户
func (s *TenantService) DeleteTenant(id uint) error {
	return db.GetDB().Delete(&model.Tenant{}, id).Error
}

// UpdateTenantStatus 更新租户状态
func (s *TenantService) UpdateTenantStatus(id uint, status int) error {
	tenant, err := s.GetTenant(id)
	if err != nil {
		return err
	}
	return db.GetDB().Model(tenant).Update("status", status).Error
}

// AddUserToTenant 添加用户到租户
func (s *TenantService) AddUserToTenant(tenantID, userID uint) error {
	tenant, err := s.GetTenant(tenantID)
	if err != nil {
		return errors.New("租户不存在")
	}

	if tenant.Status != 1 {
		return errors.New("租户已停用或已过期")
	}

	// 检查用户数配额
	var userCount int64
	db.GetDB().Model(&model.User{}).Where("tenant_id = ?", tenantID).Count(&userCount)
	if int(userCount) >= tenant.MaxUsers {
		return errors.New("已达到租户最大用户数限制")
	}

	// 更新用户的 tenant_id
	return db.GetDB().Model(&model.User{}).Where("id = ?", userID).Update("tenant_id", tenantID).Error
}

// RemoveUserFromTenant 从租户移除用户
func (s *TenantService) RemoveUserFromTenant(tenantID, userID uint) error {
	return db.GetDB().Model(&model.User{}).Where("id = ? AND tenant_id = ?", userID, tenantID).Update("tenant_id", nil).Error
}

// GetTenantUsers 获取租户用户列表
func (s *TenantService) GetTenantUsers(tenantID uint, page, pageSize int) (*response.PageData, error) {
	var users []model.User
	var total int64

	q := db.GetDB().Model(&model.User{}).Where("tenant_id = ?", tenantID)
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := q.Scopes(s.paginateScope(page, pageSize)).Find(&users).Error; err != nil {
		return nil, err
	}

	userItems := make([]model.TenantUserItem, len(users))
	for i, u := range users {
		userItems[i] = model.TenantUserItem{
			ID:        u.ID,
			Username:  u.Username,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Status:    u.Status,
			CreatedAt: u.CreatedAt,
		}
	}

	return &response.PageData{
		List:     userItems,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetTenantStats 获取租户统计
func (s *TenantService) GetTenantStats() (*model.TenantStats, error) {
	stats := &model.TenantStats{}

	db.GetDB().Model(&model.Tenant{}).Count(&stats.Total)
	db.GetDB().Model(&model.Tenant{}).Where("status = ?", 1).Count(&stats.Active)
	db.GetDB().Model(&model.Tenant{}).Where("status = ?", 0).Count(&stats.Inactive)
	db.GetDB().Model(&model.Tenant{}).Where("status = ?", 2).Count(&stats.Expired)

	// 按套餐类型统计
	type planCount struct {
		PlanType string
		Count    int64
	}
	var results []planCount
	db.GetDB().Select("plan_type, COUNT(*) as count").Model(&model.Tenant{}).Group("plan_type").Find(&results)
	stats.ByPlan = make(map[string]int64)
	for _, r := range results {
		stats.ByPlan[r.PlanType] = r.Count
	}

	return stats, nil
}

func (s *TenantService) setListScope(query model.TenantListQuery) func(db *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		if query.Name != "" {
			d = d.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Code != "" {
			d = d.Where("code = ?", query.Code)
		}
		if query.PlanType != "" {
			d = d.Where("plan_type = ?", query.PlanType)
		}
		if query.Status != nil {
			d = d.Where("status = ?", *query.Status)
		}
		return d
	}
}

func (s *TenantService) paginateScope(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return d.Offset(offset).Limit(pageSize)
	}
}
