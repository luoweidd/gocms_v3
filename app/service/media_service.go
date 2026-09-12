package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gocms_v3/app/db"
	"gocms_v3/app/model"
	"gocms_v3/app/response"
)

// MediaService 媒体服务
type MediaService struct{}

// NewMediaService 创建媒体服务
func NewMediaService() *MediaService {
	return &MediaService{}
}

// ========== 媒体资源相关 ==========

// CreateMediaAssetRequest 创建媒体资源请求
type CreateMediaAssetRequest struct {
	UserID       uint     `json:"-"`
	UploadID     string   `json:"upload_id" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	OriginalName string   `json:"original_name"`
	FilePath     string   `json:"file_path" binding:"required"`
	FileType     string   `json:"file_type" binding:"required"`
	MimeType     string   `json:"mime_type"`
	FileSize     int64    `json:"file_size" binding:"required"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Duration     int      `json:"duration"`
	Thumbnail    string   `json:"thumbnail"`
	Description  string   `json:"description"`
	Tags         []string `json:"tags"`
}

// CreateMediaAsset 创建媒体资源
func (s *MediaService) CreateMediaAsset(req CreateMediaAssetRequest) (*model.MediaAsset, error) {
	log.Printf("[SERVICE] 创建媒体资源请求 - 用户ID: %d, 文件名: %s", req.UserID, req.Name)

	// 序列化标签
	tagsJSON := "[]"
	if len(req.Tags) > 0 {
		b, _ := json.Marshal(req.Tags)
		tagsJSON = string(b)
	}

	asset := &model.MediaAsset{
		UserID:       uint64(req.UserID),
		UploadID:     req.UploadID,
		Name:         req.Name,
		OriginalName: req.OriginalName,
		FilePath:     req.FilePath,
		FileType:     req.FileType,
		MimeType:     req.MimeType,
		FileSize:     req.FileSize,
		Width:        req.Width,
		Height:       req.Height,
		Duration:     req.Duration,
		Thumbnail:    req.Thumbnail,
		Description:  req.Description,
		Tags:         tagsJSON,
		Status:       0, // 默认待审核状态
	}

	if err := db.GetDB().Create(asset).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 创建媒体资源失败 - 错误: %v", err)
		return nil, fmt.Errorf("创建媒体资源失败: %w", err)
	}

	// 注意：content_audits 表有外键约束，要求 content_id 必须存在于 articles 或 videos 表中。
	// 媒体资源（image、file等）的 ID 不在 articles/videos 表中，所以不能使用 content_audits 表。
	// 未来可以考虑创建独立的 media_asset_audits 表来管理媒体资源的审核。
	// 当前跳过为媒体资源创建审核记录。
	if asset.FileType == "video" {
		// 只有视频类型可以创建 content_audits 记录（因为 videos 表存在）
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
			ContentID:     asset.ID,
			ContentTitle:  asset.Name,
			SubmitterID:   uint(asset.UserID),
			SubmitterName: fmt.Sprintf("用户%d", asset.UserID),
			NextStatus:    0, // 待审核
		})
		if err != nil {
			log.Printf("[SERVICE-WARN] 创建媒体资源审核记录失败 - 资源ID: %d, 错误: %v", asset.ID, err)
		} else {
			log.Printf("[SERVICE] 媒体资源审核记录已创建 - 资源ID: %d, 类型: %s", asset.ID, asset.FileType)
		}
	} else {
		log.Printf("[SERVICE] 跳过媒体资源审核记录创建 - 资源ID: %d, 类型: %s (使用独立审核表管理)", asset.ID, asset.FileType)
	}

	log.Printf("[SERVICE] 创建媒体资源成功 - 资源ID: %d, 文件名: %s", asset.ID, asset.Name)
	return asset, nil
}

// UpdateMediaAssetRequest 更新媒体资源请求
type UpdateMediaAssetRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Tags        []string `json:"tags"`
	AlbumID     *uint64  `json:"album_id"`
	Status      *int8    `json:"status"`
}

