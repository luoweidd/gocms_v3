package service

import (
	"errors"
	"fmt"
	"log"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"

	"gorm.io/gorm"
)

// VideoService 视频内容服务
type VideoService struct{}

// NewVideoService 获取视频服务实例
func NewVideoService() *VideoService {
	return &VideoService{}
}

// CreateVideoRequest 创建视频请求
type CreateVideoRequest struct {
	Title    string `json:"title" binding:"required"`
	Summary  string `json:"summary"`
	Content  string `json:"content"`
	URL      string `json:"url" binding:"required"`
	Cover    string `json:"cover"`
	Format   string `json:"format"`
	Duration int    `json:"duration"`
	FileSize int64  `json:"file_size"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Bitrate  int    `json:"bitrate"`

	// 影视属性
	Year           int     `json:"year"`
	ReleaseDate    string  `json:"release_date"`
	Director       string  `json:"director"`
	Actors         string  `json:"actors"`
	Genre          string  `json:"genre"`
	Region         string  `json:"region"`
	Language       string  `json:"language"`
	Rating         float64 `json:"rating"`
	TotalEpisodes  int     `json:"total_episodes"`
	CurrentEpisode int     `json:"current_episode"`

	// 管理
	CategoryID  uint   `json:"category_id"`
	Status      int    `json:"status"`
	IsTop       bool   `json:"is_top"`
	PublishedAt string `json:"published_at"`
	TagIDs      []uint `json:"tag_ids"`
	AuthorID    uint   `json:"-"`
}

// UpdateVideoRequest 更新视频请求
type UpdateVideoRequest struct {
	Title    *string `json:"title"`
	Summary  *string `json:"summary"`
	Content  *string `json:"content"`
	URL      *string `json:"url"`
	Cover    *string `json:"cover"`
	Format   *string `json:"format"`
	Duration *int    `json:"duration"`
	FileSize *int64  `json:"file_size"`
	Width    *int    `json:"width"`
	Height   *int    `json:"height"`
	Bitrate  *int    `json:"bitrate"`

	// 影视属性
	Year           *int     `json:"year"`
	ReleaseDate    *string  `json:"release_date"`
	Director       *string  `json:"director"`
	Actors         *string  `json:"actors"`
	Genre          *string  `json:"genre"`
	Region         *string  `json:"region"`
	Language       *string  `json:"language"`
	Rating         *float64 `json:"rating"`
	TotalEpisodes  *int     `json:"total_episodes"`
	CurrentEpisode *int     `json:"current_episode"`

	// 管理
	CategoryID  *uint   `json:"category_id"`
	Status      *int    `json:"status"`
	IsTop       *bool   `json:"is_top"`
	PublishedAt *string `json:"published_at"`
	TagIDs      []uint  `json:"tag_ids"`
}

// VideoQueryParams 视频查询参数
type VideoQueryParams struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"page_size"`
	CategoryID uint   `form:"category_id"`
	Status     *int   `form:"status"` // 使用指针类型，nil 表示不过滤
	Keyword    string `form:"keyword"`
	AuthorID   uint   `form:"author_id"`
	Year       int    `form:"year"`
	Genre      string `form:"genre"`
	Region     string `form:"region"`
	Language   string `form:"language"`
	Director   string `form:"director"`
	SortBy     string `form:"sort_by"`
	SortOrder  string `form:"sort_order"`
}

