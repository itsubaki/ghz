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
	ProjectID   string
	Trace       string
	ErrorClient *errorreporting.Client
	Request     *http.Request
}

func New(projectID string, traceID ...string) *Logger {
	trace := ""
	if len(traceID) > 0 {
		trace = fmt.Sprintf("projects/%v/traces/%v", projectID, traceID[0])
	}

	return &Logger{
		ProjectID: projectID,
		Trace:     trace,
	}
}

func (l *Logger) LogWith(spanID, severity, format string, a ...interface{}) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.Trace,
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

func (l *Logger) Default(format string, a ...interface{}) {
	l.Log(DEFAULT, format, a...)
}

func (l *Logger) Debug(format string, a ...interface{}) {
	l.Log(DEBUG, format, a...)
}

func (l *Logger) Info(format string, a ...interface{}) {
	l.Log(INFO, format, a...)
}

func (l *Logger) Notice(format string, a ...interface{}) {
	l.Log(NOTICE, format, a...)
}

func (l *Logger) Warning(format string, a ...interface{}) {
	l.Log(WARNING, format, a...)
}

func (l *Logger) Error(format string, a ...interface{}) {
	l.Error(format, a...)
}

func (l *Logger) Critical(format string, a ...interface{}) {
	l.Log(CRITICAL, format, a...)
}

func (l *Logger) Alert(format string, a ...interface{}) {
	l.Log(ALERT, format, a...)
}

func (l *Logger) Emergency(format string, a ...interface{}) {
	l.Log(EMERGENCY, format, a...)
}

func (l *Logger) NewReport(ctx context.Context, req *http.Request) *Logger {
	c, err := errorreporting.NewClient(ctx, l.ProjectID, errorreporting.Config{})
	if err != nil {
		l.Error("new error report client: %v", err)
		return l
	}

	l.ErrorClient = c
	l.Request = req
	return l
}

func (l *Logger) ReportWith(spanID, severity, format string, a ...interface{}) {
	l.LogWith(spanID, severity, format, a...)
	if l.ErrorClient == nil {
		return
	}

	for _, aa := range a {
		switch err := aa.(type) {
		case error:
			l.ErrorClient.Report(errorreporting.Entry{
				Error: err,
				Req:   l.Request,
			})
		}
	}
}

func (l *Logger) ErrorReportWith(spanID, format string, a ...interface{}) {
	l.ReportWith(spanID, ERROR, format, a...)
}

func (l *Logger) CriticalReportWith(spanID, format string, a ...interface{}) {
	l.ReportWith(spanID, CRITICAL, format, a...)
}

func (l *Logger) AlertReportWith(spanID, format string, a ...interface{}) {
	l.ReportWith(spanID, ALERT, format, a...)
}

func (l *Logger) EmergencyReportWith(spanID, format string, a ...interface{}) {
	l.ReportWith(spanID, EMERGENCY, format, a...)
}

func (l *Logger) Report(severity, format string, a ...interface{}) {
	l.ReportWith("", severity, format, a...)
}

func (l *Logger) ErrorReport(format string, a ...interface{}) {
	l.Report(ERROR, format, a...)
}

func (l *Logger) CriticalReport(format string, a ...interface{}) {
	l.Report(CRITICAL, format, a...)
}

func (l *Logger) AlertReport(format string, a ...interface{}) {
	l.Report(ALERT, format, a...)
}

func (l *Logger) EmergencyReport(format string, a ...interface{}) {
	l.Report(EMERGENCY, format, a...)
}
