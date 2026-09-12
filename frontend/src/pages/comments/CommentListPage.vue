<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">评论管理</h1>
        <p class="text-sm text-gray-500">管理系统评论内容</p>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总评论数</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-blue-500">
        <p class="text-sm text-gray-500">待审核</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.pending }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4-emerald-500">
        <p class="text-sm text-gray-500">已通过</p>
        <p class="text-2xl font-bold text-emerald-600">{{ stats.approved }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
        <p class="text-sm text-gray-500">已驳回</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.rejected }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-gray-500">
        <p class="text-sm text-gray-500">已删除</p>
        <p class="text-2xl font-bold text-gray-600">{{ stats.deleted }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-purple-500">
        <p class="text-sm text-gray-500">今日新增</p>
        <p class="text-2xl font-bold text-purple-600">{{ stats.by_article + stats.by_video }}</p>
      </div>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索评论..." class="ta-input pl-10 w-full" @input="debounceSearch" />
        </div>
        <select v-model="searchParams.status" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="0">待审核</option>
          <option :value="1">已通过</option>
          <option :value="2">已驳回</option>
          <option :value="3">已删除</option>
        </select>
        <select v-model="searchParams.type" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部类型</option>
          <option value="article">文章评论</option>
          <option value="video">视频评论</option>
        </select>
      </div>
    </TailCard>

    <!-- Tab Headers -->
    <div class="bg-white rounded-t-lg shadow-sm border-b border-gray-200">
      <nav class="flex -mb-px" aria-label="Tabs">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          :class="activeTab === tab.key
            ? 'border-indigo-500 text-indigo-600'
            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'"
          class="py-4 px-6 border-b-2 font-medium text-sm transition-colors flex items-center gap-2"
        >
          {{ tab.label }}
          <span v-if="tab.count !== undefined" :class="{
            'bg-blue-100 text-blue-600': tab.key === 'pending',
            'bg-emerald-100 text-emerald-600': tab.key === 'approved',
            'bg-red-100 text-red-600': tab.key === 'rejected',
          }" class="px-2 py-0.5 rounded-full text-xs font-medium">{{ tab.count }}</span>
        </button>
      </nav>
    </div>

    <!-- Toolbar -->
    <div class="bg-white border-x shadow-sm px-4 py-3 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <label class="flex items-center gap-2 cursor-pointer select-none">
          <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll"
            class="w-4 h-4 text-indigo-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500" />
          <span class="text-sm text-gray-700">全选</span>
        </label>
        <template v-if="selectedIds.length > 0">
          <span class="text-sm text-gray-500">已选 {{ selectedIds.length }} 项</span>
          <TailButton type="success" variant="outline" size="sm" @click="batchApprove">批量通过</TailButton>
          <TailButton type="danger" variant="outline" size="sm" @click="showRejectDialog">批量驳回</TailButton>
          <TailButton type="danger" variant="outline" size="sm" @click="batchDelete">批量删除</TailButton>
        </template>
      </div>
    </div>

    <!-- Comments Table -->
    <TailCard class="rounded-t-none border-x border-b shadow-sm overflow-hidden flex-1 min-h-[400px] flex flex-col">
      <div class="overflow-x-auto flex-1">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
                <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll"
                  class="w-4 h-4 text-indigo-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500" />
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">评论内容</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="comments.length === 0" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">暂无评论数据</td>
            </tr>
            <tr 
              v-for="comment in comments" 
              :key="comment.id" 
              :class="isSelected(comment.id) ? 'bg-indigo-50' : 'hover:bg-gray-50'"
              class="transition-colors"
            >
              <td class="px-4 py-4">
                <input type="checkbox" :value="comment.id" v-model="selectedIds"
                  class="w-4 h-4 text-indigo-600 bg-gray-100 border-gray-300 rounded focus:ring-indigo-500" />
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-xs font-medium">
                    {{ getAvatarLetter(comment.nickname) }}
                  </div>
                  <span class="text-sm font-medium text-gray-900">{{ comment.nickname || '匿名用户' }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <p class="text-sm text-gray-600 max-w-md line-clamp-2">{{ comment.content }}</p>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-blue-100 text-blue-800': comment.article_id > 0,
                  'bg-purple-100 text-purple-800': comment.video_id > 0
                }" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium">
                  {{ comment.article_id > 0 ? '文章' : '视频' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-emerald-100 text-emerald-800': comment.status === 1,
                  'bg-amber-100 text-amber-800': comment.status === 0,
                  'bg-red-100 text-red-800': comment.status === 2,
                  'bg-gray-100 text-gray-800': comment.status === 3
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ ['待审核', '已通过', '已驳回', '已删除'][comment.status] || '未知' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(comment.created_at) }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <TailButton v-if="comment.status === 0" type="success" variant="ghost" size="sm" @click="approveComment(comment)">通过</TailButton>
                  <TailButton v-if="comment.status === 0" type="warning" variant="ghost" size="sm" @click="showSingleReject(comment)">驳回</TailButton>
                  <TailButton type="danger" variant="ghost" size="sm" @click="deleteComment(comment)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200" v-if="comments.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page = 1; handleSearch()">首页</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
          <TailButton 
            v-for="p in visiblePages" 
            :key="p"
            type="default" 
            size="sm" 
            class="w-9 h-9 p-0 flex items-center justify-center"
            :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
            @click="pagination.page = p; handleSearch()">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Reject Dialog -->
    <Teleport to="body">
      <div v-if="showRejectDialogModal" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showRejectDialogModal = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">批量驳回</h3>
          <p class="text-sm text-gray-600 mb-2">
            将对 {{ isSingleReject ? '1' : selectedIds.length }} 条评论执行驳回操作
          </p>
          <textarea v-model="rejectReason" placeholder="请输入驳回原因（必填）"
            class="w-full h-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 resize-none text-sm" />
          <div class="flex justify-end gap-3 mt-4">
            <TailButton variant="outline" size="sm" @click="showRejectDialogModal = false">取消</TailButton>
            <TailButton type="danger" size="sm" @click="batchReject">确认驳回</TailButton>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { 
  getCommentList, 
  approveComment as auditComment, 
  deleteComment as deleteCommentApi,
  batchDeleteComments,
  batchApproveComments,
  batchRejectComments,
  getCommentAuditStats 
} from '@/api/comment'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'

interface CommentItem {
  id: number
  nickname: string
  content: string
  article_id: number
  video_id: number
  status: number
  created_at: string
}

const comments = ref<CommentItem[]>([])
const loading = ref(false)
const selectedIds = ref<number[]>([])
const activeTab = ref<string>('all')
const showRejectDialogModal = ref(false)
const rejectReason = ref('')
const singleRejectComment = ref<CommentItem | null>(null)

const pagination = reactive({
  page: 1,
  pageSize: 15,
  total: 0
})

const searchParams = reactive({
  keyword: '',
  status: undefined as number | undefined,
  type: '' as string
})

const stats = ref({
  total: 0,
  pending: 0,
  approved: 0,
  rejected: 0,
  deleted: 0,
  by_article: 0,
  by_video: 0
})

const tabs = computed(() => [
  { key: 'all', label: '全部评论', count: stats.value.total },
  { key: 'pending', label: '待审核', count: stats.value.pending },
  { key: 'approved', label: '已通过', count: stats.value.approved },
  { key: 'rejected', label: '已驳回', count: stats.value.rejected }
])

const isSingleReject = computed(() => !!singleRejectComment.value)

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))

