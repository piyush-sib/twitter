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
	Level          LogLevel  `json:"level"`
	Timestamp      time.Time `json:"timestamp"`
	TimeProcessing int       `json:"time_processing"`
	Message        string    `json:"message"`
	HTTPRoute      string    `json:"http_route"`
	Error          error     `json:"error,omitempty"`
	UserID         int       `json:"user_id,omitempty"`
}
