use std::net::SocketAddr;
use tokio::net::TcpListener;
use tracing::{debug, error, info};

/// Broker network server
#[derive(Debug)]
pub struct BrokerServer {
    broker_id: u32,
    listen_addr: Option<SocketAddr>,
    listener: Option<TcpListener>,
}

impl BrokerServer {
    pub fn new(broker_id: u32) -> Self {
        Self {
            broker_id,
            listen_addr: None,
            listener: None,
        }
    }
    
    pub fn broker_id(&self) -> u32 {
        self.broker_id
    }
    
    /// Start the broker server
    pub async fn start(&mut self, addr: SocketAddr) -> anyhow::Result<()> {
        info!(
            "Starting broker {} server on {}",
            self.broker_id, addr
        );
        
        // Create TCP listener
        let listener = TcpListener::bind(addr).await?;
        info!("Broker {} listening on {}", self.broker_id, addr);
        
        self.listen_addr = Some(addr);
        self.listener = Some(listener);
        
        debug!("Broker {} server started successfully", self.broker_id);
        Ok(())
    }
    
    /// Accept connections and handle them
    pub async fn run(&self) -> anyhow::Result<()> {
        if let Some(listener) = &self.listener {
            loop {
                match listener.accept().await {
                    Ok((socket, addr)) => {
                        debug!(
                            "Broker {} accepted connection from {}",
                            self.broker_id, addr
                        );
                        // TODO: Spawn task to handle client connection
                        // - Parse broker protocol messages
                        // - Process produce/consume requests
                        // - Send responses
                        drop(socket); // For now, close the connection
                    }
                    Err(e) => {
                        error!("Broker {} accept error: {}", self.broker_id, e);
                    }
                }
            }
        } else {
            Err(anyhow::anyhow!(
                "Broker {} server not initialized",
                self.broker_id
            ))
        }
    }
    
    /// Shutdown the broker server
    pub async fn shutdown(&mut self) -> anyhow::Result<()> {
        info!("Shutting down broker {} server", self.broker_id);
        self.listen_addr = None;
        self.listener = None;
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    #[test]
    fn test_server_creation() {
        let server = BrokerServer::new(1);
        assert_eq!(server.broker_id(), 1);
    }
}
