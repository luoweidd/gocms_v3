<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">仪表盘</h1>
        <p class="text-sm text-gray-500">欢迎回来！以下是系统概览信息。</p>
      </div>
      <div class="flex items-center gap-3">
        <span class="text-sm text-gray-400">{{ lastUpdateTime }}</span>
        <button @click="refreshData" :disabled="loading" class="px-4 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors disabled:opacity-50 flex items-center gap-2">
          <svg v-if="!loading" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          刷新
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
      <!-- Articles Card -->
      <div class="ta-card p-5 flex items-center gap-4 transition-all duration-200 hover:shadow-cardHover hover:-translate-y-0.5">
        <div class="stat-icon bg-primary/10 text-primary">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <div class="flex-1">
          <div class="text-2xl font-bold text-gray-900">{{ formatNumber(stats.total_articles) }}</div>
          <div class="text-sm text-gray-500">文章总数</div>
          <div class="text-xs text-green-500 mt-1">+{{ stats.today_new_articles }} 今日</div>
        </div>
      </div>

      <!-- Videos Card -->
      <div class="ta-card p-5 flex items-center gap-4 transition-all duration-200 hover:shadow-cardHover hover:-translate-y-0.5">
        <div class="stat-icon bg-secondary/10 text-secondary">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.718v6.564a1 1 0 01-.553.982l-4.553 2.276a1 1 0 01-.907 0l-3.094-1.547a1 1 0 01-.553-.982V8.718a1 1 0 01.553-.982L9.093 5.969a1 1 0 01.907 0l4.553 2.276a1 1 0 01.553.982v.001z" />
          </svg>
        </div>
        <div class="flex-1">
          <div class="text-2xl font-bold text-gray-900">{{ formatNumber(stats.total_videos) }}</div>
          <div class="text-sm text-gray-500">视频总数</div>
          <div class="text-xs text-green-500 mt-1">+{{ stats.today_new_videos }} 今日</div>
        </div>
      </div>

      <!-- Comments Card -->
      <div class="ta-card p-5 flex items-center gap-4 transition-all duration-200 hover:shadow-cardHover hover:-translate-y-0.5">
        <div class="stat-icon bg-amber-500/10 text-amber-600">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z" />
          </svg>
        </div>
        <div class="flex-1">
          <div class="text-2xl font-bold text-gray-900">{{ formatNumber(stats.total_comments) }}</div>
          <div class="text-sm text-gray-500">评论总数</div>
          <div class="text-xs text-green-500 mt-1">+{{ stats.today_new_comments }} 今日</div>
        </div>
      </div>

      <!-- Users Card -->
      <div class="ta-card p-5 flex items-center gap-4 transition-all duration-200 hover:shadow-cardHover hover:-translate-y-0.5">
        <div class="stat-icon bg-purple-500/10 text-purple-600">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-12 0v1zm-3 4a1 1 0 11-2 0 1 1 0 012 0z" />
          </svg>
        </div>
        <div class="flex-1">
          <div class="text-2xl font-bold text-gray-900">{{ formatNumber(stats.total_users) }}</div>
          <div class="text-sm text-gray-500">用户总数</div>
          <div class="text-xs text-green-500 mt-1">+{{ stats.today_new_users }} 今日</div>
        </div>
      </div>
    </div>

    <!-- Charts Row -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6 mb-6">
      <!-- Content Growth Trend -->
      <div class="ta-card p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">内容增长趋势</h3>
          <select v-model="trendDays" @change="fetchTrendData" class="text-sm border border-gray-200 rounded-md px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-primary/20">
            <option value="7">最近7天</option>
            <option value="14">最近14天</option>
            <option value="30">最近30天</option>
          </select>
        </div>
        <div ref="trendChartRef" class="h-[320px]"></div>
      </div>

      <!-- Content Status Distribution -->
      <div class="ta-card p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">内容状态分布</h3>
        </div>
        <div class="grid grid-cols-2 gap-4 h-[320px]">
          <div ref="articleStatusChartRef" class="h-[140px] xl:h-[155px]"></div>
          <div ref="videoStatusChartRef" class="h-[140px] xl:h-[155px]"></div>
          <div ref="commentStatusChartRef" class="h-[140px] xl:h-[155px] col-span-2"></div>
        </div>
      </div>
    </div>

    <!-- Hardware Resources Row -->
    <div class="ta-card p-5 mb-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-base font-semibold text-gray-900">服务器资源监控</h3>
        <span class="text-xs px-2 py-1 bg-green-100 text-green-700 rounded-full flex items-center gap-1">
          <span class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></span>
          实时更新
        </span>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6">
        <!-- CPU -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">CPU 使用率</span>
            <span class="text-sm font-semibold" :class="cpuClass">{{ cpuPercent }}%</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000"
                 :class="cpuBarClass"
                 :style="{ width: `${cpuPercent}%` }"></div>
          </div>
          <div class="text-xs text-gray-400">{{ resourceDetail.cpuCores }} 核心</div>
        </div>

        <!-- Memory -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">内存使用率</span>
            <span class="text-sm font-semibold" :class="memoryClass">{{ memoryPercent }}%</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000"
                 :class="memoryBarClass"
                 :style="{ width: `${memoryPercent}%` }"></div>
          </div>
          <div class="text-xs text-gray-400">{{ formatBytes(resourceDetail.memoryUsed) }} / {{ formatBytes(resourceDetail.memoryTotal) }}</div>
        </div>

        <!-- Disk -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">磁盘使用率</span>
            <span class="text-sm font-semibold" :class="diskClass">{{ diskPercent }}%</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div class="h-full rounded-full transition-all duration-1000"
                 :class="diskBarClass"
                 :style="{ width: `${diskPercent}%` }"></div>
          </div>
          <div class="text-xs text-gray-400">{{ formatBytes(resourceDetail.diskUsed) }} / {{ formatBytes(resourceDetail.diskTotal) }}</div>
        </div>

        <!-- Goroutines -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Goroutines</span>
            <span class="text-sm font-semibold text-blue-600">{{ resourceDetail.goroutines }}</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div class="h-full bg-gradient-to-r from-blue-400 to-blue-600 rounded-full transition-all duration-1000"
                 :style="{ width: `${Math.min(resourceDetail.goroutines / 5, 100)}%` }"></div>
          </div>
          <div class="text-xs text-gray-400">运行时间: {{ formatUptime(resourceDetail.uptime) }}</div>
        </div>
      </div>
    </div>

    <!-- Category Distribution + Recent Activity -->
    <div class="grid grid-cols-1 xl:grid-cols-2 gap-6 mb-6">
      <!-- Article Category Distribution -->
      <div class="ta-card p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">文章分类分布</h3>
        </div>
        <div ref="articleCategoryChartRef" class="h-[280px]"></div>
      </div>

      <!-- Video Category Distribution -->
      <div class="ta-card p-5">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-base font-semibold text-gray-900">视频分类分布</h3>
        </div>
        <div ref="videoCategoryChartRef" class="h-[280px]"></div>
      </div>
    </div>

    <!-- Recent Activity -->
    <div class="ta-card p-5">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-base font-semibold text-gray-900">最近活动</h3>
        <button @click="fetchRecentActivity" class="text-sm text-primary hover:underline">刷新</button>
      </div>
      <div class="space-y-3">
        <div v-for="item in activities" :key="item.id" class="flex items-center gap-4 py-2 border-b last:border-b-0 border-gray-100">
          <div class="w-8 h-8 flex items-center justify-center rounded-full flex-shrink-0"
               :class="getTypeIconBg(item.type)">
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path v-if="item.type === 'article'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              <path v-else-if="item.type === 'video'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.718v6.564a1 1 0 01-.553.982l-4.553 2.276a1 1 0 01-.907 0l-3.094-1.547a1 1 0 01-.553-.982V8.718a1 1 0 01.553-.982L9.093 5.969a1 1 0 01.907 0l4.553 2.276a1 1 0 01.553.982v.001z" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
            </svg>
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-700 truncate">{{ item.title || item.content }}</p>
          </div>
          <span class="text-xs px-2 py-1 rounded-full" :class="getStatusClass(item.status)">
            {{ getStatusText(item.status) }}
          </span>
          <span class="text-xs text-gray-400 whitespace-nowrap">{{ formatTime(item.created_at || item.time) }}</span>
        </div>
        <div v-if="activities.length === 0" class="text-center py-8 text-gray-400">
          暂无活动数据
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="ta-card p-5 mt-6">
      <h3 class="text-base font-semibold text-gray-900 mb-4">快速操作</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-3">
        <router-link to="/videos/create" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-primary/10 text-primary rounded-lg flex items-center justify-center group-hover:bg-primary group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">上传视频</span>
        </router-link>
        <router-link to="/articles/create" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-secondary/10 text-secondary rounded-lg flex items-center justify-center group-hover:bg-secondary group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9l-2.829-2.829a2 2 0 012.829-2.829z" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">新建文章</span>
        </router-link>
        <router-link to="/users" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-purple-500/10 text-purple-600 rounded-lg flex items-center justify-center group-hover:bg-purple-600 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">用户管理</span>
        </router-link>
        <router-link to="/comments" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-green-500/10 text-green-600 rounded-lg flex items-center justify-center group-hover:bg-green-600 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">评论管理</span>
        </router-link>
        <router-link to="/media/files" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-orange-500/10 text-orange-600 rounded-lg flex items-center justify-center group-hover:bg-orange-600 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">文件管理</span>
        </router-link>
        <router-link to="/system/config" class="flex flex-col items-center gap-2 p-3 rounded-lg hover:bg-gray-50 transition-colors group">
          <div class="w-10 h-10 bg-gray-500/10 text-gray-600 rounded-lg flex items-center justify-center group-hover:bg-gray-600 group-hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
          </div>
          <span class="text-xs font-medium text-gray-700">系统设置</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import { getDashboardStats, getStatTrend, getRecentActivities } from '@/api/dashboard'

