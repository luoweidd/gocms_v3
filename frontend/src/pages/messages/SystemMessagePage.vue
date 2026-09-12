<template>
  <div class="page-container">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900 mb-1">系统消息管理</h1>
          <p class="text-sm text-gray-500">管理系统消息的发布、查看和删除</p>
        </div>
        <TailButton type="primary" @click="openPublishDialog()">
          <svg class="w-4 h-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          发布消息
        </TailButton>
      </div>
    </div>

    <!-- Search Bar -->
    <TailCard class="mb-4">
      <div class="flex gap-4 items-end">
        <div class="flex-1">
          <input v-model="searchKeyword" @input="handleSearch" placeholder="搜索消息标题或内容..." class="w-full border border-gray-200 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20" />
        </div>
        <select v-model="messageType" @change="handleSearch" class="border border-gray-200 rounded-md px-4 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20">
          <option value="">全部类型</option>
          <option value="system">系统通知</option>
          <option value="publish">发布通知</option>
          <option value="comment_reply">评论回复</option>
          <option value="audit_passed">审核通过</option>
          <option value="audited_failed">审核未通过</option>
        </select>
        <TailButton variant="outline" size="sm" @click="handleSearch">搜索</TailButton>
      </div>
    </TailCard>

    <!-- Message List -->
    <TailCard>
      <div v-if="loading" class="py-12 text-center text-gray-400">加载中...</div>
      <div v-else-if="messageList.length === 0" class="py-12 text-center text-gray-500">暂无消息</div>
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">标题</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">类型</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">优先级</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="msg in messageList" :key="msg.id" class="hover:bg-gray-50">
              <td class="px-6 py-4">
                <div class="text-sm font-medium text-gray-900 cursor-pointer hover:text-primary" @click="viewDetail(msg)">{{ msg.title }}</div>
                <div class="text-xs text-gray-500 truncate max-w-xs">{{ msg.content }}</div>
              </td>
              <td class="px-6 py-4"><span class="px-2 py-1 text-xs rounded-full" :class="getMessageTypeClass(msg.message_type) as string">{{ getMessageTypeName(msg.message_type) }}</span></td>
              <td class="px-6 py-4"><span :class="{ 'text-red-600': msg.priority >= 2, 'text-yellow-600': msg.priority === 1 }">{{ getPriorityName(msg.priority) }}</span></td>
              <td class="px-6 py-4"><span class="px-2 py-1 text-xs rounded-full" :class="getStatusClass(msg.status)">{{ getStatusName(msg.status) }}</span></td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ formatTime(msg.created_at) }}</td>
              <td class="px-6 py-4 text-sm">
                <button class="text-primary hover:text-primary-dark mr-3" @click="viewDetail(msg)">查看</button>
                <button v-if="canWithdraw(msg)" class="text-orange-600 hover:text-orange-700 mr-3" @click="withdrawMessage(msg.id)">撤回</button>
                <button class="text-red-600 hover:text-red-700" @click="deleteMessage(msg.id)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="messageList.length > 0" class="flex items-center justify-between px-6 py-4 border-t border-gray-200">
        <span class="text-sm text-gray-500">共 {{ pagination.total }} 条记录</span>
        <div class="flex gap-1">
          <TailButton type="default" variant="outline" size="xs" :disabled="pagination.page <= 1" @click="changePage(pagination.page - 1)">上一页</TailButton>
          <TailButton v-for="p in totalPages" :key="p" type="default" size="xs" class="w-8 h-8 p-0 flex items-center justify-center" :class="p === pagination.page ? 'bg-primary text-white border-primary' : ''" @click="changePage(p)">{{ p }}</TailButton>
          <TailButton type="default" variant="outline" size="xs" :disabled="pagination.page >= totalPages" @click="changePage(pagination.page + 1)">下一页</TailButton>
        </div>
      </div>
    </TailCard>

    <!-- Publish Dialog -->
    <div v-if="showPublishDialog" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-3xl mx-4 max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200 sticky top-0 bg-white z-10">
          <h3 class="text-lg font-semibold text-gray-900">发布系统消息</h3>
          <button @click="showPublishDialog = false" class="text-gray-400 hover:text-gray-600"><svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg></button>
        </div>
        <div class="p-6 space-y-4">
          <div><label class="block text-sm font-medium text-gray-700 mb-1">消息标题</label><input v-model="publishForm.title" placeholder="请输入消息标题" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20" /></div>
          <div><label class="block text-sm font-medium text-gray-700 mb-1">消息内容</label><textarea v-model="publishForm.content" rows="6" placeholder="请输入消息内容" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20 resize-none"></textarea></div>
          <div class="grid grid-cols-2 gap-4">
            <div><label class="block text-sm font-medium text-gray-700 mb-1">消息类型</label><select v-model="publishForm.message_type" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20"><option value="system">系统通知</option><option value="publish">发布通知</option><option value="comment_reply">评论回复</option><option value="audit_passed">审核通过</option><option value="audited_failed">审核未通过</option></select></div>
            <div><label class="block text-sm font-medium text-gray-700 mb-1">优先级</label><select v-model="publishForm.priority" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20"><option :value="0">普通</option><option :value="1">重要</option><option :value="2">紧急</option></select></div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">目标用户</label>
            <select v-model="publishForm.target_type" @change="handleTargetTypeChange" class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20">
              <option value="all">全部用户</option>
              <option value="user">指定用户</option>
              <option value="role">指定角色</option>
            </select>
          </div>
          
          <!-- 指定用户时的可选用户 -->
          <div v-if="publishForm.target_type === 'user'" class="space-y-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">选择用户</label>
            <div class="max-h-48 overflow-y-auto border border-gray-200 rounded-md p-3 bg-gray-50">
              <div v-if="availableUsers.length === 0" class="text-sm text-gray-400 text-center py-4">加载中...</div>
              <div v-else-if="availableUsers.length === 0" class="text-sm text-gray-400 text-center py-4">暂无可用用户</div>
              <div v-else>
                <div v-for="user in availableUsers" :key="user.id" class="flex items-center mb-2 last:mb-0">
                  <input
                    type="checkbox"
                    :value="user.id"
                    v-model="selectedUserIds"
                    :id="'user-' + user.id"
                    class="mr-2 rounded border-gray-300 text-primary focus:ring-primary"
                  />
                  <label :for="'user-' + user.id" class="text-sm text-gray-700 cursor-pointer select-none">
                    {{ user.username }}
                    <span v-if="user.role_name" class="text-xs text-gray-400 ml-2">({{ user.role_name }})</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
          
          <!-- 指定角色时的可选角色 -->
          <div v-if="publishForm.target_type === 'role'" class="space-y-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">选择角色</label>
            <div class="max-h-48 overflow-y-auto border border-gray-200 rounded-md p-3 bg-gray-50">
              <div v-if="availableRoles.length === 0" class="text-sm text-gray-400 text-center py-4">加载中...</div>
              <div v-else-if="availableRoles.length === 0" class="text-sm text-gray-400 text-center py-4">暂无可用角色</div>
              <div v-else>
                <div v-for="role in availableRoles" :key="role.id" class="flex items-center mb-2 last:mb-0">
                  <input
                    type="checkbox"
                    :value="role.id"
                    v-model="selectedRoleIds"
                    :id="'role-' + role.id"
                    class="mr-2 rounded border-gray-300 text-primary focus:ring-primary"
                  />
                  <label :for="'role-' + role.id" class="text-sm text-gray-700 cursor-pointer select-none">
                    {{ role.name }}
                    <span class="text-xs text-gray-400 ml-2">({{ role.code }})</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-3 p-6 border-t border-gray-200 sticky bottom-0 bg-white z-10">
          <TailButton variant="outline" @click="showPublishDialog = false">取消</TailButton>
          <TailButton type="primary" :loading="publishing" @click="handlePublish">发布</TailButton>
        </div>
      </div>
    </div>

    <!-- Detail Dialog -->
    <div v-if="showDetailDialog && selectedMessage" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">消息详情</h3>
          <button @click="showDetailDialog = false" class="text-gray-400 hover:text-gray-600"><svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" /></svg></button>
        </div>
        <div class="p-6 space-y-4">
          <div><label class="text-sm font-medium text-gray-500">标题</label><p class="text-gray-900">{{ selectedMessage.title }}</p></div>
          <div><label class="text-sm font-medium text-gray-500">内容</label><p class="text-gray-900 whitespace-pre-wrap">{{ selectedMessage.content }}</p></div>
          <div class="grid grid-cols-2 gap-4">
            <div><label class="text-sm font-medium text-gray-500">类型</label><p class="text-gray-900">{{ getMessageTypeName(selectedMessage.message_type) }}</p></div>
            <div><label class="text-sm font-medium text-gray-500">优先级</label><p class="text-gray-900">{{ getPriorityName(selectedMessage.priority) }}</p></div>
          </div>
          <div><label class="text-sm font-medium text-gray-500">发送者</label><p class="text-gray-900">{{ selectedMessage.sender_name || '系统' }}</p></div>
          <div><label class="text-sm font-medium text-gray-500">创建时间</label><p class="text-gray-900">{{ formatTime(selectedMessage.created_at) }}</p></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import { getMessages, deleteMessage as apiDeleteMessage, withdrawMessage as apiWithdrawMessage, publishMessage as apiPublishMessage } from '@/api/notification'
