package service

import (
	"fmt"
	"log"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

// DashboardService 仪表盘服务
type DashboardService struct{}

// NewDashboardService 创建仪表盘服务
func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

// GetStats 获取仪表盘综合统计数据
func (s *DashboardService) GetStats() (*model.DashboardStats, error) {
	stats := &model.DashboardStats{
		ArticleStatus:     []model.StatusCount{},
		VideoStatus:       []model.StatusCount{},
		CommentStatus:     []model.StatusCount{},
		ArticleTrend:      []model.TrendData{},
		VideoTrend:        []model.TrendData{},
		UserTrend:         []model.TrendData{},
		ArticleCategories: []model.CategoryDist{},
		VideoCategories:   []model.CategoryDist{},
	}

	// 获取文件统计
	if err := s.getFileStats(stats); err != nil {
		log.Printf("[SERVICE-WARN] 获取文件统计失败: %v", err)
		// 文件统计失败不影响其他统计
	}

	// 获取概览统计数据（包含用户总数等）
	overview, err := s.GetOverview()
	if err == nil && overview != nil {
		stats.TotalUsers = overview.TotalUsers
		stats.TodayNewUsers = overview.TodayNewUsers
		stats.TotalArticles = overview.TotalArticles
		stats.TodayNewArticles = overview.TodayNewArticles
		stats.TotalVideos = overview.TotalVideos
		stats.TodayNewVideos = overview.TodayNewVideos
		stats.TotalComments = overview.TotalComments
		stats.TodayNewComments = overview.TodayNewComments
	}

	// 获取硬件资源使用情况
	resourceUsage, err2 := s.GetResourceUsage()
	if err2 == nil && resourceUsage != nil {
		stats.ResourceUsage = resourceUsage
	}

	// 文章状态分布
	var articleStatusCounts []struct {
		Status int
		Count  int64
	}
	if err := db.GetDB().Model(&model.Article{}).Select("status, COUNT(*) as count").Group("status").Find(&articleStatusCounts).Error; err != nil {
		return nil, fmt.Errorf("统计文章状态分布失败：%w", err)
	}
	for _, asc := range articleStatusCounts {
		statusName := "未知"
		switch asc.Status {
		case 0:
			statusName = "草稿"
		case 1:
			statusName = "已发布"
		case 2:
			statusName = "已下架"
		case 3:
			statusName = "已删除"
		}
		stats.ArticleStatus = append(stats.ArticleStatus, model.StatusCount{
			Status: statusName,
			Count:  asc.Count,
		})
	}

	// 视频状态分布
	var videoStatusCounts []struct {
		Status int
		Count  int64
	}
	if err := db.GetDB().Model(&model.Video{}).Select("status, COUNT(*) as count").Group("status").Find(&videoStatusCounts).Error; err != nil {
		return nil, fmt.Errorf("统计视频状态分布失败：%w", err)
	}
	for _, vsc := range videoStatusCounts {
		statusName := "未知"
		switch vsc.Status {
		case 0:
			statusName = "草稿"
		case 1:
			statusName = "已发布"
		case 2:
			statusName = "已下架"
		case 3:
			statusName = "已删除"
		}
		stats.VideoStatus = append(stats.VideoStatus, model.StatusCount{
			Status: statusName,
			Count:  vsc.Count,
		})
	}

	// 评论状态分布
	var commentStatusCounts []struct {
		Status int
		Count  int64
	}
	if err := db.GetDB().Model(&model.Comment{}).Select("status, COUNT(*) as count").Group("status").Find(&commentStatusCounts).Error; err != nil {
		return nil, fmt.Errorf("统计评论状态分布失败：%w", err)
	}
	for _, csc := range commentStatusCounts {
		statusName := "未知"
		switch csc.Status {
		case 0:
			statusName = "待审核"
		case 1:
			statusName = "已通过"
		case 2:
			statusName = "已拒绝"
		case 3:
			statusName = "已删除"
		}
		stats.CommentStatus = append(stats.CommentStatus, model.StatusCount{
			Status: statusName,
			Count:  csc.Count,
		})
	}

	return stats, nil
}

// GetOverview 获取仪表盘概览数据（简化版）
func (s *DashboardService) GetOverview() (*model.DashboardOverview, error) {
	overview := &model.DashboardOverview{}

	today := time.Now().Truncate(24 * time.Hour)

	// 总用户数
	if err := db.GetDB().Model(&model.User{}).Count(&overview.TotalUsers).Error; err != nil {
		return nil, fmt.Errorf("统计用户总数失败：%w", err)
	}

	// 今日新增用户
	db.GetDB().Model(&model.User{}).Where("created_at >= ?", today).Count(&overview.TodayNewUsers)

	// 总文章数
	if err := db.GetDB().Model(&model.Article{}).Count(&overview.TotalArticles).Error; err != nil {
		return nil, fmt.Errorf("统计文章总数失败：%w", err)
	}

	// 今日新增文章
	db.GetDB().Model(&model.Article{}).Where("created_at >= ?", today).Count(&overview.TodayNewArticles)

	// 总视频数
	if err := db.GetDB().Model(&model.Video{}).Count(&overview.TotalVideos).Error; err != nil {
		return nil, fmt.Errorf("统计视频总数失败：%w", err)
	}

	// 今日新增视频
	db.GetDB().Model(&model.Video{}).Where("created_at >= ?", today).Count(&overview.TodayNewVideos)

	// 总评论数
	if err := db.GetDB().Model(&model.Comment{}).Count(&overview.TotalComments).Error; err != nil {
		return nil, fmt.Errorf("统计评论总数失败：%w", err)
	}

	// 今日新增评论
	db.GetDB().Model(&model.Comment{}).Where("created_at >= ?", today).Count(&overview.TodayNewComments)

	// 资源概览
	if memInfo, err := mem.VirtualMemory(); err == nil {
		overview.MemoryUsage = memInfo.UsedPercent
	}
	if diskInfo, err := disk.Usage("/"); err == nil {
		overview.DiskUsage = diskInfo.UsedPercent
	}
	if cpuPercent, err := cpu.Percent(0, false); err == nil && len(cpuPercent) > 0 {
		overview.CPUUsage = cpuPercent[0]
	}

	return overview, nil
}

// GetResourceUsage 获取硬件资源使用情况
func (s *DashboardService) GetResourceUsage() (*model.ResourceUsage, error) {
	usage := &model.ResourceUsage{}

	// CPU 信息
	cpuPercent, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercent) > 0 {
		usage.CPU = model.CPUInfo{
			UsagePercent: cpuPercent[0],
		}
	}

	// 内存信息
	if memInfo, err := mem.VirtualMemory(); err == nil {
		usage.Memory = model.MemoryInfo{
			Total: memInfo.Total,
			Used:  memInfo.Used,
		}
	}

	// 磁盘信息
	if diskInfo, err := disk.Usage("/"); err == nil {
		usage.Disk = model.DiskInfo{
			Total: diskInfo.Total,
			Used:  diskInfo.Used,
		}
	}

	return usage, nil
}

