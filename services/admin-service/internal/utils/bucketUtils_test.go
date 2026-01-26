package utils

import (
	"os"
	"testing"
)

func TestGetBucketCount(t *testing.T) {
	os.Setenv("BUCKET_COUNT", "5")
	defer os.Unsetenv("BUCKET_COUNT")

	got, err := GetBucketCount()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestGetBucketCount_MissingEnvVar(t *testing.T) {
	os.Unsetenv("BUCKET_COUNT")

	_, err := GetBucketCount()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetBucketCount_InvalidValue(t *testing.T) {
	os.Setenv("BUCKET_COUNT", "invalid")
	defer os.Unsetenv("BUCKET_COUNT")

	_, err := GetBucketCount()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
