<template>
  <header class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 shadow-header fixed top-0 z-40 transition-all duration-300"
    :class="isCollapsed ? 'left-[72px] w-[calc(100%-72px)]' : 'left-[260px] w-[calc(100%-260px)]'">
    <!-- 左侧：搜索框 -->
    <div class="flex items-center gap-4 flex-1">
      <div class="relative max-w-sm w-64">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          placeholder="搜索..."
          class="w-full pl-10 pr-4 py-2 text-sm border border-gray-200 rounded-lg bg-gray-50 focus:bg-white transition-colors outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary hover:bg-white"
        />
      </div>
    </div>

    <!-- 右侧：工具栏 -->
    <div class="flex items-center gap-2">
        <!-- 消息铃铛组件 -->
        <MessageBell />

      <!-- 分割线 -->
      <div class="w-px h-6 bg-gray-200 mx-1"></div>

      <!-- 用户菜单 -->
      <div data-user-menu class="relative">
        <button 
          @click="toggleUserMenu"
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <div class="w-8 h-8 bg-primary/10 text-primary rounded-full flex items-center justify-center font-medium text-sm">
            {{ avatarChar }}
          </div>
          <div class="hidden lg:block text-left">
            <p class="text-sm font-medium text-gray-700">{{ displayName }}</p>
          </div>
          <svg class="w-4 h-4 text-gray-500 transition-transform" :class="{ 'rotate-180': isOpenUserMenu }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>

        <!-- 下拉菜单 -->
        <div 
          v-if="isOpenUserMenu"
          class="absolute right-0 mt-2 w-56 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50"
        >
          <div class="px-4 py-3 border-b border-gray-100">
            <p class="text-sm font-medium text-gray-900">{{ displayName }}</p>
            <p class="text-xs text-gray-500 truncate">{{ authStore.userInfo?.email || '' }}</p>
          </div>
          <div class="py-1">
            <button 
              @click="goToProfile"
              class="w-full flex items-center gap-3 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              个人信息
            </button>
            <hr class="my-1 border-gray-100">
            <button 
              @click="handleLogout"
              class="w-full flex items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
              </svg>
              退出登录
            </button>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import MessageBell from '@/components/MessageBell.vue'

interface Props {
  isCollapsed: boolean
}

defineProps<Props>()

const router = useRouter()
const authStore = useAuthStore()

const isOpenUserMenu = ref(false)

// 显示名称
const displayName = computed(() => {
  return authStore.userInfo?.username || '用户'
})

// 头像字符
const avatarChar = computed(() => {
  return displayName.value.charAt(0).toUpperCase()
})

const toggleUserMenu = () => {
  isOpenUserMenu.value = !isOpenUserMenu.value
}

// 跳转到个人信息页面
const goToProfile = () => {
  router.push('/users/profile')
  isOpenUserMenu.value = false
}

// 跳转到消息通知页面
const goToMessages = () => {
  router.push('/message/notifications')
  isOpenUserMenu.value = false
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
  isOpenUserMenu.value = false
}

// 点击外部关闭用户菜单
document.addEventListener('click', (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('[data-user-menu]')) {
    isOpenUserMenu.value = false
  }
})
</script>