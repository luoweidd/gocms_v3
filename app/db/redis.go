package db

import (
	"context"
	"fmt"
	"log"

	"gocms_v3/app/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis 初始化 Redis 连接池
func InitRedis(cfg config.RedisConfig) error {
	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	ctx := context.Background()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Printf("[REDIS-ERROR] Redis连接失败 [%s:%d]: %v", cfg.Host, cfg.Port, err)
		return fmt.Errorf("Redis连接失败 [%s:%d]: %w", cfg.Host, cfg.Port, err)
	}

	log.Printf("[REDIS] Redis连接成功: %s:%d (DB: %d, 最大连接数: %d, 空闲连接数: %d)",
		cfg.Host, cfg.Port, cfg.DB, cfg.PoolSize, cfg.MinIdleConns)

	return nil
}

// GetRDB 获取 Redis 客户端
func GetRDB() *redis.Client {
	if RDB == nil {
		panic("Redis未初始化")
	}
	return RDB
}