import { getUserList, getRoleList } from '@/api/user'

interface MessageItem {
  id: number
  title: string
  content: string
  message_type: string
  priority: number
  status: number
  sender_name?: string
  created_at: string
}

const loading = ref(false)
const publishing = ref(false)
const messageList = ref<MessageItem[]>([])
const searchKeyword = ref('')
const messageType = ref('')
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })
const showPublishDialog = ref(false)
const showDetailDialog = ref(false)
const selectedMessage = ref<MessageItem | null>(null)

const publishForm = reactive({ title: '', content: '', message_type: 'system', priority: 0, target_type: 'all' })

// 可选用户和角色列表
const availableUsers = ref<any[]>([])
const availableRoles = ref<any[]>([])
const selectedUserIds = ref<number[]>([])
const selectedRoleIds = ref<number[]>([])
const usersLoading = ref(false)
const rolesLoading = ref(false)

// 加载可选用户列表
const loadAvailableUsers = async () => {
  try {
    usersLoading.value = true
    const res = await getUserList({ page: 1, page_size: 1000, status: 1 })
    // 响应拦截器返回 { ...response, data: data }，res.data 是 UserListResponse
    availableUsers.value = res.data?.list || []
  } catch (error) {
    console.error('加载用户列表失败:', error)
    availableUsers.value = []
  } finally {
    usersLoading.value = false
  }
}