// ============ State ============
const loading = ref(false)
const trendDays = ref(7)
const lastUpdateTime = ref('')
const activities = ref<any[]>([])

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const stats: Record<string, any> = reactive({
  total_users: 0,
  today_new_users: 0,
  total_articles: 0,
  today_new_articles: 0,
  total_videos: 0,
  today_new_videos: 0,
  total_comments: 0,
  today_new_comments: 0,
  article_status: [] as any[],
  video_status: [] as any[],
  comment_status: [] as any[],
  resource_usage: null,
  article_categories: [] as any[],
  video_categories: [] as any[],
})

const resourceDetail = reactive({
  cpuCores: 0,
  cpuUsage: 0,
  memoryTotal: 0,
  memoryUsed: 0,
  diskTotal: 0,
  diskUsed: 0,
  goroutines: 0,
  uptime: 0,
})

const cpuPercent = ref(0)
const memoryPercent = ref(0)
const diskPercent = ref(0)

// ============ Chart refs ============
const trendChartRef = ref<HTMLElement>()
const articleStatusChartRef = ref<HTMLElement>()
const videoStatusChartRef = ref<HTMLElement>()
const commentStatusChartRef = ref<HTMLElement>()
const articleCategoryChartRef = ref<HTMLElement>()
const videoCategoryChartRef = ref<HTMLElement>()

