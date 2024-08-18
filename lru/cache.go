package lru

import (
	"context"
	"github.com/pkg/errors"
	"sync"
	"time"
)

// node 双向链表节点
type node struct {
	key        string
	value      interface{}
	expiration int64
	prev       *node
	next       *node
}

// Cache LRU 缓存结构体
type Cache struct {
	capacity int
	items    map[string]*node
	head     *node
	tail     *node
	mu       sync.Mutex
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*node),
	}
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-ctx.Done():
		return nil, false, errors.New("context canceled in Get operation")
	default:
		if n, found := c.items[key]; found {
			if n.expiration > 0 && n.expiration < time.Now().UnixNano() {
				// 如果已过期，需要删除
				c.remove(n)
				delete(c.items, key)
				return nil, false, nil
			}
			// 移动到链表头部
			c.moveToFront(n)
			return n.value, true, nil
		}
		return nil, false, nil
	}
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.New("context canceled in Set operation")
	default:
		if n, found := c.items[key]; found {
			// 更新现有节点
			n.value = value
			n.expiration = time.Now().Add(ttl).UnixNano()
			c.moveToFront(n)
		} else {
			// 添加新节点
			n = &node{
				key:        key,
				value:      value,
				expiration: time.Now().Add(ttl).UnixNano(),
			}
			c.items[key] = n
			c.addToFront(n)
			if len(c.items) > c.capacity {
				// 超出容量，移除最久未使用的节点
				delete(c.items, c.tail.key)
				c.remove(c.tail)
			}
		}
		return nil
	}
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-ctx.Done():
		return errors.New("context canceled in Delete operation")
	default:
		if n, found := c.items[key]; found {
			c.remove(n)
			delete(c.items, key)
		}
		return nil
	}
}

// ------------------- Node func -------------------

// 将节点移到链表头部
func (c *Cache) moveToFront(n *node) {
	if n == c.head {
		return
	}
	c.remove(n)
	c.addToFront(n)
}

// 从链表中移除节点
func (c *Cache) remove(n *node) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		c.head = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	} else {
		c.tail = n.prev
	}
}

// 讲节点添加到链表头部
func (c *Cache) addToFront(n *node) {
	n.next = c.head
	n.prev = nil
	if c.head != nil {
		c.head.prev = n
	}
	c.head = n
	if c.tail == nil {
		c.tail = n
	}
}
