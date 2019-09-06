package cache

import (
	"context"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"

	"google.golang.org/api/option"
)

type AppEngineCache struct {
	client *datastore.Client
}

type AppEngineRCache struct {
	*AppEngineCache
	r *http.Request
}

type key struct {
	Value string
}

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func init() {
	if os.Getenv("USING_APPENGINE") != "" {
		client, err := datastore.NewClient(
			context.Background(),
			projectID,
			option.WithCredentialsFile("appengine.json"),
		)
		if err != nil {
			panic(err)
		}

		Cache = &AppEngineCache{client: client}
	}
}

func (c *AppEngineCache) New(r *http.Request) RCacher {
	return &AppEngineRCache{AppEngineCache: c, r: r}
}

func (c *AppEngineRCache) Set(ids ...string) ([]string, error) {
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

		idKeys[i] = datastore.NameKey("ID", newIDs[i], nil)
	}
	_, err = c.client.PutMulti(c.r.Context(), idKeys, idEntitys)
	if err != nil {
		return nil, err
	}

	return newIDs, nil
}

func (c *AppEngineRCache) Get(id string) (string, bool) {
	idKey := datastore.NameKey("ID", id, nil)

	var idEntity key
	err := c.client.Get(c.r.Context(), idKey, &idEntity)
	if err != nil {
		return "", false
	}

	return idEntity.Value, true
}
