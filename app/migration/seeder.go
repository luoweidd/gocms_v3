package migration

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gocms_v3/app/domain/entity"
	"gocms_v3/app/model"
	"gocms_v3/app/model/ngac"

	"gorm.io/gorm"
)

// Seeder 数据种子初始化器
type Seeder struct {
	db *gorm.DB
}

// SeedContentAudits 为现有内容批量生成审核记录
func (s *Seeder) SeedContentAudits() error {
	log.Println("[SEEDER] 开始为现有内容批量生成审核记录...")

	// 检查 content_audits 表是否存在
	if !s.db.Migrator().HasTable("content_audits") {
		log.Println("[SEEDER] content_audits 表不存在，跳过")
		return nil
	}

	affectedRows := 0
	maxRows := 100

	// 为现有文章生成审核记录（只创建没有审核记录的）
	func() {
		var articles []model.Article
		if err := s.db.Where("id NOT IN (SELECT content_id FROM content_audits WHERE content_type = 'article')").Limit(100).Find(&articles).Error; err != nil {
			log.Printf("[SEEDER] 查询未审核文章失败: %v", err)
			return
		}
		for _, article := range articles {
			if affectedRows >= maxRows {
				break
			}
			auditNo := s.generateAuditNo("article")
			audit := model.ContentAudit{
				AuditNo:        auditNo,
				ContentType:    "article",
				ContentID:      article.ID,
				ContentTitle:   article.Title,
				SubmitterID:    article.AuthorID,
				SubmitterName:  fmt.Sprintf("用户%d", article.AuthorID),
				Status:         0, // 待审核
				PreviousStatus: -1,
				NextStatus:     0,
			}
			if err := s.db.Create(&audit).Error; err != nil {
				log.Printf("[SEEDER] 为文章 %d 创建审核记录失败: %v", article.ID, err)
			} else {
				log.Printf("[SEEDER] 为文章 %d (%s) 创建审核记录: %s", article.ID, article.Title, auditNo)
				affectedRows++
			}
		}
	}()

	// 为现有视频生成审核记录
	func() {
		var videos []model.Video
		if err := s.db.Where("id NOT IN (SELECT content_id FROM content_audits WHERE content_type = 'video')").Limit(100).Find(&videos).Error; err != nil {
			log.Printf("[SEEDER] 查询未审核视频失败: %v", err)
			return
		}
		for _, video := range videos {
			if affectedRows >= maxRows {
				break
			}
			auditNo := s.generateAuditNo("video")
			audit := model.ContentAudit{
				AuditNo:        auditNo,
				ContentType:    "video",
				ContentID:      video.ID,
				ContentTitle:   video.Title,
				SubmitterID:    video.AuthorID,
				SubmitterName:  fmt.Sprintf("用户%d", video.AuthorID),
				Status:         0, // 待审核
				PreviousStatus: -1,
				NextStatus:     0,
			}
			if err := s.db.Create(&audit).Error; err != nil {
				log.Printf("[SEEDER] 为视频 %d 创建审核记录失败: %v", video.ID, err)
			} else {
				log.Printf("[SEEDER] 为视频 %d (%s) 创建审核记录: %s", video.ID, video.Title, auditNo)
				affectedRows++
			}
		}
	}()

	log.Println("[SEEDER] 为媒体资源创建审核记录...")

	// 为现有媒体资产生成审核记录（需要临时禁用外键检查）
	func() {
		// 临时禁用外键检查
		s.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		defer s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")

		var mediaAssets []model.MediaAsset
		if err := s.db.Where("id NOT IN (SELECT content_id FROM content_audits WHERE content_type = 'image')").Limit(100).Find(&mediaAssets).Error; err != nil {
			log.Printf("[SEEDER] 查询未审核媒体资产失败: %v", err)
			return
		}
		for idx, asset := range mediaAssets {
			if affectedRows >= maxRows {
				break
			}
			// 为每个媒体资源生成唯一的审核编号（使用索引后缀）
			auditNo := s.generateAuditNo("media")
			if idx > 0 {
				// 如果不是第一条，手动添加索引后缀以避免重复
				auditNo = auditNo[:len(auditNo)-4] + fmt.Sprintf("%04d", idx+1)
			}
			audit := model.ContentAudit{
				AuditNo:        auditNo,
				ContentType:    asset.FileType, // image, video, file 等
				ContentID:      uint(asset.ID),
				ContentTitle:   asset.Name,
				SubmitterID:    uint(asset.UserID),
				SubmitterName:  "用户1", // 默认用户
				Status:         1,     // 已通过（种子数据直接通过）
				PreviousStatus: -1,
				NextStatus:     1,
				AuditComment:   "系统种子数据，自动通过",
			}
			if err := s.db.Create(&audit).Error; err != nil {
				log.Printf("[SEEDER] ⚠ 为媒体资源 %d (%s) 创建审核记录失败: %v", asset.ID, asset.Name, err)
			} else {
				log.Printf("[SEEDER] ✓ 为媒体资源 %d (%s) 创建审核记录: %s", asset.ID, asset.Name, auditNo)
				affectedRows++
			}
		}
	}()

	log.Printf("[SEEDER] 批量生成审核记录完成，共新增 %d 条记录", affectedRows)
	return nil
}

// generateAuditNo 生成审核编号
func (s *Seeder) generateAuditNo(contentType string) string {
	now := time.Now()
	dateStr := now.Format("20060102")
	prefix := "AUD" + dateStr

	// 尝试获取同类型+同日期的最大编号
	var maxNo string
	query := s.db.Where("audit_no LIKE ?", prefix+"%")
	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}
	query.Order("audit_no DESC").Pluck("audit_no", &maxNo)

	if maxNo != "" {
		suffix := maxNo[len(prefix):]
		if num := parseSuffixToInt(suffix); num > 0 {
			return prefix + fmt.Sprintf("%04d", num+1)
		}
	}
	return prefix + "0001"
}

// parseSuffixToInt 解析编号后缀
func parseSuffixToInt(s string) int {
	num := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		} else {
			break
		}
	}
	return num
}

