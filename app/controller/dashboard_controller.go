package controller

import (
	"fmt"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardController struct {
	dashboardService *service.DashboardService
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		dashboardService: service.NewDashboardService(),
	}
}

// GetStats 获取仪表盘统计数据
func (ctrl *DashboardController) GetStats(c *gin.Context) {
	stats, err := ctrl.dashboardService.GetStats()
	if err != nil {
		zap.L().Error("获取仪表盘统计数据失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, stats)
}

// GetOverview 获取仪表盘概览数据
func (ctrl *DashboardController) GetOverview(c *gin.Context) {
	overview, err := ctrl.dashboardService.GetOverview()
	if err != nil {
		zap.L().Error("获取仪表盘概览数据失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, overview)
}

// GetTrend 获取统计趋势数据 API: GET /dashboard/trend
func (ctrl *DashboardController) GetTrend(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "7")
	var days int
	fmt.Sscanf(daysStr, "%d", &days)
	if days <= 0 || days > 90 {
		days = 7
	}

	trend, err := ctrl.dashboardService.GetTrend(days)
	if err != nil {
		zap.L().Error("获取统计趋势数据失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, trend)
}

// GetRecentActivity 获取最近活动 API: GET /dashboard/recent-activity
func (ctrl *DashboardController) GetRecentActivity(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	var limit int
	fmt.Sscanf(limitStr, "%d", &limit)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	items, err := ctrl.dashboardService.GetRecentActivityByLimit(limit)
	if err != nil {
		zap.L().Error("获取最近活动失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, items)
}

// GetResourceUsage 获取硬件资源使用情况 API: GET /dashboard/resource-usage
func (ctrl *DashboardController) GetResourceUsage(c *gin.Context) {
	usage, err := ctrl.dashboardService.GetResourceUsage()
	if err != nil {
		zap.L().Error("获取硬件资源使用情况失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, usage)
}
