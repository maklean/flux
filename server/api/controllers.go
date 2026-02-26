package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// getAllMetrics returns all the metrics from the 'metrics' table in the database.
func getAllMetrics(c *gin.Context) {
	dbConn := GetDB()

	// retrieve all metric rows
	metricRows, err := dbConn.Query(context.Background(), selectAllFromMetricsTable)
	if err != nil {
		log.Fatalf("failed to select all from metrics table: %v", err)
	}
	defer metricRows.Close()

	// collect all metrics into slice of Metric
	metrics, err := pgx.CollectRows(metricRows, pgx.RowToStructByPos[Metric])
	if err != nil {
		log.Fatalf("failed to collect metrics into rows: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// getAllEncoderMetrics returns all metrics from the 'metrics' table with a specific encoder_id.
func getAllEncoderMetrics(c *gin.Context) {
	encoderId := c.Param("encoderId")

	dbConn := GetDB()

	// check if encoder exists
	var dbEncoderId string
	err := dbConn.QueryRow(context.Background(), SelectFromEncodersTable, encoderId).Scan(&dbEncoderId)
	if err != nil {
		// encoder doesn't exist in database
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("failed to find encoder with id '%s'", encoderId),
			})
			return
		}

		log.Fatalf("failed to query database for encoder: %v", err)
	}

	// retrieve all metric rows with id 'encoderId'
	metricRows, err := dbConn.Query(context.Background(), selectAllMetricsFromEncoderId, dbEncoderId)
	if err != nil {
		log.Fatalf("failed to retrieve metric rows: %v", err)
	}

	defer metricRows.Close()

	// collect metricRows into slice of Metric
	metrics, err := pgx.CollectRows(metricRows, pgx.RowToStructByPos[Metric])
	if err != nil {
		log.Fatalf("failed to collect metrics into rows: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// getAllEncoders returns all encoders from the 'encoders' table in the database.
func getAllEncoders(c *gin.Context) {
	dbConn := GetDB()

	// retrieve all encoder rows
	encoderRows, err := dbConn.Query(context.Background(), selectAllFromEncodersTable)
	if err != nil {
		log.Fatalf("failed to select all from encoders table: %v", err)
	}

	defer encoderRows.Close()

	// collect encoder rows into slice of Encoder
	encoders, err := pgx.CollectRows(encoderRows, pgx.RowToStructByPos[Encoder])
	if err != nil {
		log.Fatalf("failed to collect encoders into rows: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"encoders": encoders,
	})
}