// ============ Chart instances ============
let trendChart: ECharts | null = null
let articleStatusChart: ECharts | null = null
let videoStatusChart: ECharts | null = null
let commentStatusChart: ECharts | null = null
let articleCategoryChart: ECharts | null = null
let videoCategoryChart: ECharts | null = null

// ============ Auto refresh timer ============
let refreshTimer: ReturnType<typeof setInterval> | null = null

// ============ Helpers ============
const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatNumber = (n: number): string => {
  if (!n) return '0'
  if (n >= 10000) return (n / 10000).toFixed(1) + '万'
  return n.toLocaleString()
}

const formatTime = (dateStr: string): string => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  if (diff < 60 * 1000) return '刚刚'
  if (diff < 60 * 60 * 1000) return Math.floor(diff / (60 * 1000)) + ' 分钟前'
  if (diff < 24 * 60 * 60 * 1000) return Math.floor(diff / (60 * 60 * 1000)) + ' 小时前'
  return date.toLocaleDateString()
}

const formatUptime = (seconds: number): string => {
  if (!seconds) return '0 秒'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  if (days > 0) return `${days} 天 ${hours} 时`
  return `${hours} 小时`
}

const getStatusClass = (status: number): string => {
  switch (status) {
    case 0: return 'bg-yellow-100 text-yellow-700'
    case 1: return 'bg-green-100 text-green-700'
    case 2: return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}

const getStatusText = (status: number): string => {
  switch (status) {
    case 0: return '草稿'
    case 1: return '已发布'
    case 2: return '已下架'
    default: return '未知'
  }
}

const getTypeIconBg = (type: string): string => {
  switch (type) {
    case 'article': return 'bg-primary/10 text-primary'
    case 'video': return 'bg-secondary/10 text-secondary'
    default: return 'bg-amber-500/10 text-amber-600'
  }
}

const cpuClass = () => {
  if (cpuPercent.value < 50) return 'text-green-600'
  if (cpuPercent.value < 80) return 'text-yellow-600'
  return 'text-red-600'
}
const cpuBarClass = () => {
  if (cpuPercent.value < 50) return 'bg-gradient-to-r from-green-400 to-green-500'
  if (cpuPercent.value < 80) return 'bg-gradient-to-r from-yellow-400 to-yellow-500'
  return 'bg-gradient-to-r from-red-400 to-red-500'
}
const memoryClass = () => {
  if (memoryPercent.value < 60) return 'text-green-600'
  if (memoryPercent.value < 85) return 'text-yellow-600'
  return 'text-red-600'
}
const memoryBarClass = () => {
  if (memoryPercent.value < 60) return 'bg-gradient-to-r from-blue-400 to-blue-500'
  if (memoryPercent.value < 85) return 'bg-gradient-to-r from-yellow-400 to-yellow-500'
  return 'bg-gradient-to-r from-red-400 to-red-500'
}
const diskClass = () => {
  if (diskPercent.value < 60) return 'text-green-600'
  if (diskPercent.value < 85) return 'text-yellow-600'
  return 'text-red-600'
}
const diskBarClass = () => {
  if (diskPercent.value < 60) return 'bg-gradient-to-r from-purple-400 to-purple-500'
  if (diskPercent.value < 85) return 'bg-gradient-to-r from-yellow-400 to-yellow-500'
  return 'bg-gradient-to-r from-red-400 to-red-500'
}

// ============ Chart init ============
const initCharts = () => {
  if (trendChartRef.value) {
    trendChart = echarts.init(trendChartRef.value)
  }
  if (articleStatusChartRef.value) {
    articleStatusChart = echarts.init(articleStatusChartRef.value)
  }
  if (videoStatusChartRef.value) {
    videoStatusChart = echarts.init(videoStatusChartRef.value)
  }
  if (commentStatusChartRef.value) {
    commentStatusChart = echarts.init(commentStatusChartRef.value)
  }
  if (articleCategoryChartRef.value) {
    articleCategoryChart = echarts.init(articleCategoryChartRef.value)
  }
  if (videoCategoryChartRef.value) {
    videoCategoryChart = echarts.init(videoCategoryChartRef.value)
  }
}

const destroyCharts = () => {
  trendChart?.dispose()
  articleStatusChart?.dispose()
  videoStatusChart?.dispose()
  commentStatusChart?.dispose()
  articleCategoryChart?.dispose()
  videoCategoryChart?.dispose()
  trendChart = null
  articleStatusChart = null
  videoStatusChart = null
  commentStatusChart = null
  articleCategoryChart = null
  videoCategoryChart = null
}

const updateTimeDisplay = () => {
  lastUpdateTime.value = new Date().toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// ============ Chart data update ============
const updateTrendChart = (trendData: any) => {
  if (!trendChart) return
  
  const dates = trendData?.article?.map((t: any) => t.date) || []
  
  const option: any = {
    tooltip: {
      trigger: 'axis' as const,
      backgroundColor: 'rgba(255, 255, 255, 0.95)',
      borderColor: '#eee',
      textStyle: { color: '#333' },
    },
    legend: {
      data: ['文章', '视频', '评论'],
      bottom: 0,
      textStyle: { fontSize: 12 }
    },
    grid: { left: '3%', right: '4%', bottom: '10%', top: '10%', containLabel: true },
    xAxis: {
      type: 'category' as const,
      boundaryGap: false,
      data: dates,
      axisLine: { lineStyle: { color: '#e5e7eb' } },
      axisLabel: { color: '#6b7280', fontSize: 11 },
    },
    yAxis: {
      type: 'value' as const,
      splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } },
      axisLabel: { color: '#6b7280', fontSize: 11 },
    },
    series: [
      {
        name: '文章',
        type: 'line' as const,
        data: trendData?.article?.map((t: any) => t.count) || [],
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2.5 },
        itemStyle: { color: '#3b82f6' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(59, 130, 246, 0.3)' },
            { offset: 1, color: 'rgba(59, 130, 246, 0.05)' }
          ])
        }
      },
      {
        name: '视频',
        type: 'line' as const,
        data: trendData?.video?.map((t: any) => t.count) || [],
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2.5 },
        itemStyle: { color: '#f59e0b' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(245, 158, 11, 0.3)' },
            { offset: 1, color: 'rgba(245, 158, 11, 0.05)' }
          ])
        }
      },
      {
        name: '评论',
        type: 'line' as const,
        data: trendData?.comment?.map((t: any) => t.count) || [],
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2.5 },
        itemStyle: { color: '#10b981' },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
            { offset: 1, color: 'rgba(16, 185, 129, 0.05)' }
          ])
        }
      }
    ]
  }
  
  trendChart.setOption(option)
}

