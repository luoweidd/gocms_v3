<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">租户管理</h1>
        <p class="text-sm text-gray-500">管理系统租户配置与配额</p>
      </div>
      <TailButton type="primary" @click="showCreateDialog = true">
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新建租户
      </TailButton>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-2 md:grid-cols-5 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总租户数</p>
        <p class="text-2xl font-bold text-gray-900">{{ stats.total || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-emerald-500">
        <p class="text-sm text-gray-500">已启用</p>
        <p class="text-2xl font-bold text-emerald-600">{{ stats.active || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-gray-500">
        <p class="text-sm text-gray-500">已停用</p>
        <p class="text-2xl font-bold text-gray-600">{{ stats.inactive || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-red-500">
        <p class="text-sm text-gray-500">已过期</p>
        <p class="text-2xl font-bold text-red-600">{{ stats.expired || 0 }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-purple-500">
        <p class="text-sm text-gray-500">Premium</p>
        <p class="text-2xl font-bold text-purple-600">{{ stats.by_plan?.premium || 0 }}</p>
      </div>
    </div>

    <!-- Filter Bar -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchParams.name" type="text" placeholder="搜索租户..." class="ta-input pl-10 w-full" @input="debounceSearch" />
        </div>
        <select v-model="searchParams.plan_type" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部套餐</option>
          <option value="free">免费版</option>
          <option value="standard">标准版</option>
          <option value="premium">高级版</option>
          <option value="enterprise">企业版</option>
        </select>
        <select v-model="searchParams.status" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="1">启用</option>
          <option :value="0">停用</option>
          <option :value="2">过期</option>
        </select>
      </div>
    </TailCard>

    <!-- Tenant Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">租户信息</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">套餐类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户配额</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">存储配额</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">到期时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">加载中...</td>
            </tr>
            <tr v-else-if="tenants.length === 0" class="animate-pulse">
              <td colspan="7" class="px-6 py-12 text-center text-gray-400">暂无租户数据</td>
            </tr>
            <tr v-for="tenant in tenants" :key="tenant.id" class="hover:bg-gray-50 transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-lg bg-indigo-100 text-indigo-600 flex items-center justify-center text-sm font-bold">
                    {{ tenant.name.charAt(0) }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ tenant.name }}</p>
                    <p class="text-xs text-gray-500">{{ tenant.code }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-blue-100 text-blue-800': tenant.plan_type === 'free',
                  'bg-emerald-100 text-emerald-800': tenant.plan_type === 'standard',
                  'bg-purple-100 text-purple-800': tenant.plan_type === 'premium',
                  'bg-gray-100 text-gray-800': tenant.plan_type === 'enterprise'
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ planTypeLabel[tenant.plan_type] || tenant.plan_type }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ tenant.user_count }} / {{ tenant.max_users }}</div>
                <div class="w-24 h-1.5 bg-gray-200 rounded-full mt-1">
                  <div class="h-full bg-indigo-500 rounded-full" :style="{ width: Math.min(100, (tenant.user_count / tenant.max_users) * 100) + '%' }"></div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ formatStorage(tenant.used_storage) }} / {{ formatStorage(tenant.max_storage) }}</div>
                <div class="w-24 h-1.5 bg-gray-200 rounded-full mt-1">
                  <div class="h-full bg-emerald-500 rounded-full" :style="{ width: Math.min(100, (tenant.used_storage / tenant.max_storage) * 100) + '%' }"></div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">
                {{ tenant.expire_at ? formatDate(tenant.expire_at) : '永不过期' }}
              </td>
              <td class="px-6 py-4">
                <span :class="{
                  'bg-emerald-100 text-emerald-800': tenant.status === 1,
                  'bg-gray-100 text-gray-800': tenant.status === 0,
                  'bg-red-100 text-red-800': tenant.status === 2
                }" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ ['', '启用', '停用', '过期'][tenant.status] || '未知' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <TailButton type="info" variant="ghost" size="sm" @click="viewUsers(tenant)">查看用户</TailButton>
                  <TailButton type="warning" variant="ghost" size="sm" @click="editTenant(tenant)">编辑</TailButton>
                  <TailButton v-if="tenant.status === 1" type="danger" variant="ghost" size="sm" @click="toggleStatus(tenant, 0)">停用</TailButton>
                  <TailButton v-else type="success" variant="ghost" size="sm" @click="toggleStatus(tenant, 1)">启用</TailButton>
                  <TailButton type="danger" variant="ghost" size="sm" @click="deleteTenantItem(tenant)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200" v-if="tenants.length > 0">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page = 1; handleSearch()">首页</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="pagination.page <= 1" @click="pagination.page--; handleSearch()">上一页</TailButton>
          <TailButton 
            v-for="p in visiblePages" 
            :key="p"
            type="default" 
            size="sm" 
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
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ showEditDialog ? '编辑租户' : '新建租户' }}</h3>
          <div class="space-y-4 max-h-[60vh] overflow-y-auto">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">租户名称 <span class="text-red-500">*</span></label>
              <input v-model="formData.name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="请输入租户名称" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">租户编码 <span class="text-red-500">*</span></label>
              <input v-model="formData.code" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="请输入租户编码（唯一）" :disabled="showEditDialog" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">联系人</label>
                <input v-model="formData.contact_name" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="联系人姓名" />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">联系电话</label>
                <input v-model="formData.contact_phone" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="联系电话" />
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">联系邮箱</label>
              <input v-model="formData.contact_email" type="email" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="联系邮箱" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">绑定域名</label>
              <input v-model="formData.domain" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" placeholder="例: tenant.example.com" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">套餐类型 <span class="text-red-500">*</span></label>
                <select v-model="formData.plan_type" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500">
                  <option value="free">免费版</option>
                  <option value="standard">标准版</option>
                  <option value="premium">高级版</option>
                  <option value="enterprise">企业版</option>
                </select>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">到期时间</label>
                <input v-model="formData.expire_at" type="date" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500" />
              </div>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <TailButton variant="outline" size="sm" @click="closeDialogs">取消</TailButton>
            <TailButton type="primary" size="sm" @click="saveTenant">{{ showEditDialog ? '更新' : '创建' }}</TailButton>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Users Dialog -->
    <Teleport to="body">
      <div v-if="showUsersDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showUsersDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl p-6 m-4">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">{{ currentTenant?.name }} - 用户列表</h3>
            <TailButton variant="ghost" size="sm" @click="showUsersDialog = false">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </TailButton>
          </div>
          <div class="mb-4">
            <div class="flex gap-3">
              <input v-model="userSearch" type="text" placeholder="搜索用户..." class="ta-input flex-1" @input="loadUsers" />
              <TailButton type="primary" size="sm" @click="addUserToTenantDialog = true">添加用户</TailButton>
            </div>
          </div>
          <div class="max-h-[40vh] overflow-y-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50 sticky top-0">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">用户名</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">昵称</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">邮箱</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                  <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-if="usersLoading" class="animate-pulse">
                  <td colspan="6" class="px-4 py-8 text-center text-gray-400">加载中...</td>
                </tr>
                <tr v-else-if="users.length === 0">
                  <td colspan="6" class="px-4 py-8 text-center text-gray-400">暂无用户</td>
                </tr>
                <tr v-for="user in users" :key="user.id">
                  <td class="px-4 py-3 text-sm text-gray-900">{{ user.id }}</td>
                  <td class="px-4 py-3 text-sm text-gray-900">{{ user.username }}</td>
                  <td class="px-4 py-3 text-sm text-gray-600">{{ user.nickname }}</td>
                  <td class="px-4 py-3 text-sm text-gray-600">{{ user.email }}</td>
                  <td class="px-4 py-3">
                    <span :class="{ 'bg-emerald-100 text-emerald-800': user.status === 1, 'bg-gray-100 text-gray-800': user.status === 0 }" class="inline-flex items-center px-2 py-0.5 rounded text-xs">
                      {{ user.status === 1 ? '启用' : '停用' }}
                    </span>
                  </td>
                  <td class="px-4 py-3 text-right">
                    <TailButton type="danger" variant="ghost" size="sm" @click="removeUser(user.id)">移除</TailButton>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Add User Dialog -->
    <Teleport to="body">
      <div v-if="addUserToTenantDialog" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="addUserToTenantDialog = false">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6 m-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">添加用户到 {{ currentTenant?.name }}</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">选择用户</label>
              <select v-model="selectedUserId" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500">
                <option value="">请选择用户</option>
                <option v-for="u in availableUsers" :key="u.id" :value="u.id">{{ u.nickname }} ({{ u.username }})</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <TailButton variant="outline" size="sm" @click="addUserToTenantDialog = false">取消</TailButton>
            <TailButton type="primary" size="sm" @click="confirmAddUser" :disabled="!selectedUserId">确认添加</TailButton>
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
import { 
  getTenantList, 
  getTenantStats, 
  getTenant,
  createTenant,
  updateTenant,
  deleteTenant,
  updateTenantStatus,
  getTenantUsers,
  addUserToTenant,
  removeUserFromTenant
} from '@/api/tenant'

