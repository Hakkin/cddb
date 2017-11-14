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

func (c *RCache) Set(ids ...string) ([]string, error) {
	idCount := len(ids)
	newIDs := make([]string, idCount)
	idEntitys := make([]*key, idCount)
	idKeys := make([]*datastore.Key, idCount)
	var err error
	for i, id := range ids {
		newIDs[i], err = translate(id)
		if err != nil {
			return nil, err
		}
		idEntitys[i] = &key{Value: id}

		idKeys[i] = datastore.NewKey(appengine.NewContext(c.r), "ID", newIDs[i], 0, nil)
	}
	_, err = datastore.PutMulti(appengine.NewContext(c.r), idKeys, idEntitys)
	if err != nil {
		return nil, err
	}

	return newIDs, nil
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
