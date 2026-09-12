<template>
  <div class="article-editor">
    <!-- 编辑器头部 -->
    <div class="editor-header">
      <a-tabs v-model:activeKey="editorMode" type="card" @change="switchMode">
        <a-tab-key key="rich">
          <span><i class="fa fa-font"></i> 富文本模式</span>
        </a-tab-key>
        <a-tab-key key="markdown">
          <span><i class="fa fa-markdown"></i> Markdown模式</span>
        </a-tab-key>
      </a-tabs>
      
      <div class="editor-actions">
        <span class="auto-save-status" v-if="autoSaveEnabled">
          <a-spin size="small" v-if="isSaving" />
          <span v-else class="save-indicator">{{ saveStatusText }}</span>
        </span>
        <a-button size="small" @click="toggleFullscreen">
          {{ isFullscreen ? '退出全屏' : '全屏编辑' }}
        </a-button>
      </div>
    </div>

    <!-- 编辑器主体 -->
    <div :class="['editor-body', { 'fullscreen': isFullscreen }]">
      <!-- 富文本编辑器 -->
      <RichTextEditor
        v-if="editorMode === 'rich'"
        ref="richEditorRef"
        v-model="currentContent"
        :placeholder="placeholder"
        :read-only="readOnly"
        :height="height"
        :show-footer="true"
        @change="onContentChange"
      />
      
      <!-- Markdown编辑器 -->
      <MarkdownEditor
        v-else
        ref="markdownEditorRef"
        v-model="currentContent"
        :placeholder="placeholder"
        :read-only="readOnly"
        :height="height"
        :show-footer="true"
        @change="onContentChange"
      />
    </div>

    <!-- 底部工具栏 -->
    <div class="editor-footer" v-if="!isFullscreen">
      <div class="footer-left">
        <a-button size="small" @click="openMediaPicker">
          插入图片
        </a-button>
        <a-button size="small" @click="insertLink">
          插入链接
        </a-button>
        <a-dropdown v-if="editorMode === 'markdown'">
          <a-button size="small">
            插入代码 <i class="anticon anticon-down"></i>
          </a-button>
          <template #overlay>
            <a-menu>
              <a-menu-item @click="insertCode('javascript')">JavaScript</a-menu-item>
              <a-menu-item @click="insertCode('python')">Python</a-menu-item>
              <a-menu-item @click="insertCode('css')">CSS</a-menu-item>
              <a-menu-item @click="insertCode('html')">HTML</a-menu-item>
              <a-menu-item @click="insertCode('sql')">SQL</a-menu-item>
              <a-menu-item @click="insertCode('bash')">Bash</a-menu-item>
            </a-menu>
          </template>
        </a-dropdown>
      </div>
      
      <div class="footer-right">
        <span class="word-count">字数: {{ wordCount }}</span>
        <span class="read-time">预计阅读: {{ readTime }}</span>
      </div>
    </div>

    <!-- 媒体选择器弹窗 -->
    <MediaPicker
      v-model:visible="showMediaPicker"
      :multiple="false"
      @select="onMediaSelect"
    />

    <!-- 链接输入弹窗 -->
    <a-modal
      v-model:visible="showLinkDialog"
      title="插入链接"
      @ok="confirmInsertLink"
      @cancel="cancelLinkDialog"
    >
      <a-form :model="linkForm" layout="vertical">
        <a-form-item label="链接地址">
          <a-input v-model:value="linkForm.url" placeholder="https://example.com" />
        </a-form-item>
        <a-form-item label="显示文本">
          <a-input v-model:value="linkForm.text" placeholder="点击访问" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { message } from 'ant-design-vue'
import RichTextEditor from './RichTextEditor.vue'
import MarkdownEditor from './MarkdownEditor.vue'
import MediaPicker from './MediaPicker.vue'

