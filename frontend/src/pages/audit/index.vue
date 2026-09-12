<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">内容审核</h1>
        <p class="text-sm text-gray-500">管理内容审核流程与状态</p>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-gray-500">
        <p class="text-sm text-gray-500">总审核</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-yellow-500">
        <p class="text-sm text-gray-500">待审核</p>
        <p class="text-2xl font-bold text-yellow-600">{{ stats.pending || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-green-500">
        <p class="text-sm text-gray-500">已通过</p>
        <p class="text-2xl font-bold text-green-600">{{ stats.passed || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
        <p class="text-sm text-gray-500">已驳回</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.rejected || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-orange-500">
        <p class="text-sm text-gray-500">需修改</p>
        <p class="text-2xl font-bold text-orange-600">{{ stats.needs_modify || 0 }}</p>
      </div>
    </div>

    <!-- Filter Bar -->
    <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索审核编号/标题/提交者..." class="ta-input pl-10 w-full" @input="debounceSearch" />
        </div>
        <select v-model="searchParams.contentType" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部类型</option>
          <option value="article">文章</option>
          <option value="video">视频</option>
          <option value="image">图片</option>
          <option value="file">文件</option>
        </select>
        <select v-model="searchParams.status" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="0">待审</option>
          <option :value="1">通过</option>
          <option :value="2">驳回</option>
          <option :value="3">需修改</option>
        </select>
        <div class="ml-auto flex gap-2">
          <button @click="loadData" class="px-4 py-2 text-sm border rounded-lg hover:bg-gray-50 flex items-center gap-2">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
            刷新
          </button>
        </div>
      </div>
    </div>

    <!-- Audit Table -->
    <div class="bg-white rounded-lg shadow-sm overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left w-10">
                <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-4 h-4 rounded" />
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">审核编号</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">内容信息</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">提交者</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">审核者</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">审核时间</th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="8" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="audits.length === 0" class="animate-pulse">
              <td colspan="8" class="px-6 py-12 text-center text-gray-400">暂无审核记录</td>
            </tr>
            <tr v-for="audit in audits" :key="audit.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <input type="checkbox" :checked="selectedIds.includes(audit.id)" @change="toggleSelect(audit.id)" class="w-4 h-4 rounded" />
              </td>
              <td class="px-6 py-4">
                <span class="text-sm font-mono text-gray-900">{{ audit.audit_no }}</span>
              </td>
              <td class="px-6 py-4">
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ audit.content_title || '-' }}</p>
                  <p class="text-xs text-gray-500">
                    <span :class="getContentTypeClass(audit.content_type)" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium mr-1">
                      {{ getContentTypeLabel(audit.content_type) }}
                    </span>
                    ID: {{ audit.content_id }}
                  </p>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ audit.submitter_name || '-' }}</td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ audit.auditor_name || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="getAuditStatusClass(audit.status)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ getAuditStatusLabel(audit.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ audit.audited_at ? formatDate(audit.audited_at) : '-' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-right">
                <button @click="viewDetail(audit)" class="text-indigo-600 hover:text-indigo-900 text-sm mr-3">详情</button>
                <template v-if="audit.status === 0">
                  <button @click="approveItem(audit)" class="text-green-600 hover:text-green-900 text-sm mr-3">通过</button>
                  <button @click="rejectItem(audit)" class="text-red-600 hover:text-red-900 text-sm mr-3">驳回</button>
                  <button @click="needsModifyItem(audit)" class="text-orange-600 hover:text-orange-900 text-sm">需修改</button>
                </template>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200" v-if="audits.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page = 1; handleSearch()">首页</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</button>
          <button v-for="p in visiblePages" :key="p" class="w-8 h-8 text-sm border rounded flex items-center justify-center" :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : 'hover:bg-gray-50'" @click="pagination.page = p; handleSearch()">{{ p }}</button>
          <button class="px-3 py-1 text-sm border rounded hover:bg-gray-50" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</button>
        </div>
      </div>
    </div>

    <!-- Batch Operations -->
    <div v-if="selectedIds.length > 0" class="fixed bottom-6 left-1/2 -translate-x-1/2 bg-gray-900 text-white px-6 py-3 rounded-full shadow-lg flex items-center gap-4 z-40">
      <span class="text-sm">{{ selectedIds.length }} 项已选</span>
      <button @click="batchApprove" class="text-sm hover:text-green-300 flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" /></svg>
        批量通过
      </button>
      <button @click="batchReject" class="text-sm hover:text-red-300 flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
        批量驳回
      </button>
      <button @click="selectedIds = []" class="text-sm text-gray-400 hover:text-white">取消</button>
    </div>

    <!-- Detail Dialog -->
    <Teleport to="body">
      <div v-if="showDetailDialog && currentAudit" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showDetailDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">审核详情</h3>
            <button @click="showDetailDialog = false" class="text-gray-400 hover:text-gray-600">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
          </div>
          <div class="space-y-4">
            <div><span class="text-sm text-gray-500">审核编号:</span> <span class="text-sm font-mono">{{ currentAudit.audit_no }}</span></div>
            <div><span class="text-sm text-gray-500">内容类型:</span> <span :class="getContentTypeClass(currentAudit.content_type)" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium">{{ getContentTypeLabel(currentAudit.content_type) }}</span></div>
            <div><span class="text-sm text-gray-500">内容标题:</span> <span class="text-sm">{{ currentAudit.content_title || '-' }}</span></div>
            <div><span class="text-sm text-gray-500">提交者:</span> <span class="text-sm">{{ currentAudit.submitter_name || '-' }}</span></div>
            <div><span class="text-sm text-gray-500">审核者:</span> <span class="text-sm">{{ currentAudit.auditor_name || '-' }}</span></div>
            <div><span class="text-sm text-gray-500">状态:</span> <span :class="getAuditStatusClass(currentAudit.status)" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ml-1">{{ getAuditStatusLabel(currentAudit.status) }}</span></div>
            <div v-if="currentAudit.audit_comment"><span class="text-sm text-gray-500">审核意见:</span> <p class="text-sm mt-1">{{ currentAudit.audit_comment }}</p></div>
            <div v-if="currentAudit.required_modification"><span class="text-sm text-gray-500">要求修改:</span> <p class="text-sm mt-1 text-orange-600">{{ currentAudit.required_modification }}</p></div>
            <div><span class="text-sm text-gray-500">提交时间:</span> <span class="text-sm">{{ formatDate(currentAudit.created_at) }}</span></div>
            <div v-if="currentAudit.audited_at"><span class="text-sm text-gray-500">审核时间:</span> <span class="text-sm">{{ formatDate(currentAudit.audited_at) }}</span></div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Reject Dialog -->
    <Teleport to="body">
      <div v-if="showRejectDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showRejectDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6 m-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">驳回审核</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">审核意见</label>
              <textarea v-model="rejectForm.audit_comment" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-red-500 focus:border-red-500" placeholder="请输入驳回原因"></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">要求修改内容</label>
              <textarea v-model="rejectForm.required_modification" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-red-500 focus:border-red-500" placeholder="请说明需要修改的内容（可选）"></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <button @click="showRejectDialog = false" class="px-4 py-2 text-sm border rounded-lg hover:bg-gray-50">取消</button>
            <button @click="confirmReject" class="px-4 py-2 text-sm bg-red-600 text-white rounded-lg hover:bg-red-700">确认驳回</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Approve Dialog -->
    <Teleport to="body">
      <div v-if="showApproveDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showApproveDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6 m-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">通过审核</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">审核意见（可选）</label>
              <textarea v-model="approveForm.audit_comment" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-green-500 focus:border-green-500" placeholder="请输入审核意见"></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <button @click="showApproveDialog = false" class="px-4 py-2 text-sm border rounded-lg hover:bg-gray-50">取消</button>
            <button @click="confirmApprove" class="px-4 py-2 text-sm bg-green-600 text-white rounded-lg hover:bg-green-700">确认通过</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import {
  getAuditList,
  getAuditStats,
  approveAudit as apiApprove,
  rejectAudit as apiReject,
  needsModificationAudit as apiNeedsModify,
  batchAudit,
  getAuditStatusLabel,
  getAuditStatusClass,
  getContentTypeLabel,
  formatDate
} from '@/api/audit'

