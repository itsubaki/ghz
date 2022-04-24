package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/errorreporting"
)

const (
	DEFAULT   = "Default"
	DEBUG     = "Debug"
	INFO      = "Info"
	NOTICE    = "Notice"
	WARNING   = "Warning"
	ERROR     = "Error"
	CRITICAL  = "Critical"
	ALERT     = "Alert"
	EMERGENCY = "Emergency"
)

type LogEntry struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
	Trace    string    `json:"logging.googleapis.com/trace"`
	SpanID   string    `json:"logging.googleapis.com/spanId,omitempty"`
}

type Logger struct {
	projectID string
	trace     string
	errC      *errorreporting.Client
	req       *http.Request
}

func New(projectID, traceID string) *Logger {
	trace := ""
	if len(traceID) > 0 {
		trace = fmt.Sprintf("projects/%v/traces/%v", projectID, traceID[0])
	}

	return &Logger{
		projectID: projectID,
		trace:     trace,
	}
}

func (l *Logger) LogWith(spanID, severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.trace,
		SpanID:   spanID,
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
	}); err != nil {
		log.Printf("encode log entry: %v", err)
	}
}

func (l *Logger) DebugWith(spanID, format string, a ...interface{}) {
	l.LogWith(spanID, DEBUG, format, a...)
}

func (l *Logger) ErrorWith(spanID, format string, a ...interface{}) {
	l.LogWith(spanID, ERROR, format, a...)
}

func (l *Logger) Log(severity, format string, a ...interface{}) {
	l.LogWith("", severity, format, a...)
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.Log(DEBUG, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log(INFO, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Error(format, a...)
}

func (l *Logger) NewReport(ctx context.Context, req *http.Request) *Logger {
	c, err := errorreporting.NewClient(ctx, l.projectID, errorreporting.Config{})
	if err != nil {
		l.Error("new error report client: %v", err)
		return l
	}

	l.errC = c
	l.req = req
	return l
}

func (l *Logger) ReportWith(spanID, severity, format string, a ...interface{}) {
	l.LogWith(spanID, severity, format, a...)
	if l.errC == nil {
		return
	}

	for _, aa := range a {
		switch err := aa.(type) {
		case error:
			l.errC.Report(errorreporting.Entry{
				Error: err,
				Req:   l.req,
			})
		}
	}
}

func (l *Logger) ErrorReportWith(spanID, format string, a ...interface{}) {
	l.ReportWith(spanID, ERROR, format, a...)
}

func (l *Logger) Report(severity, format string, a ...interface{}) {
	l.ReportWith("", severity, format, a...)
}

func (l *Logger) ErrorReport(format string, a ...interface{}) {
	l.Report(ERROR, format, a...)
}
