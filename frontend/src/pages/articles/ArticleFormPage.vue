<template>
  <div class="page-container max-w-6xl">
    <!-- Page Header -->
    <div class="mb-6">
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-2">
          <button class="text-gray-500 hover:text-gray-700 transition-colors" @click="$router.back()">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </button>
          <div>
            <h1 class="text-2xl font-bold text-gray-900 mb-1">{{ isEdit ? '编辑文章' : '发布文章' }}</h1>
            <p class="text-sm text-gray-500">{{ isEdit ? '修改文章内容' : '创建新的文章内容' }}</p>
          </div>
        </div>
        <div class="flex gap-2">
          <!-- Editor Type Toggle -->
          <div class="flex items-center bg-gray-100 rounded-lg p-1">
            <button 
              :class="editorType === 'markdown' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-600 hover:text-gray-900'" 
              class="px-4 py-2 rounded-md text-sm font-medium transition-all flex items-center gap-2"
              @click="switchEditorType('markdown')"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Markdown
            </button>
            <button 
              :class="editorType === 'rich' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-600 hover:text-gray-900'" 
              class="px-4 py-2 rounded-md text-sm font-medium transition-all flex items-center gap-2"
              @click="switchEditorType('rich')"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L9.21 19.21a2.5 2.5 0 01-3.536-3.536L15.232 5.232z" />
              </svg>
              富文本
            </button>
          </div>
          <!-- Preview Toggle -->
          <button 
            :class="showPreview ? 'bg-indigo-600 text-white' : 'bg-white text-gray-700 border border-gray-300'" 
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all hover:shadow-md flex items-center gap-2"
            @click="showPreview = !showPreview"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            {{ showPreview ? '隐藏预览' : '预览' }}
          </button>
          <router-link 
            v-if="articleId && article" 
            :to="`/articles/${article.id}`"
            class="px-4 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
            在新窗口打开
          </router-link>
        </div>
      </div>
    </div>

    <!-- Form Card -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div class="flex" :class="showPreview ? '' : ''">
        <!-- Editor Area -->
        <div class="flex-1 p-6" :class="showPreview ? '' : ''">
          <form @submit.prevent="handleSubmit">
            <!-- Title -->
            <div class="mb-5">
              <label class="ta-label">文章标题 <span class="text-red-500">*</span></label>
              <input v-model="formData.title" type="text" class="ta-input w-full" placeholder="请输入文章标题" required />
            </div>

            <!-- Category and Status -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-5 mb-5">
              <div>
                <label class="ta-label">分类 <span class="text-red-500">*</span></label>
                <select v-model="formData.categoryId" class="ta-input w-full" required>
                  <option :value="0">请选择分类</option>
                  <option v-for="c in categoryTree" :key="c.id" :value="c.id">{{ getIndent(c.depth || 0) }}{{ c.name }}</option>
                </select>
              </div>
              <div>
                <label class="ta-label">状态</label>
                <select v-model="formData.status" class="ta-input w-full">
                  <option :value="0">草稿</option>
                  <option :value="1">已发布</option>
                  <option :value="2">已下架</option>
                </select>
              </div>
            </div>

            <!-- Content Editor -->
            <div class="mb-5">
              <label class="ta-label">文章内容 <span class="text-red-500">*</span></label>
              
              <!-- Vditor Editor (Markdown) -->
              <div v-if="editorType === 'markdown'" ref="vditorRef" class="w-full h-[400px]"></div>
              
              <!-- WangEditor Editor (Rich Text) - Using Vue Component -->
              <div v-else-if="editorType === 'rich'" class="wang-editor-container w-full">
                <!-- HTML 源码切换按钮（固定在编辑器上方） -->
                <div class="html-toggle-bar p-2 bg-gray-50 border border-gray-200 rounded-lg mb-2 flex gap-2 items-center">
                  <button 
                    type="button"
                    @click="toggleHtmlSource"
                    class="px-3 py-1.5 text-sm bg-indigo-600 text-white rounded hover:bg-indigo-700 transition-colors"
                  >
                    {{ showHtmlSource ? '📝 可视化编辑' : '💻 HTML 源码' }}
                  </button>
                  <!-- 富文本工具栏（固定在HTML按钮后面，始终可见） -->
                  <div class="flex gap-1 flex-wrap">
                    <!-- Headings -->
                    <button type="button" @click="insertHeading(1)" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="一级标题">H1</button>
                    <button type="button" @click="insertHeading(2)" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="二级标题">H2</button>
                    <button type="button" @click="insertHeading(3)" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="三级标题">H3</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Text Formatting -->
                    <button type="button" @click="toggleBold" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50 font-bold" title="加粗">B</button>
                    <button type="button" @click="toggleItalic" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50 italic" title="斜体">I</button>
                    <button type="button" @click="toggleUnderline" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50 underline" title="下划线">U</button>
                    <button type="button" @click="toggleStrikeThrough" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50 line-through" title="删除线">S</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Text Color -->
                    <button type="button" @click="changeTextColor" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="文字颜色">
                      <span class="text-red-500">A</span>
                    </button>
                    <button type="button" @click="changeBgColor" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="背景颜色">
                      <span class="bg-yellow-400 px-1 rounded">B</span>
                    </button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- List -->
                    <button type="button" @click="insertUnorderedList" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="无序列表">• 列表</button>
                    <button type="button" @click="insertOrderedList" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="有序列表">1. 列表</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Alignment -->
                    <button type="button" @click="alignLeft" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="左对齐">≡←</button>
                    <button type="button" @click="alignCenter" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="居中对齐">≡↔</button>
                    <button type="button" @click="alignRight" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="右对齐">→≡</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Insert -->
                    <button type="button" @click="insertLink" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="插入链接">🔗</button>
                    <button type="button" @click="insertImage" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="插入图片">🖼️</button>
                    <button type="button" @click="insertCodeBlock" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="代码块">{ }</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Table -->
                    <button type="button" @click="insertTable" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="插入表格">▦</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Other -->
                    <button type="button" @click="insertQuote" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="引用">❝</button>
                    <button type="button" @click="insertHorizontalRule" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="分割线">—</button>
                    <button type="button" @click="clearStyles" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="清除格式">🚫</button>
                    <span class="mx-1 text-gray-300">|</span>
                    <!-- Undo/Redo -->
                    <button type="button" @click="undo" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="撤销">↩</button>
                    <button type="button" @click="redo" class="px-2 py-1 text-sm bg-white border rounded hover:bg-gray-50" title="重做">↪</button>
                  </div>
                </div>
                <!-- HTML 源码编辑器 -->
                <textarea
                  v-if="showHtmlSource"
                  v-model="formData.content"
                  class="wang-editor-container w-full h-[400px] p-4 font-mono text-sm bg-gray-900 text-green-400 border-none rounded-lg resize-none focus:outline-none"
                  placeholder="请输入 HTML 代码..."
                ></textarea>
                <!-- 可视化编辑器（包含自带工具栏） -->
                <div v-else class="border border-gray-200 rounded-lg overflow-hidden">
                  <Editor 
                    v-model="formData.content" 
                    :defaultConfig="editorDefaultConfig" 
                    @onChange="handleEditorChange"
                    @onCreated="handleEditorCreatedCustom"
                  />
                </div>
              </div>
              
              <!-- Fallback textarea for when editors are loading -->
              <textarea 
                v-model="formData.content" 
                class="ta-input w-full hidden" 
                rows="12" 
                placeholder="请输入文章内容（支持 Markdown）"
                ref="fallbackTextarea"
              ></textarea>
            </div>

            <!-- Tags -->
            <div class="mb-5">
              <label class="ta-label">标签</label>
              <input v-model="tagInput" type="text" class="ta-input w-full" placeholder="输入标签后按回车添加" @keydown.enter.prevent="addTag" />
              <div class="flex flex-wrap gap-2 mt-2">
                <span v-for="(tag, i) in formData.tags" :key="i" class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
                  {{ tag }}
                  <button type="button" @click="formData.tags.splice(i, 1)" class="ml-2 text-red-500 hover:text-red-700 font-bold">&times;</button>
                </span>
              </div>
            </div>

            <!-- Cover Image -->
            <div class="mb-5">
              <label class="ta-label">封面图</label>
              <div class="flex items-center gap-4">
                <div v-if="formData.coverUrl" class="w-32 h-20 rounded-lg overflow-hidden bg-gray-100 border border-gray-200">
                  <img :src="formData.coverUrl" class="w-full h-full object-cover" alt="封面图" />
                </div>
                <button type="button" class="ta-btn ta-btn-outline" @click="handleCoverUpload">
                  <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  {{ formData.coverUrl ? '更换封面' : '上传封面' }}
                </button>
              </div>
            </div>

            <!-- Submit Buttons -->
            <div class="flex gap-3 pt-4 border-t border-gray-200">
              <button type="submit" class="ta-btn ta-btn-primary">
                {{ isEdit ? '保存修改' : '发布文章' }}
              </button>
              <button type="button" class="ta-btn ta-btn-outline" @click="$router.back()">取消</button>
            </div>
          </form>
        </div>

        <!-- Preview Area -->
        <div v-if="showPreview" class="w-1/2 border-l border-gray-200 p-6 bg-gray-50 overflow-y-auto">
          <h3 class="text-lg font-semibold text-gray-900 mb-4 flex items-center gap-2">
            <svg class="w-5 h-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            实时预览
          </h3>
          <div class="prose prose-lg max-w-none bg-white p-6 rounded-lg shadow-sm" v-html="renderedContent"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getCategoryTree, getArticle, createArticle, updateArticle } from '@/api/article'
