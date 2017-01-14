// +build appengine

package abstract

import (
	"google.golang.org/appengine/log"

	"golang.org/x/net/context"
)

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.Infof(ctx, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Errorf(ctx, format, args...)
}