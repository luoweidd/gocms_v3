package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLogConfig 操作日志中间件配置
type OperationLogConfig struct {
	// Enabled 是否启用（默认true）
	Enabled bool
	// ExcludePaths 不需要记录的路径前缀列表
	ExcludePrefix []string
	// ModuleOverride 按路径前缀覆盖模块名
	ModuleOverrides map[string]string
}

// DefaultOperationLogConfig 默认配置
func DefaultOperationLogConfig() OperationLogConfig {
	return OperationLogConfig{
		Enabled: true,
		ExcludePrefix: []string{
			"/api/dashboard/resource-usage", // 忽略资源监控请求（高频）
			"/api/ws/",                      // 忽略WebSocket
			"/favicon.ico",
			"/static/",
		},
		ModuleOverrides: map[string]string{
			"/api/articles/":       "article",
			"/api/video/":          "video",
			"/api/videos/":         "video",
			"/api/comments/":       "comment",
			"/api/users/":          "user",
			"/api/user/":           "user",
			"/api/roles/":          "role",
			"/api/menus/":          "menu",
			"/api/upload/":         "upload",
			"/api/ngac/":           "ngac",
			"/api/graphql":         "graphql",
			"/api/system/":         "system",
			"/api/data-dicts/":     "data_dict",
			"/api/system-configs/": "system_config",
		},
	}
}

// OperationLog 操作日志中间件
func OperationLog(config ...OperationLogConfig) gin.HandlerFunc {
	cfg := DefaultOperationLogConfig()
	if len(config) > 0 {
		if !config[0].Enabled {
			// 如果禁用，直接通过
			return func(c *gin.Context) {
				c.Next()
			}
		}
		if config[0].ExcludePrefix != nil {
			cfg.ExcludePrefix = config[0].ExcludePrefix
		}
		if config[0].ModuleOverrides != nil {
			cfg.ModuleOverrides = config[0].ModuleOverrides
		}
	}

	return func(c *gin.Context) {
		// 1. 检查是否需要排除
		path := c.Request.URL.Path
		for _, prefix := range cfg.ExcludePrefix {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

		// 2. 记录请求开始时间
		start := time.Now()

		// 3. 捕获响应体
		bodyBuf := &bytes.Buffer{}
		c.Writer = &responseWriter{
			ResponseWriter: c.Writer,
			writer:         bodyBuf,
		}

		// 4. 读取请求body（非文件上传）
		var requestBody string
		if c.Request.ContentLength > 0 && !isMultipartFormData(c) {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			requestBody = sanitizeJSON(bodyBytes)
		}

		// 5. 执行请求
		c.Next()

		// 6. 记录操作日志
		duration := time.Since(start).Milliseconds()
		statusCode := c.Writer.Status()
		responseBody := bodyBuf.String()

		// 确定操作类型
		action := determineAction(c)

		// 确定模块名称
		module := determineModule(path, cfg.ModuleOverrides)

		// 获取用户信息
		userID := uint64(0)
		username := ""
		if claims, exists := c.Get("claims"); exists && claims != nil {
			if claimMap, ok := claims.(map[string]interface{}); ok {
				if uid, ok := claimMap["user_id"].(float64); ok {
					userID = uint64(uid)
				}
				if un, ok := claimMap["username"].(string); ok {
					username = un
				}
			}
		}

		// 获取IP地址
		ipAddress := c.ClientIP()
		if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
			ipAddress = realIP
		} else if forwardFor := c.GetHeader("X-Forwarded-For"); forwardFor != "" {
			ipAddress = splitIPs(forwardFor)[0]
		}

		// 构造操作描述
		description := fmt.Sprintf("%s %s", action, module)
		if id := c.Param("id"); id != "" {
			description += fmt.Sprintf("#%s", id)
		}

		// 判断成功/失败
		logStatus := int8(1)
		errorMsg := ""
		if statusCode >= 400 {
			logStatus = 0
			// 尝试从响应中提取错误信息
			var resp map[string]interface{}
			if err := json.Unmarshal([]byte(responseBody), &resp); err == nil {
				if msg, ok := resp["message"].(string); ok {
					errorMsg = truncate(msg, 500)
				}
			}
			if errorMsg == "" {
				errorMsg = fmt.Sprintf("HTTP %d", statusCode)
			}
		}

		// 截断过长的字段
		requestBody = truncate(requestBody, 65535)
		responseBody = truncate(responseBody, 65535)

		// 保存到数据库（异步，不阻塞响应）
		go func() {
			logEntry := model.OperationLog{
				UserID:        userID,
				Username:      truncate(username, 100),
				Action:        action,
				Module:        module,
				Description:   truncate(description, 500),
				IPAddress:     ipAddress,
				UserAgent:     truncate(c.GetHeader("User-Agent"), 500),
				RequestMethod: c.Request.Method,
				RequestPath:   truncate(path, 255),
				RequestData:   requestBody,
				ResponseBody:  responseBody,
				DurationMs:    int(duration),
				Status:        logStatus,
				ErrorMessage:  errorMsg,
			}

			if err := db.GetDB().Create(&logEntry).Error; err != nil {
				zap.L().Error("保存操作日志失败",
					zap.Error(err),
					zap.Uint64("user_id", userID),
					zap.String("action", action),
					zap.String("module", module),
				)
			}
		}()

		// 7. 可选：记录快速日志到控制台（DEBUG级别）
		if zap.L().Core().Enabled(zap.DebugLevel) {
			zap.L().Debug("操作日志",
				zap.Uint64("user_id", userID),
				zap.String("username", username),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Int("status", statusCode),
				zap.Int64("duration_ms", duration),
				zap.String("action", action),
				zap.String("module", module),
			)
		}
	}
}

