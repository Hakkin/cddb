package cache

import (
	"net/http"
	"sync"
)

type StandardCache struct {
	mu  sync.RWMutex
	ids map[string]string
}

type StandardRCache struct {
	*StandardCache
	r *http.Request
}

func (c *StandardCache) New(r *http.Request) RCacher {
	return &StandardRCache{StandardCache: c, r: r}
}

func (c *StandardRCache) Set(ids ...string) ([]string, error) {
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

func (c *StandardRCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	oldID, ok := c.ids[id]
	return oldID, ok
}
