#include "bridge.h"
#include "log_engine.h"
#include <cstring>
#include <iostream>
#include <new>

using namespace streamforge::storage;

extern "C" {

LogEngineHandle cpp_storage_init(const char* data_dir) {
    if (!data_dir) {
        std::cerr << "Invalid data_dir pointer" << std::endl;
        return nullptr;
    }
    
    try {
        LogEngineConfig config;
        config.data_dir = std::string(data_dir);
        config.segment_size = 1024 * 1024 * 1024; // 1GB
        
        LogEngine* engine = new LogEngine(config);
        if (engine->init() != 0) {
            delete engine;
            return nullptr;
        }
        
        return static_cast<LogEngineHandle>(engine);
    } catch (const std::exception& e) {
        std::cerr << "Error initializing storage: " << e.what() << std::endl;
        return nullptr;
    }
}

int32_t cpp_storage_write(
    LogEngineHandle handle,
    const char* key,
    const uint8_t* value,
    size_t value_len
) {
    if (!handle || !key || !value) {
        std::cerr << "Invalid arguments to cpp_storage_write" << std::endl;
        return -1;
    }
    
    try {
        LogEngine* engine = static_cast<LogEngine*>(handle);
        return engine->write(std::string(key), value, value_len);
    } catch (const std::exception& e) {
        std::cerr << "Error writing to storage: " << e.what() << std::endl;
        return -1;
    }
}

int32_t cpp_storage_read(
    LogEngineHandle handle,
    const char* key,
    uint8_t** value,
    size_t* value_len
) {
    if (!handle || !key || !value || !value_len) {
        std::cerr << "Invalid arguments to cpp_storage_read" << std::endl;
        return -1;
    }
    
    try {
        LogEngine* engine = static_cast<LogEngine*>(handle);
        int result = engine->read(std::string(key), value, value_len);
        
        if (result == 0 && *value) {
            // Rust will manage the memory with the pointer we return
            return 0;
        }
        
        return result;
    } catch (const std::exception& e) {
        std::cerr << "Error reading from storage: " << e.what() << std::endl;
        return -1;
    }
}

int32_t cpp_storage_sync(LogEngineHandle handle) {
    if (!handle) {
        std::cerr << "Invalid handle for cpp_storage_sync" << std::endl;
        return -1;
    }
    
    try {
        LogEngine* engine = static_cast<LogEngine*>(handle);
        return engine->sync();
    } catch (const std::exception& e) {
        std::cerr << "Error syncing storage: " << e.what() << std::endl;
        return -1;
    }
}

void cpp_storage_cleanup(LogEngineHandle handle) {
    if (!handle) {
        return;
    }
    
    try {
        LogEngine* engine = static_cast<LogEngine*>(handle);
        engine->shutdown();
        delete engine;
    } catch (const std::exception& e) {
        std::cerr << "Error cleaning up storage: " << e.what() << std::endl;
    }
}

} // extern "C"
