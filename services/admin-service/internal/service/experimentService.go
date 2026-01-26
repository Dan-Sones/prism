package service

import (
	"admin-service/internal/errors"
	"admin-service/internal/repository"
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

func (s *ExperimentService) CreateExperiment(ctx context.Context, experiment model.Experiment) error {

	if experiment.Name == "" {
		return &errors.MissingFieldError{Field: "Name"}
	}

	err := s.experimentRepository.CreateNewExperiment(ctx, experiment)
	if err != nil {
		return err
	}

	return nil
}
