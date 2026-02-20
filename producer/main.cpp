#include <vector>
#include "include/VideoEncoder.h"
#include <iostream>
#include <memory>
#include <thread>
#include <chrono>
#include <nlohmann/json.hpp>
#include <fstream>

using json = nlohmann::json;

int main() {
    std::vector<std::shared_ptr<VideoEncoder>> encoders;

    // load encoders from config
    std::ifstream f("config.json");
    json config = json::parse(f);

    for(const auto& enc : config["encoders"]) {
        std::shared_ptr<VideoEncoder> ve = std::make_shared<VideoEncoder>(
            enc["id"].get<std::string>(), 
            enc["base_bitrate"].get<double>(),
            enc["base_temperature"].get<double>(),
            enc["failure_bitrate_threshold"].get<double>(),
            enc["failure_bitrate_probability"].get<double>()
        );

        encoders.push_back(ve);

        // start simulation for video encoder
        std::thread ve_thread([ve]() {
            ve->run_simulation();
        });

        ve_thread.detach();
    }

    while(true) {
        for(const auto& ve : encoders) {
            std::optional<flux::TelemetryRequest> reqOpt( ve->fetch_metrics() );

            if(reqOpt) {
                std::cout << "LOG: Got TelemetryRequest(id = " << reqOpt->encoder_id() << ", bitrate = " << reqOpt->bitrate_mbps() << ", temperature = " << reqOpt->temperature() << ", timestamp = " << reqOpt->timestamp() << ")\n";
            }
        }

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    return 0;
}