interface Props {
  modelValue?: string
  placeholder?: string
  readOnly?: boolean
  height?: string
  autoSave?: boolean
  autoSaveInterval?: number
  uploadUrl?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '开始创作你的作品...',
  readOnly: false,
  height: '500px',
  autoSave: true,
  autoSaveInterval: 30000, // 30秒自动保存
  uploadUrl: '/api/media/upload'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
  (e: 'save', value: { content: string, mode: string }): void
}>()

// 编辑器模式
const editorMode = ref<'rich' | 'markdown'>('rich')

// 内容绑定
const currentContent = ref(props.modelValue)

// 子组件引用
const richEditorRef = ref<any>(null)
const markdownEditorRef = ref<any>(null)

// UI状态
const isFullscreen = ref(false)
const showMediaPicker = ref(false)
const showLinkDialog = ref(false)
const isSaving = ref(false)
const saveStatusText = ref('已保存')

// 链接表单
const linkForm = ref({
  url: '',
  text: ''
})

// 字数统计
const wordCount = computed(() => {
  return currentContent.value?.length || 0
})

// 预计阅读时间（每分钟300字）
const readTime = computed(() => {
  const minutes = Math.ceil(wordCount.value / 300)
  return minutes < 1 ? '<1分钟' : `${minutes}分钟`
})

// 自动保存状态
const autoSaveEnabled = computed(() => props.autoSave && !props.readOnly)
let saveTimer: ReturnType<typeof setInterval> | null = null
let draftTimer: ReturnType<typeof setInterval> | null = null

// 初始化
onMounted(() => {
  // 加载草稿
  loadDraft()
  
  // 启动自动保存
  if (autoSaveEnabled.value) {
    startAutoSave()
  }
})

// 监听外部内容变化
watch(() => props.modelValue, (newVal) => {
  if (newVal !== currentContent.value) {
    currentContent.value = newVal
  }
})

// 切换编辑器模式
const switchMode = (mode: string) => {
  editorMode.value = mode as 'rich' | 'markdown'
  // 保存当前模式偏好
  localStorage.setItem('article-editor-mode', mode)
}

// 内容变化处理
const onContentChange = (content: string) => {
  currentContent.value = content
  saveDraft()
}

// 全屏切换
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value
}

// 打开媒体选择器
const openMediaPicker = () => {
  showMediaPicker.value = true
}

// 媒体选择回调
const onMediaSelect = (media: any[]) => {
  if (media.length === 0) return
  
  const item = media[0]
  const url = item.url || item.path
  
  if (editorMode.value === 'rich') {
    richEditorRef.value?.insertImage(url, item.name || '图片')
  } else {
    markdownEditorRef.value?.insertImage(url, item.name || '图片')
  }
  
  message.success('图片已插入')
}

// 插入链接
const insertLink = () => {
  showLinkDialog.value = true
  linkForm.value = {
    url: '',
    text: ''
  }
}

// 确认插入链接
const confirmInsertLink = () => {
  if (!linkForm.value.url) {
    message.warning('请输入链接地址')
    return
  }
  
  const text = linkForm.value.text || linkForm.value.url
  
  if (editorMode.value === 'rich') {
    richEditorRef.value?.insertLink(linkForm.value.url, text)
  } else {
    markdownEditorRef.value?.insertLink(linkForm.value.url, text)
  }
  
  showLinkDialog.value = false
  message.success('链接已插入')
}

// 取消链接对话框
const cancelLinkDialog = () => {
  showLinkDialog.value = false
}

// 插入代码块
const insertCode = (language: string) => {
  const code = prompt(`请输入 ${language} 代码:`)
  if (!code) return
  
  if (editorMode.value === 'rich') {
    richEditorRef.value?.insertText(`\`\`\`${language}\n${code}\n\`\`\``)
  } else {
    markdownEditorRef.value?.insertCodeBlock(code, language)
  }
}

// 草稿管理
const DRAFT_KEY = 'article-editor-draft'

