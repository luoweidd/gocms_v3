package service

import (
	"fmt"
	"strings"

	"gocms_v3/app/db"
	"gocms_v3/app/model/ngac"

	"gorm.io/gorm"
)

type NgacService struct{}

func NewNgacService() *NgacService {
	return &NgacService{}
}

// ==================== 角色管理 ====================

func (s *NgacService) CreateRole(req ngac.CreateRoleReq) (*ngac.Role, error) {
	role := ngac.Role{
		Name: req.Name,
		Code: req.Code,
		Description: func() string {
			if req.Description != nil {
				return *req.Description
			}
			return ""
		}(),
		ParentID: req.ParentID,
		Sort: func() int {
			if req.Sort != nil {
				return *req.Sort
			}
			return 0
		}(),
		Status: func() int {
			if req.Status != nil && *req.Status != 0 {
				return *req.Status
			}
			return 1
		}(),
		Metadata: func() string {
			if req.Metadata != nil {
				return *req.Metadata
			}
			return ""
		}(),
	}

	if req.ParentID != nil && *req.ParentID > 0 {
		var parent ngac.Role
		if err := db.GetDB().First(&parent, *req.ParentID).Error; err != nil {
			return nil, fmt.Errorf("父角色不存在: %w", err)
		}
		role.Level = parent.Level + 1
		role.Path = parent.Path + "/" + fmt.Sprint(parent.ID)
	} else {
		role.Level = 1
		role.Path = "0"
	}

	if err := db.GetDB().Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *NgacService) ListRoles(req ngac.RoleListReq) ([]ngac.Role, int64, error) {
	var roles []ngac.Role
	var total int64
	query := db.GetDB().Model(&ngac.Role{})
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	query.Count(&total)
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("sort ASC, id DESC").Offset(offset).Limit(req.PageSize).Find(&roles).Error
	return roles, total, err
}

func (s *NgacService) GetRole(id uint) (*ngac.Role, error) {
	var role ngac.Role
	if err := db.GetDB().First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *NgacService) UpdateRole(id uint, req ngac.UpdateRoleReq) error {
	var role ngac.Role
	if err := db.GetDB().First(&role, id).Error; err != nil {
		return err
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
		var parent ngac.Role
		if err := db.GetDB().First(&parent, *req.ParentID).Error; err == nil {
			updates["level"] = parent.Level + 1
			updates["path"] = parent.Path + "/" + fmt.Sprint(parent.ID)
		}
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}
	return db.GetDB().Model(&role).Updates(updates).Error
}

func (s *NgacService) DeleteRole(id uint) error {
	return db.GetDB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", id).Delete(&ngac.RolePermission{}).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&ngac.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Delete(&ngac.Role{}, id).Error
	})
}

func (s *NgacService) BuildRoleTree() ([]ngac.RoleTreeVO, error) {
	var roles []ngac.Role
	if err := db.GetDB().Where("status = 1").Order("sort ASC, id ASC").Find(&roles).Error; err != nil {
		return nil, err
	}
	roleMap := make(map[uint]*ngac.RoleTreeVO)
	var roots []ngac.RoleTreeVO
	for i := range roles {
		vo := ngac.RoleTreeVO{ID: roles[i].ID, Name: roles[i].Name, Code: roles[i].Code, Level: roles[i].Level, Sort: roles[i].Sort, Status: roles[i].Status}
		roleMap[roles[i].ID] = &vo
		if roles[i].ParentID == nil || *roles[i].ParentID == 0 {
			roots = append(roots, vo)
		}
	}
	for i := range roles {
		if roles[i].ParentID != nil && *roles[i].ParentID > 0 {
			if p, ok := roleMap[*roles[i].ParentID]; ok {
				p.Children = append(p.Children, *roleMap[roles[i].ID])
			}
		}
	}
	return roots, nil
}

// ==================== 权限管理 ====================

