package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenRefreshRequest Token刷新请求
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// TokenRefreshResponse Token刷新响应
type TokenRefreshResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, username string) (string, string, error) {
	accessToken, err := GenerateToken(userID, username, nil)
	if err != nil {
		return "", "", fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken := fmt.Sprintf("refresh_%d_%s_%d", userID, username, time.Now().Unix())
	_ = time.Now().Add(7 * 24 * time.Hour) // 7天有效

	return accessToken, refreshToken, nil
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(refreshToken string) (uint, error) {
	if len(refreshToken) < 10 || refreshToken[:8] != "refresh_" {
		return 0, fmt.Errorf("无效的刷新令牌格式")
	}
	// 简化实现：实际项目中应查询数据库验证
	return 0, nil
}

// RevokeRefreshToken 撤销刷新令牌
func RevokeRefreshToken(refreshToken string) error {
	// 简化实现：实际项目中应将令牌加入黑名单
	if refreshToken == "" {
		return fmt.Errorf("令牌不能为空")
	}
	return nil
}

// RefreshTokenMiddleware Token刷新中间件
func RefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否需要刷新令牌
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.Next()
			return
		}

		// TODO: 解析令牌并检查是否即将过期
		// 如果接近过期，可以在响应头中添加X-Token-Refresh提示

		c.Next()
	}
}

// RevokeAllUserTokens 撤销用户所有刷新令牌
func RevokeAllUserTokens(userID uint) error {
	// 简化实现：实际项目中应查询数据库并撤销所有相关令牌
	if userID == 0 {
		return fmt.Errorf("用户ID不能为零")
	}
	return nil
}

// IsTokenValid 检查令牌是否有效
func IsTokenValid(refreshToken string) bool {
	return len(refreshToken) >= 10 && refreshToken[:8] == "refresh_"
}
