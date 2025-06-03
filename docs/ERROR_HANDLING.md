# Error Handling Strategy

This document describes the error handling strategy used in the Terraform Test Framework.

## Overview

The framework uses a structured approach to error handling, with different error types for different categories of errors. This makes it easier to understand and handle errors appropriately.

## Error Types

The framework defines the following error types:

- **ConfigError**: Errors related to configuration, such as invalid configuration values or missing required configuration.
- **ValidationError**: Errors related to validation, such as invalid input values or failed validation checks.
- **TerraformError**: Errors from Terraform operations, such as failed Terraform commands or invalid Terraform state.
- **AssertionError**: Errors from test assertions, such as failed assertions or invalid assertion parameters.
- **InternalError**: Internal framework errors, such as unexpected conditions or internal inconsistencies.

## Creating Errors

To create a new error, use the appropriate constructor function:

```go
import "github.com/caylent-solutions/terraform-test-framework/internal/errors"

// Create a configuration error
err := errors.NewConfigError("Invalid configuration value", cause)

// Create a validation error
err := errors.NewValidationError("Invalid input value", cause)

// Create a Terraform error
err := errors.NewTerraformError("Terraform command failed", cause)

// Create an assertion error
err := errors.NewAssertionError("Assertion failed", cause)

// Create an internal error
err := errors.NewInternalError("Unexpected condition", cause)
```

## Handling Errors

When handling errors, you can check the error type using type assertions:

```go
import "github.com/caylent-solutions/terraform-test-framework/internal/errors"

if err != nil {
    if configErr, ok := err.(*errors.FrameworkError); ok && configErr.Type == errors.ConfigError {
        // Handle configuration error
    } else if terraformErr, ok := err.(*errors.FrameworkError); ok && terraformErr.Type == errors.TerraformError {
        // Handle Terraform error
    } else {
        // Handle other errors
    }
}
```

## Logging Errors

When logging errors, include the error message and any relevant context:

```go
import (
    "github.com/caylent-solutions/terraform-test-framework/internal/errors"
    "github.com/caylent-solutions/terraform-test-framework/internal/logging"
)

if err != nil {
    logging.Error("Failed to run Terraform command: %v", err)
}
```

## Best Practices

1. **Use Specific Error Types**: Use the most specific error type that applies to the situation.

2. **Include Context**: Include enough context in error messages to understand what went wrong.

3. **Chain Errors**: When wrapping errors, use the `cause` parameter to chain errors together.

4. **Handle Errors Appropriately**: Handle errors at the appropriate level, and don't swallow errors without good reason.

5. **Log Errors**: Log errors with appropriate context and log level.

6. **Return Errors**: Return errors to the caller when appropriate, rather than handling them locally.

7. **Don't Panic**: Use errors rather than panics for expected error conditions.

## Example

Here's an example of using the error handling strategy:

```go
import (
    "github.com/caylent-solutions/terraform-test-framework/internal/errors"
    "github.com/caylent-solutions/terraform-test-framework/internal/logging"
)

func runTerraformCommand(command string) error {
    // Validate input
    if command == "" {
        return errors.NewValidationError("Command cannot be empty", nil)
    }

    // Run command
    output, err := exec.Command("terraform", command).CombinedOutput()
    if err != nil {
        return errors.NewTerraformError(
            fmt.Sprintf("Terraform command '%s' failed: %s", command, string(output)),
            err,
        )
    }

    return nil
}

func main() {
    err := runTerraformCommand("apply")
    if err != nil {
        logging.Error("Failed to apply Terraform: %v", err)
        os.Exit(1)
    }
}
```