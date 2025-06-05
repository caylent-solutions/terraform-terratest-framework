# Logging Framework

This document describes the logging framework used in the Terraform Test Framework.

## Overview

The framework uses a structured logging approach with different log levels for different types of messages. This makes it easier to control the verbosity of logs and filter messages based on their importance.

## Log Levels

The framework defines the following log levels, in order of increasing severity:

- **DEBUG**: Detailed information, typically useful only for diagnosing problems.
- **INFO**: General information about the progress of the application.
- **WARN**: Warning messages that indicate potential issues.
- **ERROR**: Error messages that indicate problems that prevented an operation from completing.
- **FATAL**: Critical errors that cause the application to terminate.

## Using the Logger

### Creating a Logger

To create a new logger:

```go
import "github.com/caylent-solutions/terraform-terratest-framework/internal/logging"

// Create a logger with INFO level
logger := logging.New(logging.INFO)

// Create a logger with a prefix
logger := logging.NewWithPrefix(logging.INFO, "TestRunner")
```

### Setting Log Level

To set the log level:

```go
// Set log level to DEBUG
logger.SetLevel(logging.DEBUG)

// Set log level from a string
level, err := logging.ParseLogLevel("DEBUG")
if err == nil {
    logger.SetLevel(level)
}
```

### Logging Messages

To log messages:

```go
// Log a debug message
logger.Debug("Processing item %d", itemID)

// Log an info message
logger.Info("Test started for example %s", exampleName)

// Log a warning message
logger.Warn("Deprecated feature used: %s", featureName)

// Log an error message
logger.Error("Failed to run Terraform command: %v", err)

// Log a fatal message and exit
logger.Fatal("Critical error: %v", err)
```

### Using the Default Logger

The framework provides a default logger that can be used without creating a new logger:

```go
import "github.com/caylent-solutions/terraform-terratest-framework/internal/logging"

// Set the default log level
logging.SetDefaultLogLevel(logging.DEBUG)

// Set the default prefix
logging.SetDefaultPrefix("TestRunner")

// Log messages using the default logger
logging.Debug("Debug message")
logging.Info("Info message")
logging.Warn("Warning message")
logging.Error("Error message")
logging.Fatal("Fatal message")
```

## Best Practices

1. **Use Appropriate Log Levels**: Use the appropriate log level for each message based on its importance.

2. **Include Context**: Include enough context in log messages to understand what's happening.

3. **Use Structured Logging**: Use structured logging for complex data structures.

4. **Be Concise**: Keep log messages concise and to the point.

5. **Use Prefixes**: Use prefixes to identify the source of log messages.

6. **Log Errors**: Log errors with appropriate context and log level.

7. **Don't Log Sensitive Information**: Don't log sensitive information such as passwords or access keys.

## Example

Here's an example of using the logging framework:

```go
import (
    "github.com/caylent-solutions/terraform-terratest-framework/internal/logging"
)

func runTest(exampleName string) error {
    logger := logging.NewWithPrefix(logging.INFO, "TestRunner")
    
    logger.Info("Starting test for example %s", exampleName)
    
    // Run test
    if err := runTerraformInit(); err != nil {
        logger.Error("Failed to initialize Terraform: %v", err)
        return err
    }
    
    logger.Debug("Terraform initialized successfully")
    
    if err := runTerraformApply(); err != nil {
        logger.Error("Failed to apply Terraform: %v", err)
        return err
    }
    
    logger.Info("Test completed successfully for example %s", exampleName)
    return nil
}
```