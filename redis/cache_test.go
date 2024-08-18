package redis

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "39.106.26.64:17956",
		Password: "bBc$3II&crQA",
	})
}

func TestRedisCache_SetGet(t *testing.T) {
	client := newTestRedisClient()
	cache := NewCache(client)
	ctx := context.Background()

	// 测试设置和获取缓存
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.NoError(t, err)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value1", value)
}

func TestRedisCache_GetNotFound(t *testing.T) {
	client := newTestRedisClient()
	cache := NewCache(client)
	ctx := context.Background()

	// 测试获取不存在的键
	value, found, err := cache.Get(ctx, "keyNotExist")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestRedisCache_Delete(t *testing.T) {
	client := newTestRedisClient()
	cache := NewCache(client)
	ctx := context.Background()

	// 测试删除缓存
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.NoError(t, err)

	err = cache.Delete(ctx, "key1")
	assert.NoError(t, err)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestRedisCache_Expiration(t *testing.T) {
	client := newTestRedisClient()
	cache := NewCache(client)
	ctx := context.Background()

	// 测试键的过期
	err := cache.Set(ctx, "key1", "value1", 1*time.Second)
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestRedisCache_ContextCancel(t *testing.T) {
	client := newTestRedisClient()
	cache := NewCache(client)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// 测试取消的上下文
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.Error(t, err)
}
