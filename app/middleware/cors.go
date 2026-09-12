package middleware

import (
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS配置
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig 默认CORS配置 - 使用白名单而非"*"
var DefaultCORSConfig = CORSConfig{
	AllowOrigins: []string{
		"http://localhost:3000",
		"http://localhost:5173",
		"http://localhost:8080",
		"https://example.com",
	},
	AllowMethods: []string{
		"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
	},
	AllowHeaders: []string{
		"Origin", "Content-Type", "Accept", "Authorization",
		"X-Requested-With", "X-Sign", "X-Timestamp", "X-Tenant-ID",
	},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           86400,
}

// CORS CORS中间件 - 使用白名单替代"*"
func CORS(cfg CORSConfig) gin.HandlerFunc {
	// 如果没有配置允许的来源，使用默认配置
	if len(cfg.AllowOrigins) == 0 {
		cfg = DefaultCORSConfig
	}

	log.Printf("[CORS] CORS中间件已启用, 配置的允许源: %v", cfg.AllowOrigins)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// 如果没有Origin头（可能是同源请求或浏览器内部请求），直接放行
		if origin == "" {
			c.Next()
			return
		}

		// 检查来源是否在允许列表中
		allowed := false
		var allowedOrigin string
		for _, allow := range cfg.AllowOrigins {
			if allow == "*" || allow == origin {
				allowed = true
				allowedOrigin = origin
				break
			}
			if strings.HasSuffix(allow, "*") {
				prefix := strings.Split(allow, "*")[0]
				if strings.HasPrefix(origin, prefix) {
					allowed = true
					allowedOrigin = origin
					break
				}
			}
		}

		// 如果来源不在允许列表中，拒绝请求
		if !allowed {
			log.Printf("[CORS] 拒绝来源: %s (当前配置: %v)", origin, cfg.AllowOrigins)
			c.AbortWithStatus(403)
			return
		}

		// 处理 OPTIONS 预检请求
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowMethods, ", "))
			c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowHeaders, ", "))
			if cfg.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			if cfg.MaxAge > 0 {
				c.Header("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
			}
			c.AbortWithStatus(204)
			return
		}

		// 为正常请求添加CORS头
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		if cfg.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Next()
	}
}