// UpdateMediaAsset 更新媒体资源
func (s *MediaService) UpdateMediaAsset(id uint, req UpdateMediaAssetRequest) (*model.MediaAsset, error) {
	log.Printf("[SERVICE] 更新媒体资源请求 - 资源ID: %d", id)

	var asset model.MediaAsset
	if err := db.GetDB().First(&asset, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 更新媒体资源失败 - 资源ID: %d, 原因: 资源不存在", id)
		return nil, fmt.Errorf("媒体资源不存在")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Tags != nil {
		b, _ := json.Marshal(req.Tags)
		updates["tags"] = string(b)
	}
	if req.AlbumID != nil {
		updates["album_id"] = *req.AlbumID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := db.GetDB().Model(&asset).Updates(updates).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 更新媒体资源失败 - 资源ID: %d, 错误: %v", id, err)
		return nil, fmt.Errorf("更新媒体资源失败: %w", err)
	}

	log.Printf("[SERVICE] 更新媒体资源成功 - 资源ID: %d", id)
	return &asset, nil
}

// DeleteMediaAsset 删除媒体资源
func (s *MediaService) DeleteMediaAsset(id uint) error {
	log.Printf("[SERVICE] 删除媒体资源请求 - 资源ID: %d", id)

	var asset model.MediaAsset
	if err := db.GetDB().First(&asset, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 删除媒体资源失败 - 资源ID: %d, 原因: 资源不存在", id)
		return fmt.Errorf("媒体资源不存在")
	}

	// 删除物理文件
	s.deletePhysicalFile(asset.FilePath)
	if asset.Thumbnail != "" {
		s.deletePhysicalFile(asset.Thumbnail)
	}

	if err := db.GetDB().Delete(&model.MediaAsset{}, id).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 删除媒体资源失败 - 资源ID: %d, 错误: %v", id, err)
		return fmt.Errorf("删除媒体资源失败: %w", err)
	}

	log.Printf("[SERVICE] 删除媒体资源成功 - 资源ID: %d", id)
	return nil
}

// BatchDeleteMediaAssets 批量删除媒体资源
func (s *MediaService) BatchDeleteMediaAssets(ids []uint) error {
	log.Printf("[SERVICE] 批量删除媒体资源请求 - ID数量: %d", len(ids))

	if len(ids) == 0 {
		return nil
	}

	// 先获取所有资源用于删除物理文件
	var assets []model.MediaAsset
	db.GetDB().Where("id in ?", ids).Find(&assets)
	for _, asset := range assets {
		s.deletePhysicalFile(asset.FilePath)
		if asset.Thumbnail != "" {
			s.deletePhysicalFile(asset.Thumbnail)
		}
	}

	if err := db.GetDB().Delete(&[]model.MediaAsset{}, "id in ?", ids).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 批量删除媒体资源失败 - 错误: %v", err)
		return fmt.Errorf("批量删除媒体资源失败: %w", err)
	}

	log.Printf("[SERVICE] 批量删除媒体资源成功 - 删除数量: %d", len(ids))
	return nil
}

