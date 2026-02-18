package service

import (
	"admin-service/internal/problems"
	"admin-service/internal/repository"
	"admin-service/internal/validators"
	"context"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model"
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

func (s *ExperimentService) CreateExperiment(ctx context.Context, experiment model.Experiment) ([]problems.Violation, error) {
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
