import request from './index'

export interface MenuItem {
  id: number
  name: string
  icon?: string
  path?: string
  parent_id?: number
  sort: number
  type: number          // 后端返回：1=菜单 2=目录 3=外链
  component?: string
  permission?: string
  visible: number
  status: number
  title?: string        // 菜单显示名称（用于前端显示）
  depth?: number        // 层级深度
  meta?: string         // 元数据（JSON字符串）
  children?: MenuItem[]
  created_at: string
  updated_at: string
}

export interface MenuListParams {
  page: number
  page_size: number
  keyword?: string
  type?: string
  status?: number
}

export interface MenuTreeResponse {
  list: MenuItem[]
  total: number
}

/**
 * 获取菜单列表
 */
export function getMenuList(params: MenuListParams) {
  return request.get<MenuTreeResponse>('/menus', { params })
}

/**
 * 获取菜单树
 */
export function getMenuTree() {
  return request.get<MenuItem[]>('/menus/tree')
}

/**
 * 创建菜单
 */
export function createMenu(data: Partial<MenuItem>) {
  return request.post('/menus', data)
}

/**
 * 更新菜单
 */
export function updateMenu(id: number, data: Partial<MenuItem>) {
  return request.put(`/menus/${id}`, data)
}

/**
 * 删除菜单
 */
export function deleteMenu(id: number) {
  return request.delete(`/menus/${id}`)
}
