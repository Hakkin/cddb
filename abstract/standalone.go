// +build !appengine

package abstract

import (
	"net/http"
	"time"
)

var Request *http.Request

func GetClient() *http.Client {
	client := &http.Client{}
	client.Timeout = time.Second * 10
	return client
}
