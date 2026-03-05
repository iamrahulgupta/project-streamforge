use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Config {
    /// Broker ID in the cluster
    pub broker_id: u32,
    
    /// Server listening address
    pub listen_addr: String,
    
    /// Port to listen on
    pub port: u16,
    
    /// Cluster bootstrap servers
    pub bootstrap_servers: Vec<String>,
    
    /// Storage path
    pub data_dir: String,
    
    /// Replication factor
    pub replication_factor: u16,
    
    /// Min in-sync replicas
    pub min_insync_replicas: u16,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            broker_id: 0,
            listen_addr: "0.0.0.0".to_string(),
            port: 9092,
            bootstrap_servers: vec![],
            data_dir: "./data".to_string(),
            replication_factor: 3,
            min_insync_replicas: 2,
        }
    }
}

impl Config {
    /// Load configuration from environment variables or defaults
    pub fn from_env() -> Self {
        let mut config = Self::default();
        
        if let Ok(broker_id) = std::env::var("BROKER_ID") {
            config.broker_id = broker_id.parse().unwrap_or(0);
        }
        
        if let Ok(listen_addr) = std::env::var("LISTEN_ADDR") {
            config.listen_addr = listen_addr;
        }
        
        if let Ok(port) = std::env::var("PORT") {
            config.port = port.parse().unwrap_or(9092);
        }
        
        config
    }
}
