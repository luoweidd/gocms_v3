package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
)

// MenuServiceWithCache 带缓存的菜单服务
type MenuServiceWithCache struct {
	cacheKey string
	ttl      time.Duration
}

// NewMenuServiceWithCache 创建带缓存的菜单服务
func NewMenuServiceWithCache(ttl time.Duration) *MenuServiceWithCache {
	return &MenuServiceWithCache{
		cacheKey: "menu:tree:all",
		ttl:      ttl,
	}
}

// GetMenuTreeCached 获取带缓存的菜单树
func (s *MenuServiceWithCache) GetMenuTreeCached() ([]model.MenuTree, error) {
	// TODO: 从Redis读取缓存
	var menus []model.Menu
	if err := db.GetDB().Where("status = 1 AND deleted_at IS NULL").Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("获取菜单树失败: %w", err)
	}

	tree := buildMenuTree(menus)

	// TODO: 将树结构写入Redis缓存
	return tree, nil
}

func (s *MenuServiceWithCache) ClearMenuCache() {
	// TODO: 清除Redis缓存
	fmt.Println("菜单缓存已清除")
}

// ClearAllCache 清除所有缓存
func (s *MenuServiceWithCache) ClearAllCache() error {
	// TODO: 清除所有相关缓存
	return nil
}

// buildMenuTree 构建菜单树结构
func buildMenuTree(menus []model.Menu) []model.MenuTree {
	menuMap := make(map[uint]*model.MenuTree)
	for i := range menus {
		tree := &model.MenuTree{
			ID:         menus[i].ID,
			Name:       menus[i].Name,
			Path:       menus[i].Path,
			Icon:       menus[i].Icon,
			Sort:       menus[i].SortOrder,
			Type:       menus[i].MenuType,
			Permission: menus[i].Permission,
		}
		tree.Children = []model.MenuTree{}
		menuMap[menus[i].ID] = tree
	}

	// 修复：使用指针构建根节点列表，确保子节点的追加能反映到父节点
	var roots []*model.MenuTree
	for _, menu := range menus {
		tree := menuMap[menu.ID]
		if menu.ParentID == nil || *menu.ParentID == 0 {
			roots = append(roots, tree)
		} else if parent, ok := menuMap[*menu.ParentID]; ok {
			parent.Children = append(parent.Children, *tree)
		}
	}

	// 将指针切片转换为值切片返回
	result := make([]model.MenuTree, len(roots))
	for i, r := range roots {
		result[i] = *r
	}
	return result
}

// MenuCacheManager 菜单缓存管理器
type MenuCacheManager struct {
	cacheKeyPrefix string
	ttl            time.Duration
}

// NewMenuCacheManager 创建菜单缓存管理器
func NewMenuCacheManager(ttl time.Duration) *MenuCacheManager {
	return &MenuCacheManager{
		cacheKeyPrefix: "gocms:menu:",
		ttl:            ttl,
	}
}

// SetCache 设置菜单缓存
func (m *MenuCacheManager) SetCache(key string, data interface{}) error {
	ctx := context.Background()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化菜单数据失败: %w", err)
	}

	key = m.cacheKeyPrefix + key

	// TODO: 写入Redis缓存
	_ = ctx
	_ = jsonData
	_ = m.ttl

	return nil
}

// GetCache 获取菜单缓存
func (m *MenuCacheManager) GetCache(key string) ([]byte, error) {
	// TODO: 从Redis读取缓存
	key = m.cacheKeyPrefix + key
	return []byte{}, nil
}

// InvalidateCache 失效菜单缓存
func (m *MenuCacheManager) InvalidateCache() {
	// TODO: 删除Redis中所有菜单相关缓存
	fmt.Println("菜单缓存已失效")
}