const audits = ref<any[]>([])
const loading = ref(false)
const stats = ref<any>({ total: 0, pending: 0, passed: 0, rejected: 0, needs_modify: 0 })

const pagination = reactive({ page: 1, pageSize: 15, total: 0 })
const searchParams = reactive({ keyword: '', contentType: '', status: '' as string | number })
const selectedIds = ref<number[]>([])

const tabs = computed(() => [
  { key: '', label: '全部', count: stats.value.total || 0 },
  { key: 0, label: '待审核', count: stats.value.pending || 0 },
  { key: 1, label: '已通过', count: stats.value.passed || 0 },
  { key: 2, label: '已驳回', count: stats.value.rejected || 0 },
  { key: 3, label: '需修改', count: stats.value.needs_modify || 0 }
])

const currentTab = ref('')
const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const visiblePages = computed(() => {
  const pages: number[] = []
  const start = Math.max(1, pagination.page - 2)
  const end = Math.min(totalPages.value, pagination.page + 2)
  for (let i = start; i <= end; i++) pages.push(i)
  return pages
})

const isAllSelected = computed(() => audits.value.length > 0 && audits.value.every((a: any) => selectedIds.value.includes(a.id)))

const showDetailDialog = ref(false)
const showRejectDialog = ref(false)
const showApproveDialog = ref(false)
const currentAudit = ref<any>(null)
const rejectForm = reactive({ audit_comment: '', required_modification: '' })
const approveForm = reactive({ audit_comment: '' })

