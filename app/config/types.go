package config

import (
	"time"
)

// DatabaseConfig MySQL配置（直接包含MySQL字段，不再需要嵌套的Mysql键）
type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Username        string        `yaml:"username"`
	Password        string        `yaml:"password"`
	Database        string        `yaml:"database"`
	Charset         string        `yaml:"charset"`
	ParseTime       bool          `yaml:"parse_time"`
	Loc             string        `yaml:"loc"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
	SlowThreshold   time.Duration `yaml:"slow_threshold"`
}

// RedisConfig Redis配置（用于热门菜单树缓存和用户会话）
type RedisConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolTimeout  time.Duration `yaml:"pool_timeout"`
}

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	Addresses    []string      `yaml:"addresses"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	IndexPrefix  string        `yaml:"index_prefix"`
	MaxRetries   int           `yaml:"max_retries"`
	RetryTimeout time.Duration `yaml:"retry_timeout"`
	Sniff        bool          `yaml:"sniff"`
	HealthCheck  bool          `yaml:"health_check"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string `yaml:"secret"`
	Issuer          string `yaml:"issuer"`
	Algorithm       string `yaml:"algorithm"`
	AccessTokenTTL  string `yaml:"access_token_ttl"`
	RefreshTokenTTL string `yaml:"refresh_token_ttl"`
}

// UploadConfig 上传配置
type UploadConfig struct {
	Driver       string          `yaml:"driver"`
	Local        LocalUploadConf `yaml:"local"`
	MaxFileSize  int64           `yaml:"max_file_size"`
	AllowedMimes []string        `yaml:"allowed_mimes"`
}

// LocalUploadConf 本地上传配置
type LocalUploadConf struct {
	RootPath string `yaml:"root_path"`
	BaseURL  string `yaml:"base_url"`
}

// SamplingConfig 日志采样配置
type SamplingConfig struct {
	Initial    int `yaml:"initial"`
	Thereafter int `yaml:"thereafter"`
}

// LogFileConfig 日志文件配置
type LogFileConfig struct {
	Path       string `yaml:"path"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
}

// LogConfig 日志配置（修复层级问题）
type LogConfig struct {
	Level    string          `yaml:"level"`
	Format   string          `yaml:"format"`
	Output   string          `yaml:"output"`
	File     LogFileConfig   `yaml:"file"`
	Sampling *SamplingConfig `yaml:"sampling"`
}

// TraceConfig 链路追踪配置（从顶层合并到 log）
type TraceConfig struct {
	Enabled  bool    `yaml:"enabled"`
	Sampler  string  `yaml:"sampler"`
	Ratio    float64 `yaml:"ratio"`
	Exporter string  `yaml:"exporter"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
	Port    int    `yaml:"port"`
}

// CorsConfig CORS配置（修复命名）
type CorsConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// HTTPServerConfig HTTP服务器配置（合并重复的 server 块）
type HTTPServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

// GRPCServerConfig gRPC服务器配置（修复命名）
type GRPCServerConfig struct {
	Host           string          `yaml:"host"`
	Port           int             `yaml:"port"`
	MaxRecvMsgSize int             `yaml:"max_recv_msg_size"`
	MaxSendMsgSize int             `yaml:"max_send_msg_size"`
	KeepAlive      KeepAliveConfig `yaml:"keepalive"`
}

// KeepAliveConfig gRPC keepalive配置
type KeepAliveConfig struct {
	EnforcementPolicy KeepAlivePolicy `yaml:"enforcement_policy"`
	ServerParams      KeepAliveServer `yaml:"server_params"`
}

// KeepAlivePolicy gRPC最小时间策略
type KeepAlivePolicy struct {
	MinTime time.Duration `yaml:"min_time"`
}

// KeepAliveServer gRPC服务端keepalive参数
type KeepAliveServer struct {
	Time    time.Duration `yaml:"time"`
	Timeout time.Duration `yaml:"timeout"`
}

// ServerConfig 服务器配置（合并重复的 server/grpc/http 块）
type ServerConfig struct {
	Port         int              `yaml:"port"`
	Mode         string           `yaml:"mode"`
	ReadTimeout  time.Duration    `yaml:"read_timeout"`
	WriteTimeout time.Duration    `yaml:"write_timeout"`
	HTTP         HTTPServerConfig `yaml:"http"`
	GRPC         GRPCServerConfig `yaml:"grpc"`
}

// VideoConfig 视频配置
type VideoConfig struct {
	Enabled       bool     `yaml:"enabled"`
	MaxFileSize   int64    `yaml:"max_file_size"`
	StorageDriver string   `yaml:"storage_driver"`
	LocalRootPath string   `yaml:"local_root_path"`
	LocalBaseURL  string   `yaml:"local_base_url"`
	AllowedMIMEs  []string `yaml:"allowed_mimes"`
}

// MySQLBackupConfig MySQL备份配置
type MySQLBackupConfig struct {
	DumpPath string `yaml:"dump_path"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
}

// BackupConfig 数据备份配置
type BackupConfig struct {
	Enabled        bool              `yaml:"enabled"`
	StoragePath    string            `yaml:"storage_path"`
	MaxBackupCount int               `yaml:"max_backup_count"`
	MySQL          MySQLBackupConfig `yaml:"mysql"`
}

// ModerationConfig 内容审核配置
type ModerationConfig struct {
	Enabled              bool `yaml:"enabled"`
	AutoPublishThreshold int8 `yaml:"auto_publish_threshold"` // 低于此阈值自动发布
	RequireManualReview  int8 `yaml:"require_manual_review"`  // 达到此等级需人工审核
	BlockLevel           int8 `yaml:"block_level"`            // 达到此等级直接拦截
}

// Config 主配置结构（修复重复和层级问题）
type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"mysql"`
	Redis         RedisConfig         `yaml:"redis"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	JWT           JWTConfig           `yaml:"jwt"`
	Upload        UploadConfig        `yaml:"upload"`
	Log           LogConfig           `yaml:"log"`
	Trace         TraceConfig         `yaml:"trace"` // 从顶层合并到 log.sampling
	Metrics       MetricsConfig       `yaml:"metrics"`
	Cors          CorsConfig          `yaml:"cors"`
	Video         VideoConfig         `yaml:"video"`
	Backup        BackupConfig        `yaml:"backup"`
	Moderation    ModerationConfig    `yaml:"moderation"`
}

// GlobalConfig 全局配置实例
var GlobalConfig *Config

func init() {
	GlobalConfig = &Config{}
}
