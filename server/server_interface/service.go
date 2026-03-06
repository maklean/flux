package server_interface

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/maklean/flux/server/api"
	pb "github.com/maklean/flux/server/proto"
)

type telemetryService struct {
	pb.UnimplementedTelemetryServiceServer
}

// RecordMetrics stores the metrics in the given request in the database
func (telemetryService) RecordMetrics(ctx context.Context, tr *pb.TelemetryRequest) (*pb.TelemetryResponse, error) {
	dbConn := api.GetDB()

	// insert into encoders table if needed
	var encoderId string
	err := dbConn.QueryRow(ctx, api.SelectFromEncodersTable, tr.EncoderId).Scan(&encoderId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			_, err = dbConn.Exec(ctx, api.InsertIntoEncodersTable, tr.EncoderId)

			if err != nil {
				log.Printf("failed to insert into encoders table: %v", err)
				return nil, errors.New("internal server error")
			}
		} else {
			log.Printf("failed to query row: %v", err)
			return nil, errors.New("internal server error")
		}
	}

	// insert metric into metrics table
	timestamp := time.Unix(int64(tr.Timestamp), 0) // need to convert to insert a value with a type of TIMESTAMP on the db

	_, err = dbConn.Exec(ctx, api.InsertIntoMetricsTable, tr.BitrateMbps, tr.Temperature, tr.DroppedFrames, timestamp, tr.EncoderId)

	if err != nil {
		log.Printf("failed to insert into metrics table: %v", err)
		return nil, errors.New("internal server error")
	}

	return &pb.TelemetryResponse{
		Successful: true,
		Message:    "Successfully stored metric in database.",
	}, nil
}
