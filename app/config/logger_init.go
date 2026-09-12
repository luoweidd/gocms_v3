package config

import (
	"gocms_v3/app/logger"
)

// InitLogger 初始化日志（修复字段访问问题）
func InitLogger() {
	fileCfg := GlobalConfig.Log.File
	if fileCfg.Path == "" {
		fileCfg.Path = "./logs/app.log"
	}

	loggerCfg := logger.Config{
		Level:    GlobalConfig.Log.Level,
		Format:   GlobalConfig.Log.Format,
		Output:   GlobalConfig.Log.Output,
		FileCfg:  logger.FileConfig{Path: fileCfg.Path, MaxSize: fileCfg.MaxSize, MaxAge: fileCfg.MaxAge, MaxBackups: fileCfg.MaxBackups, Compress: fileCfg.Compress},
		Sampling: &logger.SamplingConfig{Initial: 100, Thereafter: 100},
	}

	logger.InitLogger(loggerCfg)
}
