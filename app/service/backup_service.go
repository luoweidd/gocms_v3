package service

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gocms_v3/app/config"
	"gocms_v3/app/model"

	"go.uber.org/zap"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// BackupService 备份服务
type BackupService struct {
	db  *gorm.DB
	log *zap.Logger
	cfg *config.BackupConfig
}

// NewBackupService 创建备份服务
func NewBackupService(db *gorm.DB, cfg *config.Config) *BackupService {
	return &BackupService{
		db:  db,
		log: zap.L(),
		cfg: &cfg.Backup,
	}
}

// CreateBackup 创建备份
func (s *BackupService) CreateBackup(ctx context.Context, req model.BackupCreateRequest, createdBy uint) (*model.BackupRecord, error) {
	// 确保存储目录存在
	if err := os.MkdirAll(s.cfg.StoragePath, 0755); err != nil {
		return nil, fmt.Errorf("创建备份目录失败: %w", err)
	}

	// 生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("backup_%s.sql", timestamp)
	filePath := filepath.Join(s.cfg.StoragePath, filename)

	// 创建备份记录
	record := &model.BackupRecord{
		Name:       req.Name,
		BackupType: req.BackupType,
		FilePath:   filePath,
		Status:     0, // 创建中
		CreatedBy:  createdBy,
	}

	if err := s.db.Create(record).Error; err != nil {
		return nil, fmt.Errorf("创建备份记录失败: %w", err)
	}

	// 执行备份（异步）
	go func() {
		err := s.executeBackup(req.BackupType, filePath, record.ID)
		if err != nil {
			s.log.Error("备份失败", zap.Error(err), zap.Uint("recordID", record.ID))
			s.db.Model(&model.BackupRecord{}).Where("id = ?", record.ID).Updates(map[string]interface{}{
				"status":    2,
				"error_msg": err.Error(),
			})
		} else {
			s.log.Info("备份成功", zap.Uint("recordID", record.ID), zap.String("filePath", filePath))
		}
	}()

	return record, nil
}

// executeBackup 执行实际备份操作 - 使用Go代码生成SQL dump
func (s *BackupService) executeBackup(backupType, filePath string, recordID uint) error {
	// 创建远程数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		s.cfg.MySQL.User,
		s.cfg.MySQL.Password,
		s.cfg.MySQL.Host,
		s.cfg.MySQL.Port,
		s.cfg.MySQL.Database,
	)

	remoteDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接远程数据库失败: %w", err)
	}
	defer remoteDB.Close()

	if err := remoteDB.Ping(); err != nil {
		return fmt.Errorf(" ping远程数据库失败: %w", err)
	}

	// 创建输出文件
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建备份文件失败: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// 写入SQL头部
	writer.WriteString("-- GoCMS Database Backup\n")
	writer.WriteString(fmt.Sprintf("-- Backup Time: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	writer.WriteString(fmt.Sprintf("-- Backup Type: %s\n", backupType))
	writer.WriteString("\n\n")

	// 获取所有表名
	var tables []string
	if err := s.db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return fmt.Errorf("获取表列表失败: %w", err)
	}

	// 设置SQL模式
	writer.WriteString("SET NAMES utf8mb4;\n")
	writer.WriteString("SET FOREIGN_KEY_CHECKS = 0;\n\n")

	// 遍历每个表
	tableCount := 0
	for _, table := range tables {
		// 生成CREATE TABLE语句
		createSQL, err := s.getCreateTableSQL(remoteDB, table)
		if err != nil {
			writer.WriteString(fmt.Sprintf("-- Error getting create SQL for table %s: %v\n\n", table, err))
			continue
		}
		writer.WriteString(fmt.Sprintf("-- ----------------------------\n"))
		writer.WriteString(fmt.Sprintf("-- Table structure for %s\n", table))
		writer.WriteString(fmt.Sprintf("-- ----------------------------\n"))
		writer.WriteString(createSQL)
		writer.WriteString("\n\n")

		// 生成INSERT语句
		insertSQL, err := s.getTableDataSQL(remoteDB, table)
		if err != nil {
			writer.WriteString(fmt.Sprintf("-- Error getting data for table %s: %v\n\n", table, err))
			continue
		}
		writer.WriteString(insertSQL)
		writer.WriteString("\n\n")

		tableCount++
	}

	// 写入尾部
	writer.WriteString("SET FOREIGN_KEY_CHECKS = 1;\n")

	// 更新记录
	fileInfo, _ := os.Stat(filePath)
	s.db.Model(&model.BackupRecord{}).Where("id = ?", recordID).Updates(map[string]interface{}{
		"status":       1,
		"file_size":    fileInfo.Size(),
		"tables_count": tableCount,
		"error_msg":    "",
	})

	return nil
}

