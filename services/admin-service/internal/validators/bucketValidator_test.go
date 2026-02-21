package validators

import (
	"admin-service/internal/problems"
	"testing"
)

func TestValidateBucketId(t *testing.T) {

	tests := []struct {
		name        string
		bucketId    int32
		bucketCount int32
		want        []problems.Violation
	}{
		{
			name:        "Valid bucket ID",
			bucketId:    1,
			bucketCount: 5,
			want:        nil,
		},
		{
			name:        "Negative bucket ID",
			bucketId:    -1,
			bucketCount: 5,
			want: []problems.Violation{
				{
					Field:   "bucket_id",
					Message: "must be non-negative",
				},
			},
		},
		{
			name:        "Bucket ID exceeds count",
			bucketId:    5,
			bucketCount: 5,
			want: []problems.Violation{
				{
					Field:   "bucket_id",
					Message: "exceeds maximum bucket count",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateBucketId(tt.bucketId, tt.bucketCount)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], v)
				}
			}
		})
	}

}
