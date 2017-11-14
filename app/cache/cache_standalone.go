// +build !appengine

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

func (c *RCache) Set(ids ...string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	newIDs := make([]string, len(ids))
	var err error
	for i, id := range ids {
		newIDs[i], err = translate(id)
		if err != nil {
			return nil, err
		}
		c.ids[newIDs[i]] = id
	}

	return newIDs, nil
}

func (c *RCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	oldID, ok := c.ids[id]
	return oldID, ok
}
