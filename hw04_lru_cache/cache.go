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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if v, ok := c.items[key]; ok {
		v.Value = value
		c.queue.MoveToFront(v)
		return true
	}
	newItem := c.queue.PushFront(value)
	c.items[key] = newItem
	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		c.queue.Remove(back)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	v, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if v.Next == nil && v.Prev == nil {
		return nil, false
	}
	c.queue.MoveToFront(v)
	return v.Value, true
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