const tenants = ref<any[]>([])
const loading = ref(false)
const stats = ref<any>({ total: 0, active: 0, inactive: 0, expired: 0, by_plan: {} })

const pagination = reactive({
  page: 1,
  pageSize: 15,
  total: 0
})

const searchParams = reactive({
  name: '',
  plan_type: '',
  status: undefined as number | undefined
})

const planTypeLabel: Record<string, string> = {
  free: '免费版',
  standard: '标准版',
  premium: '高级版',
  enterprise: '企业版'
}

const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showUsersDialog = ref(false)
const addUserToTenantDialog = ref(false)
const currentTenant = ref<any>(null)
const formData = reactive({
  name: '',
  code: '',
  contact_name: '',
  contact_phone: '',
  contact_email: '',
  domain: '',
  plan_type: 'free' as string,
  expire_at: ''
})

const users = ref<any[]>([])
const usersLoading = ref(false)
const userSearch = ref('')
const availableUsers = ref<any[]>([])
const selectedUserId = ref<number | null>(null)

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

const formatStorage = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateStr: string) => {
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
    loadTenants()
  }, 300)
}

const handleSearch = () => {
  pagination.page = 1
  loadTenants()
}

const loadTenants = async () => {
  loading.value = true
  try {
    const res = await getTenantList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...(searchParams.name ? { name: searchParams.name } : {}),
      ...(searchParams.plan_type ? { plan_type: searchParams.plan_type } : {}),
      ...(searchParams.status !== undefined ? { status: searchParams.status } : {})
    })
    tenants.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) {
    message.error('获取租户列表失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const res = await getTenantStats()
    stats.value = res.data || { total: 0, active: 0, inactive: 0, expired: 0, by_plan: {} }
  } catch (error) {
    console.error('获取统计失败', error)
  }
}