// GetMediaAssetList 获取媒体资源列表（分页）
func (s *MediaService) GetMediaAssetList(req model.MediaAssetListRequest) (*response.PageData, error) {
	log.Printf("[SERVICE] 获取媒体资源列表 - 页码: %d, 每页: %d", req.Page, req.PageSize)

	var assets []model.MediaAsset
	var total int64

	query := db.GetDB().Model(&model.MediaAsset{})

	// 添加筛选条件
	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.FileType != "" {
		query = query.Where("file_type = ?", req.FileType)
	}
	if req.AlbumID > 0 {
		query = query.Where("album_id = ?", req.AlbumID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name LIKE ? OR original_name LIKE ? OR description LIKE ?", keyword, keyword, keyword)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 统计媒体资源数量失败 - 错误: %v", err)
		return nil, fmt.Errorf("统计失败: %w", err)
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&assets).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 查询媒体资源列表失败 - 错误: %v", err)
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	log.Printf("[SERVICE] 获取媒体资源列表成功 - 总数: %d", total)
	return &response.PageData{
		List:     assets,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetMediaAssetByID 获取媒体资源详情
func (s *MediaService) GetMediaAssetByID(id uint) (*model.MediaAsset, error) {
	log.Printf("[SERVICE] 获取媒体资源详情 - 资源ID: %d", id)

	var asset model.MediaAsset
	if err := db.GetDB().First(&asset, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 获取媒体资源详情失败 - 资源ID: %d, 原因: 资源不存在", id)
		return nil, fmt.Errorf("媒体资源不存在")
	}

	log.Printf("[SERVICE] 获取媒体资源详情成功 - 资源ID: %d", id)
	return &asset, nil
}

// ========== 相册相关 ==========

// CreateAlbumRequest 创建相册请求
type CreateAlbumRequest struct {
	Name       string `json:"name" binding:"required"`
	Desc       string `json:"desc"`
	UserID     uint   `json:"-"`
	CoverImage string `json:"cover_image"`
	IsPublic   int8   `json:"is_public"`
}

// CreateAlbum 创建相册
func (s *MediaService) CreateAlbum(req CreateAlbumRequest) (*model.Album, error) {
	log.Printf("[SERVICE] 创建相册请求 - 用户ID: %d, 名称: %s", req.UserID, req.Name)

	album := &model.Album{
		Name:        req.Name,
		Description: req.Desc,
		UserID:      uint64(req.UserID),
		CoverImage:  req.CoverImage,
		AssetCount:  0,
		IsPublic:    req.IsPublic,
		Status:      1,
	}

	if err := db.GetDB().Create(album).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 创建相册失败 - 错误: %v", err)
		return nil, fmt.Errorf("创建相册失败: %w", err)
	}

	log.Printf("[SERVICE] 创建相册成功 - 相册ID: %d", album.ID)
	return album, nil
}

// UpdateAlbumRequest 更新相册请求
type UpdateAlbumRequest struct {
	Name       *string `json:"name"`
	Desc       *string `json:"desc"`
	CoverImage *string `json:"cover_image"`
	IsPublic   *int8   `json:"is_public"`
}

// UpdateAlbum 更新相册
func (s *MediaService) UpdateAlbum(id uint, req UpdateAlbumRequest) (*model.Album, error) {
	log.Printf("[SERVICE] 更新相册请求 - 相册ID: %d", id)

	var album model.Album
	if err := db.GetDB().First(&album, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 更新相册失败 - 相册ID: %d, 原因: 相册不存在", id)
		return nil, fmt.Errorf("相册不存在")
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Desc != nil {
		updates["description"] = *req.Desc
	}
	if req.CoverImage != nil {
		updates["cover_image"] = *req.CoverImage
	}
	if req.IsPublic != nil {
		updates["is_public"] = *req.IsPublic
	}

	if err := db.GetDB().Model(&album).Updates(updates).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 更新相册失败 - 相册ID: %d, 错误: %v", id, err)
		return nil, fmt.Errorf("更新相册失败: %w", err)
	}

	log.Printf("[SERVICE] 更新相册成功 - 相册ID: %d", id)
	return &album, nil
}

// DeleteAlbum 删除相册
func (s *MediaService) DeleteAlbum(id uint) error {
	log.Printf("[SERVICE] 删除相册请求 - 相册ID: %d", id)

	var album model.Album
	if err := db.GetDB().First(&album, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 删除相册失败 - 相册ID: %d, 原因: 相册不存在", id)
		return fmt.Errorf("相册不存在")
	}

	if err := db.GetDB().Delete(&model.Album{}, id).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 删除相册失败 - 相册ID: %d, 错误: %v", id, err)
		return fmt.Errorf("删除相册失败: %w", err)
	}

	// 将相册中的资源设置为无相册
	db.GetDB().Model(&model.MediaAsset{}).Where("album_id = ?", id).Update("album_id", 0)

	log.Printf("[SERVICE] 删除相册成功 - 相册ID: %d", id)
	return nil
}

