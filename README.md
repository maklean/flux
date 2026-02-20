# flux

[GO] Generate Go protobuf types and gRPC bindings
```bash
protoc --go_out=paths=source_relative:./server --go-grpc_out=paths=source_relative:./server proto/telemetry.proto
```

[C++] Generate C++ protobuf types and gRPC bindings
```bash
protoc --grpc_out=./producer --cpp_out=./producer --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` proto/telemetry.proto
```