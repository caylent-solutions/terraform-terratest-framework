package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{FATAL, "FATAL"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.String(), "LogLevel.String() should return the correct string representation")
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
		hasError bool
	}{
		{"DEBUG", DEBUG, false},
		{"INFO", INFO, false},
		{"WARN", WARN, false},
		{"ERROR", ERROR, false},
		{"FATAL", FATAL, false},
		{"debug", DEBUG, false}, // Test case insensitivity
		{"info", INFO, false},   // Test case insensitivity
		{"invalid", INFO, true}, // Test invalid input
	}

	for _, test := range tests {
		level, err := ParseLogLevel(test.input)
		if test.hasError {
			assert.Error(t, err, "ParseLogLevel should return an error for invalid input")
		} else {
			assert.NoError(t, err, "ParseLogLevel should not return an error for valid input")
			assert.Equal(t, test.expected, level, "ParseLogLevel should return the correct LogLevel")
		}
	}
}

func TestLogger(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger with the buffer as output
	logger := New(INFO)
	logger.SetOutput(&buf)

	// Test logging at different levels
	logger.Debug("Debug message")
	assert.Empty(t, buf.String(), "Debug message should not be logged at INFO level")

	buf.Reset()
	logger.Info("Info message")
	assert.Contains(t, buf.String(), "INFO: Info message", "Info message should be logged at INFO level")

	buf.Reset()
	logger.Warn("Warn message")
	assert.Contains(t, buf.String(), "WARN: Warn message", "Warn message should be logged at INFO level")

	buf.Reset()
	logger.Error("Error message")
	assert.Contains(t, buf.String(), "ERROR: Error message", "Error message should be logged at INFO level")

	// Test changing log level
	buf.Reset()
	logger.SetLevel(ERROR)
	logger.Info("Info message")
	assert.Empty(t, buf.String(), "Info message should not be logged at ERROR level")

	buf.Reset()
	logger.Error("Error message")
	assert.Contains(t, buf.String(), "ERROR: Error message", "Error message should be logged at ERROR level")
}

func TestDefaultLogger(t *testing.T) {
	// Test that the default logger functions don't crash
	// We can't easily test their output without modifying the default logger
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	// Test setting the default log level
	SetDefaultLogLevel(ERROR)

	// This is just to ensure the function exists and can be called
	assert.NotPanics(t, func() {
		SetDefaultLogLevel(INFO)
	}, "SetDefaultLogLevel should not panic")
}

func TestLoggerFormat(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger with the buffer as output
	logger := New(INFO)
	logger.SetOutput(&buf)

	// Test formatting
	logger.Info("Count: %d, String: %s", 42, "test")
	logOutput := buf.String()

	assert.Contains(t, logOutput, "Count: 42", "Log should contain formatted count")
	assert.Contains(t, logOutput, "String: test", "Log should contain formatted string")
}

func TestLoggerTimestamp(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger with the buffer as output
	logger := New(INFO)
	logger.SetOutput(&buf)

	// Test timestamp format
	logger.Info("Test message")
	logOutput := buf.String()

	// Check for timestamp format YYYY-MM-DD HH:MM:SS
	assert.Regexp(t, `\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\]`, logOutput, "Log should contain timestamp in correct format")
}
