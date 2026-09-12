package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware 多租户支持中间件
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取租户ID
		tenantIDStr := c.GetHeader("X-Tenant-ID")
		if tenantIDStr == "" {
			// 从查询参数获取
			tenantIDStr = c.Query("tenant_id")
		}

		if tenantIDStr != "" {
			tenantID, err := strconv.ParseUint(tenantIDStr, 10, 64)
			if err == nil && tenantID > 0 {
				c.Set("tenant_id", uint(tenantID))
			}
		}

		c.Next()
	}
}

// GetTenantID 从上下文中获取租户ID
func GetTenantID(c *gin.Context) uint {
	if v, exists := c.Get("tenant_id"); exists {
		return v.(uint)
	}
	return 0
}

// TenantIsolation 数据隔离中间件 - 自动添加tenant_id条件
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := GetTenantID(c)
		if tenantID > 0 {
			c.Set("tenant_filter", " AND tenant_id = "+strconv.FormatUint(uint64(tenantID), 10))
		} else {
			c.Set("tenant_filter", "")
		}
		c.Next()
	}
}
