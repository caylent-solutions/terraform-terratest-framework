package unit

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/cmd/tftest/logger"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    logger.LogLevel
		expected string
	}{
		{logger.DEBUG, "DEBUG"},
		{logger.INFO, "INFO"},
		{logger.WARN, "WARN"},
		{logger.ERROR, "ERROR"},
		{logger.FATAL, "FATAL"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.level.String(), "LogLevel.String() should return the correct string representation")
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected logger.LogLevel
		hasError bool
	}{
		{"DEBUG", logger.DEBUG, false},
		{"INFO", logger.INFO, false},
		{"WARN", logger.WARN, false},
		{"ERROR", logger.ERROR, false},
		{"FATAL", logger.FATAL, false},
		{"debug", logger.DEBUG, false}, // Test case insensitivity
		{"info", logger.INFO, false},   // Test case insensitivity
		{"invalid", logger.INFO, true}, // Test invalid input
	}

	for _, test := range tests {
		level, err := logger.ParseLogLevel(test.input)
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
	log := logger.New(logger.INFO)
	log.SetOutput(&buf)

	// Test logging at different levels
	log.Debug("Debug message")
	assert.Empty(t, buf.String(), "Debug message should not be logged at INFO level")

	buf.Reset()
	log.Info("Info message")
	assert.Contains(t, buf.String(), "INFO: Info message", "Info message should be logged at INFO level")

	buf.Reset()
	log.Warn("Warn message")
	assert.Contains(t, buf.String(), "WARN: Warn message", "Warn message should be logged at INFO level")

	buf.Reset()
	log.Error("Error message")
	assert.Contains(t, buf.String(), "ERROR: Error message", "Error message should be logged at INFO level")

	// Test changing log level
	buf.Reset()
	log.SetLevel(logger.ERROR)
	log.Info("Info message")
	assert.Empty(t, buf.String(), "Info message should not be logged at ERROR level")

	buf.Reset()
	log.Error("Error message")
	assert.Contains(t, buf.String(), "ERROR: Error message", "Error message should be logged at ERROR level")
}

func TestLoggerFormat(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger with the buffer as output
	log := logger.New(logger.INFO)
	log.SetOutput(&buf)

	// Test formatting
	log.Info("Count: %d, String: %s", 42, "test")
	logOutput := buf.String()

	assert.Contains(t, logOutput, "Count: 42", "Log should contain formatted count")
	assert.Contains(t, logOutput, "String: test", "Log should contain formatted string")
}

func TestLoggerTimestamp(t *testing.T) {
	// Create a buffer to capture log output
	var buf bytes.Buffer

	// Create a logger with the buffer as output
	log := logger.New(logger.INFO)
	log.SetOutput(&buf)

	// Test timestamp format
	log.Info("Test message")
	logOutput := buf.String()

	// Check for timestamp format YYYY-MM-DD HH:MM:SS
	assert.Regexp(t, `\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\]`, logOutput, "Log should contain timestamp in correct format")
}
