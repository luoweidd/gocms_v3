import http from '@/utils/http'

// 所有路径不需要加 /api 前缀，因为 VITE_API_BASE_URL 已经包含 /api
const API_PREFIX = ''

/**
 * 媒体资源信息
 */
export interface MediaAsset {
  id: number
  user_id: number
  upload_id: string
  name: string
  original_name: string
  file_path: string
  file_type: 'file' | 'image' | 'video' | 'audio'
  mime_type: string
  file_size: number
  width?: number
  height?: number
  duration?: number
  thumbnail: string
  album_id?: number
  tags: string[]
  description: string
  status: number
  created_at: string
  updated_at: string
}

/**
 * 相册信息
 */
export interface Album {
  id: number
  name: string
  description: string
  user_id: number
  cover_image: string
  asset_count: number
  is_public: number
  status: number
  created_at: string
  updated_at: string
}

/**
 * 媒体统计信息
 */
export interface MediaStats {
  total_albums: number
  total_assets: number
  total_size: number
  image_count: number
  video_count: number
  audio_count: number
  file_count: number
}

/**
 * 创建媒体资源请求
 */
export interface CreateMediaAssetRequest {
  upload_id: string
  name: string
  original_name?: string
  file_path: string
  file_type: string
  mime_type?: string
  file_size: number
  width?: number
  height?: number
  duration?: number
  thumbnail?: string
  description?: string
  tags?: string[]
}

/**
 * 更新媒体资源请求
 */
export interface UpdateMediaAssetRequest {
  name?: string
  description?: string
  tags?: string[]
  album_id?: number
  status?: number
}

/**
 * 媒体资源列表查询参数
 */
export interface MediaAssetListParams {
  page?: number
  page_size?: number
  user_id?: number
  file_type?: string
  album_id?: number
  keyword?: string
  status?: number
}

/**
 * 创建相册请求
 */
export interface CreateAlbumRequest {
  name: string
  desc?: string
  cover_image?: string
  is_public?: number
}

/**
 * 更新相册请求
 */
export interface UpdateAlbumRequest {
  name?: string
  desc?: string
  cover_image?: string
  is_public?: number
}

/**
 * 相册列表查询参数
 */
export interface AlbumListParams {
  page?: number
  page_size?: number
  user_id?: number
}

// ==================== 媒体资源 API ====================

/**
 * 创建媒体资源
 */
export function createMediaAsset(data: CreateMediaAssetRequest) {
  return http.post<any>(`${API_PREFIX}/media/assets`, data)
}

/**
 * 获取媒体资源列表
 */
export function getMediaAssetList(params?: MediaAssetListParams) {
  return http.get<any>(`${API_PREFIX}/media/assets`, { params })
}

/**
 * 获取媒体资源详情
 */
export function getMediaAssetById(id: number) {
  return http.get<any>(`${API_PREFIX}/media/assets/${id}`)
}

/**
 * 更新媒体资源
 */
export function updateMediaAsset(id: number, data: UpdateMediaAssetRequest) {
  return http.put<any>(`${API_PREFIX}/media/assets/${id}`, data)
}

/**
 * 删除媒体资源
 */
export function deleteMediaAsset(id: number) {
  return http.delete<any>(`${API_PREFIX}/media/assets/${id}`)
}

/**
 * 批量删除媒体资源
 */
export function batchDeleteMediaAssets(ids: number[]) {
  return http.delete<any>(`${API_PREFIX}/media/assets/batch`, { data: { ids } })
}

/**
 * 获取媒体统计信息
 */
export function getMediaStats(params?: { user_id?: number }) {
  return http.get<any>(`${API_PREFIX}/media/assets/stats`, { params })
}

// ==================== 相册 API ====================

/**
 * 创建相册
 */
export function createAlbum(data: CreateAlbumRequest) {
  return http.post<any>(`${API_PREFIX}/media/albums`, data)
}

/**
 * 获取相册列表
 */
export function getAlbumList(params?: AlbumListParams & { user_id?: number }) {
  return http.get<any>(`${API_PREFIX}/media/albums`, { params })
}

/**
 * 获取相册详情（含资源列表）
 */
export function getAlbumById(id: number) {
  return http.get<any>(`${API_PREFIX}/media/albums/${id}`)
}

/**
 * 更新相册
 */
export function updateAlbum(id: number, data: UpdateAlbumRequest) {
  return http.put<any>(`${API_PREFIX}/media/albums/${id}`, data)
}

/**
 * 删除相册
 */
export function deleteAlbum(id: number) {
  return http.delete<any>(`${API_PREFIX}/media/albums/${id}`)
}

/**
 * 添加资源到相册
 */
export function addAssetToAlbum(albumId: number, assetIds: number[]) {
  return http.post<any>(`${API_PREFIX}/media/albums/${albumId}/assets`, { asset_ids: assetIds })
}

/**
 * 从相册移除资源
 */
export function removeAssetFromAlbum(albumId: number, assetId: number) {
  return http.delete<any>(`${API_PREFIX}/media/albums/${albumId}/assets/${assetId}`)
}

// ==================== 工具函数 ====================

/**
 * 格式化文件大小
 */
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const rank = Math.floor(Math.log(bytes) / Math.log(1024))
  return parseFloat((bytes / Math.pow(1024, rank)).toFixed(1)) + ' ' + units[rank]
}

/**
 * 格式化视频/音频时长（秒）
 */
export function formatDuration(seconds: number): string {
  if (!seconds) return '00:00'
  const mins = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
}

/**
 * 获取文件类型显示名称
 */
export function getFileTypeLabel(fileType: string): string {
  const labels: Record<string, string> = {
    image: '图片',
    video: '视频',
    audio: '音频',
    file: '文件'
  }
  return labels[fileType] || fileType
}

/**
 * 获取MIME类型图标
 */
export function getFileIcon(mimeType: string): string {
  if (mimeType.startsWith('image/')) return 'icon-image'
  if (mimeType.startsWith('video/')) return 'icon-video'
  if (mimeType.startsWith('audio/')) return 'icon-audio'
  return 'icon-file'
}