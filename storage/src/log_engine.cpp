#include "log_engine.h"
#include "segment.h"
#include <iostream>
#include <filesystem>
#include <chrono>
#include <cstring>

namespace fs = std::filesystem;

namespace streamforge {
namespace storage {

LogEngine::LogEngine(const LogEngineConfig& config) : config_(config) {}

LogEngine::~LogEngine() {
    if (initialized_) {
        shutdown();
    }
}

int LogEngine::init() {
    try {
        // Create data directory if it doesn't exist
        if (!fs::exists(config_.data_dir)) {
            fs::create_directories(config_.data_dir);
        }
        
        // Load existing segments
        for (const auto& entry : fs::directory_iterator(config_.data_dir)) {
            if (entry.path().extension() == ".log") {
                uint64_t segment_id = std::stoull(entry.path().stem());
                auto segment = std::make_shared<Segment>(
                    segment_id,
                    entry.path().string()
                );
                if (segment->init() == 0) {
                    segments_[entry.path().stem()] = segment;
                    next_segment_id_ = std::max(next_segment_id_, segment_id + 1);
                }
            }
        }
        
        initialized_ = true;
        std::cout << "LogEngine initialized successfully" << std::endl;
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Error initializing LogEngine: " << e.what() << std::endl;
        return -1;
    }
}

int LogEngine::write(const std::string& key, const uint8_t* value, size_t value_len) {
    if (!initialized_) {
        std::cerr << "LogEngine not initialized" << std::endl;
        return -1;
    }
    
    try {
        // Get or create active segment
        std::string active_key = std::to_string(next_segment_id_);
        auto it = segments_.find(active_key);
        
        std::shared_ptr<Segment> segment;
        if (it == segments_.end()) {
            // Create new segment
            std::string path = config_.data_dir + "/" + active_key + ".log";
            segment = std::make_shared<Segment>(next_segment_id_, path);
            if (segment->init() != 0) {
                return -1;
            }
            segments_[active_key] = segment;
        } else {
            segment = it->second;
        }
        
        // Check if segment is full
        if (segment->is_full(config_.segment_size)) {
            next_segment_id_++;
            return write(key, value, value_len);
        }
        
        // Write to segment
        return segment->write(key, value, value_len);
    } catch (const std::exception& e) {
        std::cerr << "Error writing to log: " << e.what() << std::endl;
        return -1;
    }
}

int LogEngine::read(const std::string& key, uint8_t** value, size_t* value_len) {
    if (!initialized_) {
        std::cerr << "LogEngine not initialized" << std::endl;
        return -1;
    }
    
    try {
        // Search all segments for the key
        for (auto& [seg_key, segment] : segments_) {
            if (segment->read(key, value, value_len) == 0) {
                return 0;
            }
        }
        
        return -1; // Key not found
    } catch (const std::exception& e) {
        std::cerr << "Error reading from log: " << e.what() << std::endl;
        return -1;
    }
}

int LogEngine::delete_entry(const std::string& key) {
    if (!initialized_) {
        return -1;
    }
    
    // TODO: Implement tombstone-based deletion
    return 0;
}

int LogEngine::sync() {
    if (!initialized_) {
        return -1;
    }
    
    try {
        for (auto& [key, segment] : segments_) {
            if (segment->sync() != 0) {
                return -1;
            }
        }
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Error syncing log: " << e.what() << std::endl;
        return -1;
    }
}

int LogEngine::compact() {
    if (!initialized_) {
        return -1;
    }
    
    // TODO: Implement segment compaction
    return 0;
}

LogEngine::Stats LogEngine::get_stats() const {
    Stats stats;
    stats.num_segments = segments_.size();
    stats.active_segment_id = next_segment_id_;
    
    for (const auto& [key, segment] : segments_) {
        stats.total_entries += segment->entry_count();
        stats.total_bytes += segment->size();
    }
    
    return stats;
}

void LogEngine::shutdown() {
    try {
        sync();
        segments_.clear();
        initialized_ = false;
        std::cout << "LogEngine shutdown complete" << std::endl;
    } catch (const std::exception& e) {
        std::cerr << "Error shutting down LogEngine: " << e.what() << std::endl;
    }
}

} // namespace storage
} // namespace streamforge
