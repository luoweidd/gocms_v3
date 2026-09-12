package middleware

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶算法实现
type TokenBucket struct {
	tokens     float64
	maxTokens  float64
	refillRate float64 // 每秒补充的令牌数
	lastRefill time.Time
}

// RateLimiter 基于IP的速率限制器
type RateLimiter struct {
	mu           sync.RWMutex
	requestCount map[string]*TokenBucket
	maxRequests  float64
	refillRate   float64
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(maxRPS float64) *RateLimiter {
	return &RateLimiter{
		requestCount: make(map[string]*TokenBucket),
		maxRequests:  maxRPS,
		refillRate:   maxRPS / 60.0, // 每秒补充的令牌数
	}
}

// allow 检查是否允许请求
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	bucket, exists := rl.requestCount[key]
	if !exists {
		bucket = &TokenBucket{
			tokens:     rl.maxRequests,
			maxTokens:  rl.maxRequests,
			refillRate: rl.refillRate,
			lastRefill: now,
		}
		rl.requestCount[key] = bucket
	}

	// 补充令牌
	elapsed := now.Sub(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > bucket.maxTokens {
		bucket.tokens = bucket.maxTokens
	}
	bucket.lastRefill = now

	// 检查是否有可用令牌
	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		return true
	}

	return false
}

// cleanup 清理过期记录
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, bucket := range rl.requestCount {
		if now.Sub(bucket.lastRefill) > 5*time.Minute {
			delete(rl.requestCount, key)
		}
	}
}

// DefaultRateLimitConfig 默认速率限制配置
type RateLimitConfig struct {
	RPS       float64 `yaml:"rps"`
	KeyPrefix string
}

var DefaultRateLimitConfig = RateLimitConfig{
	RPS:       60,
	KeyPrefix: "ratelimit",
}

// GlobalRateLimiter 全局速率限制器实例
var GlobalRateLimiter = NewRateLimiter(60)

func init() {
	// 定时清理过期记录
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			GlobalRateLimiter.cleanup()
		}
	}()
}

// RateLimit 速率限制中间件
func RateLimit(cfg *RateLimitConfig) gin.HandlerFunc {
	if cfg == nil {
		cfg = &DefaultRateLimitConfig
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		// 跳过管理接口
		path := c.Request.URL.Path
		if path == "/swagger/*any" || path == "/health" || path == "/api/v2/openapi.json" {
			c.Next()
			return
		}

		key := fmt.Sprintf("%s:%s", cfg.KeyPrefix, clientIP)

		if !GlobalRateLimiter.allow(key) {
			response.Error(c, 429, "请求过于频繁，请稍后重试")
			c.Abort()
			return
		}

		// 设置速率限制相关头
		c.Header("X-RateLimit-Limit", strconv.FormatFloat(cfg.RPS, 'f', 0, 64))
		c.Header("X-RateLimit-Remaining", "1")
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10))

		c.Next()
	}
}

// RateLimitByRole 基于角色的速率限制中间件
func RateLimitByRole(role string, maxRequests float64) gin.HandlerFunc {
	roleLimiter := NewRateLimiter(maxRequests)

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		key := fmt.Sprintf("role:%s:%s", role, clientIP)

		if !roleLimiter.allow(key) {
			response.Error(c, 429, "请求过于频繁，请稍后重试")
			c.Abort()
			return
		}

		c.Next()
	}
}
