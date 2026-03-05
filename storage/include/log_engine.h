#pragma once

#include <string>
#include <vector>
#include <memory>
#include <unordered_map>
#include <cstdint>

namespace streamforge {
namespace storage {

/// Configuration for the log engine
struct LogEngineConfig {
    /// Path to data directory
    std::string data_dir;
    
    /// Size of each segment file in bytes (default: 1GB)
    uint64_t segment_size = 1024 * 1024 * 1024;
    
    /// Enable compression (default: false)
    bool compression_enabled = false;
    
    /// Sync interval in milliseconds (default: 5000ms)
    uint64_t sync_interval_ms = 5000;
};

/// Main log storage engine
class LogEngine {
public:
    /// Create a new log engine
    LogEngine(const LogEngineConfig& config);
    
    /// Destructor
    ~LogEngine();
    
    /// Initialize the log engine
    /// Returns 0 on success, non-zero on error
    int init();
    
    /// Write data to the log
    /// @param key - Unique key for the entry
    /// @param value - Data to write
    /// @param value_len - Length of data
    /// Returns 0 on success, non-zero on error
    int write(const std::string& key, const uint8_t* value, size_t value_len);
    
    /// Read data from the log
    /// @param key - Key to read
    /// @param value - Output buffer for read data
    /// @param value_len - Output length of read data
    /// Returns 0 on success, non-zero on error
    int read(const std::string& key, uint8_t** value, size_t* value_len);
    
    /// Delete an entry from the log
    int delete_entry(const std::string& key);
    
    /// Sync all pending writes to disk
    int sync();
    
    /// Compact the log (merge segments)
    int compact();
    
    /// Get engine statistics
    struct Stats {
        uint64_t total_entries = 0;
        uint64_t total_bytes = 0;
        uint64_t num_segments = 0;
        uint64_t active_segment_id = 0;
    };
    
    Stats get_stats() const;
    
    /// Shutdown the engine
    void shutdown();

private:
    LogEngineConfig config_;
    std::unordered_map<std::string, std::shared_ptr<class Segment>> segments_;
    uint64_t next_segment_id_ = 0;
    bool initialized_ = false;
};

} // namespace storage
} // namespace streamforge