// seedSystemMessages 初始化系统消息和通知数据（用于测试）
func (s *Seeder) seedSystemMessages() error {
	if !s.db.Migrator().HasTable("system_messages") {
		log.Println("[SEEDER] ⊘ system_messages 表不存在，跳过")
		return nil
	}

	log.Println("[SEEDER] 检查并初始化系统消息数据...")

	// 获取admin用户ID
	var adminUserID uint
	s.db.Model(&model.User{}).Where("username = ?", "admin").Pluck("id", &adminUserID)
	if adminUserID == 0 {
		log.Println("[SEEDER] ⊘ 未找到admin用户，跳过消息初始化")
		return nil
	}

	// 获取所有活跃用户ID
	var activeUserIDs []uint
	s.db.Model(&model.User{}).Where("status = ?", 1).Pluck("id", &activeUserIDs)
	if len(activeUserIDs) == 0 {
		activeUserIDs = []uint{adminUserID}
	}

	// 检查是否已有系统消息数据
	var msgCount int64
	s.db.Model(&model.SystemMessage{}).Count(&msgCount)

	// 测试消息数据
	testMessages := []struct {
		Title       string
		Content     string
		MessageType string
		Priority    int8
	}{
		{
			Title:       "系统维护通知",
			Content:     "亲爱的用户，我们将于2026年9月15日凌晨2点至6点进行系统维护升级。届时部分功能可能暂时不可用，给您带来的不便敬请谅解。",
			MessageType: "system",
			Priority:    1,
		},
		{
			Title:       "新功能上线公告",
			Content:     "媒体中心2.0版本正式上线！新增视频分类管理、图片批量上传等功能，快来体验吧。",
			MessageType: "publish",
			Priority:    1,
		},
		{
			Title:       "审核通过通知",
			Content:     "您的文章《Go语言微服务架构实践》已通过审核并正式发布，感谢你的贡献！",
			MessageType: "audit_passed",
			Priority:    0,
		},
		{
			Title:       "评论回复通知",
			Content:     "用户张三回复了你的评论：'这个教程非常详细，对我帮助很大！'",
			MessageType: "comment_reply",
			Priority:    0,
		},
		{
			Title:       "审核驳回通知",
			Content:     "您的文章《新技术趋势分析》因内容不完整被驳回，请补充完整后重新提交审核。",
			MessageType: "audited_failed",
			Priority:    1,
		},
	}

	messageIDs := make([]uint, 0)

	// 如果没有数据，则创建测试消息
	if msgCount == 0 {
		log.Println("[SEEDER] 系统消息为空，创建测试消息...")
		for i, msg := range testMessages {
			systemMsg := model.SystemMessage{
				Title:         msg.Title,
				Content:       msg.Content,
				MessageType:   msg.MessageType,
				SenderID:      adminUserID,
				SenderName:    "admin",
				Priority:      msg.Priority,
				TargetType:    "all",
				TargetUserIDs: encodeTargetUserIDsForSeeder(activeUserIDs),
				IsReadDefault: 0,
				Status:        1,
			}
			systemMsg.CreatedAt = time.Now().Add(-time.Duration(i) * 24 * time.Hour)

			if err := s.db.Create(&systemMsg).Error; err != nil {
				log.Printf("[SEEDER] ⚠ 创建消息失败: %s, error: %v", msg.Title, err)
				continue
			}
			messageIDs = append(messageIDs, systemMsg.ID)
			log.Printf("[SEEDER] ✓ 创建消息: %s (类型: %s)", msg.Title, msg.MessageType)
		}
	} else {
		log.Printf("[SEEDER] 系统消息已存在 (%d 条)，检查是否需补充...", msgCount)
		// 获取已有消息ID
		var existingIDs []uint
		s.db.Model(&model.SystemMessage{}).Pluck("id", &existingIDs)
		messageIDs = existingIDs

		// 为每条已有消息确保有用户通知记录
		for _, msgID := range messageIDs {
			var notifCount int64
			s.db.Model(&model.UserNotification{}).Where("message_id = ?", msgID).Count(&notifCount)
			if notifCount == 0 {
				// 为所有活跃用户创建通知记录
				for _, userID := range activeUserIDs {
					notification := model.UserNotification{
						UserID:    userID,
						MessageID: msgID,
						IsRead:    0,
					}
					s.db.Create(&notification)
				}
				log.Printf("[SEEDER] ✓ 为消息 %d 补充了 %d 条用户通知", msgID, len(activeUserIDs))
			}
		}
	}

	// 确保 admin 用户有足够的未读通知用于测试
	var unreadCount int64
	s.db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = ?", adminUserID, 0).Count(&unreadCount)
	if unreadCount == 0 {
		log.Println("[SEEDER] ⊘ admin 用户没有未读通知，为所有消息创建未读通知...")
		for _, msgID := range messageIDs {
			var existingNotif model.UserNotification
			err := s.db.Where("user_id = ? AND message_id = ?", adminUserID, msgID).First(&existingNotif).Error
			if err != nil {
				notification := model.UserNotification{
					UserID:    adminUserID,
					MessageID: msgID,
					IsRead:    0,
				}
				s.db.Create(&notification)
			}
		}
	}

	log.Printf("[SEEDER] ✓ 系统消息处理完成，共有 %d 条消息", len(messageIDs))
	return nil
}

// encodeTargetUserIDsForSeeder 将用户ID数组编码为JSON字符串（种子数据专用）
func encodeTargetUserIDsForSeeder(userIDs []uint) string {
	data, err := json.Marshal(userIDs)
	if err != nil {
		return "[]"
	}
	return string(data)
}

// NewSeeder 创建数据种子初始化器
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed 执行所有数据初始化
func (s *Seeder) Seed() error {
	log.Println("[SEEDER] 开始初始化数据...")

	// 检查并创建需要的表
	if err := s.ensureTablesExist(); err != nil {
		return fmt.Errorf("确保表存在失败: %w", err)
	}

	// 自动迁移 SEO 相关表（如果不存在）
	if err := s.autoMigrateSeoTables(); err != nil {
		return fmt.Errorf("迁移SEO表失败: %w", err)
	}

	// 按依赖顺序初始化数据
	// 1. 先初始化基础数据（用户、角色、权限）
	if err := s.seedAdminUser(); err != nil {
		return fmt.Errorf("初始化超管用户失败: %w", err)
	}

	if err := s.seedRoles(); err != nil {
		return fmt.Errorf("初始化角色失败: %w", err)
	}

	if err := s.seedNgacRoles(); err != nil {
		return fmt.Errorf("初始化NGAC角色失败: %w", err)
	}

	if err := s.seedPermissions(); err != nil {
		return fmt.Errorf("初始化权限失败: %w", err)
	}

	if err := s.seedRolePermissions(); err != nil {
		return fmt.Errorf("初始化角色权限关联失败: %w", err)
	}

	if err := s.seedUserRoles(); err != nil {
		return fmt.Errorf("初始化用户角色关联失败: %w", err)
	}

	// 2. 初始化菜单（依赖角色数据）
	if err := s.seedMenus(); err != nil {
		return fmt.Errorf("初始化菜单失败: %w", err)
	}

	// 3. 初始化业务数据（分类、配置、字典）
	if err := s.seedArticleCategories(); err != nil {
		return fmt.Errorf("初始化文章分类失败: %w", err)
	}

	if err := s.seedVideoCategories(); err != nil {
		return fmt.Errorf("初始化视频分类失败: %w", err)
	}

	if err := s.seedSystemConfigs(); err != nil {
		return fmt.Errorf("初始化系统配置失败: %w", err)
	}

	if err := s.seedDataDicts(); err != nil {
		return fmt.Errorf("初始化数据字典失败: %w", err)
	}

	// 初始化站点地图配置
	if err := s.seedSitemapConfigs(); err != nil {
		return fmt.Errorf("初始化站点地图配置失败: %w", err)
	}

	// 5. 初始化媒体中心数据
	if err := s.seedAlbums(); err != nil {
		return fmt.Errorf("初始化相册数据失败: %w", err)
	}

	// 6. 初始化 NGAC 对象类（放在最后）
	if err := s.seedObjectClasses(); err != nil {
		return fmt.Errorf("初始化NGAC对象类失败: %w", err)
	}

	// 初始化风险关键词库
	if err := s.seedRiskKeywords(); err != nil {
		return fmt.Errorf("初始化风险词库失败: %w", err)
	}

	// 7. 初始化系统消息和通知数据（用于测试）
	if err := s.seedSystemMessages(); err != nil {
		return fmt.Errorf("初始化系统消息失败: %w", err)
	}

	// 8. 初始化默认消息类别
	if err := s.seedDefaultMessageCategories(); err != nil {
		return fmt.Errorf("初始化消息类别失败: %w", err)
	}

	// 9. 创建测试内容（文章、视频）以确保有审核数据
	if err := s.seedTestContent(); err != nil {
		log.Printf("[SEEDER] ⚠ 创建测试内容失败: %v", err)
	}

	// 10. 为现有内容批量生成审核记录（修复审核列表为空的问题）
	if err := s.SeedContentAudits(); err != nil {
		log.Printf("[SEEDER] ⚠ 生成审核记录失败: %v", err)
		// 不阻断启动流程
	}

	log.Println("[SEEDER] 数据初始化完成!")
	return nil
}

