#include <string>
#include <vector>
#include <unordered_map>
#include <iostream>
#include <fstream>
#include <cstring>
#include <cstdint>

namespace streamforge {
namespace storage {

/// Index file manager for segment lookups
class IndexManager {
public:
    struct IndexEntry {
        std::string key;
        uint64_t offset = 0;
        uint32_t length = 0;
    };
    
    /// Build index from segment data
    static int build_index(const std::string& segment_path,
                          std::vector<IndexEntry>& index) {
        try {
            std::ifstream file(segment_path, std::ios::binary);
            if (!file.is_open()) {
                std::cerr << "Cannot open segment file: " << segment_path << std::endl;
                return -1;
            }
            
            uint64_t offset = 0;
            
            while (!file.eof()) {
                // Read key length
                uint32_t key_len;
                file.read(reinterpret_cast<char*>(&key_len), sizeof(key_len));
                if (file.gcount() == 0) break;
                
                // Read key
                std::string key(key_len, '\0');
                file.read(&key[0], key_len);
                
                // Read value length
                uint32_t value_len;
                file.read(reinterpret_cast<char*>(&value_len), sizeof(value_len));
                
                // Skip value
                file.seekg(value_len, std::ios::cur);
                
                // Skip timestamp and checksum
                file.seekg(12, std::ios::cur); // 8 bytes (timestamp) + 4 bytes (checksum)
                
                // Record entry
                IndexEntry entry;
                entry.key = key;
                entry.offset = offset;
                entry.length = sizeof(key_len) + key_len + sizeof(value_len) + 
                              value_len + 8 + 4;
                index.push_back(entry);
                
                offset += entry.length;
            }
            
            file.close();
            return 0;
        } catch (const std::exception& e) {
            std::cerr << "Error building index: " << e.what() << std::endl;
            return -1;
        }
    }
    
    /// Save index to index file
    static int save_index(const std::string& index_path,
                         const std::vector<IndexEntry>& index) {
        try {
            std::ofstream file(index_path, std::ios::binary);
            if (!file.is_open()) {
                std::cerr << "Cannot open index file: " << index_path << std::endl;
                return -1;
            }
            
            // Write number of entries
            uint32_t num_entries = index.size();
            file.write(reinterpret_cast<const char*>(&num_entries), sizeof(num_entries));
            
            // Write each entry
            for (const auto& entry : index) {
                uint32_t key_len = entry.key.length();
                file.write(reinterpret_cast<const char*>(&key_len), sizeof(key_len));
                file.write(entry.key.c_str(), key_len);
                file.write(reinterpret_cast<const char*>(&entry.offset), sizeof(entry.offset));
                file.write(reinterpret_cast<const char*>(&entry.length), sizeof(entry.length));
            }
            
            file.close();
            return 0;
        } catch (const std::exception& e) {
            std::cerr << "Error saving index: " << e.what() << std::endl;
            return -1;
        }
    }
    
    /// Load index from index file
    static int load_index(const std::string& index_path,
                         std::vector<IndexEntry>& index) {
        try {
            std::ifstream file(index_path, std::ios::binary);
            if (!file.is_open()) {
                std::cerr << "Cannot open index file: " << index_path << std::endl;
                return -1;
            }
            
            // Read number of entries
            uint32_t num_entries;
            file.read(reinterpret_cast<char*>(&num_entries), sizeof(num_entries));
            
            // Read each entry
            for (uint32_t i = 0; i < num_entries; i++) {
                uint32_t key_len;
                file.read(reinterpret_cast<char*>(&key_len), sizeof(key_len));
                
                std::string key(key_len, '\0');
                file.read(&key[0], key_len);
                
                uint64_t offset;
                uint32_t length;
                file.read(reinterpret_cast<char*>(&offset), sizeof(offset));
                file.read(reinterpret_cast<char*>(&length), sizeof(length));
                
                IndexEntry entry;
                entry.key = key;
                entry.offset = offset;
                entry.length = length;
                index.push_back(entry);
            }
            
            file.close();
            return 0;
        } catch (const std::exception& e) {
            std::cerr << "Error loading index: " << e.what() << std::endl;
            return -1;
        }
    }
};

} // namespace storage
} // namespace streamforge
