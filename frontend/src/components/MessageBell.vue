<template>
  <div class="message-bell-wrapper relative">
    <!-- 消息铃铛图标和徽标 -->
    <button
      @click="togglePopover"
      class="message-bell-btn relative p-2 text-gray-600 hover:text-primary transition-colors"
      :aria-label="`通知 ${unreadCount > 0 ? unreadCount + '条未读' : ''}`"
    >
      <!-- 徽标数 -->
      <span
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 bg-red-500 text-white text-xs w-5 h-5 flex items-center justify-center rounded-full font-medium"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
      <!-- 铃铛图标 (Heroicons Bell) -->
      <svg
        v-if="unreadCount === 0"
        class="w-6 h-6"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
        />
      </svg>
      <!-- 铃铛填充图标 (Heroicons BellFilled) -->
      <svg
        v-else
        class="w-6 h-6"
        fill="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          d="M18 3A1.993 1.993 0 0016.03 4.97l-.653.653a.25.25 0 01-.354 0l-.707-.707a.25.25 0 00-.354 0l-.707.707a.25.25 0 01-.354 0l-.653-.653A1.993 1.993 0 007.97 3L8.28 5h7.44l.32-2zM19.865 8.098a1.001 1.001 0 00-.865-.496H16.19l-.566 9.436a2.002 2.002 0 01-1.995 1.87h-3.048a2.001 2.001 0 01-1.995-1.87l-.566-9.436H5.001a1.001 1.001 0 00-.865 1.498l2.692 4.487a1 1 0 00.865.496h8.31a1 1 0 00.865-.496l2.692-4.487a1.001 1.001 0 00-.1-1.011z"
        />
      </svg>
    </button>

    <!-- 通知下拉面板 -->
    <div
      v-if="popoverVisible"
      class="absolute right-0 top-full mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 z-50"
    >
      <!-- 面板头部 -->
      <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200">
        <h3 class="font-semibold text-gray-900">通知</h3>
        <button
          v-if="unreadCount > 0"
          @click="handleMarkAllRead"
          class="text-sm text-primary hover:text-primary-dark"
        >
          全部已读
        </button>
      </div>

      <!-- 通知列表 -->
      <div class="max-h-80 overflow-y-auto relative">
        <!-- 加载遮罩 -->
        <div
          v-if="loading"
          class="absolute inset-0 bg-white/80 flex items-center justify-center z-10"
        >
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary"></div>
        </div>
        <div
          v-for="item in previewList"
          :key="item.id"
          @click="handleItemClick(item)"
          class="px-4 py-3 hover:bg-gray-50 cursor-pointer transition-colors border-b border-gray-100 last:border-b-0"
          :class="{ 'bg-blue-50': !item.is_read }"
        >
          <div class="flex items-start gap-3">
            <!-- 消息类型图标 -->
            <div class="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-primary">
              <svg
                v-if="item.message?.message_type === 'system'"
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <svg
                v-else-if="item.message?.message_type === 'publish'"
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3v4M3 5h4M6 17v4m-2-2h4m5-16l2.286 6.857L21 12l-5.714 2.143L13 21l-2.286-6.857L5 12l5.714-2.143L13 3z"/>
              </svg>
              <svg
                v-else-if="item.message?.message_type === 'comment_reply'"
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
              </svg>
              <svg
                v-else-if="item.message?.message_type === 'audit_passed'"
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <svg
                v-else-if="item.message?.message_type === 'audited_failed'"
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              <svg
                v-else
                class="w-4 h-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <!-- 消息内容 -->
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-900 truncate">{{ item.message?.title || '无标题' }}</p>
              <p class="text-xs text-gray-500 mt-1">{{ formatTime(item.created_at) }}</p>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="previewList.length === 0 && !loading" class="px-4 py-8 text-center">
          <svg class="w-12 h-12 mx-auto text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
          </svg>
          <p class="text-sm text-gray-500">暂无通知</p>
        </div>
      </div>

      <!-- 面板底部 -->
      <div class="px-4 py-3 border-t border-gray-200">
        <router-link
          to="/message/notifications"
          class="block text-center text-sm text-primary hover:text-primary-dark font-medium"
        >
          查看全部通知
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getMyNotifications, getMyNotificationStats, markAllAsRead } from '@/api/notification'

const router = useRouter()

// 响应式数据
const unreadCount = ref(0)
const previewList = ref<any[]>([])
const loading = ref(false)
const popoverVisible = ref(false)

// 格式化时间
const formatTime = (time: string | undefined) => {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  
  // 一分钟内显示具体时间
  if (diff < 60 * 1000) {
    return '刚刚'
  }
  // 一小时内显示分钟
  if (diff < 60 * 60 * 1000) {
    return `${Math.floor(diff / (60 * 1000))}分钟前`
  }
  // 其他显示时间
  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 切换弹出层
const togglePopover = () => {
  popoverVisible.value = !popoverVisible.value
  if (popoverVisible.value) {
    fetchStats()
    fetchPreviewList()
  }
}

// 获取通知统计
const fetchStats = async () => {
  try {
    const res = await getMyNotificationStats()
    unreadCount.value = res.data?.unread_count || 0
  } catch (error) {
    console.error('获取通知统计失败:', error)
  }
}

// 获取预览列表
const fetchPreviewList = async () => {
  loading.value = true
  try {
    const res = await getMyNotifications({
      page: 1,
      page_size: 5,
      is_read: 0 // 只显示未读
    })
    previewList.value = res.data?.list || []
  } catch (error) {
    console.error('获取预览列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 全部已读
const handleMarkAllRead = async () => {
  try {
    await markAllAsRead()
    unreadCount.value = 0
    previewList.value = []
    // 使用原生 alert 替代 ElMessage
    alert('已全部标记为已读')
  } catch (error) {
    console.error('全部标记失败:', error)
    alert('操作失败')
  }
}

// 点击通知项
const handleItemClick = (item: any) => {
  // 标记为已读
  if (!item.is_read) {
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }
  router.push({ path: '/message/notifications', query: { id: item.message_id } })
}

// 点击外部关闭弹出层
const handleClickOutside = (event: MouseEvent) => {
  const messageBell = document.querySelector('.message-bell-wrapper')
  if (messageBell && !messageBell.contains(event.target as Node)) {
    popoverVisible.value = false
  }
}

// 生命周期
onMounted(() => {
  fetchStats()
  document.addEventListener('click', handleClickOutside)
})

// 组件卸载时移除事件监听
import { onUnmounted } from 'vue'
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background-color: #e5e7eb;
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background-color: #d1d5db;
}
</style>