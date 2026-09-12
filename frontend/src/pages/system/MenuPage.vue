<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">菜单管理</h1>
        <p class="text-sm text-gray-500">管理系统菜单结构和权限</p>
      </div>
      <TailButton type="primary" @click="addRootMenu">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新增菜单
      </TailButton>
    </div>

    <!-- Menu Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">菜单名称</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">图标</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">路径</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">排序</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="6" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="menus.length === 0" class="animate-pulse">
              <td colspan="6" class="px-6 py-12 text-center text-gray-400">暂无菜单数据</td>
            </tr>
            <tr 
              v-for="menu in menus" 
              :key="menu.id" 
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <div v-for="i in (menu.depth || 0)" :key="i" class="w-6 shrink-0">
                    <svg class="w-4 h-4 text-gray-300" fill="currentColor" viewBox="0 0 4 4">
                      <circle cx="2" cy="2" r="1.5" />
                    </svg>
                  </div>
                  <span class="font-medium text-gray-900">{{ menu.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <span v-if="menu.icon" class="text-lg">{{ getMenuIcon(menu.icon) }}</span>
                <span v-else class="text-gray-400 text-sm">-</span>
              </td>
              <td class="px-6 py-4">
                <code class="text-xs bg-gray-100 px-2 py-1 rounded">{{ menu.path || '-' }}</code>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ menu.sort }}</td>
              <td class="px-6 py-4">
                <span :class="menu.status === 1 ? 'ta-badge ta-badge-success' : 'ta-badge ta-badge-danger'">
                  {{ menu.status === 1 ? '启用' : '禁用' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <TailButton type="default" variant="ghost" size="sm" @click="addChildMenu(menu)">子菜单</TailButton>
                  <TailButton type="default" variant="ghost" size="sm" @click="editMenu(menu)">编辑</TailButton>
                  <TailButton type="danger" variant="ghost" size="sm" @click="deleteMenu(menu.id)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </TailCard>

    <!-- Modal -->
    <TailModal :open="showModal" @update:open="showModal = $event" :title="isEdit ? '编辑菜单' : '新增菜单'" width="max-w-lg">
      <div class="space-y-4">
        <div>
          <label class="ta-label">父菜单</label>
          <select v-model.number="form.parentId" class="ta-input w-full">
            <option :value="0">顶级菜单</option>
            <option v-for="m in menus" :key="m.id" :value="m.id">{{ getIndent(m.depth || 0) }}{{ m.name }}</option>
          </select>
        </div>
        <div>
          <label class="ta-label">菜单名称 <span class="text-red-500">*</span></label>
          <input v-model="form.name" type="text" class="ta-input w-full" placeholder="请输入菜单名称" />
        </div>
        <div>
          <label class="ta-label">图标</label>
          <input v-model="form.icon" type="text" class="ta-input w-full" placeholder="如: HomeIcon, FilmIcon" />
        </div>
        <div>
          <label class="ta-label">路径</label>
          <input v-model="form.path" type="text" class="ta-input w-full" placeholder="如: /dashboard" />
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
        <TailButton type="primary" :loading="submitLoading" @click="saveMenu">{{ isEdit ? '保存' : '创建' }}</TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getMenuTree, createMenu as createMenuApi, updateMenu as updateMenuApi, deleteMenu as deleteMenuApi } from '@/api/menu'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface MenuInfo {
  id: number
  name: string
  icon?: string
  path: string
  sort: number
  status: number
  depth?: number
  parentId?: number
  children?: MenuInfo[]
}

const menus = ref<MenuInfo[]>([])
const loading = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)

const form = reactive({
  id: undefined as number | undefined,
  parentId: 0,
  name: '',
  icon: '',
  path: '',
  sort: 0,
  status: 1
})

const getMenuIcon = (iconName: string) => {
  // Simple icon mapping for display purposes
  const iconMap: Record<string, string> = {
    'Squares2X2Icon': '▦',
    'DocumentTextIcon': '📄',
    'FilmIcon': '🎬',
    'UserIcon': '👤',
    'CogIcon': '⚙️',
    'KeyIcon': '🔑'
  }
  return iconMap[iconName] || '📋'
}

const getIndent = (depth: number) => '  '.repeat(depth || 0)

const flattenMenuTree = (items: MenuInfo[], depth: number = 0): MenuInfo[] => {
  const result: MenuInfo[] = []
  items.forEach(item => {
    result.push({ ...item, depth })
    if (item.children) {
      result.push(...flattenMenuTree(item.children, depth + 1))
    }
  })
  return result
}

const loadMenus = async () => {
  loading.value = true
  try {
    const res = await getMenuTree()
    const rawMenus = (res.data as MenuInfo[]) || []
    menus.value = flattenMenuTree(rawMenus)
  } catch (error) {
    message.error('加载菜单失败')
    console.error('加载菜单失败:', error)
  } finally {
    loading.value = false
  }
}

const addRootMenu = () => {
  Object.assign(form, { id: undefined, parentId: 0, name: '', icon: '', path: '', sort: 0, status: 1 })
  isEdit.value = false
  showModal.value = true
}

const addChildMenu = (m: MenuInfo) => {
  Object.assign(form, { id: undefined, parentId: m.id, name: '', icon: '', path: '', sort: 0, status: 1 })
  isEdit.value = false
  showModal.value = true
}

const editMenu = (m: MenuInfo) => {
  Object.assign(form, { 
    id: m.id,
    parentId: m.parentId || 0, 
    name: m.name, 
    icon: m.icon || '', 
    path: m.path, 
    sort: m.sort, 
    status: m.status 
  })
  isEdit.value = true
  showModal.value = true
}

const deleteMenu = async (id: number) => {
  if (!confirm('确定删除该菜单吗？')) return
  try {
    await deleteMenuApi(id)
    message.success('删除成功')
    await loadMenus()
  } catch (error) {
    message.error('删除失败')
  }
}

const saveMenu = async () => {
  if (!form.name) {
    message.error('请输入菜单名称')
    return
  }
  submitLoading.value = true
  try {
    if (isEdit.value && form.id) {
      await updateMenuApi(form.id, form)
      message.success('更新成功')
    } else {
      await createMenuApi(form)
      message.success('创建成功')
    }
    showModal.value = false
    await loadMenus()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadMenus()
})
</script>