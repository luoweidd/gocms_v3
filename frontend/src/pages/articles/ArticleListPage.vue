p<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">文章管理</h1>
        <p class="text-sm text-gray-500">管理系统文章内容</p>
      </div>
      <router-link to="/articles/create">
        <TailButton type="primary">
          <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新建文章
        </TailButton>
      </router-link>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索文章..." class="ta-input pl-10 w-full" @input="handleSearch" />
        </div>
        <select v-model="searchParams.categoryId" class="ta-input w-auto min-w-[140px]" @change="handleSearch">
          <option value="">全部分类</option>
          <option v-for="cat in categoryTree" :key="cat.id" :value="cat.id">{{ getIndent(cat.depth || 0) }}{{ cat.name }}</option>
        </select>
        <select v-model="searchParams.status" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="0">草稿</option>
          <option :value="1">已发布</option>
          <option :value="2">已下架</option>
        </select>
      </div>
    </TailCard>

    <!-- Data Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">分类</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">浏览量</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">
                <svg class="w-8 h-8 mx-auto mb-2 animate-spin text-indigo-600" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                加载中...
              </td>
            </tr>
            <tr v-else-if="articleList.length === 0" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">
                <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                暂无文章数据
              </td>
            </tr>
            <tr 
              v-for="article in articleList" 
              :key="article.id" 
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 text-sm text-gray-500">#{{ article.id }}</td>
              <td class="px-6 py-4">
                <router-link :to="`/articles/${article.id}`" class="text-sm font-medium text-indigo-600 hover:text-indigo-700">
                  {{ article.title }}
                </router-link>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ article.category_name || '-' }}</td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-emerald-100 text-emerald-800': article.status === 1,
                  'bg-amber-100 text-amber-800': article.status === 0,
                  'bg-red-100 text-red-800': article.status === 2
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ article.status === 1 ? '已发布' : article.status === 0 ? '草稿' : '已下架' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ article.view_count || 0 }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(article.created_at) }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <router-link :to="`/articles/${article.id}`">
                    <TailButton type="success" variant="ghost" size="sm">
                      <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                      </svg>
                      预览
                    </TailButton>
                  </router-link>
                  <router-link :to="`/articles/edit/${article.id}`">
                    <TailButton type="default" variant="ghost" size="sm">编辑</TailButton>
                  </router-link>
                  <TailButton type="danger" variant="ghost" size="sm" @click="deleteArticle(article)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
          <TailButton 
            v-for="p in totalPages" 
            :key="p"
            type="default" 
            size="sm" 
            class="w-9 h-9 p-0 flex items-center justify-center"
            :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
            @click="goToPage(p)">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Create/Edit Modal -->
    <TailModal :open="dialogVisible" @update:open="dialogVisible = $event" :title="isEdit ? '编辑文章' : '新建文章'" width="max-w-3xl">
      <div class="space-y-4">
        <div>
          <label class="ta-label">标题 <span class="text-red-500">*</span></label>
          <input v-model="formData.title" type="text" class="ta-input w-full" placeholder="请输入文章标题" />
        </div>
        <div>
          <label class="ta-label">分类 <span class="text-red-500">*</span></label>
          <select v-model="formData.categoryId" class="ta-input w-full">
            <option :value="undefined">请选择分类</option>
            <option v-for="cat in categoryTree" :key="cat.id" :value="cat.id">{{ getIndent(cat.depth || 0) }}{{ cat.name }}</option>
          </select>
        </div>
        <div>
          <label class="ta-label">内容 <span class="text-red-500">*</span></label>
          <textarea v-model="formData.content" class="ta-input w-full" rows="8" placeholder="请输入文章内容"></textarea>
        </div>
        <div>
          <label class="ta-label">状态</label>
          <select v-model="formData.status" class="ta-input w-full">
            <option :value="0">草稿</option>
            <option :value="1">已发布</option>
            <option :value="2">已下架</option>
          </select>
        </div>
      </div>
      <template #footer>
        <TailButton type="default" variant="outline" @click="handleClose">取消</TailButton>
        <TailButton type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ isEdit ? '保存' : '创建' }}
        </TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getArticleList, deleteArticle as deleteArticleApi, createArticle, updateArticle, getCategoryTree } from '@/api/article'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface ArticleItem {
  id: number
  title: string
  description?: string
  content: string
  author_id: number
  author_name?: string
  category_id: number
  category_name?: string
  status: 0 | 1 | 2
  view_count: number
  created_at: string
  updated_at: string
}

interface CategoryItem {
  id: number
  name: string
  depth?: number
  children?: CategoryItem[]
}

const articleList = ref<ArticleItem[]>([])
const categoryTree = ref<CategoryItem[]>([])
const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const searchParams = reactive({
  keyword: '',
  categoryId: undefined as number | undefined,
  status: undefined as number | undefined
})

const formData = reactive({
  id: undefined as number | undefined,
  title: '',
  categoryId: undefined as number | undefined,
  content: '',
  status: 1 as 0 | 1 | 2
})

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))

const getIndent = (depth: number) => '  '.repeat(depth || 0)

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('zh-CN', { 
      year: 'numeric', 
      month: '2-digit', 
      day: '2-digit' 
    })
  } catch {
    return dateStr
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getArticleList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...(searchParams.keyword ? { keyword: searchParams.keyword } : {}),
      ...(searchParams.categoryId ? { category_id: searchParams.categoryId } : {}),
      ...(searchParams.status != null ? { status: Number(searchParams.status) } : {})
    })
    articleList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    message.error('获取文章列表失败')
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadCategories = async () => {
  try {
    const res = await getCategoryTree()
    categoryTree.value = (res.data as CategoryItem[]) || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const goToPage = (p: number) => {
  pagination.page = p
  fetchData()
}

const showCreateDialog = () => {
  isEdit.value = false
  Object.assign(formData, { id: undefined, title: '', categoryId: undefined, content: '', status: 1 })
  dialogVisible.value = true
}

const editArticle = (article: ArticleItem) => {
  isEdit.value = true
  Object.assign(formData, { 
    id: article.id, 
    title: article.title, 
    categoryId: article.category_id, 
    content: article.content, 
    status: article.status 
  })
  dialogVisible.value = true
}

const deleteArticle = async (article: ArticleItem) => {
  if (!confirm(`确定删除文章 "${article.title}" 吗？`)) return
  try {
    await deleteArticleApi(article.id)
    message.success('删除成功')
    await fetchData()
  } catch (error) {
    message.error('删除失败')
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

const handleSubmit = async () => {
  if (!formData.title) {
    message.error('请输入标题')
    return
  }
  if (!formData.categoryId) {
    message.error('请选择分类')
    return
  }
  submitLoading.value = true
  try {
    if (isEdit.value && formData.id) {
      await updateArticle(formData.id, {
        title: formData.title,
        category_id: formData.categoryId,
        content: formData.content,
        status: formData.status
      })
      message.success('更新成功')
    } else {
      await createArticle({
        title: formData.title,
        category_id: formData.categoryId as number,
        content: formData.content,
        status: formData.status
      })
      message.success('创建成功')
    }
    dialogVisible.value = false
    await fetchData()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  fetchData()
  loadCategories()
})
</script>