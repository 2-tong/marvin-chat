package cache

import (
	"fmt"
	"sync"
	"time"
)

type TimeOutCache struct {
	cache map[string]cacheItem
	mutex sync.RWMutex
}

type cacheItem struct {
	value  interface{}
	expiry time.Time
}

func NewCache() *TimeOutCache {
	return &TimeOutCache{
		cache: make(map[string]cacheItem),
	}
}

func (c *TimeOutCache) Set(key string, value interface{}, exp time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[key] = cacheItem{
		value:  value,
		expiry: exp,
	}
	fmt.Printf("Set key '%s' with value '%v' expiring at %v\n", key, value, exp)
}

func (c *TimeOutCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.cache[key]
	if !found {
		return nil, false
	}

	if item.expiry.Before(time.Now()) {
		delete(c.cache, key)
		return nil, false
	}

	return item.value, true
}
