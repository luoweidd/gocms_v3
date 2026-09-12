package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadConfig 从配置文件加载配置并返回全局实例
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	GlobalConfig = cfg

	// 初始化日志
	InitLogger()

	log.Printf("[CONFIG] 配置加载成功 - 文件: %s", configPath)
	log.Printf("[CONFIG] MySQL: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	log.Printf("[CONFIG] Redis: %s:%d", cfg.Redis.Host, cfg.Redis.Port)
	log.Printf("[CONFIG] Elasticsearch: %v", cfg.Elasticsearch.Addresses)
	log.Printf("[CONFIG] JWT Secret: %.6s...", cfg.JWT.Secret)
	log.Printf("[CONFIG] 服务器端口: %d", cfg.Server.Port)

	return cfg, nil
}
