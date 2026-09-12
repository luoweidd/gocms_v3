package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ApiResponse 统一响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMsg 带消息的成功响应
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, ApiResponse{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, ApiResponse{
		Code:    httpCode,
		Message: msg,
	})
}

// PageData 分页数据结构
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
