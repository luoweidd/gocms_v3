import http from './http'
import { AxiosRequestConfig } from 'axios'

/**
 * 请求包装器（兼容 axios 风格调用）
 * 支持 request({ url, method, data }) 和 request({ url, params }) 等格式
 */
function request(config: AxiosRequestConfig): Promise<any> {
  return http(config)
}

export default request