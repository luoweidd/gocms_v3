<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">操作日志</h1>
        <p class="text-sm text-gray-500">查看系统用户操作记录与审计信息</p>
      </div>
      <TailButton type="primary" @click="handleExport">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        导出日志
      </TailButton>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索日志..." class="ta-input pl-10 w-full" @input="debounceSearch" />
        </div>
        <select v-model="searchParams.action_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部操作</option>
          <option value="create">创建</option>
          <option value="update">更新</option>
          <option value="delete">删除</option>
          <option value="query">查询</option>
          <option value="login">登录</option>
        </select>
        <select v-model="searchParams.resource_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部资源</option>
          <option value="article">文章</option>
          <option value="video">视频</option>
          <option value="user">用户</option>
          <option value="menu">菜单</option>
          <option value="role">角色</option>
          <option value="system">系统</option>
        </select>
        <input v-model="searchParams.date_range" type="date" class="ta-input w-auto" @change="handleSearch" />
        <TailButton variant="outline" size="sm" @click="resetSearch">重置</TailButton>
      </div>
    </TailCard>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总日志数</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-blue-500">
        <p class="text-sm text-gray-500">今日操作</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.today_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-emerald-500">
        <p class="text-sm text-gray-500">成功操作</p>
        <p class="text-2xl font-bold text-emerald-600">{{ stats.success_count || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
        <p class="text-sm text-gray-500">失败操作</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.fail_count || 0 }}</p>
      </div>
    </div>

    <!-- Module & Action Distribution -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
      <TailCard class="mb-0" v-if="stats.action_stats && stats.action_stats.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">操作类型分布</h3>
        </div>
        <div class="grid grid-cols-3 md:grid-cols-5 gap-3">
          <div v-for="item in stats.action_stats" :key="item.action" 
               class="bg-gray-50 rounded-lg p-3 hover:bg-gray-100 transition-colors">
            <div class="flex items-center justify-between mb-1">
              <span class="text-xs text-gray-500">{{ getActionLabel(item.action) }}</span>
            </div>
            <p class="text-lg font-bold text-gray-900">{{ item.count }}</p>
          </div>
        </div>
      </TailCard>
      <TailCard class="mb-0" v-if="stats.module_stats && stats.module_stats.length > 0">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">模块分布</h3>
        </div>
        <div class="grid grid-cols-3 md:grid-cols-5 gap-3">
          <div v-for="item in stats.module_stats" :key="item.module" 
               class="bg-gray-50 rounded-lg p-3 hover:bg-gray-100 transition-colors">
            <div class="flex items-center justify-between mb-1">
              <span class="text-xs text-gray-500">{{ resourceTypeLabel[item.module] || item.module || '未知' }}</span>
            </div>
            <p class="text-lg font-bold text-gray-900">{{ item.count }}</p>
          </div>
        </div>
      </TailCard>
    </div>

    <!-- Log Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户信息</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作描述</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP地址</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">请求方法</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">请求路径</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">耗时(ms)</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="12" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="logs.length === 0" class="animate-pulse">
              <td colspan="12" class="px-6 py-12 text-center text-gray-400">暂无日志数据</td>
            </tr>
            <tr v-for="log in logs" :key="log.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 text-sm text-gray-500">#{{ log.id }}</td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div class="w-8 h-8 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center text-xs font-bold">
                    {{ log.username?.charAt(0) || 'U' }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ log.username || '未知用户' }}</p>
                    <p class="text-xs text-gray-500">{{ log.role_name || '无角色' }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-blue-100 text-blue-800': log.action === 'create',
                  'bg-green-100 text-green-800': log.action === 'update',
                  'bg-red-100 text-red-800': log.action === 'delete',
                  'bg-gray-100 text-gray-800': log.action === 'list' || log.action === 'query' || log.action === 'stats' || log.action === 'tree',
                  'bg-purple-100 text-purple-800': log.action === 'login'
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ getActionLabel(log.action) }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
              {{ resourceTypeLabel[log.module] || log.module || '未知' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600 max-w-xs truncate">
                {{ log.description || log.details || '无描述' }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ log.ip_address || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-blue-100 text-blue-800': log.request_method === 'GET',
                  'bg-green-100 text-green-800': log.request_method === 'POST',
                  'bg-yellow-100 text-yellow-800': log.request_method === 'PUT',
                  'bg-red-100 text-red-800': log.request_method === 'DELETE'
                }" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-bold">
                  {{ log.request_method || '-' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600 font-mono max-w-[200px] truncate" :title="log.request_path">{{ log.request_path || '-' }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ log.duration_ms || 0 }}</td>
              <td class="px-6 py-4">
                <span :class="log.status === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium">
                  {{ log.status === 1 ? '成功' : '失败' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500 whitespace-nowrap">
                {{ formatDateTime(log.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <TailButton variant="ghost" size="sm" @click="viewDetail(log)">查看详情</TailButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200" v-if="logs.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page = 1; handleSearch()">首页</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
          <TailButton 
            v-for="p in visiblePages" 
            :key="p"
            class="w-9 h-9 p-0 flex items-center justify-center"
            :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
            @click="pagination.page = p; handleSearch()">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Detail Dialog -->
    <Teleport to="body">
      <div v-if="showDetailDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showDetailDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">日志详情</h3>
            <TailButton variant="ghost" size="sm" @click="showDetailDialog = false">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </TailButton>
          </div>
          <div class="space-y-4 max-h-[60vh] overflow-y-auto">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-500">操作用户</label>
                <p class="text-sm text-gray-900">{{ currentLog?.username || '未知用户' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">操作类型</label>
                <p class="text-sm text-gray-900">{{ getActionLabel(currentLog?.action) || currentLog?.action }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">资源类型</label>
                <p class="text-sm text-gray-900">{{ resourceTypeLabel[currentLog?.module] || currentLog?.module }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">操作时间</label>
                <p class="text-sm text-gray-900">{{ formatDateTime(currentLog?.created_at) }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">IP地址</label>
                <p class="text-sm text-gray-900">{{ currentLog?.ip_address || '-' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-500">User Agent</label>
                <p class="text-sm text-gray-900 break-all">{{ currentLog?.user_agent || '-' }}</p>
              </div>
            </div>
            <div>
              <label class="text-sm font-medium text-gray-500">操作描述</label>
              <p class="text-sm text-gray-900 mt-1">{{ currentLog?.description || '无描述' }}</p>
            </div>
            <div v-if="currentLog?.request_data">
              <label class="text-sm font-medium text-gray-500">请求数据</label>
              <pre class="text-xs text-gray-600 bg-gray-50 p-3 rounded mt-1 overflow-x-auto">{{ currentLog.request_data }}</pre>
            </div>
            <div v-if="currentLog?.response_body">
              <label class="text-sm font-medium text-gray-500">响应数据</label>
              <pre class="text-xs text-gray-600 bg-gray-50 p-3 rounded mt-1 overflow-x-auto">{{ currentLog.response_body }}</pre>
            </div>
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
import TailCard from '@/components/ui/TailCard.vue'
import { getOperationLogs, getOperationLogStats } from '@/api/audit'

const logs = ref<any[]>([])
const loading = ref(false)
const stats = ref<any>({ total_count: 0, today_count: 0, success_count: 0, fail_count: 0, action_stats: [] })
const actionStatsChartLoaded = ref(false)
const currentLog = ref<any>(null)
const showDetailDialog = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 15,
  total: 0
})

const searchParams = reactive({
  keyword: '',
  action_type: '',
  resource_type: '',
  date_range: ''
})

const actionTypeLabel: Record<string, string> = {
  create: '创建',
  update: '更新',
  delete: '删除',
  query: '查询',
  login: '登录'
}

const getActionLabel = (action: string) => {
  const labels: Record<string, string> = {
    create: '创建',
    update: '更新',
    delete: '删除',
    query: '查询',
    login: '登录',
    publish: '发布',
    approve: '审批',
    reject: '驳回',
    'batch_delete': '批量删除',
    import: '导入',
    export: '导出',
    list: '列表',
    stats: '统计',
    tree: '树形',
    view: '查看'
  }
  return labels[action] || action
}

const resourceTypeLabel: Record<string, string> = {
  article: '文章',
  video: '视频',
  user: '用户',
  menu: '菜单',
  role: '角色',
  system: '系统'
}

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
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

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN')
  } catch {
    return dateStr
  }
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadLogs()
  }, 300)
}

const handleSearch = () => {
  pagination.page = 1
  loadLogs()
}

const resetSearch = () => {
  Object.assign(searchParams, { keyword: '', action_type: '', resource_type: '', date_range: '' })
  handleSearch()
}

const loadLogs = async () => {
  loading.value = true
  try {
    // 将前端的 resource_type 映射到后端的 module 字段
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchParams.keyword) params.keyword = searchParams.keyword
    if (searchParams.action_type) params.action = searchParams.action_type
    if (searchParams.resource_type) params.module = searchParams.resource_type
    if (searchParams.date_range) {
      // 日期范围需要转换为 start_date 和 end_date
      params.end_date = searchParams.date_range
    }

    const res = await getOperationLogs(params)
    console.log('操作日志响应:', res.data)
    // 响应拦截器已经解包了外层，res.data 直接指向后端 response.data
    logs.value = res.data?.list || res.data?.items || []
    pagination.total = res.data?.total || res.data?.total_count || 0
  } catch (error) {
    message.error('获取操作日志失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getOperationLogStats()
    console.log('操作日志统计响应:', res.data)
    // 响应拦截器已经解包了外层，res.data 直接指向后端 response.data
    const rawData = res.data
    if (rawData) {
      stats.value = {
        total_count: rawData.total_count || 0,
        today_count: rawData.today_count || 0,
        success_count: rawData.success_count || 0,
        fail_count: rawData.fail_count || 0,
        action_stats: rawData.action_stats || []
      }
      actionStatsChartLoaded.value = true
    }
  } catch (error) {
    console.error('获取操作日志统计失败:', error)
  }
}

const viewDetail = (log: any) => {
  currentLog.value = log
  showDetailDialog.value = true
}

const handleExport = () => {
  message.info('导出功能开发中...')
}

onMounted(() => {
  Promise.all([loadLogs(), loadStats()])
})
</script>