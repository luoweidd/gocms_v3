package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// 处理请求
		c.Next()

		// 记录请求信息
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("[REQUEST] method=%s path=%s status=%d latency=%v client_ip=%s",
			c.Request.Method,
			path,
			status,
			latency,
			c.ClientIP(),
		)

		// 如果有错误，记录错误信息
		if status >= 500 {
			log.Printf("[REQUEST-ERROR] method=%s path=%s status=%d error=%v",
				c.Request.Method,
				path,
				status,
				c.Errors.String(),
			)
		}
	}
}
