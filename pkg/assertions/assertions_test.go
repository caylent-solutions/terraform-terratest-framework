package assertions

import "testing"

// TestCountResourcesOfType exercises the resource-counting logic that backs
// AssertResourceCount, including the cases the previous hardcoded regex
// (`module\.example\.<type>\.`) silently mis-counted as 0: count/for_each
// indexed module instances and resources nested in submodules.
func TestCountResourcesOfType(t *testing.T) {
	cases := []struct {
		name         string
		stateList    string
		resourceType string
		want         int
	}{
		{
			name:         "direct child of module.example",
			stateList:    "module.example.aws_lb_listener.this",
			resourceType: "aws_lb_listener",
			want:         1,
		},
		{
			name:         "for_each indexed resource instances",
			stateList:    "module.example.aws_lb_listener.this[\"http\"]\nmodule.example.aws_lb_listener.this[\"https\"]",
			resourceType: "aws_lb_listener",
			want:         2,
		},
		{
			name:         "count/for_each indexed module instance",
			stateList:    "module.example[\"a\"].aws_lb_listener.this\nmodule.example[\"b\"].aws_lb_listener.this",
			resourceType: "aws_lb_listener",
			want:         2,
		},
		{
			name:         "nested submodule",
			stateList:    "module.example.module.inner.aws_lb_listener.this",
			resourceType: "aws_lb_listener",
			want:         1,
		},
		{
			name:         "type is matched as a full segment (aws_lb must not match aws_lb_listener)",
			stateList:    "module.example.aws_lb_listener.this\nmodule.example.aws_lb.this",
			resourceType: "aws_lb",
			want:         1,
		},
		{
			name:         "root-level fixtures outside module.example are not counted",
			stateList:    "aws_lb.fixture\naws_lb_listener.fixture\nmodule.example.aws_lb_listener.this",
			resourceType: "aws_lb_listener",
			want:         1,
		},
		{
			name:         "same-named data source under module.example is not counted as a resource",
			stateList:    "module.example.data.aws_subnet.selected\nmodule.example.aws_subnet.this",
			resourceType: "aws_subnet",
			want:         1,
		},
		{
			name:         "no matching resources",
			stateList:    "module.example.aws_s3_bucket.this",
			resourceType: "aws_lb_listener",
			want:         0,
		},
		{
			name:         "empty state list",
			stateList:    "",
			resourceType: "aws_lb_listener",
			want:         0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := countResourcesOfType(tc.stateList, tc.resourceType)
			if got != tc.want {
				t.Errorf("countResourcesOfType(resourceType=%q) = %d, want %d", tc.resourceType, got, tc.want)
			}
		})
	}
}
