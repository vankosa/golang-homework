package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	Cache // Remove me after realization.

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	element, ok := c.items[key]
	// if element doesn't exist
	if !ok {
		newItem := cacheItem{
			key:   key,
			value: value,
		}
		// add to queue
		listItem := c.queue.PushFront(newItem)
		if c.queue.Len() > c.capacity {
			// get last element
			lastElement := c.queue.Back()
			// remove last element from queue
			c.queue.Remove(lastElement)
			// remove last element from items
			lastItem := lastElement.Value.(cacheItem)
			delete(c.items, lastItem.key)
		}
		// add to items
		c.items[key] = listItem

		return false
	}
	// replace value
	item := element.Value.(cacheItem)
	item.value = value
	element.Value = item
	// move element to front
	c.queue.MoveToFront(element)
	// update value in items
	c.items[item.key] = element
	return true
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	element, ok := c.items[key]
	// if element doesn't exist
	if !ok {
		return nil, false
	}
	// move element to front
	c.queue.MoveToFront(element)

	item := element.Value.(cacheItem)
	return item.value, true
}

func (c *lruCache) Clear() {
	// clear cache
	c.queue.MoveToFront(nil)
}
