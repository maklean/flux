package server_interface

import (
	"context"

	pb "github.com/maklean/flux/server/proto"
)

type telemetryService struct {
	pb.UnimplementedTelemetryServiceServer
}

// RecordMetrics stores the metrics in the given request in the database
func (telemetryService) RecordMetrics(ctx context.Context, tr *pb.TelemetryRequest) (*pb.TelemetryResponse, error) {
	// TODO: store metrics in database

	return &pb.TelemetryResponse{
		Successful: true,
		Message:    "Successfully stored metric in database.",
	}, nil
}
