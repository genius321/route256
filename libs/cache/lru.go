package cache

import (
	"container/list"
	"sync"
)

type Item struct {
	Key   string
	Value any
}

type LRU struct {
	capacity int
	queue    *list.List
	mutex    *sync.RWMutex
	items    map[string]*list.Element
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		capacity: capacity,
		queue:    list.New(),
		mutex:    new(sync.RWMutex),
		items:    make(map[string]*list.Element, capacity),
	}
}

func (c *LRU) Add(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if element, exist := c.items[key]; exist {
		c.queue.MoveToBack(element)
		element.Value.(*Item).Value = value
		return
	}

	if c.queue.Len() == c.capacity {
		element := c.queue.Front()
		c.deleteItem(element)
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	element := c.queue.PushBack(item)
	c.items[item.Key] = element
}

func (c *LRU) Get(key string) any {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exist := c.items[key]
	if !exist {
		return nil
	}

	c.queue.MoveToBack(element)
	return element.Value.(*Item).Value
}

func (c *LRU) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.items)
}

func (c *LRU) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if element, exist := c.items[key]; exist {
		c.deleteItem(element)
		return true
	}
	return false
}

func (c *LRU) deleteItem(element *list.Element) {
	item := c.queue.Remove(element).(*Item)
	delete(c.items, item.Key)
}
