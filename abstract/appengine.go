// +build appengine

package abstract

import (
	"net/http"
	"time"
	
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var Request *http.Request

func GetClient() *http.Client {
	context := appengine.NewContext(Request)
	client := urlfetch.Client(context)
	client.Timeout = time.Second * 10
	return client
}
