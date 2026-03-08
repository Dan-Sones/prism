package service

import (
	"context"
	"experimentation-service/internal/model/graph"
	"time"
)

// The interface can be implemented by any repository with any data source. This includes datawarehouses
type EventsRepository interface {
	GetEventKeyUsageForLastXMinutesWithMinuteInterval(ctx context.Context, eventKey string, xMinutes int) ([]graph.TimeScaleDataPoint, error)
	GetEventKeyUsageForLastHourWith5MinuteInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error)
	GetEventKeyUsageForLast24HoursWith1HourInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error)
	GetEventKeyUsageForLast7DaysWithDayInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error)
	GetEventKeyUsageForLast30DaysWithDayInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error)
	GetTotalEventsPast24HoursForEventKey(ctx context.Context, eventKey string) (int64, error)
	GetTotalEventsPast7DaysForEventKey(ctx context.Context, eventKey string) (int64, error)
	GetLastReceivedTimeForEventKey(ctx context.Context, eventKey string) (time.Time, error)
}