import { message } from '@/utils/message'
import TailButton from '@/components/ui/TailButton.vue'
import TailCard from '@/components/ui/TailCard.vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
// 导入 wangEditor CSS 样式
import '@wangeditor/editor/dist/css/style.css'
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck
import { Editor } from '@wangeditor/editor-for-vue'

interface CategoryItem {
  id: number
  name: string
  depth?: number
  children?: CategoryItem[]
}

interface ArticleItem {
  id: number
  title: string
  summary?: string
  description?: string
  content: string
  cover?: string
  category_id: number
  category_name?: string
  author_id: number
  author_name?: string
  status: 0 | 1 | 2
  view_count: number
  tags?: Array<{ id: number; name: string }>
  tag_ids?: number[]
  published_at?: string
  created_at: string
  updated_at: string
}

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const articleId = computed(() => route.params.id ? Number(route.params.id) : null)

const formData = reactive({
  title: '',
  categoryId: 0,
  content: '',
  status: 0 as 0 | 1 | 2,
  coverUrl: '',
  tags: [] as string[]
})

// Reference to the article being edited
const article = ref<ArticleItem | null>(null)

const tagInput = ref('')
const categoryTree = ref<CategoryItem[]>([])
const editorType = ref<'markdown' | 'rich'>('markdown')
const showPreview = ref(false)
const showHtmlSource = ref(false)
const vditorRef = ref<HTMLElement>()
const fallbackTextarea = ref<HTMLTextAreaElement>()
let vditor: Vditor | null = null
let editor: any = null

