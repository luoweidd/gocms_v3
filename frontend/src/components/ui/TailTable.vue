<template>
  <div :class="tableContainerClasses">
    <!-- Table Header -->
    <div v-if="$slots.header" class="mb-4">
      <slot name="header" />
    </div>

    <!-- Table -->
    <div :class="{ 'overflow-x-auto': horizontalScroll }">
      <table :class="tableClasses">
        <thead :class="theadClasses">
          <tr>
            <!-- Selection Column -->
            <th v-if="selectable" class="px-4 py-3 text-center" :class="thClasses">
              <input
                type="checkbox"
                :checked="isAllSelected"
                @change="toggleAll($event)"
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              />
            </th>

            <!-- Column Headers -->
            <th
              v-for="col in columns"
              :key="col.key"
              :class="thClasses"
              class="text-left py-3 px-4 font-semibold text-sm text-gray-800 cursor-pointer select-none whitespace-nowrap"
              @click="col.sortable ? handleSort(col.key) : null"
            >
              <div class="flex items-center gap-1">
                <span>{{ col.label || col.title }}</span>
                <span v-if="col.sortable" class="inline-flex flex-col">
                  <svg 
                    class="w-3 h-3 -mb-px text-gray-400" 
                    fill="currentColor" viewBox="0 0 24 24"
                  >
                    <path d="M12 7l-7 7h14L12 7z" />
                  </svg>
                  <svg 
                    class="w-3 h-3 mt-px text-gray-400" 
                    fill="currentColor" viewBox="0 0 24 24"
                  >
                    <path d="M12 17l7-7H5l7 7z" />
                  </svg>
                </span>
              </div>
            </th>

            <!-- Actions Column -->
            <th v-if="showActions" class="text-right py-3 px-4 font-semibold text-sm text-gray-800 whitespace-nowrap">
              操作
            </th>
          </tr>
        </thead>

        <!-- Loading State -->
        <tbody v-if="loading">
          <tr>
            <td :colspan="totalColumns" class="px-4 py-8 text-center">
              <div class="flex flex-col items-center justify-center text-gray-400">
                <svg class="animate-spin h-8 w-8 mb-2" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <span class="text-sm">{{ loadingText }}</span>
              </div>
            </td>
          </tr>
        </tbody>

        <!-- Empty State -->
        <tbody v-else-if="isEmpty">
          <tr>
            <td :colspan="totalColumns" class="px-4 py-8 text-center">
              <div class="flex flex-col items-center justify-center text-gray-400">
                <svg class="h-12 w-12 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                </svg>
                <p class="text-gray-500 text-sm">{{ emptyText }}</p>
                <slot name="empty-action" />
              </div>
            </td>
          </tr>
        </tbody>

        <!-- Data Rows -->
        <tbody v-else :class="tbodyClasses">
          <tr
            v-for="(row, rowIndex) in data"
            :key="getRowKey(row, rowIndex)"
            :class="[
              rowClasses, 
              striped && rowIndex % 2 === 0 ? '' : '',
              striped && rowIndex % 2 !== 0 ? 'bg-gray-50' : '',
              selectable ? 'cursor-pointer' : '',
              props.onClickRow ? 'hover:bg-indigo-50' : '',
            ].filter(Boolean).join(' ')"
            @click="handleRowClick(row, $event)"
          >
            <!-- Selection Cell -->
            <td v-if="selectable" class="px-4 py-3 text-center" :class="tdClasses">
              <input
                type="checkbox"
                :checked="isSelected(row)"
                @change="handleSelect(row, $event)"
                @click.stop
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              />
            </td>

            <!-- Data Cells -->
            <td
              v-for="col in columns"
              :key="col.key"
              :class="tdClasses"
              class="text-sm text-gray-700"
            >
              <slot :name="`col-${col.key}`" :row="row" :value="getCellValue(row, col)">
                {{ formatCell(col, getCellValue(row, col)) }}
              </slot>
            </td>

            <!-- Actions Cell -->
            <td v-if="showActions" class="text-right py-3 px-4 whitespace-nowrap text-sm" :class="tdClasses">
              <div class="flex items-center justify-end gap-2">
                <slot name="actions" :row="row" />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div v-if="pagination" :class="paginationWrapperClasses" class="flex items-center justify-between mt-4">
      <div class="text-sm text-gray-500">
        显示 {{ paginationStart }} - {{ paginationEnd }} / 共 {{ pagination?.total || 0 }} 条
      </div>
      <div class="flex items-center gap-1">
        <button
          :disabled="currentPage === 1"
          @click="goToPage(currentPage - 1)"
          class="px-3 py-1.5 rounded-md text-sm border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          上一页
        </button>
        
        <button
          v-for="page in visiblePages"
          :key="page"
          @click="goToPage(page)"
          :class="page === currentPage ? 'bg-indigo-600 text-white border-indigo-600' : 'border-gray-300 hover:bg-gray-50'"
          class="px-3 py-1.5 rounded-md text-sm border"
        >
          {{ page }}
        </button>

        <button
          :disabled="currentPage === totalPages"
          @click="goToPage(currentPage + 1)"
          class="px-3 py-1.5 rounded-md text-sm border border-gray-300 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

export interface TableColumn {
  key: string
  label?: string
  title?: string
  width?: string
  sortable?: boolean
  formatter?: (value: any, row: Record<string, any>) => string
}

export interface TablePagination {
  currentPage: number
  pageSize: number
  total: number
}

