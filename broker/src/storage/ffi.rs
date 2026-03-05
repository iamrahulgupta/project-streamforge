//! FFI bindings to C++ storage layer
//! 
//! This module provides Rust bindings to the C++ high-performance storage backend
//! for efficient message storage and retrieval.

use std::ffi::CString;
use std::os::raw::c_char;

extern "C" {
    /// Initialize the C++ storage engine
    fn cpp_storage_init(data_dir: *const c_char) -> *mut std::ffi::c_void;
    
    /// Write data to storage
    fn cpp_storage_write(
        handle: *mut std::ffi::c_void,
        key: *const c_char,
        value: *const u8,
        value_len: usize,
    ) -> i32;
    
    /// Read data from storage
    fn cpp_storage_read(
        handle: *mut std::ffi::c_void,
        key: *const c_char,
        value: *mut *mut u8,
        value_len: *mut usize,
    ) -> i32;
    
    /// Cleanup the C++ storage engine
    fn cpp_storage_cleanup(handle: *mut std::ffi::c_void);
}

/// Wrapper around C++ storage engine
#[derive(Debug)]
pub struct CppStorage {
    handle: Option<*mut std::ffi::c_void>,
}

impl CppStorage {
    pub fn new() -> Self {
        Self { handle: None }
    }
    
    /// Initialize storage with data directory
    pub fn init(&mut self, data_dir: &str) -> anyhow::Result<()> {
        let c_path = CString::new(data_dir)?;
        
        unsafe {
            let handle = cpp_storage_init(c_path.as_ptr());
            if handle.is_null() {
                return Err(anyhow::anyhow!("Failed to initialize C++ storage"));
            }
            self.handle = Some(handle);
        }
        
        Ok(())
    }
    
    /// Write data to storage
    pub fn write(&self, key: &str, value: &[u8]) -> anyhow::Result<()> {
        let handle = self
            .handle
            .ok_or_else(|| anyhow::anyhow!("Storage not initialized"))?;
        let c_key = CString::new(key)?;
        
        unsafe {
            let result = cpp_storage_write(
                handle,
                c_key.as_ptr(),
                value.as_ptr(),
                value.len(),
            );
            
            if result != 0 {
                return Err(anyhow::anyhow!("Failed to write to storage"));
            }
        }
        
        Ok(())
    }
    
    /// Read data from storage
    pub fn read(&self, key: &str) -> anyhow::Result<Vec<u8>> {
        let handle = self
            .handle
            .ok_or_else(|| anyhow::anyhow!("Storage not initialized"))?;
        let c_key = CString::new(key)?;
        
        let mut value_ptr: *mut u8 = std::ptr::null_mut();
        let mut value_len: usize = 0;
        
        unsafe {
            let result = cpp_storage_read(
                handle,
                c_key.as_ptr(),
                &mut value_ptr,
                &mut value_len,
            );
            
            if result != 0 || value_ptr.is_null() {
                return Err(anyhow::anyhow!("Failed to read from storage"));
            }
            
            Ok(std::slice::from_raw_parts(value_ptr, value_len).to_vec())
        }
    }
}

impl Drop for CppStorage {
    fn drop(&mut self) {
        if let Some(handle) = self.handle {
            unsafe {
                cpp_storage_cleanup(handle);
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_cpp_storage_creation() {
        let storage = CppStorage::new();
        assert!(storage.handle.is_none());
    }
}
