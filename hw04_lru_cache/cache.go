package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		return true
	}
	if c.queue.Len() == c.capacity {

	}
	c.items[key] = c.queue.PushFront(value)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	return nil, false
}

func (c *lruCache) Clear() {

}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
