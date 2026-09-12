package model

// Menu 菜单模型（对应 menus 表）
type Menu struct {
	BaseModel
	Name       string `gorm:"uniqueIndex:uk_menu_name;size:100;not null" json:"name"`  // 菜单名称
	NameEn     string `gorm:"size:100" json:"name_en"`                                 // 英文名称
	Path       string `gorm:"size:255" json:"path"`                                    // 路由路径
	Icon       string `gorm:"size:100" json:"icon"`                                    // 菜单图标
	Component  string `gorm:"size:255" json:"component"`                               // 前端组件路径
	SortOrder  int    `gorm:"column:sort_order;default:0" json:"sort_order"`           // 排序号
	ParentID   *uint  `gorm:"index:idx_menu_parent;column:parent_id" json:"parent_id"` // 父级菜单ID (NULL表示顶级菜单)
	MenuType   int    `gorm:"column:menu_type;default:1" json:"menu_type"`             // 菜单类型：1菜单 2目录 3外链
	Visible    int    `gorm:"default:1" json:"visible"`                                // 是否可见：1可见 0隐藏
	Permission string `gorm:"size:100" json:"permission"`                              // 权限标识
	Status     int    `gorm:"default:1;index:idx_menu_status" json:"status"`           // 状态：1启用 0禁用
	Children   []Menu `gorm:"foreignKey:ParentID;references:ID" json:"children"`       // 子菜单
}

// TableName 指定表名
func (Menu) TableName() string {
	return "menus"
}

// MenuTree 菜单树结构
type MenuTree struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	Icon       string     `json:"icon"`
	Sort       int        `json:"sort"`
	Type       int        `json:"type"`
	Permission string     `json:"permission"`
	Title      string     `json:"title"`
	Meta       string     `json:"meta"`
	ParentID   *uint      `json:"parent_id"`       // 父菜单ID（用于前端级联选择）
	Depth      int        `json:"depth,omitempty"` // 层级深度（用于前端缩进显示）
	Children   []MenuTree `json:"children"`
}
