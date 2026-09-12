<template>
  <button
    :type="nativeType"
    :disabled="disabled || loading"
    :class="buttonClasses"
    class="inline-flex items-center justify-center font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all duration-200"
    @click="onClick"
  >
    <slot name="loading">
      <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </slot>
    <slot name="icon-left">
      <slot name="leading"></slot>
    </slot>
    <span v-if="$slots.default && !loadingIconOnly"><slot /></span>
    <slot name="icon-right">
      <slot name="trailing"></slot>
    </slot>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface TailButtonProps {
  type?: 'default' | 'primary' | 'success' | 'warning' | 'danger' | 'info'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  variant?: 'solid' | 'outline' | 'ghost'
  disabled?: boolean
  loading?: boolean
  nativeType?: 'button' | 'submit' | 'reset'
  fullWidth?: boolean
  iconOnly?: boolean
}

const props = withDefaults(defineProps<TailButtonProps>(), {
  type: 'default',
  size: 'md',
  variant: 'solid',
  disabled: false,
  loading: false,
  nativeType: 'button',
  fullWidth: false,
  iconOnly: false,
})

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

const buttonClasses = computed(() => {
  const base: string[] = []

  // Size classes
  const sizeClasses: Record<string, string> = {
    xs: 'px-2 py-1 text-xs',
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base',
  }
  base.push(sizeClasses[props.size] || sizeClasses.md)

  // Width
  if (props.fullWidth) {
    base.push('w-full')
  }

  // Icon only
  if (props.iconOnly) {
    base.push('p-2')
  }

  // Type + Variant combinations
  const typeVariant = getTypeVariantClasses()

  return [...base, ...typeVariant].filter(Boolean).join(' ')
})

const loadingIconOnly = computed(() => props.loading && props.iconOnly)

function getTypeVariantClasses(): string {
  const { type, variant } = props

  if (props.disabled || props.loading) {
    return 'bg-gray-400 text-white cursor-not-allowed'
  }

  switch (type) {
    case 'primary':
      if (variant === 'outline') {
        return 'border border-indigo-600 text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500'
      }
      if (variant === 'ghost') {
        return 'text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500'
      }
      return 'bg-indigo-600 text-white hover:bg-indigo-700 active:bg-indigo-800 focus:ring-indigo-500'

    case 'success':
      if (variant === 'outline') {
        return 'border border-emerald-600 text-emerald-600 hover:bg-emerald-50 focus:ring-emerald-500'
      }
      if (variant === 'ghost') {
        return 'text-emerald-600 hover:bg-emerald-50 focus:ring-emerald-500'
      }
      return 'bg-emerald-600 text-white hover:bg-emerald-700 active:bg-emerald-800 focus:ring-emerald-500'

    case 'warning':
      if (variant === 'outline') {
        return 'border border-amber-600 text-amber-600 hover:bg-amber-50 focus:ring-amber-500'
      }
      if (variant === 'ghost') {
        return 'text-amber-600 hover:bg-amber-50 focus:ring-amber-500'
      }
      return 'bg-amber-500 text-white hover:bg-amber-600 active:bg-amber-700 focus:ring-amber-500'

    case 'danger':
      if (variant === 'outline') {
        return 'border border-red-600 text-red-600 hover:bg-red-50 focus:ring-red-500'
      }
      if (variant === 'ghost') {
        return 'text-red-600 hover:bg-red-50 focus:ring-red-500'
      }
      return 'bg-red-600 text-white hover:bg-red-700 active:bg-red-800 focus:ring-red-500'

    case 'info':
      if (variant === 'outline') {
        return 'border border-gray-600 text-gray-600 hover:bg-gray-50 focus:ring-gray-500'
      }
      if (variant === 'ghost') {
        return 'text-gray-600 hover:bg-gray-50 focus:ring-gray-500'
      }
      return 'bg-gray-800 text-white hover:bg-gray-900 active:bg-gray-700 focus:ring-gray-500'

    default: // default
      if (variant === 'outline') {
        return 'border border-gray-300 text-gray-700 hover:bg-gray-50 focus:ring-gray-400'
      }
      if (variant === 'ghost') {
        return 'text-gray-600 hover:bg-gray-100 focus:ring-gray-400'
      }
      return 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 focus:ring-gray-400'
  }
}

function onClick(event: MouseEvent) {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>