// seedTestContent 创建测试内容（文章、视频）并生成审核记录
func (s *Seeder) seedTestContent() error {
	if !s.db.Migrator().HasTable("articles") {
		log.Println("[SEEDER] ⊘ articles 表不存在，跳过测试内容创建")
		return nil
	}
	if !s.db.Migrator().HasTable("videos") {
		log.Println("[SEEDER] ⊘ videos 表不存在，跳过测试内容创建")
		return nil
	}
	if !s.db.Migrator().HasTable("content_audits") {
		log.Println("[SEEDER] ⊘ content_audits 表不存在，跳过测试内容创建")
		return nil
	}

	// 获取 admin 用户 ID
	var adminID uint
	s.db.Model(&model.User{}).Where("username = ?", "admin").Pluck("id", &adminID)
	if adminID == 0 {
		log.Println("[SEEDER] ⊘ 未找到 admin 用户，跳过测试内容创建")
		return nil
	}

	// 检查是否已有文章数据
	var articleCount int64
	s.db.Model(&model.Article{}).Count(&articleCount)
	if articleCount > 0 {
		log.Println("[SEEDER] ⊘ 文章数据已存在，跳过测试内容创建")
		return nil
	}

	log.Println("[SEEDER] 创建测试文章、测试视频及审核记录...")

	// 获取第一个分类 ID
	var cat model.ArticleCategory
	err := s.db.Model(&model.ArticleCategory{}).Order("id ASC").First(&cat).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println("[SEEDER] ⊘ article_categories 表不存在或查询失败，跳过测试内容创建")
		return nil
	}
	if cat.ID == 0 {
		log.Println("[SEEDER] ⊘ 没有文章分类，跳过测试内容创建")
		return nil
	}

	var videoCategoryID uint
	var vcat model.VideoCategory
	err = s.db.Model(&model.VideoCategory{}).Order("id ASC").First(&vcat).Error
	if err == nil {
		videoCategoryID = vcat.ID
	}
	if videoCategoryID == 0 {
		log.Println("[SEEDER] ⊘ 没有视频分类，跳过测试内容创建")
		return nil
	}

	now := time.Now()

	// ============================================
	// 创建测试文章 1（待审核状态）
	// ============================================
	testArticle1 := model.Article{
		Title:      "Go语言微服务架构实践",
		Summary:    "本文详细介绍如何使用Go语言构建微服务架构...",
		Content:    "<h1>Go语言微服务架构实践</h1><p>本文详细介绍如何使用Go语言构建微服务架构...</p>",
		CoverImage: "/uploads/go-microservice.jpg",
		CategoryID: cat.ID,
		AuthorID:   adminID,
		Status:     0, // 待审核
	}
	if err := s.db.Create(&testArticle1).Error; err != nil {
		log.Printf("[SEEDER] ⚠ 创建测试文章 1 失败: %v", err)
	} else {
		// 为文章 1 创建待审核状态的审核记录
		auditNo1 := s.generateAuditNo("article")
		audit1 := model.ContentAudit{
			AuditNo:        auditNo1,
			ContentType:    "article",
			ContentID:      testArticle1.ID,
			ContentTitle:   testArticle1.Title,
			SubmitterID:    adminID,
			SubmitterName:  "admin",
			Status:         0, // 待审核
			PreviousStatus: -1,
			NextStatus:     0,
			AuditedAt:      nil,
		}
		audit1.CreatedAt = now
		audit1.UpdatedAt = now
		if err := s.db.Create(&audit1).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 为文章 1 创建审核记录失败: %v", err)
		} else {
			log.Printf("[SEEDER] ✓ 创建测试文章1: %s (待审核，审核编号: %s)", testArticle1.Title, auditNo1)
		}
	}

	// ============================================
	// 创建测试文章 2（已通过审核）
	// ============================================
	testArticle2 := model.Article{
		Title:      "前端性能优化技巧",
		Summary:    "分享一些实际的前端性能优化技巧...",
		Content:    "<h1>前端性能优化技巧</h1><p>分享一些实际的前端性能优化技巧...</p>",
		CoverImage: "/uploads/frontend-performance.jpg",
		CategoryID: cat.ID,
		AuthorID:   adminID,
		Status:     1, // 已发布
	}
	if err := s.db.Create(&testArticle2).Error; err != nil {
		log.Printf("[SEEDER] ⚠ 创建测试文章 2 失败: %v", err)
	} else {
		// 为文章 2 创建已通过审核状态的审核记录
		auditNo2 := s.generateAuditNo("article")
		auditedAt := now.Add(-2 * time.Hour)
		audit2 := model.ContentAudit{
			AuditNo:        auditNo2,
			ContentType:    "article",
			ContentID:      testArticle2.ID,
			ContentTitle:   testArticle2.Title,
			SubmitterID:    adminID,
			SubmitterName:  "admin",
			AuditorID:      adminID,
			AuditorName:    "admin",
			Status:         1, // 已通过
			PreviousStatus: 0,
			NextStatus:     1,
			AuditComment:   "内容合规，审核通过",
			AuditedAt:      &auditedAt,
		}
		audit2.CreatedAt = now
		audit2.UpdatedAt = now
		if err := s.db.Create(&audit2).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 为文章 2 创建审核记录失败: %v", err)
		} else {
			log.Printf("[SEEDER] ✓ 创建测试文章2: %s (已发布，审核编号: %s)", testArticle2.Title, auditNo2)
		}
	}

	// ============================================
	// 创建测试视频 1（待审核状态）
	// ============================================
	testVideo1 := model.Video{
		Title:      "Go语言从入门到精通",
		Summary:    "完整的Go语言教程，适合初学者...",
		Content:    "<p>完整的Go语言教程...</p>",
		URL:        "/videos/go-tutorial.mp4",
		Cover:      "/uploads/go-tutorial-cover.jpg",
		Format:     "mp4",
		Duration:   3600,
		FileSize:   524288000,
		Width:      1920,
		Height:     1080,
		Bitrate:    5000,
		CategoryID: videoCategoryID,
		AuthorID:   adminID,
		Status:     0, // 待审核
	}
	if err := s.db.Create(&testVideo1).Error; err != nil {
		log.Printf("[SEEDER] ⚠ 创建测试视频 1 失败: %v", err)
	} else {
		// 为视频 1 创建待审核状态的审核记录
		auditNo3 := s.generateAuditNo("video")
		audit3 := model.ContentAudit{
			AuditNo:        auditNo3,
			ContentType:    "video",
			ContentID:      testVideo1.ID,
			ContentTitle:   testVideo1.Title,
			SubmitterID:    adminID,
			SubmitterName:  "admin",
			Status:         0, // 待审核
			PreviousStatus: -1,
			NextStatus:     0,
			AuditedAt:      nil,
		}
		audit3.CreatedAt = now
		audit3.UpdatedAt = now
		if err := s.db.Create(&audit3).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 为视频 1 创建审核记录失败: %v", err)
		} else {
			log.Printf("[SEEDER] ✓ 创建测试视频1: %s (待审核，审核编号: %s)", testVideo1.Title, auditNo3)
		}
	}

	// ============================================
	// 创建测试视频 2（已通过审核）
	// ============================================
	testVideo2 := model.Video{
		Title:      "Vue3进阶教程",
		Summary:    "深入学习Vue3 Composition API...",
		Content:    "<p>深入学习Vue3...</p>",
		URL:        "/videos/vue3-tutorial.mp4",
		Cover:      "/uploads/vue3-cover.jpg",
		Format:     "mp4",
		Duration:   1800,
		FileSize:   262144000,
		Width:      1920,
		Height:     1080,
		Bitrate:    3000,
		CategoryID: videoCategoryID,
		AuthorID:   adminID,
		Status:     1, // 已发布
	}
	if err := s.db.Create(&testVideo2).Error; err != nil {
		log.Printf("[SEEDER] ⚠ 创建测试视频 2 失败: %v", err)
	} else {
		// 为视频 2 创建已通过审核状态的审核记录
		auditNo4 := s.generateAuditNo("video")
		auditedAt := now.Add(-1 * time.Hour)
		audit4 := model.ContentAudit{
			AuditNo:        auditNo4,
			ContentType:    "video",
			ContentID:      testVideo2.ID,
			ContentTitle:   testVideo2.Title,
			SubmitterID:    adminID,
			SubmitterName:  "admin",
			AuditorID:      adminID,
			AuditorName:    "admin",
			Status:         1, // 已通过
			PreviousStatus: 0,
			NextStatus:     1,
			AuditComment:   "内容合规，审核通过",
			AuditedAt:      &auditedAt,
		}
		audit4.CreatedAt = now
		audit4.UpdatedAt = now
		if err := s.db.Create(&audit4).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 为视频 2 创建审核记录失败: %v", err)
		} else {
			log.Printf("[SEEDER] ✓ 创建测试视频2: %s (已发布，审核编号: %s)", testVideo2.Title, auditNo4)
		}
	}

	return nil
}

// ensureTablesExist 确保需要的表存在
func (s *Seeder) ensureTablesExist() error {
	tables := []string{
		"users", "roles", "menus", "system_configs", "data_dicts",
		"article_categories", "video_categories", "operation_logs",
		"user_roles", "ngac_roles", "ngac_permissions", "ngac_role_permissions",
		"ngac_user_roles", "ngac_object_classes",
	}

	for _, table := range tables {
		if !s.db.Migrator().HasTable(table) {
			log.Printf("[SEEDER] ⊘ 表 %s 不存在，跳过相关数据初始化", table)
		}
	}
	return nil
}

