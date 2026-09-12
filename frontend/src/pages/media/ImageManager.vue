<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">图片管理</h1>
        <p class="text-sm text-gray-500">管理系统中的所有图片资源</p>
      </div>
      <div class="flex gap-3">
        <TailButton type="primary" @click="showUploadDialog = true">
          <svg class="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2V7a2 2 0 00-2-2h-1a2 2 0 01-2-2V3a1 1 0 00-1-1H9a1 1 0 00-1 1v2a2 2 0 01-2 2H4a2 2 0 00-2 2v2" />
          </svg>
          上传图片
        </TailButton>
        <TailButton type="success" @click="showImportDialog = true">
          <svg class="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2V7a2 2 0 00-2-2h-1a2 2 0 01-2-2V3a1 1 0 00-1-1H9a1 1 0 00-1 1v2a2 2 0 01-2 2H4a2 2 0 00-2 2v2" />
          </svg>
          导入已有文件
        </TailButton>
        <TailButton type="info" @click="showCreateAlbumDialog = true">
          <svg class="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          新建相册
        </TailButton>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总图片</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.image_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总视频</p>
        <p class="text-2xl font-bold text-purple-600">{{ stats.video_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总相册</p>
        <p class="text-2xl font-bold text-indigo-600">{{ stats.total_albums || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总资源</p>
        <p class="text-2xl font-bold text-gray-600">{{ stats.total_assets || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总大小</p>
        <p class="text-lg font-bold text-emerald-600">{{ formatFileSizeFn(stats.total_size || 0) }}</p>
      </div>
    </div>

    <!-- Main Content -->
    <div class="bg-white rounded-lg shadow-sm">
      <!-- Search and Filter bar -->
      <div class="p-4 border-b border-gray-100">
        <div class="flex gap-3 flex-wrap items-center">
          <div class="relative flex-1 min-w-[200px] max-w-md">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input v-model="searchParams.keyword" type="text" placeholder="搜索图片..." class="ta-input pl-10 w-full" @input="debounceSearch" />
          </div>
          
          <select v-model="searchParams.album_id" class="ta-input w-auto min-w-[150px]" @change="handleSearch">
            <option value="">全部相册</option>
            <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
          </select>

          <div class="flex items-center gap-2 ml-auto">
            <button @click="viewMode = 'grid'" class="p-2 rounded-lg transition-colors" :class="viewMode === 'grid' ? 'bg-indigo-100 text-indigo-600' : 'text-gray-400 hover:text-gray-600'">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
              </svg>
            </button>
            <button @click="viewMode = 'list'" class="p-2 rounded-lg transition-colors" :class="viewMode === 'list' ? 'bg-indigo-100 text-indigo-600' : 'text-gray-400 hover:text-gray-600'">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Grid View -->
      <div v-if="viewMode === 'grid'" class="p-4">
        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
          <div
            v-for="asset in assets"
            :key="asset.id"
            class="group relative bg-white rounded-lg border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
            :class="{ 'ring-2 ring-indigo-500': selectedAssets.includes(asset.id) }"
          >
            <div class="relative aspect-square bg-gray-100 overflow-hidden">
              <img
                v-if="asset.thumbnail || asset.file_path"
                :src="getImageUrl(asset.thumbnail || asset.file_path)"
                :alt="asset.name"
                class="w-full h-full object-cover"
                @error="handleImageError"
              />
              <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                <button @click="previewAsset(asset)" class="p-2 bg-white rounded-full hover:bg-gray-100" title="预览">
                  <svg class="w-4 h-4 text-gray-700" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
                <button @click="deleteAsset(asset.id)" class="p-2 bg-white rounded-full hover:bg-red-100" title="删除">
                  <svg class="w-4 h-4 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            </div>
            <div class="p-2">
              <p class="text-xs font-medium text-gray-900 truncate">{{ asset.name }}</p>
              <p class="text-xs text-gray-500">{{ formatFileSizeFn(asset.file_size) }}</p>
            </div>
          </div>
        </div>

        <div v-if="assets.length === 0 && !loading" class="py-16 text-center">
          <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无图片资源</h3>
          <p class="mt-1 text-sm text-gray-500">上传第一个图片文件开始使用</p>
        </div>

        <!-- Pagination -->
        <div class="flex items-center justify-between px-4 py-3 border-t border-gray-200" v-if="assets.length > 0">
          <span class="text-sm text-gray-500">共 {{ pagination.total }} 张图片</span>
          <div class="flex gap-1">
            <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page = 1">首页</button>
            <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page--">上一页</button>
            <button v-for="p in visiblePages" :key="p" class="w-8 h-8 text-sm border rounded" :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : 'hover:bg-gray-50'" @click="pagination.page = p; handleSearch()">{{ p }}</button>
            <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page >= totalPages" @click="pagination.page++">下一页</button>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left w-10"><input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-4 h-4 rounded" /></th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">文件</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">大小</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">相册</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="asset in assets" :key="asset.id" class="hover:bg-gray-50">
              <td class="px-4 py-3"><input type="checkbox" :checked="selectedAssets.includes(asset.id)" @change="toggleSelect(asset.id)" class="w-4 h-4 rounded" /></td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded bg-gray-100 flex items-center justify-center overflow-hidden">
                    <img v-if="asset.thumbnail" :src="getImageUrl(asset.thumbnail)" class="w-full h-full object-cover" />
                  </div>
                  <div class="min-w-0"><p class="text-sm font-medium text-gray-900 truncate">{{ asset.name }}</p><p class="text-xs text-gray-500">{{ asset.original_name }}</p></div>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ formatFileSizeFn(asset.file_size) }}</td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ getAlbumName(asset.album_id) }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(asset.created_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button @click="previewAsset(asset)" class="p-1 text-gray-400 hover:text-indigo-600" title="预览">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
                <button @click="deleteAsset(asset.id)" class="p-1 text-gray-400 hover:text-red-600" title="删除">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="assets.length === 0 && !loading" class="py-16 text-center">
          <p class="text-sm text-gray-500">暂无图片资源</p>
        </div>
      </div>
    </div>

    <!-- Batch Operations Bar -->
    <div v-if="selectedAssets.length > 0" class="fixed bottom-6 left-1/2 -translate-x-1/2 bg-gray-900 text-white px-6 py-3 rounded-full shadow-lg flex items-center gap-4 z-40">
      <span class="text-sm">{{ selectedAssets.length }} 项已选</span>
      <button @click="batchAddToAlbum" class="text-sm hover:text-indigo-300">添加到相册</button>
      <button @click="batchDelete" class="text-sm hover:text-red-300">删除</button>
      <button @click="selectedAssets = []" class="text-sm text-gray-400 hover:text-white">取消</button>
    </div>

    <!-- Upload Dialog -->
    <Teleport to="body">
      <div v-if="showUploadDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showUploadDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">上传图片</h3>
            <button @click="showUploadDialog = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div
            ref="dropZoneRef"
            class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors"
            :class="[isDragOver ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400 hover:bg-gray-50']"
            @click="triggerFileInput"
            @dragover.prevent.stop="handleDragOver"
            @dragleave.prevent.stop="handleDragLeave"
            @drop.prevent.stop="handleDrop"
          >
            <input ref="fileInputRef" type="file" multiple accept="image/*" class="hidden" @change="handleFileSelect" />
            <svg class="mx-auto w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p class="mt-2 text-sm text-gray-600"><span class="text-indigo-600 font-medium">点击选择</span> 或拖放图片到此处</p>
            <p class="mt-1 text-xs text-gray-500">支持 JPG、PNG、GIF、WebP 等格式</p>
          </div>
          <div class="flex justify-end gap-3 mt-4 pt-4 border-t border-gray-200">
            <button @click="showUploadDialog = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 text-gray-700">取消</button>
            <button @click="handleUpload" :disabled="uploading" class="px-4 py-2 text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50">
              {{ uploading ? '上传中...' : '上传' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Import Dialog -->
    <Teleport to="body">
      <div v-if="showImportDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showImportDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">导入已有文件</h3>
            <button @click="showImportDialog = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">JSON 配置</label>
            <textarea v-model="importJson" rows="8" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 font-mono text-xs" placeholder='[{"name": "file.jpg", "file_path": "/storage/uploads/file.jpg", "file_type": "image"}]'></textarea>
            <p class="mt-1 text-xs text-gray-500">格式：[{"name": 文件名, "file_path": 完整路径, "file_type": "image/video/file", "file_size": 文件大小(字节)}]</p>
          </div>
          <div class="flex justify-end gap-3 mt-4 pt-4 border-t border-gray-200">
            <button @click="showImportDialog = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 text-gray-700">取消</button>
            <button @click="handleImport" :disabled="importing" class="px-4 py-2 text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50">
              {{ importing ? '导入中...' : '导入' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Create Album Dialog -->
    <Teleport to="body">
      <div v-if="showCreateAlbumDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showCreateAlbumDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6 m-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">新建相册</h3>
          <div class="space-y-4">
            <div><label class="block text-sm font-medium text-gray-700 mb-1">相册名称</label><input v-model="albumForm.name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="请输入相册名称" /></div>
            <div><label class="block text-sm font-medium text-gray-700 mb-1">描述</label><textarea v-model="albumForm.desc" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="相册描述"></textarea></div>
            <div><label class="block text-sm font-medium text-gray-700 mb-1">可见性</label><select v-model="albumForm.is_public" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500"><option :value="0">私有</option><option :value="1">公开</option></select></div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <button @click="showCreateAlbumDialog = false" class="px-4 py-2 text-sm border rounded-lg hover:bg-gray-50">取消</button>
            <button @click="handleCreateAlbum" class="px-4 py-2 text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700">创建</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import {
  getMediaAssetList,
  getMediaStats,
  deleteMediaAsset,
  batchDeleteMediaAssets,
  getAlbumList,
  createAlbum,
  formatFileSize as apiFormatFileSize,
} from '@/api/media'

const viewMode = ref<'grid' | 'list'>('grid')
const assets = ref<any[]>([])
const albums = ref<any[]>([])
const stats = ref<any>({ image_count: 0, total_albums: 0, total_size: 0 })
const loading = ref(false)

const pagination = reactive({ page: 1, pageSize: 18, total: 0 })
const searchParams = reactive({ keyword: '', album_id: undefined as number | undefined })
const selectedAssets = ref<number[]>([])

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, pagination.page - 2)
  const end = Math.min(totalPages.value, pagination.page + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const isAllSelected = computed(() => assets.value.length > 0 && assets.value.every((a: any) => selectedAssets.value.includes(a.id)))
const formatFileSizeFn = (bytes: number) => apiFormatFileSize(bytes)

const formatDate = (dateStr: string) => { try { return new Date(dateStr).toLocaleDateString('zh-CN') } catch { return dateStr } }
// 修复图片URL：通过 /api 代理访问后端静态文件（Vite proxy 自动处理 URL 解码）
const getImageUrl = (filePath: string): string => {
  if (!filePath) return ''
  // base64 或完整 URL 直接返回
  if (filePath.startsWith('data:') || filePath.startsWith('http://') || filePath.startsWith('https://')) {
    return filePath
  }
  // file_path 格式为 /storage/uploads/filename.jpg
  // 通过 /api/storage/uploads/... 代理到 Go 后端
  if (filePath.startsWith('/storage/uploads/')) {
    return '/api' + filePath
  }
  // 其他情况直接返回
  return filePath
}
const handleImageError = (e: Event) => { 
  ;(e.target as HTMLImageElement).style.display = 'none' 
}
const getAlbumName = (albumId?: number) => {
  if (!albumId) return '未分类'
  const album = albums.value.find((a: any) => a.id === albumId)
  return album?.name || '未分类'
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => { if (searchTimer) clearTimeout(searchTimer); searchTimer = setTimeout(() => { pagination.page = 1; loadAssets() }, 300) }
const handleSearch = () => { pagination.page = 1; loadAssets() }

// 从 localStorage 获取当前用户ID
const getCurrentUserId = () => {
  const userStr = localStorage.getItem('userInfo')
  if (userStr) {
    try {
      const user = JSON.parse(userStr)
      return user.id || user.user_id || 0
    } catch {
      return 0
    }
  }
  return 0
}

const loadStats = async () => { 
   try { 
     const userId = getCurrentUserId()
     
     const params: any = {}
     if (userId > 0) params.user_id = userId
     
      const res = await getMediaStats(params) 
      
      // API 响应结构: { code: 0, message: "success", data: { ... } }
      // 响应拦截器已将 data 直接放到 res.data
      const responseData = res.data || {}
     stats.value = {
       image_count: responseData.image_count || 0,
       video_count: responseData.video_count || 0,
       total_albums: responseData.total_albums || 0,
       total_assets: responseData.total_assets || 0,
       total_size: responseData.total_size || 0,
       audio_count: responseData.audio_count || 0,
       file_count: responseData.file_count || 0,
       pdf_count: responseData.pdf_count || 0,
       office_count: responseData.office_count || 0,
     }
   } catch (error) { console.error('获取统计失败', error) } 
 }

const loadAlbums = async () => { 
   try { 
     const userId = getCurrentUserId()
     
     const params: any = { page: 1, page_size: 100 }
     if (userId > 0) {
       params.user_id = userId
     }
     
      const res = await getAlbumList(params) 
      
      // API 响应结构: { code: 0, message: "success", data: { list: [...], total: X } }
      // 响应拦截器已将 data 直接放到 res.data
      const responseData = res.data || {}
     albums.value = responseData.list || responseData.albums || []
   } catch (error) { 
     console.error('获取相册列表失败', error) 
   } 
 }

const loadAssets = async () => {
   loading.value = true
   try {
     const userId = getCurrentUserId()
     
     const params: any = { page: pagination.page, page_size: pagination.pageSize, file_type: 'image' }
     if (userId > 0) params.user_id = userId
     if (searchParams.keyword) params.keyword = searchParams.keyword
     if (searchParams.album_id) params.album_id = searchParams.album_id
     
      const res = await getMediaAssetList(params)
      
      // API 响应结构: { code: 0, message: "success", data: { list: [...], total: X } }
      // 响应拦截器已将 data 直接放到 res.data
      const responseData = res.data || {}
     assets.value = responseData.list || []
     pagination.total = responseData.total || 0
   } catch (error) { 
     console.error('获取图片列表失败', error)
     message.error('获取图片列表失败: ' + (error as Error).message) 
   } finally { 
     loading.value = false 
   }
 }

const toggleSelect = (id: number) => { const index = selectedAssets.value.indexOf(id); if (index > -1) selectedAssets.value.splice(index, 1); else selectedAssets.value.push(id) }
const toggleSelectAll = () => { if (isAllSelected.value) selectedAssets.value = []; else selectedAssets.value = assets.value.map((a: any) => a.id) }
const previewAsset = (asset: any) => { message.info(`预览: ${asset.name}`) }
const deleteAsset = async (id: number) => {
  if (!confirm('确定删除此文件吗？')) return
  try { await deleteMediaAsset(id); message.success('删除成功'); loadAssets(); loadStats() } catch (error: any) { message.error(error.message || '删除失败') }
}
const batchDelete = async () => {
  if (!confirm(`确定删除选中的 ${selectedAssets.value.length} 个文件吗？`)) return
  try { await batchDeleteMediaAssets(selectedAssets.value); message.success('批量删除成功'); selectedAssets.value = []; loadAssets(); loadStats() } catch (error: any) { message.error(error.message || '批量删除失败') }
}
const batchAddToAlbum = () => { message.info(`添加到相册: ${selectedAssets.value.length} 个文件`) }

// Dialog states
const showUploadDialog = ref(false)
const showCreateAlbumDialog = ref(false)
const showImportDialog = ref(false)
const importJson = ref('[{"name": "WIN_20250510_01_27_49_Pro.jpg", "file_path": "/storage/uploads/WIN_20250510_01_27_49_Pro.jpg", "file_type": "image", "file_size": 0}]')
const importing = ref(false)
const albumForm = reactive({ name: '', desc: '', is_public: 0 })
const dropZoneRef = ref<HTMLDivElement | null>(null)
const fileInputRef = ref<HTMLInputElement | null>(null)
const isDragOver = ref(false)
const uploadFiles = ref<any[]>([])
const uploading = ref(false)

const handleDragOver = (e: DragEvent) => { isDragOver.value = true; e.stopPropagation() }
const handleDragLeave = (e: DragEvent) => { isDragOver.value = false; e.stopPropagation() }
const handleDrop = (e: DragEvent) => {
  isDragOver.value = false
  e.stopPropagation()
  if (e.dataTransfer?.files?.length) {
    uploadFiles.value = Array.from(e.dataTransfer.files).map(f => ({ file: f, name: f.name, size: f.size, type: f.type, status: 'pending' as const, progress: 0 }))
  }
}
const triggerFileInput = () => { fileInputRef.value?.click() }
const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files?.length) {
    uploadFiles.value = Array.from(target.files).map(f => ({ file: f, name: f.name, size: f.size, type: f.type, status: 'pending' as const, progress: 0 }))
    target.value = ''
  }
}

const handleUpload = async () => {
  if (uploadFiles.value.length === 0) { message.warning('请选择要上传的文件'); return }
  uploading.value = true
  let successCount = 0
  for (let i = 0; i < uploadFiles.value.length; i++) {
    const fileItem = uploadFiles.value[i]
    uploadFiles.value[i].status = 'uploading'
    uploadFiles.value[i].progress = 0
    try {
      const token = localStorage.getItem('token')
      const fileType = fileItem.file.type.startsWith('image/') ? 'image' : 'file'
      
      // 步骤1：初始化上传
      const initiateFormData = new FormData()
      initiateFormData.append('filename', fileItem.file.name)
      initiateFormData.append('size', String(fileItem.file.size))
      initiateFormData.append('chunk_num', '0')
      initiateFormData.append('file_type', fileType)
      
      const initiateRes = await fetch('/api/upload/initiate', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: initiateFormData
      })
      
      if (!initiateRes.ok) {
        const errResult = await initiateRes.json().catch(() => ({}))
        throw new Error(errResult.message || '初始化上传失败')
      }
      
      const initiateResult = await initiateRes.json()
      const uploadID = initiateResult.data.upload_id
      uploadFiles.value[i].progress = 33
      
      // 步骤2：上传分片
      const chunkFormData = new FormData()
      chunkFormData.append('chunk', fileItem.file)
      chunkFormData.append('upload_id', uploadID)
      chunkFormData.append('chunk_num', '0')
      
      const chunkRes = await fetch('/api/upload/chunk', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: chunkFormData
      })
      
      if (!chunkRes.ok) throw new Error('上传分片失败')
      uploadFiles.value[i].progress = 66
      
      // 步骤3：完成上传
      const completeRes = await fetch('/api/upload/complete', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          upload_id: uploadID,
          total_size: fileItem.file.size,
          filename: fileItem.file.name,
          file_type: fileType
        })
      })
      
      if (!completeRes.ok) throw new Error('完成上传失败')
      const completeResult = await completeRes.json()
      
      // 更新文件信息
      uploadFiles.value[i].status = 'success'
      uploadFiles.value[i].progress = 100
      uploadFiles.value[i].url = completeResult.data.url
      successCount++
    } catch (error: any) {
      uploadFiles.value[i].status = 'error'
      uploadFiles.value[i].error = error.message || '上传失败'
    }
  }
  uploading.value = false
  
  // 刷新数据（无论是否成功都要刷新）
  try {
    await Promise.all([loadAssets(), loadStats()])
  } catch (error) {
    console.error('刷新数据失败:', error)
  }
  
  if (successCount > 0) {
    message.success(`成功上传 ${successCount} 个文件`)
  }
  
  showUploadDialog.value = false
  uploadFiles.value = []
}

const handleImport = async () => {
  if (!importJson.value.trim()) {
    message.warning('请输入 JSON 配置')
    return
  }
  importing.value = true
  try {
    const files = JSON.parse(importJson.value)
    let successCount = 0
    
    for (const item of files) {
      try {
        // 调用后端 API 创建 MediaAsset 记录（不上传文件，只注册）
        const token = localStorage.getItem('token')
        const res = await fetch('/api/upload/register', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(item)
        })
        
        if (res.ok) {
          successCount++
        }
      } catch (err) {
        // 静默失败
      }
    }
    
    message.success(`成功导入 ${successCount}/${files.length} 个文件`)
    
    // 刷新数据（无论是否成功都要刷新）
    try {
      await Promise.all([loadAssets(), loadStats()])
    } catch (error) {
      console.error('刷新数据失败:', error)
    }
    
    showImportDialog.value = false
  } catch (error: any) {
    message.error('JSON 解析失败: ' + error.message)
  } finally {
    importing.value = false
  }
}

const handleCreateAlbum = async () => {
  if (!albumForm.name) { message.warning('请填写相册名称'); return }
  try {
    await createAlbum({ name: albumForm.name, desc: albumForm.desc, is_public: albumForm.is_public })
    message.success('创建成功')
    showCreateAlbumDialog.value = false
    albumForm.name = ''
    albumForm.desc = ''
    albumForm.is_public = 0
    loadAlbums()
    loadStats()
  } catch (error: any) { message.error(error.message || '创建失败') }
}

onMounted(() => { loadAssets(); loadAlbums(); loadStats() })
</script>