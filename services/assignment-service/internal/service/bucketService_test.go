package service

import "testing"

func TestGetBucketForDeterminism(t *testing.T) {
	svc := NewBucketService("test_salt", 4)

	bucket1 := svc.GetBucketFor("user_123")
	bucket2 := svc.GetBucketFor("user_123")
	bucket3 := svc.GetBucketFor("user_123")

	if bucket1 != bucket2 || bucket2 != bucket3 {
		t.Error("Same user should always get same bucket")
	}
}

func TestSaltAffectsOutput(t *testing.T) {
	svc1 := NewBucketService("salt_one", 4)
	svc2 := NewBucketService("salt_two", 4)

	bucket1 := svc1.GetBucketFor("user_123")
	bucket2 := svc2.GetBucketFor("user_123")

	if bucket1 == bucket2 {
		t.Error("Different salts should produce different buckets")
	}
}

func TestStability(t *testing.T) {
	expected := map[string]int{
		"user_123": 4221,
		"user_456": 5419,
		"user_789": 114,
	}

	svc := NewBucketService("stable_salt", 10000)

	for userID, expectedBucket := range expected {
		actual := svc.GetBucketFor(userID)
		if actual != expectedBucket {
			t.Errorf("Bucket for %s changed: expected %d, got %d",
				userID, expectedBucket, actual)
		}
	}
}
