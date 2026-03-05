use prost_build::Config;
use std::path::PathBuf;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Configure protobuf compilation
    let mut config = Config::new();
    config.type_attribute(".", "#[derive(serde::Serialize, serde::Deserialize)]");
    
    // Compile protobuf files (if any)
    // tonic_build::compile_protos("proto/broker.proto")?;
    
    // Link the pre-built C++ storage library
    let manifest_dir = std::env::var("CARGO_MANIFEST_DIR")?;
    let storage_build_dir = PathBuf::from(&manifest_dir)
        .parent()
        .unwrap()
        .join("storage/build");
    
    println!("cargo:rustc-link-search=native={}", storage_build_dir.display());
    println!("cargo:rustc-link-lib=dylib=streamforge_storage_ffi");
    
    // Set rpath to find the dylib at runtime
    println!("cargo:rustc-link-arg=-Wl,-rpath,{}", storage_build_dir.display());
    
    // Tell cargo to invalidate the built crate whenever build.rs changes
    println!("cargo:rerun-if-changed=build.rs");
    
    Ok(())
}
