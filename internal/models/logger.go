package models

import "time"

// LogLevel defines the severity of the log.
type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// LogEntry represents a structured log entry.
type LogEntry struct {
	Level     LogLevel  `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Error     string    `json:"error,omitempty"`
	Data      any       `json:"data,omitempty"`
}
