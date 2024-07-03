package main

import (
	"database/sql"
	"errors"
	"log"
	"rello-test-task/util"

	"github.com/golang-migrate/migrate/v4"

	"rello-test-task/api"
	"rello-test-task/db"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// @title Rello Test Task
// @version 0.0.1
// @description REST API to manage buildings and apartments

// @host localhost:8080
// @basePath /
func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("failed to connect to database")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("Cannot start the server at the address {}:\n{}", config.ServerAddress, err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	m, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}
	version, _, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Fatal("failed to fetch migration version:", err)
	}

	if errors.Is(err, migrate.ErrNilVersion) {
		if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("failed to run migrate up:", err)
		}
		log.Println("db migrated successfully")
	} else {
		log.Printf("db already migrated to version: %d\n", version)
	}
}
