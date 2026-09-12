import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { message } from './message'

// ==================== 请求/响应类型定义 ====================

interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

interface RequestConfig extends AxiosRequestConfig {
  silent?: boolean // 是否静默错误（不显示错误提示）
}

// ==================== HTTP 实例创建 ====================

function createHttpInstance(): AxiosInstance {
  const instance: AxiosInstance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 15000,
    headers: {
      'Content-Type': 'application/json',
    },
  })

  // ==================== 请求拦截器 ====================

  instance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      // 添加 Token
      const token = localStorage.getItem('token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }

      // 添加请求时间戳（用于防重复提交和超时计算）
      ;(config as any).metadata = { startTime: Date.now() }

      return config
    },
    (error) => {
      console.error('Request interceptor error:', error)
      return Promise.reject(error)
    }
  )

  // ==================== 响应拦截器 ====================

    instance.interceptors.response.use(
       (response: AxiosResponse<ApiResponse>) => {
         // 如果是 blob 响应（文件下载），直接返回
         if (response.config.responseType === 'blob' || response.data instanceof Blob) {
           return response
         }

         // 如果是 text 响应（如XML内容），直接返回原始数据
         if (response.config.responseType === 'text') {
           return response
         }

         const { code, message: msg, data } = response.data

         if (code === 0 || code === 200) {
           return { ...response, data: data }
         }

       // 业务错误处理
       const errorMessage = msg || '请求失败'

       // 401 未授权
       if (code === 401) {
         localStorage.removeItem('token')
         localStorage.removeItem('user')
         window.location.href = '/login'
         message.error('登录已过期，请重新登录')
         return Promise.reject(new Error(errorMessage))
       }

       // 403 权限不足
       if (code === 403) {
         message.warning('权限不足，无法执行此操作')
         return Promise.reject(new Error(errorMessage))
       }

       // 其他业务错误
       const config = response.config as RequestConfig & { metadata?: { startTime: number } }
       if (!config?.silent) {
         message.error(errorMessage)
       }

       return Promise.reject(new Error(errorMessage))
     },
    (error) => {
      let messageText = '网络错误，请稍后重试'

      if (error.response) {
        switch (error.response.status) {
          case 400:
            messageText = '请求参数错误'
            break
          case 401:
            messageText = '未授权，请重新登录'
            localStorage.removeItem('token')
            window.location.href = '/login'
            break
          case 403:
            messageText = '拒绝访问'
            break
          case 404:
            messageText = '请求的资源不存在'
            break
          case 500:
            messageText = '服务器错误'
            break
          case 502:
            messageText = '网关错误'
            break
          case 503:
            messageText = '服务不可用'
            break
          case 504:
            messageText = '网关超时'
            break
          default:
            messageText = `连接错误 (${error.response.status})`
        }
      } else if (error.code === 'ERR_CANCELED') {
        messageText = '请求已取消'
      } else {
        messageText = error.message || '未知错误'
      }

      const config = error.config as RequestConfig & { metadata?: { startTime: number } }
      if (!config?.silent) {
        message.error(messageText)
      }

      return Promise.reject(error)
    }
  )

  return instance
}

// ==================== 创建实例并导出 ====================

const http = createHttpInstance()

// ==================== 便捷方法 ====================

/**
 * 静默 GET 请求（不显示错误提示）
 */
export function silentGet(url: string, config?: Record<string, any>): Promise<AxiosResponse<any>> {
  return http.get(url, { ...config, silent: true } as AxiosRequestConfig)
}

/**
 * 静默 POST 请求
 */
export function silentPost(url: string, data?: any, config?: Record<string, any>): Promise<AxiosResponse<any>> {
  return http.post(url, data, { ...config, silent: true } as AxiosRequestConfig)
}

// ==================== 防重复提交 ====================

const pendingRequests = new Map<string, AbortController>()

/**
 * 生成请求唯一标识
 */
function getRequestKey(config: InternalAxiosRequestConfig): string {
  const { url, method, data } = config
  return `${method}_${url}_${JSON.stringify(data)}`
}

/**
 * 取消挂起的相同请求
 */
function cancelDuplicateRequest(config: InternalAxiosRequestConfig): void {
  const key = getRequestKey(config)
  if (pendingRequests.has(key)) {
    pendingRequests.get(key)?.abort()
    pendingRequests.delete(key)
  }
}

/**
 * 添加请求到挂起列表
 */
// eslint-disable-next-line @typescript-eslint/no-unused-vars
function addPendingRequest(_config: InternalAxiosRequestConfig): void {
  cancelDuplicateRequest(_config)
  const controller = new AbortController()
  _config.signal = controller.signal
  pendingRequests.set(getRequestKey(_config), controller)
}

/**
 * 移除挂起的请求
 */
function removePendingRequest(config: InternalAxiosRequestConfig): void {
  const key = getRequestKey(config)
  pendingRequests.delete(key)
}

// 重新设置拦截器（包含防重复提交）
http.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  cancelDuplicateRequest(config)
  const controller = new AbortController()
  config.signal = controller.signal

  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  return config
})

http.interceptors.response.use(
  (response) => {
    removePendingRequest(response.config as InternalAxiosRequestConfig)
    return response
  },
  (error) => {
    if (error.config) {
      removePendingRequest(error.config as InternalAxiosRequestConfig)
    }
    return Promise.reject(error)
  }
)

export default http

export type { AxiosInstance }
export { http as defaultHttp, createHttpInstance }
