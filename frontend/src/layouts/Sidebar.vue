<template>
  <aside 
    class="fixed left-0 top-0 z-50 h-screen flex flex-col bg-sidebar text-white transition-all duration-300"
    :class="isCollapsed ? 'w-[72px]' : 'w-[260px]'"
  >
    <!-- Logo区域 -->
    <div class="flex items-center justify-center h-16 border-b border-sidebar-light/20 px-4">
      <div class="flex items-center gap-3" v-if="!isCollapsed">
        <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <span class="text-lg font-bold">GoCMS</span>
      </div>
      <div class="w-8 h-8 bg-primary rounded-lg flex items-center justify-center cursor-pointer" v-else @click="toggleSidebar">
        <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      </div>
    </div>

    <!-- 导航菜单 -->
    <nav class="flex-1 overflow-y-auto py-4">
      <ul class="space-y-1 px-3">
        <li v-for="item in menuItems" :key="item.name">
          <!-- 无子菜单项 -->
          <router-link 
            v-if="!item.children"
            :to="item.path"
            class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-200"
            :class="isActive(item.path) ? 'bg-primary/20 text-white' : 'text-gray-400 hover:bg-sidebar-hover hover:text-white'"
          >
            <span class="flex-shrink-0">
              <component :is="item.icon" class="w-5 h-5" />
            </span>
            <span v-if="!isCollapsed" class="whitespace-nowrap">{{ resolveLabel(item.label, item.i18nKey) }}</span>
          </router-link>

          <!-- 有子菜单项 -->
          <div v-else>
            <div 
              class="flex items-center justify-between px-3 py-2.5 rounded-lg text-sm cursor-pointer text-gray-400 hover:bg-sidebar-hover hover:text-white transition-all duration-200"
              @click="toggleMenu(item.name)"
            >
              <div class="flex items-center gap-3">
                <span class="flex-shrink-0">
                  <component :is="item.icon" class="w-5 h-5" />
                </span>
                <span v-if="!isCollapsed" class="whitespace-nowrap">{{ resolveLabel(item.label, item.i18nKey) }}</span>
                <!-- 折叠状态下也显示省略号表示有子菜单 -->
                <span v-else class="flex-1" />
              </div>
              <svg 
                class="w-4 h-4 transition-transform duration-200 text-gray-500" 
                :class="expandedMenus.includes(item.name) ? 'rotate-180' : ''"
                fill="none" viewBox="0 0 24 24" stroke="currentColor"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
            <!-- 子菜单 -->
            <ul 
              v-if="expandedMenus.includes(item.name) && !isCollapsed"
              class="ml-4 mt-1 space-y-1 overflow-hidden transition-all duration-300"
            >
              <li v-for="child in item.children" :key="child.name">
                <router-link 
                  :to="child.path"
                  class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-400 hover:bg-sidebar-hover hover:text-white transition-all duration-200"
                  :class="isActive(child.path) ? 'bg-primary/20 text-white' : ''"
                >
                  <span class="w-1.5 h-1.5 rounded-full bg-current flex-shrink-0" />
                  <span class="whitespace-nowrap">{{ resolveLabel(child.label, child.i18nKey) }}</span>
                </router-link>
              </li>
            </ul>
          </div>
        </li>
      </ul>
    </nav>

    <!-- 折叠按钮 -->
    <div class="p-3 border-t border-sidebar-light/20">
      <button 
        @click="toggleSidebar"
        class="w-full flex items-center justify-center px-3 py-2 rounded-lg text-gray-400 hover:bg-sidebar-hover hover:text-white transition-all duration-200"
      >
        <svg 
          class="w-5 h-5 transition-transform duration-300" 
          :class="isCollapsed ? 'rotate-180' : ''"
          fill="none" viewBox="0 0 24 24" stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
        </svg>
      </button>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { getMenuTree } from '@/api/menu'
