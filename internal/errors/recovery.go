// Package errors provides comprehensive error handling and recovery mechanisms
// for the Teapot CLI application. It includes panic recovery, error reporting,
// and graceful degradation strategies.
package errors

import (
	"fmt"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ErrorType represents different categories of errors in the application
type ErrorType int

const (
	// ErrorTypeUnknown represents an unclassified error
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeValidation represents input validation errors
	ErrorTypeValidation
	// ErrorTypeNavigation represents navigation/screen transition errors
	ErrorTypeNavigation
	// ErrorTypeUI represents UI rendering or interaction errors
	ErrorTypeUI
	// ErrorTypeSystem represents system-level errors (file I/O, etc.)
	ErrorTypeSystem
	// ErrorTypePanic represents recovered panic errors
	ErrorTypePanic
)

// String returns the string representation of an ErrorType
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeValidation:
		return "Validation"
	case ErrorTypeNavigation:
		return "Navigation"
	case ErrorTypeUI:
		return "UI"
	case ErrorTypeSystem:
		return "System"
	case ErrorTypePanic:
		return "Panic"
	default:
		return "Unknown"
	}
}

// TeapotError represents a structured error with context and recovery information
type TeapotError struct {
	// Type categorizes the error for handling decisions
	Type ErrorType
	// Message is the user-friendly error message
	Message string
	// Details contains technical details for debugging
	Details string
	// Cause is the underlying error that caused this error
	Cause error
	// Stack contains the stack trace when the error occurred
	Stack string
	// Timestamp records when the error occurred
	Timestamp time.Time
	// Recoverable indicates if the error can be recovered from
	Recoverable bool
	// RecoveryAction suggests what action should be taken
	RecoveryAction string
}

// Error implements the error interface
func (te *TeapotError) Error() string {
	return fmt.Sprintf("[%s] %s", te.Type, te.Message)
}

// NewTeapotError creates a new structured error
func NewTeapotError(errType ErrorType, message string, cause error) *TeapotError {
	stack := getStackTrace()
	
	return &TeapotError{
		Type:           errType,
		Message:        message,
		Details:        getErrorDetails(cause),
		Cause:          cause,
		Stack:          stack,
		Timestamp:      time.Now(),
		Recoverable:    isRecoverable(errType),
		RecoveryAction: getRecoveryAction(errType),
	}
}

// ErrorRecovery provides error recovery mechanisms for the application
type ErrorRecovery struct {
	// errorLog stores recent errors for debugging
	errorLog []TeapotError
	// maxLogSize limits the number of errors stored
	maxLogSize int
	// recoverPanic indicates if panic recovery is enabled
	recoverPanic bool
	// debugMode enables detailed error reporting
	debugMode bool
}

// NewErrorRecovery creates a new error recovery system
func NewErrorRecovery(maxLogSize int, recoverPanic bool) *ErrorRecovery {
	return &ErrorRecovery{
		errorLog:     make([]TeapotError, 0, maxLogSize),
		maxLogSize:   maxLogSize,
		recoverPanic: recoverPanic,
		debugMode:    os.Getenv("DEBUG") != "",
	}
}

// HandleError processes an error and determines the appropriate response
func (er *ErrorRecovery) HandleError(err error) tea.Cmd {
	if err == nil {
		return nil
	}
	
	// Convert to structured error if needed
	var teapotErr *TeapotError
	if te, ok := err.(*TeapotError); ok {
		teapotErr = te
	} else {
		teapotErr = NewTeapotError(ErrorTypeUnknown, err.Error(), err)
	}
	
	// Log the error
	er.logError(*teapotErr)
	
	// Handle based on error type and recoverability
	if teapotErr.Recoverable {
		return er.createRecoveryCommand(teapotErr)
	}
	
	return er.createErrorCommand(teapotErr)
}

// RecoverPanic provides panic recovery for Bubble Tea commands
func (er *ErrorRecovery) RecoverPanic() tea.Cmd {
	if !er.recoverPanic {
		return nil
	}
	
	return func() tea.Msg {
		if r := recover(); r != nil {
			// Create panic error
			message := fmt.Sprintf("Panic recovered: %v", r)
			stack := getStackTrace()
			
			panicErr := &TeapotError{
				Type:           ErrorTypePanic,
				Message:        "Application panic recovered",
				Details:        message,
				Stack:          stack,
				Timestamp:      time.Now(),
				Recoverable:    true,
				RecoveryAction: "Return to welcome screen",
			}
			
			er.logError(*panicErr)
			return ErrorRecoveredMsg{Error: panicErr}
		}
		return nil
	}
}

// WrapCommand wraps a Bubble Tea command with error recovery
func (er *ErrorRecovery) WrapCommand(cmd tea.Cmd) tea.Cmd {
	if cmd == nil {
		return nil
	}
	
	return func() tea.Msg {
		defer func() {
			if r := recover(); r != nil {
				message := fmt.Sprintf("Command panic: %v", r)
				stack := getStackTrace()
				
				panicErr := &TeapotError{
					Type:           ErrorTypePanic,
					Message:        "Command execution failed",
					Details:        message,
					Stack:          stack,
					Timestamp:      time.Now(),
					Recoverable:    true,
					RecoveryAction: "Retry or skip operation",
				}
				
				er.logError(*panicErr)
			}
		}()
		
		return cmd()
	}
}

