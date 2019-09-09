package cache

import "net/http"

type Cacher interface {
	New(r *http.Request) RCacher
}

type RCacher interface {
	Set(ids ...string) ([]string, error)
	Get(id string) (string, bool)
}

var Cache Cacher = &StandardCache{ids: make(map[string]string)}

func New(r *http.Request) RCacher {
	return Cache.New(r)
}
