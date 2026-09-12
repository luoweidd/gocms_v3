package service

import (
	"fmt"
	"time"

	"gocms_v3/app/db"

	"github.com/gin-gonic/gin"
)

// ContentVersionService 内容版本管理服务
type ContentVersionService struct{}

// NewContentVersionService 创建内容版本服务实例
func NewContentVersionService() *ContentVersionService {
	return &ContentVersionService{}
}

// VersionRecord 内容版本记录
type VersionRecord struct {
	ID          uint      `json:"id"`
	EntityID    uint      `json:"entity_id"`
	EntityType  string    `json:"entity_type"`
	VersionNum  int       `json:"version_num"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	DiffSummary string    `json:"diff_summary"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateVersion 创建新内容版本
func (s *ContentVersionService) CreateVersion(entityID uint, entityType, title, content, diffSummary string, createdBy uint) (*VersionRecord, error) {
	// 获取当前最新版本号
	var lastVersion VersionRecord
	err := db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ?", entityID, entityType).
		Order("version_num DESC").First(&lastVersion).Error

	versionNum := 1
	if err == nil {
		versionNum = lastVersion.VersionNum + 1
	}

	record := &VersionRecord{
		EntityID:    entityID,
		EntityType:  entityType,
		VersionNum:  versionNum,
		Title:       title,
		Content:     content,
		DiffSummary: diffSummary,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	type tableRecord struct {
		ID          uint      `gorm:"column:id"`
		EntityID    uint      `gorm:"column:entity_id"`
		EntityType  string    `gorm:"column:entity_type"`
		VersionNum  int       `gorm:"column:version_num"`
		Title       string    `gorm:"column:title"`
		Content     string    `gorm:"column:content"`
		DiffSummary string    `gorm:"column:diff_summary"`
		CreatedBy   uint      `gorm:"column:created_by"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	dbRecord := tableRecord{
		EntityID:    entityID,
		EntityType:  entityType,
		VersionNum:  versionNum,
		Title:       title,
		Content:     content,
		DiffSummary: diffSummary,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	if err := db.GetDB().Table("content_versions").Create(&dbRecord).Error; err != nil {
		return nil, fmt.Errorf("创建版本失败: %w", err)
	}

	record.ID = dbRecord.ID
	return record, nil
}

// GetVersions 获取内容的版本历史
func (s *ContentVersionService) GetVersions(entityID uint, entityType string) ([]VersionRecord, error) {
	var versions []VersionRecord
	type versionRow struct {
		ID          uint      `gorm:"column:id"`
		EntityID    uint      `gorm:"column:entity_id"`
		EntityType  string    `gorm:"column:entity_type"`
		VersionNum  int       `gorm:"column:version_num"`
		Title       string    `gorm:"column:title"`
		Content     string    `gorm:"column:content"`
		DiffSummary string    `gorm:"column:diff_summary"`
		CreatedBy   uint      `gorm:"column:created_by"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	var rows []versionRow
	if err := db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ?", entityID, entityType).
		Order("version_num DESC").Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("获取版本历史失败: %w", err)
	}

	for _, r := range rows {
		versions = append(versions, VersionRecord{
			ID:          r.ID,
			EntityID:    r.EntityID,
			EntityType:  r.EntityType,
			VersionNum:  r.VersionNum,
			Title:       r.Title,
			Content:     r.Content,
			DiffSummary: r.DiffSummary,
			CreatedBy:   r.CreatedBy,
			CreatedAt:   r.CreatedAt,
		})
	}

	return versions, nil
}

// GetVersion 获取指定版本
func (s *ContentVersionService) GetVersion(id uint) (*VersionRecord, error) {
	type versionRow struct {
		ID          uint      `gorm:"column:id"`
		EntityID    uint      `gorm:"column:entity_id"`
		EntityType  string    `gorm:"column:entity_type"`
		VersionNum  int       `gorm:"column:version_num"`
		Title       string    `gorm:"column:title"`
		Content     string    `gorm:"column:content"`
		DiffSummary string    `gorm:"column:diff_summary"`
		CreatedBy   uint      `gorm:"column:created_by"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	var row versionRow
	if err := db.GetDB().Table("content_versions").First(&row, id).Error; err != nil {
		return nil, fmt.Errorf("获取版本失败: %w", err)
	}

	return &VersionRecord{
		ID:          row.ID,
		EntityID:    row.EntityID,
		EntityType:  row.EntityType,
		VersionNum:  row.VersionNum,
		Title:       row.Title,
		Content:     row.Content,
		DiffSummary: row.DiffSummary,
		CreatedBy:   row.CreatedBy,
		CreatedAt:   row.CreatedAt,
	}, nil
}

// RollbackToVersion 回滚到指定版本
func (s *ContentVersionService) RollbackToVersion(c *gin.Context, entityID uint, entityType string, targetVersionNum int) (*VersionRecord, error) {
	type versionRow struct {
		Title     string `gorm:"column:title"`
		Content   string `gorm:"column:content"`
		CreatedBy uint   `gorm:"column:created_by"`
	}

	var targetVersion versionRow
	if err := db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ? AND version_num = ?", entityID, entityType, targetVersionNum).
		First(&targetVersion).Error; err != nil {
		return nil, fmt.Errorf("找不到目标版本: %w", err)
	}

	var lastVersion struct {
		VersionNum int `gorm:"column:version_num"`
	}
	db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ?", entityID, entityType).
		Order("version_num DESC").First(&lastVersion)

	newVersionNum := lastVersion.VersionNum + 1
	diffSummary := fmt.Sprintf("回滚到版本%d", targetVersionNum)

	record := &VersionRecord{
		EntityID:    entityID,
		EntityType:  entityType,
		VersionNum:  newVersionNum,
		Title:       targetVersion.Title,
		Content:     targetVersion.Content,
		DiffSummary: diffSummary,
		CreatedBy:   targetVersion.CreatedBy,
		CreatedAt:   time.Now(),
	}

	type tableRecord struct {
		ID          uint      `gorm:"column:id"`
		EntityID    uint      `gorm:"column:entity_id"`
		EntityType  string    `gorm:"column:entity_type"`
		VersionNum  int       `gorm:"column:version_num"`
		Title       string    `gorm:"column:title"`
		Content     string    `gorm:"column:content"`
		DiffSummary string    `gorm:"column:diff_summary"`
		CreatedBy   uint      `gorm:"column:created_by"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}

	dbRecord := tableRecord{
		EntityID:    entityID,
		EntityType:  entityType,
		VersionNum:  newVersionNum,
		Title:       targetVersion.Title,
		Content:     targetVersion.Content,
		DiffSummary: diffSummary,
		CreatedBy:   targetVersion.CreatedBy,
		CreatedAt:   time.Now(),
	}

	if err := db.GetDB().Table("content_versions").Create(&dbRecord).Error; err != nil {
		return nil, fmt.Errorf("创建回滚版本失败: %w", err)
	}

	record.ID = dbRecord.ID
	return record, nil
}

// CompareVersions 比较两个版本的差异
func (s *ContentVersionService) CompareVersions(entityID uint, entityType string, v1, v2 int) (map[string]interface{}, error) {
	type versionRow struct {
		Title   string `gorm:"column:title"`
		Content string `gorm:"column:content"`
	}

	var ver1, ver2 versionRow
	if err := db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ? AND version_num = ?", entityID, entityType, v1).First(&ver1).Error; err != nil {
		return nil, fmt.Errorf("版本%d不存在: %w", v1, err)
	}

	if err := db.GetDB().Table("content_versions").
		Where("entity_id = ? AND entity_type = ? AND version_num = ?", entityID, entityType, v2).First(&ver2).Error; err != nil {
		return nil, fmt.Errorf("版本%d不存在: %w", v2, err)
	}

	return map[string]interface{}{
		"title_diff":   compareStrings(ver1.Title, ver2.Title),
		"content_diff": compareContentDiff(ver1.Content, ver2.Content),
		"v1_lines":     len(ver1.Content),
		"v2_lines":     len(ver2.Content),
	}, nil
}

func compareStrings(s1, s2 string) map[string]string {
	return map[string]string{"original": s1, "current": s2}
}

func compareContentDiff(c1, c2 string) map[string]string {
	return map[string]string{"original": c1, "current": c2}
}

// CreateArticleVersion 创建文章版本
func (s *ContentVersionService) CreateArticleVersion(c *gin.Context, articleID uint, title, content, diffSummary string) (*VersionRecord, error) {
	createdBy := c.GetUint("user_id")
	return s.CreateVersion(articleID, "article", title, content, diffSummary, createdBy)
}

// GetArticleVersions 获取文章版本历史
func (s *ContentVersionService) GetArticleVersions(c *gin.Context, articleID uint) ([]VersionRecord, error) {
	return s.GetVersions(articleID, "article")
}

// CreateVideoVersion 创建视频版本
func (s *ContentVersionService) CreateVideoVersion(c *gin.Context, videoID uint, title, content, diffSummary string) (*VersionRecord, error) {
	createdBy := c.GetUint("user_id")
	return s.CreateVersion(videoID, "video", title, content, diffSummary, createdBy)
}

// GetVideoVersions 获取视频版本历史
func (s *ContentVersionService) GetVideoVersions(c *gin.Context, videoID uint) ([]VersionRecord, error) {
	return s.GetVersions(videoID, "video")
}
