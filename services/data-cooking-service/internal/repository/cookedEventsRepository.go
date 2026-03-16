package repository

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type CookedEventsRepository interface {
	// Should not be any!!!! we need to firm up the schema first...
	InsertBatch(ctx context.Context, events any) error
}

type CookedEventsRepositoryClickhouse struct {
	connection driver.Conn
}

func NewCookedEventsRepositoryClickhouse(connection driver.Conn) *CookedEventsRepositoryClickhouse {
	return &CookedEventsRepositoryClickhouse{
		connection: connection,
	}
}

func (r *CookedEventsRepositoryClickhouse) InsertBatch(ctx context.Context, events any) error {
	return nil
}
