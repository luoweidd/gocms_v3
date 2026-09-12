<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">用户管理</h1>
        <p class="text-sm text-gray-500">管理系统用户账号和权限</p>
      </div>
      <TailButton type="primary" @click="showCreateDialog">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新建用户
      </TailButton>
    </div>

    <!-- Search & Filter -->
    <TailCard class="mb-6">
      <div class="flex gap-3 flex-wrap items-center">
        <div class="relative flex-1 min-w-[200px] max-w-md">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchKeyword" type="text" placeholder="搜索用户..." class="ta-input pl-10 w-full" @input="handleSearch" />
        </div>
        <select v-model="roleFilter" class="ta-input w-auto min-w-[120px]" @change="handleSearch">
          <option value="">全部角色</option>
          <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
        </select>
        <select v-model="statusFilter" class="ta-input w-auto min-w-[100px]" @change="handleSearch">
          <option value="">全部状态</option>
          <option :value="1">正常</option>
          <option :value="0">禁用</option>
        </select>
      </div>
    </TailCard>

    <!-- Users Table -->
    <TailCard class="overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">邮箱</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">角色</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="6" class="px-6 py-12 text-center text-gray-400">
                <svg class="w-8 h-8 mx-auto mb-2 animate-spin text-indigo-600" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                加载中...
              </td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0" class="animate-pulse">
              <td colspan="6" class="px-6 py-12 text-center text-gray-400">
                <svg class="w-12 h-12 mx-auto mb-3 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.858M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.858M7 20H2v-2a3 3 0 015.356-1.858M7 20v-2c0-.656.126-1.283.356-1.858m0 0a9.94 9.94 0 014-1 9.94 9.94 0 014 1m-7 4a3 3 0 003 3h4a3 3 0 003-3" />
                </svg>
                暂无用户数据
              </td>
            </tr>
            <tr 
              v-for="user in paginatedUsers" 
              :key="user.id" 
              class="hover:bg-gray-50 transition-colors"
            >
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-full bg-indigo-100 text-indigo-600 flex items-center justify-center font-medium text-sm">
                    {{ user.nickname?.charAt(0) || user.username?.charAt(0) || 'U' }}
                  </div>
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ user.nickname || user.username }}</div>
                    <div class="text-xs text-gray-500">@{{ user.username }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-sm text-gray-600">{{ user.email || '-' }}</td>
              <td class="px-6 py-4">
                <span v-for="(role, idx) in (user.roles || ['普通用户'])" :key="idx" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-800 mr-1">
                  {{ role }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span :class="user.status === 1 ? 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-emerald-100 text-emerald-800' : 'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800'">
                  {{ user.status === 1 ? '正常' : '禁用' }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatDate(user.created_at) }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                  <TailButton type="default" variant="ghost" size="sm" @click="editUser(user)">编辑</TailButton>
                  <TailButton type="danger" variant="ghost" size="sm" @click="deleteUser(user)">删除</TailButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ total }} 条记录，当前第 {{ page }} 页</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="sm" :disabled="page <= 1" @click="page--">上一页</TailButton>
          <TailButton 
            v-for="p in totalPages" 
            :key="p"
            type="default" 
            size="sm" 
            class="w-9 h-9 p-0 flex items-center justify-center"
            :class="p === page ? 'bg-indigo-600 text-white border-indigo-600' : ''"
            @click="page = p"
          >{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="sm" :disabled="page >= totalPages" @click="page++">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Create/Edit Modal -->
    <TailModal :open="showModal" @update:open="showModal = $event" :title="isEdit ? '编辑用户' : '新建用户'" width="max-w-lg">
      <div class="space-y-4">
        <div>
          <label class="ta-label">用户名 <span class="text-red-500">*</span></label>
          <input v-model="form.username" type="text" class="ta-input w-full" placeholder="请输入用户名" />
        </div>
        <div v-if="!isEdit">
          <label class="ta-label">密码 <span class="text-red-500">*</span></label>
          <input v-model="form.password" type="password" class="ta-input w-full" placeholder="请输入密码" @input="validatePassword" />
          <p v-if="passwordError" class="text-xs text-red-500 mt-1">{{ passwordError }}</p>
          <p v-else class="text-xs text-gray-500 mt-1">密码要求：至少8个字符，包含字母和数字</p>
        </div>
        <div>
          <label class="ta-label">邮箱</label>
          <input v-model="form.email" type="email" class="ta-input w-full" placeholder="请输入邮箱" />
        </div>
        <div>
          <label class="ta-label">昵称</label>
          <input v-model="form.nickname" type="text" class="ta-input w-full" placeholder="请输入昵称" />
        </div>
        <div>
          <label class="ta-label">角色</label>
          <select v-model="form.role_id" class="ta-input w-full">
            <option :value="undefined">请选择角色</option>
            <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.name }}</option>
          </select>
        </div>
        <div>
          <label class="ta-label">状态</label>
          <select v-model="form.status" class="ta-input w-full">
            <option :value="1">正常</option>
            <option :value="0">禁用</option>
          </select>
        </div>
      </div>
      <template #footer>
        <TailButton type="default" variant="outline" @click="showModal = false">取消</TailButton>
        <TailButton type="primary" :loading="submitLoading" @click="saveUser">
          {{ isEdit ? '保存' : '创建' }}
        </TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getUserList, createUser, updateUser, deleteUser as deleteUserApi, getRoleList } from '@/api/user'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface UserInfo {
  id: number
  username: string
  nickname?: string
  email?: string
  phone?: string
  role_id?: number
  role_name?: string
  roles: string[]
  status: number
  created_at: string
  updated_at: string
}

