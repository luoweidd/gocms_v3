import request from './index'

// ==================== Types ====================

export interface UserInfo {
  id: number
  username: string
  email?: string
  phone?: string
  avatar?: string
  role_id?: number
  role_name?: string
  status: number // 1: active, 0: disabled
  roles?: string[]
  created_at: string
  updated_at: string
}

export interface UserListParams {
  page: number
  page_size: number
  keyword?: string
  role_id?: number
  status?: number
}

export interface UserListResponse {
  list: UserInfo[]
  total: number
  page: number
  page_size: number
}

export interface CreateUserParams {
  username: string
  password: string
  email?: string
  phone?: string
  role_id?: number
  roles?: string[]
  status?: number
}

export interface UpdateUserParams {
  username?: string
  email?: string
  phone?: string
  role_id?: number
  status?: number
}

// ==================== User CRUD APIs ====================

/**
 * 获取用户列表
 */
export function getUserList(params: UserListParams) {
  return request.get<UserListResponse>('/users', { params })
}

/**
 * 获取用户详情
 */
export function getUser(id: number) {
  return request.get<UserInfo>(`/users/${id}`)
}

/**
 * 创建用户
 */
export function createUser(data: CreateUserParams) {
  return request.post('/users', data)
}

/**
 * 更新用户信息
 */
export function updateUser(id: number, data: UpdateUserParams) {
  return request.put(`/users/${id}`, data)
}

/**
 * 删除用户
 */
export function deleteUser(id: number) {
  return request.delete(`/users/${id}`)
}

/**
 * 获取角色列表
 */
export function getRoleList(params?: { page?: number; page_size?: number; keyword?: string }) {
  return request.get<any>('/roles', { params })
}

/**
 * 创建角色
 */
export function createRole(data: { name: string; code: string; description?: string; status?: number }) {
  return request.post('/roles', data)
}

/**
 * 更新角色
 */
export function updateRole(id: number, data: { name?: string; code?: string; description?: string; status?: number }) {
  return request.put(`/roles/${id}`, data)
}

/**
 * 删除角色
 */
export function deleteRole(id: number) {
  return request.delete(`/roles/${id}`)
}

/**
 * 分配用户角色
 */
export function assignRoles(userId: number, roleIds: number[]) {
  return request.post(`/users/${userId}/roles`, { role_ids: roleIds })
}

/**
 * 获取用户拥有的角色
 */
export function getUserRoles(userId: number) {
  return request.get<any>(`/users/${userId}/roles`)
}
