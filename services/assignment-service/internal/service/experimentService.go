package service

import (
	"assignment-service/internal/repository"
	"context"
	"fmt"
	"log/slog"
)

type ExperimentService struct {
	logger               *slog.Logger
	experimentRepository *repository.ExperimentRepository
}

func NewExperimentService(experimentRepository *repository.ExperimentRepository, logger *slog.Logger) *ExperimentService {
	return &ExperimentService{
		logger:               logger,
		experimentRepository: experimentRepository,
	}
}

func (e *ExperimentService) GetVariantsForBucket(bucketId int) {

	res, err := e.experimentRepository.GetExperimentsAndVariantsForBucket(context.Background(), bucketId)
	if err != nil {
		e.logger.Error("Error getting variants for bucket", slog.String("error", err.Error()))
		return
	}

	fmt.Printf("%+v\n", res)

}
