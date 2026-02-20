#include "include/VideoEncoder.h"

void VideoEncoder::run_simulation()
{
    {
        std::lock_guard<std::mutex> lock(m_mtx);

        if (m_running)
        {
            return;
        }

        m_running = true;
    }

    std::random_device rd;

    // to generate random values
    std::default_random_engine generator(rd());

    // for temperature and bitrate
    std::uniform_real_distribution<double> distribution(-0.1, 0.1);

    while (true)
    {
        {
            // lock mutex before updating
            std::lock_guard<std::mutex> lock(m_mtx);

            m_bitrate += distribution(generator); // should simulate bitrate fluctuation
            m_temperature += (m_bitrate > s_bitrate_temp_inc_threshold) ? 0.05 : -0.02;

            m_temperature += distribution(generator) * 0.1; // small jitter on temperature so it isn't linear

            m_temperature = std::max(20.0, std::min(m_temperature, 100.0)); // clamp

            // if the failure bitrate threshold is passed, simulate dropped frames
            if (m_bitrate > m_failure_bitrate_threshold)
            {
                std::uniform_real_distribution<double> chance(0.0, 1.0);

                // drop frames
                if (chance(generator) <= m_failure_bitrate_probability)
                {
                    std::uniform_int_distribution<int> drop_count(1, 3);
                    m_dropped_frames += drop_count(generator);
                }
            }
        }

        std::this_thread::sleep_for(std::chrono::seconds(s_simulation_thread_sleep_duration));
    }
}

std::optional<flux::TelemetryRequest> VideoEncoder::fetch_metrics()
{
    std::lock_guard<std::mutex> lock(m_mtx);

    if (!m_running)
    {
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