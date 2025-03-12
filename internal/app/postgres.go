package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vit6556/ozon-internship-assignment/internal/config"
)

func InitPostgres() *pgxpool.Pool {
	ctx := context.Background()

	var postgresConfig config.PostgresConfig
	config.LoadConfig(&postgresConfig)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		postgresConfig.Username, postgresConfig.Password,
		net.JoinHostPort(postgresConfig.Host, postgresConfig.Port),
		postgresConfig.Name,
	)

	dbPool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("failed to create db pool")
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}

	return dbPool
}
