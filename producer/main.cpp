#include <vector>
#include "include/VideoEncoder.h"
#include <iostream>
#include <memory>
#include <thread>
#include <chrono>
#include <nlohmann/json.hpp>
#include <fstream>
#include <grpcpp/grpcpp.h>
#include "proto/telemetry.grpc.pb.h"

using json = nlohmann::json;

const int telemetry_loop_delay = 5;

// Loads the .json config at 'path' and returns the json instance.
json loadConfig(const std::string& path);

// Creates a vector of shared_ptrs to VideoEncoder from the config.
std::vector<std::shared_ptr<VideoEncoder>> createEncoders(const json& config);

// Runs the telemetry loop that records the metrics of each VideoEncoder every 'telemetry_loop_delay' seconds.
void runTelemetryLoop(const std::vector<std::shared_ptr<VideoEncoder>>& encoders, flux::TelemetryService::Stub& c_stub);

int main()
{
    try {
        json config = loadConfig("config.json");
        auto encoders = createEncoders(config);
        
        // establish TCP connection with server & create client stub
        auto channel = grpc::CreateChannel(config["serverAddress"].get<std::string>(), grpc::InsecureChannelCredentials());
        auto c_stub = flux::TelemetryService::NewStub(channel);

        runTelemetryLoop(encoders, *c_stub);
    } catch (const std::exception& e) {
        std::cerr << "Error: " << e.what() << '\n';
        return EXIT_FAILURE;
    }
    
    return 0;
}

json loadConfig(const std::string& path) {
    std::ifstream f(path);

    if(!f) {
        throw std::runtime_error("Failed to open config file.");
    }

    return json::parse(f);
}

std::vector<std::shared_ptr<VideoEncoder>> createEncoders(const json& config) {
    std::vector<std::shared_ptr<VideoEncoder>> encoders;

    for (const auto &enc : config["encoders"])
    {
        std::shared_ptr<VideoEncoder> ve = std::make_shared<VideoEncoder>(
            enc["id"].get<std::string>(),
            enc["base_bitrate"].get<double>(),
            enc["base_temperature"].get<double>(),
            enc["failure_bitrate_threshold"].get<double>(),
            enc["failure_bitrate_probability"].get<double>());

        encoders.push_back(ve);

        // start simulation for video encoder
        std::thread ve_thread([ve](){ ve->run_simulation(); });

        ve_thread.detach();
    }

    return encoders;
}

void runTelemetryLoop(const std::vector<std::shared_ptr<VideoEncoder>>& encoders, flux::TelemetryService::Stub& c_stub) {
    while (true)
    {
        for (const auto &enc : encoders)
        {
            std::optional<flux::TelemetryRequest> reqOpt = enc->fetch_metrics();

            // Failed to fetch metrics
            if (!reqOpt)
            {
                std::cerr << "Failed to fetch metrics for VideoEncoder(id=" << enc->getId() << ")\n";
                continue;
            }

            grpc::ClientContext ctx;
            flux::TelemetryResponse res;

            grpc::Status status = c_stub.RecordMetrics(&ctx, *reqOpt, &res);

            // Failed to send request to server
            if (!status.ok())
            {
                std::cerr << "ERROR: Failed to send request for VideoEncoder(id=" << enc->getId() << "). " << status.error_message() << '\n';
                continue;
            }

            const std::string& msg(res.message());

            // The request did not return a successful result
            if (!res.successful())
            {
                std::cerr << "ERROR: Failed to complete request for VideoEncoder(id=" << enc->getId() << "): " << msg << '\n';
                continue;
            }

            std::cout << "LOG - VideoEncoder(id=" << enc->getId() << "): " << msg << '\n';
        }

        std::this_thread::sleep_for(std::chrono::seconds(telemetry_loop_delay));
    }
}