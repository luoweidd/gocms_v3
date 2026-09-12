package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gocms_v3/app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SeoService SEO服务
type SeoService struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewSeoService 创建SEO服务
func NewSeoService(db *gorm.DB) *SeoService {
	return &SeoService{
		db:  db,
		log: zap.L(),
	}
}

// CreateSeoSetting 创建SEO设置
func (s *SeoService) CreateSeoSetting(ctx context.Context, req model.SeoCreateRequest) (*model.SeoSetting, error) {
	// 检查是否已存在
	var existing model.SeoSetting
	if err := s.db.Where("resource_type = ? AND resource_id = ?", req.ResourceType, req.ResourceID).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("该资源的SEO设置已存在")
	}

	setting := &model.SeoSetting{
		ResourceType:   req.ResourceType,
		ResourceID:     req.ResourceID,
		SEOTitle:       req.SEOTitle,
		SEODescription: req.SEODescription,
		SEOKeywords:    req.SEOKeywords,
		OgTitle:        req.OgTitle,
		OgDescription:  req.OgDescription,
		OgImage:        req.OgImage,
		CanonicalUrl:   req.CanonicalUrl,
		Robots:         req.Robots,
		CustomHeadCode: req.CustomHeadCode,
		CustomURL:      req.CustomURL,
		ChangeFreq:     req.ChangeFreq,
		Priority:       req.Priority,
	}

	if err := s.db.Create(setting).Error; err != nil {
		return nil, fmt.Errorf("创建SEO设置失败: %w", err)
	}

	return setting, nil
}

// UpdateSeoSetting 更新SEO设置
func (s *SeoService) UpdateSeoSetting(id uint, req model.SeoUpdateRequest) (*model.SeoSetting, error) {
	var setting model.SeoSetting
	if err := s.db.First(&setting, id).Error; err != nil {
		return nil, fmt.Errorf("SEO设置不存在: %w", err)
	}

	updates := make(map[string]interface{})
	if req.SEOTitle != nil {
		updates["meta_title"] = *req.SEOTitle
	}
	if req.SEODescription != nil {
		updates["meta_description"] = *req.SEODescription
	}
	if req.SEOKeywords != nil {
		updates["meta_keywords"] = *req.SEOKeywords
	}
	if req.OgTitle != nil {
		updates["og_title"] = *req.OgTitle
	}
	if req.OgDescription != nil {
		updates["og_description"] = *req.OgDescription
	}
	if req.OgImage != nil {
		updates["og_image"] = *req.OgImage
	}
	if req.CanonicalUrl != nil {
		updates["canonical_url"] = *req.CanonicalUrl
	}
	if req.Robots != nil {
		updates["robots"] = *req.Robots
	}
	if req.CustomHeadCode != nil {
		updates["custom_head_code"] = *req.CustomHeadCode
	}

	if err := s.db.Model(&setting).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新SEO设置失败: %w", err)
	}

	return &setting, nil
}

// GetSeoSetting 获取SEO设置
func (s *SeoService) GetSeoSetting(id uint) (*model.SeoSettingResponse, error) {
	var setting model.SeoSetting
	if err := s.db.First(&setting, id).Error; err != nil {
		return nil, fmt.Errorf("SEO设置不存在: %w", err)
	}

	return &model.SeoSettingResponse{
		ID:             setting.ID,
		ResourceType:   setting.ResourceType,
		ResourceID:     setting.ResourceID,
		SEOTitle:       setting.SEOTitle,
		SEODescription: setting.SEODescription,
		SEOKeywords:    setting.SEOKeywords,
		OgTitle:        setting.OgTitle,
		OgDescription:  setting.OgDescription,
		OgImage:        setting.OgImage,
		CanonicalUrl:   setting.CanonicalUrl,
		Robots:         setting.Robots,
		CustomHeadCode: setting.CustomHeadCode,
		CustomURL:      setting.CustomURL,
		ChangeFreq:     setting.ChangeFreq,
		Priority:       setting.Priority,
		CreatedAt:      setting.CreatedAt,
		UpdatedAt:      setting.UpdatedAt,
	}, nil
}

