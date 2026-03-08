package service

import (
	"context"
	"experimentation-service/internal/model/event"
	model2 "experimentation-service/internal/model/graph"
	"experimentation-service/internal/repository"
	"log/slog"
	"math/rand"
)

type EventsServiceInterface interface {
	GetEventUsageOverPeriod(ctx context.Context, scale model2.GraphTimeScale, eventKey string) ([]model2.TimeScaleDataPoint, error)
}

type EventService struct {
	eventsRepository        EventsRepository
	eventsCatalogRepository *repository.EventsCatalogRepository
	logger                  *slog.Logger
}

func NewEventsService(eventsRepository EventsRepository, eventsCatalogRepository *repository.EventsCatalogRepository, logger *slog.Logger) *EventService {
	return &EventService{
		eventsRepository:        eventsRepository,
		eventsCatalogRepository: eventsCatalogRepository,
		logger:                  logger,
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

func (e *EventService) GetLiveEventStatistics(ctx context.Context, eventKey string) (event.LiveEventStatistics, error) {
	lastReceivedTime, err := e.eventsRepository.GetLastReceivedTimeForEventKey(ctx, eventKey)
	if err != nil {
		e.logger.Error(err.Error())
		return event.LiveEventStatistics{}, err
	}

	countPast7Days, err := e.eventsRepository.GetTotalEventsPast7DaysForEventKey(ctx, eventKey)
	if err != nil {
		e.logger.Error(err.Error())
		return event.LiveEventStatistics{}, err
	}

	countPast24Hours, err := e.eventsRepository.GetTotalEventsPast24HoursForEventKey(ctx, eventKey)
	if err != nil {
		e.logger.Error(err.Error())
		return event.LiveEventStatistics{}, err
	}

	missingRates, err := e.getMissingRatesForEventType(ctx, eventKey)
	if err != nil {
		e.logger.Error(err.Error())
		return event.LiveEventStatistics{}, err
	}

	// TODO: Fetch actual missing rates
	// Fetch the fields from the database
	// Then look up in metrics to get actual missing rates for each field
	return event.LiveEventStatistics{
		MissingRates:           missingRates,
		TotalEventsPast7Days:   countPast7Days,
		TotalEventsPast24Hours: countPast24Hours,
		LastReceivedTime:       lastReceivedTime,
	}, nil

}

func (e *EventService) getMissingRatesForEventType(ctx context.Context, eventKey string) (map[string]int64, error) {
	eventType, err := e.eventsCatalogRepository.GetEventTypeByKey(ctx, eventKey)
	if err != nil {
		return nil, err
	}

	missingRates := make(map[string]int64)
	for _, field := range eventType.Fields {
		missingRates[field.FieldKey] = int64(rand.Intn(100-0) + 0)
	}

	return missingRates, nil
}
