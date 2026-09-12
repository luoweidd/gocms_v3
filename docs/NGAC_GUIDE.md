# GoCMS v3 NGAC 权限管理模块完整指南

**文档版本**: v2.0  
**最后更新**: 2026-08-30  
**适用版本**: GoCMS v3

---

## 目录

1. [概述](#概述)
2. [核心架构](#核心架构)
3. [数据模型](#数据模型)
4. [服务层说明](#服务层说明)
5. [ABAC 引擎](#abac-引擎)
6. [PDP 决策点](#pdp-决策点)
7. [API 端点参考](#api-端点参考)
8. [使用示例](#使用示例)
9. [扩展指南](#扩展指南)

---

## 概述

NGAC (New Grid Access Control) 是一个基于属性的访问控制（ABAC）扩展模型，在 ABAC 基础上增加了角色层和对象类层次结构，使系统更加灵活和可扩展。

### NGAC 核心特性

1. **角色层** - 在 ABAC 属性基础上增加角色抽象，简化权限管理
2. **对象类层次** - 定义资源分类体系，支持继承
3. **策略表达式** - 支持复杂的 ABAC 条件表达式
4. **决策点(PDP)** - 集中式策略决策引擎
5. **访问审计** - 完整的访问日志记录

### NGAC 与 ABAC 的关系

```
NGAC = ABAC + Role Layer + Object Class Hierarchy

ABAC: 属性 -> 策略评估 -> 授权决策
    ↓
NGAC: 角色 -> 属性 -> 策略评估 -> 授权决策
           ↑
        对象类层次结构
```

---

## 核心架构

### NGAC 核心组件

| 组件 | 说明 |
|------|------|
| PDP (Policy Decision Point) | 策略决策点，做出访问控制决策 |
| PEP (Policy Enforcement Point) | 策略执行点，执行 PDP 的决策 |
| PIP (Policy Information Point) | 策略信息点，提供属性信息 |
| Axiomatic PDP | 基础 PDP，使用固有属性进行决策 |
| Obligations PDP | 义务 PDP，附带操作义务的决策 |

### NGAC 授权流程

1. **Subject（主体）** - 发起请求的用户/服务
2. **PIF (Policy Information Flow)** - 收集主体、客体、环境属性
3. **PDP 评估** - 使用 ABAC 引擎评估策略表达式
4. **Access Decision** - 返回 allow/deny 决策
5. **Audit Log** - 记录访问日志

---

## 数据模型

### 核心表结构

| 表名 | 说明 |
|------|------|
| `ngac_roles` | 角色表（支持层次结构） |
| `ngac_permissions` | 权限/规则表 |
| `ngac_role_permissions` | 角色-权限关联表 |
| `ngac_user_roles` | 用户-角色关联表 |
| `ngac_object_classes` | 对象类表（支持层次结构） |
| `ngac_object_attributes` | 资源对象属性表 |
| `ngac_user_attributes` | 用户属性表 |
| `ngac_context_attributes` | 上下文/环境属性表 |
| `ngac_access_logs` | 访问审计日志表 |
| `ngac_pdp_config` | PDP引擎配置表 |

### 文件路径

| 组件 | 文件路径 | 说明 |
|------|----------|------|
| 数据模型 | `app/model/ngac/models.go` | 10个核心表定义 |
| NGAC 服务 | `app/service/ngac_service.go` | CRUD + 关联管理 |
| ABAC 引擎 | `app/service/abac_engine.go` | 属性评估引擎 |
| PDP 决策点 | `app/service/pdp_engine.go` | 策略决策引擎 |
| 控制器 | `app/controller/ngac_controller.go` | RESTful API |

---

## 服务层说明

### NgacService 主要方法

#### 角色管理

| 方法 | 说明 |
|------|------|
| `CreateRole(req)` | 创建角色 |
| `UpdateRole(id, req)` | 更新角色 |
| `DeleteRole(id)` | 删除角色 |
| `GetRole(id)` | 获取角色详情 |
| `ListRoles(req)` | 分页查询角色列表 |
| `BuildRoleTree()` | 构建角色树 |

#### 权限管理

| 方法 | 说明 |
|------|------|
| `CreatePermission(req)` | 创建权限规则 |
| `UpdatePermission(id, req)` | 更新权限规则 |
| `DeletePermission(id)` | 删除权限规则 |
| `GetPermission(id)` | 获取权限详情 |
| `ListPermissions(req)` | 查询权限列表 |

#### 角色-权限关联

| 方法 | 说明 |
|------|------|
| `AssignPermissionsToRole(roleID, permIDs)` | 分配权限给角色 |
| `RemovePermissionFromRole(roleID, permID)` | 从角色移除权限 |
| `GetRolePermissions(roleID)` | 获取角色的权限列表 |

#### 用户-角色管理

| 方法 | 说明 |
|------|------|
| `AssignRolesToUser(userID, roleIDs)` | 分配角色给用户 |
| `RemoveRolesFromUser(userID, roleIDs)` | 从用户移除角色 |
| `GetUserRoles(userID)` | 获取用户角色列表 |
| `GetUserPermissions(userID)` | 获取用户权限列表 |

#### 对象类管理

| 方法 | 说明 |
|------|------|
| `CreateObjectClass(req)` | 创建对象类 |
| `GetObjectClassTree()` | 获取对象类树 |

#### 属性管理

| 方法 | 说明 |
|------|------|
| `CreateUserAttribute(userID, name, value, type)` | 创建用户属性 |
| `GetSubjectAttributes(userID)` | 获取用户属性 |
| `CreateObjectAttribute(...)` | 创建对象属性 |
| `GetObjectAttributes(resourceType, resourceID)` | 获取对象属性 |
| `CreateContextAttribute(req)` | 创建上下文属性 |
| `GetContextAttributes()` | 获取上下文属性 |

---

## ABAC 引擎

### 功能说明

ABAC 引擎实现基于属性的访问控制策略评估：

- **Evaluate(ruleExpr, subjectAttrs, objectAttrs)** - ABAC 策略评估
- 支持逻辑运算符: AND, OR, NOT
- 支持比较运算符: ==, !=, >=, <=, >, <, =~(正则), ^(前缀匹配), contains, startswith, endswith

### 策略表达式示例

```go
// 条件1 AND (条件2 OR 条件3)
"{subject.department} == 'IT' AND ({subject.level} >= 5 OR {object.visibility} == 'public')}"
```

---

## PDP 决策点

### 功能说明

PDP (Policy Decision Point) 实现集中式策略决策：

- **EvaluateAuthz(req)** - 评估授权请求
- **InitializePDP(cfg)** - 初始化 PDP 配置
- 记录访问日志

### PDP 配置

```go
service.InitializePDP(map[string]interface{}{
    "log_enabled": true,
    "default_decision": "deny",
})
```

---

## API 端点参考

### 角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/roles` | 创建角色 |
| GET | `/api/ngac/roles` | 角色列表 |
| GET | `/api/ngac/roles/:id` | 角色详情 |
| PUT | `/api/ngac/roles/:id` | 更新角色 |
| DELETE | `/api/ngac/roles/:id` | 删除角色 |
| GET | `/api/ngac/roles/tree` | 构建角色树 |

### 权限管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/permissions` | 创建权限规则 |
| GET | `/api/ngac/permissions` | 权限列表 |
| GET | `/api/ngac/permissions/:id` | 权限详情 |
| PUT | `/api/ngac/permissions/:id` | 更新权限规则 |
| DELETE | `/api/ngac/permissions/:id` | 删除权限规则 |

### 角色-权限关联

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/roles/:role_id/permissions` | 分配权限给角色 |
| DELETE | `/api/ngac/roles/:role_id/permissions/:perm_id` | 从角色移除权限 |

### 用户-角色管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/users/:user_id/roles` | 分配角色给用户 |
| GET | `/api/ngac/users/:user_id/roles` | 获取用户角色列表 |
| GET | `/api/ngac/users/:user_id/permissions` | 获取用户权限列表 |

### 授权评估

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/authz` | 授权请求（PDP） |

### 对象类管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/object-classes` | 创建对象类 |
| GET | `/api/ngac/object-classes/tree` | 获取对象类树 |

### 用户属性管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/users/:user_id/attributes` | 创建用户属性 |
| GET | `/api/ngac/users/:user_id/attributes` | 获取用户属性 |

### 对象属性管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/object-attributes` | 创建对象属性 |
| GET | `/api/ngac/object-attributes` | 获取对象属性 |

### 上下文属性管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/ngac/context-attributes` | 创建上下文属性 |
| GET | `/api/ngac/context-attributes` | 获取上下文属性 |

---

## 使用示例

### 1. 初始化 PDP 配置

```go
service.InitializePDP(map[string]interface{}{
    "log_enabled": true,
    "default_decision": "deny",
})
```

### 2. 创建角色并分配权限

```go
// 创建根角色（如管理员）
role := ngac.CreateRoleReq{
    Name:        "admin",
    Code:        "admin",
    Description: ptrStr("系统管理员"),
}
svc.CreateRole(role)

// 创建权限规则
perm := ngac.CreatePermissionReq{
    Name:         "article:create",
    Code:         "article:create",
    ResourceType: "article",
    Action:       "create",
    Priority:     ptrInt(10),
}
svc.CreatePermission(perm)

// 分配权限给角色
svc.AssignPermissionsToRole(role.ID, []uint{perm.ID})

// 给用户分配角色
svc.AssignRolesToUser(userID, []uint{role.ID})
```

### 3. 用户授权评估

```go
resp := service.EvaluateAuthz(ngac.AuthzRequest{
    UserID:       1,
    ResourceType: "article",
    ResourceID:   100,
    Action:       "create",
})
// resp.Allowed = true/false
```

---

## 扩展指南

### 如何添加新的资源类型

1. 在 `ngac_object_classes` 表中添加新对象类
2. 为对象类定义属性
3. 创建对应的权限规则
4. 将角色与权限关联

### 如何自定义策略表达式

NGAC 支持以下运算符：

| 运算符类型 | 符号 | 说明 |
|-----------|------|------|
| 逻辑运算符 | AND, OR, NOT | 组合条件 |
| 比较运算符 | ==, !=, >=, <=, >, < | 值比较 |
| 模式匹配 | =~ (正则), ^ (前缀) | 模式匹配 |
| 字符串操作 | contains, startswith, endswith | 字符串检查 |

### NGAC 数据库迁移

系统启动时自动同步表结构：

```bash
go run cmd/server/main.go
```

---

*文档版本：v2.0*  
*最后更新：2026-08-30*