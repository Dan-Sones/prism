package repository

import (
	"context"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentRepository struct {
	pgxPool *pgxpool.Pool
	logger  *slog.Logger
}

func NewExperimentRepository(pgxPool *pgxpool.Pool, logger *slog.Logger) *ExperimentRepository {
	return &ExperimentRepository{
		pgxPool: pgxPool,
		logger:  logger,
	}
}

func (r *ExperimentRepository) CreateNewExperiment(ctx context.Context, experiment model.Experiment) error {
	sql := `INSERT INTO prism.experiments (name) VALUES ($1) RETURNING id`

	err := r.pgxPool.QueryRow(ctx, sql, experiment.Name).Scan(&experiment.ID)
	if err != nil {
		r.logger.Error("Failed to create new experiment", slog.String("error", err.Error()))
		return err
	}

	return nil
}
