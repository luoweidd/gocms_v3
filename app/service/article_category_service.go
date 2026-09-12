package service

import (
	"fmt"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"
)

type ArticleCategoryService struct{}

func NewArticleCategoryService() *ArticleCategoryService {
	return &ArticleCategoryService{}
}

type CreateArticleCategoryRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Icon        string `json:"icon"`
	Sort        int    `json:"sort"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type UpdateArticleCategoryRequest struct {
	Name        *string `json:"name"`
	Icon        *string `json:"icon"`
	Sort        *int    `json:"sort"`
	Description *string `json:"description"`
	Status      *int    `json:"status"`
}

type ArticleCategoryQueryParams struct {
	Status          int  `form:"status"`
	IncludeDisabled bool `form:"include_disabled"`
}

// ListArticleCategoryRequest 列表查询请求
type ListArticleCategoryRequest struct {
	Status int `form:"status"`
}

// TreeArticleCategoryRequest 树查询请求
type TreeArticleCategoryRequest struct {
	Type string `form:"type"`
}

// ListArticleCategory 获取分类列表（兼容controller调用）
func (s *ArticleCategoryService) ListArticleCategory(req ListArticleCategoryRequest) ([]model.ArticleCategory, error) {
	var categories []model.ArticleCategory
	query := db.GetDB().Model(&model.ArticleCategory{})
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}
	if err := query.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// TreeArticleCategory 获取分类树（兼容controller调用）
func (s *ArticleCategoryService) TreeArticleCategory(req TreeArticleCategoryRequest) ([]model.ArticleCategory, error) {
	var categories []model.ArticleCategory
	query := db.GetDB().Model(&model.ArticleCategory{}).Where("status = 1")
	if err := query.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return buildMenuTreeGeneric(categories), nil
}

func buildMenuTreeGeneric(cats []model.ArticleCategory) []model.ArticleCategory {
	type catTree struct {
		model.ArticleCategory
		Children []model.ArticleCategory `json:"children"`
	}
	return cats
}

type ArticleFilterParams struct {
	CategoryID   uint       `form:"category_id"`
	CategoryType string     `form:"category_type"`
	AuthorID     uint       `form:"author_id"`
	AuthorName   string     `form:"author_name"`
	StartTime    *time.Time `form:"start_time"`
	EndTime      *time.Time `form:"end_time"`
	Type         string     `form:"type"`
	Status       int        `form:"status"`
	IsTop        *bool      `form:"is_top"`
	Page         int        `form:"page"`
	PageSize     int        `form:"page_size"`
}

func (s *ArticleCategoryService) CreateArticleCategory(req CreateArticleCategoryRequest) (*model.ArticleCategory, error) {
	if req.Status == 0 {
		req.Status = 1
	}

	category := model.ArticleCategory{
		Name:        req.Name,
		Icon:        req.Icon,
		SortOrder:   req.Sort,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := db.GetDB().Create(&category).Error; err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}

	return &category, nil
}

func (s *ArticleCategoryService) UpdateArticleCategory(id uint, req UpdateArticleCategoryRequest) (*model.ArticleCategory, error) {
	var category model.ArticleCategory
	if err := db.GetDB().First(&category, id).Error; err != nil {
		return nil, fmt.Errorf("分类不存在")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Sort != nil {
		updates["sort_order"] = *req.Sort
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return &category, nil
	}

	if err := db.GetDB().Model(&category).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}

	db.GetDB().First(&category, id)
	return &category, nil
}

func (s *ArticleCategoryService) DeleteArticleCategory(id uint) error {
	var category model.ArticleCategory
	if err := db.GetDB().First(&category, id).Error; err != nil {
		return fmt.Errorf("分类不存在")
	}

	// 当前 article_categories 表没有 pid 字段，暂时跳过子分类检查
	// 后续如需支持树形结构，需要给表添加 pid 字段

	var articleCount int64
	db.GetDB().Model(&model.Article{}).Where("category_id = ?", id).Count(&articleCount)
	if articleCount > 0 {
		return fmt.Errorf("该分类下存在文章，无法删除")
	}

	return db.GetDB().Delete(&model.ArticleCategory{}, id).Error
}

func (s *ArticleCategoryService) GetArticleCategoryTree(params ArticleCategoryQueryParams) ([]model.ArticleCategory, error) {
	var categories []model.ArticleCategory
	query := db.GetDB().Model(&model.ArticleCategory{}).Order("sort_order ASC, id ASC")

	if !params.IncludeDisabled && params.Status == 0 {
		query = query.Where("status = 1")
	} else if params.Status > 0 {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("查询分类失败: %w", err)
	}

	// 构建树形结构（article_categories 表目前没有 pid 字段，所有分类都是顶级）
	return categories, nil
}

func (s *ArticleCategoryService) GetArticleCategoryByID(id uint) (*model.ArticleCategory, error) {
	var category model.ArticleCategory
	if err := db.GetDB().Preload("Children").First(&category, id).Error; err != nil {
		return nil, fmt.Errorf("分类不存在")
	}
	return &category, nil
}

func (s *ArticleCategoryService) GetArticleByFilters(params ArticleFilterParams) (*response.PageData, error) {
	var articles []model.Article
	var total int64
	query := db.GetDB().Model(&model.Article{}).Where("deleted_at IS NULL")

	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.AuthorID > 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	if params.AuthorName != "" {
		query = query.Where("author LIKE ?", "%"+params.AuthorName+"%")
	}
	if params.StartTime != nil {
		query = query.Where("published_at >= ?", *params.StartTime)
	}
	if params.EndTime != nil {
		query = query.Where("published_at <= ?", *params.EndTime)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status >= 0 {
		query = query.Where("status = ?", params.Status)
	}
	if params.IsTop != nil {
		query = query.Where("is_top = ?", *params.IsTop)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("统计文章数量失败: %w", err)
	}

	page := params.Page
	pageSize := params.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	if err := query.Preload("Tags").Offset(offset).Limit(pageSize).Order("is_top DESC, published_at DESC").Find(&articles).Error; err != nil {
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	return &response.PageData{List: articles, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *ArticleCategoryService) GetArticleStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var statusCount []struct {
		Status int
		Count  int64
	}
	db.GetDB().Model(&model.Article{}).Select("status, COUNT(*) as count").Group("status").Find(&statusCount)
	for _, sc := range statusCount {
		switch sc.Status {
		case 0:
			stats["draft_count"] = sc.Count
		case 1:
			stats["published_count"] = sc.Count
		case 2:
			stats["taken_down_count"] = sc.Count
		}
	}

	// 按分类统计
	var categoryCount []struct {
		CategoryID uint
		Count      int64
	}
	db.GetDB().Table("articles a").Where("a.deleted_at IS NULL").Group("a.category_id").Find(&categoryCount)
	categoryCounts := make(map[uint]int64)
	for _, cc := range categoryCount {
		categoryCounts[cc.CategoryID] = cc.Count
	}
	stats["category_counts"] = categoryCounts

	var authorCount []struct {
		Author uint
		Count  int64
	}
	db.GetDB().Model(&model.Article{}).Select("author_id, COUNT(*) as count").Where("deleted_at IS NULL").Group("author_id").Find(&authorCount)
	authorCounts := make(map[uint]int64)
	for _, ac := range authorCount {
		authorCounts[ac.Author] = ac.Count
	}
	stats["author_counts"] = authorCounts

	return stats, nil
}
