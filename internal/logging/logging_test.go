package logging

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelString(t *testing.T) {
	assert.Equal(t, "DEBUG", DEBUG.String())
	assert.Equal(t, "INFO", INFO.String())
	assert.Equal(t, "WARN", WARN.String())
	assert.Equal(t, "ERROR", ERROR.String())
	assert.Equal(t, "FATAL", FATAL.String())
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
		hasError bool
	}{
		{"DEBUG", DEBUG, false},
		{"debug", DEBUG, false},
		{"INFO", INFO, false},
		{"info", INFO, false},
		{"WARN", WARN, false},
		{"warn", WARN, false},
		{"ERROR", ERROR, false},
		{"error", ERROR, false},
		{"FATAL", FATAL, false},
		{"fatal", FATAL, false},
		{"INVALID", INFO, true},
	}

	for _, test := range tests {
		level, err := ParseLogLevel(test.input)
		if test.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, level)
		}
	}
}

func TestNew(t *testing.T) {
	logger := New(DEBUG)
	assert.Equal(t, DEBUG, logger.level)
	assert.NotNil(t, logger.writer)
	assert.Equal(t, "", logger.prefix)
}

func TestNewWithPrefix(t *testing.T) {
	logger := NewWithPrefix(INFO, "test")
	assert.Equal(t, INFO, logger.level)
	assert.NotNil(t, logger.writer)
	assert.Equal(t, "test", logger.prefix)
}

func TestSetOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DEBUG)
	logger.SetOutput(&buf)
	logger.Info("test message")
	assert.Contains(t, buf.String(), "INFO: test message")
}

func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(INFO)
	logger.SetOutput(&buf)

	// Debug shouldn't log at INFO level
	logger.Debug("debug message")
	assert.Empty(t, buf.String())

	// Change level to DEBUG
	logger.SetLevel(DEBUG)
	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "DEBUG: debug message")
}

func TestSetPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := New(INFO)
	logger.SetOutput(&buf)
	logger.SetPrefix("TEST")
	logger.Info("test message")
	assert.Contains(t, buf.String(), "[TEST] INFO: test message")
}

func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DEBUG)
	logger.SetOutput(&buf)

	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "DEBUG: debug message")
	buf.Reset()

	logger.Info("info message")
	assert.Contains(t, buf.String(), "INFO: info message")
	buf.Reset()

	logger.Warn("warn message")
	assert.Contains(t, buf.String(), "WARN: warn message")
	buf.Reset()

	logger.Error("error message")
	assert.Contains(t, buf.String(), "ERROR: error message")
	buf.Reset()

	// Can't easily test Fatal as it calls os.Exit
}

func TestDefaultLoggerFunctions(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)
	defaultLogger.SetLevel(DEBUG)

	Debug("debug message")
	assert.Contains(t, buf.String(), "DEBUG: debug message")
	buf.Reset()

	Info("info message")
	assert.Contains(t, buf.String(), "INFO: info message")
	buf.Reset()

	Warn("warn message")
	assert.Contains(t, buf.String(), "WARN: warn message")
	buf.Reset()

	Error("error message")
	assert.Contains(t, buf.String(), "ERROR: error message")
	buf.Reset()

	// Can't easily test Fatal as it calls os.Exit
}

func TestSetDefaultLogLevel(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)

	SetDefaultLogLevel(ERROR)
	Info("info message")
	assert.Empty(t, buf.String())

	Error("error message")
	assert.Contains(t, buf.String(), "ERROR: error message")
}

func TestSetDefaultPrefix(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.SetOutput(&buf)
	defaultLogger.SetLevel(INFO)

	SetDefaultPrefix("DEFAULT")
	Info("info message")
	assert.Contains(t, buf.String(), "[DEFAULT] INFO: info message")
}
