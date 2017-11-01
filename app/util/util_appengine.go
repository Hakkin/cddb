// +build appengine

package util

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

type logger struct {
	ctx context.Context
}

func Logger(r *http.Request) *logger {
	return &logger{ctx: appengine.NewContext(r)}
}

func (l *logger) Infof(format string, args ...interface{}) {
	log.Infof(l.ctx, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	log.Errorf(l.ctx, format, args...)
}

func HTTPClient(r *http.Request) *http.Client {
	ctx := appengine.NewContext(r)
	ctx, _ = context.WithTimeout(ctx, time.Second * 10)
	httpClient := urlfetch.Client(ctx)
	httpClient.Timeout = time.Second * 10
	return httpClient
}
