<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800">定时任务执行日志</h1>
      <p class="text-gray-500 mt-1">查看定时任务的执行历史记录和结果</p>
    </div>

    <!-- 筛选栏 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex items-center gap-4 flex-wrap">
        <input
          v-model="searchKeyword"
          type="text"
          placeholder="搜索任务名称..."
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
        <select v-model="statusFilter" class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500">
          <option value="">全部状态</option>
          <option value="success">成功</option>
          <option value="failed">失败</option>
          <option value="running">执行中</option>
        </select>
        <input
          v-model="dateRange.start"
          type="date"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          placeholder="开始日期"
        />
        <span class="text-gray-500">至</span>
        <input
          v-model="dateRange.end"
          type="date"
          class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
          placeholder="结束日期"
        />
        <button
          @click="handleRefresh"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          刷新
        </button>
        <button
          @click="handleClearLogs"
          class="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
        >
          清空日志
        </button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">总执行次数</p>
            <p class="text-2xl font-bold text-gray-900">{{ stats.total }}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">成功次数</p>
            <p class="text-2xl font-bold text-green-600">{{ stats.success }}</p>
          </div>
          <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">失败次数</p>
            <p class="text-2xl font-bold text-red-600">{{ stats.failed }}</p>
          </div>
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-gray-500">成功率</p>
            <p class="text-2xl font-bold text-indigo-600">{{ stats.success_rate }}%</p>
          </div>
          <div class="w-12 h-12 bg-indigo-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- 日志列表 -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">任务名称</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">执行状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">执行时长</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">开始时间</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">结束时间</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">执行结果</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="log in filteredLogs" :key="log.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ log.task_name }}</div>
              <div class="text-sm text-gray-500">{{ log.task_id }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full',
                  getStatusClass(log.status)
                ]"
              >
                {{ getStatusText(log.status) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ log.duration ? log.duration + 'ms' : '-' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatTime(log.start_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ log.end_at ? formatTime(log.end_at) : '-' }}
            </td>
            <td class="px-6 py-4 max-w-xs text-sm">
              <div class="truncate text-gray-500" v-if="log.status === 'success'">
                {{ log.output || '执行成功' }}
              </div>
              <div class="truncate text-red-600 font-medium" v-else>
                {{ log.error || '未知错误' }}
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <button
                @click="handleViewDetail(log)"
                class="text-indigo-600 hover:text-indigo-900"
              >
                详情
              </button>
            </td>
          </tr>
          <tr v-if="filteredLogs.length === 0">
            <td colspan="7" class="px-6 py-12 text-center text-gray-500">
              <svg class="w-12 h-12 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              暂无执行日志
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 分页 -->
      <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between">
        <div class="text-sm text-gray-500">
          显示 {{ filteredLogs.length }} 条记录
        </div>
        <div class="flex items-center gap-2">
          <button
            :disabled="currentPage === 1"
            @click="currentPage--"
            class="px-3 py-1 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            上一页
          </button>
          <span class="text-sm text-gray-500">
            第 {{ currentPage }} / {{ totalPages }} 页
          </span>
          <button
            :disabled="currentPage === totalPages"
            @click="currentPage++"
            class="px-3 py-1 border border-gray-300 rounded-lg text-sm disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
          >
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- 日志详情弹窗 -->
    <div v-if="showDetailModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-3xl m-4 max-h-[90vh] flex flex-col">
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900">执行日志详情</h3>
          <button @click="showDetailModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="px-6 py-4 overflow-y-auto flex-1">
          <div class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-500">任务名称</label>
                <p class="mt-1 text-sm text-gray-900">{{ detailLog?.task_name }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-500">执行状态</label>
                <p class="mt-1">
                  <span
                    :class="[
                      'px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full',
                      detailLog ? getStatusClass(detailLog.status) : ''
                    ]"
                  >
                    {{ detailLog ? getStatusText(detailLog.status) : '' }}
                  </span>
                </p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-500">执行时长</label>
                <p class="mt-1 text-sm text-gray-900">{{ detailLog?.duration ? detailLog.duration + 'ms' : '-' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-500">触发方式</label>
                <p class="mt-1 text-sm text-gray-900">{{ detailLog?.trigger_type || '定时触发' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-500">开始时间</label>
                <p class="mt-1 text-sm text-gray-900">{{ detailLog?.start_at ? formatTime(detailLog.start_at) : '-' }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-500">结束时间</label>
                <p class="mt-1 text-sm text-gray-900">{{ detailLog?.end_at ? formatTime(detailLog.end_at) : '-' }}</p>
              </div>
            </div>
            <div v-if="detailLog?.output">
              <label class="block text-sm font-medium text-gray-500 mb-1">输出日志</label>
              <pre class="p-4 bg-gray-50 rounded-lg text-sm text-green-600 overflow-auto max-h-48">{{ detailLog.output }}</pre>
            </div>
            <div v-if="detailLog?.error">
              <label class="block text-sm font-medium text-gray-500 mb-1">错误信息</label>
              <pre class="p-4 bg-red-50 rounded-lg text-sm text-red-600 overflow-auto max-h-48">{{ detailLog.error }}</pre>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end">
          <button
            @click="showDetailModal = false"
            class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200"
          >
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface CronLog {
  id: number
  task_id: number
  task_name: string
  status: 'success' | 'failed' | 'running'
  duration: number | null
  start_at: string
  end_at: string | null
  output: string | null
  error: string | null
  trigger_type: string | null
}

// 模拟数据
const mockLogs = ref<CronLog[]>([
  {
    id: 1,
    task_id: 1,
    task_name: '内容清理任务',
    status: 'success',
    duration: 2350,
    start_at: '2026-09-11 00:00:00',
    end_at: '2026-09-11 00:00:02',
    output: '成功清理过期数据 15 条',
    error: null,
    trigger_type: 'scheduled'
  },
  {
    id: 2,
    task_id: 2,
    task_name: '数据备份任务',
    status: 'success',
    duration: 45200,
    start_at: '2026-09-11 02:00:00',
    end_at: '2026-09-11 02:00:45',
    output: '数据库备份完成，文件大小: 256MB',
    error: null,
    trigger_type: 'scheduled'
  },
  {
    id: 3,
    task_id: 4,
    task_name: '统计报告生成',
    status: 'failed',
    duration: 30000,
    start_at: '2026-09-10 08:00:00',
    end_at: '2026-09-10 08:00:30',
    output: null,
    error: '网络连接超时: 无法访问报告生成服务',
    trigger_type: 'scheduled'
  },
  {
    id: 4,
    task_id: 1,
    task_name: '内容清理任务',
    status: 'success',
    duration: 1890,
    start_at: '2026-09-10 00:00:00',
    end_at: '2026-09-10 00:00:01',
    output: '成功清理过期数据 8 条',
    error: null,
    trigger_type: 'scheduled'
  },
  {
    id: 5,
    task_id: 3,
    task_name: '缓存刷新任务',
    status: 'failed',
    duration: 5000,
    start_at: '2026-09-09 14:30:00',
    end_at: '2026-09-09 14:30:05',
    output: null,
    error: 'Redis连接失败: ECONNREFUSED',
    trigger_type: 'manual'
  },
  {
    id: 6,
    task_id: 2,
    task_name: '数据备份任务',
    status: 'success',
    duration: 38100,
    start_at: '2026-09-09 02:00:00',
    end_at: '2026-09-09 02:00:38',
    output: '数据库备份完成，文件大小: 248MB',
    error: null,
    trigger_type: 'scheduled'
  }
])

const searchKeyword = ref('')
const statusFilter = ref('')
const dateRange = ref({
  start: '',
  end: ''
})
const showDetailModal = ref(false)
const detailLog = ref<CronLog | null>(null)
const currentPage = ref(1)
const totalPages = ref(1)

const stats = computed(() => {
  const logs = mockLogs.value
  const total = logs.length
  const success = logs.filter(l => l.status === 'success').length
  const failed = logs.filter(l => l.status === 'failed').length
  const success_rate = total > 0 ? Math.round((success / total) * 100) : 0
  return { total, success, failed, success_rate }
})

const filteredLogs = computed(() => {
  return mockLogs.value.filter(log => {
    const matchKeyword = !searchKeyword.value || log.task_name.toLowerCase().includes(searchKeyword.value.toLowerCase())
    const matchStatus = !statusFilter.value || log.status === statusFilter.value
    return matchKeyword && matchStatus
  })
})

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('zh-CN')
}

function getStatusClass(status: string): string {
  switch (status) {
    case 'success':
      return 'bg-green-100 text-green-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    case 'running':
      return 'bg-blue-100 text-blue-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

function getStatusText(status: string): string {
  switch (status) {
    case 'success':
      return '成功'
    case 'failed':
      return '失败'
    case 'running':
      return '执行中'
    default:
      return '未知'
  }
}

function handleViewDetail(log: CronLog) {
  detailLog.value = log
  showDetailModal.value = true
}

function handleRefresh() {
  // 刷新数据逻辑
  console.log('刷新日志数据')
}

function handleClearLogs() {
  if (confirm('确定要清空所有执行日志吗？此操作不可恢复。')) {
    mockLogs.value = []
  }
}
</script>

<style scoped>
/* 可选：添加自定义样式 */
</style>