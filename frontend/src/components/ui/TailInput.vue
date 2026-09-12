<template>
  <div :class="containerClasses">
    <!-- Label -->
    <label v-if="label" class="ta-label">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>

    <!-- Input Container -->
    <div :class="inputContainerClasses">
      <!-- Leading Icon -->
      <slot name="icon-left">
        <svg v-if="iconLeft" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </slot>

      <!-- Native Input -->
      <input
        v-if="type !== 'textarea'"
        :value="modelValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :maxlength="maxlength"
        :class="inputClasses"
        @input="onInput"
        @change="onChange"
        @focus="onFocus"
        @blur="onBlur"
      />

      <!-- Textarea -->
      <textarea
        v-if="type === 'textarea'"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :rows="rows"
        :class="inputClasses + ' resize-none'"
        @input="onInput"
        @change="onChange"
        @focus="onFocus"
        @blur="onBlur"
      />

      <!-- Trailing Icon -->
      <slot name="icon-right">
        <button
          v-if="showPasswordToggle && modelValue"
          type="button"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600"
          @click="togglePassword"
        >
          <svg v-if="isPasswordVisible" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
          </svg>
          <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        </button>
      </slot>
    </div>

    <!-- Helper text -->
    <p v-if="helpText" class="mt-1.5 text-xs text-gray-500">{{ helpText }}</p>

    <!-- Error message -->
    <p v-if="error" class="mt-1.5 text-xs text-red-600">{{ error }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

export interface TailInputProps {
  modelValue?: string | number
  type?: 'text' | 'password' | 'email' | 'number' | 'tel' | 'url' | 'textarea'
  placeholder?: string
  label?: string
  required?: boolean
  disabled?: boolean
  readonly?: boolean
  error?: string
  helpText?: string
  maxlength?: number
  rows?: number
  iconLeft?: boolean
  showPasswordToggle?: boolean
  fullWidth?: boolean
}

const props = withDefaults(defineProps<TailInputProps>(), {
  modelValue: '',
  type: 'text',
  placeholder: '',
  label: '',
  required: false,
  disabled: false,
  readonly: false,
  error: '',
  helpText: '',
  rows: 4,
  showPasswordToggle: false,
  fullWidth: true,
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  input: [event: Event]
  change: [event: Event]
  focus: [event: FocusEvent]
  blur: [event: FocusEvent]
}>()

const isPasswordVisible = ref(false)

const containerClasses = computed(() => {
  const classes = ['ta-input-wrapper']
  if (props.fullWidth) classes.push('w-full')
  return classes.join(' ')
})

const inputContainerClasses = computed(() => {
  const classes = ['relative rounded-md transition-colors duration-200']
  if (props.error) {
    classes.push('border-red-500 focus-within:border-red-500 focus-within:ring-red-200')
  } else if (props.readonly) {
    classes.push('border-gray-300 bg-gray-50')
  } else {
    classes.push('border-gray-300 focus-within:border-indigo-500 focus-within:ring-indigo-200')
  }
  return classes.join(' ')
})

const inputClasses = computed(() => {
  const base = [
    'block w-full',
    'px-3 py-2',
    'text-gray-900',
    'placeholder-gray-400',
    'rounded-md',
    'border-0',
    'bg-white',
    'text-sm',
    'leading-6',
    'transition-colors duration-200',
  ]

  if (props.disabled) {
    base.push('bg-gray-100 text-gray-500 cursor-not-allowed')
  } else if (props.readonly) {
    base.push('cursor-default')
  } else {
    base.push('focus:ring-2 focus:outline-none')
  }

  return base.join(' ')
})

function onInput(event: Event) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  emit('update:modelValue', target.value)
  emit('input', event)
}

function onChange(event: Event) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement
  emit('update:modelValue', target.value)
  emit('change', event)
}

function onFocus(event: FocusEvent) {
  emit('focus', event)
}

function onBlur(event: FocusEvent) {
  emit('blur', event)
}

function togglePassword() {
  isPasswordVisible.value = !isPasswordVisible.value
}
</script>