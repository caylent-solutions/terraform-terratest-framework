# TestCtx Package Documentation

The `testctx` package is the core of the Terraform Terratest Framework, providing essential functionality for running and managing Terraform tests.

## Overview

The `testctx` package handles:
- Test context management
- Example discovery and execution
- Parallel test execution control
- Idempotency testing

**Important**: Each use of `testctx` spins up a unique Terraform deployment that is tested and then destroyed. This means every call to `RunSingleExample`, `RunAllExamples`, etc. creates separate infrastructure instances for testing.

## Key Components

### TestContext

The `TestContext` struct combines test configuration with Terraform options:

```go
type TestContext struct {
    Config        TestConfig
    Terraform     *terraform.Options
    ExamplePath   string
    Name          string
    TerraformVars map[string]interface{}
}
```

### TestConfig

The `TestConfig` struct holds configuration for a single test:

```go
type TestConfig struct {
    Name      string
    ExtraVars map[string]interface{}
}
```

## Core Functions

### Running Examples

```go
// Run a single example
ctx := testctx.RunSingleExample(t, "../../examples", "example1", testctx.TestConfig{
    Name: "example1-test",
    ExtraVars: map[string]interface{}{
        "region": "us-west-2",
    },
})

// Run all examples
results := testctx.RunAllExamples(t, "../../examples", nil)

// Run all examples with custom tests
testctx.RunAllExamplesWithTests(t, "../../examples", nil, 
    verifyS3Bucket, 
    verifyIAMRoles,
)

// Discover and run all examples with a common test function
testctx.DiscoverAndRunAllTests(t, "../../examples", func(t *testing.T, ctx testctx.TestContext) {
    // Common assertions for all examples
    assertions.AssertOutputNotEmpty(t, ctx, "id")
})
```

### Custom Tests

```go
// Define custom test functions
verifyS3Bucket := func(t *testing.T, ctx testctx.TestContext) {
    // S3 bucket verification logic
}

verifyIAMRoles := func(t *testing.T, ctx testctx.TestContext) {
    // IAM role verification logic
}

// Run all examples
results := testctx.RunAllExamples(t, "../../examples", nil)

// Run custom tests on all examples
testctx.RunCustomTests(t, results, verifyS3Bucket)
testctx.RunCustomTests(t, results, verifyIAMRoles)
```

## Controlling Parallelism

The `testctx` package provides two levels of parallelism control:

1. **Test Fixtures**: Control whether different test fixtures run in parallel
2. **Tests Within Fixtures**: Control whether tests within a single fixture run in parallel

```go
// Check if parallel tests are enabled
isParallel := testctx.IsParallelTestsEnabled()

// Environment variable control
// TERRATEST_DISABLE_PARALLEL_TESTS=true disables parallel tests within fixtures
```

## Idempotency Testing

The package automatically runs idempotency tests for all Terraform examples:

```go
// Check if idempotency testing is enabled
isEnabled := testctx.IdempotencyEnabled()

// Environment variable control
// TERRATEST_IDEMPOTENCY=false disables idempotency testing
```

## Example Usage

### Basic Example

```go
func TestExample(t *testing.T) {
    // Run a single example
    ctx := testctx.RunSingleExample(t, "../../examples", "example1", testctx.TestConfig{
        Name: "example1-test",
    })
    
    // Use assertions
    assertions.AssertOutputEquals(t, ctx, "instance_type", "t2.micro")
    assertions.AssertOutputContains(t, ctx, "bucket_name", "my-bucket")
}
```

### Running Multiple Examples

```go
func TestAllExamples(t *testing.T) {
    // Run all examples with custom configurations
    configs := map[string]testctx.TestConfig{
        "example1": {
            Name: "example1-test",
            ExtraVars: map[string]interface{}{
                "region": "us-west-2",
            },
        },
        "example2": {
            Name: "example2-test",
            ExtraVars: map[string]interface{}{
                "region": "us-east-1",
            },
        },
    }
    
    results := testctx.RunAllExamples(t, "../../examples", configs)
    
    // Run custom tests on all examples
    testctx.RunCustomTests(t, results, func(t *testing.T, ctx testctx.TestContext) {
        // Common assertions for all examples
        assertions.AssertOutputNotEmpty(t, ctx, "id")
    })
}
```

### Automatic Discovery

```go
func TestDiscovery(t *testing.T) {
    // Automatically discover and run all examples
    testctx.DiscoverAndRunAllTests(t, "../../", func(t *testing.T, ctx testctx.TestContext) {
        // Common assertions for all examples
        assertions.AssertOutputNotEmpty(t, ctx, "id")
    })
}
```

## Test Architecture & AI-Assisted Development

### Deployment Lifecycle

Each `testctx` usage follows this lifecycle:
1. **Deploy**: Creates a unique Terraform deployment with generated resource names
2. **Test**: Runs assertions against the deployed infrastructure
3. **Destroy**: Automatically tears down all created resources

This ensures test isolation but comes with infrastructure provisioning overhead.

### AI-Assisted Development Considerations

**Recommended Approach**: Use fewer tests per `testctx` (ideally 1 test per context)

**Benefits**:
- Smaller, focused test output that AI agents can easily parse and understand
- Cleaner error messages when tests fail
- Better context management for AI-assisted debugging
- Easier to identify specific test failures

**Trade-offs**:
- Longer test cycles due to more separate deployments
- Higher infrastructure provisioning overhead
- Serial execution takes significantly longer

**Alternative Approach**: Multiple tests per `testctx` for faster execution

**Benefits**:
- Faster test cycles (single deployment, multiple tests)
- Lower infrastructure overhead
- More efficient resource utilization

**Trade-offs**:
- Complex, interleaved test output that AI agents struggle to parse
- Harder to isolate specific test failures
- Large context sizes that may overwhelm AI agents

### Choosing Your Strategy

```go
// AI-Friendly: One test per context (recommended for AI-assisted development)
func TestVPCCreation(t *testing.T) {
    ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
        Name: "vpc-creation-test",
    })
    assertions.AssertOutputNotEmpty(t, ctx, "vpc_id")
}

func TestVPCTags(t *testing.T) {
    ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
        Name: "vpc-tags-test",
    })
    assertions.AssertOutputContains(t, ctx, "vpc_tags", "Environment")
}

// Efficiency-Focused: Multiple tests per context (faster but less AI-friendly)
func TestVPCComprehensive(t *testing.T) {
    ctx := testctx.RunSingleExample(t, "../../examples", "basic", testctx.TestConfig{
        Name: "vpc-comprehensive-test",
    })
    
    // Multiple tests in single deployment
    assertions.AssertOutputNotEmpty(t, ctx, "vpc_id")
    assertions.AssertOutputContains(t, ctx, "vpc_tags", "Environment")
    assertions.AssertResourceCountExact(t, ctx, "aws_vpc", 1)
    // ... many more tests
}
```

## Best Practices

1. **Use RunSingleExample for Example-Specific Tests**: When writing tests for a specific example, use `RunSingleExample` to focus on that example.

2. **Use RunAllExamples for Common Tests**: When writing tests that should run on all examples, use `RunAllExamples` or `DiscoverAndRunAllTests`.

3. **Control Parallelism**: Use environment variables to control parallelism based on your testing needs.

4. **Clean Up Resources**: The framework automatically cleans up Terraform resources, but if your tests create additional resources, clean them up.

5. **Use Descriptive Test Names**: Set meaningful names in `TestConfig` to make test failures easier to understand.

6. **Consider Your Development Workflow**: Choose between AI-friendly (fewer tests per context) or efficiency-focused (more tests per context) approaches based on your team's needs.