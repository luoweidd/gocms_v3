package service

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
)

// ContentAuditService 内容审核服务
type ContentAuditService struct{}

// NewContentAuditService 创建内容审核服务
func NewContentAuditService() *ContentAuditService {
	return &ContentAuditService{}
}

// GenerateAuditNo 生成审核编号
func (s *ContentAuditService) GenerateAuditNo() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	prefix := "AUD" + dateStr

	var maxNo string
	db.GetDB().Model(&model.ContentAudit{}).Where("audit_no LIKE ?", prefix+"%").Order("audit_no DESC").Pluck("audit_no", &maxNo)

	if maxNo != "" {
		suffix := strings.TrimPrefix(maxNo, prefix)
		if num, err := strconv.Atoi(suffix); err == nil {
			return prefix + fmt.Sprintf("%04d", num+1)
		}
	}
	return prefix + "0001"
}

// SubmitAudit 提交审核
func (s *ContentAuditService) SubmitAudit(req struct {
	ContentType   string
	ContentID     uint
	ContentTitle  string
	SubmitterID   uint
	SubmitterName string
	NextStatus    int8
}) (*model.ContentAudit, error) {
	log.Printf("[SERVICE] 提交审核 - 类型: %s, ID: %d", req.ContentType, req.ContentID)

	audit := &model.ContentAudit{
		AuditNo:        s.GenerateAuditNo(),
		ContentType:    req.ContentType,
		ContentID:      req.ContentID,
		ContentTitle:   req.ContentTitle,
		SubmitterID:    req.SubmitterID,
		SubmitterName:  req.SubmitterName,
		Status:         0,
		PreviousStatus: -1,
		NextStatus:     req.NextStatus,
		AuditedAt:      nil,
	}

	if err := db.GetDB().Create(audit).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 提交审核失败 - 错误: %v", err)
		return nil, fmt.Errorf("提交审核失败: %w", err)
	}

	log.Printf("[SERVICE] 提交审核成功 - 编号: %s", audit.AuditNo)
	return audit, nil
}

// GetAuditList 获取审核列表
func (s *ContentAuditService) GetAuditList(req model.ContentAuditListRequest) (list []model.ContentAudit, total int64, err error) {
	query := db.GetDB().Model(&model.ContentAudit{})

	if req.ContentType != "" {
		query = query.Where("content_type = ?", req.ContentType)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("audit_no LIKE ? OR content_title LIKE ? OR submitter_name LIKE ?", keyword, keyword, keyword)
	}
	if req.AuditorID > 0 {
		query = query.Where("auditor_id = ?", req.AuditorID)
	}

	if err = query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	if err = query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("查询失败: %w", err)
	}

	return list, total, nil
}

// GetAuditByID 获取审核详情
func (s *ContentAuditService) GetAuditByID(id uint) (*model.ContentAudit, error) {
	var audit model.ContentAudit
	if err := db.GetDB().First(&audit, id).Error; err != nil {
		return nil, fmt.Errorf("审核记录不存在")
	}
	return &audit, nil
}

// Approve 通过审核
func (s *ContentAuditService) Approve(id uint, auditorID uint, auditorName string, comment string) error {
	log.Printf("[SERVICE] 通过审核 - ID: %d", id)

	var audit model.ContentAudit
	if err := db.GetDB().First(&audit, id).Error; err != nil {
		return fmt.Errorf("审核记录不存在")
	}

	if audit.Status == 1 {
		return fmt.Errorf("该审核已通过")
	}

	now := time.Now().Format(time.RFC3339)
	updates := map[string]interface{}{
		"status":        1,
		"auditor_id":    auditorID,
		"auditor_name":  auditorName,
		"audit_comment": comment,
		"next_status":   1,
		"audited_at":    now,
	}

	if err := db.GetDB().Model(&audit).Updates(updates).Error; err != nil {
		return fmt.Errorf("审核失败: %w", err)
	}

	log.Printf("[SERVICE] 通过审核成功 - 编号: %s", audit.AuditNo)
	return nil
}

