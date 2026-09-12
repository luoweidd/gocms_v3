package controller

import (
	"strconv"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"
	"gocms_v3/app/service"

	"github.com/gin-gonic/gin"
)

// CommentController 评论控制器
type CommentController struct {
	commentService *service.CommentService
}

// NewCommentController 创建评论控制器
func NewCommentController(svc *service.CommentService) *CommentController {
	return &CommentController{commentService: svc}
}

// Create 创建评论 API: POST /comments
func (ctrl *CommentController) Create(c *gin.Context) {
	var req model.CommentCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	comment, err := ctrl.commentService.CreateComment(req)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "评论成功", comment)
}

// Approve 审核评论 API: PUT /comments/:id/approve
func (ctrl *CommentController) Approve(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	var req struct {
		Status       int    `json:"status" binding:"required,oneof=1 2"`
		AuditRemark  string `json:"audit_remark"`
		ApproverName string `json:"approver_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	// 默认从上下文中获取审批人
	approverID := uint(1)
	approverName := "system"
	if name, exists := c.Get("user_id"); exists {
		approverID = name.(uint)
	}

	_ = approverName // 使用 approverName 参数

	if err := ctrl.commentService.AuditComment(uint(id), req.Status, req.AuditRemark, approverID); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "审核成功", nil)
}

// BatchDelete 批量删除评论 API: DELETE /comments/batch-delete
func (ctrl *CommentController) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if len(req.IDs) == 0 {
		response.Error(c, 400, "请至少提供一个ID")
		return
	}

	if err := ctrl.commentService.BatchDeleteComments(req.IDs); err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.SuccessWithMsg(c, "批量删除成功", nil)
}

// GetList 获取评论列表 API: GET /comments
func (ctrl *CommentController) GetList(c *gin.Context) {
	var query model.CommentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	pageData, err := ctrl.commentService.GetCommentList(query)
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, pageData)
}

// GetStats 获取评论统计 API: GET /comments/stats
func (ctrl *CommentController) GetStats(c *gin.Context) {
	stats := ctrl.commentService.GetCommentStats()
	response.Success(c, stats)
}

// Delete 删除评论 API: DELETE /comments/:id
func (ctrl *CommentController) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		response.Error(c, 400, "参数错误")
		return
	}

	if err := db.GetDB().Delete(&model.Comment{}, uint(id)).Error; err != nil {
		response.Error(c, 500, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

// BatchApprove 批量通过评论 API: POST /comments/batch-approve
func (ctrl *CommentController) BatchApprove(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if len(req.IDs) == 0 {
		response.Error(c, 400, "请至少提供一个ID")
		return
	}

	if err := ctrl.commentService.BatchApprove(req.IDs); err != nil {
		response.Error(c, 500, "批量通过失败")
		return
	}

	response.SuccessWithMsg(c, "批量通过成功", nil)
}

// BatchReject 批量驳回评论 API: POST /comments/batch-reject
func (ctrl *CommentController) BatchReject(c *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids" binding:"required,min=1"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "请求参数错误")
		return
	}

	if len(req.IDs) == 0 {
		response.Error(c, 400, "请至少提供一个ID")
		return
	}

	if err := ctrl.commentService.BatchReject(req.IDs, req.Reason); err != nil {
		response.Error(c, 500, "批量驳回失败")
		return
	}

	response.SuccessWithMsg(c, "批量驳回成功", nil)
}

// GetAuditStats 获取审核统计 API: GET /comments/audit-stats
func (ctrl *CommentController) GetAuditStats(c *gin.Context) {
	stats, err := ctrl.commentService.GetAuditStats()
	if err != nil {
		response.Error(c, 500, err.Error())
		return
	}
	response.Success(c, stats)
}
