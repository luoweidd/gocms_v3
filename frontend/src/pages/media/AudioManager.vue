<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">音频管理</h1>
        <p class="text-sm text-gray-500">管理系统中的音乐和音频资源</p>
      </div>
      <div class="flex gap-3">
        <TailButton type="primary" @click="showUploadDialog = true">
          <svg class="w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2V7a2 2 0 00-2-2h-1a2 2 0 01-2-2V3a1 1 0 00-1-1H9a1 1 0 00-1 1v2a2 2 0 01-2 2H4a2 2 0 00-2 2v2" />
          </svg>
          上传音频
        </TailButton>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总音频</p>
        <p class="text-2xl font-bold text-purple-600">{{ stats.audio_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">MP3</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.mp3_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4">
        <p class="text-sm text-gray-500">总大小</p>
        <p class="text-lg font-bold text-emerald-600">{{ formatFileSizeFn(audioStats.total_size || 0) }}</p>
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
            <input v-model="searchParams.keyword" type="text" placeholder="搜索音频..." class="ta-input pl-10 w-full" @input="debounceSearch" />
          </div>
          
          <select v-model="searchParams.audio_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
            <option value="">全部类型</option>
            <option value="mp3">MP3</option>
            <option value="wav">WAV</option>
            <option value="ogg">OGG</option>
            <option value="aac">AAC</option>
          </select>
        </div>
      </div>

      <!-- List View -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left w-10"><input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-4 h-4 rounded" /></th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">音频文件</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">类型</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">时长</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">大小</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="audio in audios" :key="audio.id" class="hover:bg-gray-50">
              <td class="px-4 py-3"><input type="checkbox" :checked="selectedAudios.includes(audio.id)" @change="toggleSelect(audio.id)" class="w-4 h-4 rounded" /></td>
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded bg-purple-50 flex items-center justify-center">
                    <svg class="w-5 h-5 text-purple-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
                    </svg>
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate">{{ audio.name }}</p>
                    <p class="text-xs text-gray-500">{{ audio.original_name }}</p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3">
                <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-purple-100 text-purple-800">
                  {{ audio.extension || 'AUDIO' }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ formatDuration(audio.duration) }}</td>
              <td class="px-4 py-3 text-sm text-gray-600">{{ formatFileSizeFn(audio.file_size) }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(audio.created_at) }}</td>
              <td class="px-4 py-3 text-right">
                <button @click="playAudio(audio)" class="p-1 text-gray-400 hover:text-purple-600 mr-2" title="播放">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
                <button @click="deleteAudio(audio.id)" class="p-1 text-gray-400 hover:text-red-600" title="删除">
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="audios.length === 0 && !loading" class="py-16 text-center">
          <svg class="mx-auto w-16 h-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">暂无音频资源</h3>
          <p class="mt-1 text-sm text-gray-500">上传第一个音频文件开始使用</p>
        </div>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-4 py-3 border-t border-gray-200" v-if="audios.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 个音频</span>
        <div class="flex gap-1">
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page = 1">首页</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page--">上一页</button>
          <button v-for="p in visiblePages" :key="p" class="w-8 h-8 text-sm border rounded" :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : 'hover:bg-gray-50'" @click="pagination.page = p; handleSearch()">{{ p }}</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page >= totalPages" @click="pagination.page++">下一页</button>
        </div>
      </div>
    </div>

    <!-- Batch Operations Bar -->
    <div v-if="selectedAudios.length > 0" class="fixed bottom-6 left-1/2 -translate-x-1/2 bg-gray-900 text-white px-6 py-3 rounded-full shadow-lg flex items-center gap-4 z-40">
      <span class="text-sm">{{ selectedAudios.length }} 项已选</span>
      <button @click="batchDelete" class="text-sm hover:text-red-300">删除</button>
      <button @click="selectedAudios = []" class="text-sm text-gray-400 hover:text-white">取消</button>
    </div>

    <!-- Upload Dialog -->
    <Teleport to="body">
      <div v-if="showUploadDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showUploadDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">上传音频</h3>
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
            @dragover.prevent="handleDragOver"
            @dragleave.prevent="handleDragLeave"
            @drop.prevent="handleDrop"
            @click="triggerFileInput"
          >
            <input ref="fileInputRef" type="file" multiple accept="audio/*" class="hidden" @change="handleFileSelect" />
            <svg class="mx-auto w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
            </svg>
            <p class="mt-2 text-sm text-gray-600"><span class="text-indigo-600 font-medium">点击选择</span> 或拖放音频到此处</p>
            <p class="mt-1 text-xs text-gray-500">支持 MP3、WAV、OGG、AAC 等格式</p>
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

    <!-- Audio Player -->
    <div v-if="currentAudio" class="fixed bottom-0 left-0 right-0 bg-white shadow-lg z-50">
      <div class="flex items-center gap-4 p-4">
        <button @click="currentAudio = null" class="text-gray-400 hover:text-gray-600">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <audio :src="`/api${currentAudio.file_path}`" controls class="flex-1" autoplay ref="audioPlayer"></audio>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import {
  getMediaAssetList,
  deleteMediaAsset,
  batchDeleteMediaAssets,
  formatFileSize as apiFormatFileSize,
  formatDuration as apiFormatDuration,
} from '@/api/media'

