import request from './index'

export interface DashboardStats {
  total_users: number
  today_new_users: number
  total_articles: number
  today_new_articles: number
  total_videos: number
  today_new_videos: number
  total_comments: number
  today_new_comments: number
  article_status: Array<{ status: string; count: number }>
  video_status: Array<{ status: string; count: number }>
  comment_status: Array<{ status: string; count: number }>
  user_status: Array<{ status: string; count: number }>
  user_role_dist: Array<{ role: string; count: number }>
  article_categories: Array<{ category_id: number; category_name: string; count: number }>
  video_categories: Array<{ category_id: number; category_name: string; count: number }>
  tag_distribution: Array<{ tag_id: number; tag_name: string; count: number }>
  article_trend: Array<{ date: string; count: number }>
  video_trend: Array<{ date: string; count: number }>
  user_trend: Array<{ date: string; count: number }>
  comment_trend: Array<{ date: string; count: number }>
  recent_articles: Array<{ id: number; title: string; type: string; author: string; status: number; created_at: string }>
  recent_videos: Array<{ id: number; title: string; type: string; author: string; status: number; created_at: string }>
  recent_comments: Array<{ id: number; title: string; type: string; author: string; status: number; created_at: string }>
  recent_users: Array<{ id: number; title: string; type: string; author: string; status: number; created_at: string }>
  top_articles: Array<{ id: number; title: string; views: number; likes: number; comments: number; tags: string[]; created_at: string }>
  top_videos: Array<{ id: number; title: string; views: number; likes: number; comments: number; tags: string[]; created_at: string }>
  top_comments: Array<{ id: number; content: string; author: string; likes: number; url: string; created_at: string }>
  active_authors: Array<{ author: string; count: number }>
  resource_usage: {
    cpu: { cores: number; usage_percent: number }
    memory: { total: number; used: number; free: number; usage_percent: number }
    disk: { total: number; used: number; free: number; usage_percent: number }
    network: { up_speed: number; down_speed: number; up_total: number; down_total: number }
    uptime: number
    goroutines: number
  } | null
}

export interface DashboardOverview {
  total_users: number
  total_articles: number
  total_videos: number
  total_comments: number
  today_new_users: number
  today_new_articles: number
  today_new_videos: number
  today_new_comments: number
  user_growth_rate: number
  article_growth_rate: number
  video_growth_rate: number
  comment_growth_rate: number
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  qps: number
}

export interface StatTrend {
  article: Array<{ date: string; count: number }>
  video: Array<{ date: string; count: number }>
  comment: Array<{ date: string; count: number }>
  user: Array<{ date: string; count: number }>
}

export interface ActivityItem {
  id: number
  title: string
  type: string
  status: number
  created_at: string
}

export interface ResourceUsage {
  cpu: { cores: number; usage_percent: number }
  memory: { total: number; used: number; free: number; usage_percent: number }
  disk: { total: number; used: number; free: number; usage_percent: number }
  network: { up_speed: number; down_speed: number; up_total: number; down_total: number }
  uptime: number
  goroutines: number
}

/**
 * 获取仪表盘综合统计数据
 */
export function getDashboardStats() {
  return request.get<DashboardStats>('/dashboard/stats')
}

/**
 * 获取统计趋势数据
 */
export function getStatTrend(days: number = 7) {
  return request.get<StatTrend>(`/dashboard/trend?days=${days}`)
}

/**
 * 获取最近活动
 */
export function getRecentActivities(limit: number = 10) {
  return request.get<ActivityItem[]>(`/dashboard/recent-activity?limit=${limit}`)
}

/**
 * 获取概览数据
 */
export function getDashboardOverview() {
  return request.get<DashboardOverview>('/dashboard/overview')
}

/**
 * 获取硬件资源使用情况
 */
export function getResourceUsage() {
  return request.get<ResourceUsage>('/dashboard/resource-usage')
}