interface RoleInfo {
  id: number
  name: string
  code: string
}

const users = ref<UserInfo[]>([])
const roles = ref<RoleInfo[]>([])
const searchKeyword = ref('')
const roleFilter = ref<number | undefined>(undefined)
const statusFilter = ref<number | undefined>(undefined)
const page = ref(1)
const pageSize = ref(10)
const showModal = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)
const loading = ref(false)
const passwordError = ref('')

const validatePassword = () => {
  const pwd = form.password
  if (!pwd) {
    passwordError.value = ''
    return
  }
  const hasLetters = /[a-zA-Z]/.test(pwd)
  const hasNumbers = /\d/.test(pwd)
  if (pwd.length < 8) {
    passwordError.value = '密码长度至少为8个字符'
  } else if (!hasLetters) {
    passwordError.value = '密码必须包含字母'
  } else if (!hasNumbers) {
    passwordError.value = '密码必须包含数字'
  } else {
    passwordError.value = ''
  }
}

const form = reactive({
  id: undefined as number | undefined,
  username: '',
  password: '',
  email: '',
  nickname: '',
  role_id: undefined as number | undefined,
  status: 1
})

// Filtered users
const filteredUsers = computed(() => {
  return users.value.filter(u => {
    const matchKeyword = !searchKeyword.value || 
      u.username.toLowerCase().includes(searchKeyword.value.toLowerCase()) || 
      (u.nickname || '').toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      (u.email || '').toLowerCase().includes(searchKeyword.value.toLowerCase())
    const matchRole = !roleFilter.value || u.role_id === roleFilter.value
    const matchStatus = (statusFilter.value == null) || u.status === statusFilter.value
    return matchKeyword && matchRole && matchStatus
  })
})

const total = computed(() => filteredUsers.value.length)
const totalPages = computed(() => Math.ceil(total.value / pageSize.value))
const paginatedUsers = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredUsers.value.slice(start, start + pageSize.value)
})

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

// Load data
const loadData = async () => {
  loading.value = true
  try {
    const res = await getUserList({
      page: page.value,
      page_size: pageSize.value,
      keyword: searchKeyword.value || undefined,
      role_id: roleFilter.value,
      status: statusFilter.value || undefined
    })
    users.value = (res.data?.list || []) as UserInfo[]
  } catch (error) {
    message.error('加载用户列表失败')
    console.error('加载用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

const loadRoles = async () => {
  try {
    const res = await getRoleList()
    roles.value = (res.data?.list || []) as RoleInfo[]
  } catch (error) {
    console.error('加载角色列表失败:', error)
  }
}

const handleSearch = () => {
  page.value = 1
}

const showCreateDialog = () => {
  isEdit.value = false
  Object.assign(form, { username: '', password: '', email: '', nickname: '', role_id: undefined, status: 1 })
  showModal.value = true
}

const editUser = (user: UserInfo) => {
  isEdit.value = true
  Object.assign(form, { 
    id: user.id,
    username: user.username, 
    password: '', 
    email: user.email || '', 
    nickname: user.nickname || '', 
    role_id: user.role_id, 
    status: user.status 
  })
  showModal.value = true
}

const deleteUser = async (user: UserInfo) => {
  if (!confirm(`确定删除用户 "${user.username}" 吗？`)) return
  try {
    await deleteUserApi(user.id)
    message.success('删除成功')
    await loadData()
  } catch (error) {
    message.error('删除失败')
  }
}

const saveUser = async () => {
  if (!form.username) {
    message.error('请输入用户名')
    return
  }
  if (!isEdit.value && !form.password) {
    message.error('请输入密码')
    return
  }
  // 前端验证密码强度
  if (!isEdit.value) {
    validatePassword()
    if (passwordError.value) {
      message.error(passwordError.value)
      return
    }
  }
  submitLoading.value = true
  try {
    if (isEdit.value && form.id) {
      await updateUser(form.id, {
        email: form.email || undefined,
        role_id: form.role_id,
        status: form.status
      })
      message.success('更新成功')
    } else {
      // 创建用户时，只发送非 undefined 的值
      const userData: any = {
        username: form.username,
        password: form.password,
        status: form.status
      }
      if (form.email) {
        userData.email = form.email
      }
      if (form.role_id) {
        userData.role_id = form.role_id
      }
      await createUser(userData)
      message.success('创建成功')
    }
    showModal.value = false
    await loadData()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadData()
  loadRoles()
})
</script>