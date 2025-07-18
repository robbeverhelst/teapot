package errors

import (
	"fmt"
	"testing"
	
	tea "github.com/charmbracelet/bubbletea"
)

func TestTeapotError_Creation(t *testing.T) {
	cause := fmt.Errorf("underlying error")
	err := NewTeapotError(ErrorTypeValidation, "validation failed", cause)
	
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected error type %v, got %v", ErrorTypeValidation, err.Type)
	}
	
	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", err.Message)
	}
	
	if err.Cause != cause {
		t.Errorf("Expected cause to be set")
	}
	
	if err.Stack == "" {
		t.Error("Expected stack trace to be set")
	}
	
	if err.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

func TestTeapotError_Error(t *testing.T) {
	err := NewTeapotError(ErrorTypeValidation, "test error", nil)
	expected := "[Validation] test error"
	
	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestErrorType_String(t *testing.T) {
	tests := []struct {
		errType  ErrorType
		expected string
	}{
		{ErrorTypeValidation, "Validation"},
		{ErrorTypeNavigation, "Navigation"},
		{ErrorTypeUI, "UI"},
		{ErrorTypeSystem, "System"},
		{ErrorTypePanic, "Panic"},
		{ErrorTypeUnknown, "Unknown"},
	}
	
	for _, tt := range tests {
		if tt.errType.String() != tt.expected {
			t.Errorf("Expected %v.String() to be '%s', got '%s'", tt.errType, tt.expected, tt.errType.String())
		}
	}
}

func TestErrorRecovery_Creation(t *testing.T) {
	er := NewErrorRecovery(10, true)
	
	if er.maxLogSize != 10 {
		t.Errorf("Expected max log size 10, got %d", er.maxLogSize)
	}
	
	if !er.recoverPanic {
		t.Error("Expected panic recovery to be enabled")
	}
	
	if len(er.errorLog) != 0 {
		t.Errorf("Expected empty error log, got %d entries", len(er.errorLog))
	}
}

func TestErrorRecovery_HandleError(t *testing.T) {
	er := NewErrorRecovery(10, true)
	
	// Test with nil error
	cmd := er.HandleError(nil)
	if cmd != nil {
		t.Error("Expected nil command for nil error")
	}
	
	// Test with regular error
	testErr := fmt.Errorf("test error")
	cmd = er.HandleError(testErr)
	if cmd == nil {
		t.Error("Expected non-nil command for error")
	}
	
	// Check that error was logged
	if len(er.errorLog) != 1 {
		t.Errorf("Expected 1 error in log, got %d", len(er.errorLog))
	}
	
	// Test with TeapotError
	teapotErr := NewTeapotError(ErrorTypeValidation, "validation error", nil)
	cmd = er.HandleError(teapotErr)
	if cmd == nil {
		t.Error("Expected non-nil command for TeapotError")
	}
	
	// Check that error was logged
	if len(er.errorLog) != 2 {
		t.Errorf("Expected 2 errors in log, got %d", len(er.errorLog))
	}
}

func TestErrorRecovery_LogSizeLimit(t *testing.T) {
	er := NewErrorRecovery(3, true) // Small log size
	
	// Add more errors than the limit
	for i := 0; i < 5; i++ {
		err := fmt.Errorf("error %d", i)
		er.HandleError(err)
	}
	
	// Check that log size is limited
	if len(er.errorLog) != 3 {
		t.Errorf("Expected log size to be limited to 3, got %d", len(er.errorLog))
	}
	
	// Check that we kept the most recent errors
	if er.errorLog[0].Message != "error 2" {
		t.Errorf("Expected first error to be 'error 2', got '%s'", er.errorLog[0].Message)
	}
}

func TestErrorRecovery_ErrorStats(t *testing.T) {
	er := NewErrorRecovery(10, true)
	
	// Add various types of errors
	er.HandleError(NewTeapotError(ErrorTypeValidation, "validation 1", nil))
	er.HandleError(NewTeapotError(ErrorTypeValidation, "validation 2", nil))
	er.HandleError(NewTeapotError(ErrorTypeUI, "ui error", nil))
	er.HandleError(NewTeapotError(ErrorTypeSystem, "system error", nil))
	
	stats := er.GetErrorStats()
	
	if stats.TotalErrors != 4 {
		t.Errorf("Expected 4 total errors, got %d", stats.TotalErrors)
	}
	
	if stats.RecoverableCount != 3 {
		t.Errorf("Expected 3 recoverable errors, got %d", stats.RecoverableCount)
	}
	
	if stats.TypeBreakdown[ErrorTypeValidation] != 2 {
		t.Errorf("Expected 2 validation errors, got %d", stats.TypeBreakdown[ErrorTypeValidation])
	}
	
	if stats.TypeBreakdown[ErrorTypeUI] != 1 {
		t.Errorf("Expected 1 UI error, got %d", stats.TypeBreakdown[ErrorTypeUI])
	}
	
	if stats.TypeBreakdown[ErrorTypeSystem] != 1 {
		t.Errorf("Expected 1 system error, got %d", stats.TypeBreakdown[ErrorTypeSystem])
	}
}

func TestErrorRecovery_ClearLog(t *testing.T) {
	er := NewErrorRecovery(10, true)
	
	// Add some errors
	er.HandleError(fmt.Errorf("error 1"))
	er.HandleError(fmt.Errorf("error 2"))
	
	if len(er.errorLog) != 2 {
		t.Errorf("Expected 2 errors before clear, got %d", len(er.errorLog))
	}
	
	// Clear log
	er.ClearErrorLog()
	
	if len(er.errorLog) != 0 {
		t.Errorf("Expected 0 errors after clear, got %d", len(er.errorLog))
	}
}

func TestErrorRecovery_WrapCommand(t *testing.T) {
	er := NewErrorRecovery(10, true)
	
	// Test wrapping nil command
	wrapped := er.WrapCommand(nil)
	if wrapped != nil {
		t.Error("Expected nil when wrapping nil command")
	}
	
	// Test wrapping normal command
	normalCmd := func() tea.Msg { return "test" }
	wrapped = er.WrapCommand(normalCmd)
	if wrapped == nil {
		t.Error("Expected non-nil when wrapping normal command")
	}
	
	// Execute wrapped command
	result := wrapped()
	if result != "test" {
		t.Errorf("Expected 'test' from wrapped command, got %v", result)
	}
}

func TestIsRecoverable(t *testing.T) {
	tests := []struct {
		errType    ErrorType
		recoverable bool
	}{
		{ErrorTypeValidation, true},
		{ErrorTypeNavigation, true},
		{ErrorTypeUI, true},
		{ErrorTypePanic, true},
		{ErrorTypeSystem, false},
		{ErrorTypeUnknown, false},
	}
	
	for _, tt := range tests {
		result := isRecoverable(tt.errType)
		if result != tt.recoverable {
			t.Errorf("Expected %v to be recoverable=%v, got %v", tt.errType, tt.recoverable, result)
		}
	}
}

func TestGetRecoveryAction(t *testing.T) {
	tests := []struct {
		errType ErrorType
		action  string
	}{
		{ErrorTypeValidation, "Show validation error and allow retry"},
		{ErrorTypeNavigation, "Return to previous screen"},
		{ErrorTypeUI, "Refresh current screen"},
		{ErrorTypePanic, "Reset to safe state"},
		{ErrorTypeSystem, "Exit gracefully"},
		{ErrorTypeUnknown, "Show error and continue"},
	}
	
	for _, tt := range tests {
		result := getRecoveryAction(tt.errType)
		if result != tt.action {
			t.Errorf("Expected recovery action '%s' for %v, got '%s'", tt.action, tt.errType, result)
		}
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test NewValidationError
	err := NewValidationError("validation failed", nil)
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected validation error type, got %v", err.Type)
	}
	
	// Test NewNavigationError
	err = NewNavigationError("navigation failed", nil)
	if err.Type != ErrorTypeNavigation {
		t.Errorf("Expected navigation error type, got %v", err.Type)
	}
	
	// Test NewUIError
	err = NewUIError("ui failed", nil)
	if err.Type != ErrorTypeUI {
		t.Errorf("Expected UI error type, got %v", err.Type)
	}
	
	// Test NewSystemError
	err = NewSystemError("system failed", nil)
	if err.Type != ErrorTypeSystem {
		t.Errorf("Expected system error type, got %v", err.Type)
	}
}