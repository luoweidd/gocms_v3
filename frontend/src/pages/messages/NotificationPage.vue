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
          <h1 class="text-2xl font-bold text-gray-900 mb-1">我的通知</h1>
          <p class="text-sm text-gray-500">查看您的个人消息通知</p>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-3xl font-bold text-primary">{{ stats.total_count || 0 }}</div>
        <div class="text-sm text-gray-500 mt-1">全部通知</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-3xl font-bold text-red-600">{{ stats.unread_count || 0 }}</div>
        <div class="text-sm text-gray-500 mt-1">未读通知</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-3xl font-bold text-green-600">{{ stats.system_count || 0 }}</div>
        <div class="text-sm text-gray-500 mt-1">系统通知</div>
      </div>
      <div class="bg-white rounded-lg shadow p-4">
        <div class="text-3xl font-bold text-purple-600">{{ stats.audit_count || 0 }}</div>
        <div class="text-sm text-gray-500 mt-1">审核通知</div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="mb-6">
      <div class="flex gap-1 border-b border-gray-200">
        <button
          v-for="tab in tabs"
          :key="tab.key"
          @click="activeTab = tab.key"
          :class="[
            'px-4 py-2 text-sm font-medium transition-colors rounded-t-lg',
            activeTab === tab.key
              ? 'bg-white text-primary border-b-2 border-primary'
              : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
          ]"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <!-- Message List -->
    <TailCard class="overflow-hidden">
      <div v-if="loading" class="py-12 text-center text-gray-400">
        <svg class="w-8 h-8 mx-auto mb-2 animate-spin text-primary" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        加载中...
      </div>
      
      <div v-else-if="notifications.length === 0" class="py-12 text-center text-gray-400">
        <svg class="w-16 h-16 mx-auto mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <p class="text-sm">暂无通知</p>
      </div>

      <div v-else class="divide-y divide-gray-100">
        <div 
          v-for="notification in notifications" 
          :key="notification.id"
          @click="markAsRead(notification.id)"
          :class="[
            'px-6 py-4 hover:bg-gray-50 cursor-pointer transition-colors',
            !notification.is_read ? 'bg-blue-50/30' : ''
          ]"
        >
          <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
                 :class="getIconBg(notification.message.message_type)">
              <svg class="w-5 h-5" :class="getIconColor(notification.message.message_type)" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path v-if="notification.message.message_type === 'system'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else-if="notification.message.message_type === 'audit_passed'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else-if="notification.message.message_type === 'audited_failed'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </div>
            
            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-gray-900">{{ notification.message.title }}</span>
                <span v-if="!notification.is_read" class="w-2 h-2 bg-primary rounded-full"></span>
              </div>
              <p class="text-sm text-gray-500 mt-1 truncate">{{ notification.message.content }}</p>
              <div class="flex items-center gap-4 mt-2">
                <span class="text-xs text-gray-400">{{ formatTime(notification.created_at) }}</span>
                <span v-if="notification.message.priority >= 2" class="text-xs px-2 py-0.5 rounded-full bg-red-100 text-red-600">紧急</span>
                <span v-else-if="notification.message.priority === 1" class="text-xs px-2 py-0.5 rounded-full bg-yellow-100 text-yellow-600">重要</span>
              </div>
            </div>

            <!-- Actions -->
            <button 
              v-if="notification.is_read"
              @click.stop="deleteNotification(notification.id)"
              class="shrink-0 p-1 text-gray-400 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="notifications.length > 0" class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条通知</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="xs" :disabled="pagination.page <= 1" @click="changePage(pagination.page - 1)">上一页</TailButton>
          <TailButton 
            v-for="p in totalPages" 
            :key="p"
            type="default" 
            size="xs"
            class="w-8 h-8 p-0 flex items-center justify-center"
            :class="p === pagination.page ? 'bg-primary text-white border-primary' : ''"
            @click="changePage(p)">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="xs" :disabled="pagination.page >= totalPages" @click="changePage(pagination.page + 1)">下一页</TailButton>
        </div>
      </div>

      <!-- Mark all as read -->
      <div v-if="hasUnread" class="flex items-center justify-center px-6 py-4 border-t border-gray-200">
        <button 
          @click="markAllAsRead"
          class="text-sm text-primary hover:text-primary-dark font-medium"
        >
          全部标为已读
        </button>
      </div>
    </TailCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import { getMyNotifications, getMyNotificationStats, markAsRead as apiMarkAsRead, markAllAsRead as apiMarkAllAsRead } from '@/api/notification'