// wangEditor 配置（不要在这里放 onChange）
const editorDefaultConfig = {
  placeholder: '请输入文章内容',
  // 注意：onChange 必须通过 @onChange 事件传递，不能放在 defaultConfig 中
  toolbarConfig: {
    toolbarKeys: [
      'headings',
      '|',
      'bold',
      'italic',
      'underline',
      'strikeThrough',
      '|',
      'color',
      'bgColor',
      'clearStyles',
      '|',
      'bulletedList',
      'orderedList',
      '|',
      'justifyLeft',
      'justifyCenter',
      'justifyRight',
      '|',
      'insertLink',
      'insertImage',
      '|',
      'insertTable',
      'codeBlock',
      '|',
      'quote',
      '|',
      'undo',
      'redo',
    ],
  },
  // 主题
  theme: 'light',
  // scroll 配置
  scroll: true,
}

const handleEditorCreatedCustom = (outerEditor: any) => {
  editor = outerEditor
  console.log('WangEditor Vue component created successfully')
  if (formData.content) {
    outerEditor.setHtml(formData.content)
  }
}

const handleEditorChange = (editor: any) => {
  console.log('Editor content changed')
}


const toggleHtmlSource = () => {
  showHtmlSource.value = !showHtmlSource.value
}

// 富文本工具栏方法 - 通过修改 HTML 内容实现
const getEditorContent = (): string => {
  if (editor && editor.getHtml) {
    return editor.getHtml()
  }
  return formData.content
}

const setEditorContent = (content: string): void => {
  formData.content = content
  if (editor && editor.setHtml) {
    editor.setHtml(content)
  }
}

const insertHeading = (level: number) => {
  const currentContent = getEditorContent()
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    const html = `<h${level}>${text}</h${level}>`
    setEditorContent(currentContent + html)
  } else {
    setEditorContent(currentContent + `<h${level}>标题</h${level}>\n`)
  }
}

