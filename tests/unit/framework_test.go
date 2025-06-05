// Copyright 2023 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unit

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/stretchr/testify/assert"
)

// Test the IdempotencyEnabled function
func TestIdempotencyEnabled(t *testing.T) {
	// Test default behavior
	assert.True(t, testctx.IdempotencyEnabled(), "Idempotency should be enabled by default")
}

// Test the TestConfig structure
func TestTestConfig(t *testing.T) {
	config := testctx.TestConfig{
		Name: "test",
		ExtraVars: map[string]interface{}{
			"key": "value",
		},
	}
	
	assert.Equal(t, "test", config.Name)
	assert.Equal(t, "value", config.ExtraVars["key"])
}

// Test the assertions package
func TestAssertions(t *testing.T) {
	// This is just a placeholder for actual unit tests
	// Real tests would mock the terraform.Output function and test each assertion
	t.Log("Assertions package unit tests would go here")
}