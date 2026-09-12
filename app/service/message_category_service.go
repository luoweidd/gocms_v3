package service

import (
	"fmt"

	"gocms_v3/app/model"

	"gorm.io/gorm"
)

// MessageCategoryService 消息类别服务
type MessageCategoryService struct {
	db *gorm.DB
}

// NewMessageCategoryService 创建消息类别服务
func NewMessageCategoryService(db *gorm.DB) *MessageCategoryService {
	return &MessageCategoryService{db: db}
}

// GetCategories 获取消息类别列表（含消息数量统计）
func (s *MessageCategoryService) GetCategories(query model.MessageCategoryListQuery) ([]model.MessageCategoryResponse, int64, error) {
	var categories []model.MessageCategory
	var total int64

	// 构建查询
	db := s.db.Model(&model.MessageCategory{})

	// 过滤条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		db = db.Where("code = ?", query.Code)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询类别总数失败: %w", err)
	}

	// 排序
	db = db.Order("sort ASC, created_at DESC")

	// 分页
	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&categories).Error; err != nil {
		return nil, 0, fmt.Errorf("查询类别列表失败: %w", err)
	}

	// 组装响应，包含消息数量统计
	resp := make([]model.MessageCategoryResponse, len(categories))
	for i, cat := range categories {
		var msgCount int64
		s.db.Model(&model.SystemMessage{}).Where("message_type = ? AND status = ?", cat.Code, 1).Count(&msgCount)

		resp[i] = model.MessageCategoryResponse{
			ID:           cat.ID,
			Name:         cat.Name,
			Code:         cat.Code,
			Description:  cat.Description,
			IconType:     cat.IconType,
			Sort:         cat.Sort,
			Status:       cat.Status,
			MessageCount: msgCount,
			CreatedAt:    cat.CreatedAt,
			UpdatedAt:    cat.UpdatedAt,
		}
	}

	return resp, total, nil
}

// GetCategoryByID 根据ID获取消息类别
func (s *MessageCategoryService) GetCategoryByID(id uint) (*model.MessageCategoryResponse, error) {
	var cat model.MessageCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, fmt.Errorf("类别不存在: %w", err)
	}

	// 获取消息数量
	var msgCount int64
	s.db.Model(&model.SystemMessage{}).Where("message_type = ? AND status = ?", cat.Code, 1).Count(&msgCount)

	return &model.MessageCategoryResponse{
		ID:           cat.ID,
		Name:         cat.Name,
		Code:         cat.Code,
		Description:  cat.Description,
		IconType:     cat.IconType,
		Sort:         cat.Sort,
		Status:       cat.Status,
		MessageCount: msgCount,
		CreatedAt:    cat.CreatedAt,
		UpdatedAt:    cat.UpdatedAt,
	}, nil
}

// CreateCategory 创建消息类别
func (s *MessageCategoryService) CreateCategory(req model.CreateMessageCategoryRequest) (*model.MessageCategory, error) {
	cat := model.MessageCategory{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		IconType:    req.IconType,
		Sort:        req.Sort,
		Status:      req.Status,
	}

	if cat.Status == 0 {
		cat.Status = 1
	}

	if err := s.db.Create(&cat).Error; err != nil {
		return nil, fmt.Errorf("创建类别失败: %w", err)
	}

	return &cat, nil
}

// UpdateCategory 更新消息类别
func (s *MessageCategoryService) UpdateCategory(id uint, req model.UpdateMessageCategoryRequest) (*model.MessageCategory, error) {
	var cat model.MessageCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, fmt.Errorf("类别不存在: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.IconType != "" {
		updates["icon_type"] = req.IconType
	}
	if req.Sort != 0 {
		updates["sort"] = req.Sort
	}
	if req.Status != 0 {
		updates["status"] = req.Status
	}

	if err := s.db.Model(&cat).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新类别失败: %w", err)
	}

	return &cat, nil
}

// DeleteCategory 删除消息类别
func (s *MessageCategoryService) DeleteCategory(id uint) error {
	var cat model.MessageCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return fmt.Errorf("类别不存在: %w", err)
	}

	// 检查是否有消息使用该类别
	var msgCount int64
	s.db.Model(&model.SystemMessage{}).Where("message_type = ? AND status = ?", cat.Code, 1).Count(&msgCount)
	if msgCount > 0 {
		return fmt.Errorf("该类别下有 %d 条消息，无法删除", msgCount)
	}

	if err := s.db.Delete(&cat, id).Error; err != nil {
		return fmt.Errorf("删除类别失败: %w", err)
	}

	return nil
}

// SeedDefaultCategories 初始化默认消息类别数据
func (s *MessageCategoryService) SeedDefaultCategories() error {
	defaultCategories := []model.CreateMessageCategoryRequest{
		{
			Name:        "系统通知",
			Code:        "system",
			Description: "系统级别的通知消息",
			IconType:    "bell",
			Sort:        1,
			Status:      1,
		},
		{
			Name:        "发布通知",
			Code:        "publish",
			Description: "内容发布相关的通知",
			IconType:    "document",
			Sort:        2,
			Status:      1,
		},
		{
			Name:        "评论回复",
			Code:        "comment_reply",
			Description: "用户评论和回复通知",
			IconType:    "user",
			Sort:        3,
			Status:      1,
		},
		{
			Name:        "审核通过",
			Code:        "audit_passed",
			Description: "内容审核通过通知",
			IconType:    "check",
			Sort:        4,
			Status:      1,
		},
		{
			Name:        "审核驳回",
			Code:        "audited_failed",
			Description: "内容审核未通过通知",
			IconType:    "alert",
			Sort:        5,
			Status:      1,
		},
	}

	for _, cat := range defaultCategories {
		var existing model.MessageCategory
		err := s.db.Where("code = ?", cat.Code).First(&existing).Error
		if err != nil {
			// 不存在则创建
			s.db.Create(&model.MessageCategory{
				Name:        cat.Name,
				Code:        cat.Code,
				Description: cat.Description,
				IconType:    cat.IconType,
				Sort:        cat.Sort,
				Status:      cat.Status,
			})
		}
	}

	return nil
}
