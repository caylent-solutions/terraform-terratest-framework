# Terraform Testing Guide

This guide explains how to use the Terraform Test Framework to test your Terraform modules.

## Directory Structure

For details on the required directory structure, see the [Directory Structure Documentation](DIRECTORY_STRUCTURE.md).

## Test Organization and Writing Tests

For details on test organization and writing tests, see the [Writing Tests Documentation](WRITING_TESTS.md).

## The TestCtx Package

The `testctx` package is the core of the framework, providing essential functionality for running and managing Terraform tests. For detailed documentation, see the [TestCtx Package Documentation](TESTCTX_PACKAGE.md).

Key functions include:

```go
// Run a single example
ctx := testctx.RunSingleExample(t, "../../examples", "example1", testctx.TestConfig{
    Name: "example1-test",
})

// Run all examples
results := testctx.RunAllExamples(t, "../../examples", nil)

// Run custom tests on all examples
testctx.RunCustomTests(t, results, verifyS3Bucket)
```

## Test Fixtures

The framework includes test fixtures that demonstrate how to use the framework:

```
tests/terraform-module-test-fixtures/
├── example-basic/           # Basic example with a simple file output
│   └── main.tf
└── example-advanced/        # Advanced example with additional outputs and tags
    └── main.tf
```

These fixtures are used in the functional tests to verify that the framework works correctly.

## Example Tests

The framework includes example tests that demonstrate how to write tests for Terraform modules:

```
examples/tests/
├── common/                  # Common tests that run on all examples
│   └── common_test.go
├── helpers/                 # Helper functions for tests
│   └── helpers.go
├── ecs-public/              # Tests for the ecs-public example
│   └── module_test.go
└── ecs-private/             # Tests for the ecs-private example
    └── module_test.go
```

These examples show how to:
- Write tests for specific examples
- Write common tests that run on all examples
- Create helper functions for reuse across tests
- Test different configurations of the same module

## Running Tests

For details on running tests with the TFTest CLI, see the [CLI Usage Documentation](CLI_USAGE.md).

## Environment Variables

For details on required environment variables, see the [CLI Usage Documentation](CLI_USAGE.md#required-environment-variables).

## Best Practices

1. **Organize Tests by Example**: Keep tests for each example in their own directory

2. **Use Common Tests for Shared Requirements**: Put tests that should run on all examples in the `common` directory

3. **Create Helper Functions**: Extract reusable test logic into helper functions in the `helpers` directory

4. **Use the TestCtx Package Effectively**: Leverage the functions provided by the `testctx` package to simplify test setup and execution

5. **Clean Up Resources**: The framework automatically cleans up Terraform resources, but if your tests create additional resources, clean them up

6. **Use Descriptive Test Names**: Name tests based on what they're verifying

7. **Error Handling**: Include proper error handling and descriptive assertion messages

8. **AWS Authentication**: Ensure AWS credentials are properly configured before running tests

9. **Run Standard Tests**: Use the standard tests provided by the framework to ensure your module follows best practices