const updateStatusCharts = () => {
  // Article status pie chart
  if (articleStatusChart && stats.article_status?.length) {
    const articleOption: any = {
      tooltip: { trigger: 'item' as const },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        label: { show: true, position: 'outside', fontSize: 11 },
        data: stats.article_status.map((s: any) => ({
          name: s.status,
          value: s.count
        })),
        color: ['#3b82f6', '#10b981', '#f59e0b', '#ef4444']
      }]
    }
    articleStatusChart.setOption(articleOption)
  }

  // Video status pie chart
  if (videoStatusChart && stats.video_status?.length) {
    const videoOption: any = {
      tooltip: { trigger: 'item' as const },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['50%', '50%'],
        avoidLabelOverlap: false,
        label: { show: true, position: 'outside', fontSize: 11 },
        data: stats.video_status.map((s: any) => ({
          name: s.status,
          value: s.count
        })),
        color: ['#f59e0b', '#10b981', '#6366f1', '#ef4444']
      }]
    }
    videoStatusChart.setOption(videoOption)
  }

  // Comment status bar chart
  if (commentStatusChart && stats.comment_status?.length) {
    const commentOption: any = {
      tooltip: { trigger: 'axis' as const },
      grid: { left: '15%', right: '5%', top: '10%', bottom: '15%' },
      xAxis: {
        type: 'category' as const,
        data: stats.comment_status.map((s: any) => s.status),
        axisLabel: { fontSize: 10 }
      },
      yAxis: {
        type: 'value' as const,
        splitLine: { lineStyle: { color: '#f3f4f6', type: 'dashed' } }
      },
      series: [{
        type: 'bar',
        data: stats.comment_status.map((s: any) => s.count),
        itemStyle: { color: '#8b5cf6', borderRadius: [4, 4, 0, 0] }
      }]
    }
    commentStatusChart.setOption(commentOption)
  }
}

