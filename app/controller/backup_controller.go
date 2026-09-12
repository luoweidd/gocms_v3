package controller

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// BackupController 备份控制器
type BackupController struct {
	backupService *service.BackupService
}

// NewBackupController 创建备份控制器
func NewBackupController() *BackupController {
	return &BackupController{}
}

// SetService 设置服务（避免循环依赖）
func (c *BackupController) SetService(svc *service.BackupService) {
	c.backupService = svc
}

// CreateBackup 创建备份
// @Summary 创建数据备份
// @Description 创建全量或增量数据备份
// @Tags 数据备份
// @Param data body model.BackupCreateRequest true "备份参数"
// @Success 200 {object} response.Response
// @Router /backup/full [post]
func (c *BackupController) CreateBackup(ginCtx *gin.Context) {
	var req model.BackupCreateRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	if req.BackupType == "" {
		req.BackupType = "full"
	}
	if req.Name == "" {
		req.Name = "手动备份_" + c.getTimeStr()
	}

	createdBy := ginCtx.GetUint("user_id")
	backup, err := c.backupService.CreateBackup(ginCtx.Request.Context(), req, createdBy)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "备份已创建，请稍后查看状态", backup)
}

// GetBackupList 获取备份列表
// @Summary 获取备份列表
// @Description 获取备份记录列表
// @Tags 数据备份
// @Param page query int true "页码"
// @Param page_size query int true "每页数量"
// @Param keyword query string false "关键词"
// @Param backup_type query string false "备份类型"
// @Param status query int false "状态"
// @Success 200 {object} response.Response
// @Router /backup/list [get]
func (c *BackupController) GetBackupList(ginCtx *gin.Context) {
	query := model.BackupRecordQuery{
		Page:       1,
		PageSize:   10,
		Keyword:    ginCtx.Query("keyword"),
		BackupType: ginCtx.Query("backup_type"),
	}

	pageStr := ginCtx.Query("page")
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			query.Page = p
		}
	}

	pageSizeStr := ginCtx.Query("page_size")
	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			query.PageSize = ps
		}
	}

	statusStr := ginCtx.Query("status")
	if statusStr != "" {
		if s, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			si := int8(s)
			query.Status = &si
		}
	}

	list, total, err := c.backupService.GetBackupList(query)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, gin.H{
		"list":     list,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetBackupDetail 获取备份详情
// @Summary 获取备份详情
// @Description 获取备份记录详情
// @Tags 数据备份
// @Param id path int true "备份ID"
// @Success 200 {object} response.Response
// @Router /backup/:id [get]
func (c *BackupController) GetBackupDetail(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的备份ID")
		return
	}

	detail, err := c.backupService.GetBackupDetail(uint(id))
	if err != nil {
		response.Error(ginCtx, 404, err.Error())
		return
	}

	response.Success(ginCtx, detail)
}

// DeleteBackup 删除备份
// @Summary 删除备份
// @Description 删除指定的备份记录
// @Tags 数据备份
// @Param id path int true "备份ID"
// @Success 200 {object} response.Response
// @Router /backup/:id [delete]
func (c *BackupController) DeleteBackup(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的备份ID")
		return
	}

	if err := c.backupService.DeleteBackup(uint(id)); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "备份删除成功", nil)
}

// DownloadBackup 下载备份文件
// @Summary 下载备份文件
// @Description 下载指定的备份文件
// @Tags 数据备份
// @Param id path int true "备份ID"
// @Success 200 {file} file "backup.sql"
// @Router /backup/:id/download [get]
func (c *BackupController) DownloadBackup(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的备份ID")
		return
	}

	filePath, err := c.backupService.DownloadBackup(uint(id))
	if err != nil {
		response.Error(ginCtx, 404, err.Error())
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.Error(ginCtx, 404, "备份文件不存在")
		return
	}

	// 获取原始文件名
	originalName := filepath.Base(filePath)

	// 设置 Content-Disposition header 强制下载并指定文件名 (使用URL编码确保正确显示)
	ginCtx.Header("Content-Disposition", "attachment; filename=\""+originalName+"\"")
	ginCtx.Header("Content-Type", "application/octet-stream")
	ginCtx.File(filePath)
}

// RestoreBackup 恢复备份
// @Summary 恢复备份
// @Description 从指定的备份文件恢复数据（需要二次确认）
// @Tags 数据备份
// @Param id path int true "备份ID"
// @Success 200 {object} response.Response
// @Router /backup/:id/restore [post]
func (c *BackupController) RestoreBackup(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的备份ID")
		return
	}

	filePath, err := c.backupService.RestoreBackup(uint(id))
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "备份文件路径已获取，请确认恢复操作", gin.H{
		"file_path": filePath,
	})
}

// GetBackupStats 获取备份统计
// @Summary 获取备份统计
// @Description 获取备份统计信息
// @Tags 数据备份
// @Success 200 {object} response.Response
// @Router /backup/stats [get]
func (c *BackupController) GetBackupStats(ginCtx *gin.Context) {
	stats, err := c.backupService.GetBackupStats()
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, stats)
}

// getTimeStr 获取时间字符串
func (c *BackupController) getTimeStr() string {
	return time.Now().Format("2006-01-02_15:04:05")
}