// GetAlbumList 获取相册列表（分页）
func (s *MediaService) GetAlbumList(req model.AlbumListRequest) (*response.PageData, error) {
	log.Printf("[SERVICE] 获取相册列表 - 页码: %d, 每页: %d", req.Page, req.PageSize)

	var albums []model.Album
	var total int64

	query := db.GetDB().Model(&model.Album{})

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if err := query.Count(&total).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 统计相册数量失败 - 错误: %v", err)
		return nil, fmt.Errorf("统计失败: %w", err)
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&albums).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 查询相册列表失败 - 错误: %v", err)
		return nil, fmt.Errorf("查询失败: %w", err)
	}

	log.Printf("[SERVICE] 获取相册列表成功 - 总数: %d", total)
	return &response.PageData{
		List:     albums,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetAlbumByID 获取相册详情（含资源列表）
func (s *MediaService) GetAlbumByID(id uint) (*model.Album, []model.MediaAsset, error) {
	log.Printf("[SERVICE] 获取相册详情 - 相册ID: %d", id)

	var album model.Album
	if err := db.GetDB().First(&album, id).Error; err != nil {
		log.Printf("[SERVICE-WARN] 获取相册详情失败 - 相册ID: %d, 原因: 相册不存在", id)
		return nil, nil, fmt.Errorf("相册不存在")
	}

	var assets []model.MediaAsset
	db.GetDB().Where("album_id = ? AND status = ?", id, 1).Order("created_at DESC").Find(&assets)

	log.Printf("[SERVICE] 获取相册详情成功 - 资源数量: %d", len(assets))
	return &album, assets, nil
}

// AddAssetToAlbum 添加资源到相册
func (s *MediaService) AddAssetToAlbum(albumID uint, assetIDs []uint) error {
	log.Printf("[SERVICE] 添加资源到相册 - 相册ID: %d, 资源数量: %d", albumID, len(assetIDs))

	var album model.Album
	if err := db.GetDB().First(&album, albumID).Error; err != nil {
		return fmt.Errorf("相册不存在")
	}

	if len(assetIDs) == 0 {
		return nil
	}

	// 更新资源的album_id
	if err := db.GetDB().Model(&model.MediaAsset{}).Where("id in ?", assetIDs).Update("album_id", albumID).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 添加资源到相册失败 - 错误: %v", err)
		return fmt.Errorf("添加资源失败: %w", err)
	}

	// 更新相册的资源计数
	var count int64
	db.GetDB().Model(&model.MediaAsset{}).Where("album_id = ?", albumID).Count(&count)
	db.GetDB().Model(&album).Update("asset_count", count)

	log.Printf("[SERVICE] 添加资源到相册成功")
	return nil
}

// RemoveAssetFromAlbum 从相册移除资源
func (s *MediaService) RemoveAssetFromAlbum(albumID uint, assetID uint) error {
	log.Printf("[SERVICE] 从相册移除资源 - 相册ID: %d, 资源ID: %d", albumID, assetID)

	var album model.Album
	if err := db.GetDB().First(&album, albumID).Error; err != nil {
		return fmt.Errorf("相册不存在")
	}

	if err := db.GetDB().Model(&model.MediaAsset{}).Where("id = ? AND album_id = ?", assetID, albumID).Update("album_id", 0).Error; err != nil {
		log.Printf("[SERVICE-ERROR] 从相册移除资源失败 - 错误: %v", err)
		return fmt.Errorf("移除资源失败: %w", err)
	}

	// 更新相册的资源计数
	var count int64
	db.GetDB().Model(&model.MediaAsset{}).Where("album_id = ?", albumID).Count(&count)
	db.GetDB().Model(&album).Update("asset_count", count)

	log.Printf("[SERVICE] 从相册移除资源成功")
	return nil
}

