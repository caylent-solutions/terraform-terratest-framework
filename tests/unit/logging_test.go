package unit

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/logging"
)

func TestInternalLogLevelString(t *testing.T) {
	assert.Equal(t, "DEBUG", logging.DEBUG.String())
	assert.Equal(t, "INFO", logging.INFO.String())
	assert.Equal(t, "WARN", logging.WARN.String())
	assert.Equal(t, "ERROR", logging.ERROR.String())
	assert.Equal(t, "FATAL", logging.FATAL.String())
}

func TestInternalParseLogLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected logging.LogLevel
		hasError bool
	}{
		{"DEBUG", logging.DEBUG, false},
		{"debug", logging.DEBUG, false},
		{"INFO", logging.INFO, false},
		{"info", logging.INFO, false},
		{"WARN", logging.WARN, false},
		{"warn", logging.WARN, false},
		{"ERROR", logging.ERROR, false},
		{"error", logging.ERROR, false},
		{"FATAL", logging.FATAL, false},
		{"fatal", logging.FATAL, false},
		{"INVALID", logging.INFO, true},
	}

	for _, test := range tests {
		level, err := logging.ParseLogLevel(test.input)
		if test.hasError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, level)
		}
	}
}

func TestInternalNew(t *testing.T) {
	logger := logging.New(logging.DEBUG)
	assert.NotNil(t, logger)
}

func TestInternalNewWithPrefix(t *testing.T) {
	logger := logging.NewWithPrefix(logging.INFO, "test")
	assert.NotNil(t, logger)
}

func TestInternalSetOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(logging.DEBUG)
	logger.SetOutput(&buf)
	logger.Info("test message")
	assert.Contains(t, buf.String(), "INFO: test message")
}

func TestInternalSetLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(logging.INFO)
	logger.SetOutput(&buf)

	// Debug shouldn't log at INFO level
	logger.Debug("debug message")
	assert.Empty(t, buf.String())

	// Change level to DEBUG
	logger.SetLevel(logging.DEBUG)
	logger.Debug("debug message")
	assert.Contains(t, buf.String(), "DEBUG: debug message")
}

func TestInternalSetPrefix(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(logging.INFO)
	logger.SetOutput(&buf)
	logger.SetPrefix("TEST")
	logger.Info("test message")
	assert.Contains(t, buf.String(), "[TEST] INFO: test message")
}

func TestInternalLogLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(logging.DEBUG)
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

func TestInternalDefaultLoggerFunctions(t *testing.T) {
	// Just test that these functions exist and can be called
	assert.NotPanics(t, func() {
		logging.SetDefaultLogLevel(logging.INFO)
		logging.SetDefaultPrefix("TEST")
		logging.Debug("debug message")
		logging.Info("info message")
		logging.Warn("warn message")
		logging.Error("error message")
		// Can't test Fatal as it calls os.Exit
	})
}
