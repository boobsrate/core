package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun/driver/pgdriver"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("creating logger: %w", err))
	}
	logger = logger.Named("migrate")
	defer logger.Sync() // nolint: errcheck

	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatal("load configuration: ", zap.Error(err))
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseDSN)))
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatal("connect to database: ", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", cfg.MigrationsDir),
		"postgres", driver)
	if err != nil {
		log.Fatal("create migrate: ", zap.Error(err))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("migrate up: ", zap.Error(err))
	}
}
