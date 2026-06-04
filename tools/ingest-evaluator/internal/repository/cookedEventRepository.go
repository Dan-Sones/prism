package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type CookedEventsRepositoryClickhouse struct {
	connection driver.Conn
}

func NewCookedEventsRepositoryClickhouse(connection driver.Conn) *CookedEventsRepositoryClickhouse {
	return &CookedEventsRepositoryClickhouse{
		connection: connection,
	}
}

func (c *CookedEventsRepositoryClickhouse) GetCountOfEventInPhaseForExperimentKey(ctx context.Context, experimentKey, phase string) (int, error) {
	isAA := true
	if phase == "ab" {
		isAA = false
	}

	query := `
		SELECT count() FROM cooked_events WHERE  is_aa = @is_aa AND experiment_key = @experiment_key;
	`

	var count uint64
	if err := c.connection.QueryRow(ctx, query,
		clickhouse.Named("experiment_key", experimentKey),
		clickhouse.Named("is_aa", isAA),
	).Scan(&count); err != nil {
		return 0, err
	}

	return int(count), nil
}
