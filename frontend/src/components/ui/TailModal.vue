<template>
  <!-- 打开状态：显示遮罩和弹窗 -->
  <Teleport to="body">
    <Transition name="modal-fade">
      <div 
        v-if="open" 
        class="fixed inset-0 z-[1000] overflow-y-auto bg-black/50"
        aria-dialog
        role="dialog"
        @click.self="closeOnBackdrop && close()"
      >
        <!-- Modal Container -->
        <div class="flex min-h-full items-center justify-center p-4">
          <div 
            ref="modalRef"
            class="bg-white rounded-xl shadow-lg transition-all duration-200 w-full"
            :class="modalClasses"
            role="dialog"
            :aria-labelledby="titleId"
            :aria-describedby="descId"
            @keydown.escape="handleEscape"
          >
            <!-- Header -->
            <div v-if="$slots.header || title" class="flex items-center justify-between p-6 pb-4 border-b border-gray-200">
              <div v-if="title" :id="titleId" class="text-lg font-semibold text-gray-800">
                {{ title }}
              </div>
              <slot name="header">
                <!-- 自定义头部内容 -->
              </slot>
              <!-- Close Button -->
              <button
                v-if="closable"
                class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
                @click.stop="close"
                aria-label="Close"
              >
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Body -->
            <div v-if="$slots.default" :id="descId" class="p-6">
              <slot />
            </div>

            <!-- Footer -->
            <div v-if="$slots.footer" class="flex items-center justify-end gap-3 p-6 pt-4">
              <slot name="footer" />
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

export interface TailModalProps {
  open: boolean
  title?: string
  width?: string
  closable?: boolean
  closeOnBackdrop?: boolean
}

const props = withDefaults(defineProps<TailModalProps>(), {
  title: '',
  width: 'max-w-2xl',
  closable: true,
  closeOnBackdrop: true,
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  close: []
}>()

const titleId = ref(`modal-title-${Date.now()}`)
const descId = ref(`modal-desc-${Date.now()}`)
const modalRef = ref<HTMLElement | null>(null)

const modalClasses = computed(() => {
  return [
    props.width,
    'max-h-[90vh]',
    'overflow-auto',
  ].join(' ')
})

function close() {
  emit('update:open', false)
}

function handleEscape(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.closable) {
    close()
  }
}
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>