interface NotificationItem {
  id: number
  message_id: number
  is_read: number
  read_at?: string
  message: {
    id: number
    title: string
    content: string
    message_type: string
    priority: number
    created_at: string
  }
  created_at: string
}

interface NotificationStats {
  total_count: number
  unread_count: number
  read_count: number
  system_count: number
  publish_count: number
  comment_count: number
  audit_count: number
}

const tabs = [
  { key: 'all', label: '全部通知' },
  { key: 'system', label: '系统通知' },
  { key: 'notification', label: '其他通知' }
] as const

type TabKey = typeof tabs[number]['key']

const loading = ref(false)
const notifications = ref<NotificationItem[]>([])
const stats = ref<NotificationStats>({
  total_count: 0,
  unread_count: 0,
  read_count: 0,
  system_count: 0,
  publish_count: 0,
  comment_count: 0,
  audit_count: 0
})

const activeTab = ref<TabKey>('all')
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const hasUnread = computed(() => stats.value.unread_count > 0)

onMounted(() => {
  fetchNotifications()
  fetchStats()
})

const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await getMyNotifications({
      page: pagination.page,
      page_size: pagination.pageSize,
      message_type: activeTab.value === 'all' ? '' : activeTab.value
    })
    notifications.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    console.error('获取通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await getMyNotificationStats()
    stats.value = res.data || stats.value
  } catch (error) {
    console.error('获取通知统计失败:', error)
  }
}

const markAsRead = async (id: number) => {
  try {
    await apiMarkAsRead(id, 1)
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.is_read = 1
    }
    fetchStats()
  } catch (error) {
    console.error('标记已读失败:', error)
  }
}

const markAllAsRead = async () => {
  try {
    await apiMarkAllAsRead()
    notifications.value.forEach(n => n.is_read = 1)
    message.success('全部已标为已读')
    fetchStats()
  } catch (error) {
    console.error('全部标记失败:', error)
  }
}

const deleteNotification = async (id: number) => {
  try {
    notifications.value = notifications.value.filter(n => n.id !== id)
    pagination.total--
    message.success('删除成功')
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const changePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  pagination.page = page
  fetchNotifications()
}

const getIconBg = (type: string) => {
  const map: Record<string, string> = {
    system: 'bg-blue-100 text-blue-600',
    publish: 'bg-green-100 text-green-600',
    comment_reply: 'bg-purple-100 text-purple-600',
    audit_passed: 'bg-emerald-100 text-emerald-600',
    audited_failed: 'bg-red-100 text-red-600'
  }
  return map[type] || 'bg-gray-100 text-gray-600'
}

const getIconColor = (type: string) => {
  const map: Record<string, string> = {
    system: 'text-blue-600',
    publish: 'text-green-600',
    comment_reply: 'text-purple-600',
    audit_passed: 'text-emerald-600',
    audited_failed: 'text-red-600'
  }
  return map[type] || 'text-gray-600'
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  if (diff < 60 * 1000) return '刚刚'
  if (diff < 60 * 60 * 1000) return `${Math.floor(diff / 60 / 1000)} 分钟前`
  if (diff < 24 * 60 * 60 * 1000) return `${Math.floor(diff / 60 / 60 / 1000)} 小时前`
  if (diff < 7 * 24 * 60 * 60 * 1000) return `${Math.floor(diff / 24 / 60 / 60 / 1000)} 天前`
  
  return date.toLocaleDateString('zh-CN')
}
</script>