func (s *NgacService) CreatePermission(req ngac.CreatePermissionReq) (*ngac.Permission, error) {
	perm := ngac.Permission{
		Name: req.Name, Code: req.Code,
		Description: func() string {
			if req.Description != nil {
				return *req.Description
			}
			return ""
		}(),
		ResourceType: req.ResourceType, Action: req.Action,
		SubjectType: func() string {
			if req.SubjectType != nil {
				return *req.SubjectType
			}
			return "user"
		}(),
		Allowed: func() bool {
			if req.Allowed != nil {
				return *req.Allowed
			}
			return true
		}(),
		Priority: func() int {
			if req.Priority != nil {
				return *req.Priority
			}
			return 0
		}(),
		Efficiency: func() string {
			if req.Efficiency != nil {
				return *req.Efficiency
			}
			return "exact"
		}(),
		RuleExpr: func() string {
			if req.RuleExpr != nil {
				return *req.RuleExpr
			}
			return ""
		}(),
		Metadata: func() string {
			if req.Metadata != nil {
				return *req.Metadata
			}
			return ""
		}(),
	}
	if err := db.GetDB().Create(&perm).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (s *NgacService) ListPermissions(req ngac.PermissionListReq) ([]ngac.Permission, int64, error) {
	var perms []ngac.Permission
	var total int64
	query := db.GetDB().Model(&ngac.Permission{})
	if req.ResourceType != "" {
		query = query.Where("resource_type = ?", req.ResourceType)
	}
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}
	query.Count(&total)
	offset := (req.Page - 1) * req.PageSize
	err := query.Order("priority DESC, id DESC").Offset(offset).Limit(req.PageSize).Find(&perms).Error
	return perms, total, err
}

func (s *NgacService) GetPermission(id uint) (*ngac.Permission, error) {
	var perm ngac.Permission
	if err := db.GetDB().First(&perm, id).Error; err != nil {
		return nil, err
	}
	return &perm, nil
}

