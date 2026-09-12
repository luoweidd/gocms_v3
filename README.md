# GoCMS v3

<div align="center">

```
一个功能丰富的企业级内容管理系统
```

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-v1.9-red.svg)](https://github.com/gin-gonic/gin)
[![Vue](https://img.shields.io/badge/Vue-3.5-green.svg)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**GoCMS v3** - 基于 Go + Vue 3 的企业级内容管理系统，支持 NGAC 细粒度权限控制、多租户、媒体管理和内容审核。

</div>

---

## 📋 目录

- [✨ 核心特性](#-核心特性)
- [🏗️ 系统架构](#️-系统架构)
- [🛠️ 技术栈](#️-技术栈)
- [📁 项目结构](#-项目结构)
- [🚀 快速开始](#-快速开始)
- [📡 API 文档](#-api-文档)
- [🗄️ 数据库设计](#️-数据库设计)
- [🔐 权限系统](#-权限系统)
- [🌍 多语言](#-多语言)
- [📦 部署](#-部署)
- [📝 开发指南](#-开发指南)

---

## ✨ 核心特性

| 特性 | 说明 | 状态 |
|------|------|------|
| 🏢 **多租户支持** | 基于 tenant_id 的数据隔离，支持多租户部署 | ✅ |
| 🔐 **NGAC 权限系统** | RBAC + ABAC 混合访问控制，细粒度权限管理 | ✅ |
| 📝 **内容管理** | 文章/视频的 CRUD、版本控制、内容审核 | ✅ |
| 🎬 **媒体中心** | 图片、音频、视频统一管理，支持分片上传 | ✅ |
| 💬 **评论系统** | 评论审核、批量操作、审计统计 | ✅ |
| 🔍 **SEO 管理** | 页面 SEO 设置、站点地图生成 | ✅ |
| 💾 **数据备份** | 自动/手动数据库备份与恢复 | ✅ |
| 📊 **仪表盘** | 实时数据统计、趋势分析、资源监控 | ✅ |
| 🌍 **多语言** | 中英文切换支持 (i18n) | ✅ |
| 🔔 **消息通知** | 系统消息、用户通知、已读追踪 | ✅ |
| 📋 **操作日志** | 完整操作审计、日志清理策略 | ✅ |
| 🔑 **OAuth2** | OAuth2 认证支持 | ✅ |
| 📡 **GraphQL** | GraphQL 查询支持 | ✅ |
| 🔄 **WebSocket** | 实时通知支持 | 🚧 |
| 🔎 **全文搜索** | Elasticsearch 集成 | 📋 待扩展 |

---

## 🏗️ 系统架构

```
┌──────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3 + TS)                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ Dashboard │  │ Content  │  │ Media    │  │ System   │    │
│  │          │  │ Management│  │ Center   │  │ Manage   │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
│                    Pinia · Tailwind CSS · Heroicons           │
└──────────────────────────────────────────────────────────────┘
                              │ HTTP/REST + WebSocket
                              │ Bearer Token (JWT)
┌──────────────────────────────────────────────────────────────┐
│                      后端 (Go + Gin)                         │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                    Middleware Layer                      │ │
│  │  CORS · JWT Auth · Operation Log · Rate Limit           │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                     Router Layer                         │ │
│  │  /api/v1/*  →  RESTful API Routes                       │ │
│  │  /graphql         →  GraphQL Endpoint                   │ │
│  │  /ws              →  WebSocket Hub                      │ │
│  │  /uploads/*       →  Static Files                       │ │
│  └─────────────────────────────────────────────────────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ Controller│→│ Service  │→│  Model   │→│  NGAC    │    │
│  │          │  │          │  │          │  │  Engine  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└──────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────┼──────────────────────────────┐
│                             ▼                              │
│  ┌────────────┐  ┌────────────┐  ┌────────────────────┐   │
│  │   MySQL    │  │   Redis    │  │   Elasticsearch    │   │
│  │  (数据持久化)│  │  (缓存/会话)│  │   (全文搜索/聚合)   │   │
│  └────────────┘  └────────────┘  └────────────────────┘   │
└──────────────────────────────────────────────────────────────┘
```

---

## 🛠️ 技术栈

### 后端技术

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **Web 框架** | [Gin](https://github.com/gin-gonic/gin) | v1.9.1 | HTTP 服务器 & 路由 |
| **ORM** | [GORM](https://gorm.io/) | v1.25.5 | 数据库操作 |
| **数据库** | MySQL | 8.0+ | 数据持久化 |
| **缓存** | [go-redis](https://github.com/redis/go-redis) | v9 | 缓存 & 会话存储 |
| **认证** | [jwt-go](https://github.com/golang-jwt/jwt) | v5.1.0 | JWT 令牌 |
| **定时任务** | [robfig/cron](https://github.com/robfig/cron) | v3.0.1 | 定时任务调度 |
| **日志** | [Zap](https://go.uber.org/zap) | v1.28.0 | 高性能日志 |
| **WebSocket** | [gorilla/websocket](https://github.com/gorilla/websocket) | v1.5.3 | 实时通信 |
| **系统指标** | [gopsutil](https://github.com/shirou/gopsutil) | v3.24.5 | 系统监控 |
| **Elasticsearch** | [elasticsearch-go](https://github.com/elastic/go-elasticsearch) | - | 全文搜索 |

### 前端技术

| 类别 | 技术 | 版本 | 用途 |
|------|------|------|------|
| **框架** | [Vue 3](https://vuejs.org/) | v3.5.40 | 渐进式 UI 框架 |
| **类型系统** | TypeScript | ~5.9.0 | 类型安全 |
| **状态管理** | [Pinia](https://pinia.vuejs.org/) | v2.3.0 | 响应式状态管理 |
| **路由** | [Vue Router](https://router.vuejs.org/) | v4.5.0 | 单页路由 |
| **HTTP 客户端** | [Axios](https://axios-http.com/) | v1.19.0 | HTTP 请求 |
| **构建工具** | [Vite](https://vitejs.dev/) | v7.1.2 | 开发与打包 |
| **CSS 框架** | [Tailwind CSS](https://tailwindcss.com/) | v3.4.17 | 原子化 CSS |
| **图标** | [@heroicons/vue](https://github.com/heroicons/heroui) | v2.2.0 | SVG 图标 |
| **UI 组件** | [Ant Design Vue](https://ant-design-vue.github.io/) | v4.2.6 | 企业级组件库 |
| **富文本编辑器** | [wangeditor](https://www.wangeditor.com/) + [Vditor](https://github.com/Vanessc/vditor) | - | 内容编辑 |
| **图表** | [ECharts](https://echarts.apache.org/) + [vue-echarts](https://github.com/apache/echarts/tree/main/package/vue-echarts) | v6.0.0 | 数据可视化 |
| **工具库** | [@vueuse/core](https://vueuse.org/) | v12.8.0 | Vue 组合式 API |
| **多语言** | [vue-i18n](https://github.com/intlify/vue-i18n) | v11.1.12 | i18n 国际化 |

---

## 📁 项目结构

```
gocms_v3/
├── app/                              # 应用核心
│   ├── application/                  # 应用层（业务编排）
│   ├── config/                       # 配置管理
│   ├── controller/                   # 控制器层 (14个)
│   │   ├── auth_controller.go       # 认证（登录/注册/登出）
│   │   ├── user_controller.go       # 用户管理
│   │   ├── role_controller.go       # 角色管理
│   │   ├── menu_controller.go       # 菜单管理
│   │   ├── article_controller.go    # 文章管理
│   │   ├── article_category_controller.go  # 文章分类
│   │   ├── video_controller.go      # 视频管理
│   │   ├── video_category_controller.go   # 视频分类
│   │   ├── comment_controller.go    # 评论管理
│   │   ├── dashboard_controller.go  # 仪表盘
│   │   ├── ngac_controller.go       # NGAC 权限管理
│   │   ├── upload_controller.go     # 文件上传
│   │   ├── batch_controller.go      # 批量操作
│   │   ├── media_controller.go      # 媒体中心
│   │   ├── audit_controller.go      # 内容审核
│   │   ├── backup_controller.go     # 数据备份
│   │   ├── seo_controller.go        # SEO管理
│   │   ├── message_controller.go    # 消息管理
│   │   ├── message_category_controller.go # 消息分类
│   │   ├── tenant_controller.go     # 租户管理
│   │   ├── tag_controller.go        # 标签管理
│   │   └── operation_log_controller.go # 操作日志
│   ├── db/                           # 数据库连接（MySQL + Redis）
│   ├── domain/                       # 领域模型定义
│   ├── dto/                          # 数据传输对象
│   ├── handler/                      # 处理器（含 GraphQL）
│   ├── i18n/                         # 国际化资源
│   ├── infrastructure/               # 基础设施
│   ├── logger/                       # 日志配置
│   ├── middleware/                   # 中间件（CORS/JWT/日志等）
│   ├── migration/                    # 迁移与数据种子
│   │   ├── manager.go               # 迁移管理器
│   │   └── seeder.go                # 数据种子
│   ├── model/                        # 数据模型 (25+个)
│   │   ├── user.go                   # 用户模型
│   │   ├── role.go                   # 角色模型
│   │   ├── menu.go                   # 菜单模型
│   │   ├── article.go                # 文章模型
│   │   ├── video.go                  # 视频模型
│   │   ├── comment.go                # 评论模型
│   │   ├── tag.go                    # 标签模型
│   │   ├── category.go               # 分类模型
│   │   ├── file_upload.go            # 文件上传
│   │   ├── dashboard.go              # 仪表盘
│   │   ├── base_model.go             # 基础模型
│   │   ├── tenant.go                 # 租户模型
│   │   ├── media.go                  # 媒体资产
│   │   ├── content_audit.go          # 内容审核
│   │   ├── system_config.go          # 系统配置
│   │   ├── operation_log.go          # 操作日志
│   │   ├── notification.go           # 通知模型
│   │   ├── message_category.go       # 消息分类
│   │   ├── seo.go                    # SEO设置
│   │   ├── backup.go                 # 备份记录
│   │   ├── data_dict.go              # 数据字典
│   │   ├── risk_keyword.go           # 风险关键词
│   │   ├── oauth_auth_code.go        # OAuth授权码
│   │   ├── api_sign_key.go           # API签名密钥
│   │   └── ngac/                     # NGAC 权限模型
│   │       ├── models.go             # 10个NGAC核心表
│   │       ├── role.go               # 角色
│   │       ├── permission.go         # 权限规则
│   │       ├── user.go               # 用户角色关联
│   │       └── ...
│   ├── response/                     # 响应封装
│   ├── router/                       # 路由配置
│   └── service/                      # 业务逻辑层 (16+个)
│       ├── ngac_service.go           # NGAC服务
│       ├── abac_engine.go            # ABAC属性评估引擎
│       ├── pdp_engine.go             # PDP策略决策点
│       ├── article_service.go        # 文章服务
│       ├── video_service.go          # 视频服务
│       ├── user_service.go           # 用户服务
│       ├── menu_service.go           # 菜单服务
│       ├── comment_service.go        # 评论服务
│       ├── dashboard_service.go      # 仪表盘服务
│       ├── notification_service.go   # 通知服务
│       ├── content_version_service.go # 内容版本服务
│       ├── export_service.go         # 导出服务
│       ├── article_category_service.go # 文章分类服务
│       ├── video_category_service.go # 视频分类服务
│       ├── media_service.go          # 媒体服务
│       ├── audit_service.go          # 审核服务
│       ├── backup_service.go         # 备份服务
│       ├── seo_service.go            # SEO服务
│       └── ...
├── bin/                              # 编译输出目录
├── cmd/server/                       # 程序入口
│   └── main.go                       # 服务器启动入口
├── configs/                          # 配置文件
│   ├── config.yaml                   # 主配置（MySQL/Redis/JWT/CORS等）
│   └── config.yaml.example           # 配置示例
├── internal/                         # 内部模块
│   ├── cron/                         # 定时任务
│   ├── crypto/                       # 加密工具
│   ├── errors/                       # 错误处理
│   ├── oauth2/                       # OAuth2 支持
│   ├── signature/                    # API签名
│   ├── utils/                        # 通用工具函数
│   ├── webhook/                      # Webhook集成
│   └── websocket/                    # WebSocket Hub
├── frontend/                         # 前端项目
│   ├── src/
│   │   ├── api/                      # API接口封装 (9个模块)
│   │   ├── assets/                   # 静态资源
│   │   ├── components/               # 公共组件
│   │   │   ├── ui/                   # UI基础组件
│   │   │   │   ├── TailButton.vue
│   │   │   │   ├── TailInput.vue
│   │   │   │   ├── TailCard.vue
│   │   │   │   ├── TailModal.vue
│   │   │   │   ├── TailTable.vue
│   │   │   │   └── TailPagination.vue
│   │   │   └── icons/                # 图标组件
│   │   ├── layouts/                  # 布局组件 (3个)
│   │   ├── locales/                  # i18n 翻译文件
│   │   ├── pages/                    # 页面组件 (25+个)
│   │   │   ├── dashboard/            # 仪表盘
│   │   │   ├── login/                # 登录页
│   │   │   ├── register/             # 注册页
│   │   │   ├── articles/             # 文章管理 (4个)
│   │   │   ├── videos/               # 视频管理 (3个)
│   │   │   ├── comments/             # 评论管理
│   │   │   ├── users/                # 用户管理 (2个)
│   │   │   ├── system/               # 系统管理 (7个)
│   │   │   ├── messages/             # 消息通知 (5个)
│   │   │   ├── audit/                # 内容审核 (3个)
│   │   │   ├── media/                # 媒体中心 (4个)
│   │   │   ├── cron/                 # 定时任务 (2个)
│   │   │   └── tenant/               # 租户管理
│   │   ├── router/                   # 路由配置
│   │   ├── stores/                   # Pinia 状态管理
│   │   ├── types/                    # TypeScript类型定义
│   │   └── utils/                    # 工具函数
│   ├── public/                       # 公共资源
│   ├── package.json                  # 依赖配置
│   ├── vite.config.ts               # Vite 配置
│   ├── tailwind.config.js           # Tailwind CSS配置
│   └── tsconfig.json                # TypeScript配置
├── docs/                             # 项目文档
│   ├── SYSTEM_ARCHITECTURE.md       # 系统架构详解
│   ├── NGAC_GUIDE.md                # NGAC权限管理指南
│   ├── frontend-backend-analysis.md # 前后端分析
│   └── ...
├── storage/                          # 文件存储
│   ├── uploads/                      # 上传文件
│   ├── videos/                       # 视频文件
│   └── backups/                      # 数据库备份
├── logs/                             # 日志目录
├── go.mod                            # Go 模块依赖
├── go.sum                            # 依赖校验
├── config.yaml                       # 配置文件
└── README.md                         # 项目说明（本文件）
```

---

## 🚀 快速开始

### 环境要求

| 组件 | 最低版本 | 推荐版本 |
|------|----------|----------|
| Go | 1.21+ | 1.21+ |
| Node.js | 18+ | 22.12.0+ |
| MySQL | 5.7 | 8.0+ |
| Redis | 6.0 | 7.0+ |
| Elasticsearch | 7.x | 8.x (可选) |

### 后端安装

```bash
# 1. 克隆项目
git clone https://github.com/your-org/gocms_v3.git
cd gocms_v3

# 2. 配置数据库
# 编辑 configs/config.yaml，修改数据库连接信息
# mysql:
#   host: localhost
#   port: 3306
#   username: root
#   password: your_password
#   database: gocms_v3

# 3. 创建数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS gocms_v3 CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 4. 下载依赖
go mod download

# 5. 编译（可选）
go build -o bin/gocms_server ./cmd/server/

# 6. 启动服务器
go run ./cmd/server/main.go
```

### 前端安装

```bash
# 1. 进入前端目录
cd frontend

# 2. 安装依赖
npm install

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env，设置 VITE_API_BASE_URL

# 4. 启动开发服务器
npm run dev
```

### 访问系统

| 环境 | URL | 说明 |
|------|-----|------|
| **前端** | http://localhost:5173 | Vite 开发服务器 |
| **后端API** | http://localhost:8084 | Gin HTTP 服务器 |
| **API文档** | http://localhost:8084/api | RESTful API |
| **GraphQL** | http://localhost:8084/graphql | GraphQL 端点 |

### 默认管理员账号

```
用户名: admin
密码: admin123
邮箱: admin@example.com
角色: super_admin
```

---

## 📡 API 文档

### 统一响应格式

**成功响应:**
```json
{
  "code": 0,
  "message": "success",
  "data": { ... }
}
```

**分页响应:**
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

**错误响应:**
```json
{
  "code": 4001,
  "message": "资源不存在"
}
```

### API 端点总览

#### 认证模块 `/api/auth`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/auth/login` | 用户登录 | ❌ |
| POST | `/auth/register` | 用户注册 | ❌ |
| POST | `/auth/logout` | 用户登出 | ✅ |
| GET | `/user/info` | 获取当前用户信息 | ✅ |
| PUT | `/user/info` | 更新用户信息 | ✅ |
| GET | `/auth/permissions` | 获取当前用户权限 | ✅ |

#### 用户管理 `/api/users`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/users` | 用户列表（分页） | ✅ |
| GET | `/users/:id` | 用户详情 | ✅ |
| POST | `/users` | 创建用户 | ✅ |
| PUT | `/users/:id` | 更新用户 | ✅ |
| DELETE | `/users/:id` | 删除用户 | ✅ |

#### 角色权限 `/api/roles`, `/api/ngac`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/roles` | 角色列表 | ✅ |
| GET | `/roles/:id` | 角色详情 | ✅ |
| POST | `/roles` | 创建角色 | ✅ |
| PUT | `/roles/:id` | 更新角色 | ✅ |
| DELETE | `/roles/:id` | 删除角色 | ✅ |
| GET | `/permissions/tree` | 权限树 | ✅ |
| POST | `/ngac/roles` | NGAC 创建角色 | ✅ |
| GET | `/ngac/roles` | NGAC 角色列表 | ✅ |
| POST | `/ngac/permissions` | 创建权限规则 | ✅ |
| POST | `/ngac/authz` | 授权评估 (PDP) | ✅ |

#### 文章模块 `/api/articles`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/articles` | 文章列表（分页） | ✅ |
| GET | `/articles/:id` | 文章详情 | ❌ |
| POST | `/articles` | 创建文章 | ✅ |
| PUT | `/articles/:id` | 更新文章 | ✅ |
| DELETE | `/articles/:id` | 删除文章 | ✅ |
| GET | `/article-categories` | 文章分类列表 | ✅ |
| GET | `/article-categories/tree` | 文章分类树 | ✅ |

#### 视频模块 `/api/videos`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/videos` | 视频列表（分页） | ✅ |
| GET | `/videos/:id` | 视频详情 | ❌ |
| GET | `/videos/published` | 已发布视频 | ❌ |
| POST | `/videos` | 创建视频 | ✅ |
| PUT | `/videos/:id` | 更新视频 | ✅ |
| DELETE | `/videos/:id` | 删除视频 | ✅ |
| POST | `/video/:id/publish` | 发布视频 | ✅ |
| POST | `/video/:id/unpublish` | 取消发布 | ✅ |
| POST | `/video/:id/top` | 置顶视频 | ✅ |
| GET | `/video/:id/rate` | 视频评分 | ✅ |

#### 评论模块 `/api/comments`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/comments` | 评论列表（分页） | ✅ |
| GET | `/comments/stats` | 评论统计 | ✅ |
| POST | `/comments` | 创建评论 | ✅ |
| PUT | `/comments/:id/approve` | 审核通过 | ✅ |
| DELETE | `/comments/:id` | 删除评论 | ✅ |
| DELETE | `/comments/batch-delete` | 批量删除 | ✅ |

#### 媒体中心 `/api/media`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/media/assets` | 媒体资产列表 | ✅ |
| GET | `/media/assets/stats` | 媒体统计 | ✅ |
| GET | `/media/assets/:id` | 媒体详情 | ✅ |
| POST | `/media/assets` | 创建媒体资产 | ✅ |
| PUT | `/media/assets/:id` | 更新媒体 | ✅ |
| DELETE | `/media/assets/:id` | 删除媒体 | ✅ |
| GET | `/media/albums` | 相册列表 | ✅ |
| POST | `/media/albums` | 创建相册 | ✅ |

#### 文件上传 `/api/upload`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/upload/initiate` | 初始化分片上传 | ✅ |
| POST | `/upload/chunk` | 上传分片 | ✅ |
| POST | `/upload/complete` | 完成上传 | ✅ |
| POST | `/upload/register` | 注册文件 | ✅ |

#### 仪表盘 `/api/dashboard`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/dashboard/stats` | 综合统计数据 | ✅ |
| GET | `/dashboard/overview` | 概览数据 | ✅ |
| GET | `/dashboard/trend` | 趋势分析 | ✅ |
| GET | `/dashboard/recent-activity` | 最近活动 | ✅ |
| GET | `/dashboard/resource-usage` | 资源使用统计 | ✅ |

#### 内容审核 `/api/audit`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/audit/submit` | 提交审核 | ✅ |
| GET | `/audit/list` | 审核列表 | ✅ |
| GET | `/audit/:id` | 审核详情 | ✅ |
| PUT | `/audit/:id/approve` | 审核通过 | ✅ |
| PUT | `/audit/:id/reject` | 审核驳回 | ✅ |
| POST | `/audit/batch` | 批量审核 | ✅ |

#### 消息通知 `/api/messages`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/messages` | 发布系统消息 | ✅ |
| GET | `/messages` | 消息列表 | ✅ |
| GET | `/messages/my` | 我的通知 | ✅ |
| GET | `/messages/my/stats` | 通知统计 | ✅ |
| PUT | `/messages/:id/read` | 标记已读 | ✅ |
| POST | `/messages/my/all-read` | 全部已读 | ✅ |

#### SEO 管理 `/api/seo`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/seo/settings` | 创建SEO设置 | ✅ |
| GET | `/seo/settings` | SEO列表 | ✅ |
| GET | `/seo/settings/:id` | SEO详情 | ✅ |
| PUT | `/seo/settings/:id` | 更新SEO设置 | ✅ |
| DELETE | `/seo/settings/:id` | 删除SEO设置 | ✅ |
| GET | `/seo/sitemap-configs` | 站点地图配置 | ✅ |
| POST | `/seo/sitemap/generate` | 生成站点地图 | ✅ |

#### 数据备份 `/api/backup`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/backup/full` | 创建完整备份 | ✅ |
| GET | `/backup/list` | 备份列表 | ✅ |
| GET | `/backup/:id` | 备份详情 | ✅ |
| GET | `/backup/:id/download` | 下载备份 | ✅ |
| DELETE | `/backup/:id` | 删除备份 | ✅ |
| POST | `/backup/:id/restore` | 恢复备份 | ✅ |

#### 租户管理 `/api/tenants`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/tenants` | 租户列表 | ✅ |
| POST | `/tenants` | 创建租户 | ✅ |
| GET | `/tenants/:id` | 租户详情 | ✅ |
| PUT | `/tenants/:id` | 更新租户 | ✅ |
| DELETE | `/tenants/:id` | 删除租户 | ✅ |
| POST | `/tenants/:id/users` | 添加用户 | ✅ |

#### 标签管理 `/api/tags`

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/tags` | 标签列表 | ✅ |
| GET | `/tags/:id` | 标签详情 | ✅ |
| POST | `/tags` | 创建标签 | ✅ |
| PUT | `/tags/:id` | 更新标签 | ✅ |
| DELETE | `/tags/:id` | 删除标签 | ✅ |

#### GraphQL

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/graphql` | GraphQL 查询/突变 | ✅ |

---

## 🗄️ 数据库设计

### 核心表结构

#### 用户与权限

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `users` | 用户表 | username, email, password_hash, tenant_id |
| `roles` | 角色表 | name, code, tenant_id |
| `user_roles` | 用户-角色关联 | user_id, role_id |
| `menus` | 菜单表 | name, path, parent_id, sort_order |
| `ngac_roles` | NGAC 角色（层次结构） | name, parent_id, object_class_id |
| `ngac_permissions` | NGAC 权限规则 | name, condition_expr, role_id |
| `ngac_user_roles` | 用户-NGAC角色关联 | user_id, role_id, attribute_context |
| `ngac_object_classes` | 对象类层次结构 | name, parent_id |
| `ngac_object_attributes` | 资源属性 | object_class_id, key, value |
| `ngac_user_attributes` | 用户属性 | user_id, key, value |
| `ngac_context_attributes` | 上下文属性 | key, value |
| `ngac_access_logs` | 访问审计日志 | subject_id, object_id, decision |

#### 内容管理

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `articles` | 文章表 | title, content, category_id, user_id |
| `article_categories` | 文章分类 | name, parent_id, sort_order |
| `videos` | 视频表 | title, video_url, thumbnail, category_id |
| `video_categories` | 视频分类 | name, parent_id, sort_order |
| `comments` | 评论表 | content, user_id, target_type, target_id |
| `tags` | 标签表 | name, article_count |
| `article_tag` | 文章-标签关联 | article_id, tag_id |

#### 媒体与文件

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `media_assets` | 媒体资产表 | file_url, file_type, file_size, tenant_id |
| `albums` | 相册表 | name, description, user_id |
| `album_assets` | 相册-媒体关联 | album_id, asset_id |
| `file_uploads` | 文件上传记录 | file_name, file_path, chunk_count |

#### 系统管理

| 表名 | 说明 | 关键字段 |
|------|------|----------|
| `operation_logs` | 操作日志表 | user_id, action, resource, ip_address |
| `system_config` | 系统配置表 | key, value, group, is_public |
| `seo_settings` | SEO设置表 | page_name, meta_title, meta_keywords |
| `sitemap_configs` | 站点地图配置 | name, url, priority, change_freq |
| `backup_records` | 备份记录表 | backup_type, file_path, status |
| `content_audits` | 内容审核表 | target_type, target_id, audit_status |
| `system_messages` | 系统消息表 | title, content, publish_status |
| `user_notifications` | 用户通知表 | user_id, message_id, is_read |
| `message_categories` | 消息分类表 | name, code, icon |
| `data_dict` | 数据字典表 | group_key, dict_key, dict_value |
| `risk_keywords` | 风险关键词表 | keyword, risk_level, action_type |
| `tenants` | 租户表 | name, domain, status, expire_at |

---

## 🔐 权限系统 (NGAC)

### 架构概述

**NGAC (New Grid Access Control)** 是一个基于属性的访问控制（ABAC）扩展模型，结合了 RBAC 和 ABAC 的优点：

```
NGAC = ABAC + Role Layer + Object Class Hierarchy
```

### NGAC 核心组件

| 组件 | 文件路径 | 说明 |
|------|----------|------|
| **数据模型** | `app/model/ngac/models.go` | 10个核心表 |
| **NGAC服务** | `app/service/ngac_service.go` | CRUD + 关联管理 |
| **ABAC引擎** | `app/service/abac_engine.go` | 属性表达式评估 |
| **PDP决策点** | `app/service/pdp_engine.go` | 策略决策引擎 |
| **控制器** | `app/controller/ngac_controller.go` | RESTful API |

### NGAC 工作流程

```
1. 用户请求资源
   ↓
2. PDP (Policy Decision Point) 接收请求
   ├── 获取用户属性 (User Attributes)
   ├── 获取对象属性 (Object Attributes)
   └── 获取上下文属性 (Context Attributes)
   ↓
3. ABAC Engine 评估条件表达式
   ├── 角色继承检查 (Role Inheritance)
   ├── 权限规则匹配 (Permission Matching)
   └── 条件表达式求值 (Condition Evaluation)
   ↓
4. 返回决策结果
   ├── ALLOW - 允许访问
   └── DENY - 拒绝访问
   ↓
5. 记录访问日志 (Access Log)
```

### 权限规则表达式

NGAC 支持灵活的属性条件表达式：

```yaml
# 示例：文章编辑权限规则
permission:
  name: "edit_article"
  effect: "ALLOW"
  subjects:    # 主体条件
    - roles: ["editor", "admin"]
      attributes:
        tenant_id: "{{resource.tenant_id}}"  # 属性引用
  
  objects:    # 对象条件
    - type: "article"
      attributes:
        status: ["draft", "pending_review"]
  
  conditions:  # 额外条件
    - "resource.tenant_id == subject.tenant_id"
    - "subject.roles contains 'admin' or subject.has_article_permission = true"
```

---

## 🌍 多语言 (i18n)

系统支持中英文切换：

### 后端

```go
// app/i18n/ 目录下
├── zh-CN.toml    # 中文翻译
└── en-US.toml    # 英文翻译
```

### 前端

```typescript
// frontend/src/locales/ 目录下
├── zh-CN.ts      // 中文翻译
└── en-US.ts      // 英文翻译
```

### 切换语言

```javascript
// 前端切换语言示例
import { useI18n } from 'vue-i18n'
const { locale } = useI18n()

// 切换到中文
locale.value = 'zh-CN'

// 切换到英文
locale.value = 'en-US'
```

---

## 📦 部署

### Docker 部署（推荐）

```dockerfile
# 后端 Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates mysql-client
WORKDIR /root/
COPY --from=builder /server .
COPY configs/config.yaml .
EXPOSE 8084
CMD ["/server"]
```

### 生产环境配置

```yaml
# configs/config.yaml
server:
  mode: release  # 生产模式

log:
  level: warn    # 降低日志级别
  output: file
  file:
    path: /var/log/gocms/app.log

mysql:
  host: your-db-host
  port: 3306
  database: gocms_v3
  
# 启用 HTTPS（使用反向代理）
# nginx / caddy / traefik
```

### Nginx 配置示例

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://127.0.0.1:8084/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:8084/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

---

## 📝 开发指南

### 后端开发

```bash
# 运行服务器（debug模式）
go run ./cmd/server/main.go

# 使用 hot reload 开发
air

# 代码格式化
go fmt ./...
go vet ./...

# 静态分析
golangci-lint run
```

### 前端开发

```bash
cd frontend

# 开发服务器（带热重载）
npm run dev

# 类型检查
npm run type-check

# 生产构建
npm run build

# 预览构建结果
npm run preview
```

### 添加新模块

#### 1. 后端新增 API

```go
// 1. 定义 Model (app/model/)
type Article struct {
    model.BaseModel
    Title   string `gorm:"size:255" json:"title"`
    Content string `gorm:"type:text" json:"content"`
}

// 2. 定义 Service (app/service/)
type ArticleService struct{}

func (s *ArticleService) List(page, pageSize int) ([]model.Article, int64, error) {
    // 业务逻辑
}

// 3. 定义 Controller (app/controller/)
type ArticleController struct{}

func (c *ArticleController) GetList(c *gin.Context) {
    service := service.NewArticleService()
    // 处理请求 & 返回响应
}

// 4. 注册路由 (app/router/router.go)
r.GET("/articles", controller.NewArticleController().GetList)
```

#### 2. 前端新增页面

```typescript
// 1. 创建 API 模块 (src/api/article.ts)
export function getArticleList(params: ArticleListParams) {
  return request.get<ArticleListResponse>('/articles', { params })
}

// 2. 创建页面组件 (src/pages/articles/ArticleListPage.vue)
<template>
  <div>文章列表</div>
</template>

// 3. 添加路由 (src/router/index.ts)
{
  path: 'articles',
  component: () => import('@/pages/articles/ArticleListPage.vue'),
}
```

### 代码规范

- **后端**: 遵循 [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- **前端**: 遵循 [Vue 3 Composition API 风格指南](https://vuejs.org/style-guide/)

---

## 📊 项目统计

| 指标 | 数量 |
|------|------|
| **控制器** | 20+ |
| **服务层** | 18+ |
| **数据模型** | 35+ |
| **API端点** | 80+ |
| **前端页面** | 25+ |
| **前端组件** | 100+ |

---

## 📜 许可证

本项目采用 [MIT](LICENSE) 许可证。

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！请先阅读项目的贡献指南。

---

<div align="center">

**GoCMS v3** - 用 Go 构建的强大内容管理系统 ❤️

Made with ❤️ by GoCMS Team

</div>