interface AudioItem {
  id: number
  name: string
  original_name: string
  extension: string
  file_size: number
  file_path: string
  duration: number
  created_at: string
}

const audios = ref<AudioItem[]>([])
const loading = ref(false)
const selectedAudios = ref<number[]>([])
const currentAudio = ref<AudioItem | null>(null)

const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const searchParams = reactive({ keyword: '', audio_type: '' })

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, pagination.page - 2)
  const end = Math.min(totalPages.value, pagination.page + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const isAllSelected = computed(() => audios.value.length > 0 && audios.value.every((a: any) => selectedAudios.value.includes(a.id)))
const formatFileSizeFn = (bytes: number) => apiFormatFileSize(bytes)
const formatDuration = (seconds: number) => apiFormatDuration(seconds)
const formatDate = (dateStr: string) => { try { return new Date(dateStr).toLocaleDateString('zh-CN') } catch { return dateStr } }

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => { if (searchTimer) clearTimeout(searchTimer); searchTimer = setTimeout(() => { pagination.page = 1; loadAudios() }, 300) }
const handleSearch = () => { pagination.page = 1; loadAudios() }

const loadAudios = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.pageSize, file_type: 'audio' }
    if (searchParams.keyword) params.keyword = searchParams.keyword
    if (searchParams.audio_type) params.audio_type = searchParams.audio_type
    const res = await getMediaAssetList(params)
    audios.value = (res.data?.list || []) as any
    pagination.total = res.data?.total || 0
  } catch (error) { message.error('获取音频列表失败') } finally { loading.value = false }
}

const toggleSelect = (id: number) => { const index = selectedAudios.value.indexOf(id); if (index > -1) selectedAudios.value.splice(index, 1); else selectedAudios.value.push(id) }
const toggleSelectAll = () => { if (isAllSelected.value) selectedAudios.value = []; else selectedAudios.value = audios.value.map((a: any) => a.id) }

const playAudio = (audio: AudioItem) => {
  currentAudio.value = audio
}

const deleteAudio = async (id: number) => {
  if (!confirm('确定删除此音频吗？')) return
  try { await deleteMediaAsset(id); message.success('删除成功'); loadAudios() } catch (error: any) { message.error(error.message || '删除失败') }
}

const batchDelete = async () => {
  if (!confirm(`确定删除选中的 ${selectedAudios.value.length} 个音频吗？`)) return
  try { await batchDeleteMediaAssets(selectedAudios.value); message.success('批量删除成功'); selectedAudios.value = []; loadAudios() } catch (error: any) { message.error(error.message || '批量删除失败') }
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
  if (uploadFiles.value.length === 0) { message.warning('请选择要上传的文件'); return }
  uploading.value = true
  let successCount = 0
  for (let i = 0; i < uploadFiles.value.length; i++) {
    const fileItem = uploadFiles.value[i]
    const formData = new FormData()
    formData.append('file', fileItem.file)
    uploadFiles.value[i].status = 'uploading'
    try {
      const token = localStorage.getItem('token')
      const response = await fetch('/api/upload/initiate', {
        method: 'POST',
        headers: { 'Authorization': `Bearer ${token}` },
        body: formData
      })
      if (response.ok) {
        const result = await response.json()
        uploadFiles.value[i].status = 'success'
        uploadFiles.value[i].progress = 100
        successCount++
      }
    } catch (error) {
      console.error('上传失败:', error)
      uploadFiles.value[i].status = 'error'
    }
  }
  uploading.value = false
  message.success(`成功上传 ${successCount} 个文件`)
  showUploadDialog.value = false
  uploadFiles.value = []
  loadAudios()
}

const audioStats = ref<any>({ total_size: 0 })
const stats = ref<any>({ audio_count: 0, mp3_count: 0 })

onMounted(() => { loadAudios() })
</script>