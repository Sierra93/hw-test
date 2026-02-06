package hw04lrucache

import "sync"

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
	mu       sync.Mutex // ДОБАВЛЕНО: мьютекс для защиты данных
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()         // Блокируем доступ
	defer c.mu.Unlock() // Разблокируем при выходе

	if _, ok := c.items[key]; !ok {
		if c.queue.Len() == c.capacity {
			// Важно: Key должен быть в ListItem, чтобы мы знали, что удалять из map
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
	c.mu.Lock()         // Блокируем доступ
	defer c.mu.Unlock() // Разблокируем при выходе

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.queue.MoveToFront(item)
	return item.Value, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()         // Блокируем доступ
	defer c.mu.Unlock() // Разблокируем при выходе

	c.items = make(map[Key]*ListItem, c.capacity)
	c.queue = NewList()
}
