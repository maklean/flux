package api

import (
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
	metricRows, err := dbConn.Query(c.Request.Context(), selectAllFromMetricsTable)
	if err != nil {
		log.Printf("failed to select all from metrics table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	defer metricRows.Close()

	// collect all metrics into slice of Metric
	metrics, err := pgx.CollectRows(metricRows, pgx.RowToStructByPos[Metric])
	if err != nil {
		log.Printf("failed to collect metrics into rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// getAllEncoderMetrics returns all metrics from the 'metrics' table with a specific encoder_id.
func getAllEncoderMetrics(c *gin.Context) {
	ctx := c.Request.Context()
	encoderId := c.Param("encoderId")

	dbConn := GetDB()

	// check if encoder exists
	var dbEncoderId string
	err := dbConn.QueryRow(ctx, SelectFromEncodersTable, encoderId).Scan(&dbEncoderId)
	if err != nil {
		// encoder doesn't exist in database
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": fmt.Sprintf("failed to find encoder with id '%s'", encoderId),
			})
			return
		}

		log.Printf("failed to query database for encoder: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	// retrieve all metric rows with id 'encoderId'
	metricRows, err := dbConn.Query(ctx, selectAllMetricsFromEncoderId, dbEncoderId)
	if err != nil {
		log.Printf("failed to retrieve metric rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer metricRows.Close()

	// collect metricRows into slice of Metric
	metrics, err := pgx.CollectRows(metricRows, pgx.RowToStructByPos[Metric])
	if err != nil {
		log.Printf("failed to collect metrics into rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}

// getAllEncoders returns all encoders from the 'encoders' table in the database.
func getAllEncoders(c *gin.Context) {
	dbConn := GetDB()

	// retrieve all encoder rows
	encoderRows, err := dbConn.Query(c.Request.Context(), selectAllFromEncodersTable)
	if err != nil {
		log.Printf("failed to select all from encoders table: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	defer encoderRows.Close()

	// collect encoder rows into slice of Encoder
	encoders, err := pgx.CollectRows(encoderRows, pgx.RowToStructByPos[Encoder])
	if err != nil {
		log.Printf("failed to collect encoders into rows: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"encoders": encoders,
	})
}
