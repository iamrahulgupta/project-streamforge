#include "segment.h"
#include <iostream>
#include <algorithm>
#include <cstring>
#include <chrono>
#include <zlib.h>

namespace streamforge {
namespace storage {

// Forward declare MmapFile
class MmapFile {
public:
    MmapFile(const std::string& path, size_t size);
    ~MmapFile();
    
    int init();
    uint8_t* get_buffer() { return buffer_; }
    size_t get_size() const { return size_; }
    int sync();
    void close();
    
private:
    std::string path_;
    size_t size_;
    uint8_t* buffer_ = nullptr;
    int fd_ = -1;
};

// MmapFile implementation
MmapFile::MmapFile(const std::string& path, size_t size)
    : path_(path), size_(size) {}

MmapFile::~MmapFile() {
    if (fd_ != -1) {
        close();
    }
}

int MmapFile::init() {
    // TODO: Implement mmap initialization
    // For now, allocate buffer on heap
    try {
        buffer_ = new uint8_t[size_];
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Failed to allocate buffer: " << e.what() << std::endl;
        return -1;
    }
}

int MmapFile::sync() {
    // TODO: Implement actual mmap sync
    return 0;
}

void MmapFile::close() {
    if (buffer_) {
        delete[] buffer_;
        buffer_ = nullptr;
    }
    if (fd_ != -1) {
        // TODO: Close actual file descriptor
        fd_ = -1;
    }
}

Segment::Segment(uint64_t segment_id, const std::string& path)
    : segment_id_(segment_id), path_(path) {}

Segment::~Segment() {
    if (initialized_) {
        close();
    }
}

int Segment::init() {
    try {
        // Create mmap file (2GB by default for reserves)
        mmap_file_ = std::make_shared<MmapFile>(path_, 2ULL * 1024 * 1024 * 1024);
        if (mmap_file_->init() != 0) {
            std::cerr << "Failed to initialize mmap file: " << path_ << std::endl;
            return -1;
        }
        
        initialized_ = true;
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Error initializing segment: " << e.what() << std::endl;
        return -1;
    }
}

int Segment::write(const std::string& key, const uint8_t* data, size_t data_len) {
    if (!initialized_) {
        return -1;
    }
    
    try {
        // Create entry
        SegmentEntry entry;
        entry.key = key;
        entry.value.assign(data, data + data_len);
        entry.offset = current_offset_;
        entry.timestamp = std::chrono::system_clock::now().time_since_epoch().count();
        
        // Write entry format:
        // [key_len:4][key][value_len:4][value][timestamp:8][checksum:4]
        
        uint8_t* buffer = mmap_file_->get_buffer() + current_offset_;
        size_t offset = 0;
        
        // Write key length and key
        uint32_t key_len = key.length();
        std::memcpy(buffer + offset, &key_len, sizeof(key_len));
        offset += sizeof(key_len);
        std::memcpy(buffer + offset, key.c_str(), key.length());
        offset += key.length();
        
        // Write value length and value
        uint32_t value_len = data_len;
        std::memcpy(buffer + offset, &value_len, sizeof(value_len));
        offset += sizeof(value_len);
        std::memcpy(buffer + offset, data, data_len);
        offset += data_len;
        
        // Write timestamp
        uint64_t ts = entry.timestamp;
        std::memcpy(buffer + offset, &ts, sizeof(ts));
        offset += sizeof(ts);
        
        // Calculate and write checksum
        uint32_t checksum = 0; // TODO: Implement real CRC32
        std::memcpy(buffer + offset, &checksum, sizeof(checksum));
        offset += sizeof(checksum);
        
        current_offset_ += offset;
        
        // Add to index
        IndexEntry idx;
        idx.key = key;
        idx.offset = entry.offset;
        idx.length = offset;
        idx.checksum = checksum;
        index_.push_back(idx);
        
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Error writing to segment: " << e.what() << std::endl;
        return -1;
    }
}

int Segment::read(const std::string& key, uint8_t** data, size_t* data_len) {
    if (!initialized_) {
        return -1;
    }
    
    try {
        // Search index for key
        auto it = std::find_if(index_.begin(), index_.end(),
            [&key](const IndexEntry& e) { return e.key == key; });
        
        if (it == index_.end()) {
            return -1; // Key not found
        }
        
        // Read from mmap buffer
        uint8_t* buffer = mmap_file_->get_buffer() + it->offset;
        size_t offset = 0;
        
        // Skip key
        uint32_t key_len = *(uint32_t*)(buffer + offset);
        offset += sizeof(key_len) + key_len;
        
        // Read value
        uint32_t value_len = *(uint32_t*)(buffer + offset);
        offset += sizeof(value_len);
        
        // Allocate output buffer
        *data = new uint8_t[value_len];
        std::memcpy(*data, buffer + offset, value_len);
        *data_len = value_len;
        
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Error reading from segment: " << e.what() << std::endl;
        return -1;
    }
}

uint64_t Segment::size() const {
    return current_offset_;
}

bool Segment::is_full(uint64_t max_size) const {
    return current_offset_ >= max_size;
}

uint64_t Segment::entry_count() const {
    return index_.size();
}

int Segment::sync() {
    if (!initialized_) {
        return -1;
    }
    return mmap_file_->sync();
}

int Segment::build_index() {
    if (!initialized_) {
        return -1;
    }
    
    // TODO: Rebuild index from file
    return 0;
}

void Segment::close() {
    if (mmap_file_) {
        mmap_file_->close();
    }
    initialized_ = false;
}

} // namespace storage
} // namespace streamforge