// GetErrorLog returns the recent error log
func (er *ErrorRecovery) GetErrorLog() []TeapotError {
	return er.errorLog
}

// ClearErrorLog clears the error log
func (er *ErrorRecovery) ClearErrorLog() {
	er.errorLog = er.errorLog[:0]
}

// GetErrorStats returns statistics about errors
func (er *ErrorRecovery) GetErrorStats() ErrorStats {
	typeCount := make(map[ErrorType]int)
	recoverableCount := 0
	
	for _, err := range er.errorLog {
		typeCount[err.Type]++
		if err.Recoverable {
			recoverableCount++
		}
	}
	
	return ErrorStats{
		TotalErrors:      len(er.errorLog),
		RecoverableCount: recoverableCount,
		TypeBreakdown:    typeCount,
	}
}

// ErrorStats provides statistics about error occurrences
type ErrorStats struct {
	TotalErrors      int
	RecoverableCount int
	TypeBreakdown    map[ErrorType]int
}

// logError adds an error to the error log
func (er *ErrorRecovery) logError(err TeapotError) {
	// Add to log
	er.errorLog = append(er.errorLog, err)
	
	// Trim log if too large
	if len(er.errorLog) > er.maxLogSize {
		er.errorLog = er.errorLog[len(er.errorLog)-er.maxLogSize:]
	}
	
	// Debug logging
	if er.debugMode {
		fmt.Fprintf(os.Stderr, "ERROR [%s]: %s\n", err.Type, err.Message)
		if err.Details != "" {
			fmt.Fprintf(os.Stderr, "  Details: %s\n", err.Details)
		}
		if err.Stack != "" {
			fmt.Fprintf(os.Stderr, "  Stack: %s\n", err.Stack)
		}
	}
}

// createRecoveryCommand creates a command for recoverable errors
func (er *ErrorRecovery) createRecoveryCommand(err *TeapotError) tea.Cmd {
	return func() tea.Msg {
		return ErrorRecoveredMsg{
			Error:   err,
			Message: fmt.Sprintf("Recovered from %s error: %s", err.Type, err.Message),
		}
	}
}

// createErrorCommand creates a command for non-recoverable errors
func (er *ErrorRecovery) createErrorCommand(err *TeapotError) tea.Cmd {
	return func() tea.Msg {
		return ErrorOccurredMsg{
			Error:   err,
			Message: fmt.Sprintf("Error occurred: %s", err.Message),
		}
	}
}

// ErrorRecoveredMsg is sent when an error has been recovered
type ErrorRecoveredMsg struct {
	Error   *TeapotError
	Message string
}

// ErrorOccurredMsg is sent when an error occurs
type ErrorOccurredMsg struct {
	Error   *TeapotError
	Message string
}

// Helper functions

// getStackTrace returns the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024*4)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// getErrorDetails extracts detailed information from an error
func getErrorDetails(err error) string {
	if err == nil {
		return ""
	}
	
	// Try to get more context from the error
	details := err.Error()
	
	// Add type information if available
	if errType := fmt.Sprintf("%T", err); errType != "*errors.errorString" {
		details = fmt.Sprintf("%s (%s)", details, errType)
	}
	
	return details
}

// isRecoverable determines if an error type can be recovered from
func isRecoverable(errType ErrorType) bool {
	switch errType {
	case ErrorTypeValidation, ErrorTypeNavigation, ErrorTypeUI:
		return true
	case ErrorTypePanic:
		return true // Most panics in UI can be recovered
	case ErrorTypeSystem:
		return false // System errors are usually fatal
	default:
		return false
	}
}

// getRecoveryAction suggests an appropriate recovery action
func getRecoveryAction(errType ErrorType) string {
	switch errType {
	case ErrorTypeValidation:
		return "Show validation error and allow retry"
	case ErrorTypeNavigation:
		return "Return to previous screen"
	case ErrorTypeUI:
		return "Refresh current screen"
	case ErrorTypePanic:
		return "Reset to safe state"
	case ErrorTypeSystem:
		return "Exit gracefully"
	default:
		return "Show error and continue"
	}
}

// Helper functions for creating specific error types

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *TeapotError {
	return NewTeapotError(ErrorTypeValidation, message, cause)
}

// NewNavigationError creates a navigation error
func NewNavigationError(message string, cause error) *TeapotError {
	return NewTeapotError(ErrorTypeNavigation, message, cause)
}

// NewUIError creates a UI error
func NewUIError(message string, cause error) *TeapotError {
	return NewTeapotError(ErrorTypeUI, message, cause)
}

// NewSystemError creates a system error
func NewSystemError(message string, cause error) *TeapotError {
	return NewTeapotError(ErrorTypeSystem, message, cause)
}