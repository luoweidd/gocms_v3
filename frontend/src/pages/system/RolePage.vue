<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 mb-1">角色权限管理</h1>
        <p class="text-sm text-gray-500">管理系统角色和权限分配</p>
      </div>
      <TailButton type="primary" @click="showModal = true">
        <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新增角色
      </TailButton>
    </div>

    <!-- Statistics Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-indigo-500">
        <p class="text-sm text-gray-500">总角色数</p>
        <p class="text-2xl font-bold text-gray-900">{{ roles.length }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-emerald-500">
        <p class="text-sm text-gray-500">启用中</p>
        <p class="text-2xl font-bold text-emerald-600">{{ activeRoles }}</p>
      </div>
      <div class="bg-white rounded-lg shadow-sm p-4 border-l-4 border-gray-400">
        <p class="text-sm text-gray-500">已禁用</p>
        <p class="text-2xl font-bold text-gray-600">{{ disabledRoles }}</p>
      </div>
    </div>

    <!-- Role Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="role in roles" :key="role.id" class="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200 p-5">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-xl bg-indigo-100 text-indigo-600 flex items-center justify-center">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div>
              <h3 class="font-semibold text-gray-900">{{ role.name }}</h3>
              <p class="text-xs text-gray-500">{{ role.code }}</p>
            </div>
          </div>
          <span :class="role.status === 1 ? 'ta-badge ta-badge-success' : 'ta-badge ta-badge-danger'">
            {{ role.status === 1 ? '启用' : '禁用' }}
          </span>
        </div>
        
        <div class="mb-4">
          <p class="text-xs text-gray-500 mb-2">权限 ({{ role.permissions?.length || 0 }})</p>
          <div class="flex flex-wrap gap-1">
            <span v-for="(perm, idx) in (role.permissions || []).slice(0, 5)" :key="idx" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-800">
              {{ perm }}
            </span>
            <span v-if="role.permissions && role.permissions.length > 5" class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-600">
              +{{ role.permissions.length - 5 }}
            </span>
          </div>
        </div>

        <div class="mb-4">
          <p class="text-xs text-gray-500 mb-2">用户数: {{ role.userCount || 0 }}</p>
        </div>

        <div class="flex gap-2 pt-4 border-t border-gray-100">
          <TailButton type="default" variant="outline" size="sm" class="flex-1" @click="editRole(role)">编辑</TailButton>
          <TailButton v-if="role.id > 3" type="danger" variant="ghost" size="sm" @click="deleteRole(role.id)">删除</TailButton>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="roles.length === 0 && !loading" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      </svg>
      <p class="text-gray-500">暂无角色数据</p>
    </div>

    <!-- Modal -->
    <TailModal :open="showModal" @update:open="showModal = $event" :title="isEdit ? '编辑角色' : '新增角色'" width="max-w-lg">
      <div class="space-y-4">
        <div>
          <label class="ta-label">角色名称 <span class="text-red-500">*</span></label>
          <input v-model="form.name" type="text" class="ta-input w-full" placeholder="请输入角色名称" />
        </div>
        <div>
          <label class="ta-label">角色编码 <span class="text-red-500">*</span></label>
          <input v-model="form.code" type="text" class="ta-input w-full" placeholder="如: admin, editor" :disabled="isEdit" />
        </div>
        <div>
          <label class="ta-label">描述</label>
          <textarea v-model="form.description" class="ta-input w-full" rows="3" placeholder="请输入角色描述"></textarea>
        </div>
        <div>
          <label class="ta-label">状态</label>
          <select v-model="form.status" class="ta-input w-full">
            <option :value="1">启用</option>
            <option :value="0">禁用</option>
          </select>
        </div>
        <div>
          <label class="ta-label">选择权限</label>
          <div class="space-y-2 max-h-48 overflow-y-auto border border-gray-200 rounded-lg p-3">
            <label v-for="perm in allPermissions" :key="perm" class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer hover:text-indigo-600">
              <input type="checkbox" v-model="form.permissions" :value="perm" class="rounded border-gray-300 text-indigo-600 focus:ring-indigo-500" />
              {{ perm }}
            </label>
          </div>
        </div>
      </div>
      <template #footer>
        <TailButton type="default" variant="outline" @click="showModal = false">取消</TailButton>
        <TailButton type="primary" :loading="submitLoading" @click="saveRole">{{ isEdit ? '保存' : '创建' }}</TailButton>
      </template>
    </TailModal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getRoleList, createRole as createRoleApi, updateRole as updateRoleApi, deleteRole as deleteRoleApi } from '@/api/role'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailModal from '@/components/ui/TailModal.vue'

interface RoleInfo {
  id: number
  name: string
  code: string
  description?: string
  status: number
  permissions: string[]
  userCount?: number
}

const roles = ref<RoleInfo[]>([])
const loading = ref(false)
const showModal = ref(false)
const isEdit = ref(false)
const submitLoading = ref(false)

const form = reactive({
  id: undefined as number | undefined,
  name: '',
  code: '',
  description: '',
  status: 1,
  permissions: [] as string[]
})

const allPermissions = [
  '用户管理', '角色管理', '菜单管理', '内容管理', 
  '评论管理', '视频管理', '系统设置', '数据分析'
]

const activeRoles = computed(() => roles.value.filter(r => r.status === 1).length)
const disabledRoles = computed(() => roles.value.filter(r => r.status === 0).length)

const loadRoles = async () => {
  loading.value = true
  try {
    const res = await getRoleList({ page: 1, page_size: 100 })
    const rawData = (res.data?.list || []) as any[]
    roles.value = rawData.map((item: any) => ({
      id: item.id,
      name: item.name,
      code: item.code,
      description: item.description || '',
      status: item.status ?? 1,
      permissions: item.permissions || [],
      userCount: item.user_count || 0
    }))
  } catch (error) {
    message.error('加载角色列表失败')
    console.error('加载角色列表失败:', error)
  } finally {
    loading.value = false
  }
}

const editRole = (r: RoleInfo) => {
  Object.assign(form, { 
    id: r.id,
    name: r.name, 
    code: r.code, 
    description: r.description || '', 
    status: r.status, 
    permissions: [...(r.permissions || [])] 
  })
  isEdit.value = true
  showModal.value = true
}

const deleteRole = async (id: number) => {
  if (!confirm('确定删除该角色吗？')) return
  try {
    await deleteRoleApi(id)
    message.success('删除成功')
    await loadRoles()
  } catch (error) {
    message.error('删除失败')
  }
}

const saveRole = async () => {
  if (!form.name) {
    message.error('请输入角色名称')
    return
  }
  if (!form.code) {
    message.error('请输入角色编码')
    return
  }
  submitLoading.value = true
  try {
    if (isEdit.value && form.id) {
      await updateRoleApi(form.id, form)
      message.success('更新成功')
    } else {
      await createRoleApi(form)
      message.success('创建成功')
    }
    showModal.value = false
    await loadRoles()
  } catch (error) {
    message.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

onMounted(() => {
  loadRoles()
})
</script>