<template>
  <div class="min-h-screen bg-gray-50 p-6">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-primary">{{ stats.total_count || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">全部消息</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-warning">{{ stats.unread_count || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">未读通知</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-primary">{{ stats.system_count || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">系统通知</div>
        </div>
      </div>
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-primary">{{ stats.audit_count || 0 }}</div>
          <div class="text-sm text-gray-500 mt-1">审核通知</div>
        </div>
      </div>
    </div>

    <!-- Tab 切换 -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200">
      <div class="border-b border-gray-200">
        <nav class="flex -mb-px">
          <button
            @click="activeTab = 'my'"
            :class="activeTab === 'my' 
              ? 'border-primary text-primary px-4 py-3 font-medium border-b-2' 
              : 'text-gray-500 hover:text-gray-700 px-4 py-3'"
            class="transition-colors"
          >
            我的通知
          </button>
          <button
            v-if="isAdmin"
            @click="activeTab = 'manage'"
            :class="activeTab === 'manage' 
              ? 'border-primary text-primary px-4 py-3 font-medium border-b-2' 
              : 'text-gray-500 hover:text-gray-700 px-4 py-3'"
            class="transition-colors"
          >
            消息管理
          </button>
        </nav>
      </div>

      <!-- 我的通知 Tab -->
      <div v-if="activeTab === 'my'" class="p-4">
        <!-- 工具栏 -->
        <div class="flex flex-wrap gap-2 mb-4 items-center">
          <div class="flex rounded-md overflow-hidden border border-gray-200">
            <button
              @click="readFilter = 'all'; pagination.page = 1; fetchNotifications()"
              :class="readFilter === 'all' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              全部
            </button>
            <button
              @click="readFilter = 'unread'; pagination.page = 1; fetchNotifications()"
              :class="readFilter === 'unread' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              未读
            </button>
            <button
              @click="readFilter = 'read'; pagination.page = 1; fetchNotifications()"
              :class="readFilter === 'read' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors"
            >
              已读
            </button>
          </div>
          <div class="flex rounded-md overflow-hidden border border-gray-200">
            <button
              @click="typeFilter = 'all'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'all' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              所有类型
            </button>
            <button
              @click="typeFilter = 'system'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'system' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              系统通知
            </button>
            <button
              @click="typeFilter = 'publish'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'publish' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              发布通知
            </button>
            <button
              @click="typeFilter = 'comment_reply'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'comment_reply' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              评论回复
            </button>
            <button
              @click="typeFilter = 'audit_passed'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'audit_passed' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors border-r border-gray-200"
            >
              审核通过
            </button>
            <button
              @click="typeFilter = 'audited_failed'; pagination.page = 1; fetchNotifications()"
              :class="typeFilter === 'audited_failed' ? 'bg-primary text-white px-3 py-1.5 text-sm' : 'bg-white text-gray-700 hover:bg-gray-50 px-3 py-1.5 text-sm'"
              class="transition-colors"
            >
              审核驳回
            </button>
          </div>
          <div class="ml-auto">
            <button
              @click="handleMarkAllRead"
              class="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded-md text-sm transition-colors"
            >
              全部标记已读
            </button>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>

        <!-- 通知列表 -->
        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">时间</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="row in notificationList"
                :key="row.id"
                @click="handleRowClick(row)"
                class="hover:bg-gray-50 cursor-pointer transition-colors"
              >
                <td class="px-4 py-3">
                  <span
                    :class="row.is_read ? 'bg-gray-100 text-gray-600' : 'bg-yellow-100 text-yellow-800'"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ row.is_read ? '已读' : '未读' }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex items-center gap-2">
                    <span
                      v-if="row.message?.priority"
                      :class="getPriorityTagClass(row.message.priority)"
                      class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium"
                    >
                      {{ getPriorityLabel(row.message.priority) }}
                    </span>
                    <span class="text-sm text-gray-900 truncate max-w-[250px]">
                      {{ row.message?.title || '无标题' }}
                    </span>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ MESSAGE_TYPE_MAP[row.message?.message_type] || row.message?.message_type }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ formatTime(row.created_at) }}
                </td>
                <td class="px-4 py-3">
                  <button
                    @click.stop="handleMarkRead(row)"
                    class="text-primary hover:text-primary-dark text-sm font-medium"
                  >
                    {{ row.is_read ? '撤销' : '标记已读' }}
                  </button>
                </td>
              </tr>
              <tr v-if="notificationList.length === 0">
                <td colspan="5" class="px-4 py-8 text-center text-gray-500">
                  暂无通知
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div class="flex justify-end items-center gap-2 mt-4" v-if="pagination.total > 0">
          <span class="text-sm text-gray-600 mr-2">共 {{ pagination.total }} 条</span>
          <button
            @click="pagination.page > 1 && (pagination.page--, fetchNotifications())"
            :disabled="pagination.page === 1"
            :class="pagination.page === 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'"
            class="px-3 py-1 border border-gray-200 rounded text-sm transition-colors"
          >
            上一页
          </button>
          <span v-for="p in visiblePages" :key="p" class="mx-1">
            <button
              @click="pagination.page = p; fetchNotifications()"
              :class="p === pagination.page ? 'bg-primary text-white border-primary' : 'border-gray-200 hover:bg-gray-50'"
              class="w-8 h-8 border rounded text-sm transition-colors"
            >
              {{ p }}
            </button>
          </span>
          <button
            @click="pagination.page < totalPages && (pagination.page++, fetchNotifications())"
            :disabled="pagination.page === totalPages"
            :class="pagination.page === totalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'"
            class="px-3 py-1 border border-gray-200 rounded text-sm transition-colors"
          >
            下一页
          </button>
          <select
            v-model="pagination.pageSize"
            @change="pagination.page = 1; fetchNotifications()"
            class="border border-gray-200 rounded text-sm px-2 py-1"
          >
            <option :value="10">10条/页</option>
            <option :value="20">20条/页</option>
            <option :value="50">50条/页</option>
          </select>
        </div>
      </div>

      <!-- 消息管理 Tab -->
      <div v-if="activeTab === 'manage'" class="p-4">
        <!-- 工具栏 -->
        <div class="flex flex-wrap gap-2 mb-4 items-center">
          <button
            @click="showPublishDialog = true"
            class="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded-md text-sm transition-colors"
          >
            发布消息
          </button>
          <div class="relative">
            <input
              v-model="keyword"
              @input="handleSearch"
              placeholder="搜索消息..."
              class="border border-gray-200 rounded-md px-4 py-2 pl-10 w-64 focus:outline-none focus:ring-2 focus:ring-primary/20"
            />
            <svg class="w-5 h-5 text-gray-400 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
            </svg>
          </div>
        </div>

        <!-- 加载状态 -->
        <div v-if="messageLoading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>

        <!-- 消息列表 -->
        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">标题</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">优先级</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">发送时间</th>
                <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr
                v-for="row in messageList"
                :key="row.id"
                class="hover:bg-gray-50 transition-colors"
              >
                <td class="px-4 py-3 text-sm text-gray-900">
                  {{ row.title }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ MESSAGE_TYPE_MAP[row.message_type] || row.message_type }}
                </td>
                <td class="px-4 py-3">
                  <span
                    :class="getPriorityTagClass(row.priority)"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ getPriorityLabel(row.priority) }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <span
                    :class="getMessageStatusClass(row.status)"
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ MESSAGE_STATUS_MAP[row.status] || row.status }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm text-gray-500">
                  {{ formatTime(row.created_at) }}
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button
                      @click="handleViewDetail(row)"
                      class="text-primary hover:text-primary-dark text-sm font-medium"
                    >
                      查看详情
                    </button>
                    <button
                      @click="handleDeleteMessage(row.id)"
                      class="text-red-600 hover:text-red-800 text-sm font-medium"
                    >
                      删除
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="messageList.length === 0">
                <td colspan="6" class="px-4 py-8 text-center text-gray-500">
                  暂无消息
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- 分页 -->
        <div class="flex justify-end items-center gap-2 mt-4" v-if="messagePagination.total > 0">
          <span class="text-sm text-gray-600 mr-2">共 {{ messagePagination.total }} 条</span>
          <button
            @click="messagePagination.page > 1 && (messagePagination.page--, fetchMessages())"
            :disabled="messagePagination.page === 1"
            :class="messagePagination.page === 1 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'"
            class="px-3 py-1 border border-gray-200 rounded text-sm transition-colors"
          >
            上一页
          </button>
          <span v-for="p in visibleMessagePages" :key="p" class="mx-1">
            <button
              @click="messagePagination.page = p; fetchMessages()"
              :class="p === messagePagination.page ? 'bg-primary text-white border-primary' : 'border-gray-200 hover:bg-gray-50'"
              class="w-8 h-8 border rounded text-sm transition-colors"
            >
              {{ p }}
            </button>
          </span>
          <button
            @click="messagePagination.page < messageTotalPages && (messagePagination.page++, fetchMessages())"
            :disabled="messagePagination.page === messageTotalPages"
            :class="messagePagination.page === messageTotalPages ? 'opacity-50 cursor-not-allowed' : 'hover:bg-gray-100'"
            class="px-3 py-1 border border-gray-200 rounded text-sm transition-colors"
          >
            下一页
          </button>
          <select
            v-model="messagePagination.pageSize"
            @change="messagePagination.page = 1; fetchMessages()"
            class="border border-gray-200 rounded text-sm px-2 py-1"
          >
            <option :value="10">10条/页</option>
            <option :value="20">20条/页</option>
            <option :value="50">50条/页</option>
          </select>
        </div>
      </div>
    </div>

    <!-- 发布消息对话框 -->
    <div
      v-if="showPublishDialog"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @close="showPublishDialog = false"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-lg mx-4">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">发布系统消息</h3>
          <button @click="showPublishDialog = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">标题 <span class="text-red-500">*</span></label>
              <input
                v-model="publishForm.title"
                placeholder="请输入消息标题"
                class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">内容 <span class="text-red-500">*</span></label>
              <textarea
                v-model="publishForm.content"
                rows="6"
                placeholder="请输入消息内容"
                class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20 resize-none"
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">消息类型</label>
              <select
                v-model="publishForm.message_type"
                class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20"
              >
                <option value="system">系统通知</option>
                <option value="publish">发布通知</option>
                <option value="comment_reply">评论回复</option>
                <option value="audit_passed">审核通过</option>
                <option value="audited_failed">审核驳回</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">优先级</label>
              <div class="flex gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.priority" :value="0" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">普通</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.priority" :value="1" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">重要</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.priority" :value="2" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">紧急</span>
                </label>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">目标类型</label>
              <div class="flex gap-4">
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.target_type" value="all" @change="handleTargetTypeChange" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">全部用户</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.target_type" value="user" @change="handleTargetTypeChange" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">指定用户</span>
                </label>
                <label class="flex items-center gap-2 cursor-pointer">
                  <input type="radio" v-model="publishForm.target_type" value="role" @change="handleTargetTypeChange" class="text-primary focus:ring-primary" />
                  <span class="text-sm text-gray-700">指定角色</span>
                </label>
              </div>
            </div>
            <div v-if="publishForm.target_type === 'user'">
              <label class="block text-sm font-medium text-gray-700 mb-1">目标用户</label>
              <select
                v-model="publishForm.target_user_ids"
                multiple
                placeholder="选择目标用户"
                class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20 min-h-[120px]"
              >
                <option
                  v-for="user in userList"
                  :key="user.id"
                  :value="user.id"
                >
                  {{ user.username }}
                </option>
              </select>
              <p class="text-xs text-gray-500 mt-1">按住 Ctrl 键可选择多个用户</p>
            </div>
            <div v-if="publishForm.target_type === 'role'">
              <label class="block text-sm font-medium text-gray-700 mb-1">目标角色</label>
              <select
                v-model="publishForm.target_role_ids"
                multiple
                placeholder="选择目标角色"
                class="w-full border border-gray-200 rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-primary/20 min-h-[120px]"
              >
                <option
                  v-for="role in roleList"
                  :key="role.id"
                  :value="role.id"
                >
                  {{ role.name }}
                </option>
              </select>
              <p class="text-xs text-gray-500 mt-1">按住 Ctrl 键可选择多个角色</p>
            </div>
          </div>
        </div>
        <div class="flex justify-end gap-2 p-6 border-t border-gray-200">
          <button
            @click="showPublishDialog = false"
            class="px-4 py-2 border border-gray-200 rounded-md text-sm hover:bg-gray-50 transition-colors"
          >
            取消
          </button>
          <button
            @click="handlePublish"
            :disabled="publishing"
            :class="publishing ? 'opacity-50 cursor-not-allowed' : ''"
            class="bg-primary hover:bg-primary-dark text-white px-4 py-2 rounded-md text-sm transition-colors disabled:opacity-50"
          >
            <span v-if="publishing" class="flex items-center gap-2">
              <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              发布中...
            </span>
            <span v-else>发布</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 消息详情对话框 -->
    <div
      v-if="showDetailDialog && detailMessage"
      class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
      @close="showDetailDialog = false"
    >
      <div class="bg-white rounded-lg shadow-xl w-full max-w-lg mx-4">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">消息详情</h3>
          <button @click="showDetailDialog = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        <div class="p-6">
          <dl class="space-y-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">标题</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ detailMessage.title }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">类型</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ MESSAGE_TYPE_MAP[detailMessage.message_type] }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">优先级</dt>
              <dd class="mt-1">
                <span
                  :class="getPriorityTagClass(detailMessage.priority)"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                >
                  {{ getPriorityLabel(detailMessage.priority) }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">状态</dt>
              <dd class="mt-1">
                <span
                  :class="getMessageStatusClass(detailMessage.status)"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                >
                  {{ MESSAGE_STATUS_MAP[detailMessage.status] }}
                </span>
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">内容</dt>
              <dd class="mt-1 text-sm text-gray-900 whitespace-pre-wrap leading-relaxed">{{ detailMessage.content }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">发送时间</dt>
              <dd class="mt-1 text-sm text-gray-900">{{ formatTime(detailMessage.created_at) }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import {
  getMyNotifications,
  getMyNotificationStats,
  markAsRead,
  markAllAsRead as apiMarkAllRead,
  deleteMessage,
  publishMessage as apiPublishMessage,
  getMessages,
  getMessageDetail
} from '@/api/notification'
import { MESSAGE_TYPE_MAP, PRIORITY_MAP, MESSAGE_STATUS_MAP } from '@/api/notification'
import { getUserList } from '@/api/user'
import { getRoleList } from '@/api/role'

// 常量定义
const PRIORITY_LABELS: Record<number, string> = { 0: '普通', 1: '重要', 2: '紧急' }
const STATUS_LABELS: Record<number, string> = { 1: '已发布', 0: '草稿', 2: '已撤回' }

// 响应式数据
const activeTab = ref('my')
const readFilter = ref<'all' | 'unread' | 'read'>('all')
const typeFilter = ref<string>('all')
const keyword = ref('')
const loading = ref(false)
const messageLoading = ref(false)
const stats = ref({
  total_count: 0,
  unread_count: 0,
  read_count: 0,
  system_count: 0,
  publish_count: 0,
  comment_count: 0,
  audit_count: 0
})

// 通知列表
const notificationList = ref<any[]>([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 消息管理
const messageList = ref<any[]>([])
const messagePagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 发布消息对话框
const showPublishDialog = ref(false)
const publishing = ref(false)
const publishForm = reactive({
  title: '',
  content: '',
  message_type: 'system',
  priority: 0,
  target_type: 'all',
  target_user_ids: [] as number[],
  target_role_ids: [] as number[]
})

// 消息详情
const showDetailDialog = ref(false)
const detailMessage = ref<any>(null)

// 用户列表（用于选择目标用户）
const userList = ref<any[]>([])
const roleList = ref<any[]>([])
const isAdmin = ref(false) // TODO: 根据用户角色判断

// 分页计算
const totalPages = computed(() => Math.ceil(pagination.total / pagination.pageSize))
const visiblePages = computed(() => {
  const pages: number[] = []
  const total = totalPages.value
  const current = pagination.page
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

const messageTotalPages = computed(() => Math.ceil(messagePagination.total / messagePagination.pageSize))
const visibleMessagePages = computed(() => {
  const pages: number[] = []
  const total = messageTotalPages.value
  const current = messagePagination.page
  const start = Math.max(1, current - 2)
  const end = Math.min(total, current + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// 方法
const getPriorityLabel = (priority: number) => PRIORITY_LABELS[priority] || '未知'
const getPriorityTagClass = (priority: number) => {
  const classes: Record<number, string> = {
    0: 'bg-gray-100 text-gray-800',
    1: 'bg-yellow-100 text-yellow-800',
    2: 'bg-red-100 text-red-800'
  }
  return classes[priority] || 'bg-gray-100 text-gray-800'
}

const getMessageStatusClass = (status: number) => {
  const classes: Record<number, string> = {
    1: 'bg-green-100 text-green-800',
    0: 'bg-gray-100 text-gray-800',
    2: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatTime = (time: string | undefined) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 获取通知统计
const fetchStats = async () => {
  try {
    const res = await getMyNotificationStats()
    // 后端返回蛇形命名 (snake_case)，映射到前端变量
    if (res.data) {
      stats.value = {
        total_count: res.data.total_count || 0,
        unread_count: res.data.unread_count || 0,
        read_count: res.data.read_count || 0,
        system_count: res.data.system_count || 0,
        publish_count: res.data.publish_count || 0,
        comment_count: res.data.comment_count || 0,
        audit_count: res.data.audit_count || 0
      }
    }
  } catch (error) {
    console.error('获取统计失败:', error)
  }
}

// 获取通知列表
const fetchNotifications = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize,
      user_id: 1 // TODO: 从 store 获取当前用户 ID
    }
    if (readFilter.value !== 'all') {
      params.is_read = readFilter.value === 'unread' ? 0 : 1
    }
    if (typeFilter.value !== 'all') {
      params.message_type = typeFilter.value
    }

    const res = await getMyNotifications(params)
    notificationList.value = res.data.data?.list || []
    pagination.total = res.data.data?.total || 0
  } catch (error) {
    console.error('获取通知失败:', error)
    alert('获取通知失败')
  } finally {
    loading.value = false
  }
}

// 获取消息列表（管理员）
const fetchMessages = async () => {
  messageLoading.value = true
  try {
    const params: any = {
      page: messagePagination.page,
      page_size: messagePagination.pageSize
    }
    if (keyword.value) {
      params.keyword = keyword.value
    }

    const res = await getMessages(params)
    messageList.value = res.data.data?.list || []
    messagePagination.total = res.data.data?.total || 0
  } catch (error) {
    console.error('获取消息失败:', error)
    alert('获取消息失败')
  } finally {
    messageLoading.value = false
  }
}

// Tab 切换
const handleTabChange = (paneName: string) => {
  if (paneName === 'my') {
    fetchStats()
    fetchNotifications()
  } else if (paneName === 'manage') {
    fetchMessages()
  }
}

// 筛选变化
const handleReadFilterChange = () => {
  pagination.page = 1
  fetchNotifications()
}

const handleTypeFilterChange = () => {
  pagination.page = 1
  fetchNotifications()
}

const handleSearch = () => {
  messagePagination.page = 1
  fetchMessages()
}

// 标记已读
const handleMarkRead = async (row: any) => {
  const newIsRead = row.is_read ? 0 : 1
  try {
    await markAsRead(row.id, newIsRead)
    row.is_read = newIsRead
    alert('操作成功')
    fetchStats()
  } catch (error) {
    console.error('标记失败:', error)
    alert('操作失败')
  }
}

// 全部标记已读
const handleMarkAllRead = async () => {
  if (!confirm('确定要将所有通知标记为已读吗？')) {
    return
  }
  try {
    await apiMarkAllRead()
    alert('全部标记成功')
    fetchStats()
    fetchNotifications()
  } catch (error) {
    console.error('全部标记失败:', error)
    alert('操作失败')
  }
}

// 行点击
const handleRowClick = (row: any) => {
  if (!row.is_read) {
    handleMarkRead(row)
  }
}

// 发布消息
const handlePublishMessage = () => {
  showPublishDialog.value = true
}

const handleTargetTypeChange = async () => {
  if (publishForm.target_type === 'user') {
    try {
      const res = await getUserList({ page: 1, page_size: 1000 })
      userList.value = res.data?.list || []
    } catch (error) {
      console.error('加载用户列表失败:', error)
    }
  } else if (publishForm.target_type === 'role') {
    try {
      const res = await getRoleList({ page: 1, page_size: 1000 })
      roleList.value = res.data?.list || []
    } catch (error) {
      console.error('加载角色列表失败:', error)
    }
  }
}

const handlePublish = async () => {
  if (!publishForm.title || !publishForm.content) {
    alert('请填写标题和内容')
    return
  }

  // 根据目标类型设置相应的ID数组
  const payload: any = { ...publishForm }
  if (publishForm.target_type === 'all') {
    payload.target_user_ids = []
    payload.target_role_ids = []
  } else if (publishForm.target_type === 'user') {
    // 获取选中的用户ID
    const selectEl = document.querySelector('select[multiple]') as HTMLSelectElement
    if (selectEl) {
      payload.target_user_ids = Array.from(selectEl.selectedOptions).map((opt: HTMLOptionElement) => parseInt(opt.value))
    }
    payload.target_role_ids = []
  } else if (publishForm.target_type === 'role') {
    // 获取选中的角色ID
    const selectEl = document.querySelector('select[multiple]') as HTMLSelectElement
    if (selectEl) {
      payload.target_role_ids = Array.from(selectEl.selectedOptions).map((opt: HTMLOptionElement) => parseInt(opt.value))
    }
    payload.target_user_ids = []
  }

  publishing.value = true
  try {
    await apiPublishMessage(payload)
    alert('消息发布成功')
    showPublishDialog.value = false
    // 重置表单
    publishForm.title = ''
    publishForm.content = ''
    publishForm.message_type = 'system'
    publishForm.priority = 0
    publishForm.target_type = 'all'
    publishForm.target_user_ids = []
    publishForm.target_role_ids = []
  } catch (error) {
    console.error('发布消息失败:', error)
    alert('发布消息失败')
  } finally {
    publishing.value = false
  }
}

// 查看详情
const handleViewDetail = async (row: any) => {
  try {
    const res = await getMessageDetail(row.id)
    detailMessage.value = res.data
    showDetailDialog.value = true
  } catch (error) {
    console.error('获取详情失败:', error)
    alert('获取详情失败')
  }
}

// 删除消息
const handleDeleteMessage = async (id: number) => {
  if (!confirm('确定要删除这条消息吗？')) {
    return
  }
  try {
    await deleteMessage(id)
    alert('删除成功')
    fetchMessages()
  } catch (error) {
    console.error('删除失败:', error)
    alert('删除失败')
  }
}

// 生命周期
onMounted(() => {
  fetchStats()
  fetchNotifications()
})
</script>

<style scoped>
:root {
  --color-primary: #3b82f6;
  --color-primary-dark: #2563eb;
  --color-warning: #f59e0b;
}

.text-primary {
  color: var(--color-primary);
}

.text-primary-dark {
  color: var(--color-primary-dark);
}

.bg-primary {
  background-color: var(--color-primary);
}

.bg-primary-dark {
  background-color: var(--color-primary-dark);
}

.hover\\:bg-primary-dark:hover {
  background-color: var(--color-primary-dark);
}

.text-warning {
  color: var(--color-warning);
}

.ring-primary\/20 {
  --tw-ring-color: rgb(59 130 246 / 0.2);
}
</style>