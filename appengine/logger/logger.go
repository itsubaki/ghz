package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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
	ProjectID string
	TraceID   string
	Trace     string
}

func New(projectID, traceID string) *Logger {
	return &Logger{
		ProjectID: projectID,
		TraceID:   traceID,
		Trace:     fmt.Sprintf("projects/%v/traces/%v", projectID, traceID),
	}
}

func (l *Logger) Log(severity, format string, a ...interface{}) {
	if l.TraceID == "" {
		return
	}

	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.Trace,
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
	}); err != nil {
		panic(err)
	}
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
	l.Log(ERROR, format, a...)
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
