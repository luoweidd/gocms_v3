<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 mb-1">消息类别管理</h1>
          <p class="text-sm text-gray-500">管理系统消息的类别和类型配置</p>
        </div>
        <TailButton type="primary" @click="showAddDialog = true">
          <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新增类别
        </TailButton>
      </div>
    </div>

    <!-- Category Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <TailCard v-for="category in categories" :key="category.id" class="hover:shadow-lg transition-shadow">
        <div class="flex items-start justify-between">
          <div class="flex-1">
            <div class="flex items-center gap-3 mb-2">
              <div class="w-10 h-10 rounded-lg flex items-center justify-center" :class="category.bg_color">
                <svg class="w-5 h-5" :class="category.text_color" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                </svg>
              </div>
              <div>
                <h3 class="font-semibold text-gray-900">{{ category.name }}</h3>
                <p class="text-xs text-gray-500">代码: {{ category.code }}</p>
              </div>
            </div>
            <p class="text-sm text-gray-600 mb-3">{{ category.description }}</p>
            <div class="flex items-center gap-4 text-xs text-gray-500">
              <span>消息数: {{ category.message_count }}</span>
              <span>状态: 
                <span class="px-2 py-0.5 rounded-full" :class="category.status === 1 ? 'bg-green-100 text-green-600' : 'bg-gray-100 text-gray-600'">
                  {{ category.status === 1 ? '启用' : '禁用' }}
                </span>
              </span>
            </div>
          </div>
          <div class="flex flex-col gap-2">
            <button @click="editCategory(category)" class="p-1 text-gray-400 hover:text-primary">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
            <button @click="deleteCategory(category.id)" class="p-1 text-gray-400 hover:text-red-500">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </TailCard>
    </div>

    <!-- Empty State -->
    <div v-if="categories.length === 0 && !loading" class="py-12 text-center text-gray-500">
      <svg class="w-16 h-16 mx-auto mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
      </svg>
      <p>暂无消息类别</p>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="py-12 text-center text-gray-400">
      <svg class="w-8 h-8 mx-auto mb-2 animate-spin text-primary" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <p>加载中...</p>
    </div>

    <!-- Add/Edit Dialog -->
    <div v-if="showAddDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">{{ editingId ? '编辑类别' : '新增类别' }}</h3>
          <button @click="closeDialog" class="text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">类别名称</label>
            <input v-model="form.name" placeholder="请输入类别名称" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">代码</label>
            <input v-model="form.code" placeholder="请输入代码（英文）" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
            <textarea v-model="form.description" rows="3" placeholder="请输入类别描述" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20 resize-none"></textarea>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">图标类型</label>
            <select v-model="form.icon_type" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20">
              <option value="bell">铃铛</option>
              <option value="document">文档</option>
              <option value="video">视频</option>
              <option value="user">用户</option>
              <option value="check">勾选</option>
              <option value="alert">警告</option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
              <input v-model.number="form.sort" type="number" placeholder="数字越小越靠前" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
              <select v-model="form.status" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20">
                <option :value="1">启用</option>
                <option :value="0">禁用</option>
              </select>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-3 p-6 border-t border-gray-200">
          <TailButton variant="outline" @click="closeDialog">取消</TailButton>
          <TailButton type="primary" :loading="saving" @click="saveCategory">保存</TailButton>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import { getMessageCategories, createMessageCategory, updateMessageCategory, deleteMessageCategory } from '@/api/notification'

interface CategoryItem {
  id: number
  name: string
  code: string
  description: string
  icon_type: string
  sort: number
  status: number
  message_count: number
  bg_color: string
  text_color: string
}

const categories = ref<CategoryItem[]>([])
const showAddDialog = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const loading = ref(false)

const defaultForm = {
  name: '',
  code: '',
  description: '',
  icon_type: 'bell',
  sort: 0,
  status: 1
}

const form = reactive({ ...defaultForm })

// 图标颜色配置
const iconColorMap: Record<string, { bg: string; text: string }> = {
  bell: { bg: 'bg-blue-100', text: 'text-blue-600' },
  document: { bg: 'bg-green-100', text: 'text-green-600' },
  video: { bg: 'bg-purple-100', text: 'text-purple-600' },
  user: { bg: 'bg-orange-100', text: 'text-orange-600' },
  check: { bg: 'bg-emerald-100', text: 'text-emerald-600' },
  alert: { bg: 'bg-red-100', text: 'text-red-600' }
}

onMounted(() => {
  fetchCategories()
})

const fetchCategories = async () => {
  loading.value = true
  try {
    const res = await getMessageCategories({ page: 1, page_size: 100 })
    const list = res.data?.list || res.data?.data?.list || []
    
    // 为每个类别添加颜色配置
    categories.value = list.map((cat: any) => {
      const colors = iconColorMap[cat.icon_type] || iconColorMap.bell
      return {
        ...cat,
        bg_color: colors.bg,
        text_color: colors.text
      }
    })
  } catch (error) {
    console.error('获取类别列表失败:', error)
  } finally {
    loading.value = false
  }
}

const editCategory = (category: CategoryItem) => {
  editingId.value = category.id
  Object.assign(form, {
    name: category.name,
    code: category.code,
    description: category.description,
    icon_type: category.icon_type,
    sort: category.sort,
    status: category.status
  })
  showAddDialog.value = true
}

const deleteCategory = async (id: number) => {
  if (!confirm('确定要删除这个类别吗？')) return
  try {
    await deleteMessageCategory(id)
    message.success('删除成功')
    fetchCategories()
  } catch (error) {
    console.error('删除失败:', error)
    message.error('删除失败')
  }
}

const saveCategory = async () => {
  if (!form.name || !form.code) {
    message.warning('请填写名称和代码')
    return
  }
  saving.value = true
  
  try {
    if (editingId.value) {
      // Update existing category
      await updateMessageCategory(editingId.value, form)
      message.success('更新成功')
    } else {
      // Add new category
      await createMessageCategory(form)
      message.success('创建成功')
    }
    closeDialog()
    fetchCategories()
  } catch (error) {
    console.error('保存失败:', error)
    message.error('保存失败')
  } finally {
    saving.value = false
  }
}

const closeDialog = () => {
  showAddDialog.value = false
  editingId.value = null
  Object.assign(form, defaultForm)
}
</script>