#pragma once

#include <string>
#include <cstdint>
#include <mutex>
#include <chrono>
#include <random>
#include <thread>
#include <algorithm>
#include "../proto/telemetry.pb.h"
#include <optional>

class VideoEncoder
{
private:
    constexpr static double s_bitrate_temp_inc_threshold = 6.7;
    constexpr static int64_t s_simulation_thread_sleep_duration = 1;

    std::string m_id;
    double m_bitrate;
    double m_temperature;
    double m_failure_bitrate_threshold;
    double m_failure_bitrate_probability;

    uint32_t m_dropped_frames{};
    std::mutex m_mtx;
    bool m_running{false};

public:
    VideoEncoder(
        const std::string &id,
        double bitrate,
        double temperature,
        double failure_bitrate_threshold,
        double failure_bitrate_probability) : m_id(id), m_bitrate(bitrate), m_temperature(temperature),
                                              m_failure_bitrate_threshold(failure_bitrate_threshold),
                                              m_failure_bitrate_probability(failure_bitrate_probability)
    {
    }

    const std::string& getId() const { return m_id; }

    // run_simulation() starts the VideoEncoder simulation
    void run_simulation();

    // fetch_metrics() packages the current VideoEncoder state into a `flux::TelemetryRequest`
    std::optional<flux::TelemetryRequest> fetch_metrics();
};