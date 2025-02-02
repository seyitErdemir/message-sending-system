package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorType string

const (
	ErrorTypeWebhook    ErrorType = "WEBHOOK_ERROR"
	ErrorTypeDatabase   ErrorType = "DATABASE_ERROR"
	ErrorTypeCache      ErrorType = "CACHE_ERROR"
	ErrorTypeCron       ErrorType = "CRON_ERROR"
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	ErrorTypeInternal   ErrorType = "INTERNAL_ERROR"
)

// AppError represents a custom application error
type AppError struct {
	Type     ErrorType
	Message  string
	Original error
	Stack    string
	Metadata map[string]interface{}
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Original != nil {
		return fmt.Sprintf("%s: %s - %v", e.Type, e.Message, e.Original)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewError creates a new AppError
func NewError(errType ErrorType, message string, original error) *AppError {
	stack := getStackTrace()
	return &AppError{
		Type:     errType,
		Message:  message,
		Original: original,
		Stack:    stack,
		Metadata: make(map[string]interface{}),
	}
}

// WithMetadata adds metadata to the error
func (e *AppError) WithMetadata(key string, value interface{}) *AppError {
	e.Metadata[key] = value
	return e
}

// IsType checks if the error is of a specific type
func IsType(err error, errType ErrorType) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Type == errType
	}
	return false
}

// getStackTrace returns the stack trace as a string
func getStackTrace() string {
	var sb strings.Builder
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		sb.WriteString(fmt.Sprintf("%s:%d\n", file, line))
	}
	return sb.String()
}

// Webhook specific errors
func NewWebhookError(message string, original error) *AppError {
	return NewError(ErrorTypeWebhook, message, original)
}

// Database specific errors
func NewDatabaseError(message string, original error) *AppError {
	return NewError(ErrorTypeDatabase, message, original)
}

// Cache specific errors
func NewCacheError(message string, original error) *AppError {
	return NewError(ErrorTypeCache, message, original)
}

// Cron specific errors
func NewCronError(message string, original error) *AppError {
	return NewError(ErrorTypeCron, message, original)
}
