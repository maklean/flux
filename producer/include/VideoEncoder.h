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
    uint32_t m_dropped_frames;
    std::mutex m_mtx;
    bool m_running{ false };

public:
    VideoEncoder(const std::string &id, double bitrate, double temperature, uint32_t dropped_frames)
        : m_id(id), m_bitrate(bitrate), m_temperature(temperature), m_dropped_frames(dropped_frames)
    {
    }

    void run_simulation() {
        {
            std::lock_guard<std::mutex> lock(m_mtx);

            if(m_running) {
                return;
            }

            m_running = true;
        }

        std::random_device rd;

        // to generate random values
        std::default_random_engine generator(rd());
        std::uniform_real_distribution<double> distribution(-0.1, 0.1);

        while(true) {
            {
                // lock mutex before updating
                std::lock_guard<std::mutex> lock(m_mtx);

                m_bitrate += distribution(generator); // should simulate bitrate fluctuation
                m_temperature += (m_bitrate > s_bitrate_temp_inc_threshold) ? 0.05 : -0.02;

                m_temperature += distribution(generator) * 0.1; // small jitter on temperature so it isn't linear

                m_temperature = std::max(20.0, std::min(m_temperature, 100.0)); // clamp
            }

            std::this_thread::sleep_for(std::chrono::seconds(s_simulation_thread_sleep_duration));
        }
    }

    std::optional<flux::TelemetryRequest> fetch_metrics() {
        std::lock_guard<std::mutex> lock(m_mtx);

        if(!m_running) {
            return std::nullopt;
        }

        // create new telemetry request
        flux::TelemetryRequest req;

        req.set_encoder_id(m_id);
        req.set_bitrate_mbps(m_bitrate);
        req.set_temperature(m_temperature);
        req.set_dropped_frames(m_dropped_frames);
        req.set_timestamp(std::chrono::seconds(std::time(NULL)).count()); // get current unix timestamp

        return req;
    }
};