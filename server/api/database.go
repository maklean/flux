package api

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbConn *pgxpool.Pool = nil

// connectDatabase connects to a PostgreSQL database and returns the connection instance
func connectDatabase() *pgxpool.Pool {
	if dbConn != nil {
		log.Panicln("cannot connect to database more than once")
	}

	var err error

	// wait for a total of 10s for setup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbConn, err = pgxpool.New(ctx, getDatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// pgxpool.New() only configures the pool, and doesn't verify the connection
	if err = dbConn.Ping(ctx); err != nil {
		log.Fatalf("failed to reach database: %v", err)
	}

	execQuery := func(query, operation string) {
		_, err := dbConn.Exec(ctx, query)
		if err != nil {
			log.Fatalf("failed to %s: %v", operation, err)
		}
	}

	// run DB setup queries
	execQuery(q_CreateEncodersTable, "create encoders table")
	execQuery(q_CreateMetricsTable, "create metrics table")
	execQuery(idx_EncoderId_MetricsTable, "create encoder_id index for metrics table")

	return dbConn
}

// GetDB pings the database (waits 5s) and returns the connection instance if it's still alive.
func GetDB() *pgxpool.Pool {
	if dbConn == nil {
		log.Fatalf("database is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := dbConn.Ping(ctx); err != nil {
		log.Fatalf("loss connection to the database: %v", err)
	}

	return dbConn
}

// getDatabaseURL returns the PSQL_DB_URL value from the environment files
func getDatabaseURL() string {
	dbUrl, ok := os.LookupEnv("PSQL_DB_URL")

	if !ok {
		log.Fatalf("failed finding PSQL_DB_URL in env")
	}

	return dbUrl
}
