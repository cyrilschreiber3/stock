package database

import (
	"context"
	"log"
	"time"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Init creates a single long-lived pgxpool.Pool and assigns it to database.Pool.
// It will log and exit the process on failure to connect.
func Init() {
	log.Println("Initializing PostgreSQL database...")

	dsn := utils.GetEnv("DB_URL", "")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v", err)
	}

	Pool = pool
}

// Close closes the global connection pool if it was created.
func Close() {
	log.Println("Closing PostgreSQL database connection...")
	if Pool != nil {
		Pool.Close()
	}
}
