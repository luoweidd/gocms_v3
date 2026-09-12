package errors

import (
	"errors"
	"fmt"
)

// ErrorCode 错误码类型
type ErrorCode string

const (
	// 通用错误码
	ErrCodeSuccess         ErrorCode = "SUCCESS"
	ErrCodeUnknown         ErrorCode = "UNKNOWN"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeInternal        ErrorCode = "INTERNAL"
	ErrCodeValidationError ErrorCode = "VALIDATION_ERROR"

	// 用户相关错误码
	ErrCodeUserNotFound       ErrorCode = "USER_NOT_FOUND"
	ErrCodeUserAlreadyExists  ErrorCode = "USER_ALREADY_EXISTS"
	ErrCodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"

	// 文章相关错误码
	ErrCodeArticleNotFound ErrorCode = "ARTICLE_NOT_FOUND"

	// 视频相关错误码
	ErrCodeVideoNotFound ErrorCode = "VIDEO_NOT_FOUND"

	// 分类相关错误码
	ErrCodeCategoryNotFound ErrorCode = "CATEGORY_NOT_FOUND"

	// 权限相关错误码
	ErrCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrCodePermissionNotFound ErrorCode = "PERMISSION_NOT_FOUND"

	// 业务错误码
	ErrCodeBusinessError ErrorCode = "BUSINESS_ERROR"
)

// AppError 应用错误类型
type AppError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Err     error       `json:"-"`
	Data    interface{} `json:"data,omitempty"`
}

// New 创建新的应用错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewWithErr 创建带底层错误的应用错误
func NewWithErr(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 实现 errors.Is 和 errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is 实现 errors.Is
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Wrap 包装错误
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

// WithData 附加额外数据
func (e *AppError) WithData(data interface{}) *AppError {
	e.Data = data
	return e
}

// IsAppError 检查是否为应用错误
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// IsNotFound 检查是否为未找到错误
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeNotFound
	}
	return false
}

// IsUnauthorized 检查是否为未授权错误
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeUnauthorized
	}
	return false
}

// IsPermissionDenied 检查是否为权限拒绝错误
func IsPermissionDenied(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodePermissionDenied
	}
	return false
}

// IsConflict 检查是否为冲突错误
func IsConflict(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == ErrCodeConflict
	}
	return false
}

// Wrap 包装底层错误为应用错误
func Wrap(err error, code ErrorCode, message string) error {
	if err == nil {
		return nil
	}
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Must 辅助函数，用于处理错误
func Must[T any](v T, err error) T {
	if err != nil {
		panic(&AppError{
			Code:    ErrCodeInternal,
			Message: err.Error(),
			Err:     err,
		})
	}
	return v
}
