package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCache_SetGet(t *testing.T) {
	cache := NewCache()
	ctx := context.Background()

	// 测试设置和获取缓存
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.NoError(t, err)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value1", value)
}

func TestMemoryCache_GetNotFound(t *testing.T) {
	cache := NewCache()
	ctx := context.Background()

	// 测试获取不存在的键
	value, found, err := cache.Get(ctx, "keyNotExist")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestMemoryCache_Delete(t *testing.T) {
	cache := NewCache()
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

func TestMemoryCache_Expiration(t *testing.T) {
	cache := NewCache()
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

func TestMemoryCache_ContextCancel(t *testing.T) {
	cache := NewCache()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// 测试取消的上下文
	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.Error(t, err)
	assert.Equal(t, "context canceled in Set operation", err.Error())
}
