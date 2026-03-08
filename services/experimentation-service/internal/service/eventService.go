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
	model2.TenMinute: 10,
	model2.HalfHour:  30,
}

func (e *EventService) GetEventUsageOverPeriod(ctx context.Context, scale model2.GraphTimeScale, eventKey string) ([]model2.TimeScaleDataPoint, error) {
	var (
		res []model2.TimeScaleDataPoint
		err error
	)

	switch scale {
	case model2.ScaleHour:
		res, err = e.eventsRepository.GetEventKeyUsageForLastHourWith5MinuteInterval(ctx, eventKey)
	case model2.ScaleDay:
		res, err = e.eventsRepository.GetEventKeyUsageForLast24HoursWith1HourInterval(ctx, eventKey)
	case model2.ScaleWeek:
		res, err = e.eventsRepository.GetEventKeyUsageForLast7DaysWithDayInterval(ctx, eventKey)
	case model2.ScaleMonth:
		res, err = e.eventsRepository.GetEventKeyUsageForLast30DaysWithDayInterval(ctx, eventKey)
	default:
		minutes, ok := scaleToMinutes[scale]
		if !ok {
			return nil, nil
		}
		res, err = e.eventsRepository.GetEventKeyUsageForLastXMinutesWithMinuteInterval(ctx, eventKey, minutes)
	}

	if err != nil {
		e.logger.Error(err.Error())
		return nil, err
	}
	return res, nil
}
