# Tests for storage layer
# 
# To build tests:
#   mkdir build && cd build
#   cmake ..
#   cmake --build . --target all
#   cd tests && ctest

enable_testing()

# Test executable
add_executable(storage_tests
    test_log_engine.cpp
)

target_link_libraries(storage_tests
    streamforge_storage
)

# Add test
add_test(
    NAME StorageTests
    COMMAND storage_tests
)
