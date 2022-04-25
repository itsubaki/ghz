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

type LoggerFactory struct {
	projectID string
	errC      *errorreporting.Client
}

func Must(f *LoggerFactory, err error) *LoggerFactory {
	if err != nil {
		panic(err)
	}

	return f
}

func New(ctx context.Context, projectID string) (*LoggerFactory, error) {
	c, err := errorreporting.NewClient(ctx, projectID, errorreporting.Config{})
	if err != nil {
		return nil, fmt.Errorf("new error reporting client: %v", err)
	}

	return &LoggerFactory{
		projectID: projectID,
		errC:      c,
	}, nil
}

func MustNew(ctx context.Context, projectID string) *LoggerFactory {
	return Must(New(ctx, projectID))
}

func (f *LoggerFactory) New(traceID string, req *http.Request) *Logger {
	trace := ""
	if len(traceID) > 0 {
		trace = fmt.Sprintf("projects/%v/traces/%v", f.projectID, traceID)
	}

	return &Logger{
		errC:  f.errC,
		trace: trace,
		req:   req,
	}
}

type Logger struct {
	trace string
	errC  *errorreporting.Client
	req   *http.Request
}

type LogEntry struct {
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Time     time.Time `json:"time"`
	Trace    string    `json:"logging.googleapis.com/trace"`
}

func (l *Logger) Log(severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Severity: severity,
		Time:     time.Now(),
		Message:  fmt.Sprintf(format, a...),
		Trace:    l.trace,
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

func (l *Logger) LogReport(severity, format string, a ...interface{}) {
	// logging
	l.Log(severity, format, a...)

	// error reporting
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
	SpanID   string    `json:"logging.googleapis.com/spanId"`
}

func (e *SpanLogEntry) Log(severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&SpanLogEntry{
		Severity: severity,
		Time:     time.Now(),
		Message:  fmt.Sprintf(format, a...),
		Trace:    e.Trace,
		SpanID:   e.SpanID,
	}); err != nil {
		log.Printf("encode log entry: %v", err)
	}
}

func (e *SpanLogEntry) Debug(format string, a ...interface{}) {
	e.Log(DEBUG, format, a...)
}
