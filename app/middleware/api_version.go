package middleware

import (
	"github.com/gin-gonic/gin"
)

// APIVersionMiddleware API版本控制中间件
func APIVersionMiddleware(version string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置API版本到上下文
		c.Set("api_version", version)

		// 在响应头中添加API版本信息
		c.Header("X-API-Version", version)

		c.Next()
	}
}

// GetAPIVersion 从上下文中获取API版本
func GetAPIVersion(c *gin.Context) string {
	if v, exists := c.Get("api_version"); exists {
		return v.(string)
	}
	return "v1"
}
