package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gocms_v3/app/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// NotificationService 通知服务
type NotificationService struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{
		db:  db,
		log: zap.L(),
	}
}

// PublishMessage 发布系统消息
func (s *NotificationService) PublishMessage(ctx context.Context, msg *model.SystemMessage, targetUserIDs []uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 保存系统消息
	msg.TargetUserIDs = encodeTargetUserIDs(targetUserIDs)
	msg.CreatedAt = time.Now()
	msg.UpdatedAt = time.Now()

	if err := tx.Create(msg).Error; err != nil {
		s.log.Error("publish message failed", zap.Error(err))
		tx.Rollback()
		return fmt.Errorf("发布消息失败: %w", err)
	}

	messageID := msg.ID

	// 为目标用户创建通知记录
	for _, userID := range targetUserIDs {
		notification := &model.UserNotification{
			UserID:    userID,
			MessageID: messageID,
			IsRead:    0,
		}

		if msg.IsReadDefault == 1 {
			notification.IsRead = 1
			notification.ReadAt = &msg.CreatedAt
		}

		if err := tx.Create(notification).Error; err != nil {
			s.log.Error("create user notification failed", zap.Error(err), zap.Uint("userID", userID))
			tx.Rollback()
			return fmt.Errorf("创建通知失败: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		s.log.Error("commit transaction failed", zap.Error(err))
		return fmt.Errorf("提交失败: %w", err)
	}

	s.log.Info("message published successfully", zap.Uint("messageID", messageID), zap.Int("targetUsers", len(targetUserIDs)))
	return nil
}

// GetMessages 获取系统消息列表（管理员视角）
func (s *NotificationService) GetMessages(query model.SystemMessageQuery) ([]model.SystemMessageResponse, int64, error) {
	var messages []model.SystemMessage
	var total int64

	queryBuilder := s.db.Model(&model.SystemMessage{}).Order("created_at DESC")

	// 添加过滤条件
	if query.Keyword != "" {
		queryBuilder = queryBuilder.Where("title LIKE ?", "%"+query.Keyword+"%")
	}
	if query.MessageType != "" {
		queryBuilder = queryBuilder.Where("message_type = ?", query.MessageType)
	}
	if query.Status != nil {
		queryBuilder = queryBuilder.Where("status = ?", *query.Status)
	}

	// 获取总数
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询消息总数失败: %w", err)
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := queryBuilder.Offset(offset).Limit(query.PageSize).Find(&messages).Error; err != nil {
		return nil, 0, fmt.Errorf("查询消息列表失败: %w", err)
	}

	// 转换为响应结构
	resp := make([]model.SystemMessageResponse, len(messages))
	for i, msg := range messages {
		resp[i] = model.SystemMessageResponse{
			ID:            msg.ID,
			Title:         msg.Title,
			Content:       msg.Content,
			MessageType:   msg.MessageType,
			SenderID:      msg.SenderID,
			SenderName:    msg.SenderName,
			Priority:      msg.Priority,
			TargetType:    msg.TargetType,
			TargetUserIDs: msg.TargetUserIDs,
			IsReadDefault: msg.IsReadDefault,
			ExpireAt:      msg.ExpireAt,
			Status:        msg.Status,
			CreatedAt:     msg.CreatedAt,
		}
	}

	return resp, total, nil
}

// GetMessageDetail 获取消息详情
func (s *NotificationService) GetMessageDetail(id uint) (*model.SystemMessageResponse, error) {
	var msg model.SystemMessage
	if err := s.db.First(&msg, id).Error; err != nil {
		return nil, fmt.Errorf("消息不存在: %w", err)
	}

	return &model.SystemMessageResponse{
		ID:            msg.ID,
		Title:         msg.Title,
		Content:       msg.Content,
		MessageType:   msg.MessageType,
		SenderID:      msg.SenderID,
		SenderName:    msg.SenderName,
		Priority:      msg.Priority,
		TargetType:    msg.TargetType,
		TargetUserIDs: msg.TargetUserIDs,
		IsReadDefault: msg.IsReadDefault,
		ExpireAt:      msg.ExpireAt,
		Status:        msg.Status,
		CreatedAt:     msg.CreatedAt,
	}, nil
}

// DeleteMessage 删除消息
func (s *NotificationService) DeleteMessage(id uint) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除消息
	if err := tx.Delete(&model.SystemMessage{}, id).Error; err != nil {
		return fmt.Errorf("删除消息失败: %w", err)
	}

	// 删除关联的用户通知
	if err := tx.Where("message_id = ?", id).Delete(&model.UserNotification{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("删除关联通知失败: %w", err)
	}

	return tx.Commit().Error
}

// GetMyNotifications 获取我的通知列表
func (s *NotificationService) GetMyNotifications(query model.UserNotificationQuery) ([]model.UserNotificationResponse, int64, error) {
	type NotificationWithMessage struct {
		ID               uint       `gorm:"column:id"`
		UserID           uint       `gorm:"column:user_id"`
		MessageID        uint       `gorm:"column:message_id"`
		IsRead           int8       `gorm:"column:is_read"`
		ReadAt           *time.Time `gorm:"column:read_at"`
		CreatedAt        time.Time  `gorm:"column:created_at"`
		Title            string     `gorm:"column:title"`
		Content          string     `gorm:"column:content"`
		MessageType      string     `gorm:"column:message_type"`
		Priority         int8       `gorm:"column:priority"`
		MessageCreatedAt time.Time  `gorm:"column:message_created_at"`
	}

	var notifications []NotificationWithMessage
	var total int64

	queryBuilder := s.db.Model(&model.UserNotification{}).
		Joins("JOIN system_messages m ON user_notifications.message_id = m.id").
		Select("user_notifications.id, user_notifications.user_id, user_notifications.message_id, user_notifications.is_read, user_notifications.read_at, user_notifications.created_at, m.title, m.content, m.message_type, m.priority, m.created_at as message_created_at").
		Order("user_notifications.created_at DESC")

	// 过滤条件
	queryBuilder = queryBuilder.Where("user_notifications.user_id = ?", query.UserID)
	if query.IsRead != nil {
		queryBuilder = queryBuilder.Where("user_notifications.is_read = ?", *query.IsRead)
	}
	if query.MessageType != "" {
		queryBuilder = queryBuilder.Where("m.message_type = ?", query.MessageType)
	}

	// 获取总数
	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询通知总数失败: %w", err)
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	if err := queryBuilder.Offset(offset).Limit(query.PageSize).Find(&notifications).Error; err != nil {
		return nil, 0, fmt.Errorf("查询通知列表失败: %w", err)
	}

	// 转换为响应结构
	resp := make([]model.UserNotificationResponse, len(notifications))
	for i, n := range notifications {
		resp[i] = model.UserNotificationResponse{
			ID:        n.ID,
			UserID:    n.UserID,
			MessageID: n.MessageID,
			IsRead:    n.IsRead,
			ReadAt:    n.ReadAt,
			CreatedAt: n.CreatedAt,
			Message: model.MessageInfo{
				ID:          n.MessageID,
				Title:       n.Title,
				Content:     n.Content,
				MessageType: n.MessageType,
				Priority:    n.Priority,
				CreatedAt:   n.MessageCreatedAt,
			},
		}
	}

	return resp, total, nil
}

// GetMyNotificationStats 获取我的通知统计
func (s *NotificationService) GetMyNotificationStats(userID uint) (*model.NotificationStats, error) {
	stats := &model.NotificationStats{}

	// 总数量
	s.db.Model(&model.UserNotification{}).Where("user_id = ?", userID).Count(&stats.TotalCount)

	// 未读数量
	s.db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = 0", userID).Count(&stats.UnreadCount)

	// 已读数量
	s.db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = 1", userID).Count(&stats.ReadCount)

	// 按消息类型统计（需要JOIN）
	type CountResult struct {
		Count int64 `json:"count"`
	}

	var result []CountResult
	s.db.Table("user_notifications un").
		Select("COUNT(*) as count").
		Joins("JOIN system_messages sm ON un.message_id = sm.id").
		Where("un.user_id = ? AND sm.message_type = ?", userID, "system").
		Scan(&result)
	if len(result) > 0 {
		stats.SystemCount = result[0].Count
	}

	result = nil
	s.db.Table("user_notifications un").
		Select("COUNT(*) as count").
		Joins("JOIN system_messages sm ON un.message_id = sm.id").
		Where("un.user_id = ? AND sm.message_type = ?", userID, "publish").
		Scan(&result)
	if len(result) > 0 {
		stats.PublishCount = result[0].Count
	}

	result = nil
	s.db.Table("user_notifications un").
		Select("COUNT(*) as count").
		Joins("JOIN system_messages sm ON un.message_id = sm.id").
		Where("un.user_id = ? AND sm.message_type = ?", userID, "comment_reply").
		Scan(&result)
	if len(result) > 0 {
		stats.CommentCount = result[0].Count
	}

	result = nil
	s.db.Table("user_notifications un").
		Select("COUNT(*) as count").
		Joins("JOIN system_messages sm ON un.message_id = sm.id").
		Where("un.user_id = ? AND (sm.message_type = ? OR sm.message_type = ?)", userID, "audit_passed", "audited_failed").
		Scan(&result)
	if len(result) > 0 {
		stats.AuditCount = result[0].Count
	}

	return stats, nil
}

// MarkAsRead 标记通知为已读/未读
func (s *NotificationService) MarkAsRead(notificationID uint, isRead int8) error {
	if err := s.db.Model(&model.UserNotification{}).Where("id = ?", notificationID).Update("is_read", isRead).Error; err != nil {
		return fmt.Errorf("标记通知状态失败: %w", err)
	}

	now := time.Now()
	if isRead == 1 {
		s.db.Model(&model.UserNotification{}).Where("id = ?", notificationID).Update("read_at", now)
	} else {
		s.db.Model(&model.UserNotification{}).Where("id = ? AND read_at IS NOT NULL", notificationID).Update("read_at", nil)
	}

	return nil
}

// MarkMessageAsRead 标记消息对所有接收者为已读
func (s *NotificationService) MarkMessageAsRead(messageID uint) error {
	if err := s.db.Model(&model.UserNotification{}).Where("message_id = ?", messageID).Updates(map[string]interface{}{
		"is_read": 1,
		"read_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("标记消息已读失败: %w", err)
	}
	return nil
}

// MarkAllAsRead 将用户所有通知标记为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	if err := s.db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = 0", userID).Updates(map[string]interface{}{
		"is_read": 1,
		"read_at": time.Now(),
	}).Error; err != nil {
		return fmt.Errorf("全部标记已读失败: %w", err)
	}
	return nil
}