// getCreateTableSQL 获取表的CREATE TABLE语句
func (s *BackupService) getCreateTableSQL(db *sql.DB, table string) (string, error) {
	var createSQL string
	err := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE `%s`", table)).Scan(&table, &createSQL)
	if err != nil {
		return "", fmt.Errorf("获取表结构失败: %w", err)
	}
	return createSQL + ";\n", nil
}

// getTableDataSQL 获取表的所有数据INSERT语句
func (s *BackupService) getTableDataSQL(db *sql.DB, table string) (string, error) {
	// 获取列信息
	columns, err := s.getTableColumns(db, table)
	if err != nil {
		return "", fmt.Errorf("获取列信息失败: %w", err)
	}

	if len(columns) == 0 {
		return "", nil
	}

	// 构建查询
	selectCols := make([]string, len(columns))
	for i, col := range columns {
		selectCols[i] = fmt.Sprintf("`%s`", col.Name)
	}
	query := fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(selectCols, ", "), table)

	// 执行查询
	rows, err := db.Query(query)
	if err != nil {
		return "", fmt.Errorf("查询数据失败: %w", err)
	}
	defer rows.Close()

	// 获取列类型
	types, err := rows.ColumnTypes()
	if err != nil {
		return "", fmt.Errorf("获取列类型失败: %w", err)
	}

	var results []string
	for rows.Next() {
		// 创建扫描目标
		values := make([]interface{}, len(types))
		for i, t := range types {
			switch t.DatabaseTypeName() {
			case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "BIGINT", "FLOAT", "DOUBLE", "DECIMAL":
				values[i] = new(sql.NullFloat64)
			case "DATE", "TIME", "DATETIME", "TIMESTAMP", "YEAR":
				values[i] = new(sql.NullString)
			case "TEXT", "VARCHAR", "CHAR", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT", "ENUM", "SET":
				values[i] = new(sql.NullString)
			default: // BLOB, LONGBLOB, LONGTEXT等
				values[i] = new(sql.NullString)
			}
		}

		if err := rows.Scan(values...); err != nil {
			continue
		}

		// 构建INSERT语句
		valuesStr := make([]string, len(values))
		for i, v := range values {
			switch val := v.(type) {
			case *sql.NullFloat64:
				if val.Valid {
					valuesStr[i] = fmt.Sprintf("%g", val.Float64)
				} else {
					valuesStr[i] = "NULL"
				}
			case *sql.NullString:
				if val.Valid {
					escaped := strings.ReplaceAll(val.String, "\\", "\\\\")
					escaped = strings.ReplaceAll(escaped, "'", "\\'")
					escaped = strings.ReplaceAll(escaped, "\"", "\\\"")
					valuesStr[i] = fmt.Sprintf("'%s'", escaped)
				} else {
					valuesStr[i] = "NULL"
				}
			default:
				if val == nil {
					valuesStr[i] = "NULL"
				} else {
					valuesStr[i] = fmt.Sprintf("'%v'", val)
				}
			}
		}

		results = append(results, fmt.Sprintf("INSERT INTO `%s` VALUES (%s);", table, strings.Join(valuesStr, ",")))
	}

	return strings.Join(results, "\n"), nil
}

// getTableColumns 获取表的列名
func (s *BackupService) getTableColumns(db *sql.DB, table string) ([]struct {
	Name string
}, error) {
	rows, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []struct {
		Name string
	}
	for rows.Next() {
		var field, typ, null, key, defaul, extra string
		var precision, scale int
		if err := rows.Scan(&field, &typ, &null, &key, &defaul, &extra, &precision, &scale); err != nil {
			// 兼容旧版本MySQL
			continue
		}
		columns = append(columns, struct{ Name string }{Name: field})
	}

	// 如果上面的Scan失败，尝试简化的方式
	if len(columns) == 0 {
		rows2, err2 := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", table))
		if err2 != nil {
			return nil, err2
		}
		defer rows2.Close()

		for rows2.Next() {
			var scanResult []interface{}
			columns2, _ := rows2.ColumnTypes()
			scanResult = make([]interface{}, len(columns2))
			for i := range scanResult {
				scanResult[i] = new(interface{})
			}
			if err := rows2.Scan(scanResult...); err != nil {
				continue
			}
			if strVal, ok := scanResult[0].(*string); ok && strVal != nil {
				columns = append(columns, struct{ Name string }{Name: *strVal})
			}
		}
	}

	return columns, nil
}

