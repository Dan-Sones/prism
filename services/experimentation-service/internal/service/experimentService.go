package service

import (
	"context"
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"
)

type ExperimentService struct {
	experimentRepository *repository.ExperimentRepository
	logger               *slog.Logger
}

func NewExperimentService(experimentRepository *repository.ExperimentRepository, logger *slog.Logger) *ExperimentService {
	return &ExperimentService{
		experimentRepository: experimentRepository,
		logger:               logger,
	}
}

func (s *ExperimentService) CreateExperiment(ctx context.Context, experiment experiment.CreateExperimentRequest) ([]problems.Violation, error) {
	violations := validators.ValidateExperiment(experiment)
	if len(violations) > 0 {
		return violations, nil
	}

	err := s.experimentRepository.CreateNewExperiment(ctx, experiment)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
