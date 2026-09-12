package model

import "time"

// SeoSetting SEO设置表
type SeoSetting struct {
	BaseModel
	ResourceType   string  `gorm:"column:resource_type;type:varchar(20);not null;index" json:"resource_type"` // article/video/page
	ResourceID     uint    `gorm:"column:resource_id;type:bigint unsigned;not null;index" json:"resource_id"`
	SEOTitle       string  `gorm:"column:meta_title;type:varchar(255)" json:"seo_title"`       // 兼容前端 seo_title
	SEODescription string  `gorm:"column:meta_description;type:text" json:"seo_description"`   // 兼容前端 seo_description
	SEOKeywords    string  `gorm:"column:meta_keywords;type:varchar(500)" json:"seo_keywords"` // 兼容前端 seo_keywords
	OgTitle        string  `gorm:"column:og_title;type:varchar(255)" json:"og_title"`
	OgDescription  string  `gorm:"column:og_description;type:text" json:"og_description"`
	OgImage        string  `gorm:"column:og_image;type:varchar(500)" json:"og_image"`
	CanonicalUrl   string  `gorm:"column:canonical_url;type:varchar(500)" json:"canonical_url"`
	Robots         string  `gorm:"column:robots;type:varchar(20)" json:"robots"` // index,follow/noindex,nofollow
	CustomHeadCode string  `gorm:"column:custom_head_code;type:text" json:"custom_head_code"`
	CustomURL      string  `gorm:"column:custom_url;type:varchar(500)" json:"custom_url"`  // 兼容前端 custom_url
	ChangeFreq     string  `gorm:"column:change_freq;type:varchar(20)" json:"change_freq"` // 兼容前端 change_freq
	Priority       float64 `gorm:"column:priority;type:decimal(3,2)" json:"priority"`      // 兼容前端 priority
}

// TableName 指定表名
func (SeoSetting) TableName() string {
	return "seo_settings"
}

// SitemapConfig 站点地图配置表
type SitemapConfig struct {
	BaseModel
	Name       string  `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Type       string  `gorm:"column:type;type:varchar(20);not null" json:"type"`    // article/video/page/category
	Priority   float64 `gorm:"column:priority;type:decimal(2,1)" json:"priority"`    // 0.0-1.0
	ChangeFreq string  `gorm:"column:changefreq;type:varchar(20)" json:"changefreq"` // always/hourly/daily/weekly/monthly/yearly/never
	IsEnabled  int8    `gorm:"column:is_enabled;type:tinyint;not null;default:1" json:"is_enabled"`
	MaxCount   int     `gorm:"column:max_count;type:int;not null;default:1000" json:"max_count"`
}

// TableName 指定表名
func (SitemapConfig) TableName() string {
	return "sitemap_configs"
}

// SitemapUrl 站点地图URL项
type SitemapUrl struct {
	Loc        string  `xml:"loc" json:"loc"`
	Lastmod    string  `xml:"lastmod" json:"lastmod"`
	ChangeFreq string  `xml:"changefreq" json:"changefreq"`
	Priority   float64 `xml:"priority" json:"priority"`
}

// SeoSettingQuery SEO设置查询参数
type SeoSettingQuery struct {
	ResourceType string `form:"resource_type"`
	Keyword      string `form:"keyword"`
	Page         int    `form:"page" binding:"required"`
	PageSize     int    `form:"page_size" binding:"required"`
}

// SeoSettingResponse SEO设置响应
type SeoSettingResponse struct {
	ID             uint      `json:"id"`
	ResourceType   string    `json:"resource_type"`
	ResourceID     uint      `json:"resource_id"`
	SEOTitle       string    `json:"seo_title"`
	SEODescription string    `json:"seo_description"`
	SEOKeywords    string    `json:"seo_keywords"`
	OgTitle        string    `json:"og_title"`
	OgDescription  string    `json:"og_description"`
	OgImage        string    `json:"og_image"`
	CanonicalUrl   string    `json:"canonical_url"`
	Robots         string    `json:"robots"`
	CustomHeadCode string    `json:"custom_head_code"`
	CustomURL      string    `json:"custom_url"`
	ChangeFreq     string    `json:"change_freq"`
	Priority       float64   `json:"priority"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// SeoCreateRequest SEO创建请求
type SeoCreateRequest struct {
	ResourceType   string  `form:"resource_type" binding:"required"`
	ResourceID     uint    `form:"resource_id" binding:"required"`
	SEOTitle       string  `form:"seo_title"`
	SEODescription string  `form:"seo_description"`
	SEOKeywords    string  `form:"seo_keywords"`
	OgTitle        string  `form:"og_title"`
	OgDescription  string  `form:"og_description"`
	OgImage        string  `form:"og_image"`
	CanonicalUrl   string  `form:"canonical_url"`
	Robots         string  `form:"robots"`
	CustomHeadCode string  `form:"custom_head_code"`
	CustomURL      string  `form:"custom_url"`
	ChangeFreq     string  `form:"change_freq"`
	Priority       float64 `form:"priority"`
}

// SeoUpdateRequest SEO更新请求
type SeoUpdateRequest struct {
	SEOTitle       *string  `form:"seo_title"`
	SEODescription *string  `form:"seo_description"`
	SEOKeywords    *string  `form:"seo_keywords"`
	OgTitle        *string  `form:"og_title"`
	OgDescription  *string  `form:"og_description"`
	OgImage        *string  `form:"og_image"`
	CanonicalUrl   *string  `form:"canonical_url"`
	Robots         *string  `form:"robots"`
	CustomHeadCode *string  `form:"custom_head_code"`
	CustomURL      *string  `form:"custom_url"`
	ChangeFreq     *string  `form:"change_freq"`
	Priority       *float64 `form:"priority"`
}

// SitemapItemResponse 站点地图项响应
type SitemapItemResponse struct {
	URLs        []SitemapUrl `json:"urls"`
	TotalCount  int64        `json:"total_count"`
	GeneratedAt string       `xml:"generated_at" json:"generated_at"`
}

// SitemapStats 站点地图统计
type SitemapStats struct {
	TotalArticles int64  `json:"total_articles"`
	TotalVideos   int64  `json:"total_videos"`
	TotalPages    int64  `json:"total_pages"`
	LastGenerated string `json:"last_generated"`
}
