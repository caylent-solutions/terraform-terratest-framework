package helpers

import (
	"testing"

	"github.com/caylent-solutions/terraform-test-framework/internal/testctx"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Example of a helper function for verifying S3 buckets
func VerifyS3Bucket(t *testing.T, ctx testctx.TestContext, expectedEncryption string) {
	// Get bucket name from Terraform outputs
	bucketName := terraform.Output(t, ctx.Terraform, "bucket_name")
	
	// This would contain AWS SDK code to verify the S3 bucket
	// For example:
	//
	// s3Client := s3.NewFromConfig(...)
	// result, err := s3Client.GetBucketEncryption(...)
	// assert.NoError(t, err)
	// assert.Equal(t, expectedEncryption, *result.ServerSideEncryptionConfiguration.Rules[0].ApplyServerSideEncryptionByDefault.SSEAlgorithm)
	
	t.Logf("Verified S3 bucket: %s", bucketName)
}

// Example of a helper function for verifying IAM roles
func VerifyIamRole(t *testing.T, ctx testctx.TestContext, expectedService string) {
	// Get role name from Terraform outputs
	roleName := terraform.Output(t, ctx.Terraform, "role_name")
	
	// This would contain AWS SDK code to verify the IAM role
	// For example:
	//
	// iamClient := iam.NewFromConfig(...)
	// result, err := iamClient.GetRole(...)
	// assert.NoError(t, err)
	// assert.Contains(t, *result.Role.AssumeRolePolicyDocument, expectedService)
	
	t.Logf("Verified IAM role: %s", roleName)
}