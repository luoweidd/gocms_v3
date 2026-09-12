package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheManager 缓存管理器
type CacheManager struct {
	client *redis.Client
	ctx    context.Context
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(ctx context.Context, addr, password string, db int) *CacheManager {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &CacheManager{
		client: client,
		ctx:    ctx,
	}
}

// Get 获取缓存
func (c *CacheManager) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

// Set 设置缓存
func (c *CacheManager) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return c.client.Set(c.ctx, key, data, ttl).Err()
}

// Delete 删除缓存
func (c *CacheManager) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists 检查键是否存在
func (c *CacheManager) Exists(key string) bool {
	result := c.client.Exists(c.ctx, key).Val()
	return result > 0
}

// SetWithTTL 设置带过期时间的缓存
func (c *CacheManager) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return c.client.Set(c.ctx, key, data, ttl).Err()
}

// GetWithExpire 获取缓存并返回剩余过期时间
func (c *CacheManager) GetWithExpire(key string) (string, time.Duration, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", 0, err
	}
	ttl, err := c.client.TTL(c.ctx, key).Result()
	return val, ttl, err
}

// SetJSON 设置JSON缓存
func (c *CacheManager) SetJSON(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return c.client.Set(c.ctx, key, data, ttl).Err()
}

// GetJSON 获取JSON缓存并反序列化
func (c *CacheManager) GetJSON(key string, target interface{}) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

// Increment 增加整数
func (c *CacheManager) Increment(key string) error {
	return c.client.Incr(c.ctx, key).Err()
}

// Decrement 减少整数
func (c *CacheManager) Decrement(key string) error {
	return c.client.Decr(c.ctx, key).Err()
}

// SetNX 设置键（仅当键不存在时）
func (c *CacheManager) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("序列化失败: %w", err)
	}
	result, err := c.client.SetNX(c.ctx, key, data, ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// Close 关闭缓存连接
func (c *CacheManager) Close() error {
	return c.client.Close()
}

// DistributedLock 分布式锁
type DistributedLock struct {
	client *redis.Client
	ctx    context.Context
	key    string
	token  string
	ttl    time.Duration
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(client *redis.Client, ctx context.Context, key string, ttl time.Duration) *DistributedLock {
	return &DistributedLock{
		client: client,
		ctx:    ctx,
		key:    key,
		token:  fmt.Sprintf("%d", time.Now().UnixNano()),
		ttl:    ttl,
	}
}

// Acquire 尝试获取锁
func (l *DistributedLock) Acquire() (bool, error) {
	result, err := l.client.SetNX(l.ctx, l.key, l.token, l.ttl).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

// Release 释放锁
func (l *DistributedLock) Release() error {
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	_, err := l.client.Eval(l.ctx, script, []string{l.key}, l.token).Result()
	return err
}

// Refresh 刷新锁过期时间
func (l *DistributedLock) Refresh() error {
	return l.client.Expire(l.ctx, l.key, l.ttl).Err()
}