// CreateVideo 创建视频
func (s *VideoService) CreateVideo(req CreateVideoRequest) (*model.Video, error) {
	var publishedAt, releaseDate *time.Time
	if req.PublishedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", req.PublishedAt); err == nil {
			publishedAt = &t
		}
	}
	if req.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", req.ReleaseDate); err == nil {
			releaseDate = &t
		}
	}

	video := model.Video{
		Title:          req.Title,
		Summary:        req.Summary,
		Content:        req.Content,
		URL:            req.URL,
		Cover:          req.Cover,
		Format:         req.Format,
		Duration:       req.Duration,
		FileSize:       req.FileSize,
		Width:          req.Width,
		Height:         req.Height,
		Bitrate:        req.Bitrate,
		Year:           req.Year,
		ReleaseDate:    releaseDate,
		Director:       req.Director,
		Actors:         req.Actors,
		Genre:          req.Genre,
		Region:         req.Region,
		Language:       req.Language,
		Rating:         req.Rating,
		TotalEpisodes:  req.TotalEpisodes,
		CurrentEpisode: req.CurrentEpisode,
		CategoryID:     req.CategoryID,
		AuthorID:       req.AuthorID,
		Status:         0, // 默认草稿状态，需要审核
		IsTop:          req.IsTop,
		PublishedAt:    publishedAt,
	}

	if err := db.GetDB().Create(&video).Error; err != nil {
		return nil, fmt.Errorf("创建视频失败: %w", err)
	}

	// 关联标签
	if len(req.TagIDs) > 0 {
		s.assignTags(&video, req.TagIDs)
	}

	// 更新分类视频数量
	if req.CategoryID > 0 {
		svc := NewVideoCategoryService()
		svc.UpdateVideoCount(req.CategoryID)
	}

	// 自动创建审核记录（待审核状态）
	auditService := NewContentAuditService()
	_, err := auditService.SubmitAudit(struct {
		ContentType   string
		ContentID     uint
		ContentTitle  string
		SubmitterID   uint
		SubmitterName string
		NextStatus    int8
	}{
		ContentType:   "video",
		ContentID:     video.ID,
		ContentTitle:  video.Title,
		SubmitterID:   req.AuthorID,
		SubmitterName: fmt.Sprintf("用户%d", req.AuthorID),
		NextStatus:    0, // 待审核
	})
	if err != nil {
		log.Printf("[SERVICE-WARN] 创建视频审核记录失败 - 视频ID: %d, 错误: %v", video.ID, err)
	} else {
		log.Printf("[SERVICE] 视频审核记录已创建 - 视频ID: %d", video.ID)
	}

	// 重新加载（含标签和分类）
	db.GetDB().Preload("Tags").Preload("Category").First(&video, video.ID)
	return &video, nil
}

// UpdateVideo 更新视频
func (s *VideoService) UpdateVideo(id uint, req UpdateVideoRequest) error {
	var video model.Video
	if err := db.GetDB().First(&video, id).Error; err != nil {
		return errors.New("视频不存在")
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Cover != nil {
		updates["cover"] = *req.Cover
	}
	if req.Format != nil {
		updates["format"] = *req.Format
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.FileSize != nil {
		updates["file_size"] = *req.FileSize
	}
	if req.Width != nil {
		updates["width"] = *req.Width
	}
	if req.Height != nil {
		updates["height"] = *req.Height
	}
	if req.Bitrate != nil {
		updates["bitrate"] = *req.Bitrate
	}

	// 影视属性
	if req.Year != nil {
		updates["year"] = *req.Year
	}
	if req.ReleaseDate != nil && *req.ReleaseDate != "" {
		if t, err := time.Parse("2006-01-02", *req.ReleaseDate); err == nil {
			updates["release_date"] = t
		}
	}
	if req.Director != nil {
		updates["director"] = *req.Director
	}
	if req.Actors != nil {
		updates["actors"] = *req.Actors
	}
	if req.Genre != nil {
		updates["genre"] = *req.Genre
	}
	if req.Region != nil {
		updates["region"] = *req.Region
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.TotalEpisodes != nil {
		updates["total_episodes"] = *req.TotalEpisodes
	}
	if req.CurrentEpisode != nil {
		updates["current_episode"] = *req.CurrentEpisode
	}

	// 管理字段
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.IsTop != nil {
		updates["is_top"] = *req.IsTop
	}
	if req.PublishedAt != nil && *req.PublishedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", *req.PublishedAt); err == nil {
			updates["published_at"] = t
		}
	}

	if err := db.GetDB().Model(&video).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新视频失败: %w", err)
	}

	// 更新标签关联
	if req.TagIDs != nil {
		s.assignTags(&video, req.TagIDs)
	}

	// 更新分类视频数量
	if req.CategoryID != nil {
		oldCat := video.CategoryID
		newCat := *req.CategoryID
		if oldCat != newCat {
			svc := NewVideoCategoryService()
			if oldCat > 0 {
				svc.UpdateVideoCount(oldCat)
			}
			if newCat > 0 {
				svc.UpdateVideoCount(newCat)
			}
		}
	}

	return nil
}

