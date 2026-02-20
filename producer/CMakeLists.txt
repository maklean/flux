cmake_minimum_required(VERSION 3.10)
project(flux-producer)

# C++17
set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(Protobuf REQUIRED)
find_package(gRPC REQUIRED)

include(FetchContent)
FetchContent_Declare(
    json 
    GIT_REPOSITORY https://github.com/nlohmann/json.git
)
FetchContent_MakeAvailable(json)

# fix for the gRPC include conflicts
include_directories(${CMAKE_CURRENT_SOURCE_DIR})

add_executable(flux_producer 
    main.cpp
    VideoEncoder.cpp
    proto/telemetry.grpc.pb.cc
    proto/telemetry.pb.cc
)

target_link_libraries(flux_producer 
    gRPC::grpc++ 
    protobuf::libprotobuf
    nlohmann_json::nlohmann_json
)