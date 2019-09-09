package log

import (
	"log"
	"net/http"
)

type StandardLogger struct{}

func (sl *StandardLogger) Infof(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (sl *StandardLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func (sl *StandardLogger) WithRequest(r *http.Request) Logger {
	return sl
}
