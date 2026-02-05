package service

import (
	"context"
	"encoding/json"
	"fmt"
)

type StubAssignmentCache struct {
	cache map[string]map[string]string
}

func (s *StubAssignmentCache) ActionFullInvalidateBucket(ctx context.Context, bucketId int32) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)
	delete(s.cache, bucketIdStr)
	return nil
}

func (s *StubAssignmentCache) ActionUpdateValueForBucketAndFlag(ctx context.Context, bucketId int32, flag, newValue string) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)
	if assignments, exists := s.cache[bucketIdStr]; exists {
		assignments[flag] = newValue
		s.cache[bucketIdStr] = assignments
	} else {
		s.cache[bucketIdStr] = map[string]string{flag: newValue}
	}
	return nil
}

func (s *StubAssignmentCache) ActionRemoveFlagFromBucket(ctx context.Context, bucketId int32, flag string) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)
	if assignments, exists := s.cache[bucketIdStr]; exists {
		delete(assignments, flag)
		s.cache[bucketIdStr] = assignments
	}
	return nil
}

func NewStubAssignmentCache() *StubAssignmentCache {
	return &StubAssignmentCache{
		cache: make(map[string]map[string]string),
	}
}

func (s *StubAssignmentCache) SetAssignmentsForBucket(ctx context.Context, bucketId int32, assignments map[string]string) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)
	s.cache[bucketIdStr] = assignments
	return nil
}

func (s *StubAssignmentCache) GetAssignmentsForBucket(ctx context.Context, bucketId int32) (map[string]string, error) {
	bucketIdStr := fmt.Sprintf("%d", bucketId)
	if assignments, exists := s.cache[bucketIdStr]; exists {
		return assignments, nil
	}
	return nil, nil
}

func (s *StubAssignmentCache) ClearCache() {
	s.cache = make(map[string]map[string]string)
}

type StubAssignmentClient struct {
	assignments map[int32]map[string]string
}

func NewStubAssignmentClient() *StubAssignmentClient {
	return &StubAssignmentClient{
		assignments: make(map[int32]map[string]string),
	}
}

func (s *StubAssignmentClient) GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) (map[string]string, error) {
	if assignments, exists := s.assignments[id]; exists {
		return assignments, nil
	}
	return nil, nil
}

func (s *StubAssignmentClient) SetAssignmentsForBucket(id int32, assignments map[string]string) {
	s.assignments[id] = assignments
}

func (s *StubAssignmentClient) ClearAssignments() {
	s.assignments = make(map[int32]map[string]string)
}

func (s *StubAssignmentClient) Close() error {
	return nil
}

type MockKafkaConsumer struct {
	Messages []KafkaMessage
}

func (m *MockKafkaConsumer) PollMessages(ctx context.Context) ([]KafkaMessage, error) {
	return m.Messages, nil
}

func (m *MockKafkaConsumer) AddMessage(msg interface{}) {
	serialized, _ := json.Marshal(msg)
	m.Messages = append(m.Messages, KafkaMessage{Value: serialized})
}

func NewMockKafkaConsumer() *MockKafkaConsumer {
	return &MockKafkaConsumer{
		Messages: []KafkaMessage{},
	}
}
