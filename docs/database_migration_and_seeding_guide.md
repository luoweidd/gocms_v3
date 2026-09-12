# GoCMS V3 - 数据库迁移与初始化指南

## 概述

本文档详细说明 GoCMS V3 的数据库结构、迁移流程和种子数据初始化机制。

## 架构设计原则

### 迁移与种子数据分离

- **SQL 迁移文件**：仅负责数据库表结构的创建和修改
- **Go Seeder**：负责初始化业务数据（角色、菜单、配置等）

这种分离设计带来以下优势：
1. 清晰的职责划分，便于维护
2. 种子数据可以使用 Go 代码逻辑，更灵活
3. 迁移文件可以独立执行，支持数据库版本管理

## 迁移文件列表

### 001_init_core_tables.sql - 核心表结构
创建系统核心功能所需的基础表：

| 表名 | 描述 |
|------|------|
| `tenants` | 租户表（多租户支持） |
| `users` | 管理员用户表（含 OAuth2 字段） |
| `comments` | 评论表 |
| `file_uploads` | 文件分片上传表 |
| `oauth_auth_codes` | OAuth 授权码表 |
| `api_sign_keys` | API 签名密钥表 |

### 002_init_system_tables.sql - 系统表结构
创建 CMS 系统业务表：

| 表名 | 描述 |
|------|------|
| `article_categories` | 文章分类表 |
| `video_categories` | 视频分类表 |
| `articles` | 文章内容表 |
| `videos` | 视频内容表 |
| `menus` | 菜单表 |
| `roles` | 角色表 |
| `user_roles` | 用户-角色关联表 |
| `system_configs` | 系统配置表 |
| `data_dicts` | 数据字典表 |
| `operation_logs` | 操作日志表 |

### 003_init_ngac_and_seed.sql - NGAC 权限表 + 种子数据
创建 NGAC（Next-generation Access Control）权限模型并提供初始数据：

| 表名 | 描述 |
|------|------|
| `ngac_roles` | NGAC 角色表（层级结构） |
| `ngac_permissions` | NGAC 权限表 |
| `ngac_role_permissions` | 角色-权限关联表 |
| `ngac_user_roles` | 用户-角色关联表 |
| `ngac_object_classes` | 对象类表 |
| `ngac_object_attributes` | 对象属性表 |
| `ngac_user_attributes` | 用户属性表 |
| `ngac_context_attributes` | 上下文属性表 |
| `ngac_access_logs` | 访问日志表 |
| `ngac_pdp_config` | PDP 引擎配置表 |

## 数据库 ER 图

```
┌──────────┐     ┌──────────────┐     ┌─────────────┐
│  users   │────<│ user_roles   │────>│   roles     │
└──────────┘     └──────────────┘     └─────────────┘
                         │                    │
                         │                    │
                         v                    v
                  ┌─────────────┐     ┌──────────────┐
                  │ ngac_roles  │────>│ngac_role_   │
                  └─────────────┘     │ permissions │
                                       └──────────────┘
                                              │
                                              v
                                    ┌──────────────────┐
                                    │ngac_permissions  │
                                    └──────────────────┘

┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   menus     │────<│   articles   │────>│article_cat. │
└─────────────┘     └──────────────┘     └─────────────┘
                                              │
                                              │
                                      ┌───────▼───────┐
                                      │video_categories│
                                      └───────────────┘

┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│ system_     │     │  data_dicts  │     │operation_   │
│  configs    │     │              │     │  logs       │
└─────────────┘     └──────────────┘     └─────────────┘
```

## 种子数据初始化

### Seeder 机制

`app/migration/seeder.go` 负责在首次部署时初始化业务数据：

#### 1. 超级管理员账户
```
用户名: admin
密码: admin123
邮箱: admin@example.com
角色: super_admin
```

