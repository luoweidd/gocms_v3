package cron

import (
	"log"
	"strconv"
	"sync"
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model"

	"github.com/robfig/cron/v3"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	cron    *cron.Cron
	tasks   map[string]*CronTask
	mu      sync.RWMutex
	running bool
}

// CronTask 定时任务定义
type CronTask struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	CronSpec string    `json:"cron_spec"` // cron表达式
	Handler  func()    `json:"-"`
	NextRun  time.Time `json:"next_run"`
	PrevRun  time.Time `json:"prev_run"`
	Disabled bool      `json:"disabled"`
}

// NewScheduler 创建调度器实例
func NewScheduler() *Scheduler {
	return &Scheduler{
		cron:  cron.New(),
		tasks: make(map[string]*CronTask),
	}
}

// GlobalScheduler 全局调度器实例
var GlobalScheduler *Scheduler

// GetScheduler 获取全局调度器单例
func GetScheduler() *Scheduler {
	if GlobalScheduler == nil {
		GlobalScheduler = NewScheduler()
	}
	return GlobalScheduler
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	s.cron.Start()
	s.running = true

	// 注册默认任务
	s.registerDefaultTasks()

	log.Println("定时任务调度器已启动")
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cron != nil {
		s.cron.Stop()
	}
	s.running = false
	log.Println("定时任务调度器已停止")
}

// AddTask 添加定时任务
func (s *Scheduler) AddTask(task *CronTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task.Disabled {
		return nil
	}

	entryID, err := s.cron.AddFunc(task.CronSpec, task.Handler)
	if err != nil {
		return err
	}

	task.ID = strconv.FormatInt(int64(entryID), 10)
	s.tasks[task.ID] = task

	log.Printf("定时任务已添加: %s cron=%s", task.Name, task.CronSpec)
	return nil
}

// RemoveTask 移除定时任务
func (s *Scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entryID, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr == nil {
		s.cron.Remove(cron.EntryID(entryID))
	}
	delete(s.tasks, id)
	log.Printf("定时任务已移除: %s", id)
}

// registerDefaultTasks 注册默认任务
func (s *Scheduler) registerDefaultTasks() {
	// 每5分钟检查待发布内容
	s.cron.AddFunc("*/5 * * * *", func() {
		s.publishPendingContent()
	})

	// 每天凌晨2点清理过期数据
	s.cron.AddFunc("0 2 * * *", func() {
		s.cleanupExpiredData()
	})

	log.Println("默认定时任务已注册")
}

// publishPendingContent 发布待发布的内容
func (s *Scheduler) publishPendingContent() {
	now := time.Now()

	// 发布到期的文章
	result := db.GetDB().Model(&model.Article{}).
		Where("status = ? AND published_at <= ?", 0, now).
		Update("status", 1)

	if result.Error == nil && result.RowsAffected > 0 {
		log.Printf("已自动发布 %d 篇文章", result.RowsAffected)
	}

	// 发布到期的视频
	result = db.GetDB().Model(&model.Video{}).
		Where("status = ? AND published_at <= ?", 0, now).
		Update("status", 1)

	if result.Error == nil && result.RowsAffected > 0 {
		log.Printf("已自动发布 %d 个视频", result.RowsAffected)
	}
}

// cleanupExpiredData 清理过期数据
func (s *Scheduler) cleanupExpiredData() {
	cutoffDate := time.Now().Add(-90 * 24 * time.Hour) // 90天前

	// TODO: 清理过期的刷新令牌（model.Token 待创建）
	// db.GetDB().Where("type = 'refresh' AND expires_at < ?", cutoffDate).Delete(&model.Token{})

	// TODO: 清理过期的操作日志（model.OperationLog 待创建）
	// db.GetDB().Where("created_at < ?", cutoffDate).Delete(&model.OperationLog{})

	log.Printf("已清理 %s 之前的过期数据", cutoffDate.Format(time.DateOnly))
}

// GetTaskList 获取所有定时任务列表
func (s *Scheduler) GetTaskList() []*CronTask {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*CronTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

// GetTaskEntry 获取任务条目信息
func (s *Scheduler) GetTaskEntry(id string) (cron.Entry, bool) {
	entryID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return cron.Entry{}, false
	}
	entry := s.cron.Entry(cron.EntryID(entryID))
	return entry, entry.ID == cron.EntryID(entryID)
}

// PauseTask 暂停任务
func (s *Scheduler) PauseTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if task, ok := s.tasks[id]; ok {
		task.Disabled = true
	}
	entryID, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr == nil {
		s.cron.Remove(cron.EntryID(entryID))
	}
	return nil
}

// ResumeTask 恢复任务
func (s *Scheduler) ResumeTask(task *CronTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.Disabled = false
	_, err := s.cron.AddFunc(task.CronSpec, task.Handler)
	return err
}
