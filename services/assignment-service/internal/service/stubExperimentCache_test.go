package service

import (
	"assignment-service/internal/model"
	"context"
	"fmt"
)

type StubExperimentCache struct {
	cache map[string]interface{}
}

func NewStubExperimentCache() *StubExperimentCache {
	return &StubExperimentCache{
		cache: make(map[string]interface{}),
	}
}

func (s *StubExperimentCache) buildKeyForBucket(bucketId int32) string {
	return fmt.Sprintf("bucket:%d:experimentKeys", bucketId)
}

func (s *StubExperimentCache) buildKeyForExperiment(bucketId string) string {
	return fmt.Sprintf("experiment:%s", bucketId)
}

func (s *StubExperimentCache) GetBucketExperimentKeys(ctx context.Context, bucketId int32) ([]string, error) {
	key := s.buildKeyForBucket(bucketId)
	if keys, exists := s.cache[key]; exists {
		return keys.([]string), nil
	}
	return nil, nil
}

func (s *StubExperimentCache) AddBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error {
	key := s.buildKeyForBucket(bucketId)
	if keys, exists := s.cache[key]; exists {
		keys = append(keys.([]string), experimentKey)
		s.cache[key] = keys
	} else {
		// If the bucket doesn't exist, create a new entry
		s.cache[key] = []string{experimentKey}
	}

	return nil
}

func (s *StubExperimentCache) RemoveBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error {
	key := s.buildKeyForBucket(bucketId)
	if keys, exists := s.cache[key]; exists {
		experimentKeys := keys.([]string)
		for i, key := range experimentKeys {
			if key == experimentKey {
				// split the slice to remove the experimentKey
				experimentKeys = append(experimentKeys[:i], experimentKeys[i+1:]...)
				s.cache[key] = experimentKeys
				break
			}
		}
	}

	return nil
}

func (s *StubExperimentCache) InvalidateBucket(ctx context.Context, bucketId int32) error {
	key := s.buildKeyForBucket(bucketId)
	delete(s.cache, key)
	return nil
}

func (s *StubExperimentCache) GetExperiment(ctx context.Context, experimentKey string) (*model.ExperimentWithVariants, error) {
	key := s.buildKeyForExperiment(experimentKey)
	if experiment, exists := s.cache[key]; exists {
		return experiment.(*model.ExperimentWithVariants), nil
	}
	return nil, fmt.Errorf("experiment not found")
}

func (s *StubExperimentCache) SetExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error {
	key := s.buildKeyForExperiment(experimentKey)
	s.cache[key] = experiment
	return nil
}

func (s *StubExperimentCache) UpdateExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error {
	key := s.buildKeyForExperiment(experimentKey)
	if _, exists := s.cache[key]; exists {
		s.cache[key] = experiment
		return nil
	}
	return fmt.Errorf("experiment not found for update")
}

func (s *StubExperimentCache) InvalidateExperiment(ctx context.Context, experimentKey string) error {
	key := s.buildKeyForExperiment(experimentKey)
	delete(s.cache, key)
	return nil
}

func (s *StubExperimentCache) GetExperimentsForBucket(ctx context.Context, bucketId int32) ([]model.ExperimentWithVariants, error) {
	key := s.buildKeyForBucket(bucketId)
	if keys, exists := s.cache[key]; exists {
		experimentKeys := keys.([]string)
		var experiments []model.ExperimentWithVariants
		for _, experimentKey := range experimentKeys {
			experiment, err := s.GetExperiment(ctx, experimentKey)
			if err != nil {
				return nil, err
			}
			experiments = append(experiments, *experiment)
		}
		return experiments, nil
	}
	return nil, nil
}

func (s *StubExperimentCache) SetExperimentsForBucket(ctx context.Context, bucketId int32, experiments []model.ExperimentWithVariants) error {
	for _, experiment := range experiments {
		err := s.SetExperiment(ctx, experiment.ExperimentKey, &experiment)
		if err != nil {
			return err
		}

		err = s.AddBucketExperimentKey(ctx, bucketId, experiment.ExperimentKey)
		if err != nil {
			return err
		}
	}
	return nil
}
