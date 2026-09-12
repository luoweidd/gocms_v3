package service

import (
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"gorm.io/gorm"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

// CreateComment 创建评论（兼容controller调用）
func (s *CommentService) CreateComment(req model.CommentCreateRequest) (*model.Comment, error) {
	return s.Create(&req, "", "")
}

// AuditComment 审核评论（兼容controller调用）
func (s *CommentService) AuditComment(id uint, status int, remark string, auditBy uint) error {
	return s.Audit(id, status, auditBy, remark)
}

// GetCommentList 获取评论列表（兼容controller调用）
func (s *CommentService) GetCommentList(query model.CommentListQuery) (*response.PageData, error) {
	comments, total, err := s.List(&query)
	if err != nil {
		return nil, err
	}
	return &response.PageData{
		List:     comments,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

// GetCommentStats 获取评论统计
func (s *CommentService) GetCommentStats() map[string]interface{} {
	stats := make(map[string]interface{})
	var total int64
	db.GetDB().Model(&model.Comment{}).Count(&total)
	stats["total"] = total

	var pending int64
	db.GetDB().Model(&model.Comment{}).Where("status = 0").Count(&pending)
	stats["pending"] = pending

	var approved int64
	db.GetDB().Model(&model.Comment{}).Where("status = 1").Count(&approved)
	stats["approved"] = approved

	return stats
}

// Create 创建评论
func (s *CommentService) Create(req *model.CommentCreateRequest, ipAddr, userAgent string) (*model.Comment, error) {
	comment := &model.Comment{
		ArticleID: func() uint {
			if req.ArticleID != nil {
				return *req.ArticleID
			}
			return 0
		}(),
		VideoID: func() uint {
			if req.VideoID != nil {
				return *req.VideoID
			}
			return 0
		}(),
		UserID: 0, // TODO: 从上下文获取用户ID
		Nickname: func() string {
			if req.Nickname != "" {
				return req.Nickname
			}
			return "匿名用户"
		}(),
		Email:     req.Email,
		Content:   req.Content,
		Status:    0, // 默认待审核
		ParentID:  req.ParentID,
		IPAddress: ipAddr,
		UserAgent: userAgent,
	}

	if err := db.GetDB().Create(comment).Error; err != nil {
		return nil, err
	}
	return comment, nil
}

// List 获取评论列表
func (s *CommentService) List(query *model.CommentListQuery) ([]model.CommentListItem, int64, error) {
	var comments []model.CommentListItem
	var total int64

	q := db.GetDB().Model(&model.Comment{}).Scopes(s.setListScope(query)).Order("created_at DESC")
	if err := q.Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// TODO: 获取总数
	return comments, total, nil
}

func (s *CommentService) setListScope(query *model.CommentListQuery) func(db *gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		if query.ArticleID > 0 {
			d = d.Where("article_id = ?", query.ArticleID)
		}
		if query.VideoID > 0 {
			d = d.Where("video_id = ?", query.VideoID)
		}
		if query.UserID > 0 {
			d = d.Where("user_id = ?", query.UserID)
		}
		return d
	}
}

// Audit 审核评论
func (s *CommentService) Audit(id uint, status int, auditBy uint, remark string) error {
	now := time.Now()
	return db.GetDB().Model(&model.Comment{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       status,
		"audit_by":     auditBy,
		"audit_remark": remark,
		"audited_at":   now,
	}).Error
}

// Delete 删除评论
func (s *CommentService) Delete(id uint) error {
	return db.GetDB().Delete(&model.Comment{}, id).Error
}

// BatchDeleteComments 批量删除评论
func (s *CommentService) BatchDeleteComments(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return db.GetDB().Where("id IN ?", ids).Delete(&model.Comment{}).Error
}

// Like 点赞评论
func (s *CommentService) Like(id uint) error {
	return db.GetDB().Model(&model.Comment{}).Where("id = ?", id).Update("likes", gorm.Expr("likes + 1")).Error
}

// BatchApprove 批量通过评论
func (s *CommentService) BatchApprove(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return db.GetDB().Model(&model.Comment{}).
		Where("id IN ? AND status = ?", ids, 0). // 只审核待审核状态的
		Updates(map[string]interface{}{
			"status":     1,
			"audited_at": now,
		}).Error
}

// BatchReject 批量驳回评论
func (s *CommentService) BatchReject(ids []uint, reason string) error {
	if len(ids) == 0 {
		return nil
	}
	now := time.Now()
	return db.GetDB().Model(&model.Comment{}).
		Where("id IN ? AND status = ?", ids, 0). // 只驳回待审核状态的
		Updates(map[string]interface{}{
			"status":       2,
			"audit_remark": reason,
			"audited_at":   now,
		}).Error
}

// CommentAuditStats 评论审核统计
type CommentAuditStats struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Approved  int64 `json:"approved"`
	Rejected  int64 `json:"rejected"`
	Deleted   int64 `json:"deleted"`
	ByArticle int64 `json:"by_article"` // 今日新增(文章评论)
	ByVideo   int64 `json:"by_video"`   // 今日新增(视频评论)
}

// GetAuditStats 获取审核统计
func (s *CommentService) GetAuditStats() (*CommentAuditStats, error) {
	stats := &CommentAuditStats{}

	// 总数
	db.GetDB().Model(&model.Comment{}).Count(&stats.Total)

	// 各状态数量
	db.GetDB().Model(&model.Comment{}).Where("status = ?", 0).Count(&stats.Pending)
	db.GetDB().Model(&model.Comment{}).Where("status = ?", 1).Count(&stats.Approved)
	db.GetDB().Model(&model.Comment{}).Where("status = ?", 2).Count(&stats.Rejected)
	db.GetDB().Model(&model.Comment{}).Where("status = ?", 3).Count(&stats.Deleted)

	// 今日新增统计
	todayStart := time.Now().Truncate(24 * time.Hour)
	var articleCount, videoCount int64
	db.GetDB().Model(&model.Comment{}).Where("article_id > 0 AND created_at >= ?", todayStart).Count(&articleCount)
	db.GetDB().Model(&model.Comment{}).Where("video_id > 0 AND created_at >= ?", todayStart).Count(&videoCount)
	stats.ByArticle = articleCount
	stats.ByVideo = videoCount

	return stats, nil
}

// BatchApproveByController 批量通过（兼容controller调用）
func (s *CommentService) BatchApproveByController(ids []uint) error {
	return s.BatchApprove(ids)
}

// BatchRejectByController 批量驳回（兼容controller调用）
func (s *CommentService) BatchRejectByController(ids []uint, reason string) error {
	return s.BatchReject(ids, reason)
}
