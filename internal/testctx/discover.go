package testctx

import (
	"os"
	"path/filepath"
	"testing"
)

// DiscoverExamples finds all examples in the given directory
func DiscoverExamples(t *testing.T, moduleRootPath string) []string {
	entries, err := os.ReadDir(moduleRootPath)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	var examples []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Only include directories that start with "example-"
		if filepath.Base(entry.Name())[:8] == "example-" {
			examples = append(examples, entry.Name())
		}
	}

	return examples
}
