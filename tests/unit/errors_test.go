package unit

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	frameworkErrors "github.com/caylent-solutions/terraform-terratest-framework/internal/errors"
)

func TestErrorString(t *testing.T) {
	// Test error with cause
	cause := errors.New("underlying error")
	err := &frameworkErrors.FrameworkError{
		Type:    frameworkErrors.ConfigError,
		Message: "config error message",
		Cause:   cause,
	}
	expected := "ConfigError: config error message (cause: underlying error)"
	assert.Equal(t, expected, err.Error())

	// Test error without cause
	err = &frameworkErrors.FrameworkError{
		Type:    frameworkErrors.ValidationError,
		Message: "validation error message",
	}
	expected = "ValidationError: validation error message"
	assert.Equal(t, expected, err.Error())
}

func TestUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &frameworkErrors.FrameworkError{
		Type:    frameworkErrors.ConfigError,
		Message: "config error message",
		Cause:   cause,
	}
	assert.Equal(t, cause, err.Unwrap())
}

func TestNewConfigError(t *testing.T) {
	cause := errors.New("underlying error")
	err := frameworkErrors.NewConfigError("config error message", cause)
	assert.Equal(t, frameworkErrors.ConfigError, err.Type)
	assert.Equal(t, "config error message", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewValidationError(t *testing.T) {
	cause := errors.New("underlying error")
	err := frameworkErrors.NewValidationError("validation error message", cause)
	assert.Equal(t, frameworkErrors.ValidationError, err.Type)
	assert.Equal(t, "validation error message", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewTerraformError(t *testing.T) {
	cause := errors.New("underlying error")
	err := frameworkErrors.NewTerraformError("terraform error message", cause)
	assert.Equal(t, frameworkErrors.TerraformError, err.Type)
	assert.Equal(t, "terraform error message", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewAssertionError(t *testing.T) {
	cause := errors.New("underlying error")
	err := frameworkErrors.NewAssertionError("assertion error message", cause)
	assert.Equal(t, frameworkErrors.AssertionError, err.Type)
	assert.Equal(t, "assertion error message", err.Message)
	assert.Equal(t, cause, err.Cause)
}

func TestNewInternalError(t *testing.T) {
	cause := errors.New("underlying error")
	err := frameworkErrors.NewInternalError("internal error message", cause)
	assert.Equal(t, frameworkErrors.InternalError, err.Type)
	assert.Equal(t, "internal error message", err.Message)
	assert.Equal(t, cause, err.Cause)
}
