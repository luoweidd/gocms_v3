<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">文件管理</h1>
        <p class="text-sm text-gray-500">管理系统中的文档和其他文件资源</p>
      </div>
      <div class="flex gap-3">
        <TailButton type="primary" @click="showUploadDialog = true">
          <svg class="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2V7a2 2 0 00-2-2h-1a2 2 0 01-2-2V3a1 1 0 00-1-1H9a1 1 0 00-1 1v2a2 2 0 01-2 2H4a2 2 0 00-2 2v2" />
          </svg>
          上传文件
        </TailButton>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总文件</p>
        <p class="text-2xl font-bold text-gray-600">{{ stats.file_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">PDF</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.pdf_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">Office</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.office_count || 0 }}</p>
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
            <input v-model="searchParams.keyword" type="text" placeholder="搜索文件..." class="ta-input pl-10 w-full" @input="debounceSearch" />
          </div>
          
          <select v-model="searchParams.file_ext" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
            <option value="">全部类型</option>
            <option value=".pdf">PDF</option>
            <option value=".doc">Word</option>
            <option value=".docx">DOCX</option>
            <option value=".xls">Excel</option>
            <option value=".xlsx">XLSX</option>
            <option value=".zip">ZIP</option>
            <option value=".rar">RAR</option>
          </select>

          <select v-model="searchParams.sort_by" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
            <option value="created_at">最新上传</option>
            <option value="size">文件大小</option>
            <option value="name">文件名</option>
          </select>
        </div>
      </div>

      <!-- List View -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left w-10"><input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-4 h-4 rounded" /></th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">文件名</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">类型</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">大小</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">上传者</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="file in files" :key="file.id" class="hover:bg-gray-50">
              <td class="px-4 py-3"><input type="checkbox" :checked="selectedFiles.includes(file.id)" @change="toggleSelect(file.id)" class="w-4 h-4 rounded" /></td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded flex items-center justify-center overflow-hidden" :class="getFileIconBg(getFileExtension(file))">
                    <svg class="w-5 h-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate">{{ file.name }}</p>
                    <p class="text-xs text-gray-500">{{ file.original_name }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span :class="getFileExtensionClass(getFileExtension(file))" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium">
                  {{ getFileExtension(file) || 'FILE' }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ formatFileSizeFn(file.file_size) }}</td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ file.uploader || '系统' }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(file.created_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button @click="previewFile(file)" class="p-1 text-gray-400 hover:text-indigo-600 mr-2" title="预览/下载">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                </button>
                <button @click="deleteFile(file.id)" class="p-1 text-gray-400 hover:text-red-600" title="删除">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="files.length === 0 && !loading" class="py-16 text-center">
          <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无文件</h3>
          <p class="mt-1 text-sm text-gray-500">上传第一个文件开始使用</p>
        </div>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-4 py-3 border-t border-gray-200" v-if="files.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 个文件</span>
        <div class="flex gap-1">
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page = 1">首页</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page--">上一页</button>
          <button v-for="p in visiblePages" :key="p" class="w-8 h-8 text-sm border rounded" :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : 'hover:bg-gray-50'" @click="pagination.page = p; handleSearch()">{{ p }}</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page >= totalPages" @click="pagination.page++">下一页</button>
        </div>
      </div>
    </div>

    <!-- Batch Operations Bar -->
    <div v-if="selectedFiles.length > 0" class="fixed bottom-6 left-1/2 -translate-x-1/2 bg-gray-900 text-white px-6 py-3 rounded-full shadow-lg flex items-center gap-4 z-40">
      <span class="text-sm">{{ selectedFiles.length }} 项已选</span>
      <button @click="batchDownload" class="text-sm hover:text-indigo-300">批量下载</button>
      <button @click="batchDelete" class="text-sm hover:text-red-300">删除</button>
      <button @click="selectedFiles = []" class="text-sm text-gray-400 hover:text-white">取消</button>
    </div>

    <!-- Upload Dialog -->
    <Teleport to="body">
      <div v-if="showUploadDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showUploadDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">上传文件</h3>
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
            <input ref="fileInputRef" type="file" multiple class="hidden" @change="handleFileSelect" />
            <svg class="mx-auto w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p class="mt-2 text-sm text-gray-600"><span class="text-indigo-600 font-medium">点击选择</span> 或拖放文件到此处</p>
            <p class="mt-1 text-xs text-gray-500">支持 PDF、Word、Excel、ZIP 等格式</p>
          </div>
          <!-- Selected files preview -->
          <div v-if="uploadFiles.length > 0" class="mt-4 max-h-40 overflow-y-auto space-y-2">
            <div v-for="(f, i) in uploadFiles" :key="i" class="flex items-center justify-between text-sm p-2 rounded bg-gray-50">
              <div class="flex items-center gap-2 min-w-0 flex-1">
                <svg v-if="f.status === 'success'" class="w-4 h-4 text-green-500 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <svg v-else-if="f.status === 'uploading'" class="w-4 h-4 text-indigo-500 flex-shrink-0 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span class="truncate">{{ f.name }}</span>
                <span class="text-gray-400 flex-shrink-0">{{ formatFileSize(f.size) }}</span>
              </div>
              <button v-if="f.status === 'pending'" @click="uploadFiles.splice(i, 1)" class="text-red-500 hover:text-red-700 ml-2 flex-shrink-0">
                &times;
              </button>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-4 pt-4 border-t border-gray-200">
            <button @click="showUploadDialog = false" class="px-4 py-2 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 text-gray-700">取消</button>
            <button @click="handleUpload" :disabled="uploading || uploadFiles.length === 0" class="px-4 py-2 text-sm bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed">
              {{ uploading ? '上传中...' : '上传 (' + uploadFiles.filter(f => f.status === 'pending').length + ')' }}
            </button>
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
  formatFileSize as apiFormatFileSize,
} from '@/api/media'