// 加载可选角色列表
const loadAvailableRoles = async () => {
  try {
    rolesLoading.value = true
    const res = await getRoleList({ page: 1, page_size: 1000 })
    // getRoleList 返回的是任意类型，需要安全访问
    availableRoles.value = (res as any)?.data?.list || (res as any)?.list || []
  } catch (error) {
    console.error('加载角色列表失败:', error)
    availableRoles.value = []
  } finally {
    rolesLoading.value = false
  }
}

// 当目标类型变化时，清空之前选择
const handleTargetTypeChange = () => {
  selectedUserIds.value = []
  selectedRoleIds.value = []
}

// 打开发布对话框时加载用户和角色列表
const openPublishDialog = () => {
  showPublishDialog.value = true
  // 重置状态
  availableUsers.value = []
  availableRoles.value = []
  selectedUserIds.value = []
  selectedRoleIds.value = []
  usersLoading.value = false
  rolesLoading.value = false
  // 异步加载数据
  loadAvailableUsers()
  loadAvailableRoles()
}

onMounted(() => { fetchMessages() })

const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await getMessages({ page: pagination.page, page_size: pagination.pageSize, keyword: searchKeyword.value, message_type: messageType.value })
    // 响应拦截器已将 data 直接放到 res.data，所以 res.data 就是 { list: [...], page, pageSize, total }
    messageList.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (error) { console.error('获取消息列表失败:', error) } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchMessages() }
