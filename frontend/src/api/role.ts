import request from './index'

// ==================== Types ====================

export interface RoleInfo {
  id: number
  name: string
  code: string
  description?: string
  status: number // 1: active, 0: disabled
  created_at: string
  updated_at: string
}

export interface RoleListParams {
  page: number
  page_size: number
  keyword?: string
  status?: number
}

export interface RoleListResponse {
  list: RoleInfo[]
  total: number
  page: number
  page_size: number
}

export interface CreateRoleParams {
  name: string
  code: string
  description?: string
  status?: number
}

export interface UpdateRoleParams {
  name?: string
  code?: string
  description?: string
  status?: number
}

// ==================== Role CRUD APIs ====================

/**
 * 获取角色列表
 */
export function getRoleList(params: RoleListParams) {
  return request.get<RoleListResponse>('/roles', { params })
}

/**
 * 获取角色详情
 */
export function getRole(id: number) {
  return request.get<RoleInfo>(`/roles/${id}`)
}

/**
 * 创建角色
 */
export function createRole(data: CreateRoleParams) {
  return request.post('/roles', data)
}

/**
 * 更新角色信息
 */
export function updateRole(id: number, data: UpdateRoleParams) {
  return request.put(`/roles/${id}`, data)
}

/**
 * 删除角色
 */
export function deleteRole(id: number) {
  return request.delete(`/roles/${id}`)
}

// ==================== Permission APIs ====================

export interface PermissionNode {
  id: number
  name: string
  code: string
  type: 'menu' | 'button'
  parent_id?: number
  children?: PermissionNode[]
}

/**
 * 获取权限树
 */
export function getPermissionTree() {
  return request.get<PermissionNode[]>('/permissions/tree')
}

/**
 * 分配角色权限
 */
export function assignPermissions(roleId: number, permissionIds: number[]) {
  return request.post(`/roles/${roleId}/permissions`, { permission_ids: permissionIds })
}

/**
 * 获取角色拥有的权限
 */
export function getRolePermissions(roleId: number) {
  return request.get<number[]>(`/roles/${roleId}/permissions`)
}
