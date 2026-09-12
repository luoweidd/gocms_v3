import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

interface LoadingState {
  [key: string]: boolean
}

export const useLoadingStore = defineStore('loading', () => {
  const loadingState = ref<LoadingState>({})

  function setLoading(key: string, value: boolean) {
    loadingState.value[key] = value
  }

  function setBulkLoading(keys: string[], value: boolean) {
    keys.forEach((key) => {
      loadingState.value[key] = value
    })
  }

  function clearLoading(key: string) {
    delete loadingState.value[key]
  }

  function clearAllLoading() {
    loadingState.value = {}
  }

  return {
    loadingState,
    setLoading,
    setBulkLoading,
    clearLoading,
    clearAllLoading
  }
})

/**
 * 获取全局 loading 状态的方法
 */
export function useGlobalLoading() {
  const loadingStore = useLoadingStore()

  return {
    /**
     * 显示全局加载状态
     */
    show: () => {
      loadingStore.setLoading('global', true)
    },
    /**
     * 隐藏全局加载状态
     */
    hide: () => {
      loadingStore.setLoading('global', false)
    }
  }
}

/**
 * 获取指定 key 的加载状态
 */
export function useLoading(key: string) {
  const loadingStore = useLoadingStore()
  return computed(() => loadingStore.loadingState[key] || false)
}