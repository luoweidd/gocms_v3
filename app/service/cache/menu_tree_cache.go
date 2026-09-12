// Package cache provides Redis caching for menu tree and user sessions.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gocms_v3/app/db"

	"github.com/redis/go-redis/v9"
)

// MenuTreeItem represents a single menu item in the tree.
type MenuTreeItem struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	ParentID   *int             `json:"parent_id,omitempty"`
	Path       string           `json:"path"`
	Icon       string           `json:"icon"`
	Sort       int              `json:"sort"`
	Children   []*MenuTreeItem  `json:"children,omitempty"`
	Permission []PermissionInfo `json:"permission,omitempty"`
}

// PermissionInfo represents permission information for a menu item.
type PermissionInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Type string `json:"type"` // button, menu, dir
}

// MenuTreeCache provides caching for menu trees.
type MenuTreeCache struct {
	rdb      *redis.Client
	cacheKey string
	ttl      time.Duration
}

// NewMenuTreeCache creates a new menu tree cache instance.
func NewMenuTreeCache() *MenuTreeCache {
	return &MenuTreeCache{
		rdb:      db.GetRDB(),
		cacheKey: "cms:menu:tree",
		ttl:      24 * time.Hour,
	}
}

// GetCachedMenuTree retrieves the cached menu tree.
func (c *MenuTreeCache) GetCachedMenuTree(ctx context.Context) ([]*MenuTreeItem, error) {
	data, err := c.rdb.Get(ctx, c.cacheKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取菜单树缓存失败: %w", err)
	}

	var tree []*MenuTreeItem
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, fmt.Errorf("解析菜单树缓存失败: %w", err)
	}

	return tree, nil
}

// SetCachedMenuTree sets the menu tree cache.
func (c *MenuTreeCache) SetCachedMenuTree(ctx context.Context, tree []*MenuTreeItem) error {
	data, err := json.Marshal(tree)
	if err != nil {
		return fmt.Errorf("序列化菜单树失败: %w", err)
	}

	return c.rdb.SetEx(ctx, c.cacheKey, data, c.ttl).Err()
}

// InvalidateMenuTree invalidates the menu tree cache.
func (c *MenuTreeCache) InvalidateMenuTree(ctx context.Context) error {
	return c.rdb.Del(ctx, c.cacheKey).Err()
}

// InvalidateAllMenuKeys invalidates all menu-related cache keys.
func (c *MenuTreeCache) InvalidateAllMenuKeys(ctx context.Context) error {
	keys := []string{
		c.cacheKey,
		"cms:menu:tree:admin",
		"cms:menu:tree:frontend",
	}

	pipeline := c.rdb.Pipeline()
	for _, key := range keys {
		pipeline.Del(ctx, key)
	}
	_, err := pipeline.Exec(ctx)
	return err
}

// ============================================
// User Session Cache
// ============================================

// SessionData represents user session data.
type SessionData struct {
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	Roles        []string  `json:"roles"`
	LastActivity time.Time `json:"last_activity"`
	TokenHash    string    `json:"token_hash"`
}

// SessionCache provides caching for user sessions.
type SessionCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewSessionCache creates a new session cache instance.
func NewSessionCache() *SessionCache {
	return &SessionCache{
		rdb: db.GetRDB(),
		ttl: 2 * time.Hour,
	}
}

// GetUserSession retrieves user session from cache.
func (c *SessionCache) GetUserSession(ctx context.Context, sessionKey string) (*SessionData, error) {
	data, err := c.rdb.Get(ctx, sessionKey).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取会话缓存失败: %w", err)
	}

	var session SessionData
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("解析会话数据失败: %w", err)
	}

	return &session, nil
}

// SetUserSession stores user session in cache.
func (c *SessionCache) SetUserSession(ctx context.Context, sessionKey string, session *SessionData) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("序列化会话数据失败: %w", err)
	}

	return c.rdb.SetEx(ctx, sessionKey, data, c.ttl).Err()
}