const toggleBold = () => {
  const currentContent = getEditorContent()
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    const html = `<strong>${text}</strong>`
    setEditorContent(currentContent + html)
  } else {
    setEditorContent(currentContent + '<strong>粗体文字</strong>')
  }
}

const toggleItalic = () => {
  const currentContent = getEditorContent()
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    const html = `<em>${text}</em>`
    setEditorContent(currentContent + html)
  } else {
    setEditorContent(currentContent + '<em>斜体文字</em>')
  }
}

const toggleUnderline = () => {
  const currentContent = getEditorContent()
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    const html = `<span style="text-decoration: underline;">${text}</span>`
    setEditorContent(currentContent + html)
  } else {
    setEditorContent(currentContent + '<span style="text-decoration: underline;">下划线文字</span>')
  }
}

const toggleStrikeThrough = () => {
  const currentContent = getEditorContent()
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    const html = `<del>${text}</del>`
    setEditorContent(currentContent + html)
  } else {
    setEditorContent(currentContent + '<del>删除线文字</del>')
  }
}

const changeTextColor = () => {
  const color = prompt('请输入颜色值（如: red, #ff0000, rgb(255,0,0)）:', 'red')
  if (color) {
    const currentContent = getEditorContent()
    const selection = window.getSelection()
    if (selection && selection.toString().length > 0) {
      const text = selection.toString()
      const html = `<span style="color: ${color};">${text}</span>`
      setEditorContent(currentContent + html)
    } else {
      setEditorContent(currentContent + `<span style="color: ${color};">彩色文字</span>`)
    }
  }
}

const changeBgColor = () => {
  const color = prompt('请输入背景颜色值（如: yellow, #ffff00）:', 'yellow')
  if (color) {
    const currentContent = getEditorContent()
    const selection = window.getSelection()
    if (selection && selection.toString().length > 0) {
      const text = selection.toString()
      const html = `<span style="background-color: ${color};">${text}</span>`
      setEditorContent(currentContent + html)
    } else {
      setEditorContent(currentContent + `<span style="background-color: ${color};">背景文字</span>`)
    }
  }
}

const insertUnorderedList = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<ul><li>列表项</li><li>列表项</li></ul>')
}

const insertOrderedList = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<ol><li>列表项</li><li>列表项</li></ol>')
}

const alignLeft = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<p style="text-align: left;">左对齐文字</p>')
}

const alignCenter = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<p style="text-align: center;">居中对齐文字</p>')
}

const alignRight = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<p style="text-align: right;">右对齐文字</p>')
}

const insertLink = () => {
  const url = prompt('请输入链接地址:', 'https://')
  if (url) {
    const currentContent = getEditorContent()
    const selection = window.getSelection()
    if (selection && selection.toString().length > 0) {
      const text = selection.toString()
      const html = `<a href="${url}" target="_blank">${text}</a>`
      setEditorContent(currentContent + html)
    } else {
      setEditorContent(currentContent + '<a href="' + url + '" target="_blank">链接文字</a>')
    }
  }
}

const insertImage = () => {
  const url = prompt('请输入图片地址:', 'https://')
  if (url) {
    const alt = prompt('请输入图片描述:', '图片')
    const currentContent = getEditorContent()
    setEditorContent(currentContent + `<img src="${url}" alt="${alt || '图片'}" />`)
  }
}

const insertCodeBlock = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<pre><code>代码块</code></pre>')
}

const insertTable = () => {
  const currentContent = getEditorContent()
  const html = `<table border="1" cellpadding="5" cellspacing="0" style="border-collapse: collapse;"><thead><tr><th>表头1</th><th>表头2</th></tr></thead><tbody><tr><td>单元格1</td><td>单元格2</td></tr><tr><td>单元格3</td><td>单元格4</td></tr></tbody></table>`
  setEditorContent(currentContent + html)
}

const insertQuote = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<blockquote>引用文字</blockquote>')
}

const insertHorizontalRule = () => {
  const currentContent = getEditorContent()
  setEditorContent(currentContent + '<hr>')
}

const clearStyles = () => {
  const currentContent = getEditorContent()
  // 移除选中内容的所有样式
  const selection = window.getSelection()
  if (selection && selection.toString().length > 0) {
    const text = selection.toString()
    setEditorContent(currentContent + text)
  } else {
    message.info('请选择文字以清除格式')
  }
}

const undo = () => {
  document.execCommand('undo')
}

const redo = () => {
  document.execCommand('redo')
}

const getIndent = (depth: number) => '  '.repeat(depth || 0)

