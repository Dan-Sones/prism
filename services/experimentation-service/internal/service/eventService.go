package service

import (
	"context"
	model2 "experimentation-service/internal/model/graph"
	"experimentation-service/internal/repository"
	"log/slog"
)

type EventsServiceInterface interface {
	GetEventUsageOverPeriod(ctx context.Context, scale model2.GraphTimeScale, eventKey string) ([]model2.TimeScaleDataPoint, error)
}

type EventService struct {
	eventsRepository *repository.EventsRepository
	logger           *slog.Logger
}

func NewEventsService(eventsRepository *repository.EventsRepository, logger *slog.Logger) *EventService {
	return &EventService{
		eventsRepository: eventsRepository,
		logger:           logger,
	}
}

var scaleToMinutes = map[model2.GraphTimeScale]int{
	model2.Minute:    1,
	model2.TenMinute: 10,
	model2.HalfHour:  30,
	model2.ScaleHour: 60,
}

func (e *EventService) GetEventUsageOverPeriod(ctx context.Context, scale model2.GraphTimeScale, eventKey string) ([]model2.TimeScaleDataPoint, error) {
	minutes, ok := scaleToMinutes[scale]
	if !ok {
		return nil, nil
	}

	res, err := e.eventsRepository.GetEventKeyUsageForLastXMinutesWithMinuteInterval(ctx, eventKey, minutes)
	if err != nil {
		e.logger.Error(err.Error())
		return nil, err
	}
	return res, nil
}
