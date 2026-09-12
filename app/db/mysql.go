package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gocms_v3/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitMySQL 初始化 MySQL 连接池
func InitMySQL(cfg config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Database, cfg.Charset)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[DB-ERROR] MySQL 连接失败 [%s:%d]: %v", cfg.Host, cfg.Port, err)
		return fmt.Errorf("MySQL 连接失败 [%s:%d]: %w", cfg.Host, cfg.Port, err)
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Printf("[DB] MySQL 连接成功：%s:%d/%s (最大连接数：%d, 空闲连接数：%d)",
		cfg.Host, cfg.Port, cfg.Database,
		cfg.MaxOpenConns, cfg.MaxIdleConns)

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// GetDBWithErr 获取数据库实例，如果未初始化则返回错误
func GetDBWithErr() (*gorm.DB, error) {
	if DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return DB, nil
}

// IsDBInitialized 检查数据库是否已初始化
func IsDBInitialized() bool {
	return DB != nil
}

// RunMigration 运行 SQL 迁移文件
func RunMigration(migrationFile string) error {
	db, err := GetDBWithErr()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取底层 database/sql.DB 失败: %w", err)
	}

	// 读取迁移文件
	filePath := filepath.Join("migrations", migrationFile)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开迁移文件失败 [%s]: %w", filePath, err)
	}
	defer file.Close()

	// 读取 SQL 内容
	scanner := bufio.NewScanner(file)
	var sqlParts []string

	for scanner.Scan() {
		line := scanner.Text()
		// 跳过注释行
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		sqlParts = append(sqlParts, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取迁移文件失败 [%s]: %w", filePath, err)
	}

	sql := strings.Join(sqlParts, "\n")

	// 按分号分割执行每条 SQL
	stmts := splitSQLStatements(sql)
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := sqlDB.Exec(stmt); err != nil {
			return fmt.Errorf("执行 SQL 失败 [%s]: %w", stmt, err)
		}
	}

	log.Printf("[MIGRATION] 迁移执行成功: %s", migrationFile)
	return nil
}

// splitSQLStatements 按分号分割 SQL 语句（处理引号内的分号）
func splitSQLStatements(sql string) []string {
	var stmts []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	i := 0

	for i < len(sql) {
		ch := sql[i]

		// 处理引号状态
		if ch == '\'' && (i == 0 || sql[i-1] != '\\') {
			inSingleQuote = !inSingleQuote
		} else if ch == '"' && (i == 0 || sql[i-1] != '\\') {
			inDoubleQuote = !inDoubleQuote
		} else if ch == ';' && !inSingleQuote && !inDoubleQuote {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
			current.Reset()
		} else {
			current.WriteByte(ch)
		}
		i++
	}

	// 添加最后一条语句
	last := strings.TrimSpace(current.String())
	if last != "" {
		stmts = append(stmts, last)
	}

	return stmts
}