const getContentTypeClass = (type: string) => {
  const classes: Record<string, string> = { article: 'bg-blue-100 text-blue-800', video: 'bg-red-100 text-red-800', image: 'bg-purple-100 text-purple-800', file: 'bg-gray-100 text-gray-800' }
  return classes[type] || ''
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => { if (searchTimer) clearTimeout(searchTimer); searchTimer = setTimeout(() => { pagination.page = 1; loadAudits() }, 300) }
const handleSearch = () => { pagination.page = 1; loadAudits() }

const loadData = async () => { await Promise.all([loadAudits(), loadStats()]) }
const loadStats = async () => { try { const res = await getAuditStats(); stats.value = res.data || {} } catch { } }

const loadAudits = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.pageSize }
    if (searchParams.keyword) params.keyword = searchParams.keyword
    if (searchParams.contentType) params.content_type = searchParams.contentType
    if (searchParams.status !== '') params.status = Number(searchParams.status)
    if (currentTab.value !== '') params.status = currentTab.value

    const res = await getAuditList(params)
    
    // 兼容不同的数据结构：支持 list 和 rows 两种字段名
    if (res.data?.list) {
      audits.value = Array.isArray(res.data.list) ? res.data.list : []
      pagination.total = Number(res.data.total) || 0
    } else if (res.data?.rows) {
      // 某些接口使用 rows 作为列表字段
      audits.value = Array.isArray(res.data.rows) ? res.data.rows : []
      pagination.total = Number(res.data.total) || 0
    } else {
      audits.value = []
      pagination.total = 0
    }
  } catch (err) { 
    message.error('获取审核列表失败') 
  } finally { loading.value = false }
}

const toggleSelect = (id: number) => { const i = selectedIds.value.indexOf(id); if (i > -1) selectedIds.value.splice(i, 1); else selectedIds.value.push(id) }
const toggleSelectAll = () => { if (isAllSelected.value) selectedIds.value = []; else selectedIds.value = audits.value.map((a: any) => a.id) }

const viewDetail = (audit: any) => { currentAudit.value = audit; showDetailDialog.value = true }

const approveItem = (audit: any) => { currentAudit.value = audit; showApproveDialog.value = true }
const confirmApprove = async () => {
  try { await apiApprove(currentAudit.value.id, approveForm); message.success('审核通过'); showApproveDialog.value = false; loadData() } catch (e: any) { message.error(e.message || '操作失败') }
}

const rejectItem = (audit: any) => { currentAudit.value = audit; showRejectDialog.value = true }
const confirmReject = async () => {
  try { await apiReject(currentAudit.value.id, rejectForm); message.success('审核驳回'); showRejectDialog.value = false; loadData() } catch (e: any) { message.error(e.message || '操作失败') }
}

const needsModifyItem = async (audit: any) => {
  if (!confirm(`标记"${audit.content_title}"为需修改？`)) return
  try {
    await apiNeedsModify(audit.id, { required_modification: '请根据审核意见修改后重新提交' })
    message.success('已标记需修改')
    loadData()
  } catch (e: any) { message.error(e.message || '操作失败') }
}

const batchApprove = async () => {
  if (!confirm(`通过选中的 ${selectedIds.value.length} 条记录？`)) return
  try { await batchAudit({ ids: selectedIds.value, action: 'approve' }); message.success('批量通过成功'); loadData() } catch (e: any) { message.error(e.message || '操作失败') }
}

const batchReject = async () => {
  if (!confirm(`驳回选中的 ${selectedIds.value.length} 条记录？`)) return
  showRejectDialog.value = true
}

onMounted(() => { loadData() })
</script>