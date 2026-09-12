<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">SEO管理</h1>
        <p class="text-sm text-gray-500">管理系统搜索引擎优化设置</p>
      </div>
      <div class="flex gap-2">
        <TailButton type="default" variant="outline" @click="handleGenerateSitemap">
          <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          生成站点地图
        </TailButton>
        <TailButton type="primary" @click="showModal = true">
          <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新增SEO设置
        </TailButton>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总设置数</p>
        <p class="text-2xl font-bold text-gray-900">{{ totalCount }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-blue-500">
        <p class="text-sm text-gray-500">文章SEO</p>
        <p class="text-2xl font-bold text-blue-600">{{ articleCount }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-green-500">
        <p class="text-sm text-gray-500">视频SEO</p>
        <p class="text-2xl font-bold text-green-600">{{ videoCount }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-purple-500">
        <p class="text-sm text-gray-500">页面SEO</p>
        <p class="text-2xl font-bold text-purple-600">{{ pageCount }}</p>
      </div>
    </div>

    <!-- Search Filter -->
    <div class="bg-white rounded-lg shadow-sm p-4 mb-6">
      <div class="flex flex-wrap items-center gap-4">
        <div class="flex-1 min-w-[200px]">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="搜索关键词..."
            class="ta-input w-full"
            @input="handleSearch"
          />
        </div>
        <select v-model="searchResourceType" class="ta-input w-40" @change="handleSearch">
          <option value="">全部类型</option>
          <option value="article">文章</option>
          <option value="video">视频</option>
          <option value="page">页面</option>
        </select>
        <TailButton variant="outline" @click="handleReset">重置</TailButton>
      </div>
    </div>

    <!-- SEO Settings Table -->
    <div class="bg-white rounded-lg shadow-sm overflow-hidden">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源类型</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">SEO标题</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">描述</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">关键词</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="item in seoList" :key="item.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getResourceTypeClass(item.resource_type)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                {{ RESOURCE_TYPE_MAP[item.resource_type] || item.resource_type }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ item.resource_id }}</td>
            <td class="px-6 py-4 text-sm text-gray-900 max-w-[200px] truncate" :title="item.seo_title">{{ item.seo_title || '-' }}</td>
            <td class="px-6 py-4 text-sm text-gray-500 max-w-[200px] truncate" :title="item.seo_description">{{ item.seo_description || '-' }}</td>
            <td class="px-6 py-4 text-sm text-gray-500 max-w-[150px] truncate" :title="item.seo_keywords">{{ item.seo_keywords || '-' }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex gap-2">
                <TailButton type="default" variant="ghost" size="sm" @click="editSeo(item)">编辑</TailButton>
                <TailButton type="danger" variant="ghost" size="sm" @click="deleteSeo(item.id)">删除</TailButton>
              </div>
            </td>
          </tr>
          <tr v-if="seoList.length === 0">
            <td colspan="6" class="px-6 py-12 text-center text-gray-500">
              暂无SEO设置数据
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="totalCount > 0" class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ totalCount }} 条</span>
        <div class="flex gap-2">
          <TailButton
            variant="outline"
            size="sm"
            :disabled="currentPage <= 1"
            @click="currentPage--; loadSeoList()"
          >
            上一页
          </TailButton>
          <TailButton
            variant="outline"
            size="sm"
            :disabled="currentPage >= totalPages"
            @click="currentPage++; loadSeoList()"
          >
            下一页
          </TailButton>
        </div>
      </div>
    </div>

    <!-- SEO Settings Modal -->
    <TailModal :open="showModal" @update:open="showModal = $event" :title="isEdit ? '编辑SEO设置' : '新增SEO设置'" width="max-w-2xl">
      <div class="space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="ta-label">资源类型 <span class="text-red-500">*</span></label>
            <select v-model="form.resource_type" class="ta-input w-full">
              <option value="">请选择</option>
              <option value="article">文章</option>
              <option value="video">视频</option>
              <option value="page">页面</option>
            </select>
          </div>
          <div>
            <label class="ta-label">资源ID <span class="text-red-500">*</span></label>
            <input v-model.number="form.resource_id" type="number" class="ta-input w-full" placeholder="如: 文章ID" />
          </div>
        </div>

        <div>
          <label class="ta-label">SEO标题 <span class="text-red-500">*</span></label>
          <input v-model="form.seo_title" type="text" class="ta-input w-full" placeholder="请输入SEO标题（建议60字符以内）" />
          <p class="text-xs text-gray-400 mt-1">当前: {{ form.seo_title?.length || 0 }} / 60</p>
        </div>

        <div>
          <label class="ta-label">SEO描述</label>
          <textarea v-model="form.seo_description" class="ta-input w-full" rows="3" placeholder="请输入SEO描述（建议150字符以内）"></textarea>
          <p class="text-xs text-gray-400 mt-1">当前: {{ form.seo_description?.length || 0 }} / 150</p>
        </div>

        <div>
          <label class="ta-label">SEO关键词</label>
          <input v-model="form.seo_keywords" type="text" class="ta-input w-full" placeholder="多个关键词用逗号分隔" />
        </div>

        <div>
          <label class="ta-label">自定义URL</label>
          <input v-model="form.custom_url" type="text" class="ta-input w-full" placeholder="如: /custom-page" />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="ta-label">Robots指令</label>
            <select v-model="form.robots" class="ta-input w-full">
              <option value="index,follow">索引+跟随</option>
              <option value="noindex,nofollow">不索引+不跟随</option>
              <option value="index,nofollow">索引+不跟随</option>
              <option value="noindex,follow">不索引+跟随</option>
            </select>
          </div>
          <div>
            <label class="ta-label">更新频率</label>
            <select v-model="form.change_freq" class="ta-input w-full">
              <option value="daily">每天</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
              <option value="yearly">每年</option>
              <option value="never">从不</option>
            </select>
          </div>
        </div>

        <div>
          <label class="ta-label">优先级</label>
          <input v-model.number="form.priority" type="range" min="0" max="1" step="0.1" class="w-full" />
          <span class="text-sm text-gray-500">{{ form.priority }}</span>
        </div>
      </div>

      <template #footer>
        <TailButton type="default" variant="outline" @click="showModal = false">取消</TailButton>
        <TailButton type="primary" :loading="submitLoading" @click="saveSeo">{{ isEdit ? '保存' : '创建' }}</TailButton>
      </template>
    </TailModal>

    <!-- Sitemap Preview Modal -->
    <TailModal :open="sitemapModalOpen" @update:open="sitemapModalOpen = $event" title="站点地图" width="max-w-3xl">
      <div class="max-h-[600px] overflow-auto">
        <pre class="p-4 bg-gray-50 rounded text-sm whitespace-pre-wrap">{{ sitemapContent || '暂无内容' }}</pre>
      </div>
      <template #footer>
        <TailButton type="default" variant="outline" @click="sitemapModalOpen = false">关闭</TailButton>
        <TailButton type="primary" @click="handleCopySitemap">复制内容</TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getSeoList, createSeoSetting, updateSeoSetting, deleteSeoSetting, generateSitemap, RESOURCE_TYPE_MAP } from '@/api/seo'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface SeoItem {
  id: number
  resource_type: string
  resource_id: number
  seo_title: string
  seo_description: string
  seo_keywords: string
  custom_url: string
  robots: string
  change_freq: string
  priority: number
  created_at: string
  updated_at: string
}

