# GoCMS v3 文档中心

## 📚 文档目录

### 系统架构

| 文档 | 说明 | 校准状态 |
|------|------|----------|
| [SYSTEM_ARCHITECTURE.md](SYSTEM_ARCHITECTURE.md) | **系统架构与前端开发指引** - 综合文档 (v4.0)，包含技术栈、项目结构(12控制器/16服务/13模型)、完整API参考、NGAC概述等 | ✅ v4.0 已校准 |
| [NGAC_GUIDE.md](NGAC_GUIDE.md) | **NGAC 权限管理完整指南** - RBAC+ABAC 权限系统详细文档 | ✅ 已校准 |

### 数据库

| 文档 | 说明 | 校准状态 |
|------|------|----------|
| [database_migration_and_seeding_guide.md](database_migration_and_seeding_guide.md) | **数据库迁移与数据种子指南** - SQL 迁移文件和数据初始化说明 | ✅ 内容准确 |

### UI/前端

| 文档 | 说明 | 校准状态 |
|------|------|----------|
| [tailadmin_full_migration_plan.md](tailadmin_full_migration_plan.md) | **TailAdmin UI 完整迁移计划** (v4.0) - Tailwind CSS 风格迁移详细方案，含Phase进度(1-5完成/6进行中) | ✅ v4.0 已校准 |

### 性能优化

| 文档 | 说明 | 校准状态 |
|------|------|----------|
| [video_like_query_optimization.md](video_like_query_optimization.md) | **视频 LIKE 查询性能优化方案** - MySQL FULLTEXT / Elasticsearch 优化 | ✅ 内容准确（描述与实际代码一致） |

### API 文档

| 文档 | 说明 |
|------|------|
| [swagger/swagger.yaml](swagger/swagger.yaml) | **Swagger API 规范文件** - RESTful API 文档定义 |

---

## 📋 文档维护说明

### 最终文档结构

```
docs/
├── README.md                          # 本文件 - 文档索引
├── SYSTEM_ARCHITECTURE.md             # 系统架构综合文档（新建）
├── NGAC_GUIDE.md                      # NGAC 权限管理指南（新建）
├── database_migration_and_seeding_guide.md
├── tailadmin_full_migration_plan.md   # TailAdmin 迁移计划（已校准）
├── video_like_query_optimization.md
└── swagger/
    └── swagger.yaml                   # API 规范文件
```

### 校准记录

| 日期 | 操作 | 说明 |
|------|------|------|
| 2026-08-30 | 文档整合 | 合并相似文档，创建统一架构文档和 NGAC 指南 |
| 2026-08-30 | 内容校准 | 验证所有文档内容与当前系统实际一致 |
| 2026-08-30 | TailAdmin计划更新 | 更新Phase完成状态（Phase 1-2已完成） |
| 2026-08-30 | SYSTEM_ARCHITECTURE v4.0 | 更新控制器/服务/模型列表，补充新增API路由(分类/上传/批量/GraphQL)，前端架构补充Locales/Stores/utils |
| 2026-08-30 | tailadmin_full v4.0 | 更新技术栈版本，Phase进度(1-5完成/6进行中)，页面迁移7/12，新增后端API路由概览章节 |

### 已归档/删除的文档

以下文档的内容已整合到新文档中，原文件已删除：

| 原文件 | 处理 |
|--------|------|
| `ngac_complete_implementation.md` | → 整合至 NGAC_GUIDE.md |
| `ngac_implementation_guide.md` | → 整合至 NGAC_GUIDE.md |
| `ngac_module_readme.md` | → 整合至 NGAC_GUIDE.md |
| `backend_frontend_comprehensive_analysis.md` | → 整合至 SYSTEM_ARCHITECTURE.md |
| `管理系统分析报告_前端开发指引.md` | → 整合至 SYSTEM_ARCHITECTURE.md |
| `tailadmin_migration_plan.md` | → 与 tailadmin_full 合并 |
| `tailadmin_ui_migration_complete.md` | → 与 tailadmin_full 合并 |
| `optimization_implementation_guide.md` | 已归档（示例代码，功能未实现） |
| `comprehensive_tool_test_report.txt` | 已归档（测试报告） |
| `tool_test_result.txt` | 已归档（测试结果） |

---

*最后更新：2026-08-30 (v4.0 全面校准)*
