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

int main()
{
    std::vector<std::shared_ptr<VideoEncoder>> encoders;

    // load encoders from config
    std::ifstream f("config.json");
    json config = json::parse(f);

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
        std::thread ve_thread([ve]()
                              { ve->run_simulation(); });

        ve_thread.detach();
    }

    // establish TCP connection with server & create client stub
    auto channel = grpc::CreateChannel(config["serverAddress"].get<std::string>(), grpc::InsecureChannelCredentials());
    std::unique_ptr<flux::TelemetryService::Stub> c_stub = flux::TelemetryService::NewStub(channel);

    // send requests for each encoder
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

            grpc::Status status = c_stub->RecordMetrics(&ctx, *reqOpt, &res);

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

        std::this_thread::sleep_for(std::chrono::seconds(2));
    }
    
    return 0;
}