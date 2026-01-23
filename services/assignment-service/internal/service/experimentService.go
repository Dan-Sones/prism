package service

import (
	"log/slog"
)

type ExperimentService struct {
	logger *slog.Logger
}

func NewExperimentService(logger *slog.Logger) *ExperimentService {
	return &ExperimentService{
		logger: logger,
	}
}

func (e *ExperimentService) GetVariantsForBucket(bucketId int) {

}
