package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/malkev1ch/kafka-task/internal/config"
	log "github.com/sirupsen/logrus"
)

type PostgresRepository struct {
	DB *pgxpool.Pool
}

func NewPostgresRepository(cfg *config.Config) *PostgresRepository {
	conn, err := pgxpool.Connect(context.Background(), cfg.PostgresURL)
	if err != nil {
		log.Fatalf("repository: connection failed - %e", err)
	}
	log.Info("repository: successfully connected")
	return &PostgresRepository{DB: conn}
}

func (rps *PostgresRepository) SendBatch(ctx context.Context, batch *pgx.Batch) error {
	res := rps.DB.SendBatch(ctx, batch)
	if err := res.Close(); err != nil {
		log.Fatalf("postgres: batch error - %e", err)
		return err
	}
	return nil
}
