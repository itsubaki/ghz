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
	"go.opentelemetry.io/otel/trace"
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

func (l *Logger) Log(severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.trace,
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
	}); err != nil {
		log.Printf("encode log entry: %v", err)
	}
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.Log(DEBUG, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log(INFO, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Log(ERROR, format, a...)
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

func (l *Logger) LogReport(severity, format string, a ...interface{}) {
	l.Log(severity, format, a...)
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

func (l *Logger) ErrorReport(format string, a ...interface{}) {
	l.LogReport(ERROR, format, a...)
}

func (l *Logger) SpanOf(spanID string) *SpanLogEntry {
	return &SpanLogEntry{
		Trace:  l.trace,
		SpanID: spanID,
	}
}

func (l *Logger) Span(span trace.Span) *SpanLogEntry {
	return &SpanLogEntry{
		Trace:  l.trace,
		SpanID: span.SpanContext().SpanID().String(),
	}
}

type SpanLogEntry struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
	Trace    string    `json:"logging.googleapis.com/trace"`
	SpanID   string    `json:"logging.googleapis.com/spanId,omitempty"`
}

func (e *SpanLogEntry) Log(severity, format string, a ...interface{}) {
	e.Severity = severity
	e.Message = fmt.Sprintf(format, a...)
	e.Time = time.Now()

	if err := json.NewEncoder(os.Stdout).Encode(e); err != nil {
		log.Printf("encode log entry: %v", err)
	}
}

func (e *SpanLogEntry) Debug(format string, a ...interface{}) {
	e.Log(DEBUG, format, a...)
}
