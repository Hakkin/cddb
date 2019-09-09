package log

import (
	"net/http"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	WithRequest(r *http.Request) Logger
}

var Log Logger = &StandardLogger{}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func WithRequest(r *http.Request) Logger {
	return Log.WithRequest(r)
}
