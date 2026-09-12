# GoCMS v3 系统架构与前端开发指引

**文档版本**: v4.0  
**最后更新**: 2026-08-30  
**适用版本**: GoCMS v3

---

## 目录

1. [系统概述](#系统概述)
2. [技术架构](#技术架构)
3. [后端架构详解](#后端架构详解)
4. [前端架构详解](#前端架构详解)
5. [核心模块 API 参考](#核心模块-api-参考)
6. [NGAC 权限管理](#ngac-权限管理)
7. [数据库迁移](#数据库迁移)
8. [前端开发指引](#前端开发指引)

---

## 系统概述

GoCMS v3 是一个基于 Go 语言开发的内容管理系统，采用 Gin 框架作为 Web 框架，支持文章、视频、用户、权限、菜单等多模块管理。系统采用 NGAC（New Grid Access Control）实现细粒度的 ABAC 权限控制。

### 系统特点

- ✅ **多租户支持** - 基于 tenant_id 的数据隔离
- ✅ **NGAC 权限系统** - RBAC + ABAC 混合访问控制
- ✅ **RESTful API** - 规范的 API 设计
- ✅ **前后端分离** - Go 后端 + Vue 3 TypeScript 前端
- ✅ **Tailwind CSS UI** - TailAdmin 风格设计系统
- ✅ **缓存优化** - Redis 缓存支持
- ✅ **WebSocket** - 实时通知支持
- ✅ **GraphQL** - GraphQL 查询支持

---

## 技术架构

```
┌─────────────────────────────────────────────────────┐
│                    前端层 (Vue 3 + TS)                 │
│         Tailwind CSS + Heroicons Vue                  │
└─────────────────────────────────────────────────────┘
                            |
                    HTTP/REST API + WebSocket
                            |
┌─────────────────────────────────────────────────────┐
│                    API 层 (Gin)                      │
│   router → controller → service → model              │
└─────────────────────────────────────────────────────┘
                            |
┌─────────────────────────────────────────────────────┐
│                  数据存储层                            │
│        MySQL + Redis + Elasticsearch (可选)           │
└─────────────────────────────────────────────────────┘
```

---

## 后端架构详解

### 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| Web 框架 | Gin | v1.9.1 |
| ORM | GORM | v1.25.5 |
| 数据库 | MySQL | 8.0+ |
| 缓存 | Redis | go-redis/v9 |
| 认证 | JWT | v5.1.0 |
| 定时任务 | robfig/cron | v3.0.1 |
| 日志 | Zap | v1.28.0 |
| WebSocket | gorilla/websocket | v1.5.3 |
| 系统指标 | gopsutil | v3.24.5 |

### 项目结构

```
gocms_v3/
├── app/                          # 应用核心
│   ├── application/              # 应用层
│   ├── config/                   # 配置管理
│   ├── controller/               # 控制器层 (12个)
│   │   ├── auth_controller.go
│   │   ├── article_controller.go
│   │   ├── video_controller.go
│   │   ├── user_controller.go
│   │   ├── menu_controller.go
│   │   ├── comment_controller.go
│   │   ├── ngac_controller.go
│   │   ├── dashboard_controller.go
│   │   ├── batch_controller.go
│   │   ├── upload_controller.go
│   │   ├── article_category_controller.go
│   │   └── video_category_controller.go
│   ├── db/                       # 数据库连接
│   ├── domain/                   # 领域模型
│   ├── dto/                      # 数据传输对象
│   ├── handler/                  # 处理器 (含 GraphQL)
│   ├── i18n/                     # 国际化
│   ├── infrastructure/           # 基础设施
│   ├── logger/                   # 日志
│   ├── middleware/               # 中间件
│   ├── migration/                # 数据库迁移
│   │   ├── manager.go            # 迁移管理器
│   │   └── seeder.go             # 数据种子
│   ├── model/                    # 数据模型 (13个)
│   │   ├── user.go
│   │   ├── article.go
│   │   ├── video.go
│   │   ├── menu.go
│   │   ├── comment.go
│   │   ├── category.go
│   │   ├── tag.go
│   │   ├── file_upload.go
│   │   ├── dashboard.go
│   │   ├── base_model.go
│   │   └── ngac/                 # NGAC 权限模型
│   │       └── models.go
│   ├── response/                 # 响应封装
│   ├── router/                   # 路由配置
│   │   └── router.go
│   └── service/                  # 业务逻辑层 (16个)
│       ├── ngac_service.go
│       ├── abac_engine.go
│       ├── pdp_engine.go
│       ├── article_service.go
│       ├── video_service.go
│       ├── user_service.go
│       ├── menu_service.go
│       ├── menu_service_cache.go
│       ├── comment_service.go
│       ├── dashboard_service.go
│       ├── content_version_service.go
│       ├── export_service.go
│       ├── notification_service.go
│       ├── article_category_service.go
│       └── video_category_service.go
├── bin/                          # 编译输出 (gocms_v3_server.exe)
├── cmd/server/                   # 程序入口
│   └── main.go
├── configs/                      # 配置文件
│   └── config.yaml
├── internal/                     # 内部模块
│   ├── cron/                     # 定时任务
│   ├── crypto/                   # 加密工具
│   ├── errors/                   # 错误处理
│   ├── oauth2/                   # OAuth2 支持
│   ├── signature/                # 签名工具
│   ├── utils/                    # 工具函数
│   ├── webhook/                  # Webhook
│   └── websocket/                # WebSocket
├── migrations/                   # SQL 迁移文件 (4个)
│   ├── 001_add_tenant_fields.sql
│   ├── 002_create_comments.sql
│   ├── 003_init_system_data.sql
│   └── 004_fix_name_column_type.sql
├── storage/                      # 文件存储
├── swagger/                      # Swagger API 文档 (待完善)
│   └── swagger.yaml
└── docs/                         # 文档
```

### 启动流程

```
服务器启动 (cmd/server/main.go)
    │
    ▼
[步骤1] 加载配置文件 (configs/config.yaml)
    │
    ▼
[步骤2] 初始化数据库连接 (MySQL: 48.48.48.126:3306)
    │
    ▼
[步骤3] 初始化 Redis 连接 (48.48.48.126:6379)
    │
    ▼
[步骤4] 执行 SQL 迁移文件 (migrations/)
    │
    ▼
[步骤5] GORM AutoMigrate 同步表结构
    │
    ▼
[步骤6] 数据种子初始化 (seeder.go)
    │   ├── 超级管理员用户 (admin/admin123)
    │   ├── 角色数据 (3个角色)
    │   ├── 权限数据 (30条)
    │   └── 菜单数据 (19条)
    │
    ▼
[步骤7] 启动 HTTP 服务器 (端口: 8084)
```

---

## 前端架构详解

### 技术栈

| 类别 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue | v3.5.40 |
| 类型系统 | TypeScript | ~6.0.0 |
| 状态管理 | Pinia | v4.0.2 |
| 路由 | Vue Router | v5.2.0 |
| HTTP 客户端 | Axios | v1.19.0 |
| 构建工具 | Vite | v8.1.5 |
| CSS 框架 | Tailwind CSS | v3.4.19 |
| 图标 | Heroicons Vue | ^2.2.0 |
| 工具库 | @vueuse/core | v14.4.0 |
| 多语言 | vue-i18n | v11.4.10 |

### 项目结构

```
frontend/
├── src/
│   ├── api/                    # API 接口封装 (9个模块)
│   │   ├── index.ts            # 统一导出 + Axios 实例
│   │   ├── auth.ts             # 认证 API
│   │   ├── user.ts             # 用户 API
│   │   ├── role.ts             # 角色 API
│   │   ├── menu.ts             # 菜单 API
│   │   ├── article.ts          # 文章 API
│   │   ├── video.ts            # 视频 API
│   │   ├── comment.ts          # 评论 API
│   │   └── dashboard.ts        # 仪表盘 API
│   ├── assets/                 # 静态资源
│   ├── components/             # 公共组件
│   │   ├── ui/                 # TailAdmin UI 组件 (6个)
│   │   │   ├── TailButton.vue
│   │   │   ├── TailInput.vue
│   │   │   ├── TailCard.vue
│   │   │   ├── TailModal.vue
│   │   │   ├── TailTable.vue
│   │   │   ├── TailPagination.vue
│   │   │   └── index.ts
│   │   ├── icons/
│   │   └── ...
│   ├── layouts/                # 布局组件 (3个)
│   │   ├── MainLayout.vue
│   │   ├── Sidebar.vue
│   │   └── TopHeader.vue
│   ├── locales/                # 多语言文件 (2个)
│   │   ├── zh-CN.ts
│   │   ├── en-US.ts
│   │   └── index.ts
│   ├── pages/                  # 页面组件 (12+个)
│   │   ├── dashboard/          # 仪表盘
│   │   │   └── DashboardPage.vue
│   │   ├── login/              # 登录页
│   │   │   └── LoginPage.vue
│   │   ├── register/           # 注册页
│   │   │   └── RegisterPage.vue
│   │   ├── users/              # 用户管理
│   │   │   └── UserListPage.vue
│   │   ├── articles/           # 文章管理 (3个)
│   │   │   ├── ArticleListPage.vue
│   │   │   ├── ArticleFormPage.vue
│   │   │   └── CategoryPage.vue
│   │   ├── videos/             # 视频管理 (3个)
│   │   │   ├── VideoListPage.vue
│   │   │   ├── VideoFormPage.vue
│   │   │   └── CategoryPage.vue
│   │   ├── comments/           # 评论管理
│   │   │   └── CommentListPage.vue
│   │   └── system/             # 系统管理 (2个)
│   │       ├── MenuPage.vue
│   │       └── RolePage.vue
│   ├── router/                 # 路由配置
│   │   └── index.ts
│   ├── stores/                 # Pinia 状态管理 (3个)
│   │   ├── auth.ts
│   │   ├── counter.ts
│   │   └── loading.ts
│   ├── types/                  # TypeScript 类型
│   └── utils/                  # 工具函数 (2个)
│       ├── message.ts          # 自定义消息提示
│       └── errorHandler.ts     # 错误处理
├── public/
├── env.development             # 开发环境配置
├── package.json
├── vite.config.ts
├── tailwind.config.js
├── postcss.config.js
└── tsconfig.json
```

---

## 核心模块 API 参考

### 认证模块 (/api/auth)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/login` | 用户登录 | 否 |
| POST | `/auth/register` | 用户注册 | 否 |
| POST | `/auth/logout` | 用户登出 | 否 |
| GET | `/user/info` | 获取当前用户信息 | 是 |
| PUT | `/user/info` | 更新用户信息 | 是 |
| GET | `/auth/permissions` | 获取当前用户权限 | 是 |

### 用户模块 (/api/users)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/users` | 用户列表（分页） | 是 |
| POST | `/user` | 创建用户 | 是 |

### 文章模块 (/api/articles)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/articles` | 文章列表（分页） | 是 |
| GET | `/articles/:id` | 文章详情 | 否 |
| POST | `/article` | 创建文章 | 是 |
| PUT | `/article/:id` | 更新文章 | 是 |
| DELETE | `/article/:id` | 删除文章 | 是 |

### 文章分类模块 (/api/article-categories)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/article-categories` | 创建分类 | 是 |
| GET | `/article-categories` | 分类树 | 是 |
| GET | `/article-categories/:id` | 分类详情 | 是 |
| PUT | `/article-categories/:id` | 更新分类 | 是 |
| DELETE | `/article-categories/:id` | 删除分类 | 是 |

### 视频模块 (/api/videos)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/videos` | 视频列表 | 是 |
| GET | `/videos/:id` | 视频详情 | 否 |
| GET | `/videos/published` | 已发布视频 | 否 |
| POST | `/video` | 创建视频 | 是 |
| PUT | `/video/:id` | 更新视频 | 是 |
| DELETE | `/video/:id` | 删除视频 | 是 |

### 视频分类模块 (/api/video-categories)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/video-categories` | 创建分类 | 是 |
| GET | `/video-categories` | 分类列表 | 是 |
| GET | `/video-categories/tree` | 分类树 | 否 |
| GET | `/video-categories/:id` | 分类详情 | 是 |
| PUT | `/video-categories/:id` | 更新分类 | 是 |
| DELETE | `/video-categories/:id` | 删除分类 | 是 |

### 视频操作模块

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/video/:id/publish` | 发布视频 | 是 |
| POST | `/video/:id/unpublish` | 取消发布 | 是 |
| POST | `/video/:id/top` | 置顶视频 | 是 |
| POST | `/video/:id/untop` | 取消置顶 | 是 |
| GET | `/video/:id/rate` | 视频评分 | 是 |

### 菜单模块 (/api/menus)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/menus/tree` | 菜单树 | 否 |
| GET | `/menus/permissions` | 菜单权限 | 是 |
| POST | `/menus` | 创建菜单 | 是 |
| GET | `/menus` | 菜单列表 | 是 |
| GET/PUT/DELETE | `/menus/:id` | 菜单详情操作 | 是 |

### Dashboard 模块 (/api/dashboard)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/dashboard/stats` | 综合统计 | 是 |
| GET | `/dashboard/overview` | 概览数据 | 是 |

### 评论模块 (/api/comments)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/comments` | 评论列表 | 是 |
| GET | `/comments/stats` | 评论统计 | 是 |
| POST | `/comments` | 创建评论 | 是 |
| PUT | `/comments/:id/audit` | 评论审核 | 是 |
| DELETE | `/comments/:id` | 删除评论 | 是 |

### 文件上传模块 (/api/upload)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/upload/initiate` | 初始化分片上传 | 是 |
| POST | `/upload/chunk` | 上传分片 | 是 |
| POST | `/upload/complete` | 完成上传 | 是 |

### 批量操作模块 (/api/batch)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/articles/batch` | 批量创建文章 | 是 |
| POST | `/videos/batch` | 批量创建视频 | 是 |

### NGAC 权限模块 (/api/ngac)

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| **角色管理** | | | |
| POST | `/ngac/roles` | 创建角色 | 是 |
| GET | `/ngac/roles` | 角色列表 | 是 |
| PUT | `/ngac/roles/:role_id` | 更新角色 | 是 |
| DELETE | `/ngac/roles/:role_id` | 删除角色 | 是 |
| **权限管理** | | | |
| POST | `/ngac/permissions` | 创建权限规则 | 是 |
| GET | `/ngac/permissions` | 权限列表 | 是 |
| PUT | `/ngac/permissions/:perm_id` | 更新权限 | 是 |
| DELETE | `/ngac/permissions/:perm_id` | 删除权限 | 是 |
| **角色-权限关联** | | | |
| POST | `/ngac/roles/:role_id/permissions` | 分配权限给角色 | 是 |
| DELETE | `/ngac/roles/:role_id/permissions/:perm_id` | 移除角色权限 | 是 |
| **用户-角色管理** | | | |
| POST | `/ngac/users/:user_id/roles` | 分配角色给用户 | 是 |
| GET | `/ngac/users/:user_id/roles` | 获取用户角色 | 是 |
| GET | `/ngac/users/:user_id/permissions` | 获取用户权限 | 是 |
| **属性管理** | | | |
| GET | `/ngac/users/:user_id/attributes` | 获取用户属性 | 是 |
| GET | `/ngac/context/attributes` | 获取上下文属性 | 是 |
| **授权评估** | | | |
| POST | `/ngac/authz` | 授权请求（PDP） | 是 |
| **对象类** | | | |
| GET | `/ngac/object-classes/tree` | 获取对象类树 | 是 |

### GraphQL

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/graphql` | GraphQL 查询 | 是 |

---

## NGAC 权限管理

### 架构概述

NGAC (New Grid Access Control) 是一个基于属性的访问控制（ABAC）扩展模型，在 ABAC 基础上增加了角色层和对象类层次结构。

```
NGAC = ABAC + Role Layer + Object Class Hierarchy
```

### 核心组件

| 组件 | 文件路径 | 说明 |
|------|----------|------|
| 数据模型 | `app/model/ngac/models.go` | 10个核心表 |
| NGAC 服务 | `app/service/ngac_service.go` | CRUD + 关联管理 |
| ABAC 引擎 | `app/service/abac_engine.go` | 属性评估引擎 |
| PDP 决策点 | `app/service/pdp_engine.go` | 策略决策引擎 |
| 控制器 | `app/controller/ngac_controller.go` | RESTful API |

### 核心表结构

| 表名 | 说明 |
|------|------|
| ngac_roles | 角色表（支持层次结构） |
| ngac_permissions | 权限/规则表 |
| ngac_role_permissions | 角色-权限关联表 |
| ngac_user_roles | 用户-角色关联表 |
| ngac_object_classes | 对象类表（支持层次结构） |
| ngac_object_attributes | 资源对象属性表 |
| ngac_user_attributes | 用户属性表 |
| ngac_context_attributes | 上下文/环境属性表 |
| ngac_access_logs | 访问审计日志表 |
| ngac_pdp_config | PDP引擎配置表 |

---

## 数据库迁移

### 迁移文件列表

| 文件 | 说明 | 状态 |
|------|------|------|
| 001_add_tenant_fields.sql | 多租户字段 + OAuth2 表 | ✅ 执行 |
| 002_create_comments.sql | 评论表创建 | ✅ 执行 |
| 003_init_system_data.sql | 初始化数据 | ⚠️ 已弃用（由 seeder.go 替代） |
| 004_fix_name_column_type.sql | 修复 name 列类型 | ✅ 执行 |

### 默认管理员账号

| 字段 | 值 |
|------|-----|
| username | admin |
| password | admin123 (bcrypt加密) |
| nickname | 系统管理员 |
| email | admin@example.com |
| roles | ["super_admin"] |

---

## 前端开发指引

### 技术栈（当前已选定）

| 类别 | 技术 | 说明 |
|------|------|------|
| 框架 | Vue 3 + TypeScript | ✅ 已采用 |
| 状态管理 | Pinia | ✅ 已采用 |
| CSS 框架 | Tailwind CSS | ✅ 已迁移完成 |
| 图标 | Heroicons Vue | ✅ 已采用 |

### 路由设计

```
公开路由:
  /login              登录
  /register           注册

认证后路由:
  /dashboard          仪表盘
  /articles           文章管理
  /videos             视频管理
  /comments           评论管理
  /users              用户管理
  /menu               菜单管理
```

### 认证机制

**Token 格式**: `Bearer <token>`

**请求头**:
```
Authorization: Bearer <jwt_token>
```

### 统一响应格式

**成功响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**错误响应**:
```json
{
  "code": 404,
  "message": "资源不存在"
}
```

**分页响应**:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [...],
    "total": 100,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 待完善功能（优先级）

### P1 - 高优先级
1. TailAdmin UI 迁移完成 (Phase 3) - 剩余页面迁移
2. GraphQL 功能完善
3. WebSocket 实时通知实现

### P2 - 中优先级
1. Swagger API 文档完善
2. Dashboard 统计增加 Redis 缓存
3. Elasticsearch 全文搜索集成

### P3 - 低优先级
1. 单元测试覆盖率提升
2. 性能优化
3. 多语言完善

---

*文档版本：v4.0*  
*最后更新：2026-08-30*