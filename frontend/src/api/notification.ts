import request from '@/utils/request'

// 发布系统消息（管理员）
export function publishMessage(data: any) {
  return request({
    url: '/messages',
    method: 'post',
    data
  })
}

// 获取系统消息列表（管理员）
export function getMessages(params: any) {
  return request({
    url: '/messages',
    method: 'get',
    params
  })
}

// 获取消息详情
export function getMessageDetail(id: number) {
  return request({
    url: `/messages/${id}`,
    method: 'get'
  })
}

// 删除消息（管理员）
export function deleteMessage(id: number) {
  return request({
    url: `/messages/${id}`,
    method: 'delete'
  })
}

// 获取我的通知列表
export function getMyNotifications(params: any) {
  return request({
    url: '/messages/my',
    method: 'get',
    params
  })
}

// 获取我的通知统计
export function getMyNotificationStats() {
  return request({
    url: '/messages/my/stats',
    method: 'get'
  })
}

// 标记通知已读
export function markAsRead(id: number, isRead: number) {
  return request({
    url: `/messages/${id}/read`,
    method: 'put',
    data: { is_read: isRead }
  })
}

// 全部标记已读
export function markAllAsRead() {
  return request({
    url: '/messages/my/all-read',
    method: 'post'
  })
}

// 批量标记已读
export function batchMarkAsRead(ids: number[]) {
  return request({
    url: '/messages/batch-read',
    method: 'post',
    data: { ids }
  })
}

// 撤回消息（管理员）
export function withdrawMessage(id: number) {
  return request({
    url: `/messages/${id}/withdraw`,
    method: 'post'
  })
}

// 消息类型映射
export const MESSAGE_TYPE_MAP: Record<string, string> = {
  system: '系统通知',
  publish: '发布通知',
  comment_reply: '评论回复',
  audit_passed: '审核通过',
  audited_failed: '审核驳回'
}

// 优先级映射
export const PRIORITY_MAP: Record<number, string> = {
  0: '普通',
  1: '重要',
  2: '紧急'
}

// 状态映射
export const MESSAGE_STATUS_MAP: Record<number, string> = {
  1: '已发布',
  0: '草稿',
  2: '已撤回'
}

// ========== 消息类别管理 API ==========

// 获取消息类别列表
export function getMessageCategories(params: any) {
  return request({
    url: '/message-categories',
    method: 'get',
    params
  })
}

// 获取消息类别详情
export function getMessageCategoryDetail(id: number) {
  return request({
    url: `/message-categories/${id}`,
    method: 'get'
  })
}

// 创建消息类别
export function createMessageCategory(data: any) {
  return request({
    url: '/message-categories',
    method: 'post',
    data
  })
}

// 更新消息类别
export function updateMessageCategory(id: number, data: any) {
  return request({
    url: `/message-categories/${id}`,
    method: 'put',
    data
  })
}

// 删除消息类别
export function deleteMessageCategory(id: number) {
  return request({
    url: `/message-categories/${id}`,
    method: 'delete'
  })
}