// StatTrend 统计趋势数据
type StatTrend struct {
	Article []model.TrendData `json:"article"`
	Video   []model.TrendData `json:"video"`
	Comment []model.TrendData `json:"comment"`
	User    []model.TrendData `json:"user"`
}

// GetTrend 获取统计趋势数据
func (s *DashboardService) GetTrend(days int) (*StatTrend, error) {
	trend := &StatTrend{
		Article: []model.TrendData{},
		Video:   []model.TrendData{},
		Comment: []model.TrendData{},
		User:    []model.TrendData{},
	}

	startDate := time.Now().AddDate(0, 0, -days)

	// 文章趋势
	var articleTrends []struct {
		Date  string
		Count int64
	}
	q := db.GetDB().Model(&model.Article{}).Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).Group("DATE(created_at)").Order("DATE(created_at) ASC")
	if err := q.Find(&articleTrends).Error; err != nil {
		zap.L().Warn("获取文章趋势失败", zap.Error(err))
	}
	for _, at := range articleTrends {
		trend.Article = append(trend.Article, model.TrendData{Date: at.Date, Count: int64(at.Count)})
	}

	// 视频趋势
	var videoTrends []struct {
		Date  string
		Count int64
	}
	q2 := db.GetDB().Model(&model.Video{}).Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).Group("DATE(created_at)").Order("DATE(created_at) ASC")
	if err := q2.Find(&videoTrends).Error; err != nil {
		zap.L().Warn("获取视频趋势失败", zap.Error(err))
	}
	for _, vt := range videoTrends {
		trend.Video = append(trend.Video, model.TrendData{Date: vt.Date, Count: int64(vt.Count)})
	}

	// 评论趋势
	var commentTrends []struct {
		Date  string
		Count int64
	}
	q3 := db.GetDB().Model(&model.Comment{}).Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).Group("DATE(created_at)").Order("DATE(created_at) ASC")
	if err := q3.Find(&commentTrends).Error; err != nil {
		zap.L().Warn("获取评论趋势失败", zap.Error(err))
	}
	for _, ct := range commentTrends {
		trend.Comment = append(trend.Comment, model.TrendData{Date: ct.Date, Count: int64(ct.Count)})
	}

	// 用户趋势
	var userTrends []struct {
		Date  string
		Count int64
	}
	q4 := db.GetDB().Model(&model.User{}).Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= ?", startDate).Group("DATE(created_at)").Order("DATE(created_at) ASC")
	if err := q4.Find(&userTrends).Error; err != nil {
		zap.L().Warn("获取用户趋势失败", zap.Error(err))
	}
	for _, ut := range userTrends {
		trend.User = append(trend.User, model.TrendData{Date: ut.Date, Count: int64(ut.Count)})
	}

	return trend, nil
}