// Heroicons 图标导入
import { 
  HomeIcon, 
  DocumentTextIcon, 
  FilmIcon, 
  ChatBubbleBottomCenterIcon,
  UsersIcon,
  Cog6ToothIcon,
  TableCellsIcon,
  Squares2X2Icon,
  ArrowRightOnRectangleIcon,
  KeyIcon,
  ShieldCheckIcon,
  TagIcon,
  ChatBubbleLeftRightIcon,
  BookOpenIcon,
  FolderIcon,
  CubeIcon,
  ClockIcon
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const route = useRoute()

// 侧边栏折叠状态
const isCollapsed = ref(false)
const expandedMenus = ref<string[]>([])
const menuItems = ref<any[]>([])

// 图标映射表：将后端存储的图标类名映射为 Heroicons 组件
const iconMap: Record<string, any> = {
  'HomeIcon': HomeIcon,
  'DocumentTextIcon': DocumentTextIcon,
  'FilmIcon': FilmIcon,
  'ChatBubbleBottomCenterIcon': ChatBubbleBottomCenterIcon,
  'UsersIcon': UsersIcon,
  'Cog6ToothIcon': Cog6ToothIcon,
  'TableCellsIcon': TableCellsIcon,
  'Squares2X2Icon': Squares2X2Icon,
  'ArrowRightOnRectangleIcon': ArrowRightOnRectangleIcon,
  'KeyIcon': KeyIcon,
  'ShieldCheckIcon': ShieldCheckIcon,
  'TagIcon': TagIcon,
  'ChatBubbleLeftRightIcon': ChatBubbleLeftRightIcon,
  'BookOpenIcon': BookOpenIcon,
  'FolderIcon': FolderIcon,
  'CubeIcon': CubeIcon,
  'ClockIcon': ClockIcon,
  'default': CubeIcon,
}

// 从后端获取菜单数据
const loadMenuData = async () => {
  try {
    const res = await getMenuTree()
    // 将后端返回的菜单树转换为前端格式
    menuItems.value = convertMenuTree(res.data || [])
    // 默认展开所有一级菜单
    expandedMenus.value = (res.data || []).map((item: any) => item.name)
  } catch (error) {
    console.error('加载菜单数据失败:', error)
    // 如果加载失败，使用空菜单
    menuItems.value = []
  }
}

// 后端路径到前端路由的映射表（修复后端数据与前端路由不一致的问题）
const backendPathToFrontendPath: Record<string, string> = {
  // 消息管理：后端可能返回 /messages/xxx，修正为 /message/xxx
  '/messages/notifications': '/message/notifications',
  '/messages/system': '/message/system',
  '/messages/categories': '/message/categories',
  '/message-categories': '/message/categories',
  
  // 媒体中心：后端可能返回 /media/gallery，修正为 /media/images
  '/media/gallery': '/media/images',
  
  // 定时任务：确保 cron 路径正确映射
  '/cron/tasks': '/cron/tasks',
  '/cron/logs': '/cron/logs',
  
   // 用户管理：/users/role 和 /users/permission 现在是有效路由（不再映射到 /system/roles）
}

// i18n key 映射表：根据路径自动映射到对应的 i18n key
const pathToI18nKey: Record<string, string> = {
  '/dashboard': 'dashboard.title',
  '/articles/list': 'article.title',
  '/articles/create': 'article.create',
  '/articles/categories': 'article.category',
  '/videos/list': 'video.title',
  '/videos/create': 'video.create',
  '/videos/categories': 'video.category',
  '/comments/list': 'comment.title',
  '/message/notifications': 'messages.myNotifications',
  '/message/system': 'messages.systemMessages',
  '/message/categories': 'messages.messageCategories',
  '/audit/workbench': 'audit.reviewWorkbench',
  '/audit/articles': 'audit.articleReview',
  '/audit/videos': 'audit.videoReview',
  '/users/list': 'user.title',
  '/users/profile': 'user.title',
  '/system/menus': 'menuManagement',
  '/system/roles': 'rolePermissions',
  '/system/logs': 'operationLog',
  '/system/backup': 'dataBackup',
  '/system/seo': 'seoManagement',
  '/system/config': 'systemConfig',
  '/ngac/permissions': 'permissionList',
  '/tenants/list': 'tenantList',
  '/cron/tasks': 'cron.tasks',
  '/cron/logs': 'cron.logs',
}

// 规范化并映射后端路径到前端路由路径
const normalizeAndMapPath = (path: string): string => {
  if (!path) return '/'
  
  // 先规范化
  let normalized = path.trim()
  if (!normalized.startsWith('/')) {
    normalized = '/' + normalized
  }
  if (normalized !== '/' && normalized.endsWith('/')) {
    normalized = normalized.slice(0, -1)
  }
  
  // 再映射后端路径到前端路径
  if (backendPathToFrontendPath[normalized]) {
    return backendPathToFrontendPath[normalized]
  }
  
  return normalized
}

// 根据路径获取 i18n key
const getI18nKeyByPath = (path: string): string => {
  if (!path) return ''
  // 精确匹配
  if (pathToI18nKey[path]) return pathToI18nKey[path]
  // 模糊匹配 (匹配父路径)
  for (const key of Object.keys(pathToI18nKey)) {
    if (path.startsWith(key + '/') || path === key) {
      return pathToI18nKey[key]
    }
  }
  return ''
}

// 规范化路径：确保路径以 / 开头，且没有多余的后缀斜杠
const normalizePath = (path: string): string => {
  if (!path) return '/'
  let normalized = path.trim()
  if (!normalized.startsWith('/')) {
    normalized = '/' + normalized
  }
  // 移除末尾多余的斜杠 (根路径 / 除外)
  if (normalized !== '/' && normalized.endsWith('/')) {
    normalized = normalized.slice(0, -1)
  }
  return normalized
}

// 将后端菜单树转换为前端格式
const convertMenuTree = (menus: any[]): any[] => {
  return menus.map(menu => {
    const iconComponent = iconMap[menu.icon] || iconMap['default']
    const mappedPath = normalizeAndMapPath(menu.path)
    const item: any = {
      name: menu.name || `menu-${menu.id}`,
      label: menu.title || menu.name,
      i18nKey: getI18nKeyByPath(mappedPath),
      icon: iconComponent,
      path: mappedPath,
    }
    if (menu.children && menu.children.length > 0) {
      item.children = convertChildMenu(menu.children)
    }
    return item
  })
}

// 递归转换子菜单
const convertChildMenu = (children: any[]): any[] => {
  return children.map(child => {
    const iconComponent = iconMap[child.icon] || iconMap['default']
    const mappedPath = normalizeAndMapPath(child.path)
    const item: any = {
      name: child.name || `menu-${child.id}`,
      label: child.title || child.name,
      i18nKey: getI18nKeyByPath(mappedPath),
      icon: iconComponent,
      path: mappedPath,
    }
    if (child.children && child.children.length > 0) {
      item.children = convertChildMenu(child.children)
    }
    return item
  })
}

// 加载菜单数据
onMounted(() => {
  loadMenuData()
})

// 检查路径是否匹配
const isActive = (path: string) => {
  return route.path === path || route.path.startsWith(path + '/')
}

// 切换菜单展开状态
const toggleMenu = (name: string) => {
  const index = expandedMenus.value.indexOf(name)
  if (index > -1) {
    expandedMenus.value.splice(index, 1)
  } else {
    expandedMenus.value.push(name)
  }
}

// 切换侧边栏折叠状态
const toggleSidebar = () => {
  isCollapsed.value = !isCollapsed.value
}

// 解析显示标签：优先使用 i18n key，如果没有则直接返回原始 label
const resolveLabel = (label: string, i18nKey: string): string => {
  if (!i18nKey || !label) return label
  // 尝试用 i18n key 翻译
  try {
    const translated = t(i18nKey)
    // 如果翻译结果与原始 label 相同，说明 key 可能不存在，返回原始 label
    if (translated === i18nKey) {
      return label
    }
    return translated
  } catch {
    // 翻译失败，返回原始 label
    return label
  }
}
</script>

<style scoped>
/* 确保活跃链接有蓝色左边框指示器 */
.router-link-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 60%;
  background: #4F46E5;
  border-radius: 0 2px 2px 0;
}

.router-link-active {
  position: relative;
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 4px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>