interface FileItem {
  id: number
  name: string
  original_name: string
  extension: string
  file_size: number
  file_path: string
  uploader: string
  created_at: string
}

const files = ref<FileItem[]>([])
const loading = ref(false)
const selectedFiles = ref<number[]>([])

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const searchParams = reactive({ keyword: '', file_ext: '', sort_by: 'created_at' })

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, pagination.page - 2)
  const end = Math.min(totalPages.value, pagination.page + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const isAllSelected = computed(() => files.value.length > 0 && files.value.every((f: any) => selectedFiles.value.includes(f.id)))
const formatFileSizeFn = (bytes: number) => apiFormatFileSize(bytes)
const formatDate = (dateStr: string) => { try { return new Date(dateStr).toLocaleDateString('zh-CN') } catch { return dateStr } }

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => { if (searchTimer) clearTimeout(searchTimer); searchTimer = setTimeout(() => { pagination.page = 1; loadFiles() }, 300) }
const handleSearch = () => { pagination.page = 1; loadFiles() }

const getFileIconBg = (ext: string) => {
  const map: Record<string, string> = {
    '.pdf': 'bg-red-50',
    '.doc': 'bg-blue-50',
    '.docx': 'bg-blue-50',
    '.xls': 'bg-green-50',
    '.xlsx': 'bg-green-50',
    '.zip': 'bg-yellow-50',
    '.rar': 'bg-yellow-50',
  }
  return map[ext] || 'bg-gray-50'
}

const getFileExtensionClass = (ext: string) => {
  const map: Record<string, string> = {
    '.pdf': 'bg-red-100 text-red-800',
    '.doc': 'bg-blue-100 text-blue-800',
    '.docx': 'bg-blue-100 text-blue-800',
    '.xls': 'bg-green-100 text-green-800',
    '.xlsx': 'bg-green-100 text-green-800',
    '.zip': 'bg-yellow-100 text-yellow-800',
    '.rar': 'bg-yellow-100 text-yellow-800',
  }
  return map[ext] || 'bg-gray-100 text-gray-800'
}

// 从文件名中提取扩展名
const getFileExtension = (item: any): string => {
  const name = item.name || item.original_name || ''
  const ext = name.substring(name.lastIndexOf('.')).toLowerCase()
  return ext
}

const loadFiles = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.pageSize }
    if (searchParams.keyword) params.keyword = searchParams.keyword
    if (searchParams.file_ext) params.file_ext = searchParams.file_ext
    // 过滤掉图片和视频，只显示文件类型
    const res = await getMediaAssetList(params)
    // 在后端返回的结果中过滤，只显示非图片非视频的文档文件
    files.value = (res.data?.list || []).filter((item: any) => {
      const fileType = item.file_type || ''
      const mimeType = item.mime_type || ''
      const name = item.name || item.original_name || ''
      
      // 排除图片和视频
      if (fileType === 'image' || fileType === 'video' || fileType === 'audio') {
        return false
      }
      
      // 检查文件扩展名
      const ext = name.substring(name.lastIndexOf('.')).toLowerCase()
      const allowedExts = ['.pdf', '.doc', '.docx', '.xls', '.xlsx', '.zip', '.rar']
      if (searchParams.file_ext) {
        return ext === searchParams.file_ext.toLowerCase()
      }
      return allowedExts.includes(ext) || ext === ''
    })
    pagination.total = res.data?.total || 0
  } catch (error) { message.error('获取文件列表失败') } finally { loading.value = false }
}

