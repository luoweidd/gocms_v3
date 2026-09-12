package model

import "time"

// ============================================
// Dashboard 数据模型 - 增强版
// 提供全面的仪表盘数据统计和监控
// ============================================

// StatusCount 状态统计
type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// RoleCount 角色计数
type RoleCount struct {
	Role  string `json:"role"`
	Count int64  `json:"count"`
}

// CategoryDist 分类分布
type CategoryDist struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Count        int64  `json:"count"`
}

// TagDist 标签分布
type TagDist struct {
	TagID   uint   `json:"tag_id"`
	TagName string `json:"tag_name"`
	Count   int64  `json:"count"`
}

// TrendData 趋势数据
type TrendData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// RecentItem 最近项目
type RecentItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Author    string    `json:"author"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// TopContent 热门内容
type TopContent struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Views     int64     `json:"views"`      // 浏览量
	Likes     int64     `json:"likes"`      // 点赞数
	Comments  int64     `json:"comments"`   // 评论数
	Tags      []string  `json:"tags"`       // 标签
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// TopAuthor 热门作者
type TopAuthor struct {
	Author string `json:"author"`
	Count  int64  `json:"count"`
}

// TopComment 热评评论
type TopComment struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Likes     int64     `json:"likes"`      // 点赞数
	URL       string    `json:"url"`        // 关联内容URL
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// ========== 硬件资源监控 ==========

// ResourceUsage 硬件资源使用信息
type ResourceUsage struct {
	CPU        CPUInfo     `json:"cpu"`        // CPU信息
	Memory     MemoryInfo  `json:"memory"`     // 内存信息
	Disk       DiskInfo    `json:"disk"`       // 磁盘信息
	Network    NetworkInfo `json:"network"`    // 网络信息
	Uptime     int64       `json:"uptime"`     // 系统运行时间（秒）
	Goroutines int         `json:"goroutines"` // Goroutine数量
}

// CPUInfo CPU信息
type CPUInfo struct {
	Cores        int     `json:"cores"`         // CPU核心数
	UsagePercent float64 `json:"usage_percent"` // CPU使用率
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total        uint64  `json:"total"`         // 总内存（字节）
	Used         uint64  `json:"used"`          // 已用内存（字节）
	Free         uint64  `json:"free"`          // 空闲内存（字节）
	UsagePercent float64 `json:"usage_percent"` // 内存使用率
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Total        uint64  `json:"total"`         // 总容量（字节）
	Used         uint64  `json:"used"`          // 已用容量（字节）
	Free         uint64  `json:"free"`          // 空闲容量（字节）
	UsagePercent float64 `json:"usage_percent"` // 磁盘使用率
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	UPSpeed   float64 `json:"up_speed"`   // 上传速度（MB/s）
	DOWNSpeed float64 `json:"down_speed"` // 下载速度（MB/s）
	UpTotal   uint64  `json:"up_total"`   // 总上传量（字节）
	DownTotal uint64  `json:"down_total"` // 总下载量（字节）
}

// ========== DashboardStats 仪表盘综合统计数据（扩展版） ==========

type DashboardStats struct {
	// ========== 概览统计 ==========
	TotalUsers       int64 `json:"total_users"`        // 总用户数
	TodayNewUsers    int64 `json:"today_new_users"`    // 今日新增用户
	TotalArticles    int64 `json:"total_articles"`     // 总文章数
	TodayNewArticles int64 `json:"today_new_articles"` // 今日新增文章
	TotalVideos      int64 `json:"total_videos"`       // 总视频数
	TodayNewVideos   int64 `json:"today_new_videos"`   // 今日新增视频
	TotalComments    int64 `json:"total_comments"`     // 总评论数
	TodayNewComments int64 `json:"today_new_comments"` // 今日新增评论

	// ========== 文件统计 ==========
	TotalFiles     int64 `json:"total_files"`      // 总文件数
	TodayNewFiles  int64 `json:"today_new_files"`  // 今日新增文件
	TotalFileSize  int64 `json:"total_file_size"`  // 总文件大小（字节）
	ImageCount     int64 `json:"image_count"`      // 图片数量
	VideoFileCount int64 `json:"video_file_count"` // 视频文件数量
	AudioFileCount int64 `json:"audio_file_count"` // 音频文件数量
	DocumentCount  int64 `json:"document_count"`   // 文档数量

	// ========== 内容状态分布 ==========
	ArticleStatus []StatusCount `json:"article_status"` // 文章状态分布
	VideoStatus   []StatusCount `json:"video_status"`   // 视频状态分布
	CommentStatus []StatusCount `json:"comment_status"` // 评论状态分布

	// ========== 用户状态分布 ==========
	UserStatus   []StatusCount `json:"user_status"`    // 用户状态分布
	UserRoleDist []RoleCount   `json:"user_role_dist"` // 用户角色分布

	// ========== 分类分布 ==========
	ArticleCategories []CategoryDist `json:"article_categories"` // 文章分类分布
	VideoCategories   []CategoryDist `json:"video_categories"`   // 视频分类分布
	TagDistribution   []TagDist      `json:"tag_distribution"`   // 标签分布

	// ========== 趋势数据 ==========
	ArticleTrend []TrendData `json:"article_trend"` // 文章趋势（30天）
	VideoTrend   []TrendData `json:"video_trend"`   // 视频趋势（30天）
	UserTrend    []TrendData `json:"user_trend"`    // 用户增长趋势（30天）
	CommentTrend []TrendData `json:"comment_trend"` // 评论趋势（30天）

	// ========== 最近活动 ==========
	RecentArticles []RecentItem `json:"recent_articles"` // 最近文章
	RecentVideos   []RecentItem `json:"recent_videos"`   // 最近视频
	RecentComments []RecentItem `json:"recent_comments"` // 最近评论
	RecentUsers    []RecentItem `json:"recent_users"`    // 最近用户

	// ========== 热门内容 ==========
	TopArticles   []TopContent `json:"top_articles"`   // 热门文章
	TopVideos     []TopContent `json:"top_videos"`     // 热门视频
	TopComments   []TopComment `json:"top_comments"`   // 热评评论
	ActiveAuthors []TopAuthor  `json:"active_authors"` // 活跃作者

	// ========== 硬件资源监控 ==========
	ResourceUsage *ResourceUsage `json:"resource_usage"` // 资源使用情况

	// ========== 内容统计详情 ==========
	TotalWords         int64 `json:"total_words"`          // 总字数
	AvgReadTime        int   `json:"avg_read_time"`        // 平均阅读时间（分钟）
	TotalVideoDuration int64 `json:"total_video_duration"` // 视频总时长（秒）
	AvgVideoDuration   int   `json:"avg_video_duration"`   // 平均视频时长（秒）
}

// ========== DashboardOverview 仪表盘概览数据 ==========

type DashboardOverview struct {
	TotalUsers       int64 `json:"total_users"`        // 总用户数
	TotalArticles    int64 `json:"total_articles"`     // 总文章数
	TotalVideos      int64 `json:"total_videos"`       // 总视频数
	TotalComments    int64 `json:"total_comments"`     // 总评论数
	TodayNewUsers    int64 `json:"today_new_users"`    // 今日新增用户
	TodayNewArticles int64 `json:"today_new_articles"` // 今日新增文章
	TodayNewVideos   int64 `json:"today_new_videos"`   // 今日新增视频
	TodayNewComments int64 `json:"today_new_comments"` // 今日新增评论

	// 增长率
	UserGrowthRate    float64 `json:"user_growth_rate"`    // 用户增长率(%)
	ArticleGrowthRate float64 `json:"article_growth_rate"` // 文章增长率(%)
	VideoGrowthRate   float64 `json:"video_growth_rate"`   // 视频增长率(%)
	CommentGrowthRate float64 `json:"comment_growth_rate"` // 评论增长率(%)

	// 资源概览
	CPUUsage    float64 `json:"cpu_usage"`    // CPU使用率(%)
	MemoryUsage float64 `json:"memory_usage"` // 内存使用率(%)
	DiskUsage   float64 `json:"disk_usage"`   // 磁盘使用率(%)
	QPS         float64 `json:"qps"`          // 每秒查询数
}

// ========== ContentStats 内容统计详情 ==========

type ContentStats struct {
	Articles *ArticleStats `json:"articles"` // 文章统计
	Videos   *VideoStats   `json:"videos"`   // 视频统计
	Comments *CommentStats `json:"comments"` // 评论统计
}

// ArticleStats 文章统计
type ArticleStats struct {
	Total         int64          `json:"total"`          // 总数
	Published     int64          `json:"published"`      // 已发布
	Draft         int64          `json:"draft"`          // 草稿
	PendingReview int64          `json:"pending_review"` // 待审核
	TotalViews    int64          `json:"total_views"`    // 总浏览量
	TotalLikes    int64          `json:"total_likes"`    // 总点赞数
	TotalComments int64          `json:"total_comments"` // 总评论数
	TodayCreate   int64          `json:"today_created"`  // 今日创建
	MonthCreate   int64          `json:"month_created"`  // 本月创建
	Categories    []CategoryDist `json:"categories"`     // 分类统计
	TotalWords    int64          `json:"total_words"`    // 总字数
	AvgReadTime   int            `json:"avg_read_time"`  // 平均阅读时间（分钟）
}

// VideoStats 视频统计
type VideoStats struct {
	Total         int64          `json:"total"`          // 总数
	Published     int64          `json:"published"`      // 已发布
	Draft         int64          `json:"draft"`          // 草稿
	Processing    int64          `json:"processing"`     // 处理中
	TotalViews    int64          `json:"total_views"`    // 总播放量
	TotalLikes    int64          `json:"total_likes"`    // 总点赞数
	TotalShares   int64          `json:"total_shares"`   // 总分享数
	TodayCreate   int64          `json:"today_created"`  // 今日创建
	MonthCreate   int64          `json:"month_created"`  // 本月创建
	TotalDuration int64          `json:"total_duration"` // 总时长（秒）
	AvgDuration   int            `json:"avg_duration"`   // 平均时长（秒）
	Categories    []CategoryDist `json:"categories"`     // 分类统计
}

// CommentStats 评论统计
type CommentStats struct {
	Total             int64       `json:"total"`               // 总数
	Approved          int64       `json:"approved"`            // 已审核
	Pending           int64       `json:"pending"`             // 待审核
	Rejected          int64       `json:"rejected"`            // 已拒绝
	TodayCount        int64       `json:"today_count"`         // 今日评论
	ArticleCommentCnt int64       `json:"article_comment_cnt"` // 文章评论数
	VideoCommentCnt   int64       `json:"video_comment_cnt"`   // 视频评论数
	TopAuthors        []TopAuthor `json:"top_authors"`         // 热门评论者
}

// ========== UserActivityStats 用户活跃度统计 ==========

type UserActivityStats struct {
	ActiveUsers      map[string]int64 `json:"active_users"`       // 活跃用户（按天）
	DailyActive      int64            `json:"daily_active"`       // 日活跃用户
	WeeklyActive     int64            `json:"weekly_active"`      // 周活跃用户
	MonthlyActive    int64            `json:"monthly_active"`     // 月活跃用户
	NewUserRetention map[string]int64 `json:"new_user_retention"` // 新用户留存率
	UserLevels       map[string]int64 `json:"user_levels"`        // 用户等级分布
}

// ========== PerformanceStats 性能统计 ==========

type PerformanceStats struct {
	APILatency  []float64 `json:"api_latency"`  // API延迟（ms）
	DBQueries   int64     `json:"db_queries"`   // 数据库查询数
	CacheHit    float64   `json:"cache_hit"`    // 缓存命中率(%)
	ServerError int64     `json:"server_error"` // 服务器错误数
}

// ========== SystemInfo 系统信息 ==========

type SystemInfo struct {
	GOVersion  string `json:"go_version"`  // Go版本
	OS         string `json:"os"`          // 操作系统
	Arch       string `json:"arch"`        // 架构
	ServerIP   string `json:"server_ip"`   // 服务器IP
	Hostname   string `json:"hostname"`    // 主机名
	StartTime  string `json:"start_time"`  // 启动时间
	Uptime     int64  `json:"uptime"`      // 运行时间（秒）
	LastBackup string `json:"last_backup"` // 最后备份时间
}
