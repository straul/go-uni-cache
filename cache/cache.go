package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (interface{}, bool, error)
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
