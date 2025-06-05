#!/bin/bash

# Load asdf
. ~/.asdf/asdf.sh || echo "Failed to load asdf"

# Set environment variables to reduce memory usage
export GOGC=off

# Function to run a command with a timeout
run_with_timeout() {
  local timeout=$1
  local cmd=$2
  
  # Run the command with timeout
  timeout --preserve-status "$timeout" bash -c "$cmd"
  return $?
}

echo "Running gofmt checks..."
GOFMT_FILES=$(gofmt -l .)
if [ ! -z "$GOFMT_FILES" ]; then
  echo "The following files need formatting:"
  echo "$GOFMT_FILES"
  EXIT_CODE=1
else
  echo "All files are properly formatted."
fi

# Run go vet directly on each package in smaller batches
echo "Running go vet on packages..."
EXIT_CODE=0

# Define packages to check in very small batches
PACKAGES=(
  "./internal/assertions"
  "./internal/benchmark"
  "./internal/errors"
  "./internal/examples"
  "./internal/idempotency"
  "./internal/logging"
  "./internal/testctx"
  "./tftest-cli/cmd"
  "./tftest-cli/logger"
  "./examples/tests/common"
  "./examples/tests/ecs-private"
  "./examples/tests/ecs-public"
  "./examples/tests/helpers"
  "./tests/functional"
  "./tests/unit"
)

for pkg in "${PACKAGES[@]}"; do
  if [ ! -d "$pkg" ]; then
    echo "Skipping $pkg (not a directory)"
    continue
  fi
  
  echo "Checking $pkg with go vet..."
  run_with_timeout "30s" "go vet $pkg"
  
  if [ $? -ne 0 ]; then
    EXIT_CODE=1
  fi
done

# Run a simplified version of errcheck
echo "Running error checking on critical files..."
find ./internal ./tftest-cli -name "*.go" | grep -v "_test.go" | while read file; do
  echo "Checking $file for unchecked errors..."
  grep -n "\\bErr\\b\\|\\berr\\b\\|error" "$file" | grep -v "if.*err\\|switch.*err\\|case.*err\\|for.*err\\|:=.*err\\|=.*err\\|func.*error\\|error)" | grep -v "return.*err\\|return.*Err"
  if [ $? -eq 0 ]; then
    echo "  Potential unchecked errors found in $file"
    EXIT_CODE=1
  fi
done

# Check for common Go issues
echo "Checking for common Go issues..."
find . -name "*.go" | xargs grep -l "fmt.Print" | grep -v "_test.go" && echo "Warning: fmt.Print found in non-test files" && EXIT_CODE=1
find . -name "*.go" | xargs grep -l "log.Print" | grep -v "_test.go" | grep -v "logger.go" && echo "Warning: log.Print found in non-test files" && EXIT_CODE=1

if [ $EXIT_CODE -eq 0 ]; then
  echo "All checks passed!"
else
  echo "Some checks failed."
fi

exit $EXIT_CODE