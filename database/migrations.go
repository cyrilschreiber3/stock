package database

import (
	"database/sql"
	"log/slog"

	"github.com/cyrilschreiber3/stock/database/migrations"
	"github.com/cyrilschreiber3/stock/database/seed"
	"github.com/cyrilschreiber3/stock/logger"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/pressly/goose/v3"
)

var gooseDB *sql.DB

func Migrate() {
	subCmd := utils.GetCommandArg(0)
	switch subCmd {
	case "up":
		slog.Info("Migrating database to the latest version")
		Up()
	case "down":
		slog.Info("Rolling back the last database migration")
		Down()
	case "down-to":
		slog.Info("Rolling back to a specific migration version")
		if len(utils.GetCommandArgs()) < 1 {
			logger.Fatal("Please provide a target migration version to rollback to")
		}
		DownTo(utils.GetCommandArgInt(0))
	case "reset":
		slog.Info("Resetting the database")
		Reset()
	case "state":
		slog.Info("Checking the current migration state")
		MigrationState()
	case "seed":
		slog.Info("Seeding the database")
		Seed()
	case "reset-seed":
		slog.Info("Resetting the database seed")
		ResetSeed()
	default:
		logger.Fatal("Unknown migration subcommand", "subcommand", subCmd)
	}
}

func Up() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(migrations.SchemaMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	err := goose.Up(gooseDB, ".")
	if err != nil {
		logger.Fatal("unable to migrate database", "error", err)
	}

	slog.Info("Database migration completed successfully")
}

func Down() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(migrations.SchemaMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	version, err := goose.EnsureDBVersion(gooseDB)
	if err != nil {
		logger.Fatal("unable to ensure database version", "error", err)
	}

	slog.Info("Rolling back database migrations from version", "version", version)

	err = goose.Down(gooseDB, ".")
	if err != nil {
		logger.Fatal("unable to rollback database migrations", "error", err)
	}

	slog.Info("Database rollback completed successfully")
}

func DownTo(targetVersion int) {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(migrations.SchemaMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	slog.Info("Rolling back database migrations to version", "version", targetVersion)

	err := goose.DownTo(gooseDB, ".", int64(targetVersion))
	if err != nil {
		logger.Fatal("unable to rollback database migrations to version", "version", targetVersion, "error", err)
	}

	slog.Info("Database rollback to version completed successfully", "version", targetVersion)
}

func Reset() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(migrations.SchemaMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	err := goose.Reset(gooseDB, ".")
	if err != nil {
		logger.Fatal("unable to reset database migrations", "error", err)
	}

	slog.Info("Database reset completed successfully")
}

func MigrationState() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(migrations.SchemaMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	version, err := goose.EnsureDBVersion(gooseDB)
	if err != nil {
		logger.Fatal("unable to ensure database version", "error", err)
	}

	slog.Info("Current database version", "version", version)

	err = goose.Status(gooseDB, ".")
	if err != nil {
		logger.Fatal("unable to get database migration status", "error", err)
	}

}

func Seed() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(seed.SeedMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	err := goose.Up(gooseDB, ".", goose.WithNoVersioning())
	if err != nil {
		logger.Fatal("unable to seed database", "error", err)
	}

	slog.Info("Database seeding completed successfully")
}

func ResetSeed() {
	initGooseDB()
	defer gooseDB.Close()
	goose.SetBaseFS(seed.SeedMigrations)

	acqireMigrationLock()
	defer releaseMigrationLock()

	err := goose.Reset(gooseDB, ".", goose.WithNoVersioning())
	if err != nil {
		logger.Fatal("unable to reset database seed", "error", err)
	}

	slog.Info("Database seed reset completed successfully")
}