const saveDraft = () => {
  try {
    localStorage.setItem(DRAFT_KEY, JSON.stringify({
      content: currentContent.value,
      mode: editorMode.value,
      timestamp: Date.now()
    }))
  } catch (e) {
    console.error('Failed to save draft:', e)
  }
}

const loadDraft = () => {
  try {
    const draft = localStorage.getItem(DRAFT_KEY)
    if (!draft) return false
    
    const data = JSON.parse(draft)
    const elapsed = Date.now() - data.timestamp
    
    // 24小时内的草稿有效
    if (elapsed < 24 * 60 * 60 * 1000 && data.content) {
      currentContent.value = data.content
      editorMode.value = data.mode || 'rich'
      
      // 提示用户是否有草稿
      if (data.content.length > 0) {
        message.info('已恢复上次未保存的内容')
      }
      
      return true
    }
  } catch (e) {
    console.error('Failed to load draft:', e)
  }
  return false
}

// 清除草稿
const clearDraft = () => {
  localStorage.removeItem(DRAFT_KEY)
}

// 自动保存
const startAutoSave = () => {
  saveTimer = setInterval(async () => {
    if (currentContent.value && currentContent.value !== props.modelValue) {
      isSaving.value = true
      saveStatusText.value = '保存中...'
      
      try {
        // TODO: 调用自动保存API
        // await autoSaveArticle(currentContent.value)
        
        saveStatusText.value = '已保存'
        message.success('自动保存成功', 1)
      } catch (error) {
        saveStatusText.value = '保存失败'
        message.error('自动保存失败')
      } finally {
        setTimeout(() => {
          isSaving.value = false
        }, 2000)
      }
    }
  }, props.autoSaveInterval)
}

// 获取当前编辑器内容
const getContent = () => {
  return currentContent.value
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
const getHTML = async (): Promise<string> => {
  if (editorMode.value === 'rich') {
    return richEditorRef.value?.getHtml() || ''
  } else {
    // Markdown转为HTML
    const markdown = markdownEditorRef.value?.getValue() || ''
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    return new Promise<string>((resolve) => {
      ;(window as any).Vditor.preview(document.createElement('div'), markdown, {
        mode: 'light'
      })
      setTimeout(() => resolve(document.querySelector('.markdown-body')?.innerHTML || ''), 100)
    })
  }
}

// 设置内容
const setContent = (content: string) => {
  currentContent.value = content
  saveDraft()
}

// 清空内容
const clear = () => {
  currentContent.value = ''
  clearDraft()
  saveStatusText.value = '已清空'
}

// 组件卸载前保存
onBeforeUnmount(() => {
  if (saveTimer) {
    clearInterval(saveTimer)
  }
  saveDraft()
})

// 暴露方法给父组件
defineExpose({
  getContent,
  getHTML,
  setContent,
  clear,
  getCurrentMode: () => editorMode.value,
  getWordCount: () => wordCount.value
})
</script>

<style scoped>
.article-editor {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: hidden;
  background-color: #fff;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #fafafa;
  border-bottom: 1px solid #e8e8e8;
}

.editor-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.auto-save-status {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #52c41a;
  margin-right: 8px;
}

.save-indicator {
  color: #52c41a;
}

.editor-body {
  min-height: 300px;
}

.editor-body.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background-color: #fff;
  display: flex;
  flex-direction: column;
}

.editor-body.fullscreen > * {
  flex: 1;
}

.editor-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 16px;
  background-color: #fafafa;
  border-top: 1px solid #e8e8e8;
  font-size: 12px;
}

.footer-left {
  display: flex;
  gap: 8px;
  align-items: center;
}

.footer-right {
  display: flex;
  gap: 16px;
  align-items: center;
  color: #999;
}

.word-count,
.read-time {
  color: #999;
}

/* 全屏模式下的样式 */
:deep(.fullscreen .editor-body > div) {
  height: 100% !important;
  min-height: unset !important;
}
</style>