// seedAdminUser 初始化超级管理员用户
func (s *Seeder) seedAdminUser() error {
	if !s.db.Migrator().HasTable("users") {
		log.Println("[SEEDER] ⊘ users 表不存在，跳过")
		return nil
	}

	var admin model.User
	result := s.db.Where("username = ?", "admin").First(&admin)
	if result.Error != nil {
		admin = model.User{
			Username: "admin",
			Password: "$2a$10$c6F2wRfxlI.AsqI9zMZGcORxYprg60F3XtWEGIxrwma9/DliH14rK", // admin123
			Nickname: "系统管理员",
			Email:    "admin@example.com",
			Status:   1,
			Roles:    model.JSONString{"super_admin"},
		}
		if err := s.db.Create(&admin).Error; err != nil {
			return fmt.Errorf("创建超管用户失败: %w", err)
		}
		log.Printf("[SEEDER] ✓ 创建超管用户: admin (密码: admin123)")
	} else {
		log.Printf("[SEEDER] ✓ 超管用户已存在: admin")
	}

	return nil
}

// seedRoles 初始化角色（roles 表）
func (s *Seeder) seedRoles() error {
	if !s.db.Migrator().HasTable("roles") {
		log.Println("[SEEDER] ⊘ roles 表不存在，跳过")
		return nil
	}

	roles := []model.Role{
		{ID: 1, Name: "超级管理员", Code: "super_admin", Description: "拥有系统所有权限", DataScope: "all", Status: 1},
		{ID: 2, Name: "内容管理员", Code: "content_admin", Description: "管理文章内容", DataScope: "own", Status: 1},
		{ID: 3, Name: "普通用户", Code: "user", Description: "普通用户权限", DataScope: "own", Status: 1},
	}

	for i := range roles {
		var existing model.Role
		result := s.db.Where("id = ?", roles[i].ID).First(&existing)
		if result.Error != nil {
			s.db.Create(&roles[i])
			log.Printf("[SEEDER] ✓ 创建角色: %s (%s)", roles[i].Name, roles[i].Code)
		} else {
			s.db.Model(&existing).Updates(model.Role{
				Name:        roles[i].Name,
				Code:        roles[i].Code,
				Description: roles[i].Description,
				Status:      roles[i].Status,
			})
			log.Printf("[SEEDER] ✓ 角色已存在: %s (%s)", existing.Name, existing.Code)
		}
	}

	return nil
}

// seedNgacRoles 初始化 NGAC 角色表
func (s *Seeder) seedNgacRoles() error {
	if !s.db.Migrator().HasTable("ngac_roles") {
		log.Println("[SEEDER] ⊘ ngac_roles 表不存在，跳过")
		return nil
	}

	roles := []ngac.Role{
		{ID: 1, Name: "超级管理员", Code: "super_admin", Description: "拥有系统所有权限", Level: 1, Sort: 1, Status: 1},
		{ID: 2, Name: "内容管理员", Code: "content_admin", Description: "管理文章内容", Level: 2, Sort: 2, Status: 1},
		{ID: 3, Name: "普通用户", Code: "user", Description: "普通用户权限", Level: 3, Sort: 3, Status: 1},
	}

	for i := range roles {
		var existing ngac.Role
		result := s.db.Where("id = ?", roles[i].ID).First(&existing)
		if result.Error != nil {
			s.db.Create(&roles[i])
			log.Printf("[SEEDER] ✓ 创建NGAC角色: %s (%s)", roles[i].Name, roles[i].Code)
		} else {
			s.db.Model(&existing).Updates(ngac.Role{
				Name:        roles[i].Name,
				Code:        roles[i].Code,
				Description: roles[i].Description,
				Status:      roles[i].Status,
			})
			log.Printf("[SEEDER] ✓ NGAC角色已存在: %s (%s)", existing.Name, existing.Code)
		}
	}

	return nil
}

