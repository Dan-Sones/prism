package repository

import (
	"context"
	"experiment-simulator/internal/model"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type EventsRepository interface {
	GetCountOfEventForVariantAndExperimentInPhase(ctx context.Context, eventKey model.EventKey, variantKey model.VariantKey, experimentKey string, phase model.ExperimentPhaseType) (int, error)
}

type EventsRepositoryClickhouse struct {
	connection driver.Conn
}

func NewEventsRepositoryClickhouse(connection driver.Conn) *EventsRepositoryClickhouse {
	return &EventsRepositoryClickhouse{
		connection: connection,
	}
}

func (c *EventsRepositoryClickhouse) GetCountOfEventForVariantAndExperimentInPhase(ctx context.Context, eventKey model.EventKey, variantKey model.VariantKey, experimentKey string, phase model.ExperimentPhaseType) (int, error) {
	// having written this, it's making me think that maybe the field should just be "phase" rather than "is_aa"?
	isAA := true
	if phase == model.ExperimentPhaseAB {
		isAA = false
	}

	query := `
SELECT count() FROM cooked_events WHERE variant_key = @variant_key AND experiment_key = @experiment_key AND event_key = @event_key AND is_aa = @is_aa;
`
	var count uint64
	if err := c.connection.QueryRow(ctx, query,
		clickhouse.Named("variant_key", variantKey),
		clickhouse.Named("experiment_key", experimentKey),
		clickhouse.Named("event_key", eventKey),
		clickhouse.Named("is_aa", isAA),
	).Scan(&count); err != nil {
		return 0, err
	}

	return int(count), nil
}
