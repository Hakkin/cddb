// +build !appengine

package abstract

import (
	"net/http"
	"time"
	
	"golang.org/x/net/context"
)

func GetContext(r *http.Request) context.Context {
	return r.Context()
}

func GetClient(ctx context.Context) *http.Client {
	client := &http.Client{}
	client.Timeout = time.Second * 10
	return client
}
