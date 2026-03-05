pub mod follower;
pub mod leader;
pub mod log_sync;

pub use follower::FollowerReplicator;
pub use leader::LeaderReplicator;
pub use log_sync::LogSync;

/// Replication state and management
#[derive(Debug)]
pub struct ReplicationManager {
    leader: LeaderReplicator,
    follower: FollowerReplicator,
    log_sync: LogSync,
}

impl ReplicationManager {
    pub fn new(broker_id: u32) -> Self {
        Self {
            leader: LeaderReplicator::new(broker_id),
            follower: FollowerReplicator::new(broker_id),
            log_sync: LogSync::new(),
        }
    }
}
