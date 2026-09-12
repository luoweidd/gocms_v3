<template>
  <div class="rich-text-editor">
    <div ref="editorRef" class="editor-container"></div>
    <div class="editor-footer" v-if="showFooter">
      <span class="word-count">字数: {{ wordCount }}</span>
      <span class="toggle-toolbar">
        <a-button size="small" type="text" @click="togglePreview">
          {{ isPreview ? '编辑' : '预览' }}
        </a-button>
      </span>
    </div>
    
    <!-- 预览模态框 -->
    <a-modal
      v-model:visible="isPreview"
      title="内容预览"
      :footer="null"
      width="80%"
      class="preview-modal"
    >
      <div ref="previewRef" class="preview-content"></div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { createEditor, type IDomEditor } from '@wangeditor/editor'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { message } from 'ant-design-vue'

interface Props {
  modelValue?: string
  placeholder?: string
  readOnly?: boolean
  height?: string
  showFooter?: boolean
  autoFocus?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '请输入内容...',
  readOnly: false,
  height: '500px',
  showFooter: true,
  autoFocus: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}>()

const editorRef = ref<HTMLDivElement>()
const previewRef = ref<HTMLDivElement>()
const editor = ref<IDomEditor | null>(null)
const isPreview = ref(false)
const wordCount = ref(0)
// eslint-disable-next-line @typescript-eslint/no-unused-vars
let initialized = false

// 工具栏配置
const toolbarKeys = [
  'headerSelect', // 标题
  'bold', // 粗体
  'italic', // 斜体
  'underline', // 下划线
  'strikeThrough', // 删除线
  '-',
  'fontSize', // 字号
  'fontFamily', // 字体
  'lineHeight', // 行高
  '-',
  'color', // 颜色
  'bgColor', // 背景色
  'clearStyles', // 清除格式
  '-',
  'bulletedList', // 无序列表
  'numberedList', // 有序列表
  'todo', // 待办列表
  '-',
  'justifyLeft', // 左对齐
  'justifyCenter', // 居中对齐
  'justifyRight', // 右对齐
  'justifyFull', // 两端对齐
  '-',
  'indent', // 增加缩进
  'outdent', // 减少缩进
  '-',
  'link', // 链接
  'image', // 图片
  'video', // 视频
  'divider', // 分割线
  '-',
  'code', // 行内代码
  'codeBlock', // 代码块
  '-',
  'table', // 表格
  'quote', // 引用
  '-',
  'undo', // 撤销
  'redo', // 重做
  'fullScreen', // 全屏
  'emotion', // 表情
]

// 编辑器配置
const editorConfig = {
  placeholder: props.placeholder,
  readOnly: props.readOnly,
  autoFocus: props.autoFocus,
  scrollStyle: {
    bottom: '5px',
    left: '5px',
    right: '5px',
    top: '5px',
    isScrollWhileMouseHover: false,
    hoverDelaySeconds: 0,
  },
  menuTooltipPosition: 'below' as const,
  // 配置 TRAILING_TEXT 为空，避免追加内容
  trailingText: undefined,
  // 配置 Z-index
  zIndex: 1,
  // 配置图片上传大小限制（默认 5M）
  uploadImgSize: 5 * 1024 * 1024,
  // 工具栏配置
  toolbarConfig: {
    toolbarKeys,
  },
  // 行高级配置
  rowLayoutAs: 'float',
  // 目录渲染配置
  headingSelectId: 'wange-heading-select',
  onCreated: (editorInstance: IDomEditor) => {
    editor.value = editorInstance
    // 初始化内容
    nextTick(() => {
      editorInstance.setHtml(props.modelValue)
      updateWordCount()
    })
  }
}

// 创建编辑器
onMounted(() => {
  if (editorRef.value) {
    // 配置 onChange 事件
    const configWithChange = {
      ...editorConfig,
      onChange: (editorInstance: IDomEditor) => {
        const html = editorInstance.getHtml()
        wordCount.value = editorInstance.getText().length
        emit('update:modelValue', html)
        emit('change', html)
      }
    }
    // 创建编辑器 - 使用类型断言绕过 TypeScript 检查
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    ;(createEditor as any)(editorRef.value, configWithChange)
    initialized = true
  }
})

// 销毁编辑器
onBeforeUnmount(() => {
  editor.value?.destroy()
  initialized = false
})

// 监听外部内容变化
watch(() => props.modelValue, (newVal) => {
  if (editor.value && newVal !== editor.value.getHtml()) {
    editor.value.setHtml(newVal)
    updateWordCount()
  }
})

// 更新字数
const updateWordCount = () => {
  if (editor.value) {
    wordCount.value = editor.value.getText().length
  }
}

// 切换预览
const togglePreview = async () => {
  isPreview.value = !isPreview.value
  if (isPreview.value) {
    await nextTick()
    if (previewRef.value) {
      previewRef.value.innerHTML = editor.value?.getHtml() || ''
    }
  }
}

// 插入图片（从媒体选择器）
const insertImage = (url: string, alt: string = '', link: string = '') => {
  if (!editor.value) return
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(editor.value as any).insertLink({
    text: alt || '图片',
    href: link || url,
    target: '_blank',
  })
}

// 插入文本
const insertText = (text: string) => {
  if (!editor.value) return
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  ;(editor.value as any).insert(text)
}

// 暴露方法给父组件
defineExpose({
  getHtml: () => editor.value?.getHtml() || '',
  getText: () => editor.value?.getText() || '',
  setHtml: (html: string) => {
    editor.value?.setHtml(html)
    updateWordCount()
  },
  clear: () => {
    editor.value?.clear()
    updateWordCount()
  },
  insertText,
  insertImage
})
</script>

<style scoped>
.rich-text-editor {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.editor-container {
  height: v-bind('height');
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

.toggle-toolbar button {
  color: #1890ff;
}

.preview-content {
  min-height: 300px;
  padding: 16px;
  background-color: #fff;
}

.preview-content :deep(img) {
  max-width: 100%;
  height: auto;
}

/* wangEditor 样式覆盖 */
:deep(.w-e-toolbar) {
  border-bottom: 1px solid #e8e8e8 !important;
  background-color: #fafafa !important;
}

:deep(.w-e-editor) {
  flex: 1;
  overflow-y: auto;
}
</style>