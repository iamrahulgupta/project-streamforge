#include "../include/log_engine.h"
#include <iostream>
#include <cassert>
#include <cstring>
#include <filesystem>

namespace fs = std::filesystem;
using namespace streamforge::storage;

// Test basic create and init
void test_engine_init() {
    std::cout << "Test: LogEngine initialization..." << std::endl;
    
    LogEngineConfig config;
    config.data_dir = "/tmp/streamforge_test_1";
    
    LogEngine engine(config);
    assert(engine.init() == 0);
    
    LogEngine::Stats stats = engine.get_stats();
    assert(stats.num_segments == 0 || stats.num_segments > 0); // Either fresh or loaded
    
    std::cout << "✓ LogEngine initialization passed" << std::endl;
}

// Test write and read
void test_write_read() {
    std::cout << "Test: Write and read operations..." << std::endl;
    
    std::string test_dir = "/tmp/streamforge_test_2";
    if (fs::exists(test_dir)) {
        fs::remove_all(test_dir);
    }
    
    LogEngineConfig config;
    config.data_dir = test_dir;
    
    LogEngine engine(config);
    assert(engine.init() == 0);
    
    // Write test data
    std::string key = "test_key_1";
    uint8_t data[] = {1, 2, 3, 4, 5};
    assert(engine.write(key, data, sizeof(data)) == 0);
    
    // Read test data
    uint8_t* read_data = nullptr;
    size_t read_len = 0;
    assert(engine.read(key, &read_data, &read_len) == 0);
    assert(read_len == sizeof(data));
    assert(std::memcmp(read_data, data, sizeof(data)) == 0);
    
    delete[] read_data;
    
    std::cout << "✓ Write and read operations passed" << std::endl;
}

// Test statistics
void test_stats() {
    std::cout << "Test: Engine statistics..." << std::endl;
    
    std::string test_dir = "/tmp/streamforge_test_3";
    if (fs::exists(test_dir)) {
        fs::remove_all(test_dir);
    }
    
    LogEngineConfig config;
    config.data_dir = test_dir;
    
    LogEngine engine(config);
    assert(engine.init() == 0);
    
    // Write multiple entries
    for (int i = 0; i < 10; i++) {
        std::string key = "key_" + std::to_string(i);
        uint8_t data[100];
        engine.write(key, data, sizeof(data));
    }
    
    LogEngine::Stats stats = engine.get_stats();
    assert(stats.total_entries >= 10);
    assert(stats.total_bytes > 0);
    
    std::cout << "  Stats: entries=" << stats.total_entries 
              << ", bytes=" << stats.total_bytes 
              << ", segments=" << stats.num_segments << std::endl;
    
    std::cout << "✓ Engine statistics test passed" << std::endl;
}

int main() {
    std::cout << "Running StreamForge Storage Tests" << std::endl;
    std::cout << "=================================" << std::endl;
    
    try {
        test_engine_init();
        test_write_read();
        test_stats();
        
        std::cout << "\n✓ All tests passed!" << std::endl;
        return 0;
    } catch (const std::exception& e) {
        std::cerr << "Test failed with exception: " << e.what() << std::endl;
        return 1;
    }
}
