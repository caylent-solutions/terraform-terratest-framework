# Benchmarking and Performance Optimization

This document describes the benchmarking tools and performance optimization strategies used in the Terraform Test Framework.

## Overview

The framework includes a benchmarking package that allows you to measure the performance of your tests and identify bottlenecks. This can help you optimize your tests and improve their execution time.

## Benchmarking Tools

### Single Benchmark

To benchmark a single function:

```go
import "github.com/caylent-solutions/terraform-test-framework/internal/benchmark"

// Benchmark a function
result := benchmark.Benchmark("terraform-apply", func() error {
    // Function to benchmark
    return terraform.Apply(t, terraformOptions)
})

// Print the result
fmt.Println(result) // Output: terraform-apply: ✅ Success (2.5s)
```

### Benchmark Suite

To benchmark multiple functions as a suite:

```go
import "github.com/caylent-solutions/terraform-test-framework/internal/benchmark"

// Create a benchmark suite
suite := benchmark.NewBenchmarkSuite("terraform-operations")

// Run benchmarks
suite.Run("terraform-init", func() error {
    return terraform.Init(t, terraformOptions)
})

suite.Run("terraform-apply", func() error {
    return terraform.Apply(t, terraformOptions)
})

suite.Run("terraform-destroy", func() error {
    return terraform.Destroy(t, terraformOptions)
})

// Print the summary
suite.PrintSummary()
```

## Performance Optimization Strategies

### 1. Parallel Test Execution

The framework supports running tests in parallel to improve overall execution time:

```go
func TestAllExamples(t *testing.T) {
    // This will run all examples and their tests in parallel
    testctx.DiscoverAndRunAllTests(t, "../..")
}
```

### 2. Caching Terraform Operations

You can cache the results of Terraform operations to avoid redundant work:

```go
// Use a cached plan if available
planOutput := terraform.InitAndPlanAndShowWithStruct(t, terraformOptions)

// Later in the test
terraform.ApplyAndIdempotent(t, terraformOptions)
```

### 3. Optimizing Resource Creation

When testing AWS resources, you can optimize resource creation by:

- Using smaller instance types for testing
- Using regions with better availability
- Creating resources in parallel
- Reusing resources across tests when possible

### 4. Minimizing Test Setup and Teardown

Minimize the work done in test setup and teardown:

```go
func TestExample(t *testing.T) {
    // Use a minimal configuration for testing
    ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
        Name: "example1",
        ExtraVars: map[string]interface{}{
            // Only include necessary variables
            "region": "us-west-2",
            "instance_type": "t2.micro",
        },
    })
    
    // Only run necessary assertions
    assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
}
```

### 5. Using Terraform Workspaces

Use Terraform workspaces to isolate tests and avoid conflicts:

```go
// Create a unique workspace for this test
workspaceName := fmt.Sprintf("test-%s", uuid.New().String())
terraform.WorkspaceSelectOrNew(t, terraformOptions, workspaceName)

// Run the test
// ...

// Clean up the workspace
terraform.WorkspaceDelete(t, terraformOptions, workspaceName)
```

## Measuring Performance

To measure the performance of your tests:

```go
import (
    "github.com/caylent-solutions/terraform-test-framework/internal/benchmark"
    "github.com/caylent-solutions/terraform-test-framework/internal/logging"
)

func TestPerformance(t *testing.T) {
    logger := logging.NewWithPrefix(logging.INFO, "PerformanceTest")
    suite := benchmark.NewBenchmarkSuite("terraform-module-test")
    
    // Benchmark example1
    suite.Run("example1", func() error {
        ctx := testctx.RunSingleExample(t, "../../", "example1", testctx.TestConfig{
            Name: "example1",
        })
        return nil
    })
    
    // Benchmark example2
    suite.Run("example2", func() error {
        ctx := testctx.RunSingleExample(t, "../../", "example2", testctx.TestConfig{
            Name: "example2",
        })
        return nil
    })
    
    // Print the summary
    suite.PrintSummary()
    
    // Log the results
    logger.Info("Performance test completed: %s", suite.Summary())
}
```

## Best Practices

1. **Benchmark Regularly**: Run benchmarks regularly to track performance over time.

2. **Compare Benchmarks**: Compare benchmark results before and after changes to identify regressions.

3. **Focus on Bottlenecks**: Focus optimization efforts on the slowest parts of your tests.

4. **Use Realistic Data**: Use realistic data and configurations in your benchmarks.

5. **Isolate Variables**: When benchmarking, isolate variables to understand their impact on performance.

6. **Document Performance Requirements**: Document performance requirements and expectations.

7. **Automate Performance Testing**: Automate performance testing as part of your CI/CD pipeline.