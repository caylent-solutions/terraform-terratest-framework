package ecs_public

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/assertions"
	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Example of an example-specific test
func TestEcsPublic(t *testing.T) {
	// Run just this specific example
	ctx := testctx.RunSingleExample(t, "../../", "ecs-public", testctx.TestConfig{
		Name: "ecs-public",
		ExtraVars: map[string]interface{}{
			"cluster_name": "test-cluster",
			"region":       "us-west-2",
		},
	})

	// Example-specific assertions
	clusterName := terraform.Output(t, ctx.Terraform, "cluster_name")
	t.Logf("ECS Cluster Name: %s", clusterName)

	// Verify outputs
	assertions.AssertOutputEquals(t, ctx, "cluster_name", "test-cluster")

	// Additional example-specific tests would go here
	// verifyEcsCluster(t, ctx)
}

// Example of an example-specific verification function
func verifyEcsCluster(t *testing.T, ctx testctx.TestContext) {
	// This would contain AWS SDK code to verify the ECS cluster
	// For example:
	//
	// clusterName := terraform.Output(t, ctx.Terraform, "cluster_name")
	// ecsClient := ecs.NewFromConfig(...)
	// result, err := ecsClient.DescribeClusters(...)
	// assert.NoError(t, err)
	// assert.Equal(t, "ACTIVE", *result.Clusters[0].Status)
}
