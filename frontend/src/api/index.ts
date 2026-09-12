import axios, { AxiosInstance, AxiosResponse } from 'axios'
import { ElMessage } from '@/utils/message'

// ==================== Axios 实例配置 ====================

const service: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// ==================== 请求拦截器 ====================

service.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// ==================== 响应拦截器 ====================

service.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data
    // P1 修复：统一成功标准为 code === 0（与后端 response.Success 保持一致）
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    return res
  },
  (error) => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          ElMessage.error('未授权，请重新登录')
          localStorage.removeItem('token')
          localStorage.removeItem('user')
          window.location.href = '/login'
          break
        case 403:
          ElMessage.error('拒绝访问')
          break
        case 404:
          ElMessage.error('请求资源不存在')
          break
        case 500:
          ElMessage.error('服务器错误')
          break
        default:
          ElMessage.error(`请求失败：${error.response.status}`)
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default service

// ==================== API 模块统一导出 ====================

/**
 * Auth API - 认证相关
 * - login, register, getCurrentUser, logout
 * - LoginParams, RegisterParams, UserInfo
 */
export {
  login,
  register,
  getCurrentUser,
  logout,
} from './auth'

export type { LoginParams, RegisterParams, UserInfo } from './auth'

/**
 * User API - 用户管理
 */
export {
  getUserList,
  getUser,
  createUser,
  updateUser,
  deleteUser,
  assignRoles,
  getUserRoles,
} from './user'

export type {
  UserInfo as UserItem,
  CreateUserParams,
  UpdateUserParams,
  UserListParams,
  UserListResponse,
} from './user'

/**
 * Role API - 角色管理
 */
export {
  getRoleList,
  getRole,
  createRole,
  updateRole,
  deleteRole,
  getPermissionTree,
  assignPermissions,
  getRolePermissions,
} from './role'

export type {
  RoleInfo,
  RoleListParams,
  RoleListResponse,
  CreateRoleParams,
  UpdateRoleParams,
  PermissionNode,
} from './role'

/**
 * Menu API - 菜单管理
 */
export {
  getMenuList,
  getMenuTree,
  createMenu,
  updateMenu,
  deleteMenu,
} from './menu'

export type { MenuItem, MenuListParams, MenuTreeResponse } from './menu'

/**
 * Article API - 文章管理
 */
export {
  getArticleList,
  getArticle,
  createArticle,
  updateArticle,
  deleteArticle,
  getCategoryTree,
  createCategory,
  updateCategory,
  deleteCategory,
} from './article'

export type { Article, ArticleListParams, ArticleListResponse } from './article'

/**
 * Video API - 视频管理
 */
export {
  getVideoList,
  getVideo,
  createVideo,
  updateVideo,
  deleteVideo,
  getVideoCategoryTree,
  createVideoCategory,
  updateVideoCategory,
  deleteVideoCategory,
} from './video'

export type { Video, VideoListParams, VideoListResponse } from './video'

/**
 * Comment API - 评论管理
 */
export {
  getCommentList,
  getComment,
  approveComment,
  deleteComment,
  batchDeleteComments,
} from './comment'

export type { CommentInfo, CommentListParams, CommentListResponse } from './comment'

/**
 * Dashboard API - 仪表盘
 */
export {
  getDashboardStats,
  getStatTrend,
  getRecentActivities,
  getDashboardOverview,
} from './dashboard'

export type { DashboardStats, StatTrend, ActivityItem } from './dashboard'