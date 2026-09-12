import request from './index'

export interface CommentInfo {
  id: number
  content: string
  user_id: number
  username?: string
  article_id?: number
  video_id?: number
  status: number // 0: pending, 1: approved, 2: rejected, 3: deleted
  ip?: string
  created_at: string
  updated_at: string
}

export interface CommentListParams {
  page: number
  page_size: number
  article_id?: number
  video_id?: number
  status?: number
  keyword?: string
  user_id?: number
}

export interface CommentListResponse {
  list: CommentInfo[]
  total: number
  page: number
  page_size: number
}

export interface CommentAuditStats {
  total: number
  pending: number
  approved: number
  rejected: number
  deleted: number
  by_article: number
  by_video: number
}

/**
 * 获取评论列表
 */
export function getCommentList(params: CommentListParams) {
  return request.get<CommentListResponse>('/comments', { params })
}

/**
 * 获取评论详情
 */
export function getComment(id: number) {
  return request.get<CommentInfo>(`/comments/${id}`)
}

/**
 * 通过/拒绝评论
 */
export function approveComment(id: number, approved: boolean) {
  return request.put(`/comments/${id}/approve`, { approved })
}

/**
 * 删除评论
 */
export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}

/**
 * 批量删除评论
 */
export function batchDeleteComments(ids: number[]) {
  return request.post('/comments/batch-delete', { ids })
}

/**
 * 批量通过评论
 */
export function batchApproveComments(ids: number[]) {
  return request.post('/comments/batch-approve', { ids })
}

/**
 * 批量驳回评论
 */
export function batchRejectComments(ids: number[], reason: string) {
  return request.post('/comments/batch-reject', { ids, reason })
}

/**
 * 获取审核统计
 */
export function getCommentAuditStats() {
  return request.get<CommentAuditStats>('/comments/audit-stats')
}
