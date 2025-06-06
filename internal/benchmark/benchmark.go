package benchmark

import (
	"fmt"
	"time"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/logging"
)

// BenchmarkResult represents the result of a benchmark
type BenchmarkResult struct {
	Name     string
	Duration time.Duration
	Success  bool
	Error    error
}

// String returns a string representation of the benchmark result
func (r *BenchmarkResult) String() string {
	status := "✅ Success"
	if !r.Success {
		status = fmt.Sprintf("❌ Failed: %v", r.Error)
	}
	return fmt.Sprintf("%s: %s (%s)", r.Name, status, r.Duration)
}

// Benchmark runs a function and measures its execution time
func Benchmark(name string, fn func() error) *BenchmarkResult {
	logger := logging.NewWithPrefix(logging.INFO, "Benchmark")
	logger.Info("Starting benchmark: %s", name)

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	result := &BenchmarkResult{
		Name:     name,
		Duration: duration,
		Success:  err == nil,
		Error:    err,
	}

	if err != nil {
		logger.Error("Benchmark failed: %s (%s): %v", name, duration, err)
	} else {
		logger.Info("Benchmark completed: %s (%s)", name, duration)
	}

	return result
}

// BenchmarkSuite represents a suite of benchmarks
type BenchmarkSuite struct {
	Name    string
	Results []*BenchmarkResult
	logger  *logging.Logger
}

// NewBenchmarkSuite creates a new benchmark suite
func NewBenchmarkSuite(name string) *BenchmarkSuite {
	return &BenchmarkSuite{
		Name:    name,
		Results: []*BenchmarkResult{},
		logger:  logging.NewWithPrefix(logging.INFO, fmt.Sprintf("BenchmarkSuite[%s]", name)),
	}
}

// Run runs a benchmark and adds its result to the suite
func (s *BenchmarkSuite) Run(name string, fn func() error) *BenchmarkResult {
	s.logger.Info("Running benchmark: %s", name)
	result := Benchmark(name, fn)
	s.Results = append(s.Results, result)
	return result
}

// Summary returns a summary of the benchmark suite
func (s *BenchmarkSuite) Summary() string {
	var totalDuration time.Duration
	successCount := 0
	failureCount := 0

	for _, result := range s.Results {
		totalDuration += result.Duration
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	return fmt.Sprintf(
		"Benchmark Suite: %s\n"+
			"Total Duration: %s\n"+
			"Total Benchmarks: %d\n"+
			"Successful: %d\n"+
			"Failed: %d",
		s.Name,
		totalDuration,
		len(s.Results),
		successCount,
		failureCount,
	)
}

// PrintSummary prints a summary of the benchmark suite
func (s *BenchmarkSuite) PrintSummary() {
	s.logger.Info("Benchmark Suite Summary:")
	s.logger.Info("------------------------")
	s.logger.Info("Suite: %s", s.Name)
	s.logger.Info("Total Benchmarks: %d", len(s.Results))

	var totalDuration time.Duration
	successCount := 0
	failureCount := 0

	for _, result := range s.Results {
		totalDuration += result.Duration
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	s.logger.Info("Total Duration: %s", totalDuration)
	s.logger.Info("Successful: %d", successCount)
	s.logger.Info("Failed: %d", failureCount)
	s.logger.Info("------------------------")

	s.logger.Info("Individual Results:")
	for _, result := range s.Results {
		if result.Success {
			s.logger.Info("✅ %s: %s", result.Name, result.Duration)
		} else {
			s.logger.Error("❌ %s: %s - %v", result.Name, result.Duration, result.Error)
		}
	}
}
