import request from './index'

// ============================================
// NGAC 权限系统前端类型定义
// ============================================

// ============================================
// 权限 (Permission)
// ============================================

export interface Permission {
  id: number
  name: string
  description?: string
  code: string
  object_class_id: number
  object_class_name?: string
  attribute_rule?: string
  created_at: string
  updated_at: string
}

export interface CreatePermissionRequest {
  name: string
  description?: string
  code: string
  object_class_id: number
  attribute_rule?: string
}

export interface UpdatePermissionRequest {
  name?: string
  description?: string
  code?: string
  object_class_id?: number
  attribute_rule?: string
}

// ============================================
// 角色 (Role)
// ============================================

export interface Role {
  id: number
  name: string
  description?: string
  role_type: string    // builtin, custom
  permission_count?: number
  member_count?: number
  created_at: string
  updated_at: string
}

export interface CreateRoleRequest {
  name: string
  description?: string
  role_type?: string
}

export interface UpdateRoleRequest {
  name?: string
  description?: string
}

// ============================================
// 对象类 (Object Class)
// ============================================

export interface ObjectClass {
  id: number
  name: string
  description?: string
  table_name?: string
  primary_key?: string
  parent_id?: number
  level_path?: string
  children?: ObjectClass[]
}

export interface ObjectClassTree extends ObjectClass {
  children: ObjectClassTree[]
}

// ============================================
// 用户角色 (User Role)
// ============================================

export interface UserRolesResponse {
  user_id: number
  roles: Role[]
}

export interface AssignRolesRequest {
  role_ids: number[]
}

// ============================================
// 用户权限 (User Permission)
// ============================================

export interface UserPermissionsResponse {
  user_id: number
  permissions: Permission[]
  direct_permissions: Permission[]
  role_derived_permissions: Permission[]
}

// ============================================
// 角色权限 (Role Permission)
// ============================================

export interface AssignPermissionsRequest {
  permission_ids: number[]
}

export interface RemovePermissionRequest {
  perm_id: number
}

// ============================================
// 主体属性 (Subject Attribute)
// ============================================

export interface SubjectAttribute {
  subject_id: number
  attribute_name: string
  attribute_value: string
  operator?: string
}

export interface UpdateSubjectAttributesRequest {
  attributes: {
    name: string
    value: string
    operator?: string
  }[]
}

// ============================================
// 上下文属性 (Context Attribute)
// ============================================

export interface ContextAttribute {
  attribute_name: string
  attribute_value: string
  scope?: string       // session, global
}

export interface UpdateContextAttributesRequest {
  attributes: ContextAttribute[]
}

// ============================================
// 授权结果 (Authorization Result)
// ============================================

export interface AuthorizationSubject {
  id: number
  name: string
  attributes: Record<string, string>
}

export interface AuthorizationResource {
  id: number
  object_class: string
  attributes: Record<string, string>
}

export interface AuthorizationContext {
  time?: string
  ip_address?: string
  attributes: Record<string, string>
}

export interface AuthzRequest {
  subject: AuthorizationSubject
  resource: AuthorizationResource
  action: string
  context?: AuthorizationContext
}

export interface AuthzResponse {
  allowed: boolean
  decision: 'Permit' | 'Deny' | 'Indeterminate'
  reason?: string
  evaluated_policies?: string[]
}

// ============================================
// NGAC API 函数
// ============================================

// ============================================
// 权限管理
// ============================================

/**
 * 获取权限列表
 */
export function listPermissions(params?: { object_class_id?: number; keyword?: string }) {
  return request.get<{ list: Permission[]; total: number }>('/ngac/permissions', { params })
}

/**
 * 创建权限
 */
export function createPermission(data: CreatePermissionRequest) {
  return request.post<Permission>('/ngac/permissions', data)
}

/**
 * 更新权限
 */
export function updatePermission(permId: number, data: UpdatePermissionRequest) {
  return request.put<any>(`/ngac/permissions/${permId}`, data)
}

/**
 * 删除权限
 */
export function deletePermission(permId: number) {
  return request.delete<any>(`/ngac/permissions/${permId}`)
}

// ============================================
// 角色管理
// ============================================

/**
 * 创建角色
 */
export function createRole(data: CreateRoleRequest) {
  return request.post<Role>('/ngac/roles', data)
}

/**
 * 获取角色列表
 */
export function listRoles(params?: { role_type?: string; keyword?: string }) {
  return request.get<{ list: Role[]; total: number }>('/ngac/roles', { params })
}

/**
 * 更新角色
 */
export function updateRole(roleId: number, data: UpdateRoleRequest) {
  return request.put<any>(`/ngac/roles/${roleId}`, data)
}

/**
 * 删除角色
 */
export function deleteRole(roleId: number) {
  return request.delete<any>(`/ngac/roles/${roleId}`)
}

// ============================================
// 角色权限管理
// ============================================

/**
 * 为角色分配权限
 */
export function assignPermissionsToRole(roleId: number, data: AssignPermissionsRequest) {
  return request.post<any>(`/ngac/roles/${roleId}/permissions`, data)
}

/**
 * 从角色移除权限
 */
export function removePermissionFromRole(roleId: number, permId: number) {
  return request.delete<any>(`/ngac/roles/${roleId}/permissions/${permId}`)
}

// ============================================
// 用户角色管理
// ============================================

/**
 * 为用户分配角色
 */
export function assignRolesToUser(userId: number, data: AssignRolesRequest) {
  return request.post<any>(`/ngac/users/${userId}/roles`, data)
}

/**
 * 获取用户的角色
 */
export function getUserRoles(userId: number) {
  return request.get<UserRolesResponse>(`/ngac/users/${userId}/roles`)
}

/**
 * 获取用户的权限
 */
export function getUserPermissions(userId: number) {
  return request.get<UserPermissionsResponse>(`/ngac/users/${userId}/permissions`)
}

// ============================================
// 对象类管理
// ============================================

/**
 * 获取对象类树
 */
export function getObjectClassTree() {
  return request.get<ObjectClassTree[]>('/ngac/object-classes/tree')
}

// ============================================
// 属性管理
// ============================================

/**
 * 获取主体属性
 */
export function getSubjectAttributes(userId: number) {
  return request.get<SubjectAttribute[]>(`/ngac/users/${userId}/attributes`)
}

/**
 * 更新主体属性
 */
export function updateSubjectAttributes(userId: number, data: UpdateSubjectAttributesRequest) {
  return request.put<any>(`/ngac/users/${userId}/attributes`, data)
}

/**
 * 获取上下文属性
 */
export function getContextAttributes() {
  return request.get<ContextAttribute[]>('/ngac/context/attributes')
}

/**
 * 更新上下文属性
 */
export function updateContextAttributes(data: UpdateContextAttributesRequest) {
  return request.put<any>('/ngac/context/attributes', data)
}

// ============================================
// 授权接口
// ============================================

/**
 * 执行授权检查
 */
export function authz(data: AuthzRequest) {
  return request.post<AuthzResponse>('/ngac/authz', data)
}
