<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">媒体中心</h1>
        <p class="text-sm text-gray-500">统一管理图片、视频、文件和音频资源</p>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-4 mb-8">
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">总资源</p>
          <svg class="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total_assets || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">{{ formatFileSizeFn(stats.total_size || 0) }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$router.push('/media/images')">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">图片</p>
          <svg class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2m0 2a2 2 0 100 4m0-4a2 2 0 110 4m12 0V2m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-5 9a2 2 0 10-4 0 2 2 0 004 0zm-7 4a2 2 0 10-4 0 2 2 0 004 0zm10-8a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-blue-600">{{ stats.image_count || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">点击图片进入</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$router.push('/videos/list')">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">视频</p>
          <svg class="w-5 h-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-red-600">{{ stats.video_count || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">点击进入视频管理</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$router.push('/media/audio')">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">音频</p>
          <svg class="w-5 h-5 text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-purple-600">{{ stats.audio_count || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">点击音频进入</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$router.push('/media/files')">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">文件</p>
          <svg class="w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-gray-600">{{ stats.file_count || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">点击进入文件</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow cursor-pointer" @click="$router.push('/media/images')">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">相册</p>
          <svg class="w-5 h-5 text-indigo-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-indigo-600">{{ stats.total_albums || 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">管理相册</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 hover:shadow-md transition-shadow">
        <div class="flex items-center justify-between mb-2">
          <p class="text-sm text-gray-500">存储使用率</p>
          <svg class="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
        </div>
        <p class="text-2xl font-bold text-emerald-600">{{ storagePercentage }}%</p>
        <div class="w-full bg-gray-200 rounded-full h-1.5 mt-2">
          <div class="bg-emerald-500 h-1.5 rounded-full" :style="{ width: storagePercentage + '%' }"></div>
        </div>
      </div>
    </div>

    <!-- Quick Access Cards -->
    <div class="mb-8">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">快速访问</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg shadow-sm p-6 text-white cursor-pointer hover:shadow-lg transition-shadow" @click="$router.push('/media/images')">
          <div class="flex items-center justify-between mb-4">
            <svg class="w-10 h-10 text-white/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2m0 2a2 2 0 100 4m0-4a2 2 0 110 4m12 0V2m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-5 9a2 2 0 10-4 0 2 2 0 004 0zm-7 4a2 2 0 10-4 0 2 2 0 004 0zm10-8a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-1">图片管理</h3>
          <p class="text-sm text-white/70">浏览、搜索和管理所有图片资源</p>
        </div>
        <div class="bg-gradient-to-br from-red-500 to-red-600 rounded-lg shadow-sm p-6 text-white cursor-pointer hover:shadow-lg transition-shadow" @click="$router.push('/videos/list')">
          <div class="flex items-center justify-between mb-4">
            <svg class="w-10 h-10 text-white/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-1">视频管理</h3>
          <p class="text-sm text-white/70">管理视频内容、分类和发布状态</p>
        </div>
        <div class="bg-gradient-to-br from-gray-500 to-gray-600 rounded-lg shadow-sm p-6 text-white cursor-pointer hover:shadow-lg transition-shadow" @click="$router.push('/media/files')">
          <div class="flex items-center justify-between mb-4">
            <svg class="w-10 h-10 text-white/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-1">文件管理</h3>
          <p class="text-sm text-white/70">管理文档、压缩包等其他文件</p>
        </div>
        <div class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg shadow-sm p-6 text-white cursor-pointer hover:shadow-lg transition-shadow" @click="$router.push('/media/audio')">
          <div class="flex items-center justify-between mb-4">
            <svg class="w-10 h-10 text-white/80" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
          </div>
          <h3 class="text-lg font-semibold mb-1">音频管理</h3>
          <p class="text-sm text-white/70">管理音乐、语音等音频资源</p>
        </div>
      </div>
    </div>

    <!-- Recent Assets -->
    <div class="bg-white rounded-lg shadow-sm">
      <div class="flex items-center justify-between p-6 border-b border-gray-100">
        <h2 class="text-lg font-semibold text-gray-900">最近上传</h2>
        <button @click="$router.push('/media/images')" class="text-sm text-indigo-600 hover:text-indigo-700">查看全部 →</button>
      </div>
      <div class="p-6">
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div v-for="item in recentAssets" :key="item.id" class="group">
            <div class="relative aspect-square rounded-lg overflow-hidden bg-gray-100 mb-2">
              <img v-if="item.thumbnail || item.file_type === 'image'" :src="getImageUrl(item.thumbnail || item.file_path)" class="w-full h-full object-cover" />
              <div v-else class="w-full h-full flex items-center justify-center" :class="getFileTypeBg(item.file_type)">
                <svg class="w-8 h-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
            </div>
            <p class="text-xs text-gray-600 truncate">{{ item.name }}</p>
            <p class="text-xs text-gray-400">{{ formatDate(item.created_at) }}</p>
          </div>
        </div>
        <div v-if="recentAssets.length === 0" class="py-8 text-center text-sm text-gray-500">暂无最近上传资源</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getMediaStats, formatFileSize as apiFormatFileSize } from '@/api/media'

interface AssetItem {
  id: number
  name: string
  thumbnail?: string
  file_path: string
  file_type: string
  created_at: string
}

const stats = ref<any>({ 
  total_assets: 0, 
  image_count: 0, 
  video_count: 0, 
  audio_count: 0, 
  file_count: 0, 
  total_albums: 0, 
  total_size: 0 
})

const recentAssets = ref<AssetItem[]>([])
const formatFileSizeFn = (bytes: number) => apiFormatFileSize(bytes)
const formatDate = (dateStr: string) => { try { return new Date(dateStr).toLocaleDateString('zh-CN') } catch { return dateStr } }
const getImageUrl = (filePath: string) => `/api${filePath}`

const getFileTypeBg = (fileType: string) => {
  const classes: Record<string, string> = { image: 'bg-blue-50', video: 'bg-red-50', audio: 'bg-purple-50', file: 'bg-gray-50' }
  return classes[fileType] || classes.file
}

const storagePercentage = computed(() => {
  const maxStorage = 10737418240 // 10GB
  const used = stats.value.total_size || 0
  return Math.min(100, Math.round((used / maxStorage) * 100))
})

const loadStats = async () => {
  try {
    const res = await getMediaStats()
    stats.value = res.data || {}
  } catch (error) {
    console.error('获取统计失败', error)
  }
}

onMounted(() => { loadStats() })
</script>