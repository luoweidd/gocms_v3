package controller

import (
	"fmt"
	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationLogController 操作日志控制器
type OperationLogController struct {
	service *service.OperationLogService
}

// NewOperationLogController 创建控制器
func NewOperationLogController() *OperationLogController {
	return &OperationLogController{
		service: service.NewOperationLogService(),
	}
}

// GetList 获取操作日志列表
func (ctrl *OperationLogController) GetList(c *gin.Context) {
	var req model.OperationLogListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, 400, "参数错误："+err.Error())
		return
	}

	logs, total, err := ctrl.service.GetList(req)
	if err != nil {
		zap.L().Error("获取操作日志列表失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"list":     logs,
		"total":    total,
		"page":     req.Page,
		"pageSize": req.PageSize,
	})
}

// GetStats 获取操作日志统计
func (ctrl *OperationLogController) GetStats(c *gin.Context) {
	stats, err := ctrl.service.GetStats()
	if err != nil {
		zap.L().Error("获取操作日志统计失败", zap.Error(err))
		response.Error(c, 500, "查询失败")
		return
	}

	response.Success(c, stats)
}

// GetByID 获取单条操作日志详情
func (ctrl *OperationLogController) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		response.Error(c, 400, "无效ID")
		return
	}

	log, err := ctrl.service.GetByID(id)
	if err != nil {
		zap.L().Warn("操作日志不存在", zap.Uint("id", id))
		response.Error(c, 404, "操作日志不存在")
		return
	}

	response.Success(c, log)
}

// Delete 删除单条操作日志
func (ctrl *OperationLogController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	var id uint
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		response.Error(c, 400, "无效ID")
		return
	}

	if err := ctrl.service.DeleteByID(id); err != nil {
		zap.L().Error("删除操作日志失败", zap.Error(err), zap.Uint("id", id))
		response.Error(c, 500, "删除失败")
		return
	}

	response.Success(c, nil)
}

// BatchDelete 批量删除操作日志
func (ctrl *OperationLogController) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "参数错误：IDs必填")
		return
	}

	affected, err := ctrl.service.BatchDelete(req.IDs)
	if err != nil {
		zap.L().Error("批量删除操作日志失败", zap.Error(err))
		response.Error(c, 500, "删除失败")
		return
	}

	response.Success(c, gin.H{"affected": affected})
}

// CleanExpired 清理过期操作日志
func (ctrl *OperationLogController) CleanExpired(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "90")
	var days int
	if _, err := fmt.Sscanf(daysStr, "%d", &days); err != nil || days <= 0 {
		days = 90
	}

	affected, err := ctrl.service.DeleteExpired(days)
	if err != nil {
		zap.L().Error("清理过期操作日志失败", zap.Error(err))
		response.Error(c, 500, "清理失败")
		return
	}

	response.Success(c, gin.H{"cleaned": affected, "keep_days": days})
}

// CleanAll 清空所有操作日志
func (ctrl *OperationLogController) CleanAll(c *gin.Context) {
	affected, err := ctrl.service.CleanAll()
	if err != nil {
		zap.L().Error("清空操作日志失败", zap.Error(err))
		response.Error(c, 500, "清空失败")
		return
	}

	response.Success(c, gin.H{"cleaned": affected})
}