const toggleSelect = (id: number) => { const index = selectedFiles.value.indexOf(id); if (index > -1) selectedFiles.value.splice(index, 1); else selectedFiles.value.push(id) }
const toggleSelectAll = () => { if (isAllSelected.value) selectedFiles.value = []; else selectedFiles.value = files.value.map((f: any) => f.id) }

const previewFile = (file: FileItem) => {
  window.open(`/api${file.file_path}`, '_blank')
}

const deleteFile = async (id: number) => {
  if (!confirm('确定删除此文件吗？')) return
  try { await deleteMediaAsset(id); message.success('删除成功'); loadFiles() } catch (error: any) { message.error(error.message || '删除失败') }
}

const batchDelete = async () => {
  if (!confirm(`确定删除选中的 ${selectedFiles.value.length} 个文件吗？`)) return
  try { await batchDeleteMediaAssets(selectedFiles.value); message.success('批量删除成功'); selectedFiles.value = []; loadFiles() } catch (error: any) { message.error(error.message || '批量删除失败') }
}

const batchDownload = () => {
  message.info(`批量下载: ${selectedFiles.value.length} 个文件`)
}

// Upload dialog states
const showUploadDialog = ref(false)
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
  const pendingFiles = uploadFiles.value.filter(f => f.status === 'pending')
  if (pendingFiles.length === 0) { message.warning('请选择要上传的文件'); return }
  
  uploading.value = true
  let successCount = 0

  for (let i = 0; i < pendingFiles.length; i++) {
    const fileItem = pendingFiles[i]
    uploadFiles.value.find(f => f === fileItem)!.status = 'uploading'
    
    try {
      const token = localStorage.getItem('token')
      
      // 步骤1：初始化上传
      const initiateFormData = new FormData()
      initiateFormData.append('filename', fileItem.file.name)
      initiateFormData.append('size', String(fileItem.file.size))
      initiateFormData.append('chunk_num', '0')
      // 根据文件类型设置 file_type
      const fileType = fileItem.file.type.startsWith('image/') ? 'image' 
        : fileItem.file.type.startsWith('video/') ? 'video'
          : fileItem.file.type.startsWith('audio/') ? 'audio' : 'file'
      initiateFormData.append('file_type', fileType)
      
      const initiateRes = await fetch('/api/upload/initiate', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: initiateFormData
      })
      
      if (!initiateRes.ok) throw new Error('初始化上传失败')
      const initiateResult = await initiateRes.json()
      const uploadID = initiateResult.data.upload_id
      
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
      
      // 更新状态
      const idx = uploadFiles.value.findIndex(f => f === fileItem)
      if (idx !== -1) {
        uploadFiles.value[idx].status = 'success'
        uploadFiles.value[idx].url = completeResult.data.url
      }
      successCount++
    } catch (error: any) {
      console.error('上传失败:', error)
      const idx = uploadFiles.value.findIndex(f => f === fileItem)
      if (idx !== -1) {
        uploadFiles.value[idx].status = 'error'
      }
      message.error(`${fileItem.file.name}: ${error.message || '上传失败'}`)
    }
  }
  
  uploading.value = false
  if (successCount > 0) {
    message.success(`成功上传 ${successCount} 个文件`)
    showUploadDialog.value = false
    uploadFiles.value = []
    loadFiles()
  }
}

// 格式化文件大小
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const rank = Math.floor(Math.log(bytes) / Math.log(1024))
  return parseFloat((bytes / Math.pow(1024, rank)).toFixed(1)) + ' ' + units[rank]
}

const stats = ref<any>({ file_count: 0, pdf_count: 0, office_count: 0, total_size: 0 })

const loadStats = async () => {
  try {
    const res = await getMediaStats()
    stats.value = res.data || { file_count: 0, pdf_count: 0, office_count: 0, total_size: 0 }
  } catch (error) {
    console.error('获取统计失败', error)
    stats.value = { file_count: 0, pdf_count: 0, office_count: 0, total_size: 0 }
  }
}

onMounted(() => { loadFiles(); loadStats() })
</script>