const updateCategoryCharts = () => {
  // Article category distribution
  if (articleCategoryChart && stats.article_categories?.length) {
    const catOption: any = {
      tooltip: { trigger: 'item' as const },
      series: [{
        type: 'pie',
        radius: '70%',
        center: ['50%', '50%'],
        label: {
          show: true,
          formatter: '{b}: {c} ({d}%)',
          fontSize: 11
        },
        data: stats.article_categories.map((c: any) => ({
          name: c.category_name,
          value: c.count
        })),
        color: ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#14b8a6']
      }]
    }
    articleCategoryChart.setOption(catOption)
  }

  // Video category distribution
  if (videoCategoryChart && stats.video_categories?.length) {
    const vcatOption: any = {
      tooltip: { trigger: 'item' as const },
      series: [{
        type: 'pie',
        radius: '70%',
        center: ['50%', '50%'],
        label: {
          show: true,
          formatter: '{b}: {c} ({d}%)',
          fontSize: 11
        },
        data: stats.video_categories.map((c: any) => ({
          name: c.category_name,
          value: c.count
        })),
        color: ['#f59e0b', '#10b981', '#6366f1', '#ef4444', '#ec4899']
      }]
    }
    videoCategoryChart.setOption(vcatOption)
  }
}

const updateResourceDisplay = () => {
  if (!stats.resource_usage) return
  
  const r = stats.resource_usage
  resourceDetail.cpuCores = r.cpu?.cores || 0
  resourceDetail.cpuUsage = r.cpu?.usage_percent || 0
  resourceDetail.memoryTotal = r.memory?.total || 0
  resourceDetail.memoryUsed = r.memory?.used || 0
  resourceDetail.diskTotal = r.disk?.total || 0
  resourceDetail.diskUsed = r.disk?.used || 0
  resourceDetail.goroutines = r.goroutines || 0
  resourceDetail.uptime = r.uptime || 0
  
  cpuPercent.value = Math.round(r.cpu?.usage_percent || 0)
  memoryPercent.value = Math.round((r.memory?.used / r.memory?.total) * 100 || 0)
  diskPercent.value = Math.round((r.disk?.used / r.disk?.total) * 100 || 0)
}

