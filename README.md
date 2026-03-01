# flux

A telemetry system for monitoring video encoder instances in real-time. A multi-threaded C++ producer simulates video encoders and streams metrics via gRPC to a Go server, which persists data to a PostgreSQL database and exposes it through a REST API.

## Prerequisites

- CMake ≥ 3.10
- gRPC and Protobuf (C++ and Go)
- Go ≥ 1.21
- PostgreSQL

## Setup

### 1. Environment
- Rename `.env.example` to `.env` and set environment variables:
```bash
export API_SERVER_PORT=8080
export GRPC_SERVER_PORT=9000
export PSQL_DB_URL="postgres://username:password@localhost:5432/database_name"
```

### 2. Generate Protobuf bindings
Generate the server interface and client stub code for producer and server:

**Go:**
```bash
protoc --go_out=paths=source_relative:./server --go-grpc_out=paths=source_relative:./server proto/telemetry.proto
```

**C++:**
```bash
protoc --grpc_out=./producer --cpp_out=./producer --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` proto/telemetry.proto
```

### 3. Build and run the Go server
```bash
cd server
go run main.go
```

### 4. Build and run the C++ producer
```bash
cd producer
cmake CMakeLists.txt
make
./flux_producer
```

## API

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/metrics` | Returns all stored telemetry metrics |
| GET | `/api/metrics/:encoderId` | Returns all stored telemetry metrics for a specific encoder |
| GET | `/api/encoders` | Returns all stored encoders |

## Todo

- [X] Add `GET /api/metrics/:encoder_id` to query metrics by encoder
- [X] Add `GET /api/encoders` to list all registered encoders
- [ ] Switch from a single `pgx.Conn` to a connection pool (`pgxpool`) for concurrent access
- [ ] Add a dashboard to visualize live encoder metrics (might have to switch to websockets for this)
