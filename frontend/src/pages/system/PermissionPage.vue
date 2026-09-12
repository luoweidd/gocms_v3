<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">权限列表</h1>
        <p class="text-sm text-gray-500">管理系统权限与NGAC策略配置</p>
      </div>
      <div class="flex gap-2">
        <TailButton variant="outline" size="sm" @click="handleSync">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          同步策略
        </TailButton>
        <TailButton type="primary" @click="showCreateDialog = true">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新建权限
        </TailButton>
      </div>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.keyword" type="text" placeholder="搜索权限..." class="ta-input pl-10 w-full" @input="debounceSearch" />
        </div>
        <select v-model="searchParams.resource_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部资源</option>
          <option value="system">系统</option>
          <option value="menu">菜单</option>
          <option value="role">角色</option>
          <option value="article">文章</option>
          <option value="video">视频</option>
          <option value="user">用户</option>
        </select>
        <select v-model="searchParams.permission_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部类型</option>
          <option value="read">读取</option>
          <option value="write">写入</option>
          <option value="delete">删除</option>
          <option value="admin">管理</option>
        </select>
        <TailButton variant="outline" size="sm" @click="resetSearch">重置</TailButton>
      </div>
    </TailCard>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总权限数</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-blue-500">
        <p class="text-sm text-gray-500">读取权限</p>
        <p class="text-2xl font-bold text-blue-600">{{ stats.read || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-emerald-500">
        <p class="text-sm text-gray-500">写入权限</p>
        <p class="text-2xl font-bold text-emerald-600">{{ stats.write || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
        <p class="text-sm text-gray-500">删除权限</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.delete || 0 }}</p>
      </div>
    </div>

    <!-- Permission Tree/Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">权限名称</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">权限编码</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">资源类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">权限类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">关联角色</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="8" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="permissions.length === 0" class="animate-pulse">
              <td colspan="8" class="px-6 py-12 text-center text-gray-400">暂无权限数据</td>
            </tr>
            <tr v-for="perm in permissions" :key="perm.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <svg v-if="perm.has_children" class="w-4 h-4 text-gray-400 cursor-pointer" @click="handleExpand(perm)" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="expandedPerms.includes(perm.id) ? 'M19 9l-7 7-7-7' : 'M9 5l7 7-7 7'" />
                  </svg>
                  <span class="text-sm font-medium text-gray-900">{{ perm.name }}</span>
                </div>
              </td>
              <td class="px-6 py-4">
                <code class="text-xs bg-gray-100 px-2 py-1 rounded">{{ perm.code }}</code>
              </td>
              <td class="px-6 py-4">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  {{ resourceTypeLabel[perm.resource_type] || perm.resource_type }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-blue-100 text-blue-800': perm.permission_type === 'read',
                  'bg-green-100 text-green-800': perm.permission_type === 'write',
                  'bg-red-100 text-red-800': perm.permission_type === 'delete',
                  'bg-purple-100 text-purple-800': perm.permission_type === 'admin'
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ permissionTypeLabel[perm.permission_type] || perm.permission_type }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="flex gap-1 flex-wrap">
                  <span v-for="role in perm.roles" :key="role.id" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                    {{ role.name }}
                  </span>
                </div>
              </td>
              <td class="px-6 py-4">
                <button 
                  @click="toggleStatus(perm)"
                  :class="perm.enabled ? 'bg-green-600' : 'bg-gray-300'"
                  class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors"
                >
                  <span class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform" :class="perm.enabled ? 'translate-x-6' : 'translate-x-1'" />
                </button>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">
                {{ formatDate(perm.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex gap-2">
                  <TailButton variant="ghost" size="sm" @click="editPermission(perm)">编辑</TailButton>
                  <TailButton v-if="!perm.has_children" variant="ghost" size="sm" type="danger" @click="deletePermission(perm)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200" v-if="permissions.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page = 1; handleSearch()">首页</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
          <TailButton 
            v-for="p in visiblePages" 
            :key="p"
            class="w-9 h-9 p-0 flex items-center justify-center"
            :class="p === pagination.page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
            @click="pagination.page = p; handleSearch()">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page >= totalPages" @click="pagination.page++; handleSearch()">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Create/Edit Dialog -->
    <Teleport to="body">
      <div v-if="showCreateDialog || showEditDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="closeDialogs">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-lg p-6 m-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ showEditDialog ? '编辑权限' : '新建权限' }}</h3>
          <div class="space-y-4 max-h-[60vh] overflow-y-auto">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">权限名称 <span class="text-red-500">*</span></label>
              <input v-model="formData.name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="请输入权限名称" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">权限编码 <span class="text-red-500">*</span></label>
              <input v-model="formData.code" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="例: system:menu:list" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">资源类型 <span class="text-red-500">*</span></label>
                <select v-model="formData.resource_type" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500">
                  <option value="system">系统</option>
                  <option value="menu">菜单</option>
                  <option value="role">角色</option>
                  <option value="article">文章</option>
                  <option value="video">视频</option>
                  <option value="user">用户</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">权限类型 <span class="text-red-500">*</span></label>
                <select v-model="formData.permission_type" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500">
                  <option value="read">读取</option>
                  <option value="write">写入</option>
                  <option value="delete">删除</option>
                  <option value="admin">管理</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">关联角色</label>
              <div class="space-y-2 max-h-40 overflow-y-auto border border-gray-200 rounded-lg p-3">
                <label v-for="role in allRoles" :key="role.id" class="flex items-center gap-2">
                  <input type="checkbox" :value="role.id" v-model="formData.role_ids" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
                  <span class="text-sm text-gray-700">{{ role.name }}</span>
                </label>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">描述</label>
              <textarea v-model="formData.description" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="权限描述"></textarea>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <TailButton variant="outline" size="sm" @click="closeDialogs">取消</TailButton>
            <TailButton type="primary" size="sm" @click="savePermission">{{ showEditDialog ? '更新' : '创建' }}</TailButton>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'

const permissions = ref<any[]>([])
const allRoles = ref<any[]>([])
const loading = ref(false)
const stats = ref<any>({ total: 0, read: 0, write: 0, delete: 0 })
const expandedPerms = ref<number[]>([])
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const currentPerm = ref<any>(null)

const pagination = reactive({
  page: 1,
  pageSize: 15,
  total: 0
})

const searchParams = reactive({
  keyword: '',
  resource_type: '',
  permission_type: ''
})

const resourceTypeLabel: Record<string, string> = {
  system: '系统',
  menu: '菜单',
  role: '角色',
  article: '文章',
  video: '视频',
  user: '用户'
}

const permissionTypeLabel: Record<string, string> = {
  read: '读取',
  write: '写入',
  delete: '删除',
  admin: '管理'
}

const formData = reactive({
  name: '',
  code: '',
  resource_type: 'system',
  permission_type: 'read',
  role_ids: [] as number[],
  description: ''
})

const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))
const visiblePages = computed(() => {
  const total = totalPages.value
  const current = pagination.page
  const pages = []
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('zh-CN')
  } catch {
    return dateStr
  }
}

let searchTimer: ReturnType<typeof setTimeout> | null = null
const debounceSearch = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    pagination.page = 1
    loadPermissions()
  }, 300)
}