// Simple markdown to HTML converter
const renderedContent = computed(() => {
  if (!formData.content || formData.content.trim() === '') return ''
  let html = formData.content
  
  // Headers
  html = html.replace(/^### (.*)$/gm, '<h3 class="text-xl font-bold mt-6 mb-3">$1</h3>')
  html = html.replace(/^## (.*)$/gm, '<h2 class="text-2xl font-bold mt-6 mb-3">$1</h2>')
  html = html.replace(/^# (.*)$/gm, '<h1 class="text-3xl font-bold mt-6 mb-3">$1</h1>')
  
  // Bold and italic
  html = html.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/\*(.*?)\*/g, '<em>$1</em>')
  html = html.replace(/~~(.*?)~~/g, '<del>$1</del>')
  
  // Links
  html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" class="text-indigo-600 hover:text-indigo-800 underline" target="_blank">$1</a>')
  
  // Images
  html = html.replace(/!\[(.*?)\]\((.*?)\)/g, '<img src="$2" alt="$1" class="max-w-full h-auto rounded-lg my-4" />')
  
  // Code blocks
  html = html.replace(/```([\s\S]*?)```/g, '<pre class="bg-gray-900 text-green-400 rounded-lg p-4 my-4 overflow-x-auto"><code>$1</code></pre>')
  html = html.replace(/`([^`]+)`/g, '<code class="bg-gray-100 text-red-600 px-1.5 py-0.5 rounded text-sm">$1</code>')
  
  // Blockquotes
  html = html.replace(/^> (.*)$/gm, '<blockquote class="border-l-4 border-indigo-500 pl-4 italic text-gray-600 my-4">$1</blockquote>')
  
  // Lists
  html = html.replace(/^- (.*)$/gm, '<li class="ml-4 list-disc">$1</li>')
  html = html.replace(/^\d+\. (.*)$/gm, '<li class="ml-4 list-decimal">$1</li>')
  
  // Horizontal rules
  html = html.replace(/^---$/gm, '<hr class="my-6 border-gray-200">')
  
  // Paragraphs
  html = html.split('\n\n').map(p => {
    if (p.startsWith('<h') || p.startsWith('<blockquote') || p.startsWith('<pre') || p.startsWith('<hr') || p.startsWith('<li')) {
      return p
    }
    if (p.trim()) {
      return `<p class="my-4">${p}</p>`
    }
    return p
  }).join('\n')
  
  return html
})

// Switch editor type
const switchEditorType = (type: 'markdown' | 'rich') => {
  console.log('Switching to editor type:', type)
  
  // Save current content before switching
  const currentContent = formData.content
  
  // Destroy current editor first, then change type
  if (vditor !== null) {
    try {
      vditor.destroy()
    } catch (e) {
      console.warn('Vditor destroy error:', e)
    }
    vditor = null
  }
  if (editor !== null) {
    editor.destroy()
    editor = null
  }
  
  // First switch editor type (Vue will update DOM reactively)
  editorType.value = type
  
  // Then wait for Vue to update the DOM (v-else-if needs time to switch)
  nextTick(() => {
    // Wait another tick for v-else-if to fully render
    nextTick(() => {
      // Restore content - Vue v-model will handle this automatically
      if (type === 'markdown') {
        initVditor()
      }
      // Rich text editor uses v-model, no need to manually set content
    })
  })
}

// Initialize Vditor (Markdown Editor)
const initVditor = () => {
  const el = vditorRef.value
  if (!el) {
    console.warn('Vditor: Invalid element')
    return
  }

  // Destroy existing editor if exists
  if (vditor) {
    try {
      vditor.destroy()
    } catch (e) {
      console.warn('Vditor destroy error:', e)
    }
    vditor = null
  }

  vditor = new Vditor(el, {
    mode: 'sv',
    height: 'auto',
    theme: 'classic',
    icon: 'ant',
    cache: {
      enable: false
    },
    // 使用本地 CDN 资源（通过 Vite 中间件提供）
    cdn: '/vditor',
    lang: 'zh_CN',
    // Vditor v3.10.8 使用回调函数处理事件
    input: (value: string) => {
      formData.content = value
    },
    after: () => {
      // Set initial content
      if (formData.content) {
        vditor?.setValue(formData.content)
      }
    }
  })
}

const loadCategories = async () => {
  try {
    const res = await getCategoryTree()
    categoryTree.value = (res.data as CategoryItem[]) || []
  } catch (error) {
    console.error('加载分类失败:', error)
  }
}

const addTag = () => {
  if (tagInput.value && !formData.tags.includes(tagInput.value)) {
    formData.tags.push(tagInput.value)
  }
  tagInput.value = ''
}

const handleCoverUpload = () => {
  message.info('文件上传功能待实现')
}

const handleSubmit = async () => {
  if (!formData.title) {
    message.error('请输入标题')
    return
  }
  if (!formData.categoryId || formData.categoryId === 0) {
    message.error('请选择分类')
    return
  }
  if (!formData.content) {
    message.error('请输入内容')
    return
  }
  
  try {
    if (isEdit.value && articleId.value) {
      await updateArticle(articleId.value, {
        title: formData.title,
        category_id: formData.categoryId,
        content: formData.content,
        status: formData.status
      })
      message.success('更新成功')
    } else {
      await createArticle({
        title: formData.title,
        category_id: formData.categoryId as number,
        content: formData.content,
        status: formData.status
      })
      message.success('创建成功')
    }
    router.push('/articles')
  } catch (error) {
    message.error('操作失败')
  }
}

const loadArticleData = async (id: number) => {
  try {
    const res = await getArticle(id)
    const articleData = res.data as ArticleItem | null
    if (articleData) {
      article.value = articleData
      formData.title = articleData.title || ''
      formData.categoryId = articleData.category_id || 0
      formData.content = articleData.content || ''
      formData.status = articleData.status ?? 0
      formData.coverUrl = articleData.cover || ''
      if (articleData.tags && Array.isArray(articleData.tags)) {
        formData.tags = articleData.tags.map((tag: any) => tag.name || tag)
      }
      
      // 同步内容到对应的编辑器
      // 等待 DOM 更新后再设置内容
      await nextTick()
      if (editorType.value === 'markdown' && vditor) {
        vditor.setValue(articleData.content || '')
      } else if (editorType.value === 'rich' && editor) {
        // 富文本编辑器通过 setHtml 设置内容
        editor.setHtml(articleData.content || '')
      }
    }
  } catch (error) {
    message.error('加载文章数据失败')
    console.error('加载文章数据失败:', error)
  }
}

onMounted(() => {
  loadCategories()
  // Initialize default editor
  nextTick(() => {
    initVditor()
  })
  // If editing mode, load article data
  if (isEdit.value && articleId.value) {
    loadArticleData(articleId.value)
  }
})

// Watch for route changes
watch(() => route.params.id, (newId) => {
  if (newId && isEdit.value) {
    loadArticleData(Number(newId))
  }
})

onUnmounted(() => {
  if (vditor) {
    vditor.destroy()
    vditor = null
  }
  if (editor) {
    editor.destroy()
    editor = null
  }
})
</script>

<style scoped>
/* Editor container styles */
:deep(.vditor) {
  min-height: 400px;
}

/* wangEditor 容器样式 */
:deep(.wang-editor-container) {
  width: 100%;
}

/* wangEditor 工具栏样式 */
:deep(.w-e-toolbar) {
  border-bottom: 1px solid #e8e8e8 !important;
  background-color: #fff !important;
}

/* wangEditor 编辑器区域样式 */
:deep(.w-e-editor) {
  height: 400px !important;
}

.prose h1 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: #1f2937;
  border-bottom: 2px solid #e5e7eb;
  padding-bottom: 0.75rem;
}

.prose h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #374151;
}

.prose h3 {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #374151;
}

.prose p {
  line-height: 1.8;
  margin-bottom: 1rem;
  color: #4b5563;
}

.prose code {
  background-color: #f3f4f6;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  color: #dc2626;
}

.prose pre {
  background-color: #1f2937;
  color: #10b981;
  padding: 1rem;
  border-radius: 0.5rem;
  overflow-x: auto;
}

.prose blockquote {
  border-left: 4px solid #6366f1;
  padding-left: 1rem;
  margin-left: 0;
  color: #6b7280;
  font-style: italic;
}

.prose img {
  max-width: 100%;
  height: auto;
  border-radius: 0.5rem;
  margin-bottom: 1rem;
}

.prose ul {
  list-style-type: disc;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

.prose ol {
  list-style-type: decimal;
  padding-left: 1.5rem;
  margin-bottom: 1rem;
}

.prose hr {
  border-top: 1px solid #e5e7eb;
  margin: 1.5rem 0;
}
</style>