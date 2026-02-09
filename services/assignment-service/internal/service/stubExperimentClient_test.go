package service

import (
	"assignment-service/internal/model"
	"context"
)

type StubExperimentClient struct {
	experiments map[int32][]model.ExperimentWithVariants
}

func (s *StubExperimentClient) GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) ([]model.ExperimentWithVariants, error) {
	if experiments, ok := s.experiments[id]; ok {
		return experiments, nil
	}
	return nil, nil
}

func (s *StubExperimentClient) Close() error {
	return nil
}

func (s *StubExperimentClient) SetExperimentsForBucket(bucketId int32, experiments []model.ExperimentWithVariants) {
	s.experiments[bucketId] = append(s.experiments[bucketId], experiments...)
}

func NewStubExperimentClient() *StubExperimentClient {
	return &StubExperimentClient{
		experiments: map[int32][]model.ExperimentWithVariants{},
	}
}
