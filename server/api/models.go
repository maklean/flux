package api

import "time"

type Metric struct {
	Id            int
	BitrateMbps   float64
	Temperature   float64
	DroppedFrames uint32
	Timestamp     time.Time
	EncoderId     string
}
