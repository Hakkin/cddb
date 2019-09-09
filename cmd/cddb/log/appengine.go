package log

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
var serviceID = os.Getenv("GAE_SERVICE")
var versionID = os.Getenv("GAE_VERSION")

func init() {
	if os.Getenv("USING_APPENGINE") != "" {
		client, err := logging.NewClient(
			context.Background(),
			os.Getenv("GOOGLE_CLOUD_PROJECT"),
			option.WithCredentialsFile("appengine.json"),
		)
		if err != nil {
			panic(err)
		}

		logger := client.Logger(
			"stderr",
			logging.CommonResource(&monitoredres.MonitoredResource{
				Type: "gae_app",
				Labels: map[string]string{
					"project_id": projectID,
					"module_id":  serviceID,
					"version_id": versionID,
				},
			}),
		)

		Log = &AppEngineLogger{
			client: client,
			logger: logger,
		}
	}
}

type AppEngineLogger struct {
	client *logging.Client
	logger *logging.Logger
	r      *http.Request
}

func getTrace(r *http.Request) string {
	if r == nil {
		return ""
	}

	traceID := strings.SplitN(r.Header.Get("X-Cloud-Trace-Context"), "/", 2)[0]
	return fmt.Sprintf("projects/%s/traces/%s", projectID, traceID)
}

func (ael *AppEngineLogger) Infof(format string, args ...interface{}) {
	ael.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf(format, args...),
		Severity: logging.Info,
		Trace:    getTrace(ael.r),
	})
}

func (ael *AppEngineLogger) Errorf(format string, args ...interface{}) {
	ael.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf(format, args...),
		Severity: logging.Error,
		Trace:    getTrace(ael.r),
	})
}

func (ael *AppEngineLogger) WithRequest(r *http.Request) Logger {
	return &AppEngineLogger{
		client: ael.client,
		logger: ael.logger,
		r:      r,
	}
}
