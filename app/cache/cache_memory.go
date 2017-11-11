// +build !appengine,!sqlite,!bolt

package cache

import (
	"net/http"
	"sync"
)

type Cache struct {
	mu  sync.RWMutex
	ids map[string]string
}

type RCache struct {
	*Cache
	r *http.Request
}

func New() *Cache {
	return &Cache{ids: make(map[string]string)}
}

func (c *Cache) New(r *http.Request) *RCache {
	return &RCache{Cache: c, r: r}
}

func (c *RCache) Set(id string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newID, err := translate(id)
	if err != nil {
		return "", err
	}
	c.ids[newID] = id
	return newID, nil
}

func (c *RCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	oldID, ok := c.ids[id]
	return oldID, ok
}
