use anyhow::Result;
use std::net::SocketAddr;
use tracing::info;

pub mod cluster;
pub mod config;
pub mod metrics;
pub mod network;
pub mod partition;
pub mod replication;
pub mod storage;

use network::NetworkManager;

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    tracing_subscriber::fmt::init();

    eprintln!("[STARTUP] Starting StreamForge Broker");
    info!("Starting StreamForge Broker");

    // Load configuration from environment variables
    eprintln!("[STARTUP] Loading configuration from environment variables");
    
    let broker_id: u32 = std::env::var("BROKER_ID")
        .unwrap_or_else(|_| "1".to_string())
        .parse()
        .unwrap_or(1);
    eprintln!("[STARTUP] BROKER_ID={}", broker_id);
    
    let listen_addr: String = std::env::var("LISTEN_ADDR")
        .unwrap_or_else(|_| "0.0.0.0".to_string());
    eprintln!("[STARTUP] LISTEN_ADDR={}", listen_addr);
    
    let port: u16 = std::env::var("PORT")
        .unwrap_or_else(|_| "9092".to_string())
        .parse()
        .unwrap_or(9092);
    eprintln!("[STARTUP] PORT={}", port);
    
    let addr_str = format!("{}:{}", listen_addr, port);
    eprintln!("[STARTUP] Parsing address: {}", addr_str);
    
    let addr: SocketAddr = match addr_str.parse() {
        Ok(a) => {
            eprintln!("[STARTUP] Address parsed successfully: {}", a);
            a
        }
        Err(e) => {
            eprintln!("[STARTUP] Failed to parse socket address: {}", e);
            return Err(anyhow::anyhow!("Failed to parse socket address: {}", e));
        }
    };

    eprintln!("[STARTUP] Creating NetworkManager for broker {}", broker_id);
    info!("Initializing broker {} on {}", broker_id, addr);

    // Initialize broker components
    let mut network_mgr = NetworkManager::new(broker_id);
    eprintln!("[STARTUP] NetworkManager created");
    
    // Start network server
    eprintln!("[STARTUP] Starting broker server on {}...", addr);
    if let Err(e) = network_mgr.server.start(addr).await {
        eprintln!("[STARTUP] Failed to start broker server: {}", e);
        return Err(anyhow::anyhow!("Failed to start broker server: {}", e));
    }
    
    eprintln!("[STARTUP] Broker server started successfully");
    info!("Broker {} is running on {}", broker_id, addr);
    
    eprintln!("[STARTUP] Entering main event loop...");
    // Keep the broker running - accept and handle connections
    if let Err(e) = network_mgr.server.run().await {
        eprintln!("[STARTUP] Broker server error: {}", e);
        return Err(anyhow::anyhow!("Broker server error: {}", e));
    }

    eprintln!("[STARTUP] Broker shutting down normally");
    Ok(())
}

