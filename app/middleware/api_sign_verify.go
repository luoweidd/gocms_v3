package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gocms_v3/internal/signature"

	"github.com/gin-gonic/gin"
)

// APISignConfig API签名验证配置
type APISignConfig struct {
	Enabled   bool
	SecretKey string
	MaxAge    int64
	Required  bool
}

// DefaultAPISignConfig 默认配置
var DefaultAPISignConfig = APISignConfig{
	Enabled:  true,
	Required: true,
	MaxAge:   300,
}

// APISignVerifier API签名验证器
var APISignVerifier *signature.Verifier

// InitAPISignVerifier 初始化API签名验证器
func InitAPISignVerifier(secretKey string) {
	APISignVerifier = signature.NewVerifier(signature.VerifierConfig{
		SecretKey:  secretKey,
		MaxAge:     DefaultAPISignConfig.MaxAge,
		RequiredTS: true,
	})
	log.Printf("API签名验证器已初始化")
}

// APISignVerify API签名验证中间件
func APISignVerify(config APISignConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		if path == "/health" || path == "/ready" || path == "/swagger/*any" {
			c.Next()
			return
		}

		params := make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			if len(v) > 0 && k != "sign" && k != "timestamp" {
				params[k] = v[0]
			}
		}

		ts := c.GetHeader("X-Timestamp")
		if ts == "" {
			ts = c.Query("timestamp")
		}
		params["timestamp"] = ts

		sign := c.GetHeader("X-Sign")
		if sign == "" {
			sign = c.Query("sign")
		}
		if sign == "" && config.Required {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "缺少API签名验证信息",
			})
			return
		}

		if ts != "" {
			tsInt, err := strconv.ParseInt(ts, 10, 64)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "无效的时间戳",
				})
				return
			}
			now := time.Now().Unix()
			if now-tsInt > config.MaxAge || tsInt-now > config.MaxAge {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": fmt.Sprintf("请求已过期，允许的最大时间差为%d秒", config.MaxAge),
				})
				return
			}
		}

		if sign != "" && APISignVerifier != nil {
			if err := APISignVerifier.Verify(params, sign); err != nil {
				log.Printf("API签名验证失败: %v", err)
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "签名验证失败",
				})
				return
			}
		}

		c.Next()
	}
}

// ReadBody 读取请求体内容并返回哈希
func ReadBody(c *gin.Context) (string, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "", err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	return fmt.Sprintf("%x", body), nil
}

// APISignMiddleware 统一的API签名中间件入口
func APISignMiddleware(secretKey string, enabled bool) gin.HandlerFunc {
	if secretKey != "" && enabled {
		InitAPISignVerifier(secretKey)
		return APISignVerify(APISignConfig{
			Enabled:   true,
			SecretKey: secretKey,
			MaxAge:    DefaultAPISignConfig.MaxAge,
			Required:  true,
		})
	}
	return func(c *gin.Context) {
		c.Next()
	}
}
