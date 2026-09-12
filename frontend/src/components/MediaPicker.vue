<template>
  <a-modal
    v-model:visible="visible"
    title="选择媒体文件"
    :width="900"
    @ok="confirmSelect"
    @cancel="cancelSelect"
  >
    <!-- 搜索栏 -->
    <div class="search-bar">
      <a-input-search
        v-model:value="searchParams.keyword"
        placeholder="搜索媒体文件..."
        search-button
        @search="fetchMediaList"
        allow-clear
      />
      <a-select
        v-model:value="searchParams.type"
        placeholder="文件类型"
        style="width: 120px"
        allow-clear
        @change="fetchMediaList"
      >
        <a-select-option value="">全部</a-select-option>
        <a-select-option value="image">图片</a-select-option>
        <a-select-option value="video">视频</a-select-option>
        <a-select-option value="file">文件</a-select-option>
      </a-select>
      <a-select
        v-model:value="searchParams.sortBy"
        placeholder="排序方式"
        style="width: 120px"
        @change="fetchMediaList"
      >
        <a-select-option value="created_at">最新上传</a-select-option>
        <a-select-option value="updated_at">最近修改</a-select-option>
        <a-select-option value="name">名称</a-select-option>
        <a-select-option value="size">大小</a-select-option>
      </a-select>
    </div>

    <!-- 视图切换 -->
    <div class="view-toggle">
      <a-button-group>
        <a-button :type="viewMode === 'grid' ? 'primary' : 'default'" @click="viewMode = 'grid'">
          <i class="anticon anticon-appstore"></i>
        </a-button>
        <a-button :type="viewMode === 'list' ? 'primary' : 'default'" @click="viewMode = 'list'">
          <i class="anticon anticon-list"></i>
        </a-button>
      </a-button-group>
      
      <span class="upload-btn">
        <a-button size="small" type="primary" @click="showUploadModal = true">
          上传文件
        </a-button>
      </span>
    </div>

    <!-- 网格视图 -->
    <div v-if="viewMode === 'grid'" class="media-grid">
      <div
        v-for="item in mediaList"
        :key="item.id"
        class="media-item grid-item"
        :class="{ selected: isSelected(item) }"
        @click="toggleSelect(item)"
        @dblclick="handleDoubleClick(item)"
      >
        <div class="media-preview">
          <img
            v-if="isImage(item)"
            :src="item.url || item.path"
            :alt="item.name"
            @error="handleImageError"
          />
          <video
            v-else-if="isVideo(item)"
            :src="item.url || item.path"
            muted
          ></video>
          <div v-else class="file-icon">
            <i :class="getFileIcon(item.type)"></i>
          </div>
        </div>
        <div class="media-info">
          <p class="media-name" :title="item.name">{{ item.name }}</p>
          <span class="media-size">{{ formatSize(item.size) }}</span>
        </div>
        <div v-if="isSelected(item)" class="selected-badge">
          <i class="anticon anticon-check"></i>
        </div>
      </div>
      
      <!-- 空状态 -->
      <a-empty v-if="mediaList.length === 0" description="暂无媒体文件" />
    </div>

    <!-- 列表视图 -->
    <div v-else class="media-list">
      <a-table
        :columns="listColumns"
        :data-source="mediaList"
        :loading="loading"
        row-key="id"
        :row-selection="rowSelection"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'name'">
            <div class="file-name-cell">
              <i :class="getFileIcon(record.type)" class="file-icon-sm"></i>
              <span :title="record.name">{{ record.name }}</span>
            </div>
          </template>
          <template v-if="column.dataIndex === 'preview'">
            <a-image
              v-if="isImage(record)"
              :src="record.url || record.path"
              :width="40"
              :height="40"
              style="object-fit: cover; border-radius: 4px;"
            />
            <span v-else class="preview-text">无预览</span>
          </template>
        </template>
      </a-table>
    </div>

    <!-- 分页 -->
    <a-pagination
      v-model:current="pagination.current"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      show-size-changer
      @change="handlePageChange"
      style="margin-top: 16px; justify-content: flex-end;"
    />

    <!-- 上传模态框 -->
    <a-modal
      v-model:visible="showUploadModal"
      title="上传文件"
      :footer="null"
      width="500px"
    >
      <a-upload
        v-model:file-list="uploadFileList"
        :multiple="true"
        :before-upload="beforeUpload"
        :remove="handleRemove"
        list-type="picture-card"
      >
        <div>
          <i class="anticon anticon-plus"></i>
          <div style="margin-top: 8px">上传</div>
        </div>
      </a-upload>
    </a-modal>
  </a-modal>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { message } from 'ant-design-vue'
import type { UploadFile } from 'ant-design-vue/es/upload/interface'
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { getMediaAssetList, getAlbumList } from '@/api/media'

interface MediaItem {
  id: number
  name: string
  type: string
  url: string
  path: string
  size?: number
  album_id?: number
  created_at: string
  updated_at: string
}

interface Props {
  visible: boolean
  multiple?: boolean
  acceptTypes?: string
}

const props = withDefaults(defineProps<Props>(), {
  multiple: false,
  acceptTypes: 'image/*,video/*,.pdf,.doc,.docx,.xls,.xlsx'
})

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'select', value: MediaItem[]): void
}>()

