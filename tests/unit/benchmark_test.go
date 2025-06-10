package unit

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/benchmark"
)

func TestBenchmarkResultString(t *testing.T) {
	// Test successful benchmark
	result := &benchmark.BenchmarkResult{
		Name:     "test-benchmark",
		Duration: 123 * time.Millisecond,
		Success:  true,
	}

	expected := "test-benchmark: ✅ Success (123ms)"
	assert.Equal(t, expected, result.String())

	// Test failed benchmark
	result = &benchmark.BenchmarkResult{
		Name:     "failed-benchmark",
		Duration: 456 * time.Millisecond,
		Success:  false,
		Error:    errors.New("test error"),
	}

	expected = "failed-benchmark: ❌ Failed: test error (456ms)"
	assert.Equal(t, expected, result.String())
}

func TestBenchmarkFunction(t *testing.T) {
	// Test successful benchmark
	result := benchmark.Benchmark("test-success", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	assert.Equal(t, "test-success", result.Name)
	assert.True(t, result.Duration >= 10*time.Millisecond)
	assert.True(t, result.Success)
	assert.Nil(t, result.Error)

	// Test failed benchmark
	testErr := errors.New("test error")
	result = benchmark.Benchmark("test-failure", func() error {
		time.Sleep(10 * time.Millisecond)
		return testErr
	})

	assert.Equal(t, "test-failure", result.Name)
	assert.True(t, result.Duration >= 10*time.Millisecond)
	assert.False(t, result.Success)
	assert.Equal(t, testErr, result.Error)
}

func TestNewBenchmarkSuite(t *testing.T) {
	suite := benchmark.NewBenchmarkSuite("test-suite")
	assert.NotNil(t, suite)
	assert.Equal(t, "test-suite", suite.Name)
	assert.Empty(t, suite.Results)
}

func TestBenchmarkSuiteRun(t *testing.T) {
	suite := benchmark.NewBenchmarkSuite("test-suite")

	// Add a successful benchmark
	result := suite.Run("test-benchmark", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	})

	assert.Equal(t, 1, len(suite.Results))
	assert.Equal(t, result, suite.Results[0])
	assert.Equal(t, "test-benchmark", result.Name)
	assert.True(t, result.Duration >= 10*time.Millisecond)
	assert.True(t, result.Success)

	// Add a failed benchmark
	testErr := errors.New("test error")
	result = suite.Run("failed-benchmark", func() error {
		return testErr
	})

	assert.Equal(t, 2, len(suite.Results))
	assert.Equal(t, result, suite.Results[1])
	assert.Equal(t, "failed-benchmark", result.Name)
	assert.False(t, result.Success)
	assert.Equal(t, testErr, result.Error)
}

func TestBenchmarkSuiteSummary(t *testing.T) {
	suite := benchmark.NewBenchmarkSuite("test-suite")

	// Add benchmarks
	suite.Run("test1", func() error { return nil })
	suite.Run("test2", func() error { return errors.New("error") })
	suite.Run("test3", func() error { return nil })

	summary := suite.Summary()
	assert.Contains(t, summary, "Benchmark Suite: test-suite")
	assert.Contains(t, summary, "Total Benchmarks: 3")
	assert.Contains(t, summary, "Successful: 2")
	assert.Contains(t, summary, "Failed: 1")
}

func TestBenchmarkSuitePrintSummary(t *testing.T) {
	// This test is mostly to ensure the function doesn't panic
	suite := benchmark.NewBenchmarkSuite("test-suite")

	// Add benchmarks
	suite.Run("test1", func() error { return nil })
	suite.Run("test2", func() error { return errors.New("error") })

	// Should not panic
	suite.PrintSummary()
}
