package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	queue    List
	items    map[Key]*listItem
	capacity int
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		queue:    NewList(),
		items:    make(map[Key]*listItem),
		capacity: capacity,
	}
}

// Set saves item with key and value into items.
// It returns true if item with key was in items
// before setting, else false.
// If items already reached its capacity and we try
// to add new item into items, the oldest element
// will be removed to get a memory.
func (c *lruCache) Set(key Key, value interface{}) bool {
	if el, ok := c.items[key]; ok {
		item := el.Value.(*cacheItem)
		item.Value = value

		c.queue.MoveToFront(el)

		return true
	}

	c.items[key] = c.queue.PushFront(&cacheItem{Key: key, Value: value})

	if c.capacity < c.queue.Len() {
		last := c.queue.Back()
		c.queue.Remove(last)
		delete(c.items, last.Value.(*cacheItem).Key)
	}

	return false
}

// Get returns items value and true if item with
// key found in items, else nil and false.
func (c *lruCache) Get(key Key) (interface{}, bool) {
	if el, ok := c.items[key]; ok {
		c.queue.MoveToFront(el)

		return el.Value.(*cacheItem).Value, true
	}

	return nil, false
}

// Clear clears items.
func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*listItem)
}
