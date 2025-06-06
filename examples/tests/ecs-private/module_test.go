package ecs_private

import (
	"testing"

	"github.com/caylent-solutions/terraform-terratest-framework/internal/testctx"
	"github.com/caylent-solutions/terraform-terratest-framework/tests/unit"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Example of an example-specific test
func TestEcsPrivate(t *testing.T) {
	// Run just this specific example
	ctx := testctx.RunSingleExample(t, "../../", "ecs-private", testctx.TestConfig{
		Name: "ecs-private",
		ExtraVars: map[string]interface{}{
			"cluster_name": "private-cluster",
			"region":       "us-west-2",
			"is_private":   true,
		},
	})

	// Example-specific assertions
	clusterName := terraform.Output(t, ctx.Terraform, "cluster_name")
	t.Logf("ECS Cluster Name: %s", clusterName)

	// Verify outputs
	unit.AssertOutputEquals(t, ctx, "cluster_name", "private-cluster")
	unit.AssertOutputEquals(t, ctx, "is_private", "true")

	// Additional example-specific tests would go here
	// verifyPrivateEcsCluster(t, ctx)
}

// Example of an example-specific verification function
func verifyPrivateEcsCluster(t *testing.T, ctx testctx.TestContext) {
	// This would contain AWS SDK code to verify the private ECS cluster
	// For example:
	//
	// clusterName := terraform.Output(t, ctx.Terraform, "cluster_name")
	// ecsClient := ecs.NewFromConfig(...)
	// result, err := ecsClient.DescribeClusters(...)
	// assert.NoError(t, err)
	// assert.Equal(t, "ACTIVE", *result.Clusters[0].Status)
	//
	// // Verify private subnet configuration
	// vpcConfig := terraform.OutputMap(t, ctx.Terraform, "vpc_config")
	// assert.Equal(t, "private", vpcConfig["subnet_type"])
}
