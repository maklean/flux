#include <vector>
#include "include/VideoEncoder.h"
#include <iostream>
#include <memory>
#include <thread>
#include <chrono>

int main() {
    std::vector<std::shared_ptr<VideoEncoder>> ves;
    std::vector<std::string> encoder_ids{ "ENC-01", "ENC-02", "ENC-03" };

    for(const auto& id : encoder_ids) {
        std::shared_ptr<VideoEncoder> ve = std::make_shared<VideoEncoder>(id, 5.0, 53.2, 10);

        ves.push_back(ve);

        std::thread simulation_thread([ve]() {
            ve->run_simulation();
        });

        simulation_thread.detach();
    }

    while(true) {
        for(const auto& ve : ves) {
            std::optional<flux::TelemetryRequest> reqOpt( ve->fetch_metrics() );

            if(reqOpt) {
                std::cout << "LOG: Got TelemetryRequest(id = " << reqOpt->encoder_id() << ", bitrate = " << reqOpt->bitrate_mbps() << ", temperature = " << reqOpt->temperature() << ", timestamp = " << reqOpt->timestamp() << ")\n";
            }
        }

        std::this_thread::sleep_for(std::chrono::seconds(1));
    }

    return 0;
}