const changePage = (page: number) => { if (page >= 1 && page <= totalPages.value) { pagination.page = page; fetchMessages() } }
const totalPages = computed(() => Math.max(1, Math.ceil(pagination.total / pagination.pageSize)))

const viewDetail = (msg: MessageItem) => { selectedMessage.value = msg; showDetailDialog.value = true }

const canWithdraw = (msg: MessageItem) => { if (!msg.id) return false; const now = new Date(); const created = new Date(msg.created_at); return (now.getTime() - created.getTime()) < 30 * 60 * 1000 }

const withdrawMessage = async (id: number) => {
  if (!confirm('确定要撤回这条消息吗？')) return
  try { await apiWithdrawMessage(id); message.success('消息已撤回'); fetchMessages() } catch (error) { console.error('撤回失败:', error) }
}

const deleteMessage = async (id: number) => {
  if (!confirm('确定要删除这条消息吗？')) return
  try { await apiDeleteMessage(id); message.success('删除成功'); fetchMessages() } catch (error) { console.error('删除失败:', error) }
}

const handlePublish = async () => {
  if (!publishForm.title || !publishForm.content) { message.warning('请填写标题和内容'); return }
  
  // 验证指定用户或指定角色时必须选择至少一个
  if (publishForm.target_type === 'user' && selectedUserIds.value.length === 0) {
    message.warning('请至少选择一个用户')
    return
  }
  if (publishForm.target_type === 'role' && selectedRoleIds.value.length === 0) {
    message.warning('请至少选择一个角色')
    return
  }
  
  publishing.value = true
  try {
    // 根据目标类型构建请求数据
    const requestData: any = {
      title: publishForm.title,
      content: publishForm.content,
      message_type: publishForm.message_type,
      priority: publishForm.priority,
      target_type: publishForm.target_type,
    }
    
    // 添加对应的目标ID
    if (publishForm.target_type === 'user') {
      requestData.target_user_ids = selectedUserIds.value
    } else if (publishForm.target_type === 'role') {
      requestData.target_role_ids = selectedRoleIds.value
    }
    
    await apiPublishMessage(requestData)
    message.success('消息发布成功')
    showPublishDialog.value = false
    Object.assign(publishForm, { title: '', content: '', message_type: 'system', priority: 0, target_type: 'all' })
    selectedUserIds.value = []
    selectedRoleIds.value = []
    fetchMessages()
  } catch (error) { console.error('发布失败:', error) } finally { publishing.value = false }
}

const getMessageTypeClass = (type: string) => ({ 'bg-blue-100 text-blue-600': type === 'system', 'bg-green-100 text-green-600': type === 'publish', 'bg-purple-100 text-purple-600': type === 'comment_reply', 'bg-emerald-100 text-emerald-600': type === 'audit_passed', 'bg-red-100 text-red-600': type === 'audited_failed', 'bg-gray-100 text-gray-600': true }[type] || 'bg-gray-100 text-gray-600')
const getMessageTypeName = (type: string) => ({ system: '系统通知', publish: '发布通知', comment_reply: '评论回复', audit_passed: '审核通过', audited_failed: '审核驳回' }[type] || type)
const getPriorityName = (priority: number) => ({ 0: '普通', 1: '重要', 2: '紧急' }[priority] || '普通')
const getStatusClass = (status: number) => ({ 1: 'bg-green-100 text-green-600', 0: 'bg-gray-100 text-gray-600', 2: 'bg-orange-100 text-orange-600' }[status] || 'bg-gray-100 text-gray-600')
const getStatusName = (status: number) => ({ 1: '已发布', 0: '草稿', 2: '已撤回' }[status] || '未知')

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr); const now = new Date(); const diff = now.getTime() - date.getTime()
  if (diff < 60 * 1000) return '刚刚'
  if (diff < 60 * 60 * 1000) return `${Math.floor(diff / 60 / 1000)} 分钟前`
  if (diff < 24 * 60 * 60 * 1000) return `${Math.floor(diff / 60 / 60 / 1000)} 小时前`
  if (diff < 7 * 24 * 60 * 60 * 1000) return `${Math.floor(diff / 24 / 60 / 60 / 1000)} 天前`
  return date.toLocaleDateString('zh-CN')
}
</script>