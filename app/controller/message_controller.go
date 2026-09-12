package controller

import (
	"context"
	"strconv"
	"time"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// MessageController 消息控制器
type MessageController struct {
	notificationService *service.NotificationService
}

// NewMessageController 创建消息控制器
func NewMessageController(notificationService *service.NotificationService) *MessageController {
	return &MessageController{
		notificationService: notificationService,
	}
}

// PublishMessage 发布系统消息
func (c *MessageController) PublishMessage(ginCtx *gin.Context) {
	var req model.PublishMessageRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	msg := &model.SystemMessage{
		Title:         req.Title,
		Content:       req.Content,
		MessageType:   req.MessageType,
		Priority:      req.Priority,
		TargetType:    req.TargetType,
		IsReadDefault: req.IsReadDefault,
		Status:        1,
	}

	// 处理过期时间
	if req.ExpireAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.ExpireAt); err == nil {
			msg.ExpireAt = &t
		} else if t, err := time.Parse("2006-01-02", req.ExpireAt); err == nil {
			expireAt := t.Add(24 * time.Hour).Add(-time.Second)
			msg.ExpireAt = &expireAt
		}
	}

	authorID := ginCtx.GetUint("user_id")
	msg.SenderID = authorID
	msg.SenderName = ginCtx.GetString("username")

	targetUserIDs := req.TargetUserIDs
	if targetUserIDs == nil {
		targetUserIDs = []uint{}
	}

	ctx := context.Background()
	if err := c.notificationService.PublishMessage(ctx, msg, targetUserIDs); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "消息发布成功", gin.H{"id": msg.ID})
}

// GetMessages 获取系统消息列表（管理员）
func (c *MessageController) GetMessages(ginCtx *gin.Context) {
	query := model.SystemMessageQuery{
		Page:        1,
		PageSize:    10,
		Keyword:     ginCtx.Query("keyword"),
		MessageType: ginCtx.Query("message_type"),
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

	messages, total, err := c.notificationService.GetMessages(query)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, gin.H{
		"list":     messages,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetMessageDetail 获取消息详情
func (c *MessageController) GetMessageDetail(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的消息ID")
		return
	}

	msg, err := c.notificationService.GetMessageDetail(uint(id))
	if err != nil {
		response.Error(ginCtx, 404, err.Error())
		return
	}

	response.Success(ginCtx, msg)
}

// DeleteMessage 删除消息
func (c *MessageController) DeleteMessage(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的消息ID")
		return
	}

	if err := c.notificationService.DeleteMessage(uint(id)); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "消息删除成功", nil)
}

// GetMyNotifications 获取我的通知列表
func (c *MessageController) GetMyNotifications(ginCtx *gin.Context) {
	userID := ginCtx.GetUint("user_id")

	query := model.UserNotificationQuery{
		UserID:      userID,
		Page:        1,
		PageSize:    10,
		IsRead:      nil,
		MessageType: ginCtx.Query("message_type"),
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

	isReadStr := ginCtx.Query("is_read")
	if isReadStr != "" {
		if ir, err := strconv.ParseInt(isReadStr, 10, 8); err == nil {
			iri := int8(ir)
			query.IsRead = &iri
		}
	}

	notifications, total, err := c.notificationService.GetMyNotifications(query)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, gin.H{
		"list":     notifications,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

// GetMyNotificationStats 获取我的通知统计
func (c *MessageController) GetMyNotificationStats(ginCtx *gin.Context) {
	userID := ginCtx.GetUint("user_id")

	stats, err := c.notificationService.GetMyNotificationStats(userID)
	if err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.Success(ginCtx, stats)
}

// MarkAsRead 标记通知已读
func (c *MessageController) MarkAsRead(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的通知ID")
		return
	}

	var req model.MarkAsReadRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	if err := c.notificationService.MarkAsRead(uint(id), req.IsRead); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "标记成功", nil)
}

// MarkAllAsRead 全部标记已读
func (c *MessageController) MarkAllAsRead(ginCtx *gin.Context) {
	userID := ginCtx.GetUint("user_id")

	if err := c.notificationService.MarkAllAsRead(userID); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "全部标记成功", nil)
}

// BatchMarkAsRead 批量标记已读
func (c *MessageController) BatchMarkAsRead(ginCtx *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		response.Error(ginCtx, 400, "参数错误: "+err.Error())
		return
	}

	userID := ginCtx.GetUint("user_id")
	for _, id := range req.IDs {
		c.notificationService.MarkAsRead(id, 1)
	}
	_ = userID // 用于后续扩展：记录操作日志

	response.SuccessWithMsg(ginCtx, "批量标记成功", nil)
}

// WithdrawMessage 撤回消息
func (c *MessageController) WithdrawMessage(ginCtx *gin.Context) {
	idStr := ginCtx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		response.Error(ginCtx, 400, "无效的消息ID")
		return
	}

	userID := ginCtx.GetUint("user_id")

	if err := c.notificationService.WithdrawMessage(uint(id), userID); err != nil {
		response.Error(ginCtx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ginCtx, "消息已撤回", nil)
}