// MediaStats 媒体统计信息
type MediaStats struct {
	TotalAlbums int64 `json:"total_albums"`
	TotalAssets int64 `json:"total_assets"`
	TotalSize   int64 `json:"total_size"`
	ImageCount  int64 `json:"image_count"`
	VideoCount  int64 `json:"video_count"`
	AudioCount  int64 `json:"audio_count"`
	FileCount   int64 `json:"file_count"`
	PDFCount    int64 `json:"pdf_count"`
	OfficeCount int64 `json:"office_count"`
}

// GetMediaStats 获取媒体统计信息
func (s *MediaService) GetMediaStats(userID *uint64) (*MediaStats, error) {
	stats := &MediaStats{}

	log.Printf("[SERVICE] 获取媒体统计请求 - userID: %d", func() uint64 {
		if userID != nil {
			return *userID
		}
		return 0
	}())

	if userID != nil && *userID > 0 {
		// 按用户过滤的统计
		// 相册总数
		if err := db.GetDB().Model(&model.Album{}).Where("user_id = ?", *userID).Count(&stats.TotalAlbums).Error; err != nil {
			return nil, fmt.Errorf("统计相册失败: %w", err)
		}

		// 资源总数
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND user_id = ?", 1, *userID).Count(&stats.TotalAssets).Error; err != nil {
			log.Printf("[SERVICE-ERROR] 统计资源总数失败 - 错误: %v", err)
			return nil, fmt.Errorf("统计资源总数失败: %w", err)
		}
		log.Printf("[SERVICE] 资源总数统计结果 - TotalAssets: %d", stats.TotalAssets)

		// 总大小
		var totalSize int64
		if err := db.GetDB().Table("media_assets").Where("status = ? AND user_id = ?", 1, *userID).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize).Error; err != nil {
			return nil, fmt.Errorf("统计总大小失败: %w", err)
		}
		stats.TotalSize = totalSize
		log.Printf("[SERVICE] 总大小统计结果 - TotalSize: %d", stats.TotalSize)

		// 图片数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ? AND user_id = ?", "image", 1, *userID).Count(&stats.ImageCount).Error; err != nil {
			return nil, fmt.Errorf("统计图片数量失败: %w", err)
		}
		log.Printf("[SERVICE] 图片数量统计结果 - ImageCount: %d", stats.ImageCount)

		// 视频数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ? AND user_id = ?", "video", 1, *userID).Count(&stats.VideoCount).Error; err != nil {
			return nil, fmt.Errorf("统计视频数量失败: %w", err)
		}
		log.Printf("[SERVICE] 视频数量统计结果 - VideoCount: %d", stats.VideoCount)

		// 音频数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ? AND user_id = ?", "audio", 1, *userID).Count(&stats.AudioCount).Error; err != nil {
			return nil, fmt.Errorf("统计音频数量失败: %w", err)
		}
		log.Printf("[SERVICE] 音频数量统计结果 - AudioCount: %d", stats.AudioCount)

		// 文件数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ? AND user_id = ?", "file", 1, *userID).Count(&stats.FileCount).Error; err != nil {
			return nil, fmt.Errorf("统计文件数量失败: %w", err)
		}
		log.Printf("[SERVICE] 文件数量统计结果 - FileCount: %d", stats.FileCount)

		// PDF数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND mime_type LIKE ? AND user_id = ?", 1, "%pdf%", *userID).Count(&stats.PDFCount).Error; err != nil {
			return nil, fmt.Errorf("统计PDF数量失败: %w", err)
		}
		log.Printf("[SERVICE] PDF数量统计结果 - PDFCount: %d", stats.PDFCount)

		// Office文档数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND (mime_type LIKE '%ms-word%' OR mime_type LIKE '%word%' OR mime_type LIKE '%ms-excel%' OR mime_type LIKE '%excel%') AND user_id = ?", 1, *userID).Count(&stats.OfficeCount).Error; err != nil {
			return nil, fmt.Errorf("统计Office文档数量失败: %w", err)
		}
		log.Printf("[SERVICE] Office文档数量统计结果 - OfficeCount: %d", stats.OfficeCount)
	} else {
		// 全局统计
		// 相册总数
		if err := db.GetDB().Model(&model.Album{}).Count(&stats.TotalAlbums).Error; err != nil {
			return nil, fmt.Errorf("统计相册失败: %w", err)
		}

		// 资源总数
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ?", 1).Count(&stats.TotalAssets).Error; err != nil {
			return nil, fmt.Errorf("统计资源总数失败: %w", err)
		}

		// 总大小
		var totalSize int64
		if err := db.GetDB().Table("media_assets").Where("status = ?", 1).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize).Error; err != nil {
			return nil, fmt.Errorf("统计总大小失败: %w", err)
		}
		stats.TotalSize = totalSize

		// 图片数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "image", 1).Count(&stats.ImageCount).Error; err != nil {
			return nil, fmt.Errorf("统计图片数量失败: %w", err)
		}

		// 视频数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "video", 1).Count(&stats.VideoCount).Error; err != nil {
			return nil, fmt.Errorf("统计视频数量失败: %w", err)
		}

		// 音频数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "audio", 1).Count(&stats.AudioCount).Error; err != nil {
			return nil, fmt.Errorf("统计音频数量失败: %w", err)
		}

		// 文件数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("file_type = ? AND status = ?", "file", 1).Count(&stats.FileCount).Error; err != nil {
			return nil, fmt.Errorf("统计文件数量失败: %w", err)
		}

		// PDF数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND mime_type LIKE ?", 1, "%pdf").Count(&stats.PDFCount).Error; err != nil {
			return nil, fmt.Errorf("统计PDF数量失败: %w", err)
		}

		// Office文档数量
		if err := db.GetDB().Model(&model.MediaAsset{}).Where("status = ? AND (mime_type LIKE '%ms-word%' OR mime_type LIKE '%word%' OR mime_type LIKE '%ms-excel%' OR mime_type LIKE '%excel%')", 1).Count(&stats.OfficeCount).Error; err != nil {
			return nil, fmt.Errorf("统计Office文档数量失败: %w", err)
		}
	}

	return stats, nil
}

