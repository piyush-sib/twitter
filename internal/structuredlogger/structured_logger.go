package structuredlogger

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"go.uber.org/dig"
	"twitter/internal/closer"
	"twitter/internal/models"
)

type flagVal struct {
	dig.In

	StructuredLogFile string `name:"structured-log-file"`
}

// JSONLogger is a logger for structured JSON logs.
type JSONLogger struct {
	logger *log.Logger
	file   *os.File
}

// NewStructuredLogger creates a new instance of JSONLogger and writes logs to the specified file.
func NewStructuredLogger(flg flagVal) (*JSONLogger, closer.CloserResult, error) {
	// Open the log file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(flg.StructuredLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, closer.CloserResult{}, err
	}
	// this will ensure that the file is closed when the application is closed
	cl := func() error {
		return file.Close()
	}
	return &JSONLogger{
		logger: log.New(file, "", 0),
		file:   file,
	}, closer.CloserResult{Close: cl}, nil
}

// Info logs an informational message.
func (l *JSONLogger) Info(message string, data any) {
	l.log(models.LevelInfo, message, nil, data)
}

// Warn logs a warning message.
func (l *JSONLogger) Warn(message string, data any) {
	l.log(models.LevelWarn, message, nil, data)
}

// Error logs an error message.
func (l *JSONLogger) Error(message string, err error, data any) {
	l.log(models.LevelError, message, err, data)
}

// log creates and writes a structured log entry.
func (l *JSONLogger) log(level models.LogLevel, message string, err error, data any) {
	entry := models.LogEntry{
		Level:     level,
		Timestamp: time.Now(),
		Message:   message,
	}

	if err != nil {
		entry.Error = err.Error()
	}
	if data != nil {
		entry.Data = data
	}

	jsonData, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		// Fallback to default logging in case of a marshaling failure
		l.logger.Printf("Failed to marshal log entry: %v", jsonErr)
		return
	}

	l.logger.Println(string(jsonData))
}
