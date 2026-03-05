pub mod heartbeat;
pub mod membership;
pub mod metadata;

pub use heartbeat::HeartbeatManager;
pub use membership::ClusterMembership;
pub use metadata::MetadataManager;

/// Cluster state and management
#[derive(Debug)]
pub struct ClusterManager {
    membership: ClusterMembership,
    heartbeat: HeartbeatManager,
    metadata: MetadataManager,
}

impl ClusterManager {
    pub fn new(broker_id: u32) -> Self {
        Self {
            membership: ClusterMembership::new(broker_id),
            heartbeat: HeartbeatManager::new(),
            metadata: MetadataManager::new(),
        }
    }
}
