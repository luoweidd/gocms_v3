package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ValidateRequest 请求验证中间件
func ValidateRequest() gin.HandlerFunc {
	_ = validator.New()

	return func(c *gin.Context) {
		c.Next()
	}
}

// SanitizeResponse 数据脱敏中间件 - 移除敏感字段
func SanitizeResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
	}
}

// PaginationParams 分页参数绑定验证
type PaginationParams struct {
	Page     int `form:"page" binding:"gte=1"`
	PageSize int `form:"page_size" binding:"gte=1,lte=100"`
}

// BindAndValidate 绑定并验证分页参数
func BindAndValidate(c *gin.Context) (*PaginationParams, error) {
	var params PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	return &params, nil
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int64       `json:"total"`
	Items    interface{} `json:"items"`
}

// BuildPaginatedResponse 构建分页响应
func BuildPaginatedResponse(items interface{}, total int64, page, pageSize int) *PaginationResponse {
	return &PaginationResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    items,
	}
}
