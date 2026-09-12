package i18n

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// Translations 多语言翻译映射
var Translations = make(map[string]map[string]string)

var initOnce sync.Once

func init() {
	initOnce.Do(func() {
		Translations["zh"] = map[string]string{
			"article.created": "文章创建成功",
			"article.updated": "文章更新成功",
			"article.deleted": "文章删除成功",
			"video.created":   "视频创建成功",
			"video.updated":   "视频更新成功",
			"video.deleted":   "视频删除成功",
			"comment.audited": "评论审核成功",
			"user.created":    "用户创建成功",
			"user.updated":    "用户更新成功",
			"user.deleted":    "用户删除成功",
			"menu.tree":       "菜单树结构",
			"export.success":  "导出成功",
			"upload.success":  "上传成功",
			"upload.chunk":    "分片上传成功",
			"upload.complete": "合并完成",
			"auth.login":      "登录成功",
			"auth.register":   "注册成功",
			"auth.logout":     "退出成功",
		}

		Translations["en"] = map[string]string{
			"article.created": "Article created successfully",
			"article.updated": "Article updated successfully",
			"article.deleted": "Article deleted successfully",
			"video.created":   "Video created successfully",
			"video.updated":   "Video updated successfully",
			"video.deleted":   "Video deleted successfully",
			"comment.audited": "Comment audited successfully",
			"user.created":    "User created successfully",
			"user.updated":    "User updated successfully",
			"user.deleted":    "User deleted successfully",
			"menu.tree":       "Menu tree structure",
			"export.success":  "Export successful",
			"upload.success":  "Upload successful",
			"upload.chunk":    "Chunk upload successful",
			"upload.complete": "Merge complete",
			"auth.login":      "Login successful",
			"auth.register":   "Register successful",
			"auth.logout":     "Logout successful",
		}
	})
}

// T 翻译函数
func T(locale string, key string) string {
	if msgs, ok := Translations[locale]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	if enMsgs, ok := Translations["en"]; ok {
		if msg, ok := enMsgs[key]; ok {
			return msg
		}
	}
	return key
}

// GetLocale 从请求中获取语言环境
func GetLocale(c *gin.Context) string {
	if locale := c.GetHeader("X-Locale"); locale != "" {
		return locale
	}
	if locale := c.Query("locale"); locale != "" {
		return locale
	}
	return "zh"
}
