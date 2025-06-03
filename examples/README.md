# Terraform Terratest Framework Examples

This directory contains example code showing how to use the Terraform Terratest Framework.

## Directory Structure

```
examples/
└── tests/                     # Example test implementations
    ├── common/                # Example common tests
    │   └── common_test.go     # Tests that run on all examples
    ├── helpers/               # Example helper functions
    │   └── aws_helpers.go     # Helper functions for AWS resources
    └── ecs-public/            # Example-specific tests
        └── module_test.go     # Tests for the ecs-public example
```

## Example Tests

### Common Tests

The `tests/common/common_test.go` file shows how to:
- Run all examples in parallel
- Write tests that run on all examples
- Organize common test functions

### Helper Functions

The `tests/helpers/aws_helpers.go` file shows how to:
- Create reusable helper functions
- Verify AWS resources
- Share test logic across multiple test files

### Example-Specific Tests

The `tests/ecs-public/module_test.go` file shows how to:
- Run a specific example
- Configure example-specific variables
- Write tests for a specific example
- Use helper functions

## Using These Examples

These examples demonstrate the recommended structure and patterns for testing Terraform modules with the framework. You can use them as templates for your own tests.

To see how these examples would be used in a real Terraform module, refer to the [Testing Guide](../docs/TESTING_GUIDE.md).