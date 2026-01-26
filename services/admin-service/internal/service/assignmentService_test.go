package service

import (
	errors2 "admin-service/internal/errors"
	"context"
	"errors"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type mockExperimentRepository struct {
	getFunc func(ctx context.Context, bucketId int32) ([]*model.ExperimentVariant, error)
}

func (m *mockExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentVariant, error) {
	return m.getFunc(ctx, bucketId)
}

func TestGetExperimentsAndVariantsForBucket_RepoErrorsShouldPropagate(t *testing.T) {
	mockRepo := &mockExperimentRepository{
		getFunc: func(ctx context.Context, bucketId int32) ([]*model.ExperimentVariant, error) {
			return nil, errors.New("repo error")
		},
	}
	service := &AssignmentService{
		experimentRepository: mockRepo,
		bucketCount:          100,
	}

	_, err := service.GetExperimentsAndVariantsForBucket(context.Background(), 10)
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
		name      string
		bucketId  int32
		wantField string
		wantErr   bool
	}{
		{
			name:      "Negative bucketId returns ValidationError",
			bucketId:  -1,
			wantField: "bucket_id",
			wantErr:   true,
		},
		{
			name:      "bucketId exceeds bucketCount returns ValidationError",
			bucketId:  101,
			wantField: "bucket_id",
			wantErr:   true,
		},
	}

	service := GetNewAssignmentServiceForTest(100)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetExperimentsAndVariantsForBucket(context.Background(), tt.bucketId)
			var ve *errors2.ValidationError
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if errors.As(err, &ve) {
					if ve.Field != tt.wantField {
						t.Errorf("Expected validation error on field '%s', got '%s'", tt.wantField, ve.Field)
					}
				} else {
					t.Errorf("Expected ValidationError, got %v", err)
				}
			} else if err != nil {
				t.Errorf("Did not expect error, got %v", err)
			}
		})
	}
}
