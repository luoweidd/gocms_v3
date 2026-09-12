import request from './index'

export interface TenantCreateReq {
  name: string
  code: string
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  domain?: string
  plan_type: 'free' | 'standard' | 'premium' | 'enterprise'
  max_users?: number
  max_storage?: number
  expire_at?: string
}

export interface TenantUpdateReq {
  name?: string
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  domain?: string
  plan_type?: 'free' | 'standard' | 'premium' | 'enterprise'
  max_users?: number
  max_storage?: number
  expire_at?: string
  status?: number
  settings?: string
}

export interface TenantListItem {
  id: number
  name: string
  code: string
  contact_name: string
  contact_phone: string
  contact_email: string
  domain: string
  plan_type: string
  max_users: number
  max_storage: number
  used_storage: number
  expire_at?: string
  status: number
  user_count: number
  created_at: string
  updated_at: string
}

export interface TenantStats {
  total: number
  active: number
  inactive: number
  expired: number
  by_plan: Record<string, number>
}

export interface TenantUserItem {
  id: number
  username: string
  nickname: string
  email: string
  status: number
  created_at: string
}

/**
 * 获取租户列表
 */
export function getTenantList(params?: any) {
  return request.get<{ list: TenantListItem[]; total: number; page: number; page_size: number }>('/tenants', { params })
}

/**
 * 获取租户统计
 */
export function getTenantStats() {
  return request.get<TenantStats>('/tenants/stats')
}

/**
 * 获取租户详情
 */
export function getTenant(id: number) {
  return request.get<TenantListItem>(`/tenants/${id}`)
}

/**
 * 创建租户
 */
export function createTenant(data: TenantCreateReq) {
  return request.post('/tenants', data)
}

/**
 * 更新租户
 */
export function updateTenant(id: number, data: TenantUpdateReq) {
  return request.put(`/tenants/${id}`, data)
}

/**
 * 删除租户
 */
export function deleteTenant(id: number) {
  return request.delete(`/tenants/${id}`)
}

/**
 * 更新租户状态
 */
export function updateTenantStatus(id: number, status: number) {
  return request.put(`/tenants/${id}/status`, { status })
}

/**
 * 添加用户到租户
 */
export function addUserToTenant(tenantId: number, userId: number) {
  return request.post(`/tenants/${tenantId}/users`, { user_id: userId })
}

/**
 * 获取租户用户列表
 */
export function getTenantUsers(tenantId: number, params?: any) {
  return request.get<{ list: TenantUserItem[]; total: number; page: number; page_size: number }>(`/tenants/${tenantId}/users`, { params })
}

/**
 * 从租户移除用户
 */
export function removeUserFromTenant(tenantId: number, userId: number) {
  return request.delete(`/tenants/${tenantId}/users/${userId}`)
}
