package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// getAllMetrics returns all the metrics from the 'metrics' table in the database.
func getAllMetrics(c *gin.Context) {
	dbConn := GetDB()

	metricRows, err := dbConn.Query(context.Background(), selectAllFromMetricsTable)
	if err != nil {
		log.Fatalf("failed to select all from metrics table: %v", err)
	}
	defer metricRows.Close()

	metrics, err := pgx.CollectRows(metricRows, pgx.RowToStructByPos[Metric])
	if err != nil {
		log.Fatalf("failed to collect metrics into rows: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

func getAllEncoders(c *gin.Context) {
	dbConn := GetDB()

	encoderRows, err := dbConn.Query(context.Background(), selectAllFromEncodersTable)
	if err != nil {
		log.Fatalf("failed to select all from encoders table: %v", err)
	}

	defer encoderRows.Close()

	encoders, err := pgx.CollectRows(encoderRows, pgx.RowToStructByPos[Encoder])
	if err != nil {
		log.Fatalf("failed to collect encoders into rows: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"encoders": encoders,
	})
}
