package cache

import (
	"context"
	"sync"
	"time"
)

type MultiLevelCache struct {
	levels []Cache
	mu     sync.RWMutex
}

func NewMultiLevelCache(level ...Cache) *MultiLevelCache {
	return &MultiLevelCache{
		levels: level,
	}
}

func (c *MultiLevelCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, level := range c.levels {
		value, found, err := level.Get(ctx, key)
		if err != nil {
			return nil, false, err
		}
		if found {
			return value, true, nil
		}
	}
	return nil, false, nil
}

func (c *MultiLevelCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, level := range c.levels {
		if err := level.Set(ctx, key, value, ttl); err != nil {
			return err
		}
	}
	return nil
}

func (c *MultiLevelCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, level := range c.levels {
		if err := level.Delete(ctx, key); err != nil {
			return err
		}
	}
	return nil
}
