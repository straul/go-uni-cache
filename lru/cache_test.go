package lru

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache_SetGet(t *testing.T) {
	cache := NewCache(2)
	ctx := context.Background()

	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.NoError(t, err)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value1", value)
}

func TestLRUCache_Capacity(t *testing.T) {
	cache := NewCache(2)
	ctx := context.Background()

	_ = cache.Set(ctx, "key1", "value1", 5*time.Second)
	_ = cache.Set(ctx, "key2", "value2", 5*time.Second)
	_ = cache.Set(ctx, "key3", "value3", 5*time.Second)

	// key1 应该被移除，因为容量是 2
	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)

	// key2 和 key3 应该还在
	value, found, err = cache.Get(ctx, "key2")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value2", value)

	value, found, err = cache.Get(ctx, "key3")
	assert.NoError(t, err)
	assert.True(t, found)
	assert.Equal(t, "value3", value)
}

func TestLRUCache_Expiration(t *testing.T) {
	cache := NewCache(2)
	ctx := context.Background()

	_ = cache.Set(ctx, "key1", "value1", 1*time.Second)
	time.Sleep(2 * time.Second)

	value, found, err := cache.Get(ctx, "key1")
	assert.NoError(t, err)
	assert.False(t, found)
	assert.Nil(t, value)
}

func TestLRUCache_ContextCancel(t *testing.T) {
	cache := NewCache(2)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := cache.Set(ctx, "key1", "value1", 5*time.Second)
	assert.Error(t, err)
}