// 状态
const visible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const mediaList = ref<MediaItem[]>([])
const loading = ref(false)
const viewMode = ref<'grid' | 'list'>('grid')
const selectedItems = ref<Set<number>>(new Set())
const showUploadModal = ref(false)
const uploadFileList = ref<UploadFile[]>([])

// 搜索参数
const searchParams = ref({
  keyword: '',
  type: '' as string | undefined,
  sortBy: 'created_at' as string
})

// 分页
const pagination = ref({
  current: 1,
  pageSize: 24,
  total: 0
})

// 列表列定义
const listColumns = [
  { title: '预览', dataIndex: 'preview', width: 80 },
  { title: '名称', dataIndex: 'name', ellipsis: true },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: '大小', dataIndex: 'size', width: 100, customRender: ({ record }: any) => formatSize(record.size) },
  { title: '上传时间', dataIndex: 'created_at', width: 180 }
]

// 表格选择配置
const rowSelection = computed(() => ({
  selectedRowKeys: Array.from(selectedItems.value),
  onChange: (selectedRowKeys: any[]) => {
    selectedItems.value = new Set(selectedRowKeys)
  }
}))

// 获取媒体列表
const fetchMediaList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.current,
      page_size: pagination.value.pageSize,
      keyword: searchParams.value.keyword,
      type: searchParams.value.type,
      sort_by: searchParams.value.sortBy
    }
    
    // TODO: 调用后端API
    // const res = await getMediaList(params)
    // mediaList.value = res.data.records
    // pagination.value.total = res.data.total
    
    // 模拟数据（待替换为真实API）
    mediaList.value = []
    pagination.value.total = 0
  } catch (error) {
    message.error('获取媒体列表失败')
  } finally {
    loading.value = false
  }
}

// 判断文件类型
const isImage = (item: MediaItem) => item.type?.startsWith('image/')
const isVideo = (item: MediaItem) => item.type?.startsWith('video/')

// 获取文件图标
const getFileIcon = (type?: string) => {
  if (!type) return 'anticon anticon-file'
  if (type.startsWith('image/')) return 'anticon anticon-file-image'
  if (type.startsWith('video/')) return 'anticon anticon-file-video'
  if (type.includes('pdf')) return 'anticon anticon-file-pdf'
  if (type.includes('word') || type.includes('document')) return 'anticon anticon-file-excel'
  return 'anticon anticon-file'
}

// 格式化文件大小
const formatSize = (bytes?: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// 判断是否选中
const isSelected = (item: MediaItem) => selectedItems.value.has(item.id)

// 切换选中状态
const toggleSelect = (item: MediaItem) => {
  if (props.multiple) {
    if (selectedItems.value.has(item.id)) {
      selectedItems.value.delete(item.id)
    } else {
      selectedItems.value.add(item.id)
    }
  } else {
    selectedItems.value.clear()
    selectedItems.value.add(item.id)
  }
}

// 双击选中
const handleDoubleClick = (item: MediaItem) => {
  if (!props.multiple) {
    selectedItems.value.clear()
    selectedItems.value.add(item.id)
  }
  confirmSelect()
}

// 确认选择
const confirmSelect = () => {
  const selected = mediaList.value.filter(item => selectedItems.value.has(item.id))
  emit('select', selected)
  visible.value = false
}

// 取消选择
const cancelSelect = () => {
  selectedItems.value.clear()
  visible.value = false
}

// 图片加载失败处理
const handleImageError = (event: Event) => {
  const target = event.target as HTMLImageElement
  target.setAttribute('src', '')
}

// 分页变化
const handlePageChange = (page: number, pageSize: number) => {
  pagination.value.current = page
  pagination.value.pageSize = pageSize
  fetchMediaList()
}

// 表格变化
const handleTableChange = (pag: any) => {
  pagination.value.current = pag.current
  pagination.value.pageSize = pag.pageSize
}

// 上传前检查
const beforeUpload = (file: File) => {
  // TODO: 实现上传逻辑
  message.info(`正在上传: ${file.name}`)
  return false // 阻止默认上传
}

// 移除文件
const handleRemove = (file: UploadFile) => {
  uploadFileList.value = uploadFileList.value.filter(f => f.uid !== file.uid)
}

// 监听visible变化
watch(() => props.visible, (val) => {
  if (val) {
    selectedItems.value.clear()
    fetchMediaList()
  }
})
</script>

<style scoped>
.search-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.view-toggle {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.upload-btn button {
  background-color: #1890ff;
  color: #fff;
}

.media-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
  max-height: 500px;
  overflow-y: auto;
  padding: 8px;
}

.media-item {
  position: relative;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background-color: #fff;
}

.media-item:hover {
  border-color: #1890ff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.media-item.selected {
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24, 144, 255, 0.2);
}

.grid-item {
  display: flex;
  flex-direction: column;
}

.media-preview {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
}

.media-preview img,
.media-preview video {
  max-width: 100%;
  max-height: 100%;
  object-fit: cover;
}

.file-icon {
  font-size: 48px;
  color: #999;
}

.media-info {
  padding: 8px;
}

.media-name {
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 12px;
  color: #333;
}

.media-size {
  font-size: 11px;
  color: #999;
}

.selected-badge {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background-color: #1890ff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
}

.media-list {
  max-height: 500px;
  overflow-y: auto;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon-sm {
  font-size: 16px;
  color: #999;
}

.preview-text {
  color: #999;
  font-size: 12px;
}

::v-deep(.ant-modal-body) {
  max-height: 600px;
  overflow-y: auto;
}
</style>