const openCreateDialog = () => {
  Object.assign(formData, { name: '', code: '', contact_name: '', contact_phone: '', contact_email: '', domain: '', plan_type: 'free', expire_at: '' })
  showCreateDialog.value = true
}

const editTenant = (tenant: any) => {
  currentTenant.value = tenant
  Object.assign(formData, {
    name: tenant.name,
    code: tenant.code,
    contact_name: tenant.contact_name,
    contact_phone: tenant.contact_phone,
    contact_email: tenant.contact_email,
    domain: tenant.domain,
    plan_type: tenant.plan_type,
    expire_at: tenant.expire_at ? tenant.expire_at.split('T')[0] : ''
  })
  showEditDialog.value = true
}

const closeDialogs = () => {
  showCreateDialog.value = false
  showEditDialog.value = false
}

const saveTenant = async () => {
  if (!formData.name || !formData.code) {
    message.warning('请填写必填项')
    return
  }
  try {
    if (showEditDialog.value && currentTenant.value) {
      await updateTenant(currentTenant.value.id, formData as any)
      message.success('更新成功')
    } else {
      await createTenant(formData as any)
      message.success('创建成功')
    }
    closeDialogs()
    loadTenants()
    loadStats()
  } catch (error: any) {
    message.error(error.message || '操作失败')
  }
}

const toggleStatus = async (tenant: any, status: number) => {
  const action = status === 1 ? '启用' : '停用'
  if (!confirm(`确定${action}该租户吗？`)) return
  try {
    await updateTenantStatus(tenant.id, status)
    message.success(`${action}成功`)
    loadTenants()
    loadStats()
  } catch (error: any) {
    message.error(error.message || `${action}失败`)
  }
}

const deleteTenantItem = async (tenant: any) => {
  if (!confirm(`确定删除租户"${tenant.name}"吗？此操作不可恢复`)) return
  try {
    await deleteTenant(tenant.id)
    message.success('删除成功')
    loadTenants()
    loadStats()
  } catch (error: any) {
    message.error(error.message || '删除失败')
  }
}

const viewUsers = async (tenant: any) => {
  currentTenant.value = tenant
  showUsersDialog.value = true
  await loadUsers()
}

const loadUsers = async () => {
  if (!currentTenant.value) return
  usersLoading.value = true
  try {
    const res = await getTenantUsers(currentTenant.value.id, { page: 1, page_size: 50, ...(userSearch.value ? { keyword: userSearch.value } : {}) })
    users.value = res.data?.list || []
  } catch (error) {
    console.error('获取用户列表失败', error)
  } finally {
    usersLoading.value = false
  }
}

const removeUser = async (userId: number) => {
  if (!currentTenant.value) return
  if (!confirm('确定从该租户移除此用户吗？')) return
  try {
    await removeUserFromTenant(currentTenant.value.id, userId)
    message.success('移除成功')
    await loadUsers()
    loadTenants()
  } catch (error: any) {
    message.error(error.message || '移除失败')
  }
}

const loadAvailableUsers = async () => {
  try {
    const res = await getTenantList({ page: 1, page_size: 1000 })
    availableUsers.value = (res.data?.list || []).filter((t: any) => t.id !== currentTenant.value?.id).flatMap(t => [])
  } catch (error) {
    console.error('获取可用用户失败', error)
  }
}

const confirmAddUser = async () => {
  if (!selectedUserId.value || !currentTenant.value) return
  try {
    await addUserToTenant(currentTenant.value.id, selectedUserId.value)
    message.success('添加成功')
    addUserToTenantDialog.value = false
    loadUsers()
    loadTenants()
  } catch (error: any) {
    message.error(error.message || '添加失败')
  }
}

onMounted(() => {
  loadTenants()
  loadStats()
})
</script>