package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
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

func (l *Logger) Log(severity, message string) {
	if err := json.NewEncoder(os.Stdout).Encode(&LogEntry{
		Time:     time.Now(),
		Trace:    l.Trace,
		Severity: severity,
		Message:  message,
	}); err != nil {
		panic(err)
	}
}

func (l *Logger) Default(message string) {
	l.Log("Default", message)
}

func (l *Logger) Debug(message string) {
	l.Log("Debug", message)
}

func (l *Logger) Info(message string) {
	l.Log("Info", message)
}

func (l *Logger) Notice(message string) {
	l.Log("Notice", message)
}

func (l *Logger) Warning(message string) {
	l.Log("Warning", message)
}

func (l *Logger) Error(message string) {
	l.Log("Error", message)
}

func (l *Logger) Critical(message string) {
	l.Log("Critical", message)
}

func (l *Logger) Alert(message string) {
	l.Log("Alert", message)
}

func (l *Logger) Emergency(message string) {
	l.Log("Emergency", message)
}