// responseWriter 包装gin.ResponseWriter以捕获响应体
type responseWriter struct {
	gin.ResponseWriter
	writer *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.writer.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteString(s string) (int, error) {
	rw.writer.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}

// determineAction 根据方法和路径判断操作类型
func determineAction(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path
	action := "list" // 默认

	// GET /api/articles -> list
	if method == "GET" {
		if strings.Contains(path, "/stats") || strings.Contains(path, "/count") {
			return "stats"
		}
		if strings.HasSuffix(path, "/tree") {
			return "tree"
		}
		if id := c.Param("id"); id != "" {
			action = "get" // 查看详情
		} else {
			action = "list"
		}
		return action
	}

	// POST -> create / publish 等
	if method == "POST" {
		if strings.Contains(path, "/publish") {
			return "publish"
		}
		if strings.Contains(path, "/top") && !strings.Contains(path, "/untop") {
			return "top"
		}
		if strings.Contains(path, "/approve") || strings.Contains(path, "/audit") {
			return "approve"
		}
		if strings.Contains(path, "/batch") {
			return "batch_create"
		}
		if c.Param("action") != "" {
			return c.Param("action")
		}
		// 检查body中的action字段
		var body map[string]interface{}
		if err := json.Unmarshal(peekBody(c), &body); err == nil {
			if a, ok := body["action"].(string); ok {
				return a
			}
		}
		return "create"
	}

	// PUT/PATCH -> update
	if method == "PUT" || method == "PATCH" {
		if strings.Contains(path, "/unpublish") {
			return "unpublish"
		}
		if strings.Contains(path, "/untop") {
			return "untop"
		}
		if strings.Contains(path, "/batch") {
			return "batch_update"
		}
		return "update"
	}

	// DELETE -> delete
	if method == "DELETE" {
		if strings.Contains(path, "/batch") {
			return "batch_delete"
		}
		return "delete"
	}

	return action
}

// peekBody 查看请求body（不消费）
func peekBody(c *gin.Context) []byte {
	bodyBytes, _ := c.GetRawData()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

// determineModule 根据路径确定模块名
func determineModule(path string, overrides map[string]string) string {
	// 检查覆盖映射（按最长前缀匹配）
	var longestPrefix int
	module := "unknown"
	for prefix, mod := range overrides {
		if strings.HasPrefix(path, prefix) && len(prefix) > longestPrefix {
			longestPrefix = len(prefix)
			module = mod
		}
	}
	return module
}

// sanitizeJSON 清理JSON字符串（防止注入）
func sanitizeJSON(data []byte) string {
	// 限制最大长度
	maxLen := 65535
	if len(data) > maxLen {
		data = data[:maxLen]
	}
	// 移除可能的恶意内容（null byte）
	clean := bytes.ReplaceAll(data, []byte{0x00}, []byte{})
	return string(clean)
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// splitIPs 分割IP列表（取第一个）
func splitIPs(s string) []string {
	result := strings.Split(s, ",")
	ips := make([]string, 0, len(result))
	for _, ip := range result {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips
}

// isMultipartFormData 检查是否为multipart表单
func isMultipartFormData(c *gin.Context) bool {
	contentType := c.Request.Header.Get("Content-Type")
	return strings.Contains(contentType, "multipart/form-data")
}
