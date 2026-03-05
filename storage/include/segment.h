#pragma once

#include <string>
#include <vector>
#include <memory>
#include <cstdint>

namespace streamforge {
namespace storage {

/// Forward declaration
class MmapFile;

/// Entry in a segment
struct SegmentEntry {
    std::string key;
    std::vector<uint8_t> value;
    uint64_t offset = 0;    /// Offset in segment file
    uint64_t timestamp = 0; /// Creation timestamp
};

/// Index entry for fast lookups
struct IndexEntry {
    std::string key;
    uint64_t offset = 0;  /// File offset
    uint64_t length = 0;  /// Entry length
    uint32_t checksum = 0; /// CRC32 checksum
};

/// A segment file for log storage
class Segment {
public:
    /// Create a new segment
    /// @param segment_id - Unique segment identifier
    /// @param path - Path to segment file
    Segment(uint64_t segment_id, const std::string& path);
    
    ~Segment();
    
    /// Initialize the segment
    int init();
    
    /// Get segment ID
    uint64_t id() const { return segment_id_; }
    
    /// Write an entry to the segment
    int write(const std::string& key, const uint8_t* data, size_t data_len);
    
    /// Read an entry from the segment
    int read(const std::string& key, uint8_t** data, size_t* data_len);
    
    /// Get the current size of the segment
    uint64_t size() const;
    
    /// Check if segment is full
    bool is_full(uint64_t max_size) const;
    
    /// Get the number of entries
    uint64_t entry_count() const;
    
    /// Sync segment to disk
    int sync();
    
    /// Build index for faster lookups
    int build_index();
    
    /// Close the segment
    void close();

private:
    uint64_t segment_id_;
    std::string path_;
    std::shared_ptr<MmapFile> mmap_file_;
    std::vector<IndexEntry> index_;
    uint64_t current_offset_ = 0;
    bool initialized_ = false;
};

} // namespace storage
} // namespace streamforge