// InvalidateUserSession invalidates user session from cache.
func (c *SessionCache) InvalidateUserSession(ctx context.Context, sessionKey string) error {
	return c.rdb.Del(ctx, sessionKey).Err()
}

// AddUserSession adds a session for a user.
func (c *SessionCache) AddUserSession(ctx context.Context, userID uint, sessionKey string) error {
	userKey := fmt.Sprintf("cms:session:user:%d", userID)
	return c.rdb.HSet(ctx, userKey, sessionKey, 1).Err()
}

// RemoveUserSession removes a session for a user.
func (c *SessionCache) RemoveUserSession(ctx context.Context, userID uint, sessionKey string) error {
	userKey := fmt.Sprintf("cms:session:user:%d", userID)
	return c.rdb.HDel(ctx, userKey, sessionKey).Err()
}

// ============================================
// Popular Menu Cache
// ============================================

// PopularMenuCache provides caching for popular/popular menu tree and user sessions.
type PopularMenuCache struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewPopularMenuCache creates a new popular menu cache instance.
func NewPopularMenuCache() *PopularMenuCache {
	return &PopularMenuCache{
		rdb: db.GetRDB(),
		ttl: 24 * time.Hour,
	}
}

// GetPopularMenu retrieves the cached popular menu tree.
func (c *PopularMenuCache) GetPopularMenu(ctx context.Context) ([]*MenuTreeItem, error) {
	data, err := c.rdb.Get(ctx, "cms:popular:menu:tree").Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("获取热门菜单缓存失败: %w", err)
	}

	var menu []*MenuTreeItem
	if err := json.Unmarshal(data, &menu); err != nil {
		return nil, fmt.Errorf("解析热门菜单数据失败: %w", err)
	}

	return menu, nil
}

// SetPopularMenu caches the popular menu tree.
func (c *PopularMenuCache) SetPopularMenu(ctx context.Context, menu []*MenuTreeItem) error {
	data, err := json.Marshal(menu)
	if err != nil {
		return fmt.Errorf("序列化热门菜单失败: %w", err)
	}

	return c.rdb.SetEx(ctx, "cms:popular:menu:tree", data, c.ttl).Err()
}

// InvalidatePopularMenu invalidates the popular menu cache.
func (c *PopularMenuCache) InvalidatePopularMenu(ctx context.Context) error {
	return c.rdb.Del(ctx, "cms:popular:menu:tree").Err()
}

// ============================================
// General Cache Helper Functions
// ============================================

// SetString sets a string value in Redis with TTL.
func SetString(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return db.GetRDB().SetEx(ctx, key, value, ttl).Err()
}

// GetString gets a string value from Redis.
func GetString(ctx context.Context, key string) (string, error) {
	return db.GetRDB().Get(ctx, key).Result()
}

// SetJSON sets a JSON value in Redis with TTL.
func SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化失败: %w", err)
	}
	return db.GetRDB().SetEx(ctx, key, data, ttl).Err()
}

// GetJSON gets a JSON value from Redis and unmarshals it.
func GetJSON(ctx context.Context, key string, target interface{}) error {
	data, err := db.GetRDB().Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return fmt.Errorf("获取缓存失败: %w", err)
	}
	return json.Unmarshal([]byte(data), target)
}

// InvalidateKey invalidates a cache key.
func InvalidateKey(ctx context.Context, key string) error {
	return db.GetRDB().Del(ctx, key).Err()
}

// InvalidatePattern invalidates all keys matching a pattern.
func InvalidatePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	pipeline := db.GetRDB().Pipeline()

	for {
		keys, newCursor, err := db.GetRDB().Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("扫描键失败: %w", err)
		}

		if len(keys) > 0 {
			for _, key := range keys {
				pipeline.Del(ctx, key)
			}
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	_, err := pipeline.Exec(ctx)
	return err
}
