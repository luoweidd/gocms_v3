package service

import (
	"fmt"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
)

// OperationLogService 操作日志服务
type OperationLogService struct{}

// NewOperationLogService 创建操作日志服务
func NewOperationLogService() *OperationLogService {
	return &OperationLogService{}
}

// GetList 获取操作日志列表
func (s *OperationLogService) GetList(req model.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 200 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	query := db.GetDB().Model(&model.OperationLog{}).Order("created_at DESC")

	// 条件筛选
	if req.Module != "" {
		query = query.Where("module LIKE ?", "%"+req.Module+"%")
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	if req.Username != "" {
		query = query.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.IPAddress != "" {
		query = query.Where("ip_address LIKE ?", "%"+req.IPAddress+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartDate != "" {
		if startDate, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			query = query.Where("created_at >= ?", startDate)
		}
	}
	if req.EndDate != "" {
		if endDate, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			endDate = endDate.Add(24 * time.Hour)
			query = query.Where("created_at < ?", endDate)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计操作日志总数失败：%w", err)
	}

	var logs []model.OperationLog
	if err := query.Limit(req.PageSize).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, fmt.Errorf("查询操作日志失败：%w", err)
	}

	return logs, total, nil
}

// GetStats 获取操作日志统计
func (s *OperationLogService) GetStats() (*model.OperationLogStats, error) {
	stats := &model.OperationLogStats{}

	// 总数
	db.GetDB().Model(&model.OperationLog{}).Count(&stats.TotalCount)

	// 今日数
	today := time.Now().Truncate(24 * time.Hour)
	db.GetDB().Model(&model.OperationLog{}).Where("created_at >= ?", today).Count(&stats.TodayCount)

	// 成功/失败
	db.GetDB().Model(&model.OperationLog{}).Where("status = 1").Count(&stats.SuccessCount)
	db.GetDB().Model(&model.OperationLog{}).Where("status = 0").Count(&stats.FailCount)

	// 按模块统计
	type moduleStat struct {
		Module string
		Count  int64
	}
	var moduleStats []moduleStat
	db.GetDB().Model(&model.OperationLog{}).Select("module, COUNT(*) as count").Group("module").Order("count DESC").Limit(10).Find(&moduleStats)
	for _, ms := range moduleStats {
		stats.ModuleStats = append(stats.ModuleStats, model.ModuleStat{
			Module: ms.Module,
			Count:  ms.Count,
		})
	}

	// 按操作类型统计
	type actionStat struct {
		Action string
		Count  int64
	}
	var actionStats []actionStat
	db.GetDB().Model(&model.OperationLog{}).Select("action, COUNT(*) as count").Group("action").Order("count DESC").Limit(10).Find(&actionStats)
	for _, ast := range actionStats {
		stats.ActionStats = append(stats.ActionStats, model.ActionStat{
			Action: ast.Action,
			Count:  ast.Count,
		})
	}

	return stats, nil
}

// GetByID 获取单条操作日志
func (s *OperationLogService) GetByID(id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	if err := db.GetDB().Where("id = ?", id).First(&log).Error; err != nil {
		return nil, fmt.Errorf("操作日志不存在：%w", err)
	}
	return &log, nil
}

// DeleteExpired 删除过期日志（保留最近N天）
func (s *OperationLogService) DeleteExpired(days int) (int64, error) {
	if days <= 0 {
		days = 90 // 默认保留90天
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	result := db.GetDB().Where("created_at < ?", cutoff).Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// DeleteByID 删除单条操作日志
func (s *OperationLogService) DeleteByID(id uint) error {
	return db.GetDB().Delete(&model.OperationLog{}, id).Error
}

// BatchDelete 批量删除操作日志
func (s *OperationLogService) BatchDelete(ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := db.GetDB().Where("id IN ?", ids).Delete(&model.OperationLog{})
	return result.RowsAffected, result.Error
}

// CleanAll 清空操作日志（危险操作）
func (s *OperationLogService) CleanAll() (int64, error) {
	result := db.GetDB().Exec("DELETE FROM operation_logs")
	return result.RowsAffected, result.Error
}
