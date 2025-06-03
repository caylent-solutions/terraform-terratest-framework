# Directory Structure

This document describes the required directory structure for using the Terraform Test Framework.

## Required Directory Structure

The framework is opinionated about directory structure and expects:

```
terraform-module/
├── examples/                # Required: Contains all examples
│   ├── example1/            # Required: Each example in its own directory
│   │   ├── main.tf
│   │   └── ...
│   └── ...
├── tests/                   # Required: Contains all tests
│   ├── common/              # Optional: Tests that run on all examples
│   │   └── common_test.go
│   ├── helpers/             # Optional: Helper functions for tests
│   │   └── helpers.go
│   ├── example1/            # Required: Tests for each example (must match example name)
│   │   └── module_test.go
│   └── ...
└── ...
```

If your module doesn't follow this structure, the tool will not work correctly.

## Framework Repository Structure

The framework repository itself uses this structure:

```
terraform-terratest-framework/
├── examples/                      # Example code for users
│   └── tests/                     # Example test implementations
│       ├── common/                # Example common tests
│       ├── helpers/               # Example helper functions
│       └── ecs-public/            # Example specific tests
├── tests/                         # Tests for the framework itself
│   ├── unit/                      # Unit tests for framework components
│   ├── functional/                # Functional tests for framework
│   └── terraform-module-test-fixtures/  # Test fixtures for framework tests
└── ...
```

## Test Organization

The framework supports a structured approach to testing:

1. **Example-Specific Tests**: Each example has its own test directory with the same name
2. **Common Tests**: Tests in the `common` directory run on all examples
3. **Helper Functions**: Reusable test helpers in the `helpers` directory