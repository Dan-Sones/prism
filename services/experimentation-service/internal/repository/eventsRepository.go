package repository

import (
	"context"
	"experimentation-service/internal/model/graph"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type EventsRepository struct {
	connection driver.Conn
}

func NewEventsRepository(connection driver.Conn) *EventsRepository {
	return &EventsRepository{
		connection: connection,
	}
}

func (e *EventsRepository) GetEventKeyUsageForLastXMinutesWithMinuteInterval(ctx context.Context, eventKey string, xMinutes int) ([]graph.TimeScaleDataPoint, error) {
	query := fmt.Sprintf(`
      SELECT
          toStartOfMinute(received_at) AS timestamp_aggregation,
          count() AS event_count
      FROM events
      WHERE received_at >= toStartOfMinute(now() - INTERVAL %d MINUTE)
        AND event_key = @event_key
      GROUP BY timestamp_aggregation
      ORDER BY timestamp_aggregation WITH FILL3
          FROM toStartOfMinute(now() - INTERVAL %d MINUTE)
          TO toStartOfMinute(now())
          STEP INTERVAL 1 MINUTE
  `, xMinutes, xMinutes)

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
			Value: int64(eventCount),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dataPoints, nil
}
