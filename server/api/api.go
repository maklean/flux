package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StartAPIServer starts a REST API server at port API_SERVER_PORT
func StartAPIServer() {
	// Connect to database
	conn := connectDatabase()
	defer conn.Close(context.Background())
	log.Println("Connected to database successfully...")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Register routes
	main := r.Group("/api")
	main.GET("/metrics", getAllMetrics)
	main.GET("/encoders", getEncoders)

	// Run API server
	port := getPort()

	if gin.Mode() == gin.ReleaseMode {
		log.Printf("Running API Server on port %d...", port)
	}

	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("failed to start api server")
	}
}

// getPort returns the API_SERVER_PORT value from the environment files
func getPort() int {
	portStr, ok := os.LookupEnv("API_SERVER_PORT")
	if !ok {
		log.Fatalf("failed finding API_SERVER_PORT in env")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to convert port to int")
	}

	return port
}
