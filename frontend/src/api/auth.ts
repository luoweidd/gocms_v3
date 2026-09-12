import request from './index'

export interface LoginParams {
  username: string
  password: string
}

export interface RegisterParams {
  username: string
  password: string
  email?: string
}

export interface UserInfo {
  id: number
  username: string
  nickname?: string
  email?: string
  avatar?: string
  roles?: string[]
  phone?: string
  real_name?: string
  bio?: string
  role_name?: string
}

/**
 * 用户登录
 */
export function login(data: LoginParams) {
  return request.post('/auth/login', data)
}

/**
 * 用户注册
 */
export function register(data: RegisterParams) {
  return request.post('/auth/register', data)
}

/**
 * 获取当前用户信息
 */
export function getCurrentUser() {
  return request.get<UserInfo>('/user/info')
}

/**
 * 退出登录
 */
export function logout() {
  return request.post('/auth/logout')
}
