package model

// MediaAsset 媒体资源
// 注意：MySQL 不支持部分外键约束(partitioned foreign keys)，所以 MediaAsset
// 不能直接关联到 ContentAudit。审核记录需要通过其他方式管理。
type MediaAsset struct {
	BaseModel
	UserID       uint64 `gorm:"index" json:"user_id"`
	UploadID     string `gorm:"size:100;index" json:"upload_id"`
	Name         string `gorm:"size:255" json:"name"`
	OriginalName string `gorm:"size:255" json:"original_name"`
	FilePath     string `gorm:"size:500" json:"file_path"`
	FileType     string `gorm:"size:20;index" json:"file_type"` // file/image/video/audio
	MimeType     string `gorm:"size:100" json:"mime_type"`
	FileSize     int64  `json:"file_size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Duration     int    `json:"duration"`
	Thumbnail    string `gorm:"size:500" json:"thumbnail"`
	AlbumID      uint64 `gorm:"index" json:"album_id"`
	Tags         string `gorm:"type:text" json:"tags"` // JSON array
	Description  string `gorm:"type:text" json:"description"`
	Status       int8   `gorm:"default:1;index" json:"status"`
}

func (MediaAsset) TableName() string {
	return "media_assets"
}

// Album 相册
type Album struct {
	BaseModel
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	UserID      uint64 `gorm:"index;not null" json:"user_id"`
	CoverImage  string `gorm:"size:500" json:"cover_image"`
	AssetCount  int    `gorm:"default:0" json:"asset_count"`
	IsPublic    int8   `gorm:"default:0" json:"is_public"`
	Status      int8   `gorm:"default:1" json:"status"`
}

func (Album) TableName() string {
	return "albums"
}

// MediaAssetListRequest 媒体资源列表查询请求
type MediaAssetListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	UserID   uint64 `form:"user_id"`
	FileType string `form:"file_type"`
	AlbumID  uint64 `form:"album_id"`
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
}

// AlbumListRequest 相册列表查询请求
type AlbumListRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	UserID   uint64 `form:"user_id"`
	Name     string `form:"name"`
}
