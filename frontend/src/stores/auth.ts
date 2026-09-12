import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { UserInfo } from '@/api'

interface AuthState {
  token: string | null
  userInfo: UserInfo | null
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<UserInfo | null>(
    (() => {
      const stored = localStorage.getItem('userInfo')
      return stored ? JSON.parse(stored) : null
    })()
  )

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)
  const currentUsername = computed(() => userInfo.value?.username || '')
  const userRoles = computed(() => userInfo.value?.roles || [])
  const userId = computed(() => userInfo.value?.id || 0)
  const userAvatar = computed(() => userInfo.value?.avatar || '')

  // Actions
  function setAuthData(newToken: string, newUserInfo: UserInfo) {
    token.value = newToken
    userInfo.value = newUserInfo
    localStorage.setItem('token', newToken)
    localStorage.setItem('userInfo', JSON.stringify(newUserInfo))
  }

  function updateUserInfo(updatedInfo: Partial<UserInfo>) {
    if (userInfo.value) {
      userInfo.value = { ...userInfo.value, ...updatedInfo }
      localStorage.setItem('userInfo', JSON.stringify(userInfo.value))
    }
  }

  async function logout() {
    try {
      // 尝试调用后端登出接口（可选）
      const axios = (await import('axios')).default
      await axios.post('/api/auth/logout')
        .catch(() => {}) // 忽略错误，仍然清除本地状态
    } catch {
      // 静默处理错误
    } finally {
      token.value = null
      userInfo.value = null
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
    }
  }

  function getToken(): string | null {
    if (!token.value) {
      token.value = localStorage.getItem('token')
    }
    return token.value
  }

  // Token 过期检查（默认 2 小时）
  const tokenExpiryTime = 2 * 60 * 60 * 1000
  function isTokenExpired(expiryTimestamp: number): boolean {
    return Date.now() > expiryTimestamp
  }

  return {
    // State
    token,
    userInfo,
    // Getters
    isAuthenticated,
    isLoggedIn,
    currentUsername,
    userRoles,
    userId,
    userAvatar,
    // Actions
    setAuthData,
    updateUserInfo,
    logout,
    getToken,
    isTokenExpired
  }
})