#### 2. 系统角色（roles 表）
| ID | 名称 | 代码 | 描述 |
|----|------|------|------|
| 1 | 超级管理员 | super_admin | 拥有系统所有权限 |
| 2 | 内容管理员 | content_admin | 管理文章内容 |
| 3 | 普通用户 | user | 普通用户权限 |

#### 3. NGAC 角色（ngac_roles 表）
与 roles 表保持一致，供 NGAC 策略决策引擎使用。

#### 4. NGAC 对象类（ngac_object_classes）
| ID | 名称 | Code |
|----|------|------|
| 1-12 | 文章、分类、视频、评论、文件、图库、用户、角色、配置、菜单、日志、字典 | article, video, user, role... |

#### 5. NGAC 权限（ngac_permissions）
包含 29 个预定义权限，覆盖所有核心功能：
- `dashboard:view` - 查看仪表盘
- `content:article:*` - 文章 CRUD 操作
- `content:video:*` - 视频 CRUD 操作
- `user:list`, `user:create`, `user:update`, `user:delete` - 用户管理
- `system:config`, `system:menu`, `system:log`, `system:dict` - 系统管理

#### 6. 超级管理员关联
- super_admin 角色关联所有 29 个权限
- admin 用户关联 super_admin 角色

#### 7. 菜单数据（menus 表）
包含完整的后台管理菜单结构：
- 仪表盘
- 内容管理 → 文章管理、文章分类、视频管理、视频分类、评论管理
- 媒体中心 → 文件管理、图库管理
- 用户管理 → 用户列表、角色管理、权限管理
- 系统设置 → 系统配置、菜单管理、操作日志、数据字典

#### 8. 文章/视频分类
- 文章分类：技术分享、产品资讯、业界动态、开发教程、开源项目
- 视频分类：视频教程、直播回放、课程专区、技术分享

#### 9. 系统配置（system_configs）
包含 20 项默认配置：站点信息、上传设置、安全设置、CDN、缓存、国际化等。

#### 10. 数据字典（data_dicts）
- `article_status`: 草稿、已发布、已下架
- `video_status`: 草稿、已发布、审核中
- `comment_status`: 待审核、已通过、已拒绝、已删除
- `user_status`: 启用、禁用

## 执行流程

### 1. 执行迁移
```go
manager := migration.NewManager()
manager.RunAllMigrations()  // 执行所有 SQL 迁移文件
seeder := migration.NewSeeder(db)
seeder.Seed()               // 执行种子数据初始化
```

### 2. 幂等性保证
- 所有 SQL 使用 `CREATE TABLE IF NOT EXISTS`
- INSERT 使用 `ON DUPLICATE KEY UPDATE`
- Seeder 检查表和数据是否存在再决定是否初始化
- 迁移管理器自动捕获并跳过 "Duplicate" 相关错误

## 重置数据库

如需完全重置数据库：
1. 清空所有表数据
2. 重新执行迁移
3. 重新执行种子数据

```sql
-- 清空数据（注意顺序，考虑外键约束）
TRUNCATE TABLE ngac_role_permissions;
TRUNCATE TABLE ngac_user_roles;
TRUNCATE TABLE user_roles;
TRUNCATE TABLE menus;
TRUNCATE TABLE roles;
TRUNCATE TABLE ngac_roles;
-- ... 清空其他表

-- 重新执行
-- Go 代码自动处理迁移和种子初始化
```

## 注意事项

1. **不要直接修改迁移文件**：已执行的迁移不应再修改，如需变更请创建新迁移
2. **密码安全**：admin 默认密码为 admin123，首次登录后请立即修改
3. **角色表双写**：当前同时使用 `roles` 和 `ngac_roles` 表，两者保持一致
4. **字段类型一致性**：确保 Go model 中的字段与数据库迁移文件定义一致

## 未来改进建议

1. 移除 `roles` 表，统一使用 `ngac_roles`
2. 添加数据库版本控制机制（如 migrations 表记录已执行版本）
3. 增加更多环境特定的种子数据（开发/测试/生产）
4. 支持自定义种子数据的扩展机制