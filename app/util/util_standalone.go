// +build !appengine

package util

import (
	"log"
	"net/http"
	"time"
)

type logger struct{}

func Logger(r *http.Request) *logger {
	return &logger{}
}

func (l *logger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func HTTPClient(r *http.Request) *http.Client {
	return &http.Client{Timeout: time.Second * 30}
}