const isAllSelected = computed(() => {
  if (comments.value.length === 0) return false
  return comments.value.every(c => selectedIds.value.includes(c.id))
})

const visiblePages = computed(() => {
  const total = totalPages.value
  const current = pagination.page
  const pages = []
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const getAvatarLetter = (nickname: string) => {
  if (!nickname || nickname === '匿名用户') return 'U'
  return nickname.charAt(0).toUpperCase()
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN', { 
      year: 'numeric', 
      month: '2-digit', 
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const fetchData = async () => {
  let statusFilter = searchParams.status
  if (statusFilter === undefined && activeTab.value !== 'all') {
    const tabStatusMap: Record<string, number> = {
      'pending': 0,
      'approved': 1,
      'rejected': 2
    }
    statusFilter = tabStatusMap[activeTab.value]
  }

  loading.value = true
  try {
    const res = await getCommentList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...(searchParams.keyword ? { keyword: searchParams.keyword } : {}),
      ...(statusFilter !== undefined ? { status: statusFilter } : {}),
      ...(searchParams.type === 'article' ? { article_id: 0 } : {}),
      ...(searchParams.type === 'video' ? { video_id: 0 } : {})
    })
    const rawList = (res.data?.list || []) as any[]
    comments.value = rawList.map((item: any) => ({
      id: item.id || 0,
      nickname: item.nickname || item.user_name || '',
      content: item.content || '',
      article_id: item.article_id || 0,
      video_id: item.video_id || 0,
      status: item.status ?? 0,
      created_at: item.created_at || new Date().toISOString()
    }))
    pagination.total = res.data?.total || 0
  } catch (error) {
    message.error('获取评论列表失败')
    console.error('获取评论列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getCommentAuditStats()
    stats.value = {
      total: res.data?.total || 0,
      pending: res.data?.pending || 0,
      approved: res.data?.approved || 0,
      rejected: res.data?.rejected || 0,
      deleted: res.data?.deleted || 0,
      by_article: res.data?.by_article || 0,
      by_video: res.data?.by_video || 0
    }
  } catch (error) {
    console.error('获取审核统计失败:', error)
  }
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    fetchData()
  }, 300)
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = comments.value.map(c => c.id)
  }
}

