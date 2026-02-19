# flux

Generate type definitions & server interface
```bash
protoc --go_out=paths=source_relative:./server --go-grpc_out=paths=source_relative:./server proto/telemetry.proto
```