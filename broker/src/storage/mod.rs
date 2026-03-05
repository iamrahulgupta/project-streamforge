pub mod ffi;

pub use ffi::CppStorage;

/// Storage layer management
#[derive(Debug)]
pub struct StorageManager {
    cpp_storage: CppStorage,
}

impl StorageManager {
    pub fn new() -> Self {
        Self {
            cpp_storage: CppStorage::new(),
        }
    }
}