// GetSeoSettingByResource 根据资源获取SEO设置
func (s *SeoService) GetSeoSettingByResource(resourceType string, resourceID uint) (*model.SeoSettingResponse, error) {
	var setting model.SeoSetting
	if err := s.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).First(&setting).Error; err != nil {
		return nil, nil // 不存在返回nil而非错误
	}

	return &model.SeoSettingResponse{
		ID:             setting.ID,
		ResourceType:   setting.ResourceType,
		ResourceID:     setting.ResourceID,
		SEOTitle:       setting.SEOTitle,
		SEODescription: setting.SEODescription,
		SEOKeywords:    setting.SEOKeywords,
		OgTitle:        setting.OgTitle,
		OgDescription:  setting.OgDescription,
		OgImage:        setting.OgImage,
		CanonicalUrl:   setting.CanonicalUrl,
		Robots:         setting.Robots,
		CustomHeadCode: setting.CustomHeadCode,
		CustomURL:      setting.CustomURL,
		ChangeFreq:     setting.ChangeFreq,
		Priority:       setting.Priority,
		CreatedAt:      setting.CreatedAt,
		UpdatedAt:      setting.UpdatedAt,
	}, nil
}

// GetSeoList 获取SEO设置列表
func (s *SeoService) GetSeoList(query model.SeoSettingQuery) ([]model.SeoSettingResponse, int64, error) {
	var settings []model.SeoSetting
	var total int64

	queryBuilder := s.db.Model(&model.SeoSetting{}).Order("updated_at DESC")

	if query.ResourceType != "" {
		queryBuilder = queryBuilder.Where("resource_type = ?", query.ResourceType)
	}
	if query.Keyword != "" {
		queryBuilder = queryBuilder.Where("meta_title LIKE ? OR meta_description LIKE ?", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询SEO总数失败: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize
	if err := queryBuilder.Offset(offset).Limit(query.PageSize).Find(&settings).Error; err != nil {
		return nil, 0, fmt.Errorf("查询SEO列表失败: %w", err)
	}

	resp := make([]model.SeoSettingResponse, len(settings))
	for i, s := range settings {
		resp[i] = model.SeoSettingResponse{
			ID:             s.ID,
			ResourceType:   s.ResourceType,
			ResourceID:     s.ResourceID,
			SEOTitle:       s.SEOTitle,
			SEODescription: s.SEODescription,
			SEOKeywords:    s.SEOKeywords,
			OgTitle:        s.OgTitle,
			OgDescription:  s.OgDescription,
			OgImage:        s.OgImage,
			CanonicalUrl:   s.CanonicalUrl,
			Robots:         s.Robots,
			CustomHeadCode: s.CustomHeadCode,
			CustomURL:      s.CustomURL,
			ChangeFreq:     s.ChangeFreq,
			Priority:       s.Priority,
			CreatedAt:      s.CreatedAt,
			UpdatedAt:      s.UpdatedAt,
		}
	}

	return resp, total, nil
}

// DeleteSeoSetting 删除SEO设置
func (s *SeoService) DeleteSeoSetting(id uint) error {
	return s.db.Delete(&model.SeoSetting{}, id).Error
}

// GetSitemapConfigList 获取站点地图配置列表
func (s *SeoService) GetSitemapConfigList() ([]model.SitemapConfig, error) {
	var configs []model.SitemapConfig
	if err := s.db.Where("is_enabled = ?", 1).Order("type, id").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("查询站点地图配置失败: %w", err)
	}
	return configs, nil
}

// GenerateSitemap 生成站点地图XML
func (s *SeoService) GenerateSitemap() (string, error) {
	var configs []model.SitemapConfig
	if err := s.db.Where("is_enabled = ?", 1).Order("type, priority DESC").Find(&configs).Error; err != nil {
		s.log.Warn("查询站点地图配置失败", zap.Error(err))
		// 即使没有配置，也返回基本的XML
		return s.generateBasicSitemap(), nil
	}

	var urls []model.SitemapUrl

	for _, config := range configs {
		var items []map[string]interface{}
		switch config.Type {
		case "article":
			type ArticleItem struct {
				ID        uint      `gorm:"column:id"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}
			var articles []ArticleItem
			err := s.db.Table("articles").Select("id, updated_at").Order("updated_at DESC").Limit(config.MaxCount).Find(&articles).Error
			if err != nil {
				s.log.Warn("查询文章失败", zap.Error(err), zap.String("type", config.Type))
				continue
			}
			for _, a := range articles {
				items = append(items, map[string]interface{}{
					"loc":        fmt.Sprintf("/api/articles/%d", a.ID),
					"lastmod":    a.UpdatedAt.Format("2006-01-02"),
					"changefreq": config.ChangeFreq,
					"priority":   config.Priority,
				})
			}
		case "video":
			type VideoItem struct {
				ID        uint      `gorm:"column:id"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}
			var videos []VideoItem
			err := s.db.Table("videos").Select("id, updated_at").Order("updated_at DESC").Limit(config.MaxCount).Find(&videos).Error
			if err != nil {
				s.log.Warn("查询视频失败", zap.Error(err), zap.String("type", config.Type))
				continue
			}
			for _, v := range videos {
				items = append(items, map[string]interface{}{
					"loc":        fmt.Sprintf("/api/videos/%d", v.ID),
					"lastmod":    v.UpdatedAt.Format("2006-01-02"),
					"changefreq": config.ChangeFreq,
					"priority":   config.Priority,
				})
			}
		case "page":
			// 页面是虚拟的，暂时跳过或可以从其他表查询
			s.log.Info("页面类型暂未实现")
		case "category":
			type CategoryItem struct {
				ID        uint      `gorm:"column:id"`
				UpdatedAt time.Time `gorm:"column:updated_at"`
			}
			var categories []CategoryItem
			err := s.db.Table("article_categories").Select("id, updated_at").Order("updated_at DESC").Limit(config.MaxCount).Find(&categories).Error
			if err != nil {
				s.log.Warn("查询分类失败", zap.Error(err), zap.String("type", config.Type))
				continue
			}
			for _, c := range categories {
				items = append(items, map[string]interface{}{
					"loc":        fmt.Sprintf("/api/article-categories/%d", c.ID),
					"lastmod":    c.UpdatedAt.Format("2006-01-02"),
					"changefreq": config.ChangeFreq,
					"priority":   config.Priority,
				})
			}
		default:
			s.log.Info("未知的站点地图类型", zap.String("type", config.Type))
			continue
		}

		for _, item := range items {
			urls = append(urls, model.SitemapUrl{
				Loc:        item["loc"].(string),
				Lastmod:    item["lastmod"].(string),
				ChangeFreq: item["changefreq"].(string),
				Priority:   item["priority"].(float64),
			})
		}
	}

	xml := s.generateSitemapXML(urls)
	return xml, nil
}

// generateBasicSitemap 生成基本的站点地图XML
func (s *SeoService) generateBasicSitemap() string {
	s.log.Info("生成基本站点地图（无配置）")

	// 添加一些基本的URL
	var urls []model.SitemapUrl

	// 添加首页
	urls = append(urls, model.SitemapUrl{
		Loc:        "/",
		Lastmod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// 添加文章列表页
	var articleCount int64
	s.db.Model(&model.Article{}).Count(&articleCount)
	if articleCount > 0 {
		var latestArticle model.Article
		s.db.Order("updated_at DESC").First(&latestArticle)
		urls = append(urls, model.SitemapUrl{
			Loc:        "/api/articles",
			Lastmod:    latestArticle.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "daily",
			Priority:   0.8,
		})
	}

	return s.generateSitemapXML(urls)
}

// generateSitemapXML 生成XML格式
func (s *SeoService) generateSitemapXML(urls []model.SitemapUrl) string {
	var sb strings.Builder

	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, url := range urls {
		sb.WriteString("  <url>\n")
		sb.WriteString(fmt.Sprintf("    <loc>%s</loc>\n", url.Loc))
		if url.Lastmod != "" {
			sb.WriteString(fmt.Sprintf("    <lastmod>%s</lastmod>\n", url.Lastmod))
		}
		if url.ChangeFreq != "" {
			sb.WriteString(fmt.Sprintf("    <changefreq>%s</changefreq>\n", url.ChangeFreq))
		}
		if url.Priority > 0 {
			sb.WriteString(fmt.Sprintf("    <priority>%.1f</priority>\n", url.Priority))
		}
		sb.WriteString("  </url>\n")
	}

	sb.WriteString(`</urlset>`)
	return sb.String()
}

// GetSitemapStats 获取站点地图统计
func (s *SeoService) GetSitemapStats() (*model.SitemapStats, error) {
	stats := &model.SitemapStats{}

	s.db.Model(&model.Article{}).Count(&stats.TotalArticles)
	s.db.Model(&model.Video{}).Count(&stats.TotalVideos)

	// Pages是虚拟的，统计顶级页面数
	stats.TotalPages = 0

	return stats, nil
}

// GetSeoByResourceId 根据资源ID和类型获取SEO设置
func (s *SeoService) GetSeoByResourceId(resourceType string, resourceID uint) (*model.SeoSetting, error) {
	var setting model.SeoSetting
	if err := s.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
