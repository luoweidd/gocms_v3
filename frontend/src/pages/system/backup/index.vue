<template>
  <div class="max-w-7xl mx-auto px-4 py-6">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-gray-900">{{ stats.totalCount || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">总备份数</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-green-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-green-600">{{ stats.successCount || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">成功备份</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-red-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-red-600">{{ stats.failCount || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">失败备份</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-2xl font-bold text-gray-900">{{ formatFileSize(stats.totalSize || 0) }}</div>
          <div class="text-sm text-gray-500 mt-1">总存储空间</div>
        </div>
      </div>
    </div>

    <!-- 操作栏 -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
      <div class="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3">
        <div class="flex items-center gap-2">
          <button 
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors flex items-center gap-2 disabled:opacity-50"
            @click="handleCreateBackup" 
            :disabled="creating"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            创建备份
          </button>
          <button 
            class="px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors flex items-center gap-2"
            @click="fetchData"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            刷新
          </button>
        </div>
        <div class="flex items-center">
          <input
            v-model="keyword"
            type="text"
            placeholder="搜索备份..."
            class="w-64 px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            @input="handleSearch"
          />
        </div>
      </div>
    </div>

    <!-- 备份列表 -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200" v-loading="loading">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">名称</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">大小</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">表数量</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="row in backupList" :key="row.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <svg class="w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <span class="text-sm text-gray-900 font-medium">{{ row.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-medium bg-blue-100 text-blue-800 rounded-full">
                  {{ BACKUP_TYPE_MAP[row.backup_type] || row.backup_type }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="{
                    'bg-gray-100 text-gray-800': row.status === 0,
                    'bg-green-100 text-green-800': row.status === 1,
                    'bg-red-100 text-red-800': row.status === 2
                  }"
                >
                  {{ BACKUP_STATUS_MAP[row.status] || '未知' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatFileSize(row.file_size) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ row.tables_count }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatTime(row.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex items-center gap-2">
                  <button
                    v-if="String(row.status) === '1'"
                    class="text-blue-600 hover:text-blue-900 transition-colors"
                    @click="handleDownload(row)"
                  >
                    下载
                  </button>
                  <button
                    v-if="String(row.status) === '1'"
                    class="text-yellow-600 hover:text-yellow-900 transition-colors"
                    @click="handleRestore(row)"
                  >
                    恢复
                  </button>
                  <button 
                    v-if="row.id !== undefined"
                    class="text-red-600 hover:text-red-900 transition-colors"
                    @click="handleDelete(row)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && backupList.length === 0">
              <td colspan="7" class="px-6 py-12 text-center text-gray-500">
                暂无备份数据
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-6 py-4 border-t border-gray-200 flex items-center justify-between" v-if="pagination.total > 0">
        <div class="text-sm text-gray-500">
          共 {{ pagination.total }} 条记录
        </div>
        <div class="flex items-center gap-2">
          <select
            v-model="pagination.pageSize"
            @change="fetchData"
            class="px-2 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option :value="10">10条/页</option>
            <option :value="20">20条/页</option>
            <option :value="50">50条/页</option>
          </select>
          <button
            @click="pagination.page > 1 && (pagination.page--, fetchData())"
            :disabled="pagination.page === 1"
            class="px-3 py-1 border border-gray-300 rounded text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            上一页
          </button>
          <span class="text-sm text-gray-600 px-2">
            {{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}
          </span>
          <button
            @click="pagination.page < Math.ceil(pagination.total / pagination.pageSize) && (pagination.page++, fetchData())"
            :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
            class="px-3 py-1 border border-gray-300 rounded text-sm hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
          >
            下一页
          </button>
        </div>
      </div>
    </div>

    <!-- 创建备份对话框 -->
    <teleport to="body">
      <div v-if="showCreateDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showCreateDialog = false">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        <div class="relative bg-white rounded-xl shadow-2xl w-full max-w-md p-6 animate-fade-in">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">创建数据备份</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">备份名称</label>
              <input 
                v-model="createForm.name" 
                type="text" 
                placeholder="可留空，自动生成" 
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">备份类型</label>
              <div class="flex gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input 
                    type="radio" 
                    value="full" 
                    v-model="createForm.backup_type"
                    class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                  />
                  <span class="text-sm text-gray-700">全量备份</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input 
                    type="radio" 
                    value="incremental" 
                    v-model="createForm.backup_type"
                    class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
                  />
                  <span class="text-sm text-gray-700">增量备份</span>
                </label>
              </div>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200">
            <button 
              @click="showCreateDialog = false"
              class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm font-medium"
            >
              取消
            </button>
            <button 
              @click="handleConfirmCreate" 
              :disabled="creating"
              class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors text-sm font-medium disabled:opacity-50"
            >
              {{ creating ? '创建中...' : '创建备份' }}
            </button>
          </div>
        </div>
      </div>

    <!-- 删除确认对话框 -->
    <div v-if="showDeleteDialog && deleteTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showDeleteDialog = false">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
      <div class="relative bg-white rounded-xl shadow-2xl w-full max-w-md p-6 animate-fade-in">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">确认删除备份</h3>
        <div class="space-y-4">
          <div class="bg-red-50 border border-red-200 rounded-lg p-4">
            <div class="flex items-start gap-3">
              <svg class="w-5 h-5 text-red-600 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div class="text-sm text-red-800">
                <p class="font-medium mb-1">警告</p>
                <p>删除操作将永久移除该备份文件，请确保已做好当前数据的备份！</p>
              </div>
            </div>
          </div>
          <div>
            <p class="text-sm text-gray-600">将要删除的备份：<span class="font-medium text-gray-900">{{ deleteTarget.name }}</span></p>
            <p class="text-sm text-gray-600 mt-1">创建时间：<span class="font-medium text-gray-900">{{ formatTime(deleteTarget.created_at) }}</span></p>
          </div>
        </div>
        <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200">
          <button 
            @click="showDeleteDialog = false"
            class="px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm font-medium"
          >
            取消
          </button>
          <button 
            @click="handleConfirmDelete" 
            :disabled="deleting"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors text-sm font-medium disabled:opacity-50"
          >
            {{ deleting ? '删除中...' : '确认删除' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 恢复确认对话框 -->
    <div v-if="showRestoreDialog && restoreTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="showRestoreDialog = false">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        <div class="relative bg-white rounded-xl shadow-2xl w-full max-w-md p-6 animate-fade-in">
          <h3 class="text-xl font-semibold text-gray-900 mb-4">确认恢复备份</h3>
          <div class="space-y-4">
            <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div class="flex items-start gap-3">
                <svg class="w-5 h-5 text-yellow-600 mt-0.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div class="text-sm text-yellow-800">
                  <p class="font-medium mb-1">警告</p>
                  <p>恢复操作将覆盖当前数据库中的所有数据，请确保已做好当前数据的备份！</p>
                </div>
              </div>
            </div>
            <div>
              <p class="text-sm text-gray-600">将要恢复的备份：<span class="font-medium text-gray-900">{{ restoreTarget.name }}</span></p>
              <p class="text-sm text-gray-600 mt-1">创建时间：<span class="font-medium text-gray-900">{{ formatTime(restoreTarget.created_at) }}</span></p>
            </div>
          </div>
          <div class="flex justify-end mt-6 pt-4 border-t border-gray-200">
          </div>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message, ElMessageBox } from '@/utils/message'
import {
  createBackup,
  getBackupList,
  deleteBackup,
  downloadBackup,
  restoreBackup,
  getBackupStats,
  BACKUP_TYPE_MAP,
  BACKUP_STATUS_MAP
} from '@/api/backup'

// 响应式数据
const loading = ref(false)
const creating = ref(false)
const restoring = ref(false)
const keyword = ref('')
const backupList = ref<any[]>([])
const stats = ref({
  totalCount: 0,
  successCount: 0,
  failCount: 0,
  totalSize: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 创建备份对话框
const showCreateDialog = ref(false)
const createForm = reactive({
  name: '',
  backup_type: 'full'
})

// 删除对话框
const showDeleteDialog = ref(false)
const deleteTarget = ref<any>(null)
const deleting = ref(false)

// 恢复对话框
const showRestoreDialog = ref(false)
const restoreTarget = ref<any>(null)

// 常量映射
const STATUS_MAP = { ...BACKUP_STATUS_MAP }
const TYPE_MAP = { ...BACKUP_TYPE_MAP }

// 方法
const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatTime = (time: string | undefined): string => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const res = await getBackupList(params)
    console.log('备份列表响应:', res.data)
    backupList.value = res.data?.data?.list || res.data?.list || []
    pagination.total = res.data?.data?.total || res.data?.total || 0
  } catch (error) {
    console.error('获取备份列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await getBackupStats()
    console.log('备份统计响应:', res.data)
    // 后端返回的是下划线命名（total_count, success_count, fail_count, total_size）
    // 前端期望的是驼峰命名（totalCount, successCount, failCount, totalSize）
    const rawData = res.data?.data || res.data
    if (rawData) {
      stats.value = {
        totalCount: rawData.total_count ?? rawData.totalCount ?? 0,
        successCount: rawData.success_count ?? rawData.successCount ?? 0,
        failCount: rawData.fail_count ?? rawData.failCount ?? 0,
        totalSize: rawData.total_size ?? rawData.totalSize ?? 0
      }
    }
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

// 创建备份
const handleCreateBackup = () => {
  createForm.name = ''
  createForm.backup_type = 'full'
  showCreateDialog.value = true
}

const handleConfirmCreate = async () => {
  creating.value = true
  try {
    await createBackup(createForm)
    message.success('备份已创建，请稍后查看状态')
    showCreateDialog.value = false
    fetchData()
  } catch (error) {
    console.error('创建备份失败:', error)
  } finally {
    creating.value = false
  }
}

// 下载备份
const handleDownload = async (row: any) => {
  try {
    const res = await downloadBackup(row.id)
    // res.data 是 blob 对象 (axios responseType: 'blob' 配置)
    console.log('下载响应:', res)
    console.log('blob 数据:', res.data)
    const blob = new Blob([res.data], { type: 'application/sql' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    let fileName = 'backup.sql'
    if (row.file_path) {
      const parts = row.file_path.split('/')
      fileName = parts[parts.length - 1]
    }
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    message.success('下载成功')
  } catch (error) {
    console.error('下载失败:', error)
    message.error('下载失败')
  }
}

// 恢复备份
const handleRestore = (row: any) => {
  restoreTarget.value = row
  showRestoreDialog.value = true
}

const handleConfirmRestore = async () => {
  if (!restoreTarget.value) return
  restoring.value = true
  try {
    await ElMessageBox.confirm(
      `确定要恢复备份 "${restoreTarget.value.name}" 吗？此操作将覆盖当前数据！`,
      '危险操作确认',
      { type: 'error' as any }
    )
    await restoreBackup(restoreTarget.value.id)
    message.success('恢复请求已提交，请稍后查看状态')
    showRestoreDialog.value = false
    fetchData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('恢复失败:', error)
    }
  } finally {
    restoring.value = false
  }
}

// 删除备份
const handleDelete = (row: any) => {
  deleteTarget.value = row
  showDeleteDialog.value = true
}

const handleConfirmDelete = async () => {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await deleteBackup(deleteTarget.value.id)
    message.success('删除成功')
    showDeleteDialog.value = false
    await fetchData()
    await fetchStats()
  } catch (error) {
    console.error('删除失败:', error)
  } finally {
    deleting.value = false
  }
}

// 生命周期
onMounted(() => {
  fetchData()
  fetchStats()
})
</script>

<style scoped lang="scss">
.v-loading {
  position: relative;
}

.v-loading::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fade-in 0.2s ease-out;
}
</style>