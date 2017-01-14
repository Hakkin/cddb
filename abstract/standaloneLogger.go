// +build !appengine

package abstract

import (
	"log"

	"golang.org/x/net/context"
)

func Infof(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[INFO] " + format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.Printf("[ERROR] " + format, args...)
}