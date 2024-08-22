// cache/multi_level_cache_test.go
package cache

import (
	"context"
	"sync"
	"testing"
	"time"
)

// MockMemoryCache simulates an in-memory cache for testing purposes.
type MockMemoryCache struct {
	store map[string]interface{}
	ttl   map[string]time.Time
	mu    sync.RWMutex
}

func NewMockMemoryCache() *MockMemoryCache {
	return &MockMemoryCache{
		store: make(map[string]interface{}),
		ttl:   make(map[string]time.Time),
	}
}

func (c *MockMemoryCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.store[key]
	if !found {
		return nil, false, nil
	}
	if time.Now().After(c.ttl[key]) {
		return nil, false, nil
	}
	return value, true, nil
}

func (c *MockMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = value
	c.ttl[key] = time.Now().Add(ttl)
	return nil
}

func (c *MockMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.store, key)
	delete(c.ttl, key)
	return nil
}

// MockRedisCache simulates a Redis cache using MockMemoryCache.
type MockRedisCache struct {
	MockMemoryCache
}

func NewMockRedisCache() *MockRedisCache {
	return &MockRedisCache{
		MockMemoryCache: *NewMockMemoryCache(),
	}
}

// TestMultiLevelCache tests the MultiLevelCache with memory and Redis caches.
func TestMultiLevelCache(t *testing.T) {
	memCache := NewMockMemoryCache()
	redisCache := NewMockRedisCache() // Simulating Redis with MockRedisCache
	multiLevelCache := NewMultiLevelCache(memCache, redisCache)

	ctx := context.Background()

	// Test setting and getting values
	if err := multiLevelCache.Set(ctx, "key1", "value1", 10*time.Second); err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	value, found, err := multiLevelCache.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}
	if !found || value != "value1" {
		t.Fatalf("Expected value 'value1', got %v", value)
	}

	// Test deletion
	if err := multiLevelCache.Delete(ctx, "key1"); err != nil {
		t.Fatalf("Error deleting value: %v", err)
	}

	value, found, err = multiLevelCache.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}
	if found {
		t.Fatalf("Expected key 'key1' to be deleted, but found %v", value)
	}
}
