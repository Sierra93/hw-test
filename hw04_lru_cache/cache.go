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
	if _, ok := c.items[key]; !ok {
		if c.queue.Len() == c.capacity {
			delete(c.items, c.queue.Back().Key)
			c.queue.Remove(c.queue.Back())
		}
		newItem := c.queue.PushFront(value)
		newItem.Key = key
		c.items[key] = newItem
		return false
	}

	c.items[key].Value = value
	c.queue.MoveToFront(c.items[key])
	return true
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := c.items[key]; !ok {
		return nil, false
	}

	c.queue.MoveToFront(c.items[key])
	return c.items[key].Value, true
}

func (c *lruCache) Clear() {
	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
