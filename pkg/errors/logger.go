package errors

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var (
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.LUTC)
)

type ErrorLog struct {
	Timestamp time.Time              `json:"timestamp"`
	Type      ErrorType              `json:"type"`
	Message   string                 `json:"message"`
	Stack     string                 `json:"stack"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// LogError logs an error with all its details
func LogError(err error) {
	if err == nil {
		return
	}

	var errorLog ErrorLog
	if appErr, ok := err.(*AppError); ok {
		errorLog = ErrorLog{
			Timestamp: time.Now().UTC(),
			Type:      appErr.Type,
			Message:   appErr.Error(),
			Stack:     appErr.Stack,
			Metadata:  appErr.Metadata,
		}
	} else {
		errorLog = ErrorLog{
			Timestamp: time.Now().UTC(),
			Type:      ErrorTypeInternal,
			Message:   err.Error(),
			Stack:     getStackTrace(),
			Metadata:  make(map[string]interface{}),
		}
	}

	logJSON, _ := json.Marshal(errorLog)
	errorLogger.Println(string(logJSON))
}

// SetErrorLogger allows setting a custom logger
func SetErrorLogger(logger *log.Logger) {
	errorLogger = logger
}
