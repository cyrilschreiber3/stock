package database

import (
	"context"
	"time"

	"github.com/cyrilschreiber3/stock/logger"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

var migrationLock *SessionLock

func initGooseDB() {
	Init()

	gooseDB = stdlib.OpenDBFromPool(Pool)

	err := goose.SetDialect("postgres")
	if err != nil {
		logger.Fatal("unable to set goose dialect", "error", err)
	}
}

func acqireMigrationLock() {
	if migrationLock != nil {
		logger.Fatal("migration advisory lock is already held by this process")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	lock, err := AcquireSessionLock(ctx, Pool, AdvisoryLockKeyMigration, AdvisoryLockTypeExclusive)
	if err != nil {
		logger.Fatal("unable to acquire migration advisory lock", "error", err)
	}

	migrationLock = lock
}

func releaseMigrationLock() {
	if migrationLock == nil {
		logger.Fatal("migration advisory lock is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := migrationLock.Release(ctx); err != nil {
		logger.Fatal("unable to release migration advisory lock", "error", err)
	}

	migrationLock = nil
}
