package api

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn = nil

// connectDatabase connects to a PostgreSQL database and returns the connection instance
func connectDatabase() *pgx.Conn {
	if dbConn != nil {
		log.Panicln("cannot connect to database more than once")
	}

	var err error

	// TODO: probably don't use context.Background(), make global app context for DB queries.
	dbConn, err = pgx.Connect(context.Background(), getDatabaseURL())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// initialize tables
	_, err = dbConn.Exec(context.Background(), q_CreateEncodersTable)
	if err != nil {
		log.Fatalf("failed to create encoders table: %v", err)
	}

	_, err = dbConn.Exec(context.Background(), q_CreateMetricsTable)
	if err != nil {
		log.Fatalf("failed to create metrics table: %v", err)
	}

	_, err = dbConn.Exec(context.Background(), idx_EncoderId_MetricsTable)
	if err != nil {
		log.Fatalf("failed to encoder_id index for metrics table: %v", err)
	}

	return dbConn
}

// GetDB pings the database and returns the connection instance if it's still alive.
func GetDB() *pgx.Conn {
	if dbConn == nil {
		log.Fatalf("database is not initialized")
	}

	if err := dbConn.Ping(context.Background()); err != nil {
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
