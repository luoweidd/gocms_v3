package middleware

import (
	"net/http"
	"strconv"
	"time"

	"gocms_v3/app/response"
	"gocms_v3/internal/signature"

	"github.com/gin-gonic/gin"
)

// APIConfig API签名验证配置
type APIConfig struct {
	Enabled   bool   `yaml:"enabled"`
	SecretKey string `yaml:"secret_key"`
	MaxAge    int64  `yaml:"max_age"` // 默认5分钟
	KeyPrefix string `yaml:"key_prefix"`
}

// DefaultAPIConfig 默认配置
var DefaultAPIConfig = APIConfig{
	Enabled:   false,
	SecretKey: "change-me-to-secure-key",
	MaxAge:    300, // 5分钟
	KeyPrefix: "api_sign",
}

// APISigner 全局API签名验证器
var APISigner *signature.Signer

// InitAPISigner 初始化API签名验证器
func InitAPISigner(cfg APIConfig) {
	if cfg.SecretKey != "" {
		APISigner = signature.NewSigner(cfg.SecretKey)
	}
}

// APISignValidator API签名验证中间件
func APISignValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !DefaultAPIConfig.Enabled {
			c.Next()
			return
		}

		// 跳过不需要签名的接口
		path := c.Request.URL.Path
		if path == "/swagger/*any" || path == "/health" || path == "/api/v2/openapi.json" {
			c.Next()
			return
		}

		// 获取请求参数
		sign := c.GetHeader("X-API-Sign")
		timestampStr := c.GetHeader("X-API-Timestamp")
		appID := c.GetHeader("X-API-App-ID")

		if sign == "" || timestampStr == "" {
			response.Error(c, http.StatusUnauthorized, "缺少签名参数")
			c.Abort()
			return
		}

		// 验证时间戳
		timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, "无效的时间戳")
			c.Abort()
			return
		}

		// 检查时间是否在允许范围内
		if time.Now().Unix()-timestamp > DefaultAPIConfig.MaxAge {
			response.Error(c, http.StatusUnauthorized, "请求已过期")
			c.Abort()
			return
		}

		// 构建签名字典
		params := make(map[string]string)
		params["app_id"] = appID
		params["timestamp"] = timestampStr

		// 获取请求体参数
		if c.Request.URL.RawQuery != "" {
			for _, part := range splitQueryString(c.Request.URL.RawQuery) {
				kv := splitKV(part)
				if len(kv) == 2 {
					params[kv[0]] = kv[1]
				}
			}
		}

		// 验证签名
		if err := APISigner.VerifySign(params, sign, DefaultAPIConfig.MaxAge); err != nil {
			response.Error(c, http.StatusUnauthorized, "签名验证失败")
			c.Abort()
			return
		}

		// 设置APP ID到上下文
		c.Set("app_id", appID)
		c.Next()
	}
}

// splitQueryString 分割查询字符串
func splitQueryString(query string) []string {
	var parts []string
	start := 0
	for i, c := range query {
		if c == '&' {
			parts = append(parts, query[start:i])
			start = i + 1
		}
	}
	if start < len(query) {
		parts = append(parts, query[start:])
	}
	return parts
}

// splitKV 分割键值对
func splitKV(s string) []string {
	for i, c := range s {
		if c == '=' {
			return []string{s[:i], s[i+1:]}
		}
	}
	return []string{s, ""}
}

// GenerateAPISign 生成API签名（供调用方使用）
func GenerateAPISign(params map[string]string, secretKey string) (string, int64) {
	signer := signature.NewSigner(secretKey)
	return signer.GenerateTimestampedSign(params)
}

// VerifyAPISign 验证API签名
func VerifyAPISign(params map[string]string, sign string, secretKey string, maxAge int64) error {
	signer := signature.NewSigner(secretKey)
	return signer.VerifySign(params, sign, maxAge)
}

// tenantContextTenant上下文中间件
func TenantContext() gin.HandlerFunc {
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

// DecryptAES AES解密中间件
func DecryptAES() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否需要解密
		encryptedBody := c.GetHeader("X-Encrypted-Body")
		if encryptedBody == "" {
			c.Next()
			return
		}

		// 获取解密密钥
		key := c.GetHeader("X-Decryption-Key")
		if key == "" {
			response.Error(c, http.StatusBadRequest, "缺少解密密钥")
			c.Abort()
			return
		}

		// TODO: 实现AES解密逻辑
		// decrypted, err := crypto.NewAESCipher(key).Decrypt(encryptedBody)
		// if err != nil {
		//     response.Error(c, http.StatusBadRequest, "解密失败")
		//     c.Abort()
		//     return
		// }
		// 将解密后的内容设置到请求体
		// c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(decrypted)))

		c.Next()
	}
}
