pub mod client;
pub mod protocol;
pub mod server;

pub use client::BrokerClient;
pub use protocol::Protocol;
pub use server::BrokerServer;

/// Network communication layer
#[derive(Debug)]
pub struct NetworkManager {
    pub server: BrokerServer,
    pub client: BrokerClient,
}

impl NetworkManager {
    pub fn new(broker_id: u32) -> Self {
        Self {
            server: BrokerServer::new(broker_id),
            client: BrokerClient::new(broker_id),
        }
    }
}
