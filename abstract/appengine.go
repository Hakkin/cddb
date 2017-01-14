// +build appengine

package abstract

import (
	"net/http"
	"time"
	
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func GetContext(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

func GetClient(ctx context.Context) *http.Client {
	client := urlfetch.Client(ctx)
	client.Timeout = time.Second * 10
	return client
}
