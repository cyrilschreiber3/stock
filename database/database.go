package database

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/cyrilschreiber3/stock/logger"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Init() {
	if Pool != nil {
		slog.Info("Database pool is already initialized")
		return
	}

	slog.Info("Initializing PostgreSQL database...")
	dsn := utils.GetEnv("DB_URL", "")
	if dsn == "" {
		logger.Fatal("unable to initialize pgx pool: DB_URL is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		logger.Fatal("unable to create pgx pool", "error", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Fatal("unable to connect to the database", "error", err)
	}

	Pool = pool
}

func Close() {
	slog.Info("Closing PostgreSQL database connection...")
	if Pool != nil {
		Pool.Close()
		Pool = nil
	}
}

func RollbackTransaction(ctx context.Context, tx pgx.Tx, retErr *error) {
	if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && *retErr == nil {
		if !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			*retErr = rollbackErr
		}
	}
}