func (s *NgacService) UpdatePermission(id uint, req ngac.UpdatePermissionReq) error {
	var perm ngac.Permission
	if err := db.GetDB().First(&perm, id).Error; err != nil {
		return err
	}
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.ResourceType != nil {
		updates["resource_type"] = *req.ResourceType
	}
	if req.Action != nil {
		updates["action"] = *req.Action
	}
	if req.SubjectType != nil {
		updates["subject_type"] = *req.SubjectType
	}
	if req.ResourceObjectClassID != nil {
		updates["resource_object_class_id"] = *req.ResourceObjectClassID
	}
	if req.Allowed != nil {
		updates["allowed"] = *req.Allowed
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Efficiency != nil {
		updates["efficiency"] = *req.Efficiency
	}
	if req.RuleExpr != nil {
		updates["rule_expr"] = *req.RuleExpr
	}
	if req.Metadata != nil {
		updates["metadata"] = *req.Metadata
	}
	return db.GetDB().Model(&perm).Updates(updates).Error
}

func (s *NgacService) DeletePermission(id uint) error {
	return db.GetDB().Where("id = ?", id).Delete(&ngac.Permission{}).Error
}

func (s *NgacService) AssignPermissionsToRole(roleID uint, permIDs []uint) error {
	for _, pid := range permIDs {
		rp := ngac.RolePermission{RoleID: roleID, PermissionID: pid}
		if err := db.GetDB().Create(&rp).Error; err != nil && !strings.Contains(err.Error(), "Duplicate") {
			return err
		}
	}
	return nil
}

func (s *NgacService) RemovePermissionFromRole(roleID, permID uint) error {
	return db.GetDB().Where("role_id = ? AND permission_id = ?", roleID, permID).Delete(&ngac.RolePermission{}).Error
}

func (s *NgacService) GetRolePermissions(roleID uint) ([]ngac.Permission, error) {
	var perms []ngac.Permission
	err := db.GetDB().Table("ngac_permissions").Joins("JOIN ngac_role_permissions ON ngac_permissions.id = ngac_role_permissions.permission_id").Where("ngac_role_permissions.role_id = ?", roleID).Find(&perms).Error
	return perms, err
}

// ==================== 用户-角色管理 ====================

func (s *NgacService) AssignRolesToUser(userID uint, roleIDs []uint) error {
	for _, rid := range roleIDs {
		ur := ngac.UserRole{UserID: userID, RoleID: rid}
		if err := db.GetDB().Create(&ur).Error; err != nil && !strings.Contains(err.Error(), "Duplicate") {
			return err
		}
	}
	return nil
}

func (s *NgacService) RemoveRolesFromUser(userID uint, roleIDs []uint) error {
	return db.GetDB().Where("user_id = ? AND role_id IN ?", userID, roleIDs).Delete(&ngac.UserRole{}).Error
}

func (s *NgacService) GetUserRoles(userID uint) ([]ngac.Role, error) {
	var roles []ngac.Role
	err := db.GetDB().Table("ngac_roles").Joins("JOIN ngac_user_roles ON ngac_roles.id = ngac_user_roles.role_id").Where("ngac_user_roles.user_id = ?", userID).Find(&roles).Error
	return roles, err
}

func (s *NgacService) GetUserPermissions(userID uint) ([]ngac.Permission, error) {
	var perms []ngac.Permission
	err := db.GetDB().Table("ngac_permissions").Joins("JOIN ngac_role_permissions ON ngac_permissions.id = ngac_role_permissions.permission_id").Joins("JOIN ngac_user_roles ON ngac_role_permissions.role_id = ngac_user_roles.role_id").Where("ngac_user_roles.user_id = ?", userID).Find(&perms).Error
	return perms, err
}

// ==================== 对象类管理 ====================

func (s *NgacService) CreateObjectClass(req ngac.CreateObjectClassReq) (*ngac.ObjectClass, error) {
	oc := ngac.ObjectClass{
		Name: req.Name, Code: req.Code,
		Description: func() string {
			if req.Description != nil {
				return *req.Description
			}
			return ""
		}(),
		ParentID: req.ParentID,
		Sort: func() int {
			if req.Sort != nil {
				return *req.Sort
			}
			return 0
		}(),
		Status: func() int {
			if req.Status != nil {
				return *req.Status
			}
			return 1
		}(),
		Attributes: func() string {
			if req.Attributes != nil {
				return *req.Attributes
			}
			return ""
		}(),
	}
	if req.ParentID != nil && *req.ParentID > 0 {
		var parent ngac.ObjectClass
		if err := db.GetDB().First(&parent, *req.ParentID).Error; err == nil {
			oc.Level = parent.Level + 1
			oc.Path = parent.Path + "/" + fmt.Sprint(parent.ID)
		}
	} else {
		oc.Level = 1
		oc.Path = "0"
	}
	if err := db.GetDB().Create(&oc).Error; err != nil {
		return nil, err
	}
	return &oc, nil
}

func (s *NgacService) GetObjectClassTree() ([]ngac.ObjectClassTreeVO, error) {
	var ocs []ngac.ObjectClass
	if err := db.GetDB().Where("status = 1").Order("sort ASC, id ASC").Find(&ocs).Error; err != nil {
		return nil, err
	}
	ocMap := make(map[uint]*ngac.ObjectClassTreeVO)
	var roots []ngac.ObjectClassTreeVO
	for i := range ocs {
		vo := ngac.ObjectClassTreeVO{ID: ocs[i].ID, Name: ocs[i].Name, Code: ocs[i].Code, Level: ocs[i].Level, Sort: ocs[i].Sort, Status: ocs[i].Status}
		ocMap[ocs[i].ID] = &vo
		if ocs[i].ParentID == nil || *ocs[i].ParentID == 0 {
			roots = append(roots, vo)
		}
	}
	for i := range ocs {
		if ocs[i].ParentID != nil && *ocs[i].ParentID > 0 {
			if p, ok := ocMap[*ocs[i].ParentID]; ok {
				p.Children = append(p.Children, *ocMap[ocs[i].ID])
			}
		}
	}
	return roots, nil
}

// ==================== 属性管理 ====================

func (s *NgacService) CreateUserAttribute(userID uint, name, value string, attrType *string) error {
	attr := ngac.UserAttribute{UserID: userID, Name: name, Value: value, Type: func() string {
		if attrType != nil {
			return *attrType
		}
		return "string"
	}()}
	return db.GetDB().Create(&attr).Error
}

func (s *NgacService) GetSubjectAttributes(userID uint) map[string]string {
	var attrs []ngac.UserAttribute
	if err := db.GetDB().Where("user_id = ?", userID).Find(&attrs).Error; err != nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("subject.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}

func (s *NgacService) CreateObjectAttribute(name, value, attrType, resourceType string, resourceID, objClassID uint) error {
	attr := ngac.ObjectAttribute{Name: name, Value: value, Type: attrType, ResourceType: resourceType, ResourceID: resourceID, ObjectClassID: objClassID}
	return db.GetDB().Create(&attr).Error
}

func (s *NgacService) GetObjectAttributes(resourceType string, resourceID uint) map[string]string {
	var attrs []ngac.ObjectAttribute
	if err := db.GetDB().Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Find(&attrs).Error; err != nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("object.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}

func (s *NgacService) CreateContextAttribute(req ngac.CreateContextAttributeReq) (*ngac.ContextAttribute, error) {
	ctxAttr := ngac.ContextAttribute{Name: req.Name, Value: req.Value}
	if err := db.GetDB().Create(&ctxAttr).Error; err != nil {
		return nil, err
	}
	return &ctxAttr, nil
}

func (s *NgacService) GetContextAttributes() map[string]string {
	var attrs []ngac.ContextAttribute
	if err := db.GetDB().Where("is_global = true").Find(&attrs).Error; err != nil {
		return make(map[string]string)
	}
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("context.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}