const isSelected = (id: number) => selectedIds.value.includes(id)

const approveComment = async (comment: CommentItem) => {
  try {
    await auditComment(comment.id, true)
    message.success('审核通过')
    await fetchData()
    loadStats()
  } catch (error) {
    message.error('审核失败')
  }
}

const deleteComment = async (comment: CommentItem) => {
  if (!confirm(`确定删除该评论吗？`)) return
  try {
    await deleteCommentApi(comment.id)
    message.success('删除成功')
    selectedIds.value = selectedIds.value.filter(id => id !== comment.id)
    await fetchData()
    loadStats()
  } catch (error) {
    message.error('删除失败')
  }
}

const batchApprove = async () => {
  if (selectedIds.value.length === 0) {
    message.warning('请先选择评论')
    return
  }
  try {
    await batchApproveComments(selectedIds.value)
    message.success(`批量通过 ${selectedIds.value.length} 条评论`)
    selectedIds.value = []
    await fetchData()
    loadStats()
  } catch (error) {
    message.error('批量通过失败')
  }
}

const showSingleReject = (comment: CommentItem) => {
  singleRejectComment.value = comment
  rejectReason.value = ''
  showRejectDialogModal.value = true
}

const showRejectDialog = () => {
  if (selectedIds.value.length === 0) {
    message.warning('请先选择评论')
    return
  }
  singleRejectComment.value = null
  rejectReason.value = ''
  showRejectDialogModal.value = true
}

const batchReject = async () => {
  if (!rejectReason.value.trim()) {
    message.warning('请输入驳回原因')
    return
  }
  const idsToReject = singleRejectComment.value 
    ? [singleRejectComment.value.id] 
    : selectedIds.value
  try {
    await batchRejectComments(idsToReject, rejectReason.value)
    message.success(`批量驳回 ${idsToReject.length} 条评论`)
    selectedIds.value = selectedIds.value.filter(id => !idsToReject.includes(id))
    showRejectDialogModal.value = false
    await fetchData()
    loadStats()
  } catch (error) {
    message.error('批量驳回失败')
  }
}

const batchDelete = async () => {
  if (selectedIds.value.length === 0) {
    message.warning('请先选择评论')
    return
  }
  if (!confirm(`确定删除选中的 ${selectedIds.value.length} 条评论吗？此操作不可恢复`)) return
  try {
    await batchDeleteComments(selectedIds.value)
    message.success(`批量删除 ${selectedIds.value.length} 条评论`)
    selectedIds.value = []
    await fetchData()
    loadStats()
  } catch (error) {
    message.error('批量删除失败')
  }
}

onMounted(() => {
  fetchData()
  loadStats()
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