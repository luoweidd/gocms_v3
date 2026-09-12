package main

import (
	"fmt"
	"log"

	"gocms_v3/app/config"
	"gocms_v3/app/db"
	"gocms_v3/app/migration"
	"gocms_v3/app/model"
	"gocms_v3/app/model/ngac"
	"gocms_v3/app/router"
)

// 注: migration包仍用于数据种子(seeder)，SQL迁移逻辑已清理

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	if err := db.InitMySQL(cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis
	if err := db.InitRedis(cfg.Redis); err != nil {
		log.Printf("警告: Redis初始化失败 [%s:%d], 继续运行: %v", cfg.Redis.Host, cfg.Redis.Port, err)
	}

	// ============================================
	// 步骤1: GORM AutoMigrate 同步表结构
	// ============================================
	log.Println("========================================")
	log.Println("[STEP 1] 开始同步表结构...")
	log.Println("========================================")

	db.GetDB().AutoMigrate(
		// 核心模型
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Menu{},
		&model.Article{},
		&model.ArticleCategory{},
		&model.Tag{},
		&model.Video{},
		&model.VideoCategory{},
		&model.Comment{},
		&model.OperationLog{},
		&model.SystemConfig{},
		&model.DataDict{},
		// 租户模型
		&model.Tenant{},
		// 媒体模型
		&model.MediaAsset{},
		&model.Album{},
		// 文件上传
		&model.FileUpload{},
		// OAuth
		&model.OAuthAuthCode{},
		// API签名
		&model.APISignKey{},
		// SEO与站点地图
		&model.SeoSetting{},
		&model.SitemapConfig{},
		// 备份
		&model.BackupRecord{},
		// 审核
		&model.ContentAudit{},
		// 消息与通知
		&model.SystemMessage{},
		&model.UserNotification{},
		&model.MessageCategory{},
		// NGAC模型
		&ngac.Role{},
		&ngac.Permission{},
		&ngac.RolePermission{},
		&ngac.UserRole{},
		&ngac.ObjectClass{},
		&ngac.ObjectAttribute{},
		&ngac.UserAttribute{},
		&ngac.ContextAttribute{},
		&ngac.AccessLog{},
		&ngac.PDPEngineConfig{},
		// 标签关联表（many2many join table，由 GORM 自动管理，不需要 AutoMigrate）
	)
	log.Println("[STEP 1] 表结构同步完成")

	// ============================================
	// 步骤2: 初始化数据种子
	// ============================================
	log.Println("========================================")
	log.Println("[STEP 2] 开始初始化数据种子...")
	log.Println("========================================")

	seeder := migration.NewSeeder(db.GetDB())
	if err := seeder.Seed(); err != nil {
		log.Printf("警告: 数据种子初始化失败: %v", err)
	}

	// ============================================
	// 启动服务器
	// ============================================
	r := router.NewRouter(cfg)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Println("========================================")
	log.Println("=== GoCMS V3 系统启动成功 ===")
	log.Printf("=== 监听地址: http://localhost%s ===", addr)
	log.Println("========================================")
	log.Printf("[INFO] 默认管理员账号: admin / admin123")

	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
