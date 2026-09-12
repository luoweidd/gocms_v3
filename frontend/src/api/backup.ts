import request from '@/utils/request'

// 创建备份
export function createBackup(data: { name?: string; backup_type: string }) {
  return request({
    url: '/backup/full',
    method: 'post',
    data
  })
}

// 获取备份列表
export function getBackupList(params: any) {
  return request({
    url: '/backup/list',
    method: 'get',
    params
  })
}

// 获取备份详情
export function getBackupDetail(id: number) {
  return request({
    url: '/backup/' + id,
    method: 'get'
  })
}

// 删除备份
export function deleteBackup(id: number) {
  return request({
    url: '/backup/' + id,
    method: 'delete'
  })
}

// 下载备份文件
export function downloadBackup(id: number) {
  return request({
    url: '/backup/' + id + '/download',
    method: 'get',
    responseType: 'blob'
  })
}

// 恢复备份
export function restoreBackup(id: number) {
  return request({
    url: '/backup/' + id + '/restore',
    method: 'post'
  })
}

// 获取备份统计
export function getBackupStats() {
  return request({
    url: '/backup/stats',
    method: 'get'
  })
}

// 备份类型映射
export const BACKUP_TYPE_MAP: Record<string, string> = {
  full: '全量备份',
  incremental: '增量备份'
}

// 状态映射
export const BACKUP_STATUS_MAP: Record<number, string> = {
  0: '创建中',
  1: '成功',
  2: '失败'
}

// 格式化文件大小
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}