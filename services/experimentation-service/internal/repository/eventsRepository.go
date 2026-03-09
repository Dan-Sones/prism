package repository

import (
	"context"
	"experimentation-service/internal/model/graph"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type ClickHouseEventsRepository struct {
	connection driver.Conn
}

func NewClickHouseEventsRepository(connection driver.Conn) *ClickHouseEventsRepository {
	return &ClickHouseEventsRepository{
		connection: connection,
	}
}

func (e *ClickHouseEventsRepository) GetEventKeyUsageForLastXMinutesWithMinuteInterval(ctx context.Context, eventKey string, xMinutes int) ([]graph.TimeScaleDataPoint, error) {
	query := fmt.Sprintf(`
      SELECT
          toStartOfMinute(received_at) AS timestamp_aggregation,
          count() AS event_count
      FROM events
      WHERE received_at >= toStartOfMinute(now() - INTERVAL %d MINUTE)
        AND event_key = @event_key
      GROUP BY timestamp_aggregation
      ORDER BY timestamp_aggregation WITH FILL
          FROM toStartOfMinute(now() - INTERVAL %d MINUTE)
          TO toStartOfMinute(now())
          STEP INTERVAL 1 MINUTE
  `, xMinutes, xMinutes)

	return e.handleDataPointsResult(ctx, query, eventKey)
}

func (e *ClickHouseEventsRepository) GetEventKeyUsageForLastHourWith5MinuteInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error) {
	query := `
	SELECT
	toStartOfFiveMinute(received_at) AS timestamp_aggregation,
	count() AS event_count
	FROM events
	WHERE received_at >= toStartOfFiveMinute(now() - INTERVAL 1 Hour)
	  AND event_key = @event_key
	GROUP BY timestamp_aggregation
	ORDER BY timestamp_aggregation WITH FILL
	FROM toStartOfFiveMinute(now() - INTERVAL 1 Hour)
		TO toStartOfFiveMinute(now())
		STEP INTERVAL 5 MINUTE
  `

	return e.handleDataPointsResult(ctx, query, eventKey)
}

func (e *ClickHouseEventsRepository) GetEventKeyUsageForLast24HoursWith1HourInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error) {
	query := `
	SELECT
		toStartOfHour(received_at) AS timestamp_aggregation,
		count() AS event_count
	FROM events
	WHERE received_at >= toStartOfHour(now() - INTERVAL 1 Day)
	  AND event_key = @event_key
	GROUP BY timestamp_aggregation
	ORDER BY timestamp_aggregation WITH FILL
	FROM toStartOfHour(now() - INTERVAL 1 DAY)
		TO toStartOfHour(now())
    	STEP INTERVAL 1 Hour
  `

	return e.handleDataPointsResult(ctx, query, eventKey)
}

func (e *ClickHouseEventsRepository) GetEventKeyUsageForLast7DaysWithDayInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error) {
	query := `
		SELECT
		toStartOfDay(received_at) AS timestamp_aggregation,
		count() AS event_count
		FROM events
		WHERE received_at >= toStartOfDay(now() - INTERVAL 1 WEEK)
		  AND event_key = @event_key
		GROUP BY timestamp_aggregation
		ORDER BY timestamp_aggregation WITH FILL
		FROM toStartOfDay(now() - INTERVAL 1 WEEK)
			TO toStartOfDay(now())
			STEP INTERVAL 1 Day
  `

	return e.handleDataPointsResult(ctx, query, eventKey)
}

func (e *ClickHouseEventsRepository) GetEventKeyUsageForLast30DaysWithDayInterval(ctx context.Context, eventKey string) ([]graph.TimeScaleDataPoint, error) {
	query := `
	SELECT
	toStartOfDay(received_at) AS timestamp_aggregation,
	count() AS event_count
	FROM events
	WHERE received_at >= toStartOfDay(now() - INTERVAL 1 MONTH )
	  AND event_key = @event_key
	GROUP BY timestamp_aggregation
	ORDER BY timestamp_aggregation WITH FILL
	FROM toStartOfDay(now() - INTERVAL 1 MONTH)
		TO toStartOfDay(now())
		STEP INTERVAL 1 Day
  `

	return e.handleDataPointsResult(ctx, query, eventKey)
}

func (e *ClickHouseEventsRepository) handleDataPointsResult(ctx context.Context, query, eventKey string) ([]graph.TimeScaleDataPoint, error) {

	rows, err := e.connection.Query(ctx, query, clickhouse.Named("event_key", eventKey))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dataPoints []graph.TimeScaleDataPoint
	for rows.Next() {
		var (
			timestamp  time.Time
			eventCount uint64
		)

		if err := rows.Scan(&timestamp, &eventCount); err != nil {
			return nil, err
		}

		dataPoints = append(dataPoints, graph.TimeScaleDataPoint{
			Time:  timestamp,
			Value: int(eventCount),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataPoints, nil
}

func (e *ClickHouseEventsRepository) GetTotalEventsPast24HoursForEventKey(ctx context.Context, eventKey string) (int, error) {
	query := `
	SELECT
		count() AS event_count
	FROM events
	WHERE received_at >= now() - INTERVAL 1 Day
	  AND event_key = @event_key
  `

	var eventCount uint64
	err := e.connection.QueryRow(ctx, query, clickhouse.Named("event_key", eventKey)).Scan(&eventCount)
	if err != nil {
		return 0, err
	}

	return int(eventCount), nil
}

func (e *ClickHouseEventsRepository) GetTotalEventsPast7DaysForEventKey(ctx context.Context, eventKey string) (int, error) {
	query := `
	SELECT
		count() AS event_count
	FROM events
	WHERE received_at >= now() - INTERVAL 7 Day
	  AND event_key = @event_key
  `

	var eventCount uint64
	err := e.connection.QueryRow(ctx, query, clickhouse.Named("event_key", eventKey)).Scan(&eventCount)
	if err != nil {
		return 0, err
	}

	return int(eventCount), nil
}

func (e *ClickHouseEventsRepository) GetLastReceivedTimeForEventKey(ctx context.Context, eventKey string) (time.Time, error) {
	query := `
	SELECT
	max(received_at) AS last_received_time
	FROM events
	WHERE event_key = @event_key
  `

	var lastReceivedTime time.Time
	err := e.connection.QueryRow(ctx, query, clickhouse.Named("event_key", eventKey)).Scan(&lastReceivedTime)
	if err != nil {
		return time.Time{}, err
	}

	return lastReceivedTime, nil
}
