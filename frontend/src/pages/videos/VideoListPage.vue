<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">视频管理</h1>
        <p class="text-sm text-gray-500">管理系统视频内容</p>
      </div>
      <TailButton type="primary" @click="$router.push('/videos/create')">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新建视频
      </TailButton>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索视频..." class="ta-input pl-10 w-full" @input="handleSearch" />
        </div>
        <select v-model="searchParams.categoryId" class="ta-input w-auto min-w-[140px]" @change="handleSearch">
          <option value="">全部分类</option>
          <option v-for="cat in categoryTree" :key="cat.id" :value="cat.id">{{ getIndent(cat.depth || 0) }}{{ cat.name }}</option>
        </select>
        <select v-model="searchParams.status" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="0">草稿</option>
          <option :value="1">已发布</option>
          <option :value="2">已下架</option>
        </select>
      </div>
    </TailCard>

    <!-- Video Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      <div v-for="video in videoList" :key="video.id" class="bg-white rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-200">
        <div class="relative h-48 bg-gradient-to-br from-indigo-50 to-purple-50 flex items-center justify-center">
          <svg class="w-16 h-16 text-indigo-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <span :class="{
            'bg-emerald-100 text-emerald-800': video.status === 1,
            'bg-amber-100 text-amber-800': video.status === 0,
            'bg-red-100 text-red-800': video.status === 2
          }" class="absolute top-2 right-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium">
            {{ video.status === 1 ? '已发布' : video.status === 0 ? '草稿' : '已下架' }}
          </span>
        </div>
        <div class="p-4">
          <h3 class="text-sm font-medium text-gray-900 mb-2 line-clamp-2">{{ video.title }}</h3>
          <p class="text-xs text-gray-500 mb-1">分类: {{ video.category_name || '未分类' }}</p>
          <p class="text-xs text-gray-500 mb-3">创建: {{ formatDate(video.created_at) }}</p>
          <div class="flex gap-2 pt-3 border-t border-gray-100">
            <TailButton type="default" variant="ghost" size="sm" class="flex-1" @click="$router.push(`/videos/${video.id}/edit`)">编辑</TailButton>
            <TailButton type="danger" variant="ghost" size="sm" @click="deleteVideo(video)">删除</TailButton>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="videoList.length === 0 && !loading" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
      </svg>
      <p class="text-gray-500">暂无视频数据</p>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between px-6 py-4 mt-6 bg-white rounded-lg shadow-sm" v-if="videoList.length > 0">
      <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
      <div class="flex gap-1">
        <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
        <TailButton 
          v-for="p in totalPages" 
          :key="p"
          type="default" 
          size="sm" 
          class="w-9 h-9 p-0 flex items-center justify-center"
          :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
          @click="goToPage(p)">{{ p }}</TailButton>
        <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</TailButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getVideoList, deleteVideo as deleteVideoApi, getVideoCategoryTree } from '@/api/video'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'

interface VideoItem {
  id: number
  title: string
  description?: string
  cover_image: string
  video_url: string
  category_id: number
  category_name?: string
  status: 0 | 1 | 2
  view_count: number
  duration: number
  created_at: string
  updated_at: string
}

interface CategoryItem {
  id: number
  name: string
  depth?: number
  children?: CategoryItem[]
}

const videoList = ref<VideoItem[]>([])
const categoryTree = ref<CategoryItem[]>([])
const loading = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 12,
  total: 0
})

const searchParams = reactive({
  keyword: '',
  categoryId: undefined as number | undefined,
  status: undefined as number | undefined
})

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))

const getIndent = (depth: number) => '  '.repeat(depth || 0)

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: '2-digit', 
      day: '2-digit' 
    })
  } catch {
    return dateStr
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getVideoList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...(searchParams.keyword ? { keyword: searchParams.keyword } : {}),
      ...(searchParams.categoryId ? { category_id: searchParams.categoryId } : {}),
      ...(searchParams.status !== undefined && String(searchParams.status) !== '' ? { status: Number(searchParams.status) } : {})
    })
    videoList.value = (res.data?.list || []) as any
    pagination.total = res.data?.total || 0
  } catch (error) {
    message.error('获取视频列表失败')
    console.error('获取视频列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await getVideoCategoryTree()
    categoryTree.value = (res.data as CategoryItem[]) || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const goToPage = (p: number) => {
  pagination.page = p
  fetchData()
}

const deleteVideo = async (video: VideoItem) => {
  if (!confirm(`确定删除视频 "${video.title}" 吗？`)) return
  try {
    await deleteVideoApi(video.id)
    message.success('删除成功')
    await fetchData()
  } catch (error) {
    message.error('删除失败')
  }
}

onMounted(() => {
  fetchData()
  loadCategories()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>