package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/straul/go-uni-cache/cache"
)

type Cache struct {
	client *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil // 键不存在，返回 nil 值和 false
		}
		return nil, false, errors.Wrap(err, "Redis Get operation failed")
	}
	return val, true, nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return errors.Wrap(err, "Redis Set operation failed")
	}
	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return errors.Wrap(err, "Redis Delete operation failed")
	}
	return nil
}

var _ cache.Cache = (*Cache)(nil)
