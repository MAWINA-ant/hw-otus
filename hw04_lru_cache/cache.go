package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type pair struct {
	first  Key
	second interface{}
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	p := pair{key, value}
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		v.Value = p
		return true
	}
	if c.queue.Len() == c.capacity {
		back := c.queue.Back()
		backValue, ok := back.Value.(pair)
		if ok {
			delete(c.items, backValue.first)
		}
		c.queue.Remove(c.queue.Back())
	}
	c.items[key] = c.queue.PushFront(p)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		front := c.queue.Front()
		frontValue, ok := front.Value.(pair)
		if ok {
			return frontValue.second, true
		}
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
