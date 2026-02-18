package validators

import (
	"admin-service/internal/problems"
)

func ValidateBucketId(bucketId int32, bucketCount int32) []problems.Violation {
	var violations []problems.Violation

	if bucketId < 0 {
		violations = append(violations, problems.Violation{
			Field:   "bucket_id",
			Message: "must be non-negative",
		})
	}

	if bucketId >= bucketCount {
		violations = append(violations, problems.Violation{
			Field:   "bucket_id",
			Message: "exceeds maximum bucket count",
		})
	}

	return violations
}
