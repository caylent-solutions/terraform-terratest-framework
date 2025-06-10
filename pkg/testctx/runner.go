package testctx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// InitTerraform creates terraform options for the given path and config
func InitTerraform(path string, config TestConfig) *terraform.Options {
	return &terraform.Options{
		TerraformDir: path,
		Vars:         config.ExtraVars,
	}
}

// Run initializes a test context for a single example
func Run(path string, config TestConfig) TestContext {
	tfOptions := InitTerraform(path, config)
	return TestContext{
		Config:      config,
		Terraform:   tfOptions,
		ExamplePath: path,
		Name:        config.Name,
	}
}

// RunExample runs a single terraform example with the given config
// and automatically performs an idempotency test unless disabled via TERRATEST_IDEMPOTENCY=false
func RunExample(t *testing.T, examplePath string, config TestConfig) TestContext {
	ctx := Run(examplePath, config)
	terraform.InitAndApply(t, ctx.Terraform)

	// Run idempotency test by default unless explicitly disabled
	if IdempotencyEnabled() {
		t.Log("Running idempotency test...")
		planOutput := terraform.Plan(t, ctx.Terraform)
		// Check if the plan output contains "No changes" or "no changes"
		if strings.Contains(planOutput, "No changes") || strings.Contains(planOutput, "no changes") {
			t.Log("Idempotency test passed")
		} else {
			t.Fatalf("Idempotency test failed: Terraform plan would make changes: %s", planOutput)
		}
	} else {
		t.Log("Idempotency testing disabled via TERRATEST_IDEMPOTENCY=false")
	}

	// Register cleanup to ensure resources are destroyed
	t.Cleanup(func() {
		terraform.Destroy(t, ctx.Terraform)
	})

	return ctx
}

// RunCustomTests runs a custom test function on all examples in the results map
func RunCustomTests(t *testing.T, results map[string]TestContext, testFunc func(t *testing.T, ctx TestContext)) {
	for _, ctx := range results {
		testFunc(t, ctx)
	}
}

// RunAllExamplesWithTests runs all examples and then runs multiple custom test functions on each example
func RunAllExamplesWithTests(t *testing.T, moduleRootPath string, configs map[string]TestConfig, testFuncs ...func(t *testing.T, ctx TestContext)) map[string]TestContext {
	// Run all examples
	results := RunAllExamples(t, moduleRootPath, configs)

	// Run each test function on all examples
	for _, testFunc := range testFuncs {
		RunCustomTests(t, results, testFunc)
	}

	return results
}

// DiscoverAndRunAllTests runs all examples in the examples directory and executes a custom test function on each
func DiscoverAndRunAllTests(t *testing.T, moduleRootPath string, testFunc func(t *testing.T, ctx TestContext)) map[string]TestContext {
	// Find all example directories
	entries, err := os.ReadDir(moduleRootPath)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	// Create default configs for all examples
	configs := make(map[string]TestConfig)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Only include directories that start with "example-"
		name := entry.Name()
		if len(name) >= 8 && name[:8] == "example-" {
			configs[name] = TestConfig{
				Name:      name,
				ExtraVars: map[string]interface{}{},
			}
		}
	}

	// Run all examples
	results := RunAllExamples(t, moduleRootPath, configs)

	// If a test function is provided, run it on each example
	if testFunc != nil {
		for _, ctx := range results {
			testFunc(t, ctx)
		}
	}

	return results
}

// RunAllExamples runs all examples in the examples directory
// If configs is nil or empty, it will generate default configs for all examples
// Parallelism is controlled by the TERRATEST_DISABLE_PARALLEL_TESTS environment variable
func RunAllExamples(t *testing.T, moduleRootPath string, configs map[string]TestConfig) map[string]TestContext {
	entries, err := os.ReadDir(moduleRootPath)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	// If no configs provided, create default configs for all examples
	if configs == nil || len(configs) == 0 {
		configs = make(map[string]TestConfig)
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// Only include directories that start with "example-"
			name := entry.Name()
			if len(name) >= 8 && name[:8] == "example-" {
				configs[name] = TestConfig{
					Name:      name,
					ExtraVars: map[string]interface{}{},
				}
			}
		}
	}

	var wg sync.WaitGroup
	results := make(map[string]TestContext)
	resultsMutex := sync.Mutex{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		exampleName := entry.Name()

		// Skip if not an example directory
		if len(exampleName) < 8 || exampleName[:8] != "example-" {
			continue
		}

		examplePath := filepath.Join(moduleRootPath, exampleName)

		// Skip if no config provided for this example
		config, exists := configs[exampleName]
		if !exists {
			t.Logf("Skipping example %s: no config provided", exampleName)
			continue
		}

		// Run tests in parallel or sequentially based on environment variable
		if IsParallelTestsEnabled() {
			wg.Add(1)
			go func(name, path string, cfg TestConfig) {
				defer wg.Done()

				// Run the test in a subtest to isolate failures
				t.Run(fmt.Sprintf("Example_%s", name), func(t *testing.T) {
					ctx := RunExample(t, path, cfg)

					// Store the result
					resultsMutex.Lock()
					results[name] = ctx
					resultsMutex.Unlock()

					// Note: RunExample now registers its own cleanup function
					// so we don't need to call terraform.Destroy here
				})
			}(exampleName, examplePath, config)
		} else {
			// Run sequentially
			t.Run(fmt.Sprintf("Example_%s", exampleName), func(t *testing.T) {
				ctx := RunExample(t, examplePath, config)

				// Store the result
				resultsMutex.Lock()
				results[exampleName] = ctx
				resultsMutex.Unlock()
			})
		}
	}

	wg.Wait()
	return results
}

// RunSingleExample runs a specific example from the examples directory
// This is useful when tests are organized by example (one test folder per example)
func RunSingleExample(t *testing.T, moduleRootPath string, exampleName string, config TestConfig) TestContext {
	var examplePath string

	// Check if exampleName is "." which means use moduleRootPath directly
	if exampleName == "." {
		examplePath = moduleRootPath
	} else {
		examplePath = filepath.Join(moduleRootPath, exampleName)
	}

	// Check if example exists
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Fatalf("Example %s not found at path %s", exampleName, examplePath)
	}

	// If no config provided, create a default one
	if config.Name == "" {
		config = TestConfig{
			Name:      exampleName,
			ExtraVars: map[string]interface{}{},
		}
	}

	// Run the example
	ctx := RunExample(t, examplePath, config)

	// Note: RunExample now registers its own cleanup function
	// so we don't need to call terraform.Destroy here

	return ctx
}
