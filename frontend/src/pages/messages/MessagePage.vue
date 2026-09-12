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
          <h1 class="text-2xl font-bold text-gray-900 mb-1">消息通知</h1>
          <p class="text-sm text-gray-500">查看您的系统消息和通知</p>
        </div>
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
          <span v-if="(tab.key as string) === 'all' || (tab.key as string) === 'system' || (tab.key as string) === 'mention'" 
                class="ml-1.5 px-1.5 py-0.5 text-xs rounded-full"
                :class="unreadMap[tab.key] > 0 ? 'bg-red-100 text-red-600' : 'bg-gray-100 text-gray-500'">
            {{ unreadMap[tab.key] > 99 ? '99+' : unreadMap[tab.key] }}
          </span>
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
      
      <div v-else-if="messages.length === 0" class="py-12 text-center text-gray-400">
        <svg class="w-16 h-16 mx-auto mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        <p class="text-sm">暂无消息</p>
      </div>

      <div v-else class="divide-y divide-gray-100">
        <div 
          v-for="message in messages" 
          :key="message.id"
          @click="markAsRead(message.id)"
          :class="[
            'px-6 py-4 hover:bg-gray-50 cursor-pointer transition-colors',
            !message.is_read ? 'bg-blue-50/30' : ''
          ]"
        >
          <div class="flex items-start gap-4">
            <!-- Icon -->
            <div class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
                 :class="getIconBg(message.type)">
              <svg class="w-5 h-5" :class="getIconColor(message.type)" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path v-if="message.type === 'system'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else-if="message.type === 'success'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else-if="message.type === 'warning'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                <path v-else-if="message.type === 'error'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                      d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
            </div>
            
            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-sm font-medium text-gray-900">{{ message.title }}</span>
                <span v-if="!message.is_read" class="w-2 h-2 bg-primary rounded-full"></span>
              </div>
              <p class="text-sm text-gray-500 mt-1 truncate">{{ message.content }}</p>
              <div class="flex items-center gap-4 mt-2">
                <span class="text-xs text-gray-400">{{ formatTime(message.created_at) }}</span>
                <span v-if="message.source" class="text-xs text-gray-400">{{ message.source }}</span>
              </div>
            </div>

            <!-- Actions -->
            <button 
              v-if="message.is_read"
              @click.stop="deleteMessage(message.id)"
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
      <div v-if="messages.length > 0" class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条消息</span>
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
          class="text-sm text-primary hover:text-primary-700 font-medium"
        >
          全部标为已读
        </button>
      </div>
    </TailCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'

interface MessageItem {
  id: number
  title: string
  content: string
  type: 'system' | 'success' | 'warning' | 'error' | 'notification'
  is_read: boolean
  source?: string
  created_at: string
}

const tabs = [
  { key: 'all', label: '全部消息' },
  { key: 'system', label: '系统消息' },
  { key: 'notification', label: '通知' }
] as const

type TabKey = typeof tabs[number]['key']

const router = useRouter()
const activeTab = ref<TabKey>('all')
const loading = ref(false)
const messages = ref<MessageItem[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))

// 模拟消息数据
const mockMessages: MessageItem[] = [
  {
    id: 1,
    title: '系统更新通知',
    content: 'GoCMS v3.2 已发布，包含性能优化和安全修复，请及时更新。',
    type: 'system',
    is_read: false,
    source: '系统',
    created_at: new Date(Date.now() - 1000 * 60 * 5).toISOString() // 5 minutes ago
  },
  {
    id: 2,
    title: '文章审核通过',
    content: '您的文章《如何使用 Go 构建 RESTful API》已通过审核并发布。',
    type: 'success',
    is_read: false,
    source: '内容审核',
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString() // 2 hours ago
  },
  {
    id: 3,
    title: '存储空间即将不足',
    content: '您的存储空间使用率已超过 80%，请及时清理或扩容。',
    type: 'warning',
    is_read: true,
    source: '系统监控',
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString() // 1 day ago
  },
  {
    id: 4,
    title: '登录异常提醒',
    content: '检测到您的账户在新设备上登录，如非本人操作请及时修改密码。',
    type: 'error',
    is_read: true,
    source: '安全中心',
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 24 * 3).toISOString() // 3 days ago
  },
  {
    id: 5,
    title: '新评论通知',
    content: '用户对您的文章《Vue 3 实战教程》发表了一条评论。',
    type: 'notification',
    is_read: true,
    source: '评论系统',
    created_at: new Date(Date.now() - 1000 * 60 * 60 * 24 * 5).toISOString() // 5 days ago
  }
]

// 初始化消息
onMounted(() => {
  messages.value = mockMessages
  pagination.total = mockMessages.length
})

const unreadMap = computed(() => ({
  all: messages.value.filter(m => !m.is_read).length,
  system: messages.value.filter(m => !m.is_read && m.type === 'system').length,
  notification: messages.value.filter(m => !m.is_read && m.type === 'notification').length
}))

const hasUnread = computed(() => unreadMap.value.all > 0)

const getIconBg = (type: string) => {
  const map: Record<string, string> = {
    system: 'bg-blue-100 text-blue-600',
    success: 'bg-green-100 text-green-600',
    warning: 'bg-yellow-100 text-yellow-600',
    error: 'bg-red-100 text-red-600',
    notification: 'bg-purple-100 text-purple-600'
  }
  return map[type] || 'bg-gray-100 text-gray-600'
}

const getIconColor = (type: string) => {
  const map: Record<string, string> = {
    system: 'text-blue-600',
    success: 'text-green-600',
    warning: 'text-yellow-600',
    error: 'text-red-600',
    notification: 'text-purple-600'
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

const markAsRead = async (id: number) => {
  const msg = messages.value.find(m => m.id === id)
  if (msg) {
    msg.is_read = true
    message.success('已标记为已读')
  }
}

const markAllAsRead = async () => {
  messages.value.forEach(m => m.is_read = true)
  message.success('全部已标为已读')
}

const deleteMessage = async (id: number) => {
  messages.value = messages.value.filter(m => m.id !== id)
  pagination.total--
  message.success('删除成功')
}

const changePage = (page: number) => {
  if (page < 1 || page > totalPages.value) return
  pagination.page = page
  // TODO: 实际应该重新请求数据
}
</script>

<style scoped>
.group:hover .group-hover\:opacity-100 {
  opacity: 1;
}
</style>