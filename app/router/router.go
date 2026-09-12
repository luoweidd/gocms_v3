package router

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"gocms_v3/app/config"
	"gocms_v3/app/controller"
	"gocms_v3/app/db"
	"gocms_v3/app/handler"
	"gocms_v3/app/middleware"
	"gocms_v3/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NewRouter 创建路由
func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// ========== 跨域配置（必须最先注册） ==========
	corsCfg := middleware.CORSConfig{
		AllowOrigins:     cfg.Cors.AllowOrigins,
		AllowMethods:     cfg.Cors.AllowMethods,
		AllowHeaders:     cfg.Cors.AllowHeaders,
		ExposeHeaders:    cfg.Cors.ExposeHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		MaxAge:           cfg.Cors.MaxAge,
	}

	// 全局注册CORS中间件（必须在所有路由之前）
	r.Use(middleware.CORS(corsCfg))

	// ========== 静态文件服务配置 ==========
	var uploadDir string
	execPath, err := os.Getwd()
	if err == nil {
		uploadDir = filepath.Join(execPath, "storage", "uploads")
		videoDir := filepath.Join(execPath, "storage", "videos")

		// 确保目录存在
		os.MkdirAll(uploadDir, 0755)
		os.MkdirAll(videoDir, 0755)

		fmt.Printf("[STATIC] Upload directory: %s\n", uploadDir)
		fmt.Printf("[STATIC] Video directory: %s\n", videoDir)

		// 配置上传文件静态服务（无 /api 前缀，直接由浏览器访问）
		r.Static("/uploads", uploadDir)
		// 配置视频文件静态服务
		r.Static("/videos", videoDir)
		r.Static("/storage/videos", videoDir)
	}

	// ========== 公开路由组（无需认证） ==========
	public := r.Group("/api")
	{
		// 处理 URL 编码文件名的静态服务（如空格 %20）- 必须在 public 和 uploadDir 都可用后注册
		if err == nil {
			public.GET("/storage/uploads/*filePath", func(c *gin.Context) {
				// c.Param("filePath") 返回 "/filename.png"（包含前导 /）
				filePath := c.Param("filePath")

				// 移除前导 /
				if len(filePath) > 0 && filePath[0] == '/' {
					filePath = filePath[1:]
				}

				// URL 解码文件名（处理空格 %20 等特殊字符）
				decodedPath, unescapeErr := url.QueryUnescape(filePath)
				if unescapeErr != nil {
					decodedPath = filePath
				}

				// 确保路径正确拼接（使用 / 作为分隔符）
				fullPath := filepath.Join(uploadDir, filepath.ToSlash(decodedPath))

				// 验证文件是否存在
				_, statErr := os.Stat(fullPath)
				if statErr != nil || os.IsNotExist(statErr) {
					c.JSON(404, gin.H{"error": "file not found", "path": fullPath})
					c.Abort()
					return
				}

				// 直接返回文件
				http.ServeFile(c.Writer, c.Request, fullPath)
				c.Abort()
			})
		}

		// 用户认证相关（登录/注册/登出）
		public.POST("/auth/login", controller.NewAuthController().Login)
		public.POST("/auth/register", controller.NewAuthController().Register)
		public.POST("/auth/logout", controller.NewAuthController().Logout)

		// 公开视频接口
		public.GET("/videos/published", controller.NewVideoController().GetPublished)
		public.GET("/videos/:id", controller.NewVideoController().GetDetail)

		// 公开文章接口
		public.GET("/articles/:id", controller.NewArticleController().GetDetail)

		// 公开视频分类树（文章分类需要认证）
		public.GET("/video-categories/tree", controller.NewVideoCategoryController().GetTree)
	}

	// ========== API 前缀统一配置 ==========
	// 所有API路由统一使用 /api 前缀，与前端配置保持一致

	// ========== 认证路由组（需要 JWT 认证） ==========
	authed := r.Group("/api")
	authed.Use(middleware.JWTAuth())
	authed.Use(middleware.OperationLog())
	{
		// 用户管理 - P0 修复：统一使用复数 URL
		userController := controller.NewUserController()
		authed.GET("/user/info", userController.GetUserInfo)
		authed.PUT("/user/info", userController.UpdateUser)
		authed.GET("/users", userController.ListUsers)
		authed.POST("/users", userController.CreateUser)
		authed.GET("/users/:id", userController.GetByID)       // P1 修复：用户详情
		authed.DELETE("/users/:id", userController.DeleteUser) // P1 修复：用户删除

		// 文章管理 - P0 修复：统一使用复数 URL
		authed.POST("/articles", controller.NewArticleController().Create)
		authed.PUT("/articles/:id", controller.NewArticleController().Update)
		authed.DELETE("/articles/:id", controller.NewArticleController().Delete)
		authed.GET("/articles", controller.NewArticleController().GetList)

		// 文章分类管理 - 固定段路由必须在 :id 之前注册
		authed.POST("/article-categories", controller.NewArticleCategoryController().Create)
		authed.GET("/article-categories/tree", controller.NewArticleCategoryController().GetTree)
		authed.GET("/article-categories", controller.NewArticleCategoryController().GetList)
		authed.GET("/article-categories/:id", controller.NewArticleCategoryController().GetByID)
		authed.PUT("/article-categories/:id", controller.NewArticleCategoryController().Update)
		authed.DELETE("/article-categories/:id", controller.NewArticleCategoryController().Delete)

		// 视频管理 - P0 修复：统一使用复数 URL
		authed.POST("/videos", controller.NewVideoController().Create)
		authed.PUT("/videos/:id", controller.NewVideoController().Update)
		authed.DELETE("/videos/:id", controller.NewVideoController().Delete)
		authed.GET("/videos", controller.NewVideoController().GetList)

		// 视频分类管理 - 修复路由顺序
		authed.POST("/video-categories", controller.NewVideoCategoryController().Create)
		authed.GET("/video-categories", controller.NewVideoCategoryController().GetTree)
		authed.GET("/video-categories/:id", controller.NewVideoCategoryController().GetByID)
		authed.PUT("/video-categories/:id", controller.NewVideoCategoryController().Update)
		authed.DELETE("/video-categories/:id", controller.NewVideoCategoryController().Delete)

		// 视频管理 - 发布/置顶/评分操作
		authed.POST("/video/:id/publish", controller.NewVideoController().Publish)
		authed.POST("/video/:id/unpublish", controller.NewVideoController().Unpublish)
		authed.POST("/video/:id/top", controller.NewVideoController().Top)
		authed.POST("/video/:id/untop", controller.NewVideoController().UnTop)
		authed.GET("/video/:id/rate", controller.NewVideoController().Rate)

		// 菜单管理 - 注意：固定段路由必须在 :id 路由之前注册
		authed.GET("/menus/tree", controller.NewMenuController().GetTree)
		authed.GET("/menus/permissions", controller.NewMenuController().GetPermissions)
		authed.POST("/menus", controller.NewMenuController().Create)
		authed.GET("/menus", controller.NewMenuController().GetList)
		authed.GET("/menus/:id", controller.NewMenuController().GetByID)
		authed.PUT("/menus/:id", controller.NewMenuController().Update)
		authed.DELETE("/menus/:id", controller.NewMenuController().Delete)

		// Dashboard 仪表盘
		dashboardCtrl := controller.NewDashboardController()
		authed.GET("/dashboard/stats", dashboardCtrl.GetStats)
		authed.GET("/dashboard/overview", dashboardCtrl.GetOverview)
		authed.GET("/dashboard/trend", dashboardCtrl.GetTrend)
		authed.GET("/dashboard/recent-activity", dashboardCtrl.GetRecentActivity)
		authed.GET("/dashboard/resource-usage", dashboardCtrl.GetResourceUsage)

		// 评论管理
		commentController := controller.NewCommentController(service.NewCommentService())
		authed.POST("/comments", commentController.Create)
		authed.PUT("/comments/:id/approve", commentController.Approve)
		authed.GET("/comments", commentController.GetList)
		authed.GET("/comments/stats", commentController.GetStats)
		authed.DELETE("/comments/:id", commentController.Delete)

		// 评论批量操作
		authed.DELETE("/comments/batch-delete", commentController.BatchDelete)
		authed.POST("/comments/batch-approve", commentController.BatchApprove)
		authed.POST("/comments/batch-reject", commentController.BatchReject)
		authed.GET("/comments/audit-stats", commentController.GetAuditStats)

		// 文件分片上传
		uploadController := controller.NewUploadController()
		authed.POST("/upload/initiate", uploadController.InitiateChunkUpload)
		authed.POST("/upload/chunk", uploadController.UploadChunk)
		authed.POST("/upload/complete", uploadController.CompleteUpload)
		authed.POST("/upload/register", uploadController.RegisterFile)

		// 批量操作 - P0 修复：实现批量创建
		authed.POST("/articles/batch", controller.NewBatchController().BatchCreateArticles)
		authed.POST("/videos/batch", controller.NewBatchController().BatchCreateVideos)

		// 权限管理
		authed.GET("/auth/permissions", controller.NewAuthController().GetUserPermissions)

		// 租户管理
		tenantController := controller.NewTenantController()
		authed.POST("/tenants", tenantController.CreateTenant)
		authed.GET("/tenants", tenantController.GetTenantList)
		authed.GET("/tenants/stats", tenantController.GetTenantStats)
		authed.GET("/tenants/:id", tenantController.GetTenant)
		authed.PUT("/tenants/:id", tenantController.UpdateTenant)
		authed.DELETE("/tenants/:id", tenantController.DeleteTenant)
		authed.PUT("/tenants/:id/status", tenantController.UpdateTenantStatus)
		authed.POST("/tenants/:id/users", tenantController.AddUserToTenant)
		authed.GET("/tenants/:id/users", tenantController.GetTenantUsers)
		authed.DELETE("/tenants/:id/users/:user_id", tenantController.RemoveUserFromTenant)

		// P0 修复：补充角色管理路由
		roleController := controller.NewRoleController()
		authed.GET("/roles", roleController.ListRoles)                    // 角色列表
		authed.GET("/roles/:id", roleController.GetRole)                  // 角色详情
		authed.POST("/roles", roleController.CreateRole)                  // 创建角色
		authed.PUT("/roles/:id", roleController.UpdateRole)               // 更新角色
		authed.DELETE("/roles/:id", roleController.DeleteRole)            // 删除角色
		authed.GET("/permissions/tree", roleController.GetPermissionTree) // P2 修复：权限树（前端 getPermissionTree 调用）

		// 操作日志管理
		operationLogCtrl := controller.NewOperationLogController()
		authed.GET("/operation-logs", operationLogCtrl.GetList)
		authed.GET("/operation-logs/stats", operationLogCtrl.GetStats)
		authed.GET("/operation-logs/:id", operationLogCtrl.GetByID)
		authed.DELETE("/operation-logs/:id", operationLogCtrl.Delete)
		authed.POST("/operation-logs/batch-delete", operationLogCtrl.BatchDelete)
		authed.DELETE("/operation-logs/expired", operationLogCtrl.CleanExpired)
		authed.DELETE("/operation-logs/clean-all", operationLogCtrl.CleanAll)

		// P0 修复：补充标签管理路由
		tagController := controller.NewTagController()
		authed.GET("/tags", tagController.ListTags)         // 标签列表
		authed.GET("/tags/:id", tagController.GetTag)       // 标签详情
		authed.POST("/tags", tagController.CreateTag)       // 创建标签
		authed.PUT("/tags/:id", tagController.UpdateTag)    // 更新标签
		authed.DELETE("/tags/:id", tagController.DeleteTag) // 删除标签

		// 内容审核管理
		auditController := controller.NewAuditController()
		authed.POST("/audit/submit", auditController.SubmitAudit)
		authed.GET("/audit/list", auditController.GetAuditList)
		authed.GET("/audit/:id", auditController.GetAuditByID)
		authed.PUT("/audit/:id/approve", auditController.Approve)
		authed.PUT("/audit/:id/reject", auditController.Reject)
		authed.PUT("/audit/:id/needs-modification", auditController.NeedsModification)
		authed.POST("/audit/batch", auditController.BatchAudit)
		authed.GET("/audit/stats", auditController.GetAuditStats)

		// 媒体中心管理
		mediaController := controller.NewMediaController()
		authed.POST("/media/assets", mediaController.CreateMediaAsset)
		authed.GET("/media/assets", mediaController.GetMediaAssetList)
		authed.GET("/media/assets/stats", mediaController.GetMediaStats)
		authed.GET("/media/assets/:id", mediaController.GetMediaAssetByID)
		authed.PUT("/media/assets/:id", mediaController.UpdateMediaAsset)
		authed.DELETE("/media/assets/:id", mediaController.DeleteMediaAsset)
		authed.DELETE("/media/assets/batch", mediaController.BatchDeleteMediaAssets)

		// 相册管理
		authed.POST("/media/albums", mediaController.CreateAlbum)
		authed.GET("/media/albums", mediaController.GetAlbumList)
		authed.GET("/media/albums/:id", mediaController.GetAlbumByID)
		authed.PUT("/media/albums/:id", mediaController.UpdateAlbum)
		authed.DELETE("/media/albums/:id", mediaController.DeleteAlbum)
		authed.POST("/media/albums/:id/assets", mediaController.AddAssetToAlbum)
		authed.DELETE("/media/albums/:id/assets/:asset_id", mediaController.RemoveAssetFromAlbum)

		// 系统消息与通知
		notificationService := service.NewNotificationService(db.GetDB())
		messageController := controller.NewMessageController(notificationService)
		authed.POST("/messages", messageController.PublishMessage)
		authed.GET("/messages", messageController.GetMessages)
		authed.GET("/messages/:id", messageController.GetMessageDetail)
		authed.DELETE("/messages/:id", messageController.DeleteMessage)
		authed.GET("/messages/my", messageController.GetMyNotifications)
		authed.GET("/messages/my/stats", messageController.GetMyNotificationStats)
		authed.PUT("/messages/:id/read", messageController.MarkAsRead)
		authed.POST("/messages/my/all-read", messageController.MarkAllAsRead)
		authed.POST("/messages/batch-read", messageController.BatchMarkAsRead)
		authed.POST("/messages/:id/withdraw", messageController.WithdrawMessage)

		// 消息类别管理
		messageCategoryService := service.NewMessageCategoryService(db.GetDB())
		messageCategoryController := controller.NewMessageCategoryController(messageCategoryService)
		authed.GET("/message-categories", messageCategoryController.GetCategories)
		authed.GET("/message-categories/:id", messageCategoryController.GetCategoryByID)
		authed.POST("/message-categories", messageCategoryController.CreateCategory)
		authed.PUT("/message-categories/:id", messageCategoryController.UpdateCategory)
		authed.DELETE("/message-categories/:id", messageCategoryController.DeleteCategory)

		// 数据备份管理
		backupController := controller.NewBackupController()
		backupService := service.NewBackupService(db.GetDB(), config.GlobalConfig)
		backupController.SetService(backupService)
		authed.POST("/backup/full", backupController.CreateBackup)
		authed.GET("/backup/list", backupController.GetBackupList)
		authed.GET("/backup/stats", backupController.GetBackupStats)
		authed.GET("/backup/:id", backupController.GetBackupDetail)
		authed.GET("/backup/:id/download", backupController.DownloadBackup)
		authed.DELETE("/backup/:id", backupController.DeleteBackup)
		authed.POST("/backup/:id/restore", backupController.RestoreBackup)

		// SEO管理
		seoController := controller.NewSeoController()
		seoService := service.NewSeoService(db.GetDB())
		seoController.SetService(seoService)
		authed.POST("/seo/settings", seoController.CreateSeoSetting)
		authed.PUT("/seo/settings/:id", seoController.UpdateSeoSetting)
		authed.GET("/seo/settings", seoController.GetSeoList)
		authed.GET("/seo/settings/:id", seoController.GetSeoDetail)
		authed.DELETE("/seo/settings/:id", seoController.DeleteSeoSetting)
		authed.GET("/seo/sitemap-configs", seoController.GetSitemapConfigs)
		authed.POST("/seo/sitemap/generate", seoController.GenerateSitemap)
		authed.GET("/seo/sitemap-stats", seoController.GetSitemapStats)

		// NGAC 权限管理路由
		ngacController := controller.NewNgacController()

		// 角色管理
		// 注意：更具体的路由（带子路径）必须在通配符路由之前注册
		authed.POST("/ngac/roles", ngacController.CreateRole)
		authed.GET("/ngac/roles", ngacController.ListRoles)

		// 角色权限管理 - 必须在 /ngac/roles/:role_id 之前注册
		authed.POST("/ngac/roles/:role_id/permissions", ngacController.AssignPermissionsToRole)
		authed.DELETE("/ngac/roles/:role_id/permissions/:perm_id", ngacController.RemovePermissionFromRole)

		authed.PUT("/ngac/roles/:role_id", ngacController.UpdateRole)
		authed.DELETE("/ngac/roles/:role_id", ngacController.DeleteRole)

		// 权限管理
		authed.GET("/ngac/permissions", ngacController.ListPermissions)
		authed.POST("/ngac/permissions", ngacController.CreatePermission)
		authed.PUT("/ngac/permissions/:perm_id", ngacController.UpdatePermission)
		authed.DELETE("/ngac/permissions/:perm_id", ngacController.DeletePermission)

		// 用户角色管理
		authed.POST("/ngac/users/:user_id/roles", ngacController.AssignRoles)
		authed.GET("/ngac/users/:user_id/roles", ngacController.GetUserRoles)
		authed.GET("/ngac/users/:user_id/permissions", ngacController.GetUserPermissions)

		// 对象类管理
		authed.GET("/ngac/object-classes/tree", ngacController.GetObjectClassTree)

		// 属性管理
		authed.GET("/ngac/users/:user_id/attributes", ngacController.GetSubjectAttributes)
		authed.GET("/ngac/context/attributes", ngacController.GetContextAttributes)

		// 授权接口
		authed.POST("/ngac/authz", ngacController.Authz)

		// P1 修复：接入 GraphQL 路由
		authed.POST("/graphql", func(c *gin.Context) {
			var req struct {
				Query     string                 `json:"query" binding:"required"`
				Variables map[string]interface{} `json:"variables"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误：" + err.Error()})
				return
			}

			result, err := handler.NewGraphQLHandler(
				service.NewArticleService(),
				service.NewVideoService(),
				service.NewUserService(),
				service.NewCommentService(),
			).ExecuteGraphQL(req.Query, req.Variables)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"data": result})
		})
	}

	return r
}
