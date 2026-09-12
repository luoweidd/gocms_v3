package model

// FileUpload 文件分片上传模型
type FileUpload struct {
	BaseModel
	UploadID   string `gorm:"uniqueIndex;size:100" json:"upload_id"`
	Filename   string `gorm:"size:255" json:"filename"`
	Size       int64  `json:"size"`
	ChunkNum   int    `json:"chunk_num"`
	Status     int    `gorm:"default:0" json:"status"` // 0初始化 1完成 2失败
	FilePath   string `gorm:"size:500" json:"file_path"`
	MD5        string `gorm:"size:32" json:"md5"`
	ActualSize int64  `json:"actual_size"`
}

func (FileUpload) TableName() string {
	return "file_uploads"
}
