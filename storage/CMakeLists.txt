cmake_minimum_required(VERSION 3.15)
project(streamforge-storage)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Add include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include)
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/ffi)

# Source files
set(STORAGE_SOURCES
    src/log_engine.cpp
    src/segment.cpp
    src/index.cpp
    src/mmap_file.cpp
)

# Create static library
add_library(streamforge_storage STATIC ${STORAGE_SOURCES})

# FFI bridge library
add_library(streamforge_storage_ffi SHARED
    ffi/bridge.cpp
    ${STORAGE_SOURCES}
)

# Set properties
set_target_properties(streamforge_storage PROPERTIES
    PUBLIC_HEADER "include/log_engine.h;include/segment.h"
)

set_target_properties(streamforge_storage_ffi PROPERTIES
    SOVERSION 1
    VERSION 1.0.0
)

# Optional: Add tests
enable_testing()
add_subdirectory(tests EXCLUDE_FROM_ALL)

# Installation
install(TARGETS streamforge_storage streamforge_storage_ffi
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
    PUBLIC_HEADER DESTINATION include/streamforge
)
