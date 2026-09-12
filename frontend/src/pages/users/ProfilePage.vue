<template>
  <div class="page-container max-w-4xl">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center gap-2 mb-2">
        <button class="text-gray-500 hover:text-gray-700" @click="$router.back()">
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>
        <div>
          <h1 class="text-2xl font-bold text-gray-900 mb-1">个人信息</h1>
          <p class="text-sm text-gray-500">查看和编辑您的个人资料信息</p>
        </div>
      </div>
    </div>

    <!-- Profile Card -->
    <TailCard class="p-6">
      <div class="flex items-center gap-6 mb-8 pb-8 border-b border-gray-200">
        <div class="w-20 h-20 bg-primary/10 text-primary rounded-full flex items-center justify-center text-2xl font-bold">
          {{ avatarChar }}
        </div>
        <div class="flex-1">
          <h2 class="text-xl font-semibold text-gray-900">{{ displayName }}</h2>
          <p class="text-sm text-gray-500">{{ authStore.userInfo?.email || '未设置邮箱' }}</p>
          <div class="flex gap-2 mt-3">
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              {{ authStore.userInfo?.role_name || '未知角色' }}
            </span>
          </div>
        </div>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Username -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="ta-label">用户名 <span class="text-red-500">*</span></label>
            <input v-model="formData.username" type="text" class="ta-input w-full" placeholder="请输入用户名" required />
          </div>
          <div>
            <label class="ta-label">邮箱 <span class="text-red-500">*</span></label>
            <input v-model="formData.email" type="email" class="ta-input w-full" placeholder="请输入邮箱" required />
          </div>
        </div>

        <!-- Phone and Real Name -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
          <div>
            <label class="ta-label">手机号</label>
            <input v-model="formData.phone" type="tel" class="ta-input w-full" placeholder="请输入手机号" />
          </div>
          <div>
            <label class="ta-label">真实姓名</label>
            <input v-model="formData.realName" type="text" class="ta-input w-full" placeholder="请输入真实姓名" />
          </div>
        </div>

        <!-- Bio -->
        <div>
          <label class="ta-label">个人简介</label>
          <textarea v-model="formData.bio" class="ta-input w-full" rows="4" placeholder="请输入个人简介"></textarea>
        </div>

        <!-- Password Section -->
        <div class="pt-4 border-t border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">修改密码</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label class="ta-label">当前密码</label>
              <input v-model="passwordForm.currentPassword" type="password" class="ta-input w-full" placeholder="请输入当前密码" />
            </div>
            <div>
              <label class="ta-label">新密码</label>
              <input v-model="passwordForm.newPassword" type="password" class="ta-input w-full" placeholder="请输入新密码" />
            </div>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mt-5">
            <div>
              <label class="ta-label">确认新密码</label>
              <input v-model="passwordForm.confirmPassword" type="password" class="ta-input w-full" placeholder="请再次输入新密码" />
            </div>
          </div>
        </div>

        <!-- Submit Buttons -->
        <div class="flex gap-3 pt-4 border-t border-gray-200">
          <button type="submit" class="ta-btn ta-btn-primary">
            保存修改
          </button>
          <button type="button" class="ta-btn ta-btn-outline" @click="resetForm">重置</button>
        </div>
      </form>
    </TailCard>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { message } from '@/utils/message'
import TailCard from '@/components/ui/TailCard.vue'

const router = useRouter()
const authStore = useAuthStore()

const formData = reactive({
  username: '',
  email: '',
  phone: '',
  realName: '',
  bio: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const displayName = computed(() => {
  return authStore.userInfo?.username || '用户'
})

const avatarChar = computed(() => {
  return displayName.value.charAt(0).toUpperCase()
})

const loadUserInfo = async () => {
  try {
    const userInfo = authStore.userInfo
    if (userInfo) {
      formData.username = userInfo.username || ''
      formData.email = userInfo.email || ''
      formData.phone = userInfo.phone || ''
      formData.realName = userInfo.real_name || ''
      formData.bio = userInfo.bio || ''
    }
  } catch (error) {
    console.error('加载用户信息失败:', error)
  }
}

const handleSubmit = async () => {
  if (!formData.username) {
    message.error('请输入用户名')
    return
  }
  if (!formData.email) {
    message.error('请输入邮箱')
    return
  }

  try {
    // TODO: 调用更新用户信息的 API
    message.success('个人信息保存成功')
    // 清空密码表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
  } catch (error) {
    message.error('保存失败')
  }
}

const resetForm = () => {
  loadUserInfo()
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
}

onMounted(() => {
  loadUserInfo()
})
</script>