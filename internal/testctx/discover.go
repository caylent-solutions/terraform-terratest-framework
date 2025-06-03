package testctx

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// DiscoverAndRunAllTests finds all examples and their corresponding test directories,
// then runs all tests in parallel
func DiscoverAndRunAllTests(t *testing.T, moduleRootPath string) map[string]TestContext {
	// Find all examples
	examplesPath := filepath.Join(moduleRootPath, "examples")
	entries, err := os.ReadDir(examplesPath)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}
	
	var wg sync.WaitGroup
	results := make(map[string]TestContext)
	resultsMutex := sync.Mutex{}
	
	// For each example, run its tests in parallel
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		
		exampleName := entry.Name()
		examplePath := filepath.Join(examplesPath, exampleName)
		
		// Check if there's a corresponding test directory
		testPath := filepath.Join(moduleRootPath, "tests", exampleName)
		if _, err := os.Stat(testPath); os.IsNotExist(err) {
			t.Logf("No test directory found for example %s, skipping", exampleName)
			continue
		}
		
		wg.Add(1)
		go func(name, path string) {
			defer wg.Done()
			
			// Run the test in a subtest to isolate failures
			t.Run(fmt.Sprintf("Example_%s", name), func(t *testing.T) {
				// Create default config
				config := TestConfig{
					Name:      name,
					ExtraVars: map[string]interface{}{},
				}
				
				// Run the example
				ctx := RunExample(t, path, config)
				
				// Store the result
				resultsMutex.Lock()
				results[name] = ctx
				resultsMutex.Unlock()
				
				// Clean up resources
				defer terraform.Destroy(t, ctx.Terraform)
				
				// Run common tests if they exist
				commonTestPath := filepath.Join(moduleRootPath, "tests", "common")
				if _, err := os.Stat(commonTestPath); !os.IsNotExist(err) {
					t.Logf("Running common tests for example %s", name)
					// Common tests would be implemented by the user in their test files
					// and would be discovered and run by the Go test runner
				}
			})
		}(exampleName, examplePath)
	}
	
	wg.Wait()
	return results
}