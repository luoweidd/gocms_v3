import request from './index'

export interface ArticleTag {
  id: number
  name: string
}

export interface ArticleCategory {
  id: number
  name: string
  parent_id?: number
  children?: ArticleCategory[]
}

export interface Article {
  id: number
  title: string
  summary?: string        // 后端字段名，与前端 description 语义相同
  description?: string    // 兼容字段
  content: string
  cover?: string          // 封面图
  category_id: number
  category_name?: string
  author_id: number
  author_name?: string
  status: 0 | 1 | 2 // 0: 草稿，1: 已发布，2: 已下架
  view_count: number
  is_top?: boolean        // 是否置顶
  is_top_value?: number   // 兼容后端 int 类型
  tags?: ArticleTag[]     // 标签列表
  tag_ids?: number[]      // 标签 ID 列表（用于创建/更新）
  published_at?: string   // 发布时间
  created_at: string
  updated_at: string
}

export interface ArticleListParams {
  page: number
  page_size: number
  category_id?: number
  status?: number
  keyword?: string
  author_id?: number
  is_top?: boolean
}

export interface ArticleListResponse {
  list: Article[]
  total: number
  page: number
  page_size: number
}

/**
 * 获取文章列表
 */
export function getArticleList(params: ArticleListParams) {
  return request.get<ArticleListResponse>('/articles', { params })
}

/**
 * 获取文章详情
 */
export function getArticle(id: number) {
  return request.get<Article>(`/articles/${id}`)
}

/**
 * 创建文章
 */
export function createArticle(data: Partial<Article>) {
  return request.post('/articles', data)
}

/**
 * 更新文章
 */
export function updateArticle(id: number, data: Partial<Article>) {
  return request.put(`/articles/${id}`, data)
}

/**
 * 删除文章
 */
export function deleteArticle(id: number) {
  return request.delete(`/articles/${id}`)
}

/**
 * 获取分类列表（普通列表）
 */
export function getArticleCategoryList() {
  return request.get<any>('/article-categories')
}

/**
 * 获取分类树
 */
export function getCategoryTree() {
  return request.get<any>('/article-categories/tree')
}

/**
 * 创建分类
 */
export function createCategory(data: { name: string; parent_id?: number }) {
  return request.post('/article-categories', data)
}

/**
 * 更新分类
 */
export function updateCategory(id: number, data: { name: string; parent_id?: number }) {
  return request.put(`/article-categories/${id}`, data)
}

/**
 * 删除分类
 */
export function deleteCategory(id: number) {
  return request.delete(`/article-categories/${id}`)
}

/**
 * 上传封面图
 */
export function uploadCover(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return request.post<{ url: string }>('/upload/initiate', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}

/**
 * 获取标签列表
 */
export function getTags() {
  return request.get<ArticleTag[]>('/tags')
}
