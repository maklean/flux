package api

const (
	q_CreateEncodersTable = `
	CREATE TABLE IF NOT EXISTS encoders(
		id VARCHAR(255) NOT NULL PRIMARY KEY
	);`

	q_CreateMetricsTable = `
	CREATE TABLE IF NOT EXISTS metrics(
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
		bitrate_mbps DOUBLE PRECISION NOT NULL, 
		temperature DOUBLE PRECISION NOT NULL, 
		dropped_frames INT NOT NULL, 
		timestamp TIMESTAMP NOT NULL, 
		encoder_id VARCHAR(255) NOT NULL REFERENCES encoders(id)
	);`

	idx_EncoderId_MetricsTable = "CREATE INDEX IF NOT EXISTS idx_metrics_encoder_id ON metrics(encoder_id);"
)
