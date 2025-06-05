package benchmark

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBenchmark(t *testing.T) {
	// Test successful benchmark
	result := Benchmark("test-success", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	assert.Equal(t, "test-success", result.Name)
	assert.True(t, result.Success)
	assert.Nil(t, result.Error)
	assert.True(t, result.Duration >= 10*time.Millisecond)

	// Test failed benchmark
	expectedErr := errors.New("test error")
	result = Benchmark("test-failure", func() error {
		time.Sleep(10 * time.Millisecond)
		return expectedErr
	})

	assert.Equal(t, "test-failure", result.Name)
	assert.False(t, result.Success)
	assert.Equal(t, expectedErr, result.Error)
	assert.True(t, result.Duration >= 10*time.Millisecond)
}

func TestBenchmarkSuite(t *testing.T) {
	suite := NewBenchmarkSuite("test-suite")

	// Run successful benchmark
	result1 := suite.Run("test-success", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	// Run failed benchmark
	expectedErr := errors.New("test error")
	result2 := suite.Run("test-failure", func() error {
		time.Sleep(10 * time.Millisecond)
		return expectedErr
	})

	// Check results
	assert.Equal(t, 2, len(suite.Results))
	assert.Equal(t, result1, suite.Results[0])
	assert.Equal(t, result2, suite.Results[1])

	// Check summary
	summary := suite.Summary()
	assert.Contains(t, summary, "Benchmark Suite: test-suite")
	assert.Contains(t, summary, "Total Benchmarks: 2")
	assert.Contains(t, summary, "Successful: 1")
	assert.Contains(t, summary, "Failed: 1")

	// Test PrintSummary (just make sure it doesn't panic)
	assert.NotPanics(t, func() {
		suite.PrintSummary()
	})
}
