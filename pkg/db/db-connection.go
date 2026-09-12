package db

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Conn *pgxpool.Pool

func ConnectDB() error {
	dbURL := strings.TrimSpace("postgres://" +
		os.Getenv("DATABASE_USERNAME") +
		":" +
		os.Getenv("DATABASE_PASSWORD") +
		"@" +
		os.Getenv("DATABASE_URL") +
		"/" +
		os.Getenv("DATABASE_NAME") +
		"?sslmode=disable",
	)
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set; configure the PostgreSQL connection string in the deployment environment")
	}

	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("parse database config: %w", err)
	}

	config.MaxConns = 20
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	Conn, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return err
	}

	if err := Conn.Ping(context.Background()); err != nil {
		Conn.Close()
		return err
	}

	log.Println("DB Connected")
	return nil

}
