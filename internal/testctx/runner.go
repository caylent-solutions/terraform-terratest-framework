package testctx

import (
	"fmt"
	"os"
	"path/filepath"
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
		Config:     config,
		Terraform:  tfOptions,
		ExamplePath: path,
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
		planResult := terraform.Plan(t, ctx.Terraform)
		if planResult != "" {
			t.Fatalf("Idempotency test failed: Terraform plan would make changes: %s", planResult)
		}
		t.Log("Idempotency test passed")
	} else {
		t.Log("Idempotency testing disabled via TERRATEST_IDEMPOTENCY=false")
	}
	
	return ctx
}

// RunAllExamples runs all examples in the examples directory in parallel
// If configs is nil or empty, it will generate default configs for all examples
func RunAllExamples(t *testing.T, moduleRootPath string, configs map[string]TestConfig) map[string]TestContext {
	examplesPath := filepath.Join(moduleRootPath, "examples")
	entries, err := os.ReadDir(examplesPath)
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
			configs[entry.Name()] = TestConfig{
				Name:      entry.Name(),
				ExtraVars: map[string]interface{}{},
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
		examplePath := filepath.Join(examplesPath, exampleName)
		
		// Skip if no config provided for this example
		config, exists := configs[exampleName]
		if !exists {
			t.Logf("Skipping example %s: no config provided", exampleName)
			continue
		}
		
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
				
				// Clean up resources
				defer terraform.Destroy(t, ctx.Terraform)
			})
		}(exampleName, examplePath, config)
	}
	
	wg.Wait()
	return results
}

// RunSingleExample runs a specific example from the examples directory
// This is useful when tests are organized by example (one test folder per example)
func RunSingleExample(t *testing.T, moduleRootPath string, exampleName string, config TestConfig) TestContext {
	examplePath := filepath.Join(moduleRootPath, "examples", exampleName)
	
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
	
	// Clean up resources when test completes
	t.Cleanup(func() {
		terraform.Destroy(t, ctx.Terraform)
	})
	
	return ctx
}