// ============ Data fetch ============
const fetchStats = async () => {
  try {
    const res = await getDashboardStats()
    if (res.data) {
      Object.assign(stats, res.data)
      updateResourceDisplay()
      
      // Update charts
      updateTrendChart(res.data)
      updateStatusCharts()
      updateCategoryCharts()
      updateTimeDisplay()
    }
  } catch (error) {
    console.error('获取Dashboard统计数据失败:', error)
  }
}

const fetchTrendData = async () => {
  try {
    const res = await getStatTrend(trendDays.value)
    if (res.data && trendChart) {
      updateTrendChart(res.data)
    }
  } catch (error) {
    console.error('获取趋势数据失败:', error)
  }
}

const fetchRecentActivity = async () => {
  try {
    const res = await getRecentActivities(10)
    if (res.data) {
      activities.value = res.data.map((item: any) => ({
        id: item.id,
        title: item.title || '',
        type: item.type,
        status: item.status,
        created_at: item.created_at,
        time: formatTime(item.created_at)
      }))
    }
  } catch (error) {
    console.error('获取最近活动失败:', error)
  }
}

const refreshData = async () => {
  loading.value = true
  try {
    await Promise.all([fetchStats(), fetchRecentActivity()])
  } finally {
    loading.value = false
  }
}

// ============ Lifecycle ============
onMounted(() => {
  initCharts()
  refreshData()
  
  // Auto refresh every 30 seconds
  refreshTimer = setInterval(() => {
    fetchStats()
    fetchRecentActivity()
  }, 30000)
})

onUnmounted(() => {
  destroyCharts()
  if (refreshTimer) clearInterval(refreshTimer)
})

// Responsive charts
watch([trendDays], () => {
  fetchTrendData()
})

// Handle window resize
const handleResize = () => {
  trendChart?.resize()
  articleStatusChart?.resize()
  videoStatusChart?.resize()
  commentStatusChart?.resize()
  articleCategoryChart?.resize()
  videoCategoryChart?.resize()
}

window.addEventListener('resize', handleResize)
onUnmounted(() => window.removeEventListener('resize', handleResize))
</script>

<style scoped>
.stat-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  flex-shrink: 0;
}

.page-container {
  max-width: 1600px;
  margin: 0 auto;
}
</style>