// DeleteVideo 删除视频（软删除）
func (s *VideoService) DeleteVideo(id uint) error {
	var video model.Video
	if err := db.GetDB().First(&video, id).Error; err != nil {
		return errors.New("视频不存在")
	}

	if err := db.GetDB().Delete(&model.Video{}, id).Error; err != nil {
		return fmt.Errorf("删除视频失败: %w", err)
	}

	// 更新分类视频数量
	if video.CategoryID > 0 {
		svc := NewVideoCategoryService()
		svc.UpdateVideoCount(video.CategoryID)
	}

	return nil
}

// GetVideoList 获取视频列表（分页+多维筛选）
func (s *VideoService) GetVideoList(params VideoQueryParams) (*response.PageData, error) {
	DB := db.GetDB()
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化，请检查连接配置")
	}

	// 检查 videos 表是否存在（使用动态获取数据库名）
	var tableExists int
	dbName := getCurrentDatabaseName(DB)
	if err := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dbName, "videos").Scan(&tableExists).Error; err == nil && tableExists == 0 {
		return nil, fmt.Errorf("videos 表不存在，请先运行数据库迁移")
	}
	// 使用 DB.Table 明确指定表名，避免 GORM 无法推断表的问题
	query := DB.Table("videos")

	// 筛选条件
	if params.CategoryID > 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	// Status: nil 表示不过滤，非 nil 表示按状态过滤
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.AuthorID > 0 {
		query = query.Where("author_id = ?", params.AuthorID)
	}
	if params.Keyword != "" {
		keyword := fmt.Sprintf("%%%s%%", params.Keyword)
		query = query.Where("title LIKE ? OR summary LIKE ? OR director LIKE ? OR actors LIKE ?",
			keyword, keyword, keyword, keyword)
	}
	if params.Year > 0 {
		query = query.Where("year = ?", params.Year)
	}
	if params.Genre != "" {
		query = query.Where("genre LIKE ?", fmt.Sprintf("%%%s%%", params.Genre))
	}
	if params.Region != "" {
		query = query.Where("region LIKE ?", fmt.Sprintf("%%%s%%", params.Region))
	}
	if params.Language != "" {
		query = query.Where("language LIKE ?", fmt.Sprintf("%%%s%%", params.Language))
	}
	if params.Director != "" {
		query = query.Where("director LIKE ?", fmt.Sprintf("%%%s%%", params.Director))
	}

	// 统计总数（使用 Scopes 避免 Session 导致表上下文丢失）
	var total int64
	if err := DB.Table("videos").Scopes(func(db *gorm.DB) *gorm.DB {
		if params.CategoryID > 0 {
			db = db.Where("category_id = ?", params.CategoryID)
		}
		if params.Status != nil {
			db = db.Where("status = ?", *params.Status)
		}
		if params.AuthorID > 0 {
			db = db.Where("author_id = ?", params.AuthorID)
		}
		if params.Keyword != "" {
			keyword := fmt.Sprintf("%%%s%%", params.Keyword)
			db = db.Where("title LIKE ? OR summary LIKE ? OR director LIKE ? OR actors LIKE ?",
				keyword, keyword, keyword, keyword)
		}
		if params.Year > 0 {
			db = db.Where("year = ?", params.Year)
		}
		if params.Genre != "" {
			db = db.Where("genre LIKE ?", fmt.Sprintf("%%%s%%", params.Genre))
		}
		if params.Region != "" {
			db = db.Where("region LIKE ?", fmt.Sprintf("%%%s%%", params.Region))
		}
		if params.Language != "" {
			db = db.Where("language LIKE ?", fmt.Sprintf("%%%s%%", params.Language))
		}
		if params.Director != "" {
			db = db.Where("director LIKE ?", fmt.Sprintf("%%%s%%", params.Director))
		}
		return db
	}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	sortField := params.SortBy
	allowedSortFields := map[string]bool{
		"created_at": true, "updated_at": true, "view_count": true,
		"duration": true, "file_size": true, "published_at": true,
		"rating": true, "year": true, "release_date": true,
	}
	if !allowedSortFields[sortField] {
		sortField = "created_at"
	}
	sortOrder := params.SortOrder
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	offset := (params.Page - 1) * params.PageSize

	// 先获取视频列表（不预加载关联数据，避免外键错误）
	var videos []model.Video
	if err := query.
		Offset(offset).
		Limit(params.PageSize).
		Order("is_top DESC, " + sortField + " " + sortOrder).
		Find(&videos).Error; err != nil {
		return nil, fmt.Errorf("查询视频列表失败: %w", err)
	}

	// 注意：不预加载标签和分类，避免many2many关联表不存在导致的错误
	// 如需预加载，请确保 migration 已执行创建 video_tags 表

	return &response.PageData{
		List:     videos,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetVideoByID 获取视频详情
func (s *VideoService) GetVideoByID(id uint) (*model.Video, error) {
	var video model.Video
	if err := db.GetDB().Preload("Tags").Preload("Category").First(&video, id).Error; err != nil {
		return nil, errors.New("视频不存在")
	}
	return &video, nil
}

// IncrementViewCount 增加播放次数
func (s *VideoService) IncrementViewCount(id uint) error {
	return db.GetDB().Model(&model.Video{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// SubmitRating 提交评分
func (s *VideoService) SubmitRating(videoID uint, rating float64) error {
	if rating < 0 || rating > 10 {
		return errors.New("评分必须在0-10之间")
	}

	var video model.Video
	if err := db.GetDB().First(&video, videoID).Error; err != nil {
		return errors.New("视频不存在")
	}

	// 计算新的平均分
	newRating := (video.Rating*float64(video.RatingCount) + rating) / float64(video.RatingCount+1)

	return db.GetDB().Model(&video).Updates(map[string]interface{}{
		"rating":       newRating,
		"rating_count": video.RatingCount + 1,
	}).Error
}

// GetPublishedVideos 获取已发布视频（公开接口，支持筛选）
func (s *VideoService) GetPublishedVideos(page, pageSize int, categoryID uint, genre, region, language string, year int) (*response.PageData, error) {
	DB := db.GetDB()

	// 使用 Scopes 避免 Session(NewDB: true) 导致表上下文丢失
	var total int64
	if err := DB.Table("videos").Where("status = 1").Scopes(func(db *gorm.DB) *gorm.DB {
		if categoryID > 0 {
			db = db.Where("category_id = ?", categoryID)
		}
		if genre != "" {
			db = db.Where("genre LIKE ?", fmt.Sprintf("%%%s%%", genre))
		}
		if region != "" {
			db = db.Where("region LIKE ?", fmt.Sprintf("%%%s%%", region))
		}
		if language != "" {
			db = db.Where("language LIKE ?", fmt.Sprintf("%%%s%%", language))
		}
		if year > 0 {
			db = db.Where("year = ?", year)
		}
		return db
	}).Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	var videos []model.Video
	if err := DB.Table("videos").Where("status = 1").Scopes(func(db *gorm.DB) *gorm.DB {
		if categoryID > 0 {
			db = db.Where("category_id = ?", categoryID)
		}
		if genre != "" {
			db = db.Where("genre LIKE ?", fmt.Sprintf("%%%s%%", genre))
		}
		if region != "" {
			db = db.Where("region LIKE ?", fmt.Sprintf("%%%s%%", region))
		}
		if language != "" {
			db = db.Where("language LIKE ?", fmt.Sprintf("%%%s%%", language))
		}
		if year > 0 {
			db = db.Where("year = ?", year)
		}
		return db
	}).
		Offset(offset).
		Limit(pageSize).
		Order("is_top DESC, created_at DESC").
		Find(&videos).Error; err != nil {
		return nil, err
	}

	return &response.PageData{
		List:     videos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// assignTags 关联标签
func (s *VideoService) assignTags(video *model.Video, tagIDs []uint) {
	var tags []model.Tag
	for _, tagID := range tagIDs {
		var tag model.Tag
		if err := db.GetDB().First(&tag, tagID).Error; err == nil {
			tags = append(tags, tag)
		}
	}
	if len(tags) > 0 {
		db.GetDB().Model(video).Association("Tags").Replace(tags)
	}
}

// getCurrentDatabaseName 获取当前数据库名称
func getCurrentDatabaseName(DB *gorm.DB) string {
	// 尝试从连接字符串中提取数据库名
	// 如果无法获取，返回默认值 "gocms_v3"
	var dbName string
	if err := DB.Raw("SELECT DATABASE()").Scan(&dbName).Error; err != nil {
		// 如果无法获取，返回默认值
		return "gocms_v3"
	}
	if dbName == "" {
		return "gocms_v3"
	}
	return dbName
}