// seedPermissions 初始化 NGAC 权限
func (s *Seeder) seedPermissions() error {
	if !s.db.Migrator().HasTable("ngac_permissions") {
		log.Println("[SEEDER] ⊘ ngac_permissions 表不存在，跳过")
		return nil
	}

	objectClassMap := map[string]uint{
		"article":          1,
		"article_category": 2,
		"video":            3,
		"comment":          4,
		"file":             5,
		"gallery":          6,
		"user":             7,
		"role":             8,
		"system_config":    9,
		"menu":             10,
		"operation_log":    11,
		"data_dict":        12,
	}

	permissions := []ngac.Permission{
		{Name: "查看仪表盘", Code: "dashboard:view", ResourceType: "module", Action: "view", SubjectType: "role", Allowed: true, Priority: 1},
		{Name: "文章列表", Code: "content:article:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["article"], Allowed: true, Priority: 10},
		{Name: "文章创建", Code: "content:article:create", ResourceType: "resource", Action: "create", SubjectType: "role", ResourceObjectClassID: objectClassMap["article"], Allowed: true, Priority: 10},
		{Name: "文章更新", Code: "content:article:update", ResourceType: "resource", Action: "update", SubjectType: "role", ResourceObjectClassID: objectClassMap["article"], Allowed: true, Priority: 10},
		{Name: "文章删除", Code: "content:article:delete", ResourceType: "resource", Action: "delete", SubjectType: "role", ResourceObjectClassID: objectClassMap["article"], Allowed: true, Priority: 10},
		{Name: "文章分类列表", Code: "content:category:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["article_category"], Allowed: true, Priority: 20},
		{Name: "文章分类创建", Code: "content:category:create", ResourceType: "resource", Action: "create", SubjectType: "role", ResourceObjectClassID: objectClassMap["article_category"], Allowed: true, Priority: 20},
		{Name: "文章分类更新", Code: "content:category:update", ResourceType: "resource", Action: "update", SubjectType: "role", ResourceObjectClassID: objectClassMap["article_category"], Allowed: true, Priority: 20},
		{Name: "文章分类删除", Code: "content:category:delete", ResourceType: "resource", Action: "delete", SubjectType: "role", ResourceObjectClassID: objectClassMap["article_category"], Allowed: true, Priority: 20},
		{Name: "视频列表", Code: "content:video:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["video"], Allowed: true, Priority: 30},
		{Name: "视频创建", Code: "content:video:create", ResourceType: "resource", Action: "create", SubjectType: "role", ResourceObjectClassID: objectClassMap["video"], Allowed: true, Priority: 30},
		{Name: "视频更新", Code: "content:video:update", ResourceType: "resource", Action: "update", SubjectType: "role", ResourceObjectClassID: objectClassMap["video"], Allowed: true, Priority: 30},
		{Name: "视频删除", Code: "content:video:delete", ResourceType: "resource", Action: "delete", SubjectType: "role", ResourceObjectClassID: objectClassMap["video"], Allowed: true, Priority: 30},
		{Name: "评论列表", Code: "content:comment:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["comment"], Allowed: true, Priority: 40},
		{Name: "评论审核", Code: "content:comment:audit", ResourceType: "resource", Action: "audit", SubjectType: "role", ResourceObjectClassID: objectClassMap["comment"], Allowed: true, Priority: 40},
		{Name: "文件管理", Code: "media:file", ResourceType: "resource", Action: "all", SubjectType: "role", ResourceObjectClassID: objectClassMap["file"], Allowed: true, Priority: 50},
		{Name: "图库管理", Code: "media:gallery", ResourceType: "resource", Action: "all", SubjectType: "role", ResourceObjectClassID: objectClassMap["gallery"], Allowed: true, Priority: 51},
		{Name: "用户列表", Code: "user:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["user"], Allowed: true, Priority: 60},
		{Name: "用户创建", Code: "user:create", ResourceType: "resource", Action: "create", SubjectType: "role", ResourceObjectClassID: objectClassMap["user"], Allowed: true, Priority: 60},
		{Name: "用户更新", Code: "user:update", ResourceType: "resource", Action: "update", SubjectType: "role", ResourceObjectClassID: objectClassMap["user"], Allowed: true, Priority: 60},
		{Name: "用户删除", Code: "user:delete", ResourceType: "resource", Action: "delete", SubjectType: "role", ResourceObjectClassID: objectClassMap["user"], Allowed: true, Priority: 60},
		{Name: "角色列表", Code: "user:role:list", ResourceType: "resource", Action: "list", SubjectType: "role", ResourceObjectClassID: objectClassMap["role"], Allowed: true, Priority: 70},
		{Name: "角色创建", Code: "user:role:create", ResourceType: "resource", Action: "create", SubjectType: "role", ResourceObjectClassID: objectClassMap["role"], Allowed: true, Priority: 70},
		{Name: "角色更新", Code: "user:role:update", ResourceType: "resource", Action: "update", SubjectType: "role", ResourceObjectClassID: objectClassMap["role"], Allowed: true, Priority: 70},
		{Name: "角色删除", Code: "user:role:delete", ResourceType: "resource", Action: "delete", SubjectType: "role", ResourceObjectClassID: objectClassMap["role"], Allowed: true, Priority: 70},
		{Name: "系统配置", Code: "system:config", ResourceType: "resource", Action: "all", SubjectType: "role", ResourceObjectClassID: objectClassMap["system_config"], Allowed: true, Priority: 80},
		{Name: "菜单管理", Code: "system:menu", ResourceType: "resource", Action: "all", SubjectType: "role", ResourceObjectClassID: objectClassMap["menu"], Allowed: true, Priority: 81},
		{Name: "操作日志", Code: "system:log", ResourceType: "resource", Action: "view", SubjectType: "role", ResourceObjectClassID: objectClassMap["operation_log"], Allowed: true, Priority: 82},
		{Name: "数据字典", Code: "system:dict", ResourceType: "resource", Action: "all", SubjectType: "role", ResourceObjectClassID: objectClassMap["data_dict"], Allowed: true, Priority: 83},
	}

	count := 0
	for i := range permissions {
		var perm ngac.Permission
		result := s.db.Where("code = ?", permissions[i].Code).First(&perm)
		if result.Error != nil {
			s.db.Create(&permissions[i])
			count++
		}
	}

	log.Printf("[SEEDER] ✓ 初始化 %d 个权限", count)
	return nil
}

// seedRolePermissions 初始化超级管理员角色与所有权限的关联
func (s *Seeder) seedRolePermissions() error {
	if !s.db.Migrator().HasTable("ngac_role_permissions") {
		log.Println("[SEEDER] ⊘ ngac_role_permissions 表不存在，跳过")
		return nil
	}

	var superAdminRole ngac.Role
	result := s.db.Where("code = ?", "super_admin").First(&superAdminRole)
	if result.Error != nil {
		log.Println("[SEEDER] ⊘ 找不到超级管理员角色，跳过")
		return nil
	}

	var allPerms []ngac.Permission
	s.db.Find(&allPerms)

	count := 0
	for _, perm := range allPerms {
		var existsCount int64
		s.db.Model(&ngac.RolePermission{}).
			Where("role_id = ? AND permission_id = ?", superAdminRole.ID, perm.ID).
			Count(&existsCount)

		if existsCount == 0 {
			stmt := s.db.Raw(
				"INSERT IGNORE INTO `ngac_role_permissions` (`role_id`, `permission_id`) VALUES (?, ?)",
				superAdminRole.ID, perm.ID,
			)
			if stmt.Error != nil {
				log.Printf("[SEEDER] ⚠ 关联权限失败: permission_id=%d, error=%v", perm.ID, stmt.Error)
				continue
			}
			count++
		}
	}

	log.Printf("[SEEDER] ✓ 超级管理员关联 %d 个权限", count)
	return nil
}

// seedUserRoles 初始化用户-角色关联
func (s *Seeder) seedUserRoles() error {
	if !s.db.Migrator().HasTable("user_roles") {
		log.Println("[SEEDER] ⊘ user_roles 表不存在，跳过")
		return nil
	}

	var admin model.User
	result := s.db.Where("username = ?", "admin").First(&admin)
	if result.Error != nil {
		log.Println("[SEEDER] ⊘ 找不到admin用户，跳过")
		return nil
	}

	var superAdminRole model.Role
	s.db.Where("code = ?", "super_admin").First(&superAdminRole)
	if superAdminRole.ID == 0 {
		log.Println("[SEEDER] ⊘ 找不到超级管理员角色，跳过")
		return nil
	}

	var exists model.UserRole
	result = s.db.Where("user_id = ? AND role_id = ?", admin.ID, superAdminRole.ID).First(&exists)
	if result.Error != nil {
		s.db.Create(&model.UserRole{
			UserID: admin.ID,
			RoleID: superAdminRole.ID,
		})
	}

	log.Printf("[SEEDER] ✓ 用户 admin 关联超级管理员角色")
	return nil
}

// menuDef 菜单定义结构
type menuDef struct {
	ID         uint
	Name       string
	Path       string
	Icon       string
	Component  string
	Sort       int
	MenuType   int
	Status     int
	Visible    int
	Permission string
	ParentID   *uint
}

// ptrUint 辅助函数：返回 *uint
func ptrUint(u uint) *uint {
	return &u
}

// menuPathToParentMap 定义子菜单路径到父菜单路径的映射关系
var menuPathToParentMap = map[string]string{
	"/articles/list":          "/content",
	"/articles/create":        "/content",
	"/articles/categories":    "/content",
	"/videos/list":            "/content",
	"/videos/create":          "/content",
	"/videos/categories":      "/content",
	"/comments/list":          "/content",
	"/media/files":            "/media",
	"/media/gallery":          "/media",
	"/users/list":             "/users",
	"/users/profile":          "/users",
	"/users/role":             "/users",
	"/users/permission":       "/users",
	"/audit/workbench":        "/audit",
	"/audit/articles":         "/audit",
	"/audit/videos":           "/audit",
	"/system/menus":           "/system",
	"/system/roles":           "/system",
	"/system/logs":            "/system",
	"/system/backup":          "/system",
	"/system/seo":             "/system",
	"/system/config":          "/system",
	"/ngac/permissions":       "/ngac",
	"/tenants/list":           "/tenants",
	"/messages/notifications": "/message",
	"/messages/system":        "/message",
	"/message-categories":     "/message",
	"/cron/tasks":             "/cron",
	"/cron/logs":              "/cron",
}

// parentMenuPath 获取子菜单对应的父菜单路径
func parentMenuPath(childPath string) string {
	if parentPath, ok := menuPathToParentMap[childPath]; ok {
		return parentPath
	}
	return ""
}

// seedMenus 初始化菜单数据
func (s *Seeder) seedMenus() error {
	if !s.db.Migrator().HasTable("menus") {
		log.Println("[SEEDER] ⊘ menus 表不存在，跳过")
		return nil
	}

	log.Println("[SEEDER] 清理可能的脏数据...")
	s.db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	result := s.db.Exec("DELETE FROM menus WHERE 1=1")
	if result.Error != nil {
		s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		log.Printf("[SEEDER] ⚠ 清理脏数据失败: %v", result.Error)
	} else {
		log.Printf("[SEEDER] ✓ 清理了 %d 条无效记录", result.RowsAffected)
		s.db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	}

	parentMenus := []menuDef{
		{ID: 1, Name: "仪表盘", Path: "/dashboard", Icon: "el-icon-data-line", Component: "views/dashboard/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1, Permission: "dashboard:view"},
		{ID: 2, Name: "内容管理", Path: "/content", Icon: "el-icon-document", Sort: 2, MenuType: 1, Status: 1, Visible: 1},
		{ID: 3, Name: "媒体中心", Path: "/media", Icon: "el-icon-picture-outline", Sort: 3, MenuType: 1, Status: 1, Visible: 1},
		{ID: 4, Name: "用户管理", Path: "/users", Icon: "el-icon-user", Sort: 4, MenuType: 1, Status: 1, Visible: 1, Permission: "user:list"},
		{ID: 5, Name: "消息管理", Path: "/message", Icon: "el-icon-bell", Sort: 5, MenuType: 1, Status: 1, Visible: 1},
		{ID: 6, Name: "审核管理", Path: "/audit", Icon: "el-icon-check", Sort: 6, MenuType: 1, Status: 1, Visible: 1},
		{ID: 7, Name: "系统设置", Path: "/system", Icon: "el-icon-setting", Sort: 7, MenuType: 1, Status: 1, Visible: 1, Permission: "system:config"},
		{ID: 8, Name: "NGAC权限", Path: "/ngac", Icon: "el-icon-lock", Sort: 8, MenuType: 1, Status: 1, Visible: 1},
		{ID: 9, Name: "租户管理", Path: "/tenants", Icon: "el-icon-office-building", Sort: 9, MenuType: 1, Status: 1, Visible: 1},
		{ID: 10, Name: "定时任务", Path: "/cron", Icon: "el-icon-time", Sort: 10, MenuType: 1, Status: 1, Visible: 1, Permission: "system:cron"},
	}

	parentMenusWithID := make([]struct {
		def menuDef
		id  uint
	}, 0, len(parentMenus))
	for _, m := range parentMenus {
		menu := model.Menu{
			ParentID:   nil,
			Name:       m.Name,
			Path:       m.Path,
			Component:  m.Component,
			Icon:       m.Icon,
			SortOrder:  m.Sort,
			MenuType:   m.MenuType,
			Permission: m.Permission,
			Visible:    m.Visible,
			Status:     m.Status,
		}
		if err := s.db.Create(&menu).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 创建顶级菜单失败: %s (%s), error: %v", m.Name, m.Path, err)
			continue
		}
		log.Printf("[SEEDER] ✓ 创建顶级菜单: %s (%s) ID=%d", m.Name, m.Path, menu.ID)
		parentMenusWithID = append(parentMenusWithID, struct {
			def menuDef
			id  uint
		}{def: m, id: menu.ID})
	}

	if len(parentMenusWithID) == 0 {
		return nil
	}

	parentIDMap := make(map[string]uint, len(parentMenus))
	for _, item := range parentMenusWithID {
		parentIDMap[item.def.Path] = item.id
	}

	childMenus := []menuDef{
		{Name: "文章列表", Path: "/articles/list", Icon: "el-icon-document", Component: "views/content/article/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1, Permission: "content:article:list"},
		{Name: "发布文章", Path: "/articles/create", Icon: "el-icon-plus", Component: "views/content/article/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
		{Name: "分类管理", Path: "/articles/categories", Icon: "el-icon-menu", Component: "views/content/category/index", Sort: 3, MenuType: 2, Status: 1, Visible: 1, Permission: "content:category:list"},
		{Name: "视频列表", Path: "/videos/list", Icon: "el-icon-video-play", Component: "views/content/video/index", Sort: 4, MenuType: 2, Status: 1, Visible: 1, Permission: "content:video:list"},
		{Name: "上传视频", Path: "/videos/create", Icon: "el-icon-camera", Component: "views/content/video/index", Sort: 5, MenuType: 2, Status: 1, Visible: 1},
		{Name: "视频分类", Path: "/videos/categories", Icon: "el-icon-menu", Component: "views/content/video-category/index", Sort: 6, MenuType: 2, Status: 1, Visible: 1},
		{Name: "评论管理", Path: "/comments/list", Icon: "el-icon-chat-dot-round", Component: "views/content/comment/index", Sort: 7, MenuType: 2, Status: 1, Visible: 1, Permission: "content:comment:list"},
		{Name: "文件管理", Path: "/media/files", Icon: "el-icon-folder", Component: "views/media/file/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1, Permission: "media:file"},
		{Name: "图库管理", Path: "/media/gallery", Icon: "el-icon-picture", Component: "views/media/gallery/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1, Permission: "media:gallery"},
		{Name: "用户列表", Path: "/users/list", Icon: "el-icon-user", Component: "views/user/list/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1, Permission: "user:list"},
		{Name: "个人信息", Path: "/users/profile", Icon: "el-icon-user-filled", Component: "views/user/list/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
		{Name: "角色管理", Path: "/users/role", Icon: "el-icon-s-custom", Component: "views/user/role/index", Sort: 3, MenuType: 2, Status: 1, Visible: 1, Permission: "user:role:list"},
		{Name: "权限管理", Path: "/users/permission", Icon: "el-icon-lock", Component: "views/user/permission/index", Sort: 4, MenuType: 2, Status: 1, Visible: 1},
		{Name: "我的通知", Path: "/messages/notifications", Icon: "el-icon-bell", Component: "views/message/NotificationPage", Sort: 1, MenuType: 2, Status: 1, Visible: 1},
		{Name: "系统消息管理", Path: "/messages/system", Icon: "el-icon-message", Component: "views/message/SystemMessagePage", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
		{Name: "消息类别管理", Path: "/message-categories", Icon: "el-icon-collection", Component: "views/message/CategoryPage", Sort: 3, MenuType: 2, Status: 1, Visible: 1},
		{Name: "审核工作台", Path: "/audit/workbench", Icon: "el-icon-check", Component: "views/audit/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1},
		{Name: "文章审核", Path: "/audit/articles", Icon: "el-icon-document", Component: "views/audit/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
		{Name: "视频审核", Path: "/audit/videos", Icon: "el-icon-video-play", Component: "views/audit/index", Sort: 3, MenuType: 2, Status: 1, Visible: 1},
		{Name: "菜单管理", Path: "/system/menus", Icon: "el-icon-menu", Component: "views/system/menu/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1},
		{Name: "角色权限", Path: "/system/roles", Icon: "el-icon-s-custom", Component: "views/system/role/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
		{Name: "操作日志", Path: "/system/logs", Icon: "el-icon-document", Component: "views/system/log/index", Sort: 3, MenuType: 2, Status: 1, Visible: 1},
		{Name: "数据备份", Path: "/system/backup", Icon: "el-icon-download", Component: "views/system/backup/index", Sort: 4, MenuType: 2, Status: 1, Visible: 1},
		{Name: "SEO管理", Path: "/system/seo", Icon: "el-icon-data-line", Component: "views/system/seo/index", Sort: 5, MenuType: 2, Status: 1, Visible: 1},
		{Name: "系统配置", Path: "/system/config", Icon: "el-icon-setting", Component: "views/system/config/index", Sort: 6, MenuType: 2, Status: 1, Visible: 1},
		{Name: "权限列表", Path: "/ngac/permissions", Icon: "el-icon-lock", Component: "views/system/permission/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1},
		{Name: "租户列表", Path: "/tenants/list", Icon: "el-icon-office-building", Component: "views/system/tenant/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1},
		{Name: "任务列表", Path: "/cron/tasks", Icon: "el-icon-time", Component: "views/cron/task/index", Sort: 1, MenuType: 2, Status: 1, Visible: 1, Permission: "system:cron"},
		{Name: "执行日志", Path: "/cron/logs", Icon: "el-icon-document", Component: "views/cron/log/index", Sort: 2, MenuType: 2, Status: 1, Visible: 1},
	}

	validChildCount := 0
	for i := range childMenus {
		parentPath := parentMenuPath(childMenus[i].Path)
		if parentPath == "" {
			continue
		}
		parentID, ok := parentIDMap[parentPath]
		if !ok {
			continue
		}
		childMenus[i].ParentID = ptrUint(parentID)
		validChildCount++
	}

	for _, m := range childMenus {
		if m.ParentID == nil {
			continue
		}
		menu := model.Menu{
			Name:       m.Name,
			Path:       m.Path,
			Icon:       m.Icon,
			Component:  m.Component,
			SortOrder:  m.Sort,
			MenuType:   m.MenuType,
			Status:     m.Status,
			Visible:    m.Visible,
			Permission: m.Permission,
			ParentID:   m.ParentID,
		}
		if err := s.db.Create(&menu).Error; err != nil {
			log.Printf("[SEEDER] ⚠ 创建子菜单失败: %s (%s), error: %v", m.Name, m.Path, err)
			continue
		}
		log.Printf("[SEEDER] ✓ 创建子菜单: %s -> parent_id=%d", m.Name, *m.ParentID)
	}

	log.Printf("[SEEDER] ✓ 菜单初始化完成，共 %d 个菜单", len(parentMenusWithID)+validChildCount)
	return nil
}

// seedArticleCategories 初始化文章分类
func (s *Seeder) seedArticleCategories() error {
	if err := s.db.AutoMigrate(&model.ArticleCategory{}); err != nil {
		log.Printf("[SEEDER] ⚠ 迁移 article_categories 表失败: %v", err)
		return nil
	}

	var count int64
	if err := s.db.Model(&model.ArticleCategory{}).Count(&count).Error; err != nil {
		return nil
	}
	if count > 0 {
		log.Println("[SEEDER] ⊘ 文章分类数据已存在，跳过")
		return nil
	}

	categories := []struct {
		Name string
		Icon string
		Sort int
	}{
		{"技术分享", "el-icon-cpu", 1},
		{"产品资讯", "el-icon-news", 2},
		{"业界动态", "el-icon-trophy", 3},
		{"开发教程", "el-icon-school", 4},
		{"开源项目", "el-icon-box", 5},
	}

	for _, c := range categories {
		s.db.Create(&model.ArticleCategory{Name: c.Name, Icon: c.Icon, SortOrder: c.Sort, Status: 1})
	}
	log.Printf("[SEEDER] ✓ 初始化 %d 个文章分类", len(categories))
	return nil
}

// seedVideoCategories 初始化视频分类
func (s *Seeder) seedVideoCategories() error {
	if err := s.db.AutoMigrate(&model.VideoCategory{}); err != nil {
		log.Printf("[SEEDER] ⚠ 迁移 video_categories 表失败: %v", err)
		return nil
	}

	var count int64
	if err := s.db.Model(&model.VideoCategory{}).Count(&count).Error; err != nil {
		return nil
	}
	if count > 0 {
		log.Println("[SEEDER] ⊘ 视频分类数据已存在，跳过")
		return nil
	}

	categories := []struct {
		Name string
		Icon string
		Sort int
	}{
		{"视频教程", "el-icon-video-play", 1},
		{"直播回放", "el-icon-camera", 2},
		{"课程专区", "el-icon-school", 3},
		{"技术分享", "el-icon-cpu", 4},
	}

	for _, c := range categories {
		s.db.Create(&model.VideoCategory{Name: c.Name, Icon: c.Icon, Sort: c.Sort, Status: 1})
	}
	log.Printf("[SEEDER] ✓ 初始化 %d 个视频分类", len(categories))
	return nil
}

// seedSystemConfigs 初始化系统配置
func (s *Seeder) seedSystemConfigs() error {
	if err := s.db.AutoMigrate(&model.SystemConfig{}); err != nil {
		log.Printf("[SEEDER] ⚠ 迁移 system_configs 表失败: %v", err)
		return nil
	}

	var count int64
	if err := s.db.Model(&model.SystemConfig{}).Count(&count).Error; err != nil {
		return nil
	}
	if count > 0 {
		log.Println("[SEEDER] ⊘ 系统配置数据已存在，跳过")
		return nil
	}

	configs := []struct{ Key, Value, Type, Desc string }{
		{"site_name", "GoCMS V3", "string", "网站名称"},
		{"site_description", "现代化内容管理系统", "string", "网站描述"},
		{"upload_max_size", "104857600", "number", "最大上传大小"},
		{"password_min_length", "8", "number", "密码最小长度"},
		{"login_max_attempts", "5", "number", "最大登录尝试次数"},
		{"session_lifetime", "86400", "number", "会话有效期"},
		{"cdn_enabled", "0", "boolean", "是否启用CDN"},
		{"cache_enabled", "1", "boolean", "是否启用缓存"},
		{"i18n_default_locale", "zh-CN", "string", "默认语言"},
	}

	for _, cfg := range configs {
		s.db.Create(&model.SystemConfig{ConfigKey: cfg.Key, ConfigValue: cfg.Value, ConfigType: cfg.Type, Description: cfg.Desc, IsSystem: true})
	}
	log.Printf("[SEEDER] ✓ 初始化 %d 个系统配置", len(configs))
	return nil
}

// seedDataDicts 初始化数据字典
func (s *Seeder) seedDataDicts() error {
	if err := s.db.AutoMigrate(&model.DataDict{}); err != nil {
		log.Printf("[SEEDER] ⚠ 迁移 data_dicts 表失败: %v", err)
		return nil
	}

	var count int64
	if err := s.db.Model(&model.DataDict{}).Count(&count).Error; err != nil {
		return nil
	}
	if count > 0 {
		log.Println("[SEEDER] ⊘ 数据字典数据已存在，跳过")
		return nil
	}

	dicts := []struct {
		Type, Label, Value string
		Sort               int
	}{
		{"article_status", "草稿", "0", 1},
		{"article_status", "已发布", "1", 2},
		{"video_status", "草稿", "0", 1},
		{"video_status", "已发布", "1", 2},
		{"comment_status", "待审核", "0", 1},
		{"comment_status", "已通过", "1", 2},
		{"user_status", "启用", "1", 1},
		{"user_status", "禁用", "0", 2},
	}

	for _, d := range dicts {
		s.db.Create(&model.DataDict{DictType: d.Type, DictLabel: d.Label, DictValue: d.Value, SortOrder: d.Sort, Status: 1})
	}
	log.Printf("[SEEDER] ✓ 初始化 %d 个数据字典", len(dicts))
	return nil
}

// autoMigrateSeoTables 自动迁移 SEO 相关表
func (s *Seeder) autoMigrateSeoTables() error {
	if err := s.db.AutoMigrate(&model.SeoSetting{}); err != nil {
		return fmt.Errorf("迁移 seo_settings 表失败: %w", err)
	}
	if err := s.db.AutoMigrate(&model.SitemapConfig{}); err != nil {
		return fmt.Errorf("迁移 sitemap_configs 表失败: %w", err)
	}
	// 自动迁移定时任务相关表（CronJob 实体在 app/domain/entity/comment.go 中定义）
	if err := s.db.AutoMigrate(&entity.CronJob{}); err != nil {
		return fmt.Errorf("迁移 cron_jobs 表失败: %w", err)
	}
	return nil
}

// seedSitemapConfigs 初始化站点地图配置
func (s *Seeder) seedSitemapConfigs() error {
	if err := s.db.AutoMigrate(&model.SitemapConfig{}); err != nil {
		log.Printf("[SEEDER] ⚠ 迁移 sitemap_configs 表失败: %v", err)
		return nil
	}

	var count int64
	if err := s.db.Model(&model.SitemapConfig{}).Count(&count).Error; err != nil {
		return nil
	}
	if count > 0 {
		log.Println("[SEEDER] ⊘ 站点地图配置数据已存在，跳过")
		return nil
	}

	configs := []model.SitemapConfig{
		{Name: "文章站点地图", Type: "article", Priority: 0.8, ChangeFreq: "daily", IsEnabled: 1, MaxCount: 1000},
		{Name: "视频站点地图", Type: "video", Priority: 0.7, ChangeFreq: "weekly", IsEnabled: 1, MaxCount: 500},
		{Name: "文章分类站点地图", Type: "category", Priority: 0.6, ChangeFreq: "weekly", IsEnabled: 1, MaxCount: 100},
	}

	for _, cfg := range configs {
		var existing model.SitemapConfig
		if err := s.db.Where("type = ?", cfg.Type).First(&existing).Error; err != nil {
			s.db.Create(&cfg)
			log.Printf("[SEEDER] ✓ 创建站点地图配置: %s (%s)", cfg.Name, cfg.Type)
		}
	}
	log.Printf("[SEEDER] ✓ 初始化 %d 个站点地图配置", len(configs))
	return nil
}

// seedObjectClasses 初始化 NGAC 对象类
func (s *Seeder) seedObjectClasses() error {
	if !s.db.Migrator().HasTable("ngac_object_classes") {
		log.Println("[SEEDER] ⊘ ngac_object_classes 表不存在，跳过")
		return nil
	}

	var count int64
	s.db.Raw("SELECT COUNT(*) FROM ngac_object_classes").Scan(&count)
	if count > 0 {
		log.Println("[SEEDER] ⊘ NGAC对象类数据已存在，跳过")
		return nil
	}

	objectClasses := []ngac.ObjectClass{
		{ID: 1, Name: "文章", Code: "article", Description: "文章内容资源", Level: 1, Path: "/article", Sort: 1, Status: 1},
		{ID: 2, Name: "文章分类", Code: "article_category", Description: "文章分类资源", Level: 1, Path: "/article_category", Sort: 2, Status: 1},
		{ID: 3, Name: "视频", Code: "video", Description: "视频内容资源", Level: 1, Path: "/video", Sort: 3, Status: 1},
		{ID: 4, Name: "评论", Code: "comment", Description: "评论内容资源", Level: 1, Path: "/comment", Sort: 4, Status: 1},
		{ID: 5, Name: "文件", Code: "file", Description: "文件资源", Level: 1, Path: "/file", Sort: 5, Status: 1},
		{ID: 6, Name: "图库", Code: "gallery", Description: "图库资源", Level: 1, Path: "/gallery", Sort: 6, Status: 1},
		{ID: 7, Name: "用户", Code: "user", Description: "用户资源", Level: 1, Path: "/user", Sort: 7, Status: 1},
		{ID: 8, Name: "角色", Code: "role", Description: "角色资源", Level: 1, Path: "/role", Sort: 8, Status: 1},
		{ID: 9, Name: "系统配置", Code: "system_config", Description: "系统配置资源", Level: 1, Path: "/system_config", Sort: 9, Status: 1},
		{ID: 10, Name: "菜单", Code: "menu", Description: "菜单资源", Level: 1, Path: "/menu", Sort: 10, Status: 1},
	}

	s.db.Create(&objectClasses)
	log.Printf("[SEEDER] ✓ 初始化 %d 个NGAC对象类", len(objectClasses))
	return nil
}

// seedRiskKeywords 初始化风险关键词库
func (s *Seeder) seedRiskKeywords() error {
	if !s.db.Migrator().HasTable("risk_keywords") {
		log.Println("[SEEDER] ⊘ risk_keywords 表不存在，跳过")
		return nil
	}

	if err := s.db.AutoMigrate(&model.RiskKeyword{}); err != nil {
		return fmt.Errorf("迁移风险关键词表失败: %w", err)
	}

	allInserted := 0
	keywords := []struct {
		Keyword, Category, Description string
		Level                          int
	}{
		{"测试敏感词1", "politics", "政治敏感词示例", 4},
		{"测试敏感词2", "violence", "暴力敏感词示例", 3},
		{"测试敏感词3", "sex", "色情敏感词示例", 3},
	}

	for _, kw := range keywords {
		var existing model.RiskKeyword
		if err := s.db.Where("keyword = ?", kw.Keyword).First(&existing).Error; err != nil {
			s.db.Create(&model.RiskKeyword{Keyword: kw.Keyword, Level: int8(kw.Level), Category: kw.Category, Description: kw.Description, Enabled: true, MatchType: 1})
			allInserted++
		}
	}

	log.Printf("[SEEDER] ✓ 风险关键词库初始化完成，新增 %d 个关键词", allInserted)
	return nil
}

// seedAlbums 初始化媒体相册种子数据
func (s *Seeder) seedAlbums() error {
	if !s.db.Migrator().HasTable("albums") {
		log.Println("[SEEDER] ⊘ albums 表不存在，跳过")
		return nil
	}
	if !s.db.Migrator().HasTable("media_assets") {
		log.Println("[SEEDER] ⊘ media_assets 表不存在，跳过")
		return nil
	}

	var albumCount int64
	s.db.Raw("SELECT COUNT(*) FROM albums").Scan(&albumCount)
	if albumCount > 0 {
		log.Println("[SEEDER] ⊘ 相册数据已存在，跳过")
		return nil
	}

	var adminID uint64
	result := s.db.Raw("SELECT id FROM users WHERE username = 'admin' LIMIT 1").Scan(&adminID)
	if result.Error != nil {
		return nil
	}

	defaultAlbums := []struct {
		Name, Description string
		IsPublic          int8
	}{
		{"默认相册", "系统默认创建的相册", 0},
		{"我的图片", "我上传的所有图片", 0},
		{"视频收藏", "我收藏的视频资源", 0},
	}

	for _, album := range defaultAlbums {
		var count int64
		s.db.Raw("SELECT COUNT(*) FROM albums WHERE name = ?", album.Name).Scan(&count)
		if count == 0 {
			s.db.Create(&model.Album{Name: album.Name, Description: album.Description, UserID: adminID, IsPublic: album.IsPublic, Status: 1})
			log.Printf("[SEEDER] ✓ 创建相册: %s", album.Name)
		}
	}

	var assetCount int64
	s.db.Raw("SELECT COUNT(*) FROM media_assets").Scan(&assetCount)
	if assetCount > 0 {
		return nil
	}

	log.Println("[SEEDER] 初始化媒体资产种子数据...")

	testAssets := []struct {
		Name, OriginalName, FilePath, FileType, MimeType string
		FileSize                                         int64
		Width, Height, Duration                          int
		Thumbnail                                        string
	}{
		{"test_photo_1.jpg", "test_photo_1.jpg", "/storage/uploads/test_photo_1.jpg", "image", "image/jpeg", 256000, 1920, 1080, 0, "/storage/uploads/thumbnails/test_photo_1_thumb.jpg"},
		{"test_photo_2.png", "test_photo_2.png", "/storage/uploads/test_photo_2.png", "image", "image/png", 512000, 1920, 1080, 0, "/storage/uploads/thumbnails/test_photo_2_thumb.png"},
		{"demo_video.mp4", "demo_video.mp4", "/storage/videos/demo_video.mp4", "video", "video/mp4", 52428800, 1920, 1080, 180, ""},
		{"document.pdf", "documentation.pdf", "/storage/uploads/document.pdf", "file", "application/pdf", 1048576, 0, 0, 0, ""},
	}

	for _, asset := range testAssets {
		tagsJSON := `["测试","示例"]`
		if asset.FileType == "file" {
			tagsJSON = `[]`
		}
		now := time.Now().Add(-time.Duration(rand.Intn(30)) * time.Hour * 24)
		mediaAsset := model.MediaAsset{
			UserID: adminID, UploadID: fmt.Sprintf("upload_%d", rand.Intn(10000)),
			Name: asset.Name, OriginalName: asset.OriginalName, FilePath: asset.FilePath,
			FileType: asset.FileType, MimeType: asset.MimeType, FileSize: asset.FileSize,
			Width: asset.Width, Height: asset.Height, Duration: asset.Duration,
			Thumbnail: asset.Thumbnail, Description: "测试媒体资产", Tags: tagsJSON, Status: 1,
		}
		mediaAsset.CreatedAt = now
		mediaAsset.UpdatedAt = now
		s.db.Create(&mediaAsset)
		log.Printf("[SEEDER] ✓ 创建媒体资产: %s (%s)", asset.Name, asset.FileType)
	}

	log.Printf("[SEEDER] ✓ 初始化 %d 个媒体资产", len(testAssets))
	return nil
}

// seedDefaultMessageCategories 初始化默认消息类别
func (s *Seeder) seedDefaultMessageCategories() error {
	if !s.db.Migrator().HasTable("message_categories") {
		log.Println("[SEEDER] ⊘ message_categories 表不存在，跳过")
		return nil
	}

	var count int64
	s.db.Model(&model.MessageCategory{}).Count(&count)
	if count > 0 {
		log.Println("[SEEDER] ⊘ 消息类别数据已存在，跳过")
		return nil
	}

	log.Println("[SEEDER] 初始化默认消息类别...")

	defaultCategories := []struct {
		Name, Code, Description string
		Sort                    int
	}{
		{"系统通知", "system", "系统级别的通知消息", 1},
		{"发布通知", "publish", "内容发布相关的通知", 2},
		{"评论回复", "comment_reply", "用户评论和回复通知", 3},
		{"审核通过", "audit_passed", "内容审核通过通知", 4},
		{"审核驳回", "audited_failed", "内容审核未通过通知", 5},
	}

	for _, cat := range defaultCategories {
		s.db.Create(&model.MessageCategory{Name: cat.Name, Code: cat.Code, Description: cat.Description, IconType: strings.ToLower(cat.Code), Sort: cat.Sort, Status: 1})
		log.Printf("[SEEDER] ✓ 创建消息类别: %s (%s)", cat.Name, cat.Code)
	}

	log.Printf("[SEEDER] ✓ 消息类别初始化完成，共 %d 条", len(defaultCategories))
	return nil
}
