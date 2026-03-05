#pragma once

#include <cstdint>
#include <cstddef>

#ifdef __cplusplus
extern "C" {
#endif

/// Opaque handle to LogEngine
typedef void* LogEngineHandle;

/// Initialize storage with data directory
/// @param data_dir - Path to data directory
/// @return Handle to LogEngine, NULL on error
LogEngineHandle cpp_storage_init(const char* data_dir);

/// Write data to storage
/// @param handle - LogEngine handle
/// @param key - Key string (null-terminated)
/// @param value - Data to write
/// @param value_len - Length of data
/// @return 0 on success, non-zero on error
int32_t cpp_storage_write(
    LogEngineHandle handle,
    const char* key,
    const uint8_t* value,
    size_t value_len
);

/// Read data from storage
/// @param handle - LogEngine handle
/// @param key - Key string (null-terminated)
/// @param value - Output pointer to data (allocated by function)
/// @param value_len - Output length of data
/// @return 0 on success, non-zero on error
int32_t cpp_storage_read(
    LogEngineHandle handle,
    const char* key,
    uint8_t** value,
    size_t* value_len
);

/// Sync storage to disk
/// @param handle - LogEngine handle
/// @return 0 on success, non-zero on error
int32_t cpp_storage_sync(LogEngineHandle handle);

/// Cleanup and free storage handle
/// @param handle - LogEngine handle
void cpp_storage_cleanup(LogEngineHandle handle);

#ifdef __cplusplus
}
#endif
