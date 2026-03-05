#include "segment.h"
#include <sys/mman.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <cstring>
#include <iostream>
#include <filesystem>

namespace fs = std::filesystem;

namespace streamforge {
namespace storage {

/// Memory-mapped file implementation
class MmapFile {
public:
    MmapFile(const std::string& path, size_t size)
        : path_(path), size_(size), buffer_(nullptr), fd_(-1) {}
    
    ~MmapFile() {
        close();
    }
    
    int init() {
        try {
            // Open or create file
            fd_ = open(path_.c_str(), O_CREAT | O_RDWR, 0644);
            if (fd_ < 0) {
                perror("open");
                return -1;
            }
            
            // Extend file to desired size
            if (lseek(fd_, size_ - 1, SEEK_SET) < 0) {
                perror("lseek");
                close();
                return -1;
            }
            
            // Write a byte at the end
            if (write(fd_, "", 1) < 0) {
                perror("write");
                close();
                return -1;
            }
            
            // Memory map the file
            buffer_ = static_cast<uint8_t*>(
                mmap(nullptr, size_, PROT_READ | PROT_WRITE, MAP_SHARED, fd_, 0)
            );
            
            if (buffer_ == MAP_FAILED) {
                perror("mmap");
                close();
                return -1;
            }
            
            std::cout << "MmapFile initialized: " << path_ << " (" << size_ << " bytes)" << std::endl;
            return 0;
        } catch (const std::exception& e) {
            std::cerr << "Error initializing MmapFile: " << e.what() << std::endl;
            return -1;
        }
    }
    
    uint8_t* get_buffer() {
        return buffer_;
    }
    
    size_t get_size() const {
        return size_;
    }
    
    int sync() {
        if (!buffer_) {
            return -1;
        }
        
        if (msync(buffer_, size_, MS_SYNC) < 0) {
            perror("msync");
            return -1;
        }
        
        return 0;
    }
    
    void close() {
        if (buffer_ && buffer_ != MAP_FAILED) {
            munmap(buffer_, size_);
            buffer_ = nullptr;
        }
        
        if (fd_ >= 0) {
            ::close(fd_);
            fd_ = -1;
        }
    }

private:
    std::string path_;
    size_t size_;
    uint8_t* buffer_;
    int fd_;
};

} // namespace storage
} // namespace streamforge
