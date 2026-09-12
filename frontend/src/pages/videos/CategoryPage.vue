<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">视频分类管理</h1>
        <p class="text-sm text-gray-500">管理系统视频分类结构</p>
      </div>
      <TailButton type="primary" @click="handleAdd">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        添加分类
      </TailButton>
    </div>

    <!-- Category Tree -->
    <TailCard class="overflow-hidden">
      <template #header>
        <div class="flex items-center justify-between w-full">
          <h3 class="text-lg font-semibold text-gray-900">分类列表</h3>
          <div class="relative max-w-sm flex-1 ml-4">
            <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <input v-model="searchKeyword" type="text" placeholder="搜索分类..." class="ta-input pl-10 w-full" />
          </div>
        </div>
      </template>
      
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">分类名称</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">父级</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">排序</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr 
              v-for="category in filteredCategories" 
              :key="category.id" 
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div v-for="i in (category.depth || 0)" :key="i" class="w-6 shrink-0">
                    <svg class="w-4 h-4 text-gray-300" fill="currentColor" viewBox="0 0 4 4">
                      <circle cx="2" cy="2" r="1.5" />
                    </svg>
                  </div>
                  <span class="font-medium text-gray-900">{{ category.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">
                {{ category.parent_id === 0 ? '顶级分类' : parentNames[category.parent_id] || '-' }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ category.sort }}</td>
              <td class="px-6 py-4">
                <span :class="category.status === 1 ? 'ta-badge ta-badge-success' : 'ta-badge ta-badge-danger'">
                  {{ category.status === 1 ? '启用' : '禁用' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <TailButton type="success" variant="ghost" size="sm" @click="handleAddChild(category)">
                    <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                    </svg>
                    新建
                  </TailButton>
                  <TailButton type="default" variant="ghost" size="sm" @click="handleEdit(category)">编辑</TailButton>
                  <TailButton type="danger" variant="ghost" size="sm" @click="handleDelete(category)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        
        <!-- Empty State -->
        <div v-if="filteredCategories.length === 0" class="text-center py-12">
          <svg class="w-12 h-12 mx-auto text-gray-300 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 8V3h4z" />
          </svg>
          <p class="text-gray-500">暂无分类数据</p>
        </div>
      </div>
    </TailCard>

    <!-- Create/Edit Modal -->
    <TailModal 
      :open="showModal" 
      @update:open="showModal = $event"
      :title="isEdit ? '编辑分类' : '添加分类'"
      width="max-w-lg"
    >
      <div class="space-y-4">
        <div>
          <label class="ta-label">父分类</label>
          <select v-model="form.parent_id" class="ta-input w-full">
            <option :value="0">顶级分类</option>
            <option v-for="p in parentOptions" :key="p.id" :value="p.id">{{ getIndent(p.depth || 0) }}{{ p.name }}</option>
          </select>
        </div>
        <div>
          <label class="ta-label">分类名称</label>
          <input v-model="form.name" type="text" class="ta-input w-full" placeholder="请输入分类名称" />
        </div>
        <div>
          <label class="ta-label">排序</label>
          <input v-model.number="form.sort" type="number" class="ta-input w-full" placeholder="数字越小越靠前" />
        </div>
        <div>
          <label class="ta-label">状态</label>
          <select v-model="form.status" class="ta-input w-full">
            <option :value="1">启用</option>
            <option :value="0">禁用</option>
          </select>
        </div>
      </div>
      <template #footer>
        <TailButton type="default" variant="outline" @click="showModal = false">取消</TailButton>
        <TailButton type="primary" :loading="submitLoading" @click="saveCategory">
          {{ isEdit ? '保存' : '创建' }}
        </TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getVideoCategoryTree, createVideoCategory, updateVideoCategory, deleteVideoCategory } from '@/api/video'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface CategoryItem {
  id: number
  name: string
  parent_id: number
  sort: number
  status: number
  depth?: number
  children?: CategoryItem[]
}

const categories = ref<CategoryItem[]>([])
const searchKeyword = ref('')
const showModal = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const editingId = ref(0)

const form = reactive({
  name: '',
  parent_id: 0,
  sort: 0,
  status: 1
})

// Flatten category tree for table display
const flatCategories = computed(() => {
  const result: CategoryItem[] = []
  const flatten = (items: CategoryItem[], depth: number = 0) => {
    items.forEach(item => {
      result.push({ ...item, depth })
      if (item.children) {
        flatten(item.children, depth + 1)
      }
    })
  }
  flatten(categories.value)
  return result
})

const parentNames = computed(() => {
  const map: Record<number, string> = {}
  const buildMap = (items: CategoryItem[]) => {
    items.forEach(item => {
      map[item.id] = item.name
      if (item.children) buildMap(item.children)
    })
  }
  buildMap(categories.value)
  return map
})

const filteredCategories = computed(() => {
  if (!searchKeyword.value) return flatCategories.value
  return flatCategories.value.filter(c => 
    c.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const parentOptions = computed(() => {
  // Exclude the current editing category from parent options
  if (isEdit.value) {
    return categories.value.filter(c => c.id !== editingId.value)
  }
  return categories.value
})

const getIndent = (depth: number) => '  '.repeat(depth || 0)

// Load categories
const loadCategories = async () => {
  try {
    const res = await getVideoCategoryTree()
    categories.value = (res.data as CategoryItem[]) || []
  } catch (error) {
    message.error('加载分类失败')
    console.error('加载分类失败:', error)
  }
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(form, { name: '', parent_id: 0, sort: 0, status: 1 })
  showModal.value = true
}

const handleAddChild = (parentCategory: CategoryItem) => {
  isEdit.value = false
  editingId.value = 0
  Object.assign(form, { name: '', parent_id: parentCategory.id, sort: 0, status: 1 })
  showModal.value = true
}

const handleEdit = (category: CategoryItem) => {
  isEdit.value = true
  editingId.value = category.id
  Object.assign(form, { 
    name: category.name, 
    parent_id: category.parent_id, 
    sort: category.sort, 
    status: category.status 
  })
  showModal.value = true
}

const handleDelete = async (category: CategoryItem) => {
  if (!confirm(`确定删除分类 "${category.name}" 吗？`)) return
  try {
    await deleteVideoCategory(category.id)
    message.success('删除成功')
    await loadCategories()
  } catch (error) {
    message.error('删除失败')
  }
}

const saveCategory = async () => {
  if (!form.name) {
    message.error('请输入分类名称')
    return
  }
  submitLoading.value = true
  try {
      if (isEdit.value) {
        await updateVideoCategory(editingId.value, form)
        message.success('更新成功')
      } else {
        await createVideoCategory(form)
      message.success('创建成功')
    }
    showModal.value = false
    await loadCategories()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadCategories()
})
</script>