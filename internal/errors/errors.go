package errors

import (
	"fmt"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// ConfigError represents configuration errors
	ConfigError ErrorType = "ConfigError"
	
	// ValidationError represents validation errors
	ValidationError ErrorType = "ValidationError"
	
	// TerraformError represents errors from Terraform operations
	TerraformError ErrorType = "TerraformError"
	
	// AssertionError represents errors from test assertions
	AssertionError ErrorType = "AssertionError"
	
	// InternalError represents internal framework errors
	InternalError ErrorType = "InternalError"
)

// FrameworkError represents an error in the Terraform Test Framework
type FrameworkError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error returns the error message
func (e *FrameworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause of the error
func (e *FrameworkError) Unwrap() error {
	return e.Cause
}

// NewConfigError creates a new configuration error
func NewConfigError(message string, cause error) *FrameworkError {
	return &FrameworkError{
		Type:    ConfigError,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string, cause error) *FrameworkError {
	return &FrameworkError{
		Type:    ValidationError,
		Message: message,
		Cause:   cause,
	}
}

// NewTerraformError creates a new Terraform error
func NewTerraformError(message string, cause error) *FrameworkError {
	return &FrameworkError{
		Type:    TerraformError,
		Message: message,
		Cause:   cause,
	}
}

// NewAssertionError creates a new assertion error
func NewAssertionError(message string, cause error) *FrameworkError {
	return &FrameworkError{
		Type:    AssertionError,
		Message: message,
		Cause:   cause,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(message string, cause error) *FrameworkError {
	return &FrameworkError{
		Type:    InternalError,
		Message: message,
		Cause:   cause,
	}
}