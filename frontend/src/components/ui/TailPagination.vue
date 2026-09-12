
<template>
  <div v-if="total > pageSize" :class="paginationClasses">
    <!-- Info -->
    <div v-if="showInfo" class="text-sm text-gray-500">
      显示 {{ start }} - {{ end }} / 共 {{ total }} 条
    </div>

    <!-- Pages -->
    <div class="flex items-center gap-1">
      <!-- Prev Button -->
      <button
        :disabled="currentPage === 1"
        @click="handlePrev"
        class="px-3 py-1.5 rounded-md text-sm border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- First Page -->
      <button
        v-if="showFirstLast && currentPage > 3"
        @click="goTo(1)"
        class="w-8 h-8 rounded-md text-sm border border-gray-300 hover:bg-gray-50 transition-colors"
      >
        1
      </button>

      <!-- Ellipsis after first page -->
      <span
        v-if="showFirstLast && currentPage > 4"
        class="px-2 text-gray-400"
      >...</span>

      <!-- Page Numbers -->
      <button
        v-for="page in visiblePages"
        :key="page"
        @click="goTo(page)"
        :class="page === currentPage ? activeClasses : 'border-gray-300 hover:bg-gray-50'"
        class="w-8 h-8 rounded-md text-sm border transition-colors"
      >
        {{ page }}
      </button>

      <!-- Ellipsis before last page -->
      <span
        v-if="showFirstLast && currentPage < totalPages - 3"
        class="px-2 text-gray-400"
      >...</span>

      <!-- Last Page -->
      <button
        v-if="showFirstLast && currentPage < totalPages - 2"
        @click="goTo(totalPages)"
        class="w-8 h-8 rounded-md text-sm border border-gray-300 hover:bg-gray-50 transition-colors"
      >
        {{ totalPages }}
      </button>

      <!-- Next Button -->
      <button
        :disabled="currentPage === totalPages"
        @click="handleNext"
        class="px-3 py-1.5 rounded-md text-sm border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50 transition-colors"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <!-- Size Selector -->
    <div v-if="showSizeSelector && sizeOptions" class="flex items-center gap-2">
      <select
        :value="pageSize"
        @change="handleSizeChange($event)"
        class="border border-gray-300 rounded-md text-sm px-2 py-1 focus:ring-indigo-500 focus:border-indigo-500"
      >
        <option v-for="size in sizeOptions" :key="size" :value="size">
          {{ size }} 条/页
        </option>
      </select>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface TailPaginationProps {
  currentPage: number
  pageSize: number
  total: number
  maxVisible?: number
  showInfo?: boolean
  showSizeSelector?: boolean
  sizeOptions?: number[]
  showFirstLast?: boolean
  compact?: boolean
}

const props = withDefaults(defineProps<TailPaginationProps>(), {
  currentPage: 1,
  pageSize: 10,
  total: 0,
  maxVisible: 7,
  showInfo: true,
  showSizeSelector: false,
  sizeOptions: () => [10, 20, 50, 100],
  showFirstLast: true,
  compact: false,
})

const emit = defineEmits<{
  'update:currentPage': [value: number]
  'change': [page: number]
  'size-change': [size: number]
}>()

const totalPages = computed(() => Math.ceil(props.total / props.pageSize) || 1)

const start = computed(() => {
  if (props.total === 0) return 0
  return (props.currentPage - 1) * props.pageSize + 1
})

const end = computed(() => {
  if (props.total === 0) return 0
  return Math.min(props.currentPage * props.pageSize, props.total)
})

const visiblePages = computed(() => {
  const maxVisible = props.maxVisible || 7
  const total = totalPages.value
  const current = props.currentPage
  
  if (total <= maxVisible) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  
  const pages: number[] = []
  const halfVisible = Math.floor(maxVisible / 2)
  
  let start = current - halfVisible
  let end = current + halfVisible
  
  if (start < 1) {
    start = 1
    end = Math.min(maxVisible, total)
  }
  
  if (end > total) {
    end = total
    start = Math.max(1, total - maxVisible + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

const activeClasses = 'bg-indigo-600 text-white border-indigo-600'

const paginationClasses = 'flex items-center justify-between flex-wrap gap-4 py-4 border-t border-gray-200'

function goTo(page: number) {
  if (page < 1 || page > totalPages.value || page === props.currentPage) return
  emit('update:currentPage', page)
  emit('change', page)
}

function handlePrev() {
  if (props.currentPage > 1) {
    goTo(props.currentPage - 1)
  }
}

function handleNext() {
  if (props.currentPage < totalPages.value) {
    goTo(props.currentPage + 1)
  }
}

function handleSizeChange(event: Event) {
  const target = event.target as HTMLSelectElement
  const newSize = parseInt(target.value, 10)
  emit('size-change', newSize)
}
</script>