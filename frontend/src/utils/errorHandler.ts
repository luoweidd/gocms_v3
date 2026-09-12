import { ElMessage } from '@/utils/message'
import axios, { type AxiosError } from 'axios'

// 错误码映射表
const ERROR_CODE_MAP: Record<number, string> = {
  401: '认证失败，请重新登录',
  403: '没有权限访问该资源',
  404: '请求的资源不存在',
  500: '服务器内部错误',
  502: '网关错误',
  503: '服务不可用',
  504: '网关超时'
}

/**
 * 统一错误处理函数
 */
export function handleError(error: AxiosError | Error): void {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError
    const status = axiosError.response?.status
    
    if (status === 401) {
      // 清除本地认证信息
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      // 跳转到登录页
      window.location.href = '/login'
      ElMessage.error('登录已过期，请重新登录')
    } else if (status === 403) {
      ElMessage.error('没有权限执行此操作')
    } else if (status === 404) {
      ElMessage.error('请求的资源不存在')
    } else if (status && status >= 500) {
      ElMessage.error(ERROR_CODE_MAP[status] || '服务器错误，请稍后重试')
    } else {
      ElMessage.error(axiosError.message || '请求失败')
    }
  } else {
    ElMessage.error(error.message || '未知错误')
  }
}

/**
 * 优雅的错误处理包装器
 */
export async function safeCall<T>(
  fn: () => Promise<T>,
  onSuccess?: (data: T) => void,
  onError?: (error: Error) => void
): Promise<boolean> {
  try {
    const result = await fn()
    if (onSuccess) {
      onSuccess(result)
    }
    return true
  } catch (error) {
    const err = error instanceof Error ? error : new Error('未知错误')
    handleError(err)
    if (onError) {
      onError(err)
    }
    return false
  }
}

/**
 * 成功消息提示
 */
export function showSuccess(message: string): void {
  ElMessage.success(message)
}

/**
 * 警告消息提示
 */
export function showWarning(message: string): void {
  ElMessage.warning(message)
}

/**
 * 错误消息提示
 */
export function showError(message: string): void {
  ElMessage.error(message)
}