const seoList = ref<SeoItem[]>([])
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const submitLoading = ref(false)
const showModal = ref(false)
const sitemapModalOpen = ref(false)
const isEdit = ref(false)
const searchKeyword = ref('')
const searchResourceType = ref('')
const sitemapContent = ref('')

const articleCount = computed(() => seoList.value.filter(i => i.resource_type === 'article').length)
const videoCount = computed(() => seoList.value.filter(i => i.resource_type === 'video').length)
const pageCount = computed(() => seoList.value.filter(i => i.resource_type === 'page').length)
const totalPages = computed(() => Math.ceil(totalCount.value / pageSize.value))

const form = reactive({
  id: undefined as number | undefined,
  resource_type: '',
  resource_id: undefined as number | undefined,
  seo_title: '',
  seo_description: '',
  seo_keywords: '',
  custom_url: '',
  robots: 'index,follow',
  change_freq: 'daily',
  priority: 0.5
})

const getResourceTypeClass = (type: string) => {
  const classes: Record<string, string> = {
    article: 'bg-blue-100 text-blue-800',
    video: 'bg-green-100 text-green-800',
    page: 'bg-purple-100 text-purple-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const loadSeoList = async () => {
  loading.value = true
  try {
    const res = await getSeoList({
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value || undefined,
      resource_type: searchResourceType.value || undefined
    })
    const data = res.data
    seoList.value = data.list || []
    totalCount.value = data.total || 0
  } catch (error) {
    message.error('加载SEO列表失败')
    console.error('加载SEO列表失败:', error)
  } finally {
    loading.value = false
  }
}

const editSeo = (item: SeoItem) => {
  Object.assign(form, {
    id: item.id,
    resource_type: item.resource_type,
    resource_id: item.resource_id,
    seo_title: item.seo_title || '',
    seo_description: item.seo_description || '',
    seo_keywords: item.seo_keywords || '',
    custom_url: item.custom_url || '',
    robots: item.robots || 'index,follow',
    change_freq: item.change_freq || 'daily',
    priority: item.priority || 0.5
  })
  isEdit.value = true
  showModal.value = true
}

const deleteSeo = async (id: number) => {
  if (!confirm('确定删除该SEO设置吗？')) return
  try {
    await deleteSeoSetting(id)
    message.success('删除成功')
    await loadSeoList()
  } catch (error) {
    message.error('删除失败')
  }
}

const saveSeo = async () => {
  if (!form.resource_type || !form.resource_id || !form.seo_title) {
    message.error('请填写必填项')
    return
  }
  submitLoading.value = true
  try {
    if (isEdit.value && form.id) {
      await updateSeoSetting(form.id, form)
      message.success('更新成功')
    } else {
      await createSeoSetting(form)
      message.success('创建成功')
    }
    showModal.value = false
    await loadSeoList()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadSeoList()
}

const handleReset = () => {
  searchKeyword.value = ''
  searchResourceType.value = ''
  currentPage.value = 1
  loadSeoList()
}

const handleGenerateSitemap = async () => {
  try {
    const res = await generateSitemap()
    // res.data 现在是字符串类型的XML内容（因为设置了 responseType: 'text'）
    if (res.data && res.data.includes('<urlset')) {
      sitemapContent.value = transformApiUrlsToFrontendUrls(res.data)
      sitemapModalOpen.value = true
    } else {
      sitemapContent.value = res.data || '生成失败'
      sitemapModalOpen.value = true
    }
  } catch (error) {
    message.error('生成站点地图失败')
  }
}

// 将API路径转换为前端路由路径
const transformApiUrlsToFrontendUrls = (xmlContent: string): string => {
  // 替换 <loc>/api/xxx</loc> 为 <loc>/xxx</loc>
  return xmlContent.replace(/<loc>\/api\//g, '<loc>/')
}

const handleCopySitemap = async () => {
  try {
    await navigator.clipboard.writeText(sitemapContent.value)
    message.success('已复制到剪贴板')
  } catch {
    message.error('复制失败')
  }
}

onMounted(() => {
  loadSeoList()
})
</script>