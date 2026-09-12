<template>
  <div class="page-container max-w-4xl">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center gap-2 mb-2">
        <button class="text-gray-500 hover:text-gray-700" @click="$router.back()">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900 mb-1">{{ isEdit ? '编辑视频' : '上传视频' }}</h1>
          <p class="text-sm text-gray-500">{{ isEdit ? '修改视频信息' : '发布新的视频内容' }}</p>
        </div>
      </div>
    </div>

    <!-- Form Card -->
    <TailCard class="p-6">
      <form @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Cover Upload -->
        <div>
          <label class="ta-label">视频封面</label>
          <div class="flex items-center gap-4">
            <div v-if="formData.coverUrl" class="w-32 h-20 rounded-lg overflow-hidden bg-gray-100 border border-gray-200">
              <img :src="formData.coverUrl" class="w-full h-full object-cover" alt="封面图" />
            </div>
            <div class="flex-1">
              <div
                ref="coverDropZoneRef"
                class="border-2 border-dashed rounded-lg p-4 text-center cursor-pointer transition-colors"
                :class="isCoverDragOver ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400 hover:bg-gray-50'"
                @dragover.prevent.stop="handleCoverDragOver"
                @dragleave.prevent.stop="handleCoverDragLeave"
                @drop.prevent.stop="handleCoverDrop"
                @click.stop="triggerCoverFileInput"
              >
                <input ref="coverFileInputRef" type="file" accept="image/jpeg,image/png,image/gif,image/webp" class="hidden" @change="handleCoverFileSelect" />
                <p class="text-sm text-gray-600"><span class="text-indigo-600 font-medium">点击选择</span> 或拖放封面图片</p>
                <p class="mt-1 text-xs text-gray-500">支持 JPG、PNG、GIF、WebP，建议尺寸 1280x720</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Title -->
        <div>
          <label class="ta-label">视频标题 <span class="text-red-500">*</span></label>
          <input v-model="formData.title" type="text" class="ta-input w-full" placeholder="请输入视频标题" required />
        </div>

        <!-- Description -->
        <div>
          <label class="ta-label">视频描述</label>
          <textarea v-model="formData.description" class="ta-input w-full" rows="4" placeholder="请输入视频描述"></textarea>
        </div>

        <!-- Video File Upload -->
        <div>
          <label class="ta-label">上传视频文件</label>
          <div
            ref="videoDropZoneRef"
            class="border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors"
            :class="isVideoDragOver ? 'border-indigo-500 bg-indigo-50' : 'border-gray-300 hover:border-indigo-400 hover:bg-gray-50'"
            @dragover.prevent.stop="handleVideoDragOver"
            @dragleave.prevent.stop="handleVideoDragLeave"
            @drop.prevent.stop="handleVideoDrop"
            @click.stop="triggerVideoFileInput"
          >
            <input ref="videoFileInputRef" type="file" accept="video/mp4,video/webm,video/avi,video/mov,video/mkv" class="hidden" @change="handleVideoFileSelect" />
            <svg class="mx-auto w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <p class="mt-2 text-sm text-gray-600"><span class="text-indigo-600 font-medium">点击选择</span> 或拖放视频文件到此处</p>
            <p class="mt-1 text-xs text-gray-500">支持 MP4、WebM、AVI 等格式，最大 500MB</p>
            <div v-if="selectedVideoFile" class="mt-3">
              <div class="flex items-center justify-center gap-2 text-sm">
                <svg class="w-4 h-4 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="text-gray-700">{{ selectedVideoFile.name }}</span>
                <span class="text-gray-400">({{ formatFileSize(selectedVideoFile.size) }})</span>
              </div>
              <div v-if="uploadingVideo" class="mt-2">
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div class="bg-indigo-600 h-2 rounded-full animate-pulse" style="width: 60%"></div>
                </div>
                <p class="text-xs text-indigo-600 mt-1">上传中...</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Category and Video URL -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="ta-label">分类 <span class="text-red-500">*</span></label>
            <select v-model="formData.categoryId" class="ta-input w-full" required>
              <option :value="0">请选择分类</option>
              <option v-for="c in categoryTree" :key="c.id" :value="c.id">{{ getIndent(c.depth || 0) }}{{ c.name }}</option>
            </select>
          </div>
          <div>
            <label class="ta-label">视频URL（上传后自动填充）</label>
            <input v-model="formData.videoUrl" type="url" class="ta-input w-full" placeholder="https://example.com/video.mp4 或上传本地视频" />
          </div>
        </div>

        <!-- Tags -->
        <div>
          <label class="ta-label">标签</label>
          <input v-model="tagInput" type="text" class="ta-input w-full" placeholder="输入标签后按回车添加" @keydown.enter.prevent="addTag" />
          <div class="flex flex-wrap gap-2 mt-2">
            <span v-for="(tag, i) in formData.tags" :key="i" class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
              {{ tag }}
              <button type="button" @click="formData.tags.splice(i, 1)" class="ml-2 text-red-500 hover:text-red-700 font-bold">&times;</button>
            </span>
          </div>
        </div>

        <!-- Status -->
        <div class="flex items-center gap-6">
          <label class="flex items-center gap-2 cursor-pointer">
            <input v-model="formData.status" type="checkbox" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
            <span class="text-sm text-gray-700">立即发布</span>
          </label>
        </div>

        <!-- Submit Buttons -->
        <div class="flex gap-3 pt-4 border-t border-gray-200">
          <button type="submit" class="ta-btn ta-btn-primary">
            {{ isEdit ? '保存修改' : '发布视频' }}
          </button>
          <button type="button" class="ta-btn ta-btn-outline" @click="$router.back()">取消</button>
        </div>
      </form>
    </TailCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getVideoCategoryTree, getVideo, createVideo, updateVideo } from '@/api/video'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'

