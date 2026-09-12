import request from '@/utils/request'

// 创建SEO设置
export function createSeoSetting(data: any) {
  return request({
    url: '/seo/settings',
    method: 'post',
    data
  })
}

// 更新SEO设置
export function updateSeoSetting(id: number, data: any) {
  return request({
    url: `/seo/settings/${id}`,
    method: 'put',
    data
  })
}

// 获取SEO列表
export function getSeoList(params: any) {
  return request({
    url: '/seo/settings',
    method: 'get',
    params
  })
}

// 获取SEO详情
export function getSeoDetail(id: number) {
  return request({
    url: `/seo/settings/${id}`,
    method: 'get'
  })
}

// 删除SEO设置
export function deleteSeoSetting(id: number) {
  return request({
    url: `/seo/settings/${id}`,
    method: 'delete'
  })
}

// 获取站点地图配置
export function getSitemapConfigs() {
  return request({
    url: '/seo/sitemap-configs',
    method: 'get'
  })
}

// 生成站点地图
export function generateSitemap() {
  return request({
    url: '/seo/sitemap/generate',
    method: 'post',
    responseType: 'text' // XML响应需要设置为text类型
  })
}

// 获取站点地图统计
export function getSitemapStats() {
  return request({
    url: '/seo/sitemap-stats',
    method: 'get'
  })
}

// 资源类型映射
export const RESOURCE_TYPE_MAP: Record<string, string> = {
  article: '文章',
  video: '视频',
  page: '页面'
}

// Robots指令映射
export const ROBOTS_MAP: Record<string, string> = {
  'index,follow': '索引+跟随',
  'noindex,nofollow': '不索引+不跟随',
  'index,nofollow': '索引+不跟随',
  'noindex,follow': '不索引+跟随'
}

// 更新频率映射
export const CHANGE_FREQ_MAP: Record<string, string> = {
  always: '每次',
  hourly: '每小时',
  daily: '每天',
  weekly: '每周',
  monthly: '每月',
  yearly: '每年',
  never: '从不'
}