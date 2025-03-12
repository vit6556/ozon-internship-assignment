package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
)

func help() {
	log.Println("Usage:")
	log.Println("  migrator up        - Apply all migrations")
	log.Println("  migrator down      - Rollback the last migration")
}

func main() {
	var databaseConfig config.PostgresConfig
	config.LoadConfig(&databaseConfig)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		databaseConfig.Username, databaseConfig.Password,
		net.JoinHostPort(databaseConfig.Host, databaseConfig.Port),
		databaseConfig.Name,
	)

	m, err := migrate.New(fmt.Sprintf("file://%s", databaseConfig.MigrationsPath), connString)
	if err != nil {
		log.Fatalf("error initializing migrations: %v", err)
	}

	flag.Usage = help
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("error applying migrations: %v", err)
		}
		log.Println("all migrations applied successfully")

	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("error rolling back migration: %v", err)
		}
		log.Println("last migration rolled back")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
