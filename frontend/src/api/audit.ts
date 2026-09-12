import http from '@/utils/http'

// 所有路径不需要加 /api 前缀，因为 VITE_API_BASE_URL 已经包含 /api
const API_PREFIX = ''

/**
 * 内容审核记录
 */
export interface ContentAudit {
  id: number
  audit_no: string
  content_type: 'article' | 'video' | 'image' | 'file'
  content_id: number
  content_title: string
  submitter_id: number
  submitter_name: string
  auditor_id: number
  auditor_name: string
  status: 0 | 1 | 2 | 3 // 0待审 1通过 2驳回 3需修改
  audit_comment: string
  previous_status: number
  next_status: number
  required_modification: string
  created_at: string
  updated_at: string
  audited_at: string | null
}

/**
 * 审核统计
 */
export interface AuditStats {
  total: number
  pending: number
  passed: number
  rejected: number
  needs_modify: number
}

/**
 * 提交审核请求
 */
export interface SubmitAuditRequest {
  content_type: string
  content_id: number
  content_title?: string
  next_status: number
}

/**
 * 审核列表查询参数
 */
export interface AuditListParams {
  page?: number
  page_size?: number
  content_type?: string
  status?: number
  keyword?: string
  auditor_id?: number
}

/**
 * 批量审核请求
 */
export interface BatchAuditRequest {
  ids: number[]
  action: 'approve' | 'reject'
  audit_comment?: string
  required_modification?: string
}

// ==================== 内容审核 API ====================

/**
 * 提交审核
 */
export function submitAudit(data: SubmitAuditRequest) {
  return http.post<any>(`${API_PREFIX}/audit/submit`, data)
}

/**
 * 获取审核列表
 */
export function getAuditList(params?: AuditListParams) {
  return http.get<any>(`${API_PREFIX}/audit/list`, { params })
}

/**
 * 获取审核详情
 */
export function getAuditById(id: number) {
  return http.get<any>(`${API_PREFIX}/audit/${id}`)
}

/**
 * 通过审核
 */
export function approveAudit(id: number, data?: { audit_comment?: string }) {
  return http.put<any>(`${API_PREFIX}/audit/${id}/approve`, data || {})
}

/**
 * 驳回审核
 */
export function rejectAudit(id: number, data: { audit_comment?: string; required_modification?: string }) {
  return http.put<any>(`${API_PREFIX}/audit/${id}/reject`, data)
}

/**
 * 标记需修改
 */
export function needsModificationAudit(id: number, data: { audit_comment?: string; required_modification?: string }) {
  return http.put<any>(`${API_PREFIX}/audit/${id}/needs-modification`, data)
}

/**
 * 批量审核
 */
export function batchAudit(data: BatchAuditRequest) {
  return http.post<any>(`${API_PREFIX}/audit/batch`, data)
}

/**
 * 获取审核统计
 */
export function getAuditStats() {
  return http.get<any>(`${API_PREFIX}/audit/stats`)
}

// ==================== 操作日志 API ====================

/**
 * 操作日志列表查询参数
 */
export interface OperationLogListParams {
  page?: number
  page_size?: number
  keyword?: string
  action_type?: string
  resource_type?: string
  date?: string
}

/**
 * 获取操作日志列表
 */
export function getOperationLogs(params?: OperationLogListParams) {
  return http.get<any>('/operation-logs' + (params ? '?' + new URLSearchParams(params as any).toString() : ''))
}

// ==================== 操作日志统计 ====================

/**
 * 操作日志统计
 */
export interface OperationLogStats {
  total_count: number
  today_count: number
  success_count: number
  fail_count: number
  module_stats: Array<{ module: string; count: number }>
  action_stats: Array<{ action: string; count: number }>
}

/**
 * 获取操作日志统计
 */
export function getOperationLogStats() {
  return http.get<any>('/operation-logs/stats')
}

// ==================== 工具函数 ====================

/**
 * 获取审核状态标签
 */
export function getAuditStatusLabel(status: number): string {
  const labels: Record<number, string> = {
    0: '待审',
    1: '通过',
    2: '驳回',
    3: '需修改'
  }
  return labels[status] || '未知'
}

/**
 * 获取审核状态样式
 */
export function getAuditStatusClass(status: number): string {
  const classes: Record<number, string> = {
    0: 'bg-yellow-100 text-yellow-800',
    1: 'bg-green-100 text-green-800',
    2: 'bg-red-100 text-red-800',
    3: 'bg-orange-100 text-orange-800'
  }
  return classes[status] || ''
}

/**
 * 获取内容类型标签
 */
export function getContentTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    article: '文章',
    video: '视频',
    image: '图片',
    file: '文件'
  }
  return labels[type] || type
}

/**
 * 格式化日期
 */
export function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}