package controller

import (
	"strconv"

	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// AuditController 内容审核控制器
type AuditController struct {
	AuditService *service.ContentAuditService
}

// NewAuditController 创建内容审核控制器
func NewAuditController() *AuditController {
	return &AuditController{
		AuditService: service.NewContentAuditService(),
	}
}

// SubmitAudit 提交审核
func (c *AuditController) SubmitAudit(ctx *gin.Context) {
	var req struct {
		ContentType  string `json:"content_type" binding:"required"`
		ContentID    uint   `json:"content_id" binding:"required"`
		ContentTitle string `json:"content_title"`
		NextStatus   int8   `json:"next_status" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	userName := ctx.GetString("username")

	audit, err := c.AuditService.SubmitAudit(struct {
		ContentType   string
		ContentID     uint
		ContentTitle  string
		SubmitterID   uint
		SubmitterName string
		NextStatus    int8
	}{
		ContentType:   req.ContentType,
		ContentID:     req.ContentID,
		ContentTitle:  req.ContentTitle,
		SubmitterID:   userID,
		SubmitterName: userName,
		NextStatus:    req.NextStatus,
	})
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, audit)
}

// GetAuditList 获取审核列表
func (c *AuditController) GetAuditList(ctx *gin.Context) {
	var req model.ContentAuditListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 200 {
		req.PageSize = 20
	}

	list, total, err := c.AuditService.GetAuditList(req)
	if err != nil {
		response.Error(ctx, 500, "查询失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{"list": list, "total": total})
}

// GetAuditByID 获取审核详情
func (c *AuditController) GetAuditByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	audit, err := c.AuditService.GetAuditByID(uint(id))
	if err != nil {
		response.Error(ctx, 404, err.Error())
		return
	}

	response.Success(ctx, audit)
}

// Approve 通过审核
func (c *AuditController) Approve(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req struct {
		AuditComment string `json:"audit_comment"`
	}
	ctx.ShouldBindJSON(&req)

	userID := ctx.GetUint("user_id")
	userName := ctx.GetString("username")

	if err := c.AuditService.Approve(uint(id), userID, userName, req.AuditComment); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "审核通过", nil)
}

// Reject 驳回审核
func (c *AuditController) Reject(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req struct {
		AuditComment         string `json:"audit_comment"`
		RequiredModification string `json:"required_modification"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	userName := ctx.GetString("username")

	if err := c.AuditService.Reject(uint(id), userID, userName, req.AuditComment, req.RequiredModification); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "审核驳回", nil)
}

// NeedsModification 标记需修改
func (c *AuditController) NeedsModification(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(ctx, 400, "参数错误")
		return
	}

	var req struct {
		AuditComment         string `json:"audit_comment"`
		RequiredModification string `json:"required_modification"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	userID := ctx.GetUint("user_id")
	userName := ctx.GetString("username")

	if err := c.AuditService.NeedsModification(uint(id), userID, userName, req.AuditComment, req.RequiredModification); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "已标记需修改", nil)
}

// BatchAudit 批量审核
func (c *AuditController) BatchAudit(ctx *gin.Context) {
	var req struct {
		IDs                  []uint `json:"ids" binding:"required"`
		Action               string `json:"action" binding:"required"` // approve/reject
		AuditComment         string `json:"audit_comment"`
		RequiredModification string `json:"required_modification"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "参数错误: "+err.Error())
		return
	}

	if len(req.IDs) == 0 {
		response.Error(ctx, 400, "请选择要审核的记录")
		return
	}

	userID := ctx.GetUint("user_id")
	userName := ctx.GetString("username")

	if err := c.AuditService.BatchAudit(userID, userName, req.IDs, req.Action, req.AuditComment, req.RequiredModification); err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.SuccessWithMsg(ctx, "批量审核成功", nil)
}

// GetAuditStats 获取审核统计
func (c *AuditController) GetAuditStats(ctx *gin.Context) {
	stats, err := c.AuditService.GetAuditStats()
	if err != nil {
		response.Error(ctx, 500, "查询失败: "+err.Error())
		return
	}

	response.Success(ctx, stats)
}
