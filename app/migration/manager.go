package migration

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gocms_v3/app/db"
)

// Manager 迁移管理器
type Manager struct {
	migrationDir string
}

// NewManager 创建迁移管理器
func NewManager() *Manager {
	return &Manager{
		migrationDir: "migrations",
	}
}

// Migration 表示一个迁移文件
type Migration struct {
	Name     string
	Content  string
	Executed bool
}

// GetMigrations 获取所有迁移文件（按名称排序）
// 注意：迁移和种子数据分离设计
// - SQL 迁移文件负责表结构创建
// - seeder.go 负责初始化数据（角色、菜单、配置等）
func (m *Manager) GetMigrations() ([]Migration, error) {
	dir := filepath.Join(m.migrationDir)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("读取迁移目录失败: %w", err)
	}

	var migrations []Migration
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			content, err := os.ReadFile(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("读取迁移文件失败 [%s]: %w", file.Name(), err)
			}
			migrations = append(migrations, Migration{
				Name:    file.Name(),
				Content: string(content),
			})
		}
	}

	// 按文件名排序（确保执行顺序）
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

// RunMigration 执行单个迁移文件
func (m *Manager) RunMigration(filename string) error {
	dbInstance, err := db.GetDBWithErr()
	if err != nil {
		return fmt.Errorf("数据库未初始化: %w", err)
	}

	sqlDB, err := dbInstance.DB()
	if err != nil {
		return fmt.Errorf("获取底层database失败: %w", err)
	}

	filePath := filepath.Join(m.migrationDir, filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取迁移文件失败 [%s]: %w", filename, err)
	}

	// 解析并执行SQL语句
	stmts := parseSQLStatements(string(content))
	executedCount := 0
	skippedCount := 0

	// 已知的可忽略错误（继续执行后续语句）
	skipableErrors := []string{
		"Duplicate column name",                 // 列已存在
		"Duplicate key name",                    // 索引已存在
		"Duplicate entry",                       // 重复条目
		"Table already exists",                  // 表已存在
		"Duplicate foreign key",                 // 重复外键约束
		"Duplicate foreign key constraint name", // 重复外键约束名称
		"Cannot add or update a child row",      // 外键约束失败（子行不存在）
		"Foreign key constraint fails",          // 外键约束失败
		"Unknown column",                        // 列不存在（表结构不一致时）
		"doesn't have a default value",          // 字段没有默认值
	}

	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := sqlDB.Exec(stmt); err != nil {
			// 检查是否是可忽略的错误
			skipped := false
			for _, errMsg := range skipableErrors {
				if strings.Contains(err.Error(), errMsg) {
					log.Printf("[MIGRATION] ⊘ 跳过 (已存在): %s", truncateString(stmt, 80))
					skipped = true
					skippedCount++
					break
				}
			}
			if !skipped {
				return fmt.Errorf("执行SQL失败 [%s]: %w", truncateString(stmt, 100), err)
			}
		} else {
			executedCount++
		}
	}

	log.Printf("[MIGRATION] ✓ %s (成功: %d条, 跳过: %d条)", filename, executedCount, skippedCount)
	return nil
}

// RunAllMigrations 执行所有未执行的迁移
// 注意：单个迁移失败不会阻止其他迁移的执行
func (m *Manager) RunAllMigrations() error {
	migrations, err := m.GetMigrations()
	if err != nil {
		return fmt.Errorf("获取迁移列表失败: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("[MIGRATION] 未发现迁移文件")
		return nil
	}

	log.Printf("[MIGRATION] 发现 %d 个迁移文件\n", len(migrations))

	var failedMigrations []string
	executedCount := 0

	for _, migration := range migrations {
		log.Printf("[MIGRATION] 正在执行: %s ...", migration.Name)
		if err := m.RunMigration(migration.Name); err != nil {
			log.Printf("[MIGRATION] ⚠ 警告: 迁移 %s 失败: %v (将继续执行后续迁移)", migration.Name, err)
			failedMigrations = append(failedMigrations, migration.Name)
		} else {
			executedCount++
		}
	}

	if len(failedMigrations) > 0 {
		log.Printf("[MIGRATION] %d/%d 个迁移成功，%d 个失败: %v",
			executedCount, len(migrations), len(failedMigrations), failedMigrations)
		log.Println("[MIGRATION] 建议检查失败的迁移文件，可以手动执行")
	} else {
		log.Printf("[MIGRATION] 所有 %d 个迁移执行完成!", executedCount)
	}

	return nil
}

// truncateString 截断字符串用于日志显示
func truncateString(s string, maxLen int) string {
	s = strings.TrimSpace(s)
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// parseSQLStatements 解析SQL语句（处理引号内的分号）
func parseSQLStatements(sql string) []string {
	var stmts []string
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	i := 0

	for i < len(sql) {
		ch := sql[i]

		// 处理引号状态 - 检查前一个字符是否是转义字符
		prevChar := byte(0)
		if i > 0 {
			prevChar = sql[i-1]
		}

		if ch == '\'' && prevChar != '\\' {
			inSingleQuote = !inSingleQuote
			current.WriteByte(ch)
		} else if ch == '"' && prevChar != '\\' {
			inDoubleQuote = !inDoubleQuote
			current.WriteByte(ch)
		} else if ch == '\\' && inSingleQuote && i+1 < len(sql) {
			// 跳过转义字符和下一个字符
			current.WriteByte(ch)
			i++
			current.WriteByte(sql[i])
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
