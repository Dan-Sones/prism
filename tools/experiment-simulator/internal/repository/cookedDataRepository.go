package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type CookedDataRepository interface {
	GetCountOfEventForVariantAndExperiment(ctx context.Context, eventKey string, variantKey string, experimentKey string) (int, error)
}

type CookedDataRepositoryClickhouse struct {
	connection driver.Conn
}

func NewCookedDataRepositoryClickhouse(connection driver.Conn) *CookedDataRepositoryClickhouse {
	return &CookedDataRepositoryClickhouse{
		connection: connection,
	}
}

func (c *CookedDataRepositoryClickhouse) GetCountOfEventForVariantAndExperiment(ctx context.Context, eventKey string, variantKey string, experimentKey string) (int, error) {
	//TODO implement me
	query := `
SELECT count() FROM cooked_events WHERE variant_key == @variant_key AND experiment_key == @experiment_key AND event_key == @event_key;
`
	var count uint64
	if err := c.connection.QueryRow(ctx, query,
		clickhouse.Named("variant_key", variantKey),
		clickhouse.Named("experiment_key", experimentKey),
		clickhouse.Named("event_key", eventKey),
	).Scan(&count); err != nil {
		return 0, err
	}

	return int(count), nil
}
