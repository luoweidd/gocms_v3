package dto

import "time"

// ==================== 用户相关 DTO ====================

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	Token    string   `json:"token"`
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Roles    []string `json:"roles"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required,min=6"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Roles    []string `json:"roles"`
}

// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required,min=6"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Roles    []string `json:"roles"`
	RoleID   *uint    `json:"role_id"` // 支持单个角色ID（前端发送）
	Status   int      `json:"status"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   *int   `json:"status,omitempty"`
}

// UserUpdatePasswordRequest 用户更新密码请求
type UserUpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	Avatar      string     `json:"avatar"`
	Status      int        `json:"status"`
	Roles       []string   `json:"roles"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	List     []UserInfoResponse `json:"list"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

// ==================== 文章相关 DTO ====================

// ArticleCreateRequest 文章创建请求
type ArticleCreateRequest struct {
	Title       string     `json:"title" binding:"required"`
	Content     string     `json:"content" binding:"required"`
	Abstract    string     `json:"abstract"`
	CoverImage  string     `json:"cover_image"`
	CategoryID  uint       `json:"category_id" binding:"required"`
	TagIDs      []uint     `json:"tag_ids"`
	Status      int        `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

// ArticleUpdateRequest 文章更新请求
type ArticleUpdateRequest struct {
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Abstract    string     `json:"abstract"`
	CoverImage  string     `json:"cover_image"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	TagIDs      []uint     `json:"tag_ids"`
	Status      *int       `json:"status,omitempty"`
	PublishedAt *time.Time `json:"published_at"`
}

// ArticleListItem 文章列表项响应
type ArticleListItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Abstract     string    `json:"abstract"`
	CoverImage   string    `json:"cover_image"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	Status       int       `json:"status"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ArticleDetailResponse 文章详情响应
type ArticleDetailResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	Abstract     string    `json:"abstract"`
	CoverImage   string    `json:"cover_image"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	TagIDs       []uint    `json:"tag_ids"`
	Status       int       `json:"status"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ArticleListResponse 文章列表响应
type ArticleListResponse struct {
	List     []ArticleListItem `json:"list"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// ==================== 视频相关 DTO ====================

// VideoCreateRequest 视频创建请求
type VideoCreateRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	CoverImage  string     `json:"cover_image"`
	VideoURL    string     `json:"video_url" binding:"required"`
	CategoryID  uint       `json:"category_id" binding:"required"`
	TagIDs      []uint     `json:"tag_ids"`
	Duration    int        `json:"duration"`
	Status      int        `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
}

// VideoUpdateRequest 视频更新请求
type VideoUpdateRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	CoverImage  string     `json:"cover_image"`
	VideoURL    string     `json:"video_url"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	TagIDs      []uint     `json:"tag_ids"`
	Duration    *int       `json:"duration,omitempty"`
	Status      *int       `json:"status,omitempty"`
	PublishedAt *time.Time `json:"published_at"`
}

// VideoListItem 视频列表项响应
type VideoListItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverImage   string    `json:"cover_image"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	Status       int       `json:"status"`
	Duration     int       `json:"duration"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
}

// VideoDetailResponse 视频详情响应
type VideoDetailResponse struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CoverImage   string    `json:"cover_image"`
	VideoURL     string    `json:"video_url"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name,omitempty"`
	TagIDs       []uint    `json:"tag_ids"`
	Status       int       `json:"status"`
	Duration     int       `json:"duration"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ==================== 分类相关 DTO ====================

// CategoryCreateRequest 分类创建请求
type CategoryCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	ParentID  uint   `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	Status    int    `json:"status"`
}

// CategoryUpdateRequest 分类更新请求
type CategoryUpdateRequest struct {
	Name      string `json:"name"`
	ParentID  *uint  `json:"parent_id,omitempty"`
	SortOrder *int   `json:"sort_order,omitempty"`
	Status    *int   `json:"status,omitempty"`
}

// CategoryItem 分类项响应
type CategoryItem struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	ParentID  uint           `json:"parent_id"`
	SortOrder int            `json:"sort_order"`
	Status    int            `json:"status"`
	Children  []CategoryItem `json:"children,omitempty"`
}

// ==================== 评论相关 DTO ====================

// CommentCreateRequest 评论创建请求
type CommentCreateRequest struct {
	ArticleID uint   `json:"article_id"`
	VideoID   uint   `json:"video_id"`
	Content   string `json:"content" binding:"required"`
	ParentID  *uint  `json:"parent_id"`
}

// CommentUpdateRequest 评论更新请求
type CommentUpdateRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentItem 评论项响应
type CommentItem struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	ArticleID uint      `json:"article_id"`
	VideoID   uint      `json:"video_id"`
	ParentID  *uint     `json:"parent_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ==================== 菜单相关 DTO ====================

// MenuCreateRequest 菜单创建请求
type MenuCreateRequest struct {
	Name      string `json:"name" binding:"required"`
	Path      string `json:"path" binding:"required"`
	Icon      string `json:"icon"`
	ParentID  uint   `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	Status    int    `json:"status"`
}

// MenuUpdateRequest 菜单更新请求
type MenuUpdateRequest struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Icon      string `json:"icon"`
	ParentID  *uint  `json:"parent_id,omitempty"`
	SortOrder *int   `json:"sort_order,omitempty"`
	Status    *int   `json:"status,omitempty"`
}

// MenuItem 菜单项响应
type MenuItem struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Icon      string     `json:"icon"`
	ParentID  uint       `json:"parent_id"`
	SortOrder int        `json:"sort_order"`
	Status    int        `json:"status"`
	Children  []MenuItem `json:"children,omitempty"`
}

// ==================== 文件上传相关 DTO ====================

// FileUploadInitRequest 文件上传初始化请求
type FileUploadInitRequest struct {
	Filename  string `json:"filename" binding:"required"`
	Size      int64  `json:"size" binding:"required"`
	MimeType  string `json:"mime_type"`
	ChunkSize int    `json:"chunk_size"`
}

// FileUploadChunkRequest 文件分片上传请求
type FileUploadChunkRequest struct {
	FileID      string `form:"file_id" binding:"required"`
	ChunkNum    int    `form:"chunk_num" binding:"required"`
	TotalChunks int    `form:"total_chunks" binding:"required"`
	File        any    `form:"file"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID   string `json:"file_id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
}

// ==================== NGAC 相关 DTO ====================

// RoleCreateRequest 角色创建请求
type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
}

// RoleUpdateRequest 角色更新请求
type RoleUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

// PermissionCreateRequest 权限创建请求
type PermissionCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

// PermissionUpdateRequest 权限更新请求
type PermissionUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// AuthzRequest 授权评估请求
type AuthzRequest struct {
	SubjectID  uint           `json:"subject_id" binding:"required"`
	ObjectType string         `json:"object_type" binding:"required"`
	ObjectID   uint           `json:"object_id" binding:"required"`
	Action     string         `json:"action" binding:"required"`
	Context    map[string]any `json:"context"`
}

// AuthzResponse 授权评估响应
type AuthzResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason,omitempty"`
}