interface CategoryItem {
  id: number
  name: string
  depth?: number
  children?: CategoryItem[]
}

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)

const formData = reactive({
  title: '',
  description: '',
  coverUrl: '',
  videoUrl: '',
  categoryId: 0,
  tags: [] as string[],
  status: 0 as 0 | 1 | 2
})

const tagInput = ref('')
const categoryTree = ref<CategoryItem[]>([])

// ========== 视频上传相关 ==========
const videoDropZoneRef = ref<HTMLDivElement | null>(null)
const videoFileInputRef = ref<HTMLInputElement | null>(null)
const isVideoDragOver = ref(false)
const selectedVideoFile = ref<File | null>(null)
const uploadingVideo = ref(false)

// ========== 封面上传相关 ==========
const coverDropZoneRef = ref<HTMLDivElement | null>(null)
const coverFileInputRef = ref<HTMLInputElement | null>(null)
const isCoverDragOver = ref(false)
const selectedCoverFile = ref<File | null>(null)
const uploadingCover = ref(false)

const getIndent = (depth: number) => '  '.repeat(depth || 0)

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  const rank = Math.floor(Math.log(bytes) / Math.log(1024))
  return parseFloat((bytes / Math.pow(1024, rank)).toFixed(1)) + ' ' + units[rank]
}

const loadCategories = async () => {
  try {
    const res = await getVideoCategoryTree()
    categoryTree.value = (res.data as CategoryItem[]) || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const addTag = () => {
  if (tagInput.value && !formData.tags.includes(tagInput.value)) {
    formData.tags.push(tagInput.value)
  }
  tagInput.value = ''
}

// ========== 视频上传方法 ==========
const handleVideoDragOver = (e: DragEvent) => { isVideoDragOver.value = true; e.stopPropagation() }
const handleVideoDragLeave = (e: DragEvent) => { isVideoDragOver.value = false; e.stopPropagation() }
const handleVideoDrop = (e: DragEvent) => {
  isVideoDragOver.value = false
  e.stopPropagation()
  if (e.dataTransfer?.files?.length) {
    const file = e.dataTransfer.files[0]
    if (file.type.startsWith('video/') || /\.(mp4|webm|avi|mov|mkv)$/i.test(file.name)) {
      selectedVideoFile.value = file
      uploadVideoFile(file)
    } else {
      message.warning('请选择视频文件（MP4、WebM、AVI 等格式）')
    }
  }
}
const triggerVideoFileInput = () => { videoFileInputRef.value?.click() }
const handleVideoFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files?.length) {
    const file = target.files[0]
    selectedVideoFile.value = file
    uploadVideoFile(file)
    target.value = ''
  }
}

