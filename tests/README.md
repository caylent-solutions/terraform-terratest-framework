# Terraform Terratest Framework Tests

This directory contains tests for the Terraform Terratest Framework itself.

## Directory Structure

```
tests/
├── unit/                      # Unit tests for framework components
│   └── framework_test.go      # Tests for core framework functions
├── functional/                # Functional tests for framework
│   └── framework_test.go      # Tests for framework integration
└── terraform-module-test-fixtures/  # Test fixtures for framework tests
    └── basic/                 # Basic Terraform module for testing
        └── main.tf            # Simple Terraform configuration
```

## Test Categories

### Unit Tests

The `unit/` directory contains tests for individual components of the framework:
- Testing the `testctx` package
- Testing the `assertions` package
- Testing utility functions

### Functional Tests

The `functional/` directory contains tests that verify the framework works as a whole:
- Testing the `RunSingleExample` function
- Testing the `RunAllExamples` function
- Testing the discovery mechanism

### Test Fixtures

The `terraform-module-test-fixtures/` directory contains Terraform modules used for testing the framework:
- The `basic/` fixture creates a simple file and outputs its path and content
- Additional fixtures can be added for more complex testing scenarios

## Running Tests

To run the framework tests:

```bash
go test ./tests/unit/... -v       # Run unit tests
go test ./tests/functional/... -v  # Run functional tests
go test ./tests/... -v            # Run all tests
```