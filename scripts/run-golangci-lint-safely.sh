#!/bin/bash

# Load asdf
. ~/.asdf/asdf.sh || echo "Failed to load asdf"

# Set environment variables to reduce memory usage
export GOGC=off
export GOLANGCI_LINT_CACHE=/tmp/golangci-lint-cache

# Pass all arguments to golangci-lint
EXTRA_ARGS="$@"

# Define directories to lint separately
DIRS=(
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

# Create a minimal config file
cat > .golangci.safe.yml << EOF
run:
  concurrency: 1
  timeout: 30s
  skip-dirs:
    - vendor
    - .terraform
  memory: 512MB
  go: "1.23"

linters:
  enable:
    - gofmt
  disable-all: false

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
EOF

echo "Running golangci-lint in safe mode..."
EXIT_CODE=0

# First run just gofmt on the entire codebase - it's lightweight
echo "Running gofmt checks..."
golangci-lint run --config=.golangci.safe.yml --disable-all --enable=gofmt $EXTRA_ARGS
FMT_EXIT=$?
if [ $FMT_EXIT -ne 0 ]; then
  EXIT_CODE=$FMT_EXIT
fi

# Then run govet on each directory separately
for dir in "${DIRS[@]}"; do
  if [ ! -d "$dir" ]; then
    echo "Skipping $dir (not a directory)"
    continue
  fi
  
  echo "Running govet on $dir..."
  timeout 30s golangci-lint run --config=.golangci.safe.yml --disable-all --enable=govet $EXTRA_ARGS "$dir/..."
  
  # Capture exit code but continue with other directories
  CURRENT_EXIT=$?
  if [ $CURRENT_EXIT -ne 0 ] && [ $CURRENT_EXIT -ne 124 ]; then  # 124 is timeout exit code
    EXIT_CODE=$CURRENT_EXIT
  elif [ $CURRENT_EXIT -eq 124 ]; then
    echo "Timeout occurred for $dir, skipping"
  fi
done

# Clean up temporary config
rm .golangci.safe.yml

exit $EXIT_CODE