const uploadVideoFile = async (file: File) => {
  uploadingVideo.value = true
  try {
    const token = localStorage.getItem('token')
    
    // 步骤1：初始化上传
    const initiateFormData = new FormData()
    initiateFormData.append('filename', file.name)
    initiateFormData.append('size', String(file.size))
    initiateFormData.append('chunk_num', '0')
    initiateFormData.append('file_type', 'video')
    
    const initiateRes = await fetch('/api/upload/initiate', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` },
      body: initiateFormData
    })
    
    if (!initiateRes.ok) throw new Error('初始化上传失败')
    const initiateResult = await initiateRes.json()
    const uploadID = initiateResult.data.upload_id
    
    // 步骤2：上传分片（单片模式）
    const chunkFormData = new FormData()
    chunkFormData.append('chunk', file)
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
        total_size: file.size,
        filename: file.name,
        file_type: 'video'
      })
    })
    
    if (!completeRes.ok) throw new Error('完成上传失败')
    const completeResult = await completeRes.json()
    
    // 更新表单中的视频URL
    formData.videoUrl = completeResult.data.url
    message.success('视频上传成功')
    selectedVideoFile.value = null
  } catch (error: any) {
    console.error('视频上传失败:', error)
    message.error(error.message || '视频上传失败')
  } finally {
    uploadingVideo.value = false
  }
}

// ========== 封面上传方法 ==========
const handleCoverDragOver = (e: DragEvent) => { isCoverDragOver.value = true; e.stopPropagation() }
const handleCoverDragLeave = (e: DragEvent) => { isCoverDragOver.value = false; e.stopPropagation() }
const handleCoverDrop = (e: DragEvent) => {
  isCoverDragOver.value = false
  e.stopPropagation()
  if (e.dataTransfer?.files?.length) {
    const file = e.dataTransfer.files[0]
    if (file.type.startsWith('image/') || /\.(jpg|jpeg|png|gif|webp)$/i.test(file.name)) {
      selectedCoverFile.value = file
      uploadCoverFile(file)
    } else {
      message.warning('请选择图片文件（JPG、PNG、GIF、WebP 格式）')
    }
  }
}
const triggerCoverFileInput = () => { coverFileInputRef.value?.click() }
const handleCoverFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement
  if (target.files?.length) {
    const file = target.files[0]
    selectedCoverFile.value = file
    uploadCoverFile(file)
    target.value = ''
  }
}

const uploadCoverFile = async (file: File) => {
  uploadingCover.value = true
  try {
    const token = localStorage.getItem('token')
    
    // 步骤1：初始化上传
    const initiateFormData = new FormData()
    initiateFormData.append('filename', file.name)
    initiateFormData.append('size', String(file.size))
    initiateFormData.append('chunk_num', '0')
    initiateFormData.append('file_type', 'image')
    
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
    chunkFormData.append('chunk', file)
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
        total_size: file.size,
        filename: file.name,
        file_type: 'image'
      })
    })
    
    if (!completeRes.ok) throw new Error('完成上传失败')
    const completeResult = await completeRes.json()
    
    // 更新表单中的封面URL
    formData.coverUrl = completeResult.data.url
    message.success('封面上传成功')
    selectedCoverFile.value = null
  } catch (error: any) {
    console.error('封面上传失败:', error)
    message.error(error.message || '封面上传失败')
  } finally {
    uploadingCover.value = false
  }
}

const handleSubmit = async () => {
  if (!formData.title) {
    message.error('请输入视频标题')
    return
  }
  if (!formData.categoryId || formData.categoryId === 0) {
    message.error('请选择分类')
    return
  }
  
  try {
    if (isEdit.value && route.params.id) {
      await updateVideo(Number(route.params.id), {
        title: formData.title,
        description: formData.description,
        cover_image: formData.coverUrl,
        video_url: formData.videoUrl,
        category_id: formData.categoryId,
        status: formData.status
      })
      message.success('更新成功')
    } else {
      await createVideo({
        title: formData.title,
        description: formData.description,
        cover_image: formData.coverUrl,
        video_url: formData.videoUrl,
        category_id: formData.categoryId,
        status: formData.status,
        tag_ids: []
      })
      message.success('创建成功')
    }
    router.push('/videos')
  } catch (error) {
    message.error('操作失败')
  }
}

onMounted(() => {
  loadCategories()
  // 如果是编辑模式，加载视频数据
  if (isEdit.value && route.params.id) {
    loadVideoData(Number(route.params.id))
  }
})

// 监听路由参数变化
watch(() => route.params.id, (newId) => {
  if (newId && isEdit.value) {
    loadVideoData(Number(newId))
  }
})

const loadVideoData = async (id: number) => {
  try {
    const res = await getVideo(id)
    const video = res.data
    if (video) {
      formData.title = video.title || ''
      formData.description = video.description || ''
      formData.coverUrl = video.cover_image || ''
      formData.videoUrl = video.video_url || ''
      formData.categoryId = video.category_id || 0
      formData.status = video.status ?? 0
      // 加载标签数据
      if (video.tags && Array.isArray(video.tags)) {
        formData.tags = video.tags.map((tag: any) => tag.name || tag)
      }
    }
  } catch (error) {
    message.error('加载视频数据失败')
    console.error('加载视频数据失败:', error)
  }
}
</script>