const handleSearch = () => {
  pagination.page = 1
  loadPermissions()
}

const resetSearch = () => {
  Object.assign(searchParams, { keyword: '', resource_type: '', permission_type: '' })
  handleSearch()
}

const loadPermissions = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }
    if (searchParams.keyword) params.keyword = searchParams.keyword
    if (searchParams.resource_type) params.resource_type = searchParams.resource_type
    if (searchParams.permission_type) params.permission_type = searchParams.permission_type

    // 模拟数据 - 实际应该调用API
    const mockData = generateMockPermissions(params)
    permissions.value = mockData.list
    pagination.total = mockData.total
    stats.value = mockData.stats
  } catch (error) {
    message.error('获取权限列表失败')
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  // 模拟角色数据
  allRoles.value = [
    { id: 1, name: '超级管理员' },
    { id: 2, name: '内容管理员' },
    { id: 3, name: '编辑' },
    { id: 4, name: '访客' }
  ]
}

const generateMockPermissions = (params: any) => {
  const list = [
    { id: 1, name: '查看菜单', code: 'system:menu:list', resource_type: 'menu', permission_type: 'read', roles: [{ id: 1, name: '超级管理员' }, { id: 2, name: '内容管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 2, name: '编辑菜单', code: 'system:menu:edit', resource_type: 'menu', permission_type: 'write', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 3, name: '删除菜单', code: 'system:menu:delete', resource_type: 'menu', permission_type: 'delete', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 4, name: '查看角色', code: 'system:role:list', resource_type: 'role', permission_type: 'read', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 5, name: '编辑角色', code: 'system:role:edit', resource_type: 'role', permission_type: 'write', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 6, name: '查看文章', code: 'article:list', resource_type: 'article', permission_type: 'read', roles: [{ id: 1, name: '超级管理员' }, { id: 2, name: '内容管理员' }, { id: 3, name: '编辑' }], enabled: true, has_children: true, created_at: '2024-01-01T00:00:00Z' },
    { id: 7, name: '创建文章', code: 'article:create', resource_type: 'article', permission_type: 'write', roles: [{ id: 1, name: '超级管理员' }, { id: 2, name: '内容管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 8, name: '编辑文章', code: 'article:edit', resource_type: 'article', permission_type: 'write', roles: [{ id: 1, name: '超级管理员' }, { id: 2, name: '内容管理员' }, { id: 3, name: '编辑' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 9, name: '删除文章', code: 'article:delete', resource_type: 'article', permission_type: 'delete', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 10, name: '查看视频', code: 'video:list', resource_type: 'video', permission_type: 'read', roles: [{ id: 1, name: '超级管理员' }, { id: 2, name: '内容管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 11, name: '上传视频', code: 'video:upload', resource_type: 'video', permission_type: 'write', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
    { id: 12, name: '查看用户', code: 'user:list', resource_type: 'user', permission_type: 'read', roles: [{ id: 1, name: '超级管理员' }], enabled: true, has_children: false, created_at: '2024-01-01T00:00:00Z' },
  ]

  let filtered = list
  if (params.keyword) {
    const kw = params.keyword.toLowerCase()
    filtered = filtered.filter(p => p.name.toLowerCase().includes(kw) || p.code.toLowerCase().includes(kw))
  }
  if (params.resource_type) {
    filtered = filtered.filter(p => p.resource_type === params.resource_type)
  }
  if (params.permission_type) {
    filtered = filtered.filter(p => p.permission_type === params.permission_type)
  }

  return {
    list: filtered,
    total: filtered.length,
    stats: {
      total: list.length,
      read: list.filter(p => p.permission_type === 'read').length,
      write: list.filter(p => p.permission_type === 'write').length,
      delete: list.filter(p => p.permission_type === 'delete').length
    }
  }
}

const handleExpand = (perm: any) => {
  const index = expandedPerms.value.indexOf(perm.id)
  if (index > -1) {
    expandedPerms.value.splice(index, 1)
  } else {
    expandedPerms.value.push(perm.id)
  }
}

const toggleStatus = async (perm: any) => {
  const newStatus = perm.enabled ? 0 : 1
  try {
    // TODO: 调用API更新状态
    perm.enabled = newStatus === 1
    message.success('状态更新成功')
  } catch (error) {
    message.error('状态更新失败')
  }
}

const editPermission = (perm: any) => {
  currentPerm.value = perm
  Object.assign(formData, {
    name: perm.name,
    code: perm.code,
    resource_type: perm.resource_type,
    permission_type: perm.permission_type,
    role_ids: perm.roles.map((r: any) => r.id),
    description: perm.description || ''
  })
  showEditDialog.value = true
}

const deletePermission = async (perm: any) => {
  if (!confirm(`确定删除权限"${perm.name}"吗？`)) return
  try {
    // TODO: 调用API删除
    message.success('删除成功')
    loadPermissions()
  } catch (error) {
    message.error('删除失败')
  }
}

const savePermission = async () => {
  if (!formData.name || !formData.code) {
    message.warning('请填写必填项')
    return
  }
  try {
    if (showEditDialog.value) {
      // TODO: 调用API更新
      message.success('更新成功')
    } else {
      // TODO: 调用API创建
      message.success('创建成功')
    }
    closeDialogs()
    loadPermissions()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

const closeDialogs = () => {
  showCreateDialog.value = false
  showEditDialog.value = false
}

const handleSync = () => {
  message.info('策略同步功能开发中...')
}

onMounted(() => {
  loadPermissions()
  loadRoles()
})
</script>