// Reject 驳回审核
func (s *ContentAuditService) Reject(id uint, auditorID uint, auditorName string, comment string, requiredMod string) error {
	log.Printf("[SERVICE] 驳回审核 - ID: %d", id)

	var audit model.ContentAudit
	if err := db.GetDB().First(&audit, id).Error; err != nil {
		return fmt.Errorf("审核记录不存在")
	}

	now := time.Now().Format(time.RFC3339)
	updates := map[string]interface{}{
		"status":                2,
		"auditor_id":            auditorID,
		"auditor_name":          auditorName,
		"audit_comment":         comment,
		"next_status":           0,
		"required_modification": requiredMod,
		"audited_at":            now,
	}

	if err := db.GetDB().Model(&audit).Updates(updates).Error; err != nil {
		return fmt.Errorf("驳回失败: %w", err)
	}

	log.Printf("[SERVICE] 驳回审核成功 - 编号: %s", audit.AuditNo)
	return nil
}

// NeedsModification 标记需修改
func (s *ContentAuditService) NeedsModification(id uint, auditorID uint, auditorName string, comment string, requiredMod string) error {
	log.Printf("[SERVICE] 标记需修改 - ID: %d", id)

	var audit model.ContentAudit
	if err := db.GetDB().First(&audit, id).Error; err != nil {
		return fmt.Errorf("审核记录不存在")
	}

	now := time.Now().Format(time.RFC3339)
	updates := map[string]interface{}{
		"status":                3,
		"auditor_id":            auditorID,
		"auditor_name":          auditorName,
		"audit_comment":         comment,
		"next_status":           -1,
		"required_modification": requiredMod,
		"audited_at":            now,
	}

	if err := db.GetDB().Model(&audit).Updates(updates).Error; err != nil {
		return fmt.Errorf("标记失败: %w", err)
	}

	return nil
}

// BatchAudit 批量审核
func (s *ContentAuditService) BatchAudit(auditorID uint, auditorName string, ids []uint, action string, comment string, requiredMod string) error {
	log.Printf("[SERVICE] 批量审核 - 数量: %d, 操作: %s", len(ids), action)

	if len(ids) == 0 {
		return nil
	}

	now := time.Now().Format(time.RFC3339)
	var updates map[string]interface{}

	switch action {
	case "approve":
		updates = map[string]interface{}{
			"status":        1,
			"auditor_id":    auditorID,
			"auditor_name":  auditorName,
			"audit_comment": comment,
			"next_status":   1,
			"audited_at":    now,
		}
	case "reject":
		updates = map[string]interface{}{
			"status":                2,
			"auditor_id":            auditorID,
			"auditor_name":          auditorName,
			"audit_comment":         comment,
			"next_status":           0,
			"required_modification": requiredMod,
			"audited_at":            now,
		}
	default:
		return fmt.Errorf("无效的操作: %s", action)
	}

	if err := db.GetDB().Model(&model.ContentAudit{}).Where("id in ?", ids).Updates(updates).Error; err != nil {
		return fmt.Errorf("批量审核失败: %w", err)
	}

	log.Printf("[SERVICE] 批量审核成功")
	return nil
}

// GetAuditStats 获取审核统计
func (s *ContentAuditService) GetAuditStats() (*model.AuditStats, error) {
	stats := &model.AuditStats{}

	db.GetDB().Model(&model.ContentAudit{}).Count(&stats.Total)
	db.GetDB().Model(&model.ContentAudit{}).Where("status = ?", 0).Count(&stats.Pending)
	db.GetDB().Model(&model.ContentAudit{}).Where("status = ?", 1).Count(&stats.Passed)
	db.GetDB().Model(&model.ContentAudit{}).Where("status = ?", 2).Count(&stats.Rejected)
	db.GetDB().Model(&model.ContentAudit{}).Where("status = ?", 3).Count(&stats.NeedsModify)

	return stats, nil
}
