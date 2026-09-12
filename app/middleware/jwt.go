package middleware

import (
	"fmt"
	"strings"
	"time"

	"gocms_v3/app/config"
	"gocms_v3/app/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// minLen 返回两个整数中的较小值
func minLen(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// JWTClaims JWT Claims结构
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username string, roles []string) (string, error) {
	accessTokenTTL := config.GlobalConfig.JWT.AccessTokenTTL
	if accessTokenTTL == "" {
		accessTokenTTL = "2h"
	}

	expiresAt, err := time.ParseDuration(accessTokenTTL)
	if err != nil {
		expiresAt, _ = time.ParseDuration(accessTokenTTL)
	}

	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
			Issuer:    config.GlobalConfig.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GlobalConfig.JWT.Secret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GlobalConfig.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			fmt.Printf("[JWT] 未提供认证令牌 - %s %s\n", c.Request.Method, c.Request.URL.Path)
			response.Error(c, 401, "未提供认证令牌")
			c.Abort()
			return
		}

		// 支持 "Bearer " 和 "bearer " 前缀（大小写不敏感）
		var tokenString string
		lowerHeader := strings.ToLower(authHeader)
		if strings.HasPrefix(lowerHeader, "bearer ") {
			tokenString = authHeader[7:]
		} else {
			fmt.Printf("[JWT] 无效的认证头格式 - %s %s (expected 'Bearer <token>', got: '%s')\n", c.Request.Method, c.Request.URL.Path, authHeader[:minLen(len(authHeader), 50)])
			response.Error(c, 401, "无效的认证头格式")
			c.Abort()
			return
		}

		// 去除 token 首尾空白字符
		tokenString = strings.TrimSpace(tokenString)
		if tokenString == "" {
			fmt.Printf("[JWT] Token 为空 - %s %s\n", c.Request.Method, c.Request.URL.Path)
			response.Error(c, 401, "Token 不能为空")
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			fmt.Printf("[JWT] 令牌验证失败 - %s %s: %v\n", c.Request.Method, c.Request.URL.Path, err)
			response.Error(c, 401, "无效的认证令牌")
			c.Abort()
			return
		}

		fmt.Printf("[JWT] 认证成功 - 用户ID: %d, 用户名: %s - %s %s\n", claims.UserID, claims.Username, c.Request.Method, c.Request.URL.Path)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// GetTokenFromHeader 从请求头获取Token（支持大小写不敏感的 Bearer 前缀）
func GetTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	lowerHeader := strings.ToLower(authHeader)
	if strings.HasPrefix(lowerHeader, "bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	return ""
}
