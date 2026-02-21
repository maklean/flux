package server_interface

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/maklean/flux/server/proto"
	"google.golang.org/grpc"
)

// StartgRPCServer starts a gRPC server at port GRPC_SERVER_PORT
func StartgRPCServer() {
	port := getPort()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen on tcp port 9000: %v\n", err)
	}

	// start new gRPC server and register service for network requests
	grpcServer := grpc.NewServer()
	pb.RegisterTelemetryServiceServer(grpcServer, &telemetryService{})

	log.Printf("Running TelemetryService gRPC server on port %d...", port)
	grpcServer.Serve(lis)
}

// getPort returns the GRPC_SERVER_PORT value from the environment files
func getPort() int {
	portStr, ok := os.LookupEnv("GRPC_SERVER_PORT")
	if !ok {
		log.Fatalf("failed finding GRPC_SERVER_PORT in env")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("failed to convert port to int")
	}

	return port
}
