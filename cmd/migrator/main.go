package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to storage path")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "", "Path to migrations table")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migrations-path is required")
	}
	dbURL := fmt.Sprintf("%s?x-migrations-table=%s", storagePath, migrationsTable)
	m, err := migrate.New(
		"file://"+migrationsPath, dbURL,
	)

	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}
	fmt.Println("migrations applied successfully")
}
