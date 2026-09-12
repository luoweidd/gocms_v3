<template>
  <div :class="cardClasses">
    <!-- Header -->
    <div v-if="$slots.header || title || $slots['title-action']" class="flex items-center justify-between pb-4 border-b border-gray-200">
      <div class="flex items-center gap-3 flex-1">
        <!-- Icon -->
        <slot name="header-icon">
          <div v-if="headerIcon" class="w-10 h-10 rounded-lg flex items-center justify-center" :class="headerIconBg">
            <svg class="w-5 h-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
        </slot>
        
        <!-- Title -->
        <div v-if="title" class="flex-1">
          <h3 class="text-lg font-semibold text-gray-800">{{ title }}</h3>
          <p v-if="subtitle" class="text-sm text-gray-500 mt-0.5">{{ subtitle }}</p>
        </div>
      </div>

      <!-- Title Action -->
      <slot name="title-action" />
    </div>

    <!-- Body -->
    <div v-if="$slots.default" class="py-4">
      <slot />
    </div>

    <!-- Footer -->
    <div v-if="$slots.footer" class="flex items-center justify-end pt-4 border-t border-gray-200 gap-3">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface TailCardProps {
  title?: string
  subtitle?: string
  padded?: boolean
  hover?: boolean
  headerIcon?: boolean
  headerIconBg?: string
  bordered?: boolean
}

const props = withDefaults(defineProps<TailCardProps>(), {
  title: '',
  subtitle: '',
  padded: true,
  hover: false,
  headerIcon: false,
  headerIconBg: 'bg-indigo-500',
  bordered: false,
})

const cardClasses = computed(() => {
  const base = [
    'bg-white',
    'rounded-lg',
  ]

  if (props.bordered) {
    base.push('border border-gray-200')
  } else {
    base.push('shadow-sm')
  }

  if (props.hover) {
    base.push('hover:-translate-y-1 hover:shadow-md transition-all duration-200')
  }

  if (props.padded) {
    base.push('p-6')
  }

  return base.join(' ')
})
</script>