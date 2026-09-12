package service

import (
	"errors"
	"fmt"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
)

// VideoCategoryService 视频分类服务
type VideoCategoryService struct{}

// NewVideoCategoryService 获取视频分类服务实例
func NewVideoCategoryService() *VideoCategoryService {
	return &VideoCategoryService{}
}

// CreateVideoCategoryRequest 创建分类请求
type CreateVideoCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Pid         *uint  `json:"pid"`
	Sort        int    `json:"sort"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// UpdateVideoCategoryRequest 更新分类请求
type UpdateVideoCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Pid         *uint   `json:"pid"`
	Sort        *int    `json:"sort"`
	Icon        *string `json:"icon"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
}

// Create 创建分类
func (s *VideoCategoryService) Create(req CreateVideoCategoryRequest) (*model.VideoCategory, error) {
	// 检查名称是否重复
	var count int64
	db.GetDB().Model(&model.VideoCategory{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("分类名称已存在")
	}

	// 检查父分类是否存在
	if req.Pid != nil && *req.Pid > 0 {
		var parent model.VideoCategory
		if err := db.GetDB().First(&parent, *req.Pid).Error; err != nil {
			return nil, errors.New("父分类不存在")
		}
	}

	if req.Slug == "" {
		req.Slug = req.Name
	}

	category := model.VideoCategory{
		Name:        req.Name,
		Slug:        req.Slug,
		Pid:         req.Pid,
		Sort:        req.Sort,
		Icon:        req.Icon,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := db.GetDB().Create(&category).Error; err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}

	return &category, nil
}

// Update 更新分类
func (s *VideoCategoryService) Update(id uint, req UpdateVideoCategoryRequest) (*model.VideoCategory, error) {
	var category model.VideoCategory
	if err := db.GetDB().First(&category, id).Error; err != nil {
		return nil, errors.New("分类不存在")
	}

	// 不能将自己设为父分类
	if req.Pid != nil && *req.Pid == id {
		return nil, errors.New("不能将自身设为父分类")
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		// 检查名称重复
		var count int64
		db.GetDB().Model(&model.VideoCategory{}).Where("name = ? AND id != ?", *req.Name, id).Count(&count)
		if count > 0 {
			return nil, errors.New("分类名称已存在")
		}
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Pid != nil {
		updates["pid"] = *req.Pid
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := db.GetDB().Model(&category).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}

	return &category, nil
}

// Delete 删除分类（有子分类或关联视频时不可删除）
func (s *VideoCategoryService) Delete(id uint) error {
	var category model.VideoCategory
	if err := db.GetDB().First(&category, id).Error; err != nil {
		return errors.New("分类不存在")
	}

	// 检查是否有子分类
	var childCount int64
	db.GetDB().Model(&model.VideoCategory{}).Where("pid = ?", id).Count(&childCount)
	if childCount > 0 {
		return errors.New("该分类下存在子分类，无法删除")
	}

	// 检查是否有关联视频
	var videoCount int64
	db.GetDB().Model(&model.Video{}).Where("category_id = ?", id).Count(&videoCount)
	if videoCount > 0 {
		return errors.New("该分类下存在视频，无法删除")
	}

	return db.GetDB().Delete(&model.VideoCategory{}, id).Error
}

// GetByID 根据ID获取分类
func (s *VideoCategoryService) GetByID(id uint) (*model.VideoCategory, error) {
	var category model.VideoCategory
	if err := db.GetDB().Preload("Children").First(&category, id).Error; err != nil {
		return nil, errors.New("分类不存在")
	}
	return &category, nil
}

// GetTree 获取分类树
func (s *VideoCategoryService) GetTree(includeDisabled bool) ([]model.VideoCategory, error) {
	var categories []model.VideoCategory
	query := db.GetDB().Model(&model.VideoCategory{}).Order("sort ASC, id ASC")
	if !includeDisabled {
		query = query.Where("status = 1")
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	// 构建树形结构
	return buildTree(categories), nil
}

// buildTree 构建分类树
func buildTree(categories []model.VideoCategory) []model.VideoCategory {
	categoryMap := make(map[uint]*model.VideoCategory)
	var roots []model.VideoCategory

	for i := range categories {
		categoryMap[categories[i].ID] = &categories[i]
	}

	for i := range categories {
		if categories[i].Pid == nil || *categories[i].Pid == 0 {
			roots = append(roots, categories[i])
		} else {
			if parent, ok := categoryMap[*categories[i].Pid]; ok {
				parent.Children = append(parent.Children, categories[i])
			}
		}
	}

	return roots
}

// GetBySlug 根据Slug获取分类
func (s *VideoCategoryService) GetBySlug(slug string) (*model.VideoCategory, error) {
	var category model.VideoCategory
	if err := db.GetDB().Where("slug = ?", slug).First(&category).Error; err != nil {
		return nil, errors.New("分类不存在")
	}
	return &category, nil
}

// UpdateVideoCount 更新分类视频数量统计
func (s *VideoCategoryService) UpdateVideoCount(categoryID uint) error {
	var count int64
	db.GetDB().Model(&model.Video{}).Where("category_id = ?", categoryID).Count(&count)
	return db.GetDB().Model(&model.VideoCategory{}).Where("id = ?", categoryID).Update("video_count", count).Error
}

// RecalculateVideoCounts 重新计算所有分类的视频数量
func RecalculateVideoCounts() {
	var categories []model.VideoCategory
	db.GetDB().Find(&categories)
	for _, cat := range categories {
		var count int64
		db.GetDB().Model(&model.Video{}).Where("category_id = ?", cat.ID).Count(&count)
		db.GetDB().Model(&model.VideoCategory{}).Where("id = ?", cat.ID).Update("video_count", count)
	}
}

// GetVideoCount 获取分类下的视频数量
func (s *VideoCategoryService) GetVideoCount(categoryID uint) (int64, error) {
	var count int64
	err := db.GetDB().Model(&model.Video{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}
