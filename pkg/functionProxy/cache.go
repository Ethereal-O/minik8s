package functionProxy

import "sync"

var cache FifoCache

type FifoCache struct {
	lock       sync.Mutex
	items      map[string]string
	capability int
	queue      []string
	point      int
}

func (c *FifoCache) Init(cap int) {
	c.items = make(map[string]string)
	c.capability = cap
	c.queue = make([]string, cap)
	c.point = 0
}

func (c *FifoCache) PutCache(key string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.queue[c.point] != "" {
		delete(c.items, c.queue[c.point])
	}
	c.items[key] = value
	c.queue[c.point] = key
	c.point = (c.point + 1) % c.capability
}

func (c *FifoCache) GetCache(key string) string {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.items[key]
	if ok {
		return value
	}
	return ""
}
