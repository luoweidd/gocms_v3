<template>
  <div class="markdown-editor">
    <div ref="vditorRef" class="vditor-container"></div>
    <div class="editor-footer" v-if="showFooter">
      <span class="word-count">字数: {{ wordCount }}</span>
      <span class="toolbar-actions">
        <a-button size="small" type="text" @click="togglePreview">
          {{ isPreview ? '编辑' : '预览' }}
        </a-button>
        <a-button size="small" type="text" @click="toggleSource">
          {{ isSourceCode ? '渲染' : '源码' }}
        </a-button>
        <a-button size="small" type="text" @click="exportMarkdown">
          导出
        </a-button>
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

interface Props {
  modelValue?: string
  placeholder?: string
  readOnly?: boolean
  height?: string
  showFooter?: boolean
  autoFocus?: boolean
  mode?: 'ir' | 'sv' | 'wysiwyg'
  previewTheme?: 'default' | 'light' | 'dark'
  toolbarConfig?: {
    pin?: boolean
  }
  hint?: string[]
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入Markdown内容...',
  readOnly: false,
  height: '500px',
  showFooter: true,
  autoFocus: false,
  mode: 'ir',
  previewTheme: 'default',
  toolbarConfig: () => ({ pin: true }),
  hint: () => ['@', ':', '#', '/', '[', '']
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const vditorRef = ref<HTMLDivElement>()
const wordCount = ref(0)
const isPreview = ref(false)
const isSourceCode = ref(false)
// eslint-disable-next-line @typescript-eslint/no-explicit-any
let vditor: any = null

// 初始化 Vditor
onMounted(() => {
  if (!vditorRef.value) return
  
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  vditor = new Vditor(vditorRef.value, {
    mode: props.mode,
    placeholder: props.placeholder,
    height: props.height.replace('px', ''),
    readonlyEnable: props.readOnly,
    autoFocus: props.autoFocus,
    theme: 'classic',
    // 使用本地 CDN 资源（通过 Vite 中间件提供）
    cdn: '/vditor',
    lang: 'zh_CN',
    preview: {
      theme: { name: props.previewTheme },
      markdown: { render: { force: true } }
    },
    cache: { enable: false },
    toolbarConfig: props.toolbarConfig,
    hint: function (_prefix: string) {
      switch (_prefix) {
        case '@': return ['@A', '@B', '@C']
        case ':': return ['smile:', 'tumor:', 'call:', 'pause']
        case '#': return ['# 一级标题', '## 二级标题', '### 三级标题']
        case '/': return ['/image', '/video', '/audio', '/file']
        case '[': return ['[链接', '[图片']
        default: return []
      }
    },
    // Vditor v3.10.8 使用回调函数处理事件
    input: (value: string) => {
      emit('update:modelValue', value)
      emit('change', value)
      updateWordCount()
    },
    after: () => {
      // 设置初始内容
      if (props.modelValue && vditor) {
        vditor.setValue(props.modelValue)
      }
      updateWordCount()
    }
  } as any)
})

// 监听外部内容变化
watch(() => props.modelValue, (newVal) => {
  if (vditor && newVal !== vditor.getValue()) {
    vditor.setValue(newVal)
    updateWordCount()
  }
})

// 销毁 Vditor
onBeforeUnmount(() => {
  vditor?.destroy()
  vditor = null
})

// 更新字数统计
const updateWordCount = () => {
  if (vditor) {
    const text = vditor.getValue()
    wordCount.value = text.length
  }
}

// 切换预览模式
const togglePreview = () => {
  isPreview.value = !isPreview.value
  if (isPreview.value && vditor) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(vditor as any).preview({ mode: 'offset', position: 'center' })
  } else {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(vditor as any).preview?.stop()
  }
}

// 切换源码模式
const toggleSource = () => {
  isSourceCode.value = !isSourceCode.value
  if (isSourceCode.value && vditor) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(vditor as any).setMode('sv')
  } else if (!isSourceCode.value && vditor) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(vditor as any).setMode('ir')
  }
}

// 导出Markdown
const exportMarkdown = () => {
  const content = vditor?.getValue() || ''
  const blob = new Blob([content], { type: 'text/markdown;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `article-${Date.now()}.md`
  a.click()
  URL.revokeObjectURL(url)
}

// 插入图片
const insertImage = (url: string, alt: string = '') => {
  if (!vditor) return
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(vditor as any).insertText(`![${alt}](${url})`)
}

// 插入链接
const insertLink = (url: string, text: string = '') => {
  if (!vditor) return
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(vditor as any).insertText(`[${text || url}](${url})`)
}

// 插入代码块
const insertCodeBlock = (code: string, language: string = '') => {
  if (!vditor) return
  const codeText = `\`\`\`${language}\n${code}\n\`\`\``
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(vditor as any).insertText(codeText)
}

// 暴露方法给父组件
defineExpose({
  getValue: () => vditor?.getValue() || '',
  getHTML: () => {
    if (!vditor) return ''
    // Vditor 提供 marked 功能来渲染Markdown为HTML
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return (Vditor as any).markdownToCode(vditor.getValue(), {
      theme: 'default',
      method: 'preprocess'
    }) || ''
  },
  setValue: (value: string) => {
    vditor?.setValue(value)
    updateWordCount()
  },
  clear: () => {
    vditor?.clear()
    updateWordCount()
  },
  insertImage,
  insertLink,
  insertCodeBlock,
  getMarkdown: () => vditor?.getValue() || ''
})
</script>

<style scoped>
.markdown-editor {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.vditor-container {
  min-height: 300px;
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  background-color: #fafafa;
  border-top: 1px solid #e8e8e8;
  font-size: 12px;
  color: #666;
}

.word-count {
  color: #999;
}

.toolbar-actions button {
  color: #1890ff;
}
</style>