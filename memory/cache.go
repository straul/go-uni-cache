package memory

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/straul/go-uni-cache/cache"
)

type item struct {
	value      interface{} // 缓存值
	expiration int64       // 缓存 ttl
}

type Cache struct {
	data        map[string]item // 缓存本存
	mu          sync.RWMutex    // μ
	cleanupStop chan struct{}   // 停止自动清理过期缓存的信号 channel
}

func NewCache(cleanupInterval time.Duration) *Cache {
	c := &Cache{
		data:        make(map[string]item),
		cleanupStop: make(chan struct{}),
	}

	if cleanupInterval > 0 {
		// 如果设置了自动清理过期缓存的间隔时间，那么启动一个协程定时清理过期缓存
		go c.startCleanup(cleanupInterval)
	}

	return c
}

// Get 从缓存中获取键 key 对应的值
func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, false, errors.New("context canceled in Get operation")
	default:
		it, ok := c.data[key]
		if !ok || (it.expiration > 0 && it.expiration < time.Now().UnixNano()) {
			return nil, false, nil // 键不存在，返回 nil 值和 false
		}
		return it.value, true, nil
	}
}

// Set 设置键 key 对应的值 value，并设置过期时间 ttl
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.New("context canceled in Set operation")
	default:
		expiration := time.Now().Add(ttl).UnixNano()
		c.data[key] = item{
			value:      value,
			expiration: expiration,
		}
		return nil
	}
}

// Delete 删除键 key 对应的值
func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.New("context canceled in Delete operation")
	default:
		delete(c.data, key)
		return nil
	}
}

// startCleanup 启动一个协程定时清理过期缓存
func (c *Cache) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.cleanupStop:
			return
		}
	}
}

// cleanup 清理过期缓存
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, it := range c.data {
		if it.expiration > 0 && it.expiration < time.Now().UnixNano() {
			delete(c.data, key)
		}
	}
}

// StopCleanup 停止自动清理过期缓存
func (c *Cache) StopCleanup() {
	close(c.cleanupStop)
}

// 检查 Cache 是否实现了 cache.Cache 接口
var _ cache.Cache = (*Cache)(nil)