// WithdrawMessage 撤回消息
func (s *NotificationService) WithdrawMessage(id uint, userID uint) error {
	if err := s.db.Model(&model.SystemMessage{}).Where("id = ? AND sender_id = ?", id, userID).Update("status", 2).Error; err != nil {
		return fmt.Errorf("撤回消息失败: %w", err)
	}
	return nil
}

// GetNotificationsByMessageID 根据消息ID获取所有通知
func (s *NotificationService) GetNotificationsByMessageID(messageID uint) ([]model.UserNotification, error) {
	var notifications []model.UserNotification
	if err := s.db.Where("message_id = ?", messageID).Find(&notifications).Error; err != nil {
		return nil, fmt.Errorf("查询通知失败: %w", err)
	}
	return notifications, nil
}

// encodeTargetUserIDs 将用户ID数组编码为JSON字符串
func encodeTargetUserIDs(userIDs []uint) string {
	data, err := json.Marshal(userIDs)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// decodeTargetUserIDs 将JSON字符串解码为用户ID数组
func decodeTargetUserIDs(s string) []uint {
	var userIDs []uint
	if s == "" || s == "[]" {
		return userIDs
	}
	err := json.Unmarshal([]byte(s), &userIDs)
	if err != nil {
		return userIDs
	}
	return userIDs
}
