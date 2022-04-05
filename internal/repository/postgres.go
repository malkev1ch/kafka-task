package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/malkev1ch/kafka-task/internal/config"
)

type PostgresRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	conn, err := pgxpool.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{DB: conn}, nil
}

func (rps *PostgresRepository) SendBatch(ctx context.Context, batch *pgx.Batch) error {
	res := rps.DB.SendBatch(ctx, batch)
	if err := res.Close(); err != nil {
		return err
	}
	return nil
}
