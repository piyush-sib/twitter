package structuredlogger

import (
	"encoding/json"
	"go.uber.org/dig"
	"log"
	"os"
	"time"
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

// Log creates and writes a structured log entry to file for debugging later any request.
func (l *JSONLogger) Log(entry *models.LogEntry, requestTime time.Time) {
	entry.TimeProcessing = int(time.Now().Sub(requestTime).Nanoseconds() / 1e6)
	jsonData, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		// Fallback to default logging in case of a marshaling failure
		l.logger.Printf("Failed to marshal log entry: %v", jsonErr)
		return
	}

	l.logger.Println(string(jsonData))
}
