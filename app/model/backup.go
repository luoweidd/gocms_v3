package model

import "time"

// BackupRecord 备份记录表
type BackupRecord struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	BackupType  string `gorm:"column:backup_type;type:varchar(20)" json:"backup_type"` // full/incremental
	FilePath    string `gorm:"column:file_path;type:varchar(500);not null" json:"file_path"`
	FileSize    int64  `gorm:"column:file_size;type:bigint;not null;default:0" json:"file_size"`
	Status      int8   `gorm:"column:status;type:tinyint;not null;default:0" json:"status"` // 0创建中 1成功 2失败
	ErrorMsg    string `gorm:"column:error_msg;type:varchar(500)" json:"error_msg"`
	TablesCount int    `gorm:"column:tables_count;type:int;default:0" json:"tables_count"`
	CreatedBy   uint   `gorm:"column:created_by;type:bigint unsigned;default:0" json:"created_by"`
}

// TableName 指定表名
func (BackupRecord) TableName() string {
	return "backup_records"
}

// BackupRecordQuery 备份记录查询参数
type BackupRecordQuery struct {
	Keyword    string `form:"keyword"`
	BackupType string `form:"backup_type"`
	Status     *int8  `form:"status"`
	Page       int    `form:"page" binding:"required"`
	PageSize   int    `form:"page_size" binding:"required"`
}

// BackupRecordResponse 备份记录响应
type BackupRecordResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	BackupType  string    `json:"backup_type"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	Status      int8      `json:"status"`
	ErrorMsg    string    `json:"error_msg"`
	TablesCount int       `json:"tables_count"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}

// BackupCreateRequest 创建备份请求
type BackupCreateRequest struct {
	Name       string `form:"name" json:"name"`
	BackupType string `form:"backup_type" binding:"required" json:"backup_type"`
}

// BackupStats 备份统计
type BackupStats struct {
	TotalCount   int64 `json:"total_count"`
	SuccessCount int64 `json:"success_count"`
	FailCount    int64 `json:"fail_count"`
	TotalSize    int64 `json:"total_size"`
}
