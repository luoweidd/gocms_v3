package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// CDNConfig CDN配置
type CDNConfig struct {
	Enabled      bool
	Domain       string
	StaticPrefix string
	VideoPrefix  string
	CacheControl string
}

// DefaultCDNConfig 默认CDN配置
func DefaultCDNConfig() *CDNConfig {
	return &CDNConfig{
		Enabled:      false,
		Domain:       "",
		StaticPrefix: "/static/",
		VideoPrefix:  "/videos/",
		CacheControl: "public, max-age=86400",
	}
}

// CDNMiddleware CDN加速中间件
func CDNMiddleware(config *CDNConfig) gin.HandlerFunc {
	if !config.Enabled || config.Domain == "" {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 静态资源替换为CDN URL
		if strings.HasPrefix(path, config.StaticPrefix) {
			_ = replaceStaticResourceURL(path, config)
			c.Header("X-CDN-Served", "true")
			c.Header("Cache-Control", config.CacheControl)
			c.Next()
			return
		}

		// 视频资源替换为CDN URL
		if strings.HasPrefix(path, config.VideoPrefix) {
			_ = replaceVideoResourceURL(path, config)
			c.Header("X-CDN-Served", "true")
			c.Header("Cache-Control", "public, max-age=86400")
			c.Next()
			return
		}

		c.Next()
	}
}

// replaceStaticResourceURL 替换静态资源URL为CDN URL
func replaceStaticResourceURL(path string, config *CDNConfig) string {
	if !strings.HasPrefix(config.Domain, "http") {
		config.Domain = fmt.Sprintf("https://%s", config.Domain)
	}
	return strings.Replace(path, config.StaticPrefix, fmt.Sprintf("%s/static/", config.Domain), 1)
}

// replaceVideoResourceURL 替换视频资源URL为CDN URL
func replaceVideoResourceURL(path string, config *CDNConfig) string {
	if !strings.HasPrefix(config.Domain, "http") {
		config.Domain = fmt.Sprintf("https://%s", config.Domain)
	}
	return strings.Replace(path, config.VideoPrefix, fmt.Sprintf("%s/videos/", config.Domain), 1)
}

// CDNResponseHelper CDN响应辅助函数
type CDNResponseHelper struct{}

// SetCDNHeader 设置CDN响应头
func (h *CDNResponseHelper) SetCDNHeader(c *gin.Context, path string) {
	if strings.HasPrefix(path, "/static/") || strings.HasPrefix(path, "/videos/") {
		c.Header("Cache-Control", "public, max-age=86400")
		c.Header("X-CDN-Cache", "HIT")
	}
}
