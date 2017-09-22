// +build appengine

package cache

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Cache struct{}

type RCache struct {
	*Cache
	r *http.Request
}

type key struct {
	Value string
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) New(r *http.Request) *RCache {
	return &RCache{Cache: c, r: r}
}

func (c *RCache) Set(id string) (string, error) {
	newID, err := translate(id)
	if err != nil {
		return "", err
	}
	idEntity := &key{Value: id}

	idKey := datastore.NewKey(appengine.NewContext(c.r), "ID", newID, 0, nil)

	_, err = datastore.Put(appengine.NewContext(c.r), idKey, idEntity)
	if err != nil {
		return "", err
	}

	return newID, nil
}

func (c *RCache) Get(id string) (string, bool) {
	idKey := datastore.NewKey(appengine.NewContext(c.r), "ID", id, 0, nil)

	var idEntity key
	err := datastore.Get(appengine.NewContext(c.r), idKey, &idEntity)
	if err != nil {
		return "", false
	}

	return idEntity.Value, true
}
