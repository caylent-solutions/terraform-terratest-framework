package testctx

import (
	"fmt"
	"testing"
)

// CustomTestFunc defines a function type for custom tests that can be run on each example
type CustomTestFunc func(t *testing.T, ctx TestContext)

// RunCustomTests runs the provided custom test functions on all examples
func RunCustomTests(t *testing.T, contexts map[string]TestContext, testFuncs ...CustomTestFunc) {
	for name, ctx := range contexts {
		for i, testFunc := range testFuncs {
			testName := ""
			if len(testFuncs) > 1 {
				testName = fmt.Sprintf("CustomTest_%d_%s", i+1, name)
			} else {
				testName = fmt.Sprintf("CustomTest_%s", name)
			}
			
			t.Run(testName, func(t *testing.T) {
				testFunc(t, ctx)
			})
		}
	}
}

// RunAllExamplesWithTests runs all examples in parallel and then executes custom tests on each
func RunAllExamplesWithTests(t *testing.T, moduleRootPath string, configs map[string]TestConfig, testFuncs ...CustomTestFunc) map[string]TestContext {
	results := RunAllExamples(t, moduleRootPath, configs)
	RunCustomTests(t, results, testFuncs...)
	return results
}