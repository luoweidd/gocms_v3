package service

import (
	"fmt"

	"gocms_v3/app/db"
	"gocms_v3/app/model"

	"gorm.io/gorm"
)

type MenuService struct{}

func NewMenuService() *MenuService {
	return &MenuService{}
}

// CreateMenu 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	if err := db.GetDB().Create(menu).Error; err != nil {
		return fmt.Errorf("创建菜单失败: %w", err)
	}
	return nil
}

// UpdateMenu 更新菜单
func (s *MenuService) UpdateMenu(id uint, updates map[string]interface{}) error {
	if err := db.GetDB().Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}
	return nil
}

// DeleteMenu 删除菜单（级联删除子菜单）
func (s *MenuService) DeleteMenu(id uint) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		var children []model.Menu
		if err := tx.Where("parent_id = ?", id).Find(&children).Error; err != nil {
			return err
		}
		for _, child := range children {
			if err := s.deleteMenuWithChildren(tx, child.ID); err != nil {
				return err
			}
		}
		return tx.Where("id = ?", id).Delete(&model.Menu{}).Error
	})
}

// deleteMenuWithChildren 递归删除子菜单
func (s *MenuService) deleteMenuWithChildren(dbInst *gorm.DB, menuID uint) error {
	var children []model.Menu
	if err := dbInst.Where("parent_id = ?", menuID).Find(&children).Error; err != nil {
		return err
	}
	for _, child := range children {
		if err := s.deleteMenuWithChildren(dbInst, child.ID); err != nil {
			return err
		}
	}
	return dbInst.Where("id = ?", menuID).Delete(&model.Menu{}).Error
}

// GetMenuByID 根据ID获取菜单
func (s *MenuService) GetMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := db.GetDB().Preload("Children").First(&menu, id).Error; err != nil {
		return nil, fmt.Errorf("菜单不存在")
	}
	return &menu, nil
}

// ListMenus 获取菜单列表
func (s *MenuService) ListMenus(status *int, menuType *int) ([]model.Menu, error) {
	var menus []model.Menu
	query := db.GetDB().Where("deleted_at IS NULL")

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if menuType != nil {
		query = query.Where("menu_type = ?", *menuType)
	}

	if err := query.Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("获取菜单列表失败: %w", err)
	}
	return menus, nil
}

// GetMenuTree 获取菜单树形结构
func (s *MenuService) GetMenuTree() ([]model.MenuTree, error) {
	var menus []model.Menu
	if err := db.GetDB().Where("status = 1 AND deleted_at IS NULL").Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("获取菜单树失败: %w", err)
	}

	menuMap := make(map[uint]*model.MenuTree)
	parentIDMap := make(map[uint]uint) // 记录每个菜单的父节点ID
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
		if menus[i].ParentID != nil && *menus[i].ParentID > 0 {
			tree.ParentID = menus[i].ParentID
			parentIDMap[menus[i].ID] = *menus[i].ParentID
		}
		tree.Children = []model.MenuTree{}
		menuMap[menus[i].ID] = tree
	}

	// 计算深度
	calcDepth := func(id uint) int {
		depth := 0
		current := id
		visited := make(map[uint]bool) // 防止循环引用
		for pid, ok := parentIDMap[current]; ok && pid > 0 && !visited[current]; {
			depth++
			visited[current] = true
			current = pid
			var newPid uint
			newPid, ok = parentIDMap[current]
			if !ok {
				break
			}
			pid = newPid
		}
		return depth
	}

	// 修复：使用指针构建根节点列表，确保子节点的修改能反映到父节点
	var roots []*model.MenuTree
	for _, menu := range menus {
		tree := menuMap[menu.ID]
		tree.Depth = calcDepth(menu.ID)
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
	return result, nil
}

// GetPermissionsByRole 根据角色获取权限列表
func (s *MenuService) GetPermissionsByRole(role string) ([]string, error) {
	var menus []model.Menu
	if err := db.GetDB().Where("status = 1 AND menu_type = 2 AND deleted_at IS NULL").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("获取权限列表失败: %w", err)
	}

	var permissions []string
	for _, menu := range menus {
		permissions = append(permissions, menu.Permission)
	}

	return permissions, nil
}

// GetMenuByPermission 根据权限标识获取菜单
func (s *MenuService) GetMenuByPermission(permission string) (*model.Menu, error) {
	var menu model.Menu
	if err := db.GetDB().Where("permission = ? AND status = 1", permission).First(&menu).Error; err != nil {
		return nil, fmt.Errorf("权限不存在")
	}
	return &menu, nil
}