// ========== 工具函数 ==========

// getUploadDir 获取上传目录
func (s *MediaService) getUploadDir() string {
	uploadPath := os.Getenv("MEDIA_UPLOAD_DIR")
	if uploadPath != "" {
		return uploadPath
	}
	return "storage/uploads"
}

// deletePhysicalFile 删除物理文件
func (s *MediaService) deletePhysicalFile(filePath string) {
	if filePath == "" {
		return
	}
	// 将相对路径转换为绝对路径
	if !filepath.IsAbs(filePath) {
		baseDir := s.getUploadDir()
		filePath = filepath.Join(baseDir, filePath)
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		log.Printf("[SERVICE-WARN] 删除物理文件失败 - 路径: %s, 错误: %v", filePath, err)
	}
}

// GetFileTypeFromMIME 从MIME类型判断文件类型
func (s *MediaService) GetFileTypeFromMIME(mimeType string) string {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return "image"
	case strings.HasPrefix(mimeType, "video/"):
		return "video"
	case strings.HasPrefix(mimeType, "audio/"):
		return "audio"
	default:
		return "file"
	}
}

// generateThumbnail 生成缩略图（占位实现）
func (s *MediaService) generateThumbnail(sourcePath string) (string, error) {
	// TODO: 实现缩略图生成功能
	// - 图片：使用 imgkit or similar 库缩放
	// - 视频：提取第一帧作为缩略图
	return "", nil
}

// formatFileSize 格式化文件大小
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
