package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type mockExperimentRepository struct {
	getFunc func(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error)
}

func (m *mockExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error) {
	return m.getFunc(ctx, bucketId)
}

func TestGetExperimentsAndVariantsForBucket_RepoErrorsShouldPropagate(t *testing.T) {
	mockRepo := &mockExperimentRepository{
		getFunc: func(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error) {
			return nil, errors.New("repo error")
		},
	}
	service := &AssignmentService{
		experimentRepository: mockRepo,
		bucketCount:          100,
	}

	_, _, err := service.GetExperimentsAndVariantsForBucket(context.Background(), 10)
	if err == nil {
		t.Errorf("Expected error, got nil")
	} else if err.Error() != "repo error" {
		t.Errorf("Expected 'repo error', got %v", err)
	}
}

func GetNewAssignmentServiceForTest(bucketCount int32) *AssignmentService {
	return &AssignmentService{
		bucketCount: bucketCount,
	}
}

func TestGetExperimentsAndVariantsForBucket_Validation(t *testing.T) {
	tests := []struct {
		name          string
		bucketId      int32
		wantField     string
		wantViolation bool
	}{
		{
			name:          "Negative bucketId returns validation violation",
			bucketId:      -1,
			wantField:     "bucket_id",
			wantViolation: true,
		},
		{
			name:          "bucketId exceeds bucketCount returns validation violation",
			bucketId:      101,
			wantField:     "bucket_id",
			wantViolation: true,
		},
	}

	service := GetNewAssignmentServiceForTest(100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, violations, err := service.GetExperimentsAndVariantsForBucket(context.Background(), tt.bucketId)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if tt.wantViolation {
				if len(violations) == 0 {
					t.Fatal("Expected violations, got none")
				}
				found := false
				for _, v := range violations {
					if v.Field == tt.wantField {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected violation on field %q, got %v", tt.wantField, violations)
				}
			}
		})
	}
}
