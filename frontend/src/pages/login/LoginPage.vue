<template>
  <div class="min-h-screen flex auth-bg">
    <!-- Left: Branding -->
    <div class="hidden lg:flex lg:w-1/2 bg-gradient-to-br from-primary/90 to-indigo-800 relative overflow-hidden">
      <div class="absolute inset-0 bg-[url('data:image/svg+xml,%3Csvg%20width%3D%2260%22%20height%3D%2260%22%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%3E%3Cg%20fill%3D%22none%22%20fill-rule%3D%22evenodd%22%3E%3Cg%20fill%3D%22%23ffffff%22%20fill-opacity%3D%220.05%22%3E%3Cpath%20d%3D%22M36%2034v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6%2034v-4H4v4H0v2h4v4h2v-4h4v-2H6z%22%2F%3E%3C%2Fg%3E%3C%2Fsvg%3E')] opacity-50"></div>
      <div class="relative z-10 flex flex-col items-center justify-center w-full p-12 text-white">
        <div class="flex items-center gap-4 mb-8">
          <div class="w-16 h-16 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
            <svg class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <span class="text-3xl font-bold">GoCMS</span>
        </div>
        <h2 class="text-3xl font-bold mb-4">现代内容管理系统</h2>
        <p class="text-lg text-white/80 mb-12">强大、灵活、易用的企业级 CMS 平台</p>
        <div class="space-y-6 max-w-sm w-full">
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <span class="text-white/90">高性能微服务架构</span>
          </div>
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
            </div>
            <span class="text-white/90">RBAC 权限管理系统</span>
          </div>
          <div class="flex items-center gap-4">
            <div class="w-10 h-10 bg-white/20 rounded-lg flex items-center justify-center flex-shrink-0">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
              </svg>
            </div>
            <span class="text-white/90">多租户数据隔离</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Right: Login Form -->
    <div class="flex-1 flex items-center justify-center p-8">
      <div class="w-full max-w-md">
        <!-- Mobile Logo -->
        <div class="text-center mb-8 lg:hidden">
          <div class="w-14 h-14 bg-primary/10 text-primary rounded-xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h1 class="text-2xl font-bold text-gray-900">GoCMS</h1>
        </div>

        <div class="ta-card p-8">
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-gray-900 mb-2">欢迎回来</h2>
            <p class="text-gray-500 text-sm">请输入您的账号密码登录管理后台</p>
          </div>

          <!-- 登录表单 -->
          <form @submit.prevent="handleLogin" class="space-y-5">
            <!-- 用户名 -->
            <div>
              <label class="ta-label">用户名</label>
              <div class="relative">
                <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
                <input v-model="loginForm.username" type="text" placeholder="请输入用户名" class="ta-input pl-10" :class="{ 'border-red-500 focus:ring-red-200 focus:border-red-500': errors.username }" />
              </div>
              <p v-if="errors.username" class="mt-1.5 text-xs text-red-600">{{ errors.username }}</p>
            </div>

            <!-- 密码 -->
            <div>
              <label class="ta-label">密码</label>
              <div class="relative">
                <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                </svg>
                <input v-model="loginForm.password" :type="showPassword ? 'text' : 'password'" placeholder="请输入密码" class="ta-input pl-10 pr-10" :class="{ 'border-red-500 focus:ring-red-200 focus:border-red-500': errors.password }" />
                <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600">
                  <svg v-if="!showPassword" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </button>
              </div>
              <p v-if="errors.password" class="mt-1.5 text-xs text-red-600">{{ errors.password }}</p>
            </div>

            <!-- 记住我 & 忘记密码 -->
            <div class="flex items-center justify-between">
              <label class="flex items-center gap-2 cursor-pointer">
                <input v-model="loginForm.remember" type="checkbox" class="w-4 h-4 text-primary border-gray-300 rounded focus:ring-primary" />
                <span class="text-sm text-gray-600">记住我</span>
              </label>
              <a href="#" class="text-sm text-primary hover:text-primary-700 font-medium">忘记密码？</a>
            </div>

            <!-- 登录按钮 -->
            <button type="submit" class="ta-btn ta-btn-primary w-full h-11 text-base" :disabled="loading">
              <svg v-if="loading" class="animate-spin w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? '登录中...' : '登 录' }}
            </button>
          </form>

          <!-- 管理员提示 -->
          <p class="mt-6 text-center text-sm text-gray-500">
            管理后台账号请联系系统管理员分配
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const loading = ref(false)
const showPassword = ref(false)
const errors = reactive<Record<string, string>>({})

const loginForm = reactive({
  username: '',
  password: '',
  remember: false,
})

const validateForm = (): boolean => {
  Object.keys(errors).forEach(key => delete errors[key])
  
  if (!loginForm.username.trim()) {
    errors.username = '请输入用户名'
    return false
  }
  if (!loginForm.password) {
    errors.password = '请输入密码'
    return false
  }
  if (loginForm.password.length < 6) {
    errors.password = '密码长度不能少于6位'
    return false
  }
  return true
}

const handleLogin = async () => {
  if (!validateForm()) return
  
  loading.value = true
  try {
    const axios = await import('axios').then(m => m.default)
    const res = await axios.post('/api/auth/login', {
      username: loginForm.username,
      password: loginForm.password,
    })
    
    // 修复：后端 response.Success 返回的结构是 {code: 0, message: "success", data: {...}}
    // res.data 对应整个响应体，实际数据在 res.data.data 中
    const result = res.data
    if (result.code !== 0) {
      errors.username = result.message || '登录失败'
      return
    }
    
    const data = result.data
    authStore.setAuthData(data.token, {
      id: data.user_id,
      username: data.username,
      nickname: data.nickname,
      roles: data.roles || [],
    })
    
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (err: any) {
    errors.username = err.response?.data?.message || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}
</script>