export interface TailTableProps {
  data: Record<string, any>[]
  columns: TableColumn[]
  loading?: boolean
  loadingText?: string
  emptyText?: string
  selectable?: boolean
  striped?: boolean
  showActions?: boolean
  horizontalScroll?: boolean
  pagination?: TablePagination
  sortColumn?: string
  sortOrder?: 'asc' | 'desc' | null
  onClickRow?: (row: Record<string, any>, event: Event) => void
}

const props = withDefaults(defineProps<TailTableProps>(), {
  data: () => [],
  loading: false,
  loadingText: '加载中...',
  emptyText: '暂无数据',
  selectable: false,
  striped: true,
  showActions: true,
  horizontalScroll: true,
  pagination: undefined,
  sortColumn: '',
  sortOrder: null,
  onClickRow: undefined,
})

const emit = defineEmits<{
  'update:sortColumn': [value: string]
  'update:sortOrder': [value: 'asc' | 'desc' | null]
  'sort': [key: string, order: 'asc' | 'desc']
  'select': [rows: string[]]
  'page-change': [page: number]
}>()

const selectedRows = ref<Set<string>>(new Set())
const currentPage = computed(() => props.pagination?.currentPage || 1)
const totalPages = computed(() => {
  if (!props.pagination) return 1
  return Math.ceil(props.pagination.total / props.pagination.pageSize) || 1
})
const paginationStart = computed(() => {
  if (!props.pagination) return 0
  return (currentPage.value - 1) * props.pagination.pageSize + 1
})
const paginationEnd = computed(() => {
  if (!props.pagination) return 0
  return Math.min(currentPage.value * props.pagination.pageSize, props.pagination.total)
})

const totalColumns = computed(() => {
  let count = props.columns.length
  if (props.selectable) count++
  if (props.showActions) count++
  return count
})

const tableContainerClasses = computed(() => {
  return ['bg-white rounded-lg shadow-sm', props.horizontalScroll ? 'overflow-x-auto' : ''].filter(Boolean).join(' ')
})

const tableClasses = computed(() => {
  return 'min-w-full divide-y divide-gray-200'
})

const theadClasses = computed(() => {
  return 'bg-gray-50'
})

const thClasses = computed(() => {
  return 'text-left py-3 px-4 font-semibold text-sm text-gray-800 whitespace-nowrap'
})

const tbodyClasses = computed(() => {
  return 'divide-y divide-gray-200 bg-white'
})

const rowClasses = 'transition-colors duration-150'

const tdClasses = 'py-3 px-4 text-sm text-gray-700'

const paginationWrapperClasses = 'flex items-center justify-between pt-4 border-t border-gray-200'

const isEmpty = computed(() => !props.loading && props.data.length === 0)

const isAllSelected = computed(() => {
  if (!props.selectable || props.data.length === 0) return false
  return selectedRows.value.size === props.data.length
})

const visiblePages = computed(() => {
  const total = totalPages.value
  const current = currentPage.value
  const pages: number[] = []
  
  for (let i = Math.max(1, current - 2); i <= Math.min(total, current + 2); i++) {
    pages.push(i)
  }
  
  return pages
})

function getRowKey(row: Record<string, any>, index: number): string {
  if ('id' in row) return String((row as any).id)
  return `row-${index}`
}

function getCellValue(row: Record<string, any>, col: TableColumn): any {
  const keys = col.key.split('.')
  let value: any = row
  
  for (const key of keys) {
    if (value == null) return ''
    value = value[key]
  }
  
  return value ?? ''
}

function formatCell(col: TableColumn, value: any): string {
  if (col.formatter) return col.formatter(value, {} as any)
  if (typeof value === 'boolean') return value ? '是' : '否'
  if (value === null || value === undefined) return '-'
  return String(value)
}

function handleSort(key: string) {
  const currentOrder = props.sortOrder
  let newOrder: 'asc' | 'desc' | null
  
  if (currentOrder === null || currentOrder === 'desc') {
    newOrder = 'asc'
  } else {
    newOrder = 'desc'
  }
  
  emit('update:sortColumn', key)
  emit('update:sortOrder', newOrder)
  emit('sort', key, newOrder)
}

function isSelected(row: Record<string, any>): boolean {
  if (!('id' in row)) return false
  const id = String((row as any).id)
  return selectedRows.value.has(id)
}

function handleSelect(row: Record<string, any>, event: Event) {
  const target = event.target as HTMLInputElement
  const id = 'id' in row ? String((row as any).id) : null
  
  if (id !== null) {
    if (target.checked) {
      selectedRows.value.add(id)
    } else {
      selectedRows.value.delete(id)
    }
  }
  
  emit('select', Array.from(selectedRows.value))
}

function toggleAll(event: Event) {
  const target = event.target as HTMLInputElement
  
  if (target.checked) {
    props.data.forEach((row, index) => {
      selectedRows.value.add(getRowKey(row, index))
    })
  } else {
    selectedRows.value.clear()
  }
  
  emit('select', Array.from(selectedRows.value))
}

function handleRowClick(row: Record<string, any>, event: Event) {
  if (props.onClickRow) {
    props.onClickRow(row, event)
  } else if (props.selectable) {
    handleSelect(row, event)
  }
}

function goToPage(page: number) {
  if (!props.pagination) return
  const min = 1
  const max = totalPages.value
  
  if (page < min || page > max) return
  
  emit('page-change', page)
}
</script>