import request from './index'

// ============================================
// Video 视频完整类型定义 - 覆盖所有影视属性
// ============================================

export interface VideoTag {
  id: number
  name: string
}

export interface VideoCategory {
  id: number
  name: string
  parent_id?: number
  children?: VideoCategory[]
}

export interface VideoActor {
  id: number
  name: string
  role?: string
  avatar?: string
}

export interface Video {
  id: number
  title: string
  summary?: string        // 后端字段名，与前端 description 语义相同
  description?: string    // 兼容字段
  content?: string        // 视频描述/剧本内容
  cover: string           // 封面图
  cover_image?: string    // 兼容字段
  video_url?: string      // 视频播放URL
  url?: string            // 兼容字段
  
  // 视频元数据
  format?: string         // 视频格式 (mp4, avi, mkv等)
  duration?: number       // 时长（秒）
  file_size?: number      // 文件大小（字节）
  width?: number          // 视频宽度（像素）
  height?: number         // 视频高度（像素）
  bitrate?: number        // 码率
  
  // 影视属性
  year?: string | number  // 年份
  release_date?: string   // 发行日期
  director?: string       // 导演
  actors?: string | VideoActor[]  // 演员
  genre?: string          // 类型/流派
  region?: string         // 地区
  language?: string       // 语言
  rating?: number         // 评分 (0-10)
  total_episodes?: number // 总集数（剧集用）
  current_episode?: number// 当前集数
  
  // 管理字段
  category_id: number
  category_name?: string
  tag_ids?: number[]      // 标签ID列表
  tags?: VideoTag[]       // 标签列表
  author_id?: number
  author_name?: string
  status: 0 | 1 | 2 // 0: 草稿，1: 已发布，2: 已下架
  view_count: number
  like_count?: number     // 点赞数
  is_top?: boolean        // 是否置顶
  is_top_value?: number   // 兼容后端 int 类型
  published_at?: string   // 发布时间
  created_at: string
  updated_at: string
}

export interface VideoListParams {
  page: number
  page_size: number
  category_id?: number
  status?: number
  keyword?: string
  genre?: string          // 类型筛选
  region?: string         // 地区筛选
  language?: string       // 语言筛选
  year?: number | string  // 年份筛选
  director?: string       // 导演筛选
  sort_by?: string        // 排序字段
  sort_order?: 'asc' | 'desc'  // 排序方向
}

export interface VideoListResponse {
  list: Video[]
  total: number
  page: number
  page_size: number
}

// ============================================
// API 函数
// ============================================

/**
 * 获取视频列表
 */
export function getVideoList(params: VideoListParams) {
  return request.get<VideoListResponse>('/videos', { params })
}

/**
 * 获取已发布视频列表（公开接口）
 */
export function getPublishedVideos(params?: VideoListParams) {
  return request.get<VideoListResponse>('/videos/published', { params })
}

/**
 * 获取视频详情
 */
export function getVideo(id: number) {
  return request.get<Video>(`/videos/${id}`)
}

/**
 * 创建视频
 */
export function createVideo(data: Partial<Video>) {
  return request.post('/videos', data)
}

/**
 * 更新视频
 */
export function updateVideo(id: number, data: Partial<Video>) {
  return request.put(`/videos/${id}`, data)
}

/**
 * 删除视频
 */
export function deleteVideo(id: number) {
  return request.delete(`/videos/${id}`)
}

/**
 * 发布视频
 */
export function publishVideo(id: number) {
  return request.post(`/videos/${id}/publish`)
}

/**
 * 下架视频
 */
export function unpublishVideo(id: number) {
  return request.post(`/videos/${id}/unpublish`)
}

/**
 * 置顶视频
 */
export function topVideo(id: number) {
  return request.post(`/videos/${id}/top`)
}

/**
 * 取消置顶视频
 */
export function untopVideo(id: number) {
  return request.post(`/videos/${id}/untop`)
}

/**
 * 视频评分
 */
export function rateVideo(id: number, rating: number) {
  return request.get(`/videos/${id}/rate`, { params: { rating } })
}

/**
 * 获取视频分类树
 */
export function getVideoCategoryTree() {
  return request.get<any>('/video-categories/tree')
}

/**
 * 创建视频分类
 */
export function createVideoCategory(data: { name: string; parent_id?: number }) {
  return request.post('/video-categories', data)
}

/**
 * 更新视频分类
 */
export function updateVideoCategory(id: number, data: { name: string; parent_id?: number }) {
  return request.put(`/video-categories/${id}`, data)
}

/**
 * 删除视频分类
 */
export function deleteVideoCategory(id: number) {
  return request.delete(`/video-categories/${id}`)
}

/**
 * 上传视频文件
 */
export function uploadVideoFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<{ url: string; file_size: number; filename: string }>('/upload/initiate', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 上传封面图
 */
export function uploadCover(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<{ url: string }>('/upload/initiate', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 获取标签列表
 */
export function getTags() {
  return request.get<VideoTag[]>('/tags')
}
