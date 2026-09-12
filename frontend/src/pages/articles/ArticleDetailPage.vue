<template>
  <div class="page-container max-w-4xl">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <button class="text-gray-500 hover:text-gray-700 transition-colors" @click="$router.back()">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </button>
          <div>
             <h1 class="text-2xl font-bold text-gray-900 mb-1">{{ article?.title || '文章详情' }}</h1>
            <p class="text-sm text-gray-500">查看文章内容</p>
          </div>
        </div>
        <div class="flex gap-2">
          <button 
            v-if="article" 
            :class="isPreviewMode ? 'bg-indigo-600 text-white' : 'bg-white text-gray-700 border border-gray-300'" 
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all hover:shadow-md flex items-center gap-2"
            @click="isPreviewMode = !isPreviewMode"
          >
            <svg v-if="!isPreviewMode" class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            {{ isPreviewMode ? '编辑模式' : '预览模式' }}
          </button>
          <router-link 
            v-if="article" 
            :to="`/articles/edit/${article.id}`" 
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
            编辑文章
          </router-link>
        </div>
      </div>
    </div>

    <!-- Article Content -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="text-center">
        <svg class="w-12 h-12 mx-auto mb-4 animate-spin text-indigo-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-500">加载中...</p>
      </div>
    </div>

    <div v-else-if="article" class="space-y-6">
      <!-- Article Card -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <!-- Article Header -->
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-center gap-3 mb-4 flex-wrap">
            <span :class="{
              'bg-emerald-100 text-emerald-800': article.status === 1,
              'bg-amber-100 text-amber-800': article.status === 0,
              'bg-red-100 text-red-800': article.status === 2
            }" class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium">
              {{ article.status === 1 ? '已发布' : article.status === 0 ? '草稿' : '已下架' }}
            </span>
            <span v-if="article.category_name" class="inline-flex items-center px-3 py-1 rounded-full text-sm bg-indigo-100 text-indigo-800">
              {{ article.category_name }}
            </span>
            <span v-if="article.view_count" class="inline-flex items-center text-sm text-gray-500">
              <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              {{ article.view_count }} 次浏览
            </span>
          </div>
          <h1 class="text-3xl font-bold text-gray-900">{{ article.title }}</h1>
        </div>

        <!-- Article Body -->
        <div class="p-6">
          <!-- Cover Image -->
          <div v-if="article.cover" class="mb-6 rounded-xl overflow-hidden">
            <img :src="getImageUrl(article.cover)" :alt="article.title" class="w-full h-64 object-cover" />
          </div>

          <!-- Tags -->
          <div v-if="article.tags && article.tags.length > 0" class="mb-6 flex items-center gap-2 flex-wrap">
            <span class="text-sm text-gray-500">标签：</span>
            <span v-for="tag in article.tags" :key="typeof tag === 'object' ? (tag as any).id : tag" class="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-gray-200 transition-colors">
              # {{ typeof tag === 'object' ? (tag as any).name : tag }}
            </span>
          </div>

          <!-- Content -->
          <div v-if="isPreviewMode" class="article-content prose max-w-none">
            <div v-html="renderedContent" class="text-gray-800 leading-relaxed"></div>
          </div>
          <div v-else class="text-gray-800 leading-relaxed whitespace-pre-wrap">{{ article.content }}</div>
        </div>

        <!-- Article Footer -->
        <div class="px-6 py-4 bg-gray-50 border-t border-gray-200">
          <div class="flex items-center justify-between text-sm text-gray-500">
            <div class="flex items-center gap-4">
              <span v-if="article.author_name">作者：{{ article.author_name }}</span>
              <span v-if="article.published_at">发布时间：{{ formatDate(article.published_at) }}</span>
            </div>
            <div class="flex items-center gap-4">
              <span v-if="article.created_at">创建时间：{{ formatDate(article.created_at) }}</span>
              <span v-if="article.updated_at">更新时间：{{ formatDate(article.updated_at) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Markdown Source Preview -->
      <div v-if="isPreviewMode && article.content" class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
          <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <svg class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Markdown 源码
          </h3>
        </div>
        <div class="p-6">
          <pre class="bg-gray-900 text-green-400 rounded-lg p-4 overflow-x-auto text-sm"><code>{{ article.content }}</code></pre>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="flex items-center justify-center py-20">
      <div class="text-center">
        <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">文章不存在</h3>
        <p class="text-gray-500 mb-4">该文章可能已被删除或移除</p>
        <button @click="$router.back()" class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors">
          返回上一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getArticle } from '@/api/article'
import { message } from '@/utils/message'

const route = useRoute()
const articleId = computed(() => Number(route.params.id))

interface ArticleItem {
  id: number
  title: string
  summary?: string
  description?: string
  content: string
  cover?: string
  category_id: number
  category_name?: string
  author_id: number
  author_name?: string
  status: 0 | 1 | 2
  view_count: number
  tags?: Array<{ id: number; name: string }>
  tag_ids?: number[]
  published_at?: string
  created_at: string
  updated_at: string
}

const article = ref<ArticleItem | null>(null)
const loading = ref(false)
const isPreviewMode = ref(true)

const renderedContent = computed(() => {
  if (!article.value?.content) return ''
  // Simple markdown to HTML converter
  let html = article.value.content
  
  // Headers
  html = html.replace(/^### (.*)$/gm, '<h3 class="text-xl font-bold mt-6 mb-3">$1</h3>')
  html = html.replace(/^## (.*)$/gm, '<h2 class="text-2xl font-bold mt-6 mb-3">$1</h2>')
  html = html.replace(/^# (.*)$/gm, '<h1 class="text-3xl font-bold mt-6 mb-3">$1</h1>')
  
  // Bold and italic
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')
  html = html.replace(/~~(.*?)~~/g, '<del>$1</del>')
  
  // Links
  html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" class="text-indigo-600 hover:text-indigo-800 underline" target="_blank">$1</a>')
  
  // Images
  html = html.replace(/!\[(.*?)\]\((.*?)\)/g, '<img src="$2" alt="$1" class="max-w-full h-auto rounded-lg my-4" />')
  
  // Code blocks
  html = html.replace(/```([\s\S]*?)```/g, '<pre class="bg-gray-900 text-green-400 rounded-lg p-4 my-4 overflow-x-auto"><code>$1</code></pre>')
  html = html.replace(/`([^`]+)`/g, '<code class="bg-gray-100 text-red-600 px-1.5 py-0.5 rounded text-sm">$1</code>')
  
  // Blockquotes
  html = html.replace(/^> (.*)$/gm, '<blockquote class="border-l-4 border-indigo-500 pl-4 italic text-gray-600 my-4">$1</blockquote>')
  
  // Unordered lists
  html = html.replace(/^- (.*)$/gm, '<li class="ml-4 list-disc">$1</li>')
  
  // Ordered lists
  html = html.replace(/^\d+\. (.*)$/gm, '<li class="ml-4 list-decimal">$1</li>')
  
  // Horizontal rules
  html = html.replace(/^---$/gm, '<hr class="my-6 border-gray-200">')
  
  // Paragraphs
  html = html.split('\n\n').map(p => {
    if (p.startsWith('<h') || p.startsWith('<blockquote') || p.startsWith('<pre') || p.startsWith('<hr') || p.startsWith('<li')) {
      return p
    }
    if (p.trim()) {
      return `<p class="my-4">${p}</p>`
    }
    return p
  }).join('\n')
  
  return html
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

const getImageUrl = (filePath: string) => {
  if (!filePath) return ''
  // If it's already a full URL or starts with http, return as is
  if (filePath.startsWith('http://') || filePath.startsWith('https://')) {
    return filePath
  }
  // If it starts with /api, return as is
  if (filePath.startsWith('/api/')) {
    return `${import.meta.env.VITE_API_BASE_URL || ''}${filePath}`
  }
  // Otherwise prepend /api
  return `/api${filePath}`
}

const loadArticle = async () => {
  loading.value = true
  try {
    const res = await getArticle(articleId.value)
    article.value = res.data as ArticleItem | null
  } catch (error) {
    message.error('加载文章失败')
    console.error('加载文章失败:', error)
    article.value = null
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (articleId.value) {
    loadArticle()
  } else {
    message.error('无效的文章ID')
  }
})
</script>

<style scoped>
.article-content p {
  margin-bottom: 1rem;
  line-height: 1.75;
}

.article-content img {
  max-width: 100%;
  height: auto;
}

.article-content blockquote {
  border-left: 4px solid #e5e7eb;
  padding-left: 1rem;
  font-style: italic;
  color: #6b7280;
}

.article-content pre {
  background-color: #1f2937;
  color: #10b981;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
}

.article-content code {
  background-color: #f3f4f6;
  color: #dc2626;
  padding: 0.125rem 0.25rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
}

.article-content li {
  margin-left: 1.5rem;
  margin-bottom: 0.25rem;
}
</style>