// GetBackupList 获取备份列表
func (s *BackupService) GetBackupList(query model.BackupRecordQuery) ([]model.BackupRecordResponse, int64, error) {
	var records []model.BackupRecord
	var total int64

	queryBuilder := s.db.Model(&model.BackupRecord{}).Order("created_at DESC")

	if query.Keyword != "" {
		queryBuilder = queryBuilder.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.BackupType != "" {
		queryBuilder = queryBuilder.Where("backup_type = ?", query.BackupType)
	}
	if query.Status != nil {
		queryBuilder = queryBuilder.Where("status = ?", *query.Status)
	}

	if err := queryBuilder.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询备份总数失败: %w", err)
	}

	offset := (query.Page - 1) * query.PageSize
	if err := queryBuilder.Offset(offset).Limit(query.PageSize).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询备份列表失败: %w", err)
	}

	resp := make([]model.BackupRecordResponse, len(records))
	for i, r := range records {
		resp[i] = model.BackupRecordResponse{
			ID:          r.ID,
			Name:        r.Name,
			BackupType:  r.BackupType,
			FilePath:    r.FilePath,
			FileSize:    r.FileSize,
			Status:      r.Status,
			ErrorMsg:    r.ErrorMsg,
			TablesCount: r.TablesCount,
			CreatedBy:   r.CreatedBy,
			CreatedAt:   r.CreatedAt,
		}
	}

	return resp, total, nil
}

// GetBackupDetail 获取备份详情
func (s *BackupService) GetBackupDetail(id uint) (*model.BackupRecordResponse, error) {
	var record model.BackupRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, fmt.Errorf("备份记录不存在: %w", err)
	}

	return &model.BackupRecordResponse{
		ID:          record.ID,
		Name:        record.Name,
		BackupType:  record.BackupType,
		FilePath:    record.FilePath,
		FileSize:    record.FileSize,
		Status:      record.Status,
		ErrorMsg:    record.ErrorMsg,
		TablesCount: record.TablesCount,
		CreatedBy:   record.CreatedBy,
		CreatedAt:   record.CreatedAt,
	}, nil
}

// DeleteBackup 删除备份
func (s *BackupService) DeleteBackup(id uint) error {
	var record model.BackupRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return fmt.Errorf("备份记录不存在: %w", err)
	}

	// 删除文件
	if record.FilePath != "" {
		if err := os.Remove(record.FilePath); err != nil {
			s.log.Warn("删除备份文件失败", zap.Error(err), zap.String("filePath", record.FilePath))
		}
	}

	// 删除记录
	return s.db.Delete(&model.BackupRecord{}, id).Error
}

// DownloadBackup 获取备份文件完整路径（用于下载）
func (s *BackupService) DownloadBackup(id uint) (string, error) {
	var record model.BackupRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return "", fmt.Errorf("备份记录不存在: %w", err)
	}

	if record.Status != 1 {
		return "", fmt.Errorf("备份文件不可用")
	}

	// 如果路径是相对的，转换为绝对路径
	if !filepath.IsAbs(record.FilePath) {
		// 获取当前工作目录
		workDir, err := os.Getwd()
		if err == nil {
			return filepath.Join(workDir, record.FilePath), nil
		}
	}

	return record.FilePath, nil
}

// RestoreBackup 恢复备份（需要二次确认，返回需要执行的 SQL 文件路径）
func (s *BackupService) RestoreBackup(id uint) (string, error) {
	var record model.BackupRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return "", fmt.Errorf("备份记录不存在: %w", err)
	}

	if record.Status != 1 {
		return "", fmt.Errorf("备份文件不可用")
	}

	if _, err := os.Stat(record.FilePath); os.IsNotExist(err) {
		return "", fmt.Errorf("备份文件不存在")
	}

	// TODO: 实际恢复逻辑需要先导出表结构，然后导入数据
	// 这里只返回文件路径供后续处理
	return record.FilePath, nil
}

// GetBackupStats 获取备份统计
func (s *BackupService) GetBackupStats() (*model.BackupStats, error) {
	stats := &model.BackupStats{}

	s.db.Model(&model.BackupRecord{}).Count(&stats.TotalCount)
	s.db.Model(&model.BackupRecord{}).Where("status = ?", 1).Count(&stats.SuccessCount)
	s.db.Model(&model.BackupRecord{}).Where("status = ?", 2).Count(&stats.FailCount)
	s.db.Model(&model.BackupRecord{}).Select("COALESCE(SUM(file_size), 0)").Scan(&stats.TotalSize)

	return stats, nil
}

// CleanupOldBackups 清理过旧的备份，保留最近的 N 个
func (s *BackupService) CleanupOldBackups(maxCount int) ([]uint, error) {
	if maxCount <= 0 {
		maxCount = s.cfg.MaxBackupCount
	}

	var oldRecords []model.BackupRecord
	offset := maxCount
	s.db.Order("created_at DESC").Offset(offset).Find(&oldRecords)

	var ids []uint
	for _, r := range oldRecords {
		ids = append(ids, r.ID)
		// 删除文件
		if r.FilePath != "" {
			os.Remove(r.FilePath)
		}
	}

	if len(ids) > 0 {
		s.db.Where("id IN ?", ids).Delete(&model.BackupRecord{})
	}

	return ids, nil
}