// RecentItem 最近活动项
type RecentItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetRecentActivityByLimit 获取最近活动（带数量限制）
func (s *DashboardService) GetRecentActivityByLimit(limit int) ([]RecentItem, error) {
	items := []RecentItem{}

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	// 最近文章
	var recentArticles []model.Article
	if err := db.GetDB().Model(&model.Article{}).Order("created_at DESC").Limit(limit / 3).Find(&recentArticles).Error; err != nil {
		zap.L().Warn("获取最近文章失败", zap.Error(err))
	}
	for _, a := range recentArticles {
		items = append(items, RecentItem{
			ID:        a.ID,
			Title:     a.Title,
			Type:      "article",
			Status:    a.Status,
			CreatedAt: a.CreatedAt,
		})
	}

	// 最近视频
	var recentVideos []model.Video
	if err := db.GetDB().Model(&model.Video{}).Order("created_at DESC").Limit(limit / 3).Find(&recentVideos).Error; err != nil {
		zap.L().Warn("获取最近视频失败", zap.Error(err))
	}
	for _, v := range recentVideos {
		items = append(items, RecentItem{
			ID:        v.ID,
			Title:     v.Title,
			Type:      "video",
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
		})
	}

	// 最近评论
	commentLimit := limit - len(recentArticles) - len(recentVideos)
	if commentLimit <= 0 {
		commentLimit = 1
	}
	var recentComments []model.Comment
	if err := db.GetDB().Model(&model.Comment{}).Order("created_at DESC").Limit(commentLimit).Find(&recentComments).Error; err != nil {
		zap.L().Warn("获取最近评论失败", zap.Error(err))
	}
	for _, c := range recentComments {
		items = append(items, RecentItem{
			ID:        c.ID,
			Title:     c.Content,
			Type:      "comment",
			Status:    c.Status,
			CreatedAt: c.CreatedAt,
		})
	}

	return items, nil
}

// GetRecentActivity 获取最近活动（使用默认数量10）
func (s *DashboardService) GetRecentActivity() ([]RecentItem, error) {
	return s.GetRecentActivityByLimit(10)
}

// getFileStats 获取文件统计信息
func (s *DashboardService) getFileStats(stats *model.DashboardStats) error {
	// 使用 media_assets 表统计文件（因为 file_uploads 是中间表）
	// 总文件数
	db.GetDB().Model(&model.MediaAsset{}).Where("status = ?", 1).Count(&stats.TotalFiles)

	// 今日新增文件
	today := time.Now().Truncate(24 * time.Hour)
	db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND created_at >= ?", 1, today).Count(&stats.TodayNewFiles)

	// 总文件大小
	var totalSize int64
	db.GetDB().Table("media_assets").Where("status = ?", 1).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize)
	stats.TotalFileSize = totalSize

	// 按类型统计
	db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "image", 1).Count(&stats.ImageCount)
	db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "video", 1).Count(&stats.VideoFileCount)
	db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "audio", 1).Count(&stats.AudioFileCount)
	db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "file", 1).Count(&stats.DocumentCount)

	return nil
}
