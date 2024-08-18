package memory

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/straul/go-uni-cache/cache"
)

type item struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	data map[string]item
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]item),
	}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, errors.New("context canceled in Get operation")
	default:
		c.mu.RLock()
		defer c.mu.RUnlock()
		it, ok := c.data[key]
		if !ok || (it.expiration > 0 && it.expiration < time.Now().UnixNano()) {
			return nil, false, nil // 键不存在，返回 nil 值和 false
		}
		return it.value, true, nil
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled in Set operation")
	default:
		c.mu.Lock()
		defer c.mu.Unlock()
		expiration := time.Now().Add(ttl).UnixNano()
		c.data[key] = item{
			value:      value,
			expiration: expiration,
		}
		return nil
	}
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return errors.New("context canceled in Delete operation")
	default:
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.data, key)
		return nil
	}
}

var _ cache.Cache = (*Cache)(nil)
