package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
	Once   sync.Once
)

// Config 日志配置
type Config struct {
	Level    string          `yaml:"level"`    // debug, info, warn, error
	Format   string          `yaml:"format"`   // console, json
	Output   string          `yaml:"output"`   // stdout, file
	FileCfg  FileConfig      `yaml:"file"`     // 文件配置
	Sampling *SamplingConfig `yaml:"sampling"` // 采样配置
}

// FileConfig 文件配置
type FileConfig struct {
	Path       string `yaml:"path"`        // 日志文件路径
	MaxSize    int    `yaml:"max_size"`    // 单个日志文件最大大小(MB)
	MaxAge     int    `yaml:"max_age"`     // 最大保留天数
	MaxBackups int    `yaml:"max_backups"` // 最大备份文件数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}

// SamplingConfig 采样配置
type SamplingConfig struct {
	Initial    int `yaml:"initial"`    // 初始采样率
	Thereafter int `yaml:"thereafter"` // 后续采样率
}

// InitLogger 初始化日志
func InitLogger(cfg Config) {
	Once.Do(func() {
		level := parseLevel(cfg.Level)
		encoder := getEncoder(cfg.Format)

		var cores []zapcore.Core

		// 控制台输出
		if cfg.Output == "stdout" || cfg.Output == "both" {
			consoleEncoder := zapcore.Lock(os.Stdout)
			cores = append(cores, zapcore.NewCore(
				encoder,
				consoleEncoder,
				level,
			))
		}

		// 文件输出
		if cfg.Output == "file" || cfg.Output == "both" {
			if cfg.FileCfg.Path != "" {
				writeSyncer := getWriteSyncer(cfg.FileCfg)
				cores = append(cores, zapcore.NewCore(
					encoder,
					writeSyncer,
					level,
				))
			}
		}

		opts := []zap.Option{
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		}

		if cfg.Sampling != nil {
			opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(
					core,
					0, // initial delay
					defaultInitial,
					defaultThereafter,
				)
			}))
		}

		Logger = zap.New(zapcore.NewTee(cores...), opts...)
	})
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		// 默认使用开发模式
		devLogger, _ := zap.NewDevelopment()
		return devLogger
	}
	return Logger
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoder 获取编码器
func getEncoder(format string) zapcore.Encoder {
	if format == "json" {
		return zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		})
	}
	// 默认使用控制台编码器
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
}

// getWriteSyncer 获取写入器
func getWriteSyncer(fileCfg FileConfig) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileCfg.Path,
		MaxSize:    fileCfg.MaxSize,
		MaxBackups: fileCfg.MaxBackups,
		MaxAge:     fileCfg.MaxAge,
		Compress:   fileCfg.Compress,
	})
}

// 默认采样配置
var (
	defaultInitial    = 100
	defaultThereafter = 100
)
