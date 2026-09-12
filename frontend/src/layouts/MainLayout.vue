<template>
  <div class="min-h-screen bg-body">
    <!-- 侧边栏 -->
    <Sidebar />
    
    <!-- 顶部导航 -->
    <TopHeader :is-collapsed="sidebarCollapsed" />
    
    <!-- 主内容区域 -->
    <main 
      class="pt-16 transition-all duration-300"
      :class="sidebarCollapsed ? 'ml-[72px]' : 'ml-[260px]'"
    >
      <!-- 面包屑导航 -->
      <div class="bg-white border-b border-gray-200 px-6 py-3">
        <nav class="flex items-center gap-2 text-sm">
          <router-link to="/dashboard" class="text-gray-500 hover:text-primary transition-colors flex items-center gap-1">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            首页
          </router-link>
          <svg class="w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
          <span class="text-primary font-medium">{{ currentTitle }}</span>
        </nav>
      </div>
      
      <!-- 页面内容 -->
      <div class="p-6">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, RouterView, RouterLink } from 'vue-router'
import Sidebar from './Sidebar.vue'
import TopHeader from './TopHeader.vue'

const route = useRoute()
const sidebarCollapsed = ref(false)

// 当前页面标题
const currentTitle = computed(() => {
  return (route.meta.title as string) || '首页'
})
</script>

<style scoped>
.page-enter-active {
  transition: all 0.2s ease-out;
}

.page-leave-active {
  transition: all 0.15s ease-in;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(10px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}
</style>