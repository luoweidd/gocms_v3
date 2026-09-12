<template>
  <div class="p-6 bg-gray-50 min-h-screen">
    <!-- 页面头部 -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-800">定时任务管理</h1>
      <p class="text-gray-500 mt-1">管理和配置系统定时任务</p>
    </div>

    <!-- 操作栏 -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-4">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索任务名称..."
            class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          />
          <select v-model="statusFilter" class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500">
            <option value="">全部状态</option>
            <option value="1">启用</option>
            <option value="0">禁用</option>
          </select>
        </div>
        <button
          @click="handleCreate"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          创建任务
        </button>
      </div>
    </div>

    <!-- 任务列表 -->
    <div class="bg-white rounded-lg shadow overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">任务名称</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">表达式</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">执行时间</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">最后执行</th>
            <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="task in filteredTasks" :key="task.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="flex-shrink-0 w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center">
                  <svg class="w-4 h-4 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">{{ task.name }}</div>
                  <div class="text-sm text-gray-500">{{ task.handler }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <code class="px-2 py-1 bg-gray-100 rounded text-sm text-indigo-600">{{ task.cron_expression }}</code>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatNextRun(task.next_run_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full',
                  task.status === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                ]"
              >
                {{ task.status === 1 ? '启用' : '禁用' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ task.last_run_at ? formatTime(task.last_run_at) : '从未执行' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <div class="flex items-center justify-end gap-2">
                <button
                  @click="handleToggleStatus(task)"
                  :class="task.status === 1 ? 'text-yellow-600 hover:text-yellow-900' : 'text-green-600 hover:text-green-900'"
                  class="flex items-center gap-1"
                >
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path v-if="task.status === 1" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.478A1 1 0 0010 9.81v4.38a1 1 0 001.555.828l3.197-2.478a1 1 0 000-1.656z" />
                  </svg>
                  {{ task.status === 1 ? '停止' : '启动' }}
                </button>
                <button
                  @click="handleEdit(task)"
                  class="text-indigo-600 hover:text-indigo-900"
                >
                  编辑
                </button>
                <button
                  @click="handleDelete(task)"
                  class="text-red-600 hover:text-red-900"
                >
                  删除
                </button>
                <button
                  @click="handleRunOnce(task)"
                  class="text-blue-600 hover:text-blue-900"
                >
                  执行
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="filteredTasks.length === 0">
            <td colspan="6" class="px-6 py-12 text-center text-gray-500">
              <svg class="w-12 h-12 mx-auto text-gray-400 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
              </svg>
              暂无定时任务数据
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 创建/编辑弹窗 -->
    <div v-if="showModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl m-4">
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900">{{ editingTask ? '编辑任务' : '创建任务' }}</h3>
          <button @click="closeModal" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="px-6 py-4">
          <form @submit.prevent="handleSubmit" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">任务名称</label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Cron表达式</label>
              <input
                v-model="form.cron_expression"
                type="text"
                required
                placeholder="例如: 0 0 * * * (每天午夜执行)"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 font-mono"
              />
              <p class="text-xs text-gray-500 mt-1">格式: 秒 分 时 日 月 星期</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">处理器类型</label>
              <select
                v-model="form.handler_type"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="">请选择</option>
                <option value="http">HTTP请求</option>
                <option value="script">脚本执行</option>
                <option value="internal">内部任务</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">处理器配置</label>
              <textarea
                v-model="form.handler_config"
                required
                rows="4"
                placeholder='{"url": "https://api.example.com/cron", "method": "POST"}'
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 font-mono"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">超时时间（秒）</label>
              <input
                v-model.number="form.timeout"
                type="number"
                min="1"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div class="flex items-center gap-4">
              <label class="flex items-center">
                <input
                  v-model="form.enabled"
                  type="checkbox"
                  class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                />
                <span class="ml-2 text-sm text-gray-700">启用任务</span>
              </label>
              <label class="flex items-center">
                <input
                  v-model="form.allow_concurrent"
                  type="checkbox"
                  class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                />
                <span class="ml-2 text-sm text-gray-700">允许并发</span>
              </label>
            </div>
          </form>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end gap-3">
          <button
            @click="closeModal"
            class="px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
          >
            取消
          </button>
          <button
            @click="handleSubmit"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
          >
            {{ editingTask ? '保存' : '创建' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface CronTask {
  id: number
  name: string
  cron_expression: string
  handler: string
  status: number
  last_run_at: string | null
  next_run_at: string | null
}

// 模拟数据
const mockTasks = ref<CronTask[]>([
  {
    id: 1,
    name: '内容清理任务',
    cron_expression: '0 0 * * *',
    handler: 'internal:cron/clean_content',
    status: 1,
    last_run_at: '2026-09-11 00:00:00',
    next_run_at: '2026-09-12 00:00:00'
  },
  {
    id: 2,
    name: '数据备份任务',
    cron_expression: '0 2 * * *',
    handler: 'http://backup.internal/api/backup',
    status: 1,
    last_run_at: '2026-09-11 02:00:00',
    next_run_at: '2026-09-12 02:00:00'
  },
  {
    id: 3,
    name: '缓存刷新任务',
    cron_expression: '*/30 * * * *',
    handler: 'internal:cron/refresh_cache',
    status: 0,
    last_run_at: '2026-09-10 14:30:00',
    next_run_at: null
  },
  {
    id: 4,
    name: '统计报告生成',
    cron_expression: '0 8 * * 1',
    handler: 'script:/opt/scripts/report.sh',
    status: 1,
    last_run_at: '2026-09-08 08:00:00',
    next_run_at: '2026-09-15 08:00:00'
  }
])

const searchKeyword = ref('')
const statusFilter = ref('')
const showModal = ref(false)
const editingTask = ref<CronTask | null>(null)
const form = ref({
  name: '',
  cron_expression: '',
  handler_type: '',
  handler_config: '',
  timeout: 30,
  enabled: true,
  allow_concurrent: false
})

const filteredTasks = computed(() => {
  return mockTasks.value.filter(task => {
    const matchKeyword = !searchKeyword.value || task.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
    const matchStatus = !statusFilter.value || task.status.toString() === statusFilter.value
    return matchKeyword && matchStatus
  })
})

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString('zh-CN')
}

function formatNextRun(dateStr: string | null): string {
  if (!dateStr) return '未设置'
  return new Date(dateStr).toLocaleString('zh-CN')
}

function handleCreate() {
  editingTask.value = null
  form.value = {
    name: '',
    cron_expression: '',
    handler_type: '',
    handler_config: '{}',
    timeout: 30,
    enabled: true,
    allow_concurrent: false
  }
  showModal.value = true
}

function handleEdit(task: CronTask) {
  editingTask.value = task
  form.value = {
    name: task.name,
    cron_expression: task.cron_expression,
    handler_type: 'internal',
    handler_config: JSON.stringify({ handler: task.handler }, null, 2),
    timeout: 30,
    enabled: task.status === 1,
    allow_concurrent: false
  }
  showModal.value = true
}

function handleDelete(task: CronTask) {
  if (confirm(`确定要删除任务 "${task.name}" 吗？`)) {
    mockTasks.value = mockTasks.value.filter(t => t.id !== task.id)
  }
}

function handleToggleStatus(task: CronTask) {
  task.status = task.status === 1 ? 0 : 1
  if (task.status === 0) {
    task.next_run_at = null
  }
}

function handleRunOnce(task: CronTask) {
  if (confirm(`确定要立即执行任务 "${task.name}" 吗？`)) {
    task.last_run_at = new Date().toISOString().slice(0, 19).replace('T', ' ')
  }
}

function closeModal() {
  showModal.value = false
}

function handleSubmit() {
  if (editingTask.value) {
    // 编辑逻辑
    const task = mockTasks.value.find(t => t.id === editingTask.value!.id)
    if (task) {
      task.name = form.value.name
      task.cron_expression = form.value.cron_expression
      task.status = form.value.enabled ? 1 : 0
    }
  } else {
    // 创建逻辑
    const newTask: CronTask = {
      id: Date.now(),
      name: form.value.name,
      cron_expression: form.value.cron_expression,
      handler: form.value.handler_config,
      status: form.value.enabled ? 1 : 0,
      last_run_at: null,
      next_run_at: '2026-09-13 00:00:00'
    }
    mockTasks.value.push(newTask)
  }
  closeModal()
}
</script>

<style scoped>
